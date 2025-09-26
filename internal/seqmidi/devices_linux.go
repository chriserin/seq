//go:build linux

package seqmidi

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

func OpenVirtualOut(name string) (drivers.Out, error) {
	driver, err := rtmididrv.New()
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
				go func() {

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
					err = driver.Close()
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error closing driver: %v\n", err)
						return
					}
				}()
			}
		}
	}()
}

func GetIns() ([]drivers.In, error) {
	driver, err := rtmididrv.New()
	if err != nil {
		return nil, fmt.Errorf("can't open MIDI driver: %v", err)
	}
	ins, err := driver.Ins()
	if err != nil {
		return nil, fmt.Errorf("can't get MIDI ins: %v", err)
	}
	return ins, nil
}
