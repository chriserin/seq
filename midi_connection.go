package main

import (
	"sync"

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
		panic(err)
		// panic("Can't find outport")
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
		mc.connected = true
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

var playMutex = sync.Mutex{}

func (mc MidiConnection) AcquireSendFunc() SendFunc {
	if mc.outport.IsOpen() {
		sendFn, err := midi.SendTo(mc.outport)
		if err != nil {
			panic("Could not acquire send function")
		}
		return func(msg midi.Message) error {
			playMutex.Lock()
			error := sendFn(msg)
			playMutex.Unlock()
			return error
		}
	}
	panic("Midi outport is not open")
}
