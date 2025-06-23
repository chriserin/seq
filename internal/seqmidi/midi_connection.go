package seqmidi

import (
	"sync"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	midi "gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"

	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

type MidiConnection struct {
	outport   drivers.Out
	connected bool
}

type SendFunc func(msg midi.Message) error

func InitMidiConnection() MidiConnection {
	return MidiConnection{connected: false}
}

func (mc MidiConnection) HasConnection() bool {
	return mc.connected
}

func (mc *MidiConnection) Connect() error {
	outport, err := midi.OutPort(0)
	if err != nil {
		return fault.Wrap(err, fmsg.WithDesc("midi outport not available", "There are no midi out ports available."))
	}
	mc.outport = outport
	return nil
}

func (mc *MidiConnection) ConnectAndOpen() error {
	if !mc.connected {
		err := mc.Connect()
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

func (mc *MidiConnection) Close() {
	if mc.connected {
		if mc.outport.IsOpen() {
			mc.Close()
		}
	}
}

func (mc MidiConnection) IsReady() bool {
	return mc.connected && mc.outport.IsOpen()
}

var playMutex = sync.Mutex{}

func (mc MidiConnection) AcquireSendFunc() (SendFunc, error) {
	err := mc.ConnectAndOpen()
	if err != nil {
		return nil, fault.Wrap(err)
	}
	//NOTE: midi library checks IsOpen(), tries to open and returns error
	sendFn, err := midi.SendTo(mc.outport)
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
