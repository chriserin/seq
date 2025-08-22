package timing

import (
	"errors"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chriserin/seq/internal/beats"
	midi "gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
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
}

type MidiLoopMode uint8

const (
	MlmStandAlone MidiLoopMode = iota
	MlmTransmitter
	MlmReceiver
)

func Loop(mode MidiLoopMode, lockReceiverChannel, unlockReceiverChannel chan bool, programChannel chan TimingMessage, beatChannel chan beats.BeatMsg, sendFn func(tea.Msg)) error {
	timing := Timing{}
	switch mode {
	case MlmStandAlone:
		timing.StandAloneLoop(programChannel, beatChannel, sendFn)
	case MlmTransmitter:
		err := timing.TransmitterLoop(programChannel, beatChannel, sendFn)
		if err != nil {
			return fault.Wrap(err, fmsg.With("cannot start transmitter loop"))
		}
	case MlmReceiver:
		err := timing.ReceiverLoop(lockReceiverChannel, unlockReceiverChannel, programChannel, beatChannel, sendFn)
		timing.StandAloneLoop(programChannel, beatChannel, sendFn)
		if err != nil {
			// NOTE: In case the receiver loop was not setup correctly, swallow the lock/unlock messages
			go func() {
				for {
					select {
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

func (tmtr Transmitter) Start() error {
	err := tmtr.out.Send(midi.Start())
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

const TransmitterName string = "seq-transmitter"

func (t *Timing) TransmitterLoop(programChannel chan TimingMessage, beatChannel chan beats.BeatMsg, sendFn func(tea.Msg)) error {
	driver, err := rtmididrv.New()
	if err != nil {
		return fault.Wrap(err, fmsg.With("midi driver error"))
	}
	out, err := driver.OpenVirtualOut(TransmitterName)
	if err != nil {
		return fault.Wrap(err, fmsg.With("could not open virtual out"))
	}
	transmitter := Transmitter{out}
	err = transmitter.ActiveSense()
	if err != nil {
		return fault.Wrap(err)
	}

	tickChannel := make(chan Timing)
	activeSenseChannel := make(chan bool)
	var command TimingMessage

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
			case command = <-programChannel:
				switch command := command.(type) {
				case StartMsg:
					t.started = true
					t.playTime = time.Now()
					t.tempo = command.Tempo
					t.subdivisions = command.Subdivisions
					t.trackTime = time.Duration(0)
					t.pulseCount = 0
					pulse(0)
					err := transmitter.Start()
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
				case TempoMsg:
					t.tempo = command.Tempo
					t.subdivisions = command.Subdivisions
				}
			case pulseTiming := <-tickChannel:
				if t.started {
					err := transmitter.Pulse()
					if err != nil {
						wrappedErr := fault.Wrap(err)
						sendFn(ErrorMsg{wrappedErr})
					}
					pulse(t.PulseInterval())
					if t.pulseCount%(pulseTiming.subdivisions/t.subdivisions) == 0 {
						beatChannel <- beats.BeatMsg{Interval: t.BeatInterval()}
					}
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

func (t *Timing) ReceiverLoop(lockReceiverChannel, unlockReceiverChannel chan bool, programChannel chan TimingMessage, beatChannel chan beats.BeatMsg, sendFn func(tea.Msg)) (receiverError error) {
	transmitPort, err := midi.FindInPort(TransmitterName)
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
			transmitPort, err = midi.FindInPort(TransmitterName)
			if err != nil {
				receiverError = fault.Wrap(err, fmsg.WithDesc("cannot find transmitport", "Could not find a transmitter. Start a seq program with the --transmit flag before starting a receiver"))
				return
			}
			err = transmitPort.Open()
			if err != nil {
				receiverError = fault.Wrap(err, fmsg.WithDesc("cannot open transmitport", "Could not open a transmitter.  Start a seq program with the --transmit flag before starting a receiver"))
				return
			}
			receiverChannel := make(chan TimingMessage)
			tickChannel := make(chan Timing)
			activeSenseChannel := make(chan bool)
			var ReceiverFunc ListenFn = func(msg []byte, milliseconds int32) {
				midiMessage := midi.Message(msg)
				switch midiMessage.Type() {
				case midi.StartMsg:
					receiverChannel <- StartMsg{}
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
			var command TimingMessage
		inner:
			for {
				select {
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
					switch command.(type) {
					case StartMsg:
						t.started = true
						t.playTime = time.Now()
						t.trackTime = time.Duration(0)
						t.pulseCount = 0
						sendFn(UIStartMsg{})
					case StopMsg:
						t.started = false
						sendFn(UIStopMsg{})
						// m.playing should be false now.
					}
				case command = <-programChannel:
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
							beatChannel <- beats.BeatMsg{Interval: t.BeatInterval()}
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

func (t *Timing) StandAloneLoop(programChannel chan TimingMessage, beatChannel chan beats.BeatMsg, sendFn func(tea.Msg)) {
	tickChannel := make(chan Timing)
	var command TimingMessage

	var tickTimer *time.Timer
	tick := func(adjustedInterval time.Duration) {
		tickTimer = time.AfterFunc(adjustedInterval, func() {
			tickChannel <- Timing{}
		})
	}

	go func() {
		for {
			select {
			case command = <-programChannel:
				switch command := command.(type) {
				case StartMsg:
					t.started = true
					t.playTime = time.Now()
					t.tempo = command.Tempo
					t.subdivisions = command.Subdivisions
					t.trackTime = time.Duration(0)
					tick(0)
				case StopMsg:
					t.started = false
					tickTimer.Stop()
				case TempoMsg:
					t.tempo = command.Tempo
					t.subdivisions = command.Subdivisions
				}
			case <-tickChannel:
				if t.started {
					adjustedInterval := t.BeatInterval()
					tick(adjustedInterval)
					beatChannel <- beats.BeatMsg{Interval: adjustedInterval}
				}
			}
		}
	}()
}

type TimingMessage = any
type StartMsg struct {
	Tempo        int
	Subdivisions int
}

type StopMsg struct {
}

type QuitMsg struct {
}

type TempoMsg struct {
	Tempo        int
	Subdivisions int
}

type ErrorMsg struct {
	error error
}

type UIStopMsg struct{}
type UIStartMsg struct{}

type TransmitterConnectedMsg struct{}
type TransmitterNotConnectedMsg struct{}
