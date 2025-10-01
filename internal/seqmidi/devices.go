package seqmidi

import (
	"fmt"
	"os"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"gitlab.com/gomidi/midi/v2/drivers"
)

type OutDeviceInfo struct {
	IsOpen   bool
	Selected bool
	IsDaw    bool
	Out      drivers.Out
	Name     string
	Type     string
}

type InDeviceInfo struct {
	IsOpen        bool
	IsTransmitter bool
	In            drivers.In
	Name          string
	Type          string
}

func (di OutDeviceInfo) Matches(name string) bool {
	return strings.Contains(di.Name, name)
}

func (di *OutDeviceInfo) Open() {
	if !di.IsOpen {
		err := di.Out.Open()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening MIDI output %v: %v\n", di.Name, err)
			return
		}
		di.IsOpen = true
	}
}

func (di InDeviceInfo) Matches(name string) bool {
	return strings.Contains(di.Name, name)
}

func (di *InDeviceInfo) Open() {
	if !di.IsOpen {
		err := di.In.Open()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening MIDI input %v: %v\n", di.Name, err)
			return
		}
		di.IsOpen = true
	}
}

type ReceiverFunc func(msg []byte, milliseconds int32)

func (mc *MidiConnection) ListenToTransmitter(recFunc ReceiverFunc) error {
	mc.ReceiverFunc = recFunc
	for _, device := range mc.inDevices {
		if device.Matches(TransmitterName) {
			if device.In != nil && device.In.IsOpen() {
				if mc.StopFn != nil {
					mc.StopFn()
				}
				stopFn, err := device.In.Listen(recFunc, drivers.ListenConfig{TimeCode: true, ActiveSense: true})
				if err != nil {
					return fault.Wrap(err, fmsg.With("cannot listen to transmitter"))
				}
				mc.StopFn = stopFn
			}
		}
	}
	return nil
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

func (mc *MidiConnection) UpdateOutDeviceList(driver drivers.Driver) error {
	var newDevices []*OutDeviceInfo

	outs, err := driver.Outs()
	if err != nil {
		return fmt.Errorf("failed to get output ports: %v", err)
	}

	for _, out := range outs {
		// Check if we already have this device
		var foundDevice *OutDeviceInfo
		for _, currentDevice := range mc.outDevices {
			if currentDevice.Name == out.String() {
				foundDevice = currentDevice
				foundDevice.Out = out
				foundDevice.IsOpen = false
				if mc.outportName != "" && foundDevice.Matches(mc.outportName) {
					foundDevice.Open()
				}
				break
			}
		}

		if foundDevice == nil {
			newDevice := &OutDeviceInfo{
				Out:  out,
				Name: out.String(),
			}
			if mc.outportName != "" && newDevice.Matches(mc.outportName) {
				newDevice.Open()
				newDevice.Selected = true
			}
			for _, dawName := range dawOutports {
				if strings.Contains(newDevice.Name, dawName) {
					newDevice.Open()
					newDevice.IsDaw = true
				}
			}
			newDevices = append(newDevices, newDevice)
		} else {
			newDevices = append(newDevices, foundDevice)
		}
	}

	mc.outDevices = newDevices
	mc.EnsureConnection()

	return nil
}

func SeqOut() (drivers.Out, error) {
	return OpenVirtualOut(OutputName)
}

func TransmitterOut() (drivers.Out, error) {
	return OpenVirtualOut(TransmitterName)
}
