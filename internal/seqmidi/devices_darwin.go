//go:build darwin

package seqmidi

import "C"
import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/chriserin/seq/internal/macmididrv"
	"github.com/chriserin/seq/internal/macmididrv/imported/macmidi"
	"gitlab.com/gomidi/midi/v2/drivers"
)

const (
	DeviceAdded   = 2
	DeviceRemoved = 3
)

func deviceCallback(messageID int, name string) {
	fmt.Printf("Message %d %s\n", messageID, name)
	switch messageID {
	case DeviceAdded:
		fmt.Printf("Device added\n")
	case DeviceRemoved:
		fmt.Printf("Device removed\n")
	default:
		fmt.Printf("Unknown message %d", messageID)
	}
}

func OpenVirtualOut(name string) (drivers.Out, error) {
	driver, err := macmididrv.New()
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("midi driver error"))
	}
	out, err := driver.OpenVirtualOut(name)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("cannot open virtual out"))
	}
	return out, nil
}

func (mc *MidiConnection) DeviceLoop(ctx context.Context) {
	go func() {
		//NOTE: The client notification will only be called on the same thread as the client inittialization, we have to "pump the run loop" on that thread.
		runtime.LockOSThread()

		// NOTE: Currently this New call does some C singleton initialization and is called here explicitly for that.
		driver, err := macmididrv.New()
		deviceCallback := func(messageID int, name string) {
			mc.UpdateDeviceList(driver)
		}
		macmididrv.SetNotificationCallback(deviceCallback)
		if err != nil {
			fmt.Printf("Can't open MIDI driver: %v\n", err)
			return
		}
		err = mc.UpdateDeviceList(driver)
		if err != nil {
			fmt.Printf("Can't update MIDI device list: %v\n", err)
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(1 * time.Second):
				macmidi.PumpRunLoop() // Process pending run loop events
			}
		}
	}()
}
