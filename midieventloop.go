package main

import (
	"errors"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	tea "github.com/charmbracelet/bubbletea"
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

func MidiEventLoop(mode MidiLoopMode, lockReceiverChannel, unlockReceiverChannel chan bool, programChannel chan midiEventLoopMsg, program *tea.Program) error {
	timing := Timing{}
	switch mode {
	case MlmStandAlone:
		timing.StandAloneLoop(programChannel, program)
	case MlmTransmitter:
		err := timing.TransmitterLoop(programChannel, program)
		if err != nil {
			return fault.Wrap(err, fmsg.With("cannot start transmitter loop"))
		}
	case MlmReceiver:
		err := timing.ReceiverLoop(lockReceiverChannel, unlockReceiverChannel, programChannel, program)
		if err != nil {
			return fault.Wrap(err, fmsg.With("cannot start receiver loop"))
		}
		timing.StandAloneLoop(programChannel, program)
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

func (t *Timing) TransmitterLoop(programChannel chan midiEventLoopMsg, program *tea.Program) error {
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
	var command midiEventLoopMsg

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
				case startMsg:
					t.started = true
					t.playTime = time.Now()
					t.tempo = command.tempo
					t.subdivisions = command.subdivisions
					t.trackTime = time.Duration(0)
					t.pulseCount = 0
					pulse(t.PulseInterval())
					err := transmitter.Start()
					if err != nil {
						program.Send(errorMsg{err})
					}
				case stopMsg:
					t.started = false
					tickTimer.Stop()
					err := transmitter.Stop()
					if err != nil {
						program.Send(errorMsg{err})
					}
					// m.playing should be false now.
				case tempoMsg:
					t.tempo = command.tempo
					t.subdivisions = command.subdivisions
				}
			case pulseTiming := <-tickChannel:
				if t.started {
					err := transmitter.Pulse()
					if err != nil {
						wrappedErr := fault.Wrap(err)
						program.Send(errorMsg{wrappedErr})
					}
					t.pulseCount++
					if t.pulseCount%(pulseTiming.subdivisions/t.subdivisions) == 0 {
						program.Send(beatMsg{t.BeatInterval()})
					}
					pulse(t.PulseInterval())
				}
			case <-activeSenseChannel:
				if !t.started {
					err := transmitter.ActiveSense()
					if err != nil {
						wrappedErr := fault.Wrap(err, fmsg.With("activesense interrupted"))
						program.Send(errorMsg{wrappedErr})
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

func (t *Timing) ReceiverLoop(lockReceiverChannel, unlockReceiverChannel chan bool, programChannel chan midiEventLoopMsg, program *tea.Program) error {
	transmitPort, err := midi.FindInPort(TransmitterName)
	if err != nil {
		return fault.Wrap(err, fmsg.WithDesc("cannot find transmitport", "Could not find a transmitter. Start a seq program with the --transmit flag before starting a receiver"))
	}
	err = transmitPort.Open()
	if err != nil {
		return fault.Wrap(err, fmsg.WithDesc("cannot open transmitport", "Could not open a transmitter.  Start a seq program with the --transmit flag before starting a receiver"))
	}
	go func() {
		for {
			receiverChannel := make(chan midiEventLoopMsg)
			tickChannel := make(chan Timing)
			activeSenseChannel := make(chan bool)
			var ReceiverFunc ListenFn = func(msg []byte, milliseconds int32) {
				midiMessage := midi.Message(msg)
				switch midiMessage.Type() {
				case midi.StartMsg:
					receiverChannel <- startMsg{}
				case midi.StopMsg:
					receiverChannel <- stopMsg{}
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
				program.Send(errorMsg{errors.New("error in setting up midi listener for transmitter")})
			}
			timer := time.AfterFunc(330*time.Millisecond, func() {
				program.Send(uiNotConnectedMsg{})
				activeSenseChannel <- false
			})
			var command midiEventLoopMsg
		inner:
			for {
				select {
				case <-lockReceiverChannel:
					stopFn()
					err := transmitPort.Close()
					if err != nil {
						wrappedErr := fault.Wrap(err, fmsg.With("transmit port not closed"))
						program.Send(errorMsg{wrappedErr})
					}
					timer.Stop()
					<-unlockReceiverChannel
					break inner
				case command = <-receiverChannel:
					switch command.(type) {
					case startMsg:
						t.started = true
						t.playTime = time.Now()
						t.trackTime = time.Duration(0)
						t.pulseCount = 0
						program.Send(uiStartMsg{})
					case stopMsg:
						t.started = false
						program.Send(uiStopMsg{})
						// m.playing should be false now.
					}
				case command = <-programChannel:
					switch command := command.(type) {
					case tempoMsg:
						t.tempo = command.tempo
						t.subdivisions = command.subdivisions
					case quitMsg:
						timer.Stop()
						stopFn()
					}
				case pulseTiming := <-tickChannel:
					timer.Reset(330 * time.Millisecond)
					if t.started {
						t.pulseCount++
						if t.pulseCount%(pulseTiming.subdivisions/t.subdivisions) == 0 {
							program.Send(beatMsg{t.BeatInterval()})
						}
					}
				case isGood := <-activeSenseChannel:
					timer.Reset(330 * time.Millisecond)
					if isGood {
						program.Send(uiConnectedMsg{})
					} else {
						program.Send(uiNotConnectedMsg{})
					}
				}
			}
		}
	}()
	return nil
}

func (t *Timing) StandAloneLoop(programChannel chan midiEventLoopMsg, program *tea.Program) {
	tickChannel := make(chan Timing)
	var command midiEventLoopMsg

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
				case startMsg:
					t.started = true
					t.playTime = time.Now()
					t.tempo = command.tempo
					t.subdivisions = command.subdivisions
					t.trackTime = time.Duration(0)
					tick(t.BeatInterval())
				case stopMsg:
					t.started = false
					tickTimer.Stop()
				case tempoMsg:
					t.tempo = command.tempo
					t.subdivisions = command.subdivisions
				}
			case <-tickChannel:
				if t.started {
					adjustedInterval := t.BeatInterval()
					program.Send(beatMsg{adjustedInterval})
					tick(adjustedInterval)
				}
			}
		}
	}()
}

type midiEventLoopMsg = any
type startMsg struct {
	tempo        int
	subdivisions int
}

type stopMsg struct {
}
type quitMsg struct {
}
type tempoMsg struct {
	tempo        int
	subdivisions int
}
