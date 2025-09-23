package timing

import (
	"context"
	"errors"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chriserin/seq/internal/beats"
	"github.com/chriserin/seq/internal/playstate"
	"github.com/chriserin/seq/internal/seqmidi"
	midi "gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)

func (t *Timing) BeatInterval() time.Duration {
	tickInterval := t.TickInterval()
	adjuster := time.Since(t.playTime) - t.trackTime
	t.trackTime = t.trackTime + tickInterval
	next := tickInterval - adjuster
	return next
}

func (t Timing) TickInterval() time.Duration {
	return time.Minute / time.Duration(t.tempo*t.subdivisions)
}

func (t Timing) PulseInterval() time.Duration {
	return time.Minute / time.Duration(t.tempo*24)
}

type Timing struct {
	playTime     time.Time
	trackTime    time.Duration
	tempo        int
	subdivisions int
	started      bool
	pulseCount   int
	pulseLimit   int
	preRollBeats uint8
	transmitting bool
	beatsLooper  beats.BeatsLooper
	ctx          context.Context
}

type MidiLoopMode uint8

const (
	MlmStandAlone MidiLoopMode = iota
	MlmTransmitter
	MlmReceiver
)

var timingChannel chan TimingMsg

func init() {
	timingChannel = make(chan TimingMsg)
}

func GetTimingChannel() chan TimingMsg {
	return timingChannel
}

func Loop(mode MidiLoopMode, lockReceiverChannel, unlockReceiverChannel chan bool, ctx context.Context, beatsLooper beats.BeatsLooper, sendFn func(tea.Msg)) error {
	timing := Timing{beatsLooper: beatsLooper, ctx: ctx}
	switch mode {
	case MlmStandAlone:
		timing.StandAloneLoop(sendFn)
	case MlmTransmitter:
		err := timing.TransmitterLoop(sendFn)
		if err != nil {
			return fault.Wrap(err, fmsg.With("cannot start transmitter loop"))
		}
	case MlmReceiver:
		err := timing.ReceiverLoop(lockReceiverChannel, unlockReceiverChannel, sendFn)
		timing.StandAloneLoop(sendFn)
		if err != nil {
			// NOTE: In case the receiver loop was not setup correctly, swallow the lock/unlock messages
			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					case <-lockReceiverChannel:
					case <-unlockReceiverChannel:
					}
				}
			}()
			return fault.Wrap(err, fmsg.With("cannot start receiver loop"))
		}
	}
	return nil
}

type Transmitter struct {
	out drivers.Out
}

func (tmtr Transmitter) Start(loopMode playstate.LoopMode) error {
	message := midi.SPP(uint16(loopMode))
	err := tmtr.out.Send(message)
	if err != nil {
		return fault.Wrap(err, fmsg.With("cannot send midi spp pre-start"))
	}
	message = midi.Start()
	err = tmtr.out.Send(message)
	if err != nil {
		return fault.Wrap(err, fmsg.With("cannot send midi start"))
	}
	return nil
}

func (tmtr Transmitter) Stop() error {
	err := tmtr.out.Send(midi.Stop())
	if err != nil {
		return fault.Wrap(err, fmsg.With("cannot send midi stop"))
	}
	return nil
}

func (tmtr Transmitter) Pulse() error {
	err := tmtr.out.Send(midi.TimingClock())
	if err != nil {
		return fault.Wrap(err, fmsg.With("cannot send midi clock"))
	}
	return nil
}

func (tmtr Transmitter) ActiveSense() error {
	err := tmtr.out.Send(midi.Activesense())
	if err != nil {
		return fault.Wrap(err, fmsg.With("cannot send midi active sense"))
	}
	return nil
}

func (t *Timing) TransmitterLoop(sendFn func(tea.Msg)) error {
	var beatChannel = t.beatsLooper.BeatChannel
	out, err := seqmidi.TransmitterOut()
	transmitter := Transmitter{out}
	err = transmitter.ActiveSense()
	if err != nil {
		return fault.Wrap(err)
	}

	tickChannel := make(chan Timing)
	activeSenseChannel := make(chan bool)
	var command TimingMsg

	var tickTimer *time.Timer
	pulse := func(adjustedInterval time.Duration) {
		tickTimer = time.AfterFunc(adjustedInterval, func() {
			tickChannel <- Timing{subdivisions: 24}
		})
	}

	activesense := func() {
		time.AfterFunc(300*time.Millisecond, func() {
			activeSenseChannel <- true
		})
	}

	go func() {
		for {
			select {
			case <-t.ctx.Done():
				return
			case command = <-timingChannel:
				switch command := command.(type) {
				case StartMsg:
					t.started = true
					t.playTime = time.Now()
					t.tempo = command.Tempo
					t.subdivisions = command.Subdivisions
					t.trackTime = time.Duration(0)
					t.pulseCount = 0
					t.pulseLimit = 0
					t.preRollBeats = command.Prerollbeats
					t.transmitting = command.Transmitting
					pulse(0)
					err := transmitter.Start(command.LoopMode)
					if err != nil {
						sendFn(ErrorMsg{err})
					}
				case StopMsg:
					t.started = false
					tickTimer.Stop()
					err := transmitter.Stop()
					if err != nil {
						sendFn(ErrorMsg{err})
					}
					// m.playing should be false now.
				case AnticipatoryStopMsg:
					//NOTE: A receiver must not receive a Pulse message and a Stop message in immediate succession.
					//This will result in a race condition on the receiver end.  Instead, we anticipate stopping
					//and set a limit on the pulses that will be accumulated, preventing the final pulse.
					if t.pulseLimit == 0 {
						t.pulseLimit = t.pulseCount + ((24 / t.subdivisions) - 1)
					}
				case TempoMsg:
					t.tempo = command.Tempo
					t.subdivisions = command.Subdivisions
				}
			case pulseTiming := <-tickChannel:
				if t.started {
					if t.preRollBeats == 0 {
						if t.pulseLimit == 0 || t.pulseCount < t.pulseLimit {
							if t.transmitting {
								err := transmitter.Pulse()
								if err != nil {
									wrappedErr := fault.Wrap(err)
									sendFn(ErrorMsg{wrappedErr})
								}
							}
						}
						if t.pulseCount%(pulseTiming.subdivisions/t.subdivisions) == 0 {
							beatChannel <- beats.BeatMsg{Interval: t.TickInterval()}
						}
					} else {
						if t.pulseCount%(pulseTiming.subdivisions/t.subdivisions) == 0 {
							t.preRollBeats--
							if t.preRollBeats == 0 {
								t.pulseCount = -1
							}
						}
					}
					pulseInterval := t.PulseInterval()

					adjuster := time.Since(t.playTime) - t.trackTime
					t.trackTime = t.trackTime + pulseInterval
					next := pulseInterval - adjuster
					pulse(next)
					t.pulseCount++
				}
			case <-activeSenseChannel:
				if !t.started {
					err := transmitter.ActiveSense()
					if err != nil {
						wrappedErr := fault.Wrap(err, fmsg.With("activesense interrupted"))
						sendFn(ErrorMsg{wrappedErr})
					}
				}
				activesense()
			}
		}
	}()
	// Start active sense loop
	activesense()
	return nil
}

type ListenFn func(msg []byte, milliseconds int32)

func (t *Timing) ReceiverLoop(lockReceiverChannel, unlockReceiverChannel chan bool, sendFn func(tea.Msg)) (receiverError error) {

	var beatChannel = t.beatsLooper.BeatChannel
	transmitPort, err := seqmidi.FindTransmitterPort()
	if err != nil {
		receiverError = fault.Wrap(err, fmsg.WithDesc("cannot find transmitport", "Could not find a transmitter. Start a seq program with the --transmit flag before starting a receiver"))
		return
	}
	err = transmitPort.Open()
	if err != nil {
		receiverError = fault.Wrap(err, fmsg.WithDesc("cannot open transmitport", "Could not open a transmitter.  Start a seq program with the --transmit flag before starting a receiver"))
		return
	}
	go func() {
		for {
			// NOTE: transmitPort must be redeclared within to loop to avoid memory error
			transmitPort, err := seqmidi.FindTransmitterPort()
			if err != nil {
				receiverError = fault.Wrap(err, fmsg.WithDesc("cannot find transmitport", "Could not find a transmitter. Start a seq program with the --transmit flag before starting a receiver"))
				return
			}
			err = transmitPort.Open()
			if err != nil {
				receiverError = fault.Wrap(err, fmsg.WithDesc("cannot open transmitport", "Could not open a transmitter.  Start a seq program with the --transmit flag before starting a receiver"))
				return
			}
			receiverChannel := make(chan TimingMsg)
			tickChannel := make(chan Timing)
			activeSenseChannel := make(chan bool)
			var loopMode playstate.LoopMode
			var ReceiverFunc ListenFn = func(msg []byte, milliseconds int32) {
				midiMessage := midi.Message(msg)
				switch midiMessage.Type() {
				case midi.SPPMsg:
					var ref uint16
					midiMessage.GetSPP(&ref)
					loopMode = playstate.LoopMode(ref)
				case midi.StartMsg:
					receiverChannel <- StartMsg{LoopMode: loopMode}
				case midi.StopMsg:
					receiverChannel <- StopMsg{}
				case midi.TimingClockMsg:
					tickChannel <- Timing{subdivisions: 24}
				case midi.ActiveSenseMsg:
					activeSenseChannel <- true
				default:
					println("receiving unknown msg")
					println(midiMessage.Type().String())
				}
			}
			stopFn, err := transmitPort.Listen(ReceiverFunc, drivers.ListenConfig{TimeCode: true, ActiveSense: true})
			if err != nil {
				sendFn(ErrorMsg{errors.New("error in setting up midi listener for transmitter")})
			}
			timer := time.AfterFunc(330*time.Millisecond, func() {
				sendFn(TransmitterNotConnectedMsg{})
				activeSenseChannel <- false
			})
			var command TimingMsg
		inner:
			for {
				select {
				case <-t.ctx.Done():
					return
				case <-lockReceiverChannel:
					stopFn()
					err := transmitPort.Close()
					if err != nil {
						wrappedErr := fault.Wrap(err, fmsg.With("transmit port not closed"))
						sendFn(ErrorMsg{wrappedErr})
					}
					timer.Stop()
					<-unlockReceiverChannel
					break inner
				case command = <-receiverChannel:
					switch command := command.(type) {
					case StartMsg:
						t.started = true
						t.playTime = time.Now()
						t.trackTime = time.Duration(0)
						t.pulseCount = 0
						sendFn(UIStartMsg{LoopMode: command.LoopMode})
					case StopMsg:
						t.started = false
						sendFn(UIStopMsg{})
						// m.playing should be false now.
					}
				case command = <-timingChannel:
					switch command := command.(type) {
					case TempoMsg:
						t.tempo = command.Tempo
						t.subdivisions = command.Subdivisions
					case QuitMsg:
						timer.Stop()
						stopFn()
					}
				case pulseTiming := <-tickChannel:
					timer.Reset(330 * time.Millisecond)
					if t.started {
						if t.pulseCount%(pulseTiming.subdivisions/t.subdivisions) == 0 {
							beatChannel <- beats.BeatMsg{Interval: t.TickInterval()}
						}
						t.pulseCount++
					}
				case isGood := <-activeSenseChannel:
					timer.Reset(330 * time.Millisecond)
					if isGood {
						sendFn(TransmitterConnectedMsg{})
					} else {
						sendFn(TransmitterNotConnectedMsg{})
					}
				}
			}
		}
	}()
	return nil
}

func (t *Timing) StandAloneLoop(sendFn func(tea.Msg)) {
	var beatChannel = t.beatsLooper.BeatChannel
	tickChannel := make(chan Timing)
	var command TimingMsg

	var tickTimer *time.Timer
	tick := func(adjustedInterval time.Duration) {
		tickTimer = time.AfterFunc(adjustedInterval, func() {
			tickChannel <- Timing{}
		})
	}

	go func() {
		for {
			select {
			case <-t.ctx.Done():
				return
			case command = <-timingChannel:
				switch command := command.(type) {
				case StartMsg:
					t.started = true
					t.playTime = time.Now()
					t.tempo = command.Tempo
					t.subdivisions = command.Subdivisions
					t.trackTime = time.Duration(0)
					t.preRollBeats = command.Prerollbeats
					tick(0)
				case StopMsg:
					t.started = false
					if tickTimer != nil {

						tickTimer.Stop()
					}
				case TempoMsg:
					t.tempo = command.Tempo
					t.subdivisions = command.Subdivisions
				}
			case <-tickChannel:
				if t.started {
					adjustedInterval := t.BeatInterval()
					tick(adjustedInterval)
					if t.preRollBeats == 0 {
						beatChannel <- beats.BeatMsg{Interval: adjustedInterval}
					} else {
						t.preRollBeats--
					}
				}
			}
		}
	}()
}

type TimingMsg = any
type StartMsg struct {
	Transmitting bool
	LoopMode     playstate.LoopMode
	Prerollbeats uint8
	Tempo        int
	Subdivisions int
}

type StopMsg struct{}
type AnticipatoryStopMsg struct{}
type QuitMsg struct{}

type TempoMsg struct {
	Tempo        int
	Subdivisions int
}

type ErrorMsg struct {
	error error
}

type UIStopMsg struct{}
type UIStartMsg struct{ LoopMode playstate.LoopMode }

type TransmitterConnectedMsg struct{}
type TransmitterNotConnectedMsg struct{}
