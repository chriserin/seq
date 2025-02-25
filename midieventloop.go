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

func MidiEventLoop(mode MidiLoopMode, programChannel chan midiEventLoopMsg, program *tea.Program) {
	switch mode {
	case MLM_STAND_ALONE:
		StandAloneLoop(programChannel, program)
	case MLM_TRANSMITTER:
		TransmitterLoop(programChannel, program)
	case MLM_RECEIVER:
		ReceiverLoop(programChannel, program)
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

const TRANSMITTER_NAME string = "seq-transmitter"

func TransmitterLoop(programChannel chan midiEventLoopMsg, program *tea.Program) {
	driver, err := rtmididrv.New()
	if err != nil {
		panic("Could not get driver")
	}
	out, err := driver.OpenVirtualOut(TRANSMITTER_NAME)
	if err != nil {
		panic("Could not open virtual In")
	}
	out.Send(midi.Activesense())
	transmitter := Transmitter{out}

	tickChannel := make(chan Timing)
	var command midiEventLoopMsg

	pulse := func(adjustedInterval time.Duration) {
		time.AfterFunc(adjustedInterval, func() {
			tickChannel <- Timing{subdivisions: 24}
		})
	}

	go func() {
		timing := Timing{}
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
					pulse(timing.PulseInterval())
					transmitter.Pulse()
					if timing.pulseCount%(pulseTiming.subdivisions/timing.subdivisions) == 0 {
						program.Send(beatMsg{timing.BeatInterval()})
					}
					timing.pulseCount++
				}
			}
		}
	}()
}

type ListenFn func(msg []byte, milliseconds int32)

func ReceiverLoop(programChannel chan midiEventLoopMsg, program *tea.Program) {
	inport, err := midi.FindInPort(TRANSMITTER_NAME)
	if err != nil {
		panic("Could not find transmitter")
	}
	err = inport.Open()
	if err != nil {
		panic("Could not open transmitter connection")
	} else {
		println("successful open")
	}
	tickChannel := make(chan Timing)
	var ReceiverFunc ListenFn = func(msg []byte, milliseconds int32) {
		midiMessage := midi.Message(msg)
		switch midiMessage.Type() {
		case midi.StartMsg:
			programChannel <- startMsg{}
		case midi.StopMsg:
			programChannel <- stopMsg{}
		case midi.TimingClockMsg:
			tickChannel <- Timing{subdivisions: 24}
		default:
			println("receiving unknown msg")
			println(midiMessage.Type().String())
		}
	}
	inport.Listen(ReceiverFunc, drivers.ListenConfig{TimeCode: true})
	var command midiEventLoopMsg
	go func() {
		timing := Timing{}
		for {
			select {
			case command = <-programChannel:
				switch command := command.(type) {
				case startMsg:
					timing.started = true
					timing.playTime = time.Now()
					timing.trackTime = time.Duration(0)
					program.Send(uiStartMsg{})
				case stopMsg:
					timing.started = false
					// m.playing should be false now.
				case tempoMsg:
					timing.tempo = command.tempo
					timing.subdivisions = command.subdivisions
				}
			case pulseTiming := <-tickChannel:
				if timing.started {
					if timing.pulseCount%(pulseTiming.subdivisions/timing.subdivisions) == 0 {
						program.Send(beatMsg{timing.BeatInterval()})
					}
					timing.pulseCount++
				}
			}
		}
	}()
}

func StandAloneLoop(programChannel chan midiEventLoopMsg, program *tea.Program) {
	tickChannel := make(chan Timing)
	var command midiEventLoopMsg

	tick := func(adjustedInterval time.Duration) {
		time.AfterFunc(adjustedInterval, func() {
			tickChannel <- Timing{}
			program.Send(beatMsg{adjustedInterval})
		})
	}

	go func() {
		timing := Timing{}
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
type tempoMsg struct {
	tempo        int
	subdivisions int
}
