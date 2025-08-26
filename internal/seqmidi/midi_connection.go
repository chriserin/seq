// Package seqmidi provides MIDI connection management and message sending
// functionality for the sequencer. It handles virtual and physical MIDI
// port connections, thread-safe message transmission, and DAW integration
// features including record triggering for external recording systems.
package seqmidi

import (
	"strings"
	"sync"

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
}

type SendFunc func(msg midi.Message) error

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

		return MidiConnection{connected: true, outport: out}, nil
	} else {
		return MidiConnection{outportName: outportName, connected: false}, nil
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
	}
	return nil
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

var sendFn SendFunc

func (mc MidiConnection) AcquireSendFunc() (SendFunc, error) {
	if mc.Test {
		return func(msg midi.Message) error {
			return nil
		}, nil
	}

	if sendFn != nil {
		return sendFn, nil
	}

	err := mc.ConnectAndOpen()
	if err != nil {
		return nil, fault.Wrap(err)
	}
	//NOTE: midi library checks IsOpen(), tries to open and returns error
	newSendFn, err := midi.SendTo(mc.outport)
	sendFn = newSendFn
	if err != nil {
		return nil, fault.Wrap(err)
	}
	return func(msg midi.Message) error {
		// NOTE: midi send reuses a global state buffer so we need to ensure that the buffer is accessed sequentially
		playMutex.Lock()
		defer playMutex.Unlock()
		error := sendFn(msg)
		return error
	}, nil
}

var dawOutports = []string{"Logic Pro Virtual In"}

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
