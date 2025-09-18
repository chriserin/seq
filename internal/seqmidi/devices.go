package seqmidi

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

type DeviceInfo struct {
	Out      drivers.Out
	Name     string
	Type     string
	IsOpen   bool
	Selected bool
}

func (di DeviceInfo) Matches(name string) bool {
	return strings.Contains(di.Name, name)
}

func (di *DeviceInfo) Open() {
	if !di.IsOpen {
		err := di.Out.Open()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening MIDI output %v: %v\n", di.Name, err)
			return
		}
		di.IsOpen = true
	}
}

func (mc *MidiConnection) UpdateDeviceList(driver *rtmididrv.Driver) error {
	var newDevices []*DeviceInfo

	outs, err := driver.Outs()
	if err != nil {
		return fmt.Errorf("failed to get output ports: %v", err)
	}

	for _, out := range outs {
		// Check if we already have this device
		var foundDevice *DeviceInfo
		for _, currentDevice := range mc.devices {
			if currentDevice.Name == out.String() {
				foundDevice = currentDevice
				break
			}
		}

		if foundDevice == nil {
			newDevice := &DeviceInfo{
				Out:  out,
				Name: out.String(),
				Type: "Output",
			}
			if mc.outportName != "" && newDevice.Matches(mc.outportName) {
				newDevice.Open()
				newDevice.Selected = true
			}
			newDevices = append(newDevices, newDevice)
		}
	}

	mc.devices = append(mc.devices, newDevices...)

	return nil
}

func (mc *MidiConnection) DeviceLoop(ctx context.Context) {

	driver, err := rtmididrv.New()
	if err != nil {
		fmt.Printf("Can't open MIDI driver: %v\n", err)
		return
	}

	err = mc.UpdateDeviceList(driver)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting device list: %v\n", err)
		return
	}

	go func() {
		// NOTE: Setup initially so we don't have to 3 seconds

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(3 * time.Second):
				// go func() {
				err := mc.UpdateDeviceList(driver)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error getting device list: %v\n", err)
					return
				}
			}
		}
	}()
}
