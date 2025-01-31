package overlaykey

import "fmt"

type OverlayPeriodicity struct {
	Shift      uint8
	Interval   uint8
	Width      uint8
	StartCycle uint8
}

func InitOverlayKey(shift, interval uint8) OverlayPeriodicity {
	return OverlayPeriodicity{shift, interval, 0, 0}
}

func (o *OverlayPeriodicity) IncrementShift() {
	if o.Shift < 9 {
		o.Shift++
	}
}

func (o *OverlayPeriodicity) IncrementInterval() {
	if o.Interval < 9 {
		o.Interval++
	}
}

func (o *OverlayPeriodicity) DecrementShift() {
	if o.Shift > 1 {
		o.Shift--
	}
}

func (o *OverlayPeriodicity) DecrementInterval() {
	if o.Interval > 1 {
		o.Interval--
	}
}

var ROOT_OVERLAY OverlayPeriodicity = OverlayPeriodicity{1, 1, 0, 0}

func (op OverlayPeriodicity) MarshalTOML() ([]byte, error) {
	return []byte(fmt.Sprintf("%d/%d", op.Shift, op.Interval)), nil
}

func (op OverlayPeriodicity) WriteKey() string {
	return fmt.Sprintf("OverlayKey-%d/%d", op.Shift, op.Interval)
}

// Sort from most specific to least specific
func Sort(a, b OverlayPeriodicity) int {
	intervalDiff := int(b.Interval) - int(a.Interval)
	shiftDiff := int(b.Shift) - int(a.Shift)
	startDiff := int(b.StartCycle) - int(a.StartCycle)
	widthDiff := int(a.Width) - int(b.Width)

	switch {
	case intervalDiff != 0:
		return int(intervalDiff)
	case shiftDiff != 0:
		return int(shiftDiff)
	case startDiff != 0:
		return int(startDiff)
	case widthDiff != 0:
		return int(widthDiff)
	}
	return 0
}

func (op OverlayPeriodicity) DoesMatch(cycle int) bool {
	if cycle < int(op.StartCycle) {
		return false
	}

	shift, overallInterval := op.normalizeShiftInterval()
	intervalNum := (cycle - shift) / overallInterval
	intervalLevel := int(overallInterval) * intervalNum

	if shift == 0 && cycle < int(overallInterval) {
		return false
	} else if int(shift) == cycle {
		return true
	} else if shift+overallInterval == cycle {
		return true
	} else if int(shift)+intervalLevel <= cycle && shift+intervalLevel+int(op.Width) >= cycle {
		return true
	}
	return false
}

func (op OverlayPeriodicity) normalizeShiftInterval() (int, int) {

	var overallInterval = op.Interval
	var shift = op.Shift
	if op.Interval < op.Shift {
		if op.Shift%op.Interval == 0 {
			overallInterval = op.Shift / op.Interval
			shift = 0
		} else {
			overallInterval = op.Interval * (op.Shift/op.Interval + 1)
		}
	}
	return int(shift), int(overallInterval)
}
