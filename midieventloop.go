package main

import (
	"errors"
	"fmt"
	"time"

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
	MLM_STAND_ALONE MidiLoopMode = iota
	MLM_TRANSMITTER
	MLM_RECEIVER
)

func MidiEventLoop(mode MidiLoopMode, lockRecieverChannel, unlockReceiverChannel chan bool, programChannel chan midiEventLoopMsg, program *tea.Program) error {
	timing := Timing{}
	switch mode {
	case MLM_STAND_ALONE:
		timing.StandAloneLoop(programChannel, program)
	case MLM_TRANSMITTER:
		err := timing.TransmitterLoop(programChannel, program)
		if err != nil {
			return fmt.Errorf("could not setup transmitter loop: %w", err)
		}
	case MLM_RECEIVER:
		err := timing.ReceiverLoop(lockRecieverChannel, unlockReceiverChannel, programChannel, program)
		if err != nil {
			return fmt.Errorf("could not setup receiver loop: %w", err)
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
		return fmt.Errorf("could not send midi start %w", err)
	}
	return nil
}

func (tmtr Transmitter) Stop() error {
	err := tmtr.out.Send(midi.Stop())
	if err != nil {
		return fmt.Errorf("could not send midi stop %w", err)
	}
	return nil
}

func (tmtr Transmitter) Pulse() error {
	err := tmtr.out.Send(midi.TimingClock())
	if err != nil {
		return fmt.Errorf("could not send midi timing clock msg %w", err)
	}
	return nil
}

func (tmtr Transmitter) ActiveSense() error {
	err := tmtr.out.Send(midi.Activesense())
	if err != nil {
		return fmt.Errorf("could not send midi active sense msg %w", err)
	}
	return nil
}

const TRANSMITTER_NAME string = "seq-transmitter"

func (timing *Timing) TransmitterLoop(programChannel chan midiEventLoopMsg, program *tea.Program) error {
	driver, err := rtmididrv.New()
	if err != nil {
		return fmt.Errorf("could not get midi driver: %w", err)
	}
	out, err := driver.OpenVirtualOut(TRANSMITTER_NAME)
	if err != nil {
		return fmt.Errorf("could not open virtual out: %w", err)
	}
	err = out.Send(midi.Activesense())
	if err != nil {
		return fmt.Errorf("could not send active sense: %w", err)
	}
	transmitter := Transmitter{out}

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
					timing.started = true
					timing.playTime = time.Now()
					timing.tempo = command.tempo
					timing.subdivisions = command.subdivisions
					timing.trackTime = time.Duration(0)
					timing.pulseCount = 0
					pulse(timing.PulseInterval())
					err := transmitter.Start()
					if err != nil {
						program.Send(errorMsg{err})
					}
				case stopMsg:
					timing.started = false
					tickTimer.Stop()
					err := transmitter.Stop()
					if err != nil {
						program.Send(errorMsg{err})
					}
					// m.playing should be false now.
				case tempoMsg:
					timing.tempo = command.tempo
					timing.subdivisions = command.subdivisions
				}
			case pulseTiming := <-tickChannel:
				if timing.started {
					err := transmitter.Pulse()
					if err != nil {
						wrappedErr := fmt.Errorf("pulse timer interrupted %w", err)
						program.Send(errorMsg{wrappedErr})
					}
					timing.pulseCount++
					if timing.pulseCount%(pulseTiming.subdivisions/timing.subdivisions) == 0 {
						program.Send(beatMsg{timing.BeatInterval()})
					}
					pulse(timing.PulseInterval())
				}
			case <-activeSenseChannel:
				if !timing.started {
					err := transmitter.ActiveSense()
					if err != nil {
						wrappedErr := fmt.Errorf("active sense interrupted %w", err)
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

func (timing *Timing) ReceiverLoop(lockReceiverChannel, unlockReceiverChannel chan bool, programChannel chan midiEventLoopMsg, program *tea.Program) error {
	transmitPort, err := midi.FindInPort(TRANSMITTER_NAME)
	if err != nil {
		return fmt.Errorf("could not find transmitter port %w", err)
	}
	err = transmitPort.Open()
	if err != nil {
		return fmt.Errorf("could not open transmitter connection %w", err)
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
						wrappedErr := fmt.Errorf("transmit port not closed: %w", err)
						program.Send(errorMsg{wrappedErr})
					}
					timer.Stop()
					<-unlockReceiverChannel
					break inner
				case command = <-receiverChannel:
					switch command.(type) {
					case startMsg:
						timing.started = true
						timing.playTime = time.Now()
						timing.trackTime = time.Duration(0)
						timing.pulseCount = 0
						program.Send(uiStartMsg{})
					case stopMsg:
						timing.started = false
						program.Send(uiStopMsg{})
						// m.playing should be false now.
					}
				case command = <-programChannel:
					switch command := command.(type) {
					case tempoMsg:
						timing.tempo = command.tempo
						timing.subdivisions = command.subdivisions
					case quitMsg:
						timer.Stop()
						stopFn()
					}
				case pulseTiming := <-tickChannel:
					timer.Reset(330 * time.Millisecond)
					if timing.started {
						timing.pulseCount++
						if timing.pulseCount%(pulseTiming.subdivisions/timing.subdivisions) == 0 {
							program.Send(beatMsg{timing.BeatInterval()})
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

func (timing *Timing) StandAloneLoop(programChannel chan midiEventLoopMsg, program *tea.Program) {
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
					timing.started = true
					timing.playTime = time.Now()
					timing.tempo = command.tempo
					timing.subdivisions = command.subdivisions
					timing.trackTime = time.Duration(0)
					tick(timing.BeatInterval())
				case stopMsg:
					timing.started = false
					tickTimer.Stop()
				case tempoMsg:
					timing.tempo = command.tempo
					timing.subdivisions = command.subdivisions
				}
			case <-tickChannel:
				if timing.started {
					adjustedInterval := timing.BeatInterval()
					program.Send(beatMsg{adjustedInterval})
					tick(adjustedInterval)
				}
			}
		}
	}()
}

type midiEventLoopMsg = interface{}
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
