package overlays

import (
	"github.com/chriserin/seq/internal/grid"
	overlaykey "github.com/chriserin/seq/internal/overlaykey"
)

// type gridKey = int
// type note = int

type Key = overlaykey.OverlayPeriodicity

type Overlay struct {
	Key       Key
	Below     *Overlay
	Notes     grid.Pattern
	PressUp   bool
	PressDown bool
}

func (ol Overlay) Add(key Key) *Overlay {
	aboveComparison := overlaykey.Compare(key, ol.Key)
	var belowComparison int
	if ol.Below != nil {
		belowComparison = overlaykey.Compare(key, (*ol.Below).Key)
	} else {
		belowComparison = -1
	}

	if aboveComparison > 0 && belowComparison > 0 {
		ol.Below.Add(key)
	} else if aboveComparison > 0 && belowComparison < 0 {
		newOverlay := InitOverlay(key, ol.Below)
		ol.Below = newOverlay
	} else if aboveComparison < 0 {
		return InitOverlay(key, &ol)
	} else {
		panic("NOT AN OPTION")
	}

	return &ol
}

func (ol Overlay) Remove(key Key) *Overlay {
	if ol.Key == key {
		if ol.Below != nil {
			return ol.Below
		} else {
			return nil
		}
	} else {
		olBelow := (*ol.Below).Remove(key)
		(&ol).Below = olBelow
		return &ol
	}
}

func InitOverlay(key Key, below *Overlay) *Overlay {
	return &Overlay{
		Key:   key,
		Below: below,
		Notes: make(grid.Pattern),
	}
}

func (ol Overlay) CollectKeys(collection *[]Key) {
	for currentOverlay := &ol; currentOverlay != nil; currentOverlay = currentOverlay.Below {
		(*collection) = append((*collection), currentOverlay.Key)
	}
}

func (ol Overlay) CombinePattern(pattern *grid.Pattern, keyCycles int) {
	var addFunc = func(overlayPattern grid.Pattern, key Key) bool {
		for gridKey, note := range overlayPattern {
			_, hasNote := (*pattern)[gridKey]
			if !hasNote {
				(*pattern)[gridKey] = note
			}
		}
		return true
	}
	ol.combine(keyCycles, addFunc)
}

func (ol Overlay) CombineActionPattern(pattern *grid.Pattern, keyCycles int) {
	var addFunc = func(overlayPattern grid.Pattern, key Key) bool {
		for gridKey, note := range overlayPattern {
			_, hasNote := (*pattern)[gridKey]
			if !hasNote && note.Action != grid.ACTION_NOTHING {
				(*pattern)[gridKey] = note
			}
		}
		return true
	}
	ol.combine(keyCycles, addFunc)
}

func (ol Overlay) GetMatchingOverlayKeys(keys *[]Key, keyCycles int) {
	var addFunc = func(pattern grid.Pattern, key Key) bool {
		(*keys) = append((*keys), key)
		return true
	}
	ol.combine(keyCycles, addFunc)
}

type AddFunc = func(grid.Pattern, Key) bool

func (ol Overlay) combine(keyCycles int, addFunc AddFunc) {
	previousPressDown := false
	firstMatch := false
	for currentOverlay := &ol; currentOverlay != nil; currentOverlay = currentOverlay.Below {

		if previousPressDown ||
			(!firstMatch && currentOverlay.Key.DoesMatch(keyCycles)) ||
			(currentOverlay.PressUp && currentOverlay.Key.DoesMatch(keyCycles)) {
			firstMatch = true
			if !addFunc(currentOverlay.Notes, currentOverlay.Key) {
				break
			}

			previousPressDown = currentOverlay.PressDown
		}
	}
}

func (ol Overlay) CurrentBeatOverlayPattern(pattern *grid.Pattern, keyCycles int, beats []grid.GridKey) {
	var addFunc = func(overlayPattern grid.Pattern, currentKey Key) bool {
		for _, gridKey := range beats {
			_, hasNote := (*pattern)[gridKey]
			if !hasNote {
				note, hasNote := overlayPattern[gridKey]
				if hasNote {
					(*pattern)[gridKey] = note
				}
			}
		}
		return len(*pattern) < len(beats)
	}
	ol.combine(keyCycles, addFunc)
}

type OverlayNote struct {
	OverlayKey     overlaykey.OverlayPeriodicity
	Note           grid.Note
	HighestOverlay bool
}
type OverlayPattern map[grid.GridKey]OverlayNote

func (ol Overlay) CombineOverlayPattern(pattern *OverlayPattern, keyCycles int) {
	var firstMatch = true
	var addFunc = func(overlayPattern grid.Pattern, currentKey Key) bool {
		for gridKey, note := range overlayPattern {
			_, hasNote := (*pattern)[gridKey]
			if !hasNote {
				(*pattern)[gridKey] = OverlayNote{OverlayKey: currentKey, Note: note, HighestOverlay: firstMatch}
			}
		}

		firstMatch = false

		return true
	}
	ol.combine(keyCycles, addFunc)
}

func (ol Overlay) FindOverlay(key Key) *Overlay {
	for currentOverlay := &ol; currentOverlay != nil; currentOverlay = currentOverlay.Below {
		if currentOverlay.Key == key {
			return currentOverlay
		}
	}
	return nil
}

func (ol *Overlay) FindAboveOverlay(key Key) *Overlay {
	previousOverlay := ol
	for currentOverlay := ol; currentOverlay != nil; currentOverlay = currentOverlay.Below {
		if currentOverlay.Key == key {
			return previousOverlay
		} else {
			previousOverlay = currentOverlay
		}
	}
	return nil
}

func (ol Overlay) HighestMatchingOverlay(keyCycle int) *Overlay {
	for currentOverlay := &ol; currentOverlay != nil; currentOverlay = currentOverlay.Below {
		if currentOverlay.Key.DoesMatch(keyCycle) {
			return currentOverlay
		}
	}
	return nil
}

func (ol Overlay) GetNote(gridKey grid.GridKey) (grid.Note, bool) {
	for currentOverlay := &ol; currentOverlay != nil; currentOverlay = currentOverlay.Below {
		note, hasNote := currentOverlay.Notes[gridKey]
		if hasNote {
			return note, true
		}
	}
	return grid.Note{}, false
}

func (ol *Overlay) SetNote(gridKey grid.GridKey, note grid.Note) {
	(*ol).Notes[gridKey] = note
}

func (ol *Overlay) RemoveNote(gridKey grid.GridKey) {
	delete((*ol).Notes, gridKey)
}

func (ol Overlay) GridKeysInUse(gridKey grid.GridKey) {
}

func (ol *Overlay) ToggleOverlayStackOptions() {
	if !ol.PressDown && !ol.PressUp {
		ol.PressUp = true
		ol.PressDown = false
	} else if ol.PressUp {
		ol.PressUp = false
		ol.PressDown = true
	} else {
		ol.PressUp = false
		ol.PressDown = false
	}
}
