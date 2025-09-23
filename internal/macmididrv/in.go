//go:build !js
// +build !js

package macmididrv

import (
	"fmt"
	"math"

	"github.com/chriserin/seq/internal/macmididrv/imported/macmidi"
	"gitlab.com/gomidi/midi/v2/drivers"
)

type in struct {
	number int
	//sync.RWMutex
	//listenerSet bool
	driver *Driver
	name   string
	midiIn macmidi.MIDIIn
}

// IsOpen returns whether the MIDI in port is open
func (i *in) IsOpen() (open bool) {
	//	i.RLock()
	open = i.midiIn != nil
	//i.RUnlock()
	return
}

// String returns the name of the MIDI in port.
func (i *in) String() string {
	return i.name
}

// Underlying returns the underlying macmidi.MIDIIn. Use it with type casting:
//
//	rtIn := i.Underlying().(macmidi.MIDIIn)
func (i *in) Underlying() interface{} {
	return i.midiIn
}

// Number returns the number of the MIDI in port.
// Note that with macmidi, out and in ports are counted separately.
// That means there might exists out ports and an in ports that share the same number.
func (i *in) Number() int {
	return i.number
}

// Close closes the MIDI in port, after it has stopped listening.
func (i *in) Close() (err error) {
	if !i.IsOpen() {
		return nil
	}

	//	i.StopListening()
	//	i.Lock()
	err = i.midiIn.Close()
	i.midiIn = nil
	//	i.Unlock()
	return
}

// Open opens the MIDI in port
func (i *in) Open() (err error) {
	if i.IsOpen() {
		return nil
	}

	//i.Lock()

	i.midiIn, err = macmidi.NewMIDIInDefault()
	if err != nil {
		i.midiIn = nil
		//i.Unlock()
		return fmt.Errorf("can't open default MIDI in: %v", err)
	}

	err = i.midiIn.OpenPort(i.number, "")
	//i.Unlock()

	if err != nil {
		i.Close()
		return fmt.Errorf("can't open MIDI in port %v (%s): %v", i.number, i, err)
	}

	//i.driver.Lock()
	//i.midiIn.IgnoreTypes(i.driver.ignoreSysex, i.driver.ignoreTimeCode, i.driver.ignoreActiveSense)
	i.driver.opened = append(i.driver.opened, i)
	//i.driver.Unlock()

	return nil
}

/*
type readerState int

const (
	readerStateClean                readerState = 0
	readerStateWithinChannelMessage readerState = 1
	readerStateWithinSysCommon      readerState = 2
	readerStateInSysEx              readerState = 3
	readerStateWithinUnknown        readerState = 4
)

const (
	byteMIDITimingCodeMessage  = byte(0xF1)
	byteSysSongPositionPointer = byte(0xF2)
	byteSysSongSelect          = byte(0xF3)
	byteSysTuneRequest         = byte(0xF6)
)

const (
	byteProgramChange         = 0xC
	byteChannelPressure       = 0xD
	byteNoteOff               = 0x8
	byteNoteOn                = 0x9
	bytePolyphonicKeyPressure = 0xA
	byteControlChange         = 0xB
	bytePitchWheel            = 0xE
)
*/

func newIn(driver *Driver, number int, name string) drivers.In {
	return &in{driver: driver, number: number, name: name}
}

func (i *in) Listen(onMsg func(msg []byte, milliseconds int32), config drivers.ListenConfig) (stopFn func(), err error) {

	if onMsg == nil {
		return nil, fmt.Errorf("onMsg callback must not be nil")
	}

	i.midiIn.IgnoreTypes(!config.SysEx, !config.TimeCode, !config.ActiveSense)

	//var inSysEx bool
	if config.SysExBufferSize == 0 {
		config.SysExBufferSize = 1024
	}

	var rd = drivers.NewReader(config, onMsg)

	stopFn = func() {
		i.midiIn.CancelCallback()
	}

	go i.midiIn.SetCallback(func(in macmidi.MIDIIn, bt []byte, deltaSeconds float64) {

		rd.EachMessage(bt, int32(math.Round(deltaSeconds*1000)))

	})

	return stopFn, nil
}
