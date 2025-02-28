package main

import (
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

func MidiEventLoop(mode MidiLoopMode, lockRecieverChannel, unlockReceiverChannel chan bool, programChannel chan midiEventLoopMsg, program *tea.Program) {
	timing := Timing{}
	switch mode {
	case MLM_STAND_ALONE:
		timing.StandAloneLoop(programChannel, program)
	case MLM_TRANSMITTER:
		timing.TransmitterLoop(programChannel, program)
	case MLM_RECEIVER:
		timing.ReceiverLoop(lockRecieverChannel, unlockReceiverChannel, programChannel, program)
		timing.StandAloneLoop(programChannel, program)
	}
}

type Transmitter struct {
	out drivers.Out
}

func (tmtr Transmitter) Start() {
	err := tmtr.out.Send(midi.Start())
	if err != nil {
		panic("could not send start")
	}
}

func (tmtr Transmitter) Stop() {
	err := tmtr.out.Send(midi.Stop())
	if err != nil {
		panic("could not send stop")
	}
}

func (tmtr Transmitter) Pulse() {
	err := tmtr.out.Send(midi.TimingClock())
	if err != nil {
		panic("could not send timing clock")
	}
}

func (tmtr Transmitter) ActiveSense() {
	err := tmtr.out.Send(midi.Activesense())
	if err != nil {
		panic("could not send activesense")
	}
}

const TRANSMITTER_NAME string = "seq-transmitter"

func (timing *Timing) TransmitterLoop(programChannel chan midiEventLoopMsg, program *tea.Program) {
	driver, err := rtmididrv.New()
	if err != nil {
		panic("Could not get driver")
	}
	out, err := driver.OpenVirtualOut(TRANSMITTER_NAME)
	if err != nil {
		panic("Could not open virtual In")
	}
	err = out.Send(midi.Activesense())
	if err != nil {
		panic("Could not send active sense")
	}
	transmitter := Transmitter{out}

	tickChannel := make(chan Timing)
	activeSenseChannel := make(chan bool)
	var command midiEventLoopMsg

	pulse := func(adjustedInterval time.Duration) {
		time.AfterFunc(adjustedInterval, func() {
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
					transmitter.Start()
				case stopMsg:
					timing.started = false
					transmitter.Stop()
					// m.playing should be false now.
				case tempoMsg:
					timing.tempo = command.tempo
					timing.subdivisions = command.subdivisions
				}
			case pulseTiming := <-tickChannel:
				if timing.started {
					transmitter.Pulse()
					timing.pulseCount++
					if timing.pulseCount%(pulseTiming.subdivisions/timing.subdivisions) == 0 {
						program.Send(beatMsg{timing.BeatInterval()})
					}
					pulse(timing.PulseInterval())
				}
			case <-activeSenseChannel:
				if !timing.started {
					transmitter.ActiveSense()
				}
				activesense()
			}
		}
	}()
	// Start active sense loop
	activesense()
}

type ListenFn func(msg []byte, milliseconds int32)

func (timing *Timing) ReceiverLoop(lockReceiverChannel, unlockReceiverChannel chan bool, programChannel chan midiEventLoopMsg, program *tea.Program) {
	go func() {
		for {
			transmitPort, err := midi.FindInPort(TRANSMITTER_NAME)
			if err != nil {
				panic("Could not find transmitter")
			}
			err = transmitPort.Open()
			if err != nil {
				panic("Could not open transmitter connection")
			}

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
				panic("error in setting up midi listener for transmitter")
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
					transmitPort.Close()
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
}

func (timing *Timing) StandAloneLoop(programChannel chan midiEventLoopMsg, program *tea.Program) {
	tickChannel := make(chan Timing)
	var command midiEventLoopMsg

	tick := func(adjustedInterval time.Duration) {
		time.AfterFunc(adjustedInterval, func() {
			tickChannel <- Timing{}
			program.Send(beatMsg{adjustedInterval})
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
					// m.playing should be false now.
				case tempoMsg:
					timing.tempo = command.tempo
					timing.subdivisions = command.subdivisions
				}
			case <-tickChannel:
				if timing.started {
					tick(timing.BeatInterval())
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
