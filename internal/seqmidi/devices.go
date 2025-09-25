package seqmidi

import (
	"fmt"
	"os"
	"strings"

	"gitlab.com/gomidi/midi/v2/drivers"
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

func (mc *MidiConnection) UpdateDeviceList(driver drivers.Driver) error {
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
				foundDevice.Out = out
				foundDevice.IsOpen = false
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
		} else {
			newDevices = append(newDevices, foundDevice)
		}
	}

	mc.devices = newDevices

	return nil
}

func SeqOut() (drivers.Out, error) {
	return OpenVirtualOut(OutputName)
}

func TransmitterOut() (drivers.Out, error) {
	return OpenVirtualOut(OutputName)
}
