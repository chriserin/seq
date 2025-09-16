// Package seqmidi provides MIDI connection management and message sending
// functionality for the sequencer. It handles virtual and physical MIDI
// port connections, thread-safe message transmission, and DAW integration
// features including record triggering for external recording systems.
package seqmidi

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	midi "gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"

	rtmididrv "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

type MidiConnection struct {
	outportName string
	outport     drivers.Out
	connected   bool
	Test        bool
	midiChannel chan Message
	TestQueue   *[]Message
}

type Message struct {
	Delay time.Duration
	Msg   midi.Message
}

type SendFunc func(Message)

var OutputName string = "seq-cli-out"

func InitMidiConnection(createOut bool, outportName string) (MidiConnection, error) {
	if createOut {
		driver, err := rtmididrv.New()
		if err != nil {
			return MidiConnection{}, fault.Wrap(err, fmsg.With("midi driver error"))
		}
		out, err := driver.OpenVirtualOut(OutputName)
		if err != nil {
			return MidiConnection{}, fault.Wrap(err, fmsg.With("cannot open virtual out"))
		}

		return MidiConnection{connected: true, outport: out, midiChannel: make(chan Message)}, nil
	} else {
		return MidiConnection{outportName: outportName, connected: false, midiChannel: make(chan Message)}, nil
	}
}

func (mc MidiConnection) HasConnection() bool {
	return mc.connected
}

func (mc *MidiConnection) Connect(portnumber int) error {
	outport, err := midi.OutPort(portnumber)
	if err != nil {
		return fault.Wrap(err, fmsg.WithDesc("midi outport not available", "There are no midi out ports available."))
	}
	mc.outport = outport
	return nil
}

func (mc *MidiConnection) ConnectAndOpen() error {
	if !mc.connected {
		outports := midi.GetOutPorts()
		//NOTE: Default to 0 (the first midi outport) if not found or if empty string
		var outportIndex = 0
		for i, outport := range outports {
			if mc.outportName != "" && strings.Contains(outport.String(), mc.outportName) {
				outportIndex = i
			}
		}

		err := mc.Connect(outportIndex)
		if err != nil {
			return fault.Wrap(err)
		}
		mc.connected = true
	}
	if !mc.outport.IsOpen() {
		err := mc.outport.Open()
		if err != nil {
			return fault.Wrap(err, fmsg.With("cannot open midi port"))
		}
	} else {
		mc.LoopMidi()
	}
	return nil
}

func (mc *MidiConnection) LoopMidi() {
	go func() {
		for {
			msg := <-mc.midiChannel
			if msg.Delay == 0 {
				playMutex.Lock()
				err := mc.outport.Send(msg.Msg)
				playMutex.Unlock()
				if err != nil {
					panic(err)
				}
			} else {
				time.AfterFunc(msg.Delay, func() {
					playMutex.Lock()
					err := mc.outport.Send(msg.Msg)
					playMutex.Unlock()
					if err != nil {
						panic(err)
					}
				})
			}
		}
	}()
}

func (mc *MidiConnection) Panic() error {
	// NOTE: No connection means nothing to panic about
	if mc.connected {
		for i := range 127 {
			err := mc.outport.Send(midi.NoteOff(0, uint8(i)))
			if err != nil {
				return fault.Wrap(err, fmsg.With("cannot send panic note off"))
			}
		}
	}

	return nil
}

func (mc *MidiConnection) Close() {
	if mc.connected {
		if mc.outport.IsOpen() {
			err := mc.outport.Close()
			if err != nil {
				panic("Could not close connection")
			}
		}
	}
}

func (mc MidiConnection) IsReady() bool {
	return mc.connected && mc.outport.IsOpen()
}

var playMutex = sync.Mutex{}

func (mc MidiConnection) AcquireSendFunc() (SendFunc, error) {
	if mc.Test {
		return func(seqmidiMessage Message) {
			(*mc.TestQueue) = append((*mc.TestQueue), seqmidiMessage)
			fmt.Println("TEST SEND", mc, len(*mc.TestQueue), seqmidiMessage)
		}, nil
	}

	// Ensure connection is open
	err := mc.ConnectAndOpen()
	if err != nil {
		return nil, fault.Wrap(err)
	}
	sendFn := func(msg Message) {
		mc.midiChannel <- msg
	}
	return sendFn, nil
}

var dawOutports = []string{"Logic Pro Virtual In", "TESTDAW"}

func SendRecordMessage() error {
	outports := midi.GetOutPorts()

	var selectedOutport drivers.Out
foundtheport:
	for _, outport := range outports {
		for _, name := range dawOutports {
			if strings.Contains(outport.String(), name) {
				selectedOutport = outport
				break foundtheport

			}
		}
	}

	if selectedOutport == nil {
		return fault.New("could not find daw outport", fmsg.WithDesc("could not find daw outport", "Could not find a daw outport, please ensure the DAW record destination is open and then restart seq"))
	}

	err := selectedOutport.Open()
	if err != nil {
		return fault.Wrap(err, fmsg.With("cannot open midi port"))
	}
	// NOTE: The record message must be configured in Logic.
	// Logic Pro -> Control Surfaces -> Controller Assignments...
	// Create a new zone (seq), a new Mode (commands), a new control
	// Within the new control press "Learn" and then in seq use the PlayRecord mapping (`:<Space>`).
	// Then Chose `Class: Key Command` `Command: Global Commands` and then `Record` in the unlabeled parameter list
	err = selectedOutport.Send(midi.ControlChange(16, 127, 127))
	if err != nil {
		return fault.Wrap(err, fmsg.With("could not send record message"))
	}
	return nil
}
