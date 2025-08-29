// Package overlaykey provides functionality for managing overlay periodicity patterns
// in the sequencer. It defines timing patterns with shift, interval, width, and
// start cycle parameters that determine when overlay events should trigger during
// sequence playback, enabling complex arrangement variations.
package overlaykey

type OverlayPeriodicity struct {
	Shift      uint8
	Interval   uint8
	Width      uint8
	StartCycle uint8
}

func InitOverlayKey(shift, interval uint8) OverlayPeriodicity {
	return OverlayPeriodicity{shift, interval, 1, 0}
}

func (op *OverlayPeriodicity) IncrementShift() {
	if op.Shift < 99 {
		op.Shift++
	}
}

func (op *OverlayPeriodicity) IncrementInterval() {
	if op.Interval < 99 {
		op.Interval++
	}
}

func (op *OverlayPeriodicity) IncrementWidth() {
	if op.Width < 99 {
		op.Width++
	}
}

func (op *OverlayPeriodicity) IncrementStartCycle() {
	if op.StartCycle < 99 {
		op.StartCycle++
	}
}

func (op *OverlayPeriodicity) DecrementShift() {
	if op.Shift > 1 {
		op.Shift--
	}
}

func (op *OverlayPeriodicity) DecrementWidth() {
	if op.Width > 1 {
		op.Width--
	}
}

func (op *OverlayPeriodicity) DecrementStartCycle() {
	if op.StartCycle > 1 {
		op.StartCycle--
	}
}

func (op *OverlayPeriodicity) DecrementInterval() {
	if op.Interval > 1 {
		op.Interval--
	}
}

var ROOT OverlayPeriodicity = OverlayPeriodicity{1, 1, 1, 0}

// Compare from most specific to least specific
func Compare(a, b OverlayPeriodicity) int {
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
	} else if int(shift)+intervalLevel <= cycle && shift+intervalLevel+(int(op.Width)-1) >= cycle {
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

func (op OverlayPeriodicity) GetMinimumKeyCycle() int {
	for i := 1; i < 100; i++ {
		if op.DoesMatch(i) {
			return i
		}
	}
	return 100
}
