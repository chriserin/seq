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
	IsTransmitter bool
	DoNotListen   bool
	outportName   string
	seqOutport    drivers.Out
	midiChannel   chan Message
	outDevices    []*OutDeviceInfo
	inDevices     []*InDeviceInfo
	TestQueue     *[]Message
	Test          bool
	StopFn        func()
	ReceiverFunc  ReceiverFunc
}

func (mc *MidiConnection) StopReceivingFromTransmitter() {
	if mc.StopFn != nil {
		mc.StopFn()
		mc.StopFn = nil
	}
}

func (mc *MidiConnection) HasDevices() bool {
	return len(mc.outDevices) > 0
}

func (mc *MidiConnection) EnsureConnection() {
	if !mc.HasConnection() {
		if len(mc.outDevices) > 0 {
			mc.outDevices[0].Open()
			mc.outDevices[0].Selected = true
		}
	}
}

func (mc *MidiConnection) HasConnection() bool {
	hasConnection := false
	for _, device := range mc.outDevices {
		if device.IsOpen && device.Selected {
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

func InitMidiConnection(createOut bool, outportName string, isTransmitter bool, ctx context.Context) *MidiConnection {
	var midiConn MidiConnection
	if createOut {
		midiConn = MidiConnection{midiChannel: make(chan Message), IsTransmitter: isTransmitter}
	} else {
		midiConn = MidiConnection{outportName: outportName, midiChannel: make(chan Message), IsTransmitter: isTransmitter}
	}

	return &midiConn
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
	if mc.seqOutport != nil {
		if mc.seqOutport.IsOpen() {
			err := mc.seqOutport.Send(msg)
			if err != nil {
				return err
			}
		}
	} else {
		// Send to all selected devices
		for _, device := range mc.outDevices {
			if device.Selected {
				if device.Out.IsOpen() {
					err := device.Out.Send(msg)
					if err != nil {
						return err
					}
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
	for _, device := range mc.outDevices {
		if device.IsOpen {
			_ = device.Out.Close()
			device.IsOpen = false
		}
	}
}

var dawOutports = []string{"Logic Pro Virtual In", "TESTDAW"}

func (mc MidiConnection) SendRecordMessage() error {

	var selectedOutport drivers.Out
	for _, device := range mc.outDevices {
		if device.IsDaw && device.IsOpen {
			selectedOutport = device.Out
			break
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
	return FindInPort(TransmitterName)
}

func FindInPort(inPortName string) (drivers.In, error) {
	ins, err := GetIns()
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("cannot get midi ins"))
	}
	for _, in := range ins {
		if strings.Contains(in.String(), inPortName) {
			return in, nil
		}
	}
	return nil, fault.New("cannot find transmitter port", fmsg.With("cannot find transmitter port"))
}

var OutputName string = "seq-cli-out"

func (mc *MidiConnection) CreateOutport() error {
	outport, err := OpenVirtualOut(OutputName)
	if err != nil {
		return fault.Wrap(err, fmsg.With("cannot create virtual outport"))
	}
	mc.seqOutport = outport
	return nil
}
