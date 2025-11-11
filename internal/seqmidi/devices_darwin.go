//go:build darwin

package seqmidi

import "C"
import (
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/chriserin/sq/internal/macmididrv"
	"github.com/chriserin/sq/internal/macmididrv/imported/macmidi"
	"gitlab.com/gomidi/midi/v2/drivers"
)

const (
	DeviceAdded   = 2
	DeviceRemoved = 3
)

var deviceChannel chan struct{}

func init() {
	deviceChannel = make(chan struct{})
}

func (mc MidiConnection) WaitUntilDevicesQueried() {
	deviceChannel <- struct{}{}
}

func driver() (*macmididrv.Driver, error) {
	deviceChannel <- struct{}{}
	driver, err := macmididrv.New()
	return driver, err
}

func OpenVirtualOut(name string) (drivers.Out, error) {
	driver, err := driver()
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("midi driver error"))
	}
	out, err := driver.OpenVirtualOut(name)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("cannot open virtual out"))
	}
	return out, nil
}

func (mc *MidiConnection) UpdateInDeviceList(driver drivers.Driver) error {
	var newDevices []*InDeviceInfo

	ins, err := driver.Ins()
	if err != nil {
		return fmt.Errorf("failed to get input ports: %v", err)
	}

	for _, in := range ins {
		// Check if we already have this device
		var foundDevice *InDeviceInfo
		for _, currentDevice := range mc.inDevices {
			if currentDevice.Name == in.String() {
				foundDevice = currentDevice
				foundDevice.In = in
				foundDevice.IsOpen = false
				if foundDevice.IsTransmitter {
					foundDevice.Open()
				}
				break
			}
		}

		if foundDevice == nil {
			// NOTE: Don't connect to ourself
			if !(mc.IsTransmitter && strings.Contains(in.String(), TransmitterName)) {
				newDevice := &InDeviceInfo{
					Name: in.String(),
					In:   in,
				}
				if newDevice.Matches(TransmitterName) {
					newDevice.IsTransmitter = true
					newDevice.Open()
				}
				newDevices = append(newDevices, newDevice)
			}
		} else {
			newDevices = append(newDevices, foundDevice)
		}
	}

	mc.inDevices = newDevices

	if mc.ReceiverFunc != nil && !mc.DoNotListen {
		err := mc.ListenToTransmitter(mc.ReceiverFunc)
		if err != nil {
			return fault.Wrap(err, fmsg.With("cannot listen to transmitter"))
		}
	}

	return nil
}

func (mc *MidiConnection) DeviceLoop(ctx context.Context) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "Recovered in MIDI device loop from panic: %v\n", r)
				debug.PrintStack()
			}
		}()

		//NOTE: The client notification will only be called on the same thread as the client initialization, we have to "pump the run loop" on that thread.
		runtime.LockOSThread()

		// NOTE: Currently this New call does some C singleton initialization and is called here explicitly for that.
		driver, err := macmididrv.New()
		deviceCallback := func(messageID int, name string) {
			err := mc.UpdateOutDeviceList(driver)
			if err != nil {
				fmt.Printf("Can't update MIDI out device list: %v\n", err)
				return
			}
			err = mc.UpdateInDeviceList(driver)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Can't update MIDI in device list: %v\n", err)
				return
			}
		}
		macmididrv.SetNotificationCallback(deviceCallback)
		if err != nil {
			fmt.Printf("Can't open MIDI driver: %v\n", err)
			return
		}
		err = mc.UpdateOutDeviceList(driver)
		if err != nil {
			fmt.Printf("Can't update MIDI out device list: %v\n", err)
			return
		}
		err = mc.UpdateInDeviceList(driver)
		if err != nil {
			fmt.Printf("Can't update MIDI in device list: %v\n", err)
			return
		}

		for {
			select {
			case <-deviceChannel:
			// NOTE: block on this channel until loop is started
			case <-ctx.Done():
				return
			case <-time.After(1 * time.Second):
				macmidi.PumpRunLoop() // Process pending run loop events
			}
		}
	}()

	mc.WaitUntilDevicesQueried()
}

func GetIns() ([]drivers.In, error) {
	driver, err := driver()
	if err != nil {
		return nil, fmt.Errorf("can't open MIDI driver: %v", err)
	}
	ins, err := driver.Ins()
	if err != nil {
		return nil, fmt.Errorf("can't get MIDI ins: %v", err)
	}
	return ins, nil
}

func Outs() ([]drivers.Out, error) {
	driver, err := macmididrv.New()
	if err != nil {
		return nil, fmt.Errorf("can't open MIDI driver: %v", err)
	}
	return driver.Outs()
}
