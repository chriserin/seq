package main

import (
	midi "gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"

	_ "gitlab.com/gomidi/midi/v2/drivers/portmididrv"
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

func (mc *MidiConnection) Connect() bool {
	outport, err := midi.OutPort(0)
	if err != nil {
		panic("Can't find outport")
	}
	mc.outport = outport
	return true
}

func (mc *MidiConnection) ConnectAndOpen() error {
	if !mc.connected && mc.Connect() {
		err := mc.outport.Open()
		if err != nil {
			return err
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

func (mc MidiConnection) IsOpen() bool {
	return mc.connected && mc.outport.IsOpen()
}

func (mc MidiConnection) AcquireSendFunc() SendFunc {
	if mc.outport.IsOpen() {
		sendFn, err := midi.SendTo(mc.outport)
		if err != nil {
			panic("Could not acquire send function")
		}
		return sendFn
	}
	panic("Midi outport is not open")
}
