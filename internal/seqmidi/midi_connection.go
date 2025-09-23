// Package seqmidi provides MIDI connection management and message sending
// functionality for the sequencer. It handles virtual and physical MIDI
// port connections, thread-safe message transmission, and DAW integration
// features including record triggering for external recording systems.
package seqmidi

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	midi "gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)

const TransmitterName string = "seq-transmitter"

type MidiConnection struct {
	outportName string
	midiChannel chan Message
	devices     []*DeviceInfo
	TestQueue   *[]Message
	Test        bool
}

func (mc *MidiConnection) EnsureConnection() {
	if !mc.HasConnection() {
		mc.devices[0].Open()
		mc.devices[0].Selected = true
	}
}

func (mc *MidiConnection) HasConnection() bool {
	hasConnection := false
	for _, device := range mc.devices {
		if device.IsOpen {
			hasConnection = true
			break
		}
	}
	return hasConnection
}

type Message struct {
	Delay time.Duration
	Msg   midi.Message
}

var OutputName string = "seq-cli-out"

func InitMidiConnection(createOut bool, outportName string, ctx context.Context) (*MidiConnection, error) {
	var midiConn MidiConnection
	if createOut {
		out, err := SeqOut()
		if err != nil {
			return nil, fault.Wrap(err, fmsg.With("cannot create seq midi out"))
		}
		midiConn = MidiConnection{devices: []*DeviceInfo{{Out: out, Selected: true, IsOpen: true}}, midiChannel: make(chan Message)}
	} else {
		midiConn = MidiConnection{outportName: outportName, midiChannel: make(chan Message)}
	}

	return &midiConn, nil
}

var playMutex = sync.Mutex{}

func (mc *MidiConnection) LoopMidi(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-mc.midiChannel:
				if msg.Delay == 0 {
					playMutex.Lock()
					err := mc.SendMidi(msg.Msg)
					playMutex.Unlock()
					if err != nil {
						panic(err)
					}
				} else {
					time.AfterFunc(msg.Delay, func() {
						playMutex.Lock()
						err := mc.SendMidi(msg.Msg)
						playMutex.Unlock()
						if err != nil {
							panic(err)
						}
					})
				}
			}
		}
	}()
}

func (mc *MidiConnection) Send(msg Message) {
	if mc.Test {
		*mc.TestQueue = append(*mc.TestQueue, msg)
	} else {
		mc.midiChannel <- msg
	}
}

func (mc MidiConnection) SendMidi(msg midi.Message) error {
	// Send to all selected devices
	for _, device := range mc.devices {
		if device.Selected {
			if device.Out.IsOpen() {
				err := device.Out.Send(msg)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (mc *MidiConnection) Panic() error {
	// NOTE: No connection means nothing to panic about
	for i := range 127 {
		err := mc.SendMidi(midi.NoteOff(0, uint8(i)))
		if err != nil {
			return fault.Wrap(err, fmsg.With("cannot send panic note off"))
		}
	}

	return nil
}

func (mc *MidiConnection) Close() {
	for _, device := range mc.devices {
		if device.IsOpen {
			_ = device.Out.Close()
			device.IsOpen = false
		}
	}
}

var dawOutports = []string{"Logic Pro Virtual In", "TESTDAW"}

func (mc MidiConnection) SendRecordMessage() error {

	var selectedOutport drivers.Out
foundtheport:
	for _, device := range mc.devices {
		for _, dawName := range dawOutports {
			if strings.Contains(device.Name, dawName) {
				device.Open()
				selectedOutport = device.Out
				break foundtheport
			}
		}
	}

	if selectedOutport == nil {
		return fault.New("could not find daw outport", fmsg.WithDesc("could not find daw outport", "Could not find a daw outport, please ensure the DAW record destination is open and then restart seq"))
	}

	// NOTE: The record message must be configured in Logic.
	// Logic Pro -> Control Surfaces -> Controller Assignments...
	// Create a new zone (seq), a new Mode (commands), a new control
	// Within the new control press "Learn" and then in seq use the PlayRecord mapping (`:<Space>`).
	// Then Chose `Class: Key Command` `Command: Global Commands` and then `Record` in the unlabeled parameter list
	err := selectedOutport.Send(midi.ControlChange(16, 127, 127))
	if err != nil {
		return fault.Wrap(err, fmsg.With("could not send record message"))
	}
	return nil
}

func FindTransmitterPort() (drivers.In, error) {
	return midi.FindInPort(TransmitterName)
}
