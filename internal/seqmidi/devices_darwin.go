//go:build darwin

package seqmidi

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/chriserin/seq/internal/macmididrv"
	"gitlab.com/gomidi/midi/v2/drivers"
)

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

		driver, err := macmididrv.New()
		if err != nil {
			fmt.Printf("Can't open MIDI driver: %v\n", err)
			return
		}
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(3 * time.Second):
				// go func() {
				go func() {

					err = mc.UpdateDeviceList(driver)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error getting device list: %v\n", err)
						return
					}
				}()
			}
		}
	}()
}
