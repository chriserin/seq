// Package overlays provides layered pattern management for the sequencer.
// It implements a hierarchical overlay system where multiple pattern layers
// can be stacked and combined based on timing cycles, enabling complex
// arrangements with pattern variations, chord progressions, and dynamic
// pattern switching during playback.
package overlays

import (
	"fmt"

	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/overlaykey"
)

type Key = overlaykey.OverlayPeriodicity

type Overlay struct {
	Key       Key
	Below     *Overlay
	Notes     grid.Pattern
	Chords    Chords
	blockers  []grid.GridKey
	PressUp   bool
	PressDown bool
}

func (ol Overlay) String() string {
	return fmt.Sprintf("%v %d", ol.Key, len(ol.Chords))
}

func (ol Overlay) IsFresh() bool {
	return len(ol.Notes) == 0 && len(ol.Chords) == 0
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
		Key:     key,
		Below:   below,
		Notes:   make(grid.Pattern),
		PressUp: overlaykey.ROOT == key,
		Chords:  []*GridChord{},
	}
}

func (ol Overlay) CollectKeys(collection *[]Key) {
	for currentOverlay := &ol; currentOverlay != nil; currentOverlay = currentOverlay.Below {
		(*collection) = append((*collection), currentOverlay.Key)
	}
}

type OverlayChord struct {
	Overlay   *Overlay
	GridChord *GridChord
}

func (oc OverlayChord) HasValue() bool {
	return oc.GridChord != nil
}

func (oc OverlayChord) BelongsTo(anotherOverlay *Overlay) bool {
	return oc.Overlay == anotherOverlay
}

func (ol *Overlay) FindChord(position grid.GridKey) (OverlayChord, bool) {
	var currentOverlay *Overlay
	for currentOverlay = ol; currentOverlay != nil; currentOverlay = currentOverlay.Below {
		gridChord, exists := currentOverlay.Chords.FindChord(position)
		if exists {
			return OverlayChord{currentOverlay, gridChord}, true
		}
	}
	return OverlayChord{}, false
}

func (ol *Overlay) FindChordWithNote(position grid.GridKey) (*GridChord, bool) {
	var currentOverlay *Overlay
	for currentOverlay = ol; currentOverlay != nil; currentOverlay = currentOverlay.Below {
		gridChord, exists := currentOverlay.Chords.FindChordWithNote(position)
		if exists {
			return gridChord, exists
		}
	}
	return nil, false
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

var zeronote grid.Note

func (ol Overlay) CombineActionPattern(pattern *grid.Pattern, keyCycles int) {
	var addFunc = func(overlayPattern grid.Pattern, key Key) bool {
		for gridKey, note := range overlayPattern {
			_, hasNote := (*pattern)[gridKey]
			if !hasNote && (note.Action != grid.ActionNothing || note == zeronote) {
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

func (ol *Overlay) combine(keyCycles int, addFunc AddFunc) {
	previousPressDown := false
	firstMatch := false

	blockedChords := make(map[grid.GridKey]struct{})

	var currentOverlay *Overlay
	for currentOverlay = ol; currentOverlay != nil; currentOverlay = currentOverlay.Below {
		if previousPressDown ||
			(!firstMatch && currentOverlay.Key.DoesMatch(keyCycles)) ||
			(currentOverlay.PressUp && currentOverlay.Key.DoesMatch(keyCycles)) {
			firstMatch = true
			chordPattern := make(grid.Pattern)

			for _, gridChord := range currentOverlay.Chords {
				_, chordAlreadyPlaced := blockedChords[gridChord.Root]
				if !chordAlreadyPlaced {
					gridChord.ArppegiatedPattern(&chordPattern)
					blockedChords[gridChord.Root] = struct{}{}
				}
			}

			for _, gridKey := range currentOverlay.blockers {
				blockedChords[gridKey] = struct{}{}
			}

			addFunc(chordPattern, currentOverlay.Key)
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

func (ol *Overlay) CombineOverlayPattern(pattern *OverlayPattern, keyCycles int) {
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

func (ol *Overlay) FindOverlay(key Key) *Overlay {
	var currentOverlay *Overlay

	for currentOverlay = ol; currentOverlay != nil; currentOverlay = currentOverlay.Below {
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

func (ol *Overlay) MoveNoteTo(gridKey grid.GridKey, note grid.Note) {
	if ol.Below == nil {
		if note != zeronote {
			(*ol).Notes[gridKey] = note
		} else {
			delete((*ol).Notes, gridKey)
		}
	} else {
		_, exists := ol.Below.GetNote(gridKey)
		if exists || note != zeronote {
			(*ol).Notes[gridKey] = note
		} else {
			delete((*ol).Notes, gridKey)
		}
	}
}

func (ol *Overlay) SetNote(gridKey grid.GridKey, note grid.Note) {
	_, exists := (*ol).Notes[gridKey]
	if exists {
		(*ol).Notes[gridKey] = note
		return
	}
	chord, exists := ol.Chords.FindChordWithNote(gridKey)
	if exists {
		chord.SetChordNote(gridKey, note)
	} else {
		(*ol).Notes[gridKey] = note
	}
}

func (ol *Overlay) SetChord(gridChord *GridChord) *GridChord {
	ol.blockers = append(ol.blockers, gridChord.Root)
	newGridChord := *gridChord
	chordRef := &newGridChord
	ol.Chords = append(ol.Chords, chordRef)
	return chordRef
}

func (ol *Overlay) RemoveNote(gridKey grid.GridKey) {
	delete((*ol).Notes, gridKey)
}

func (ol *Overlay) RemoveChord(overlayChord OverlayChord) {
	if overlayChord.Overlay == ol {
		ol.Chords = ol.Chords.Remove(overlayChord.GridChord)
	} else {
		ol.blockers = append(ol.blockers, overlayChord.GridChord.Root)
	}
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

func (gc GridChord) ChordNotes(pattern *grid.Pattern) {
	for i, interval := range gc.Chord.Intervals() {
		beatnote := gc.Notes[i]
		(*pattern)[gc.Key(interval, beatnote)] = beatnote.Note
	}
}

func (gc GridChord) ArppegiatedPattern(pattern *grid.Pattern) {
	for i, interval := range gc.ArppegioIntervals() {
		beatnote := gc.Notes[i]
		(*pattern)[gc.Key(interval, beatnote)] = beatnote.Note
	}
}

func (gc GridChord) Key(interval uint8, beatnote BeatNote) grid.GridKey {
	b := uint8(int(gc.Root.Beat) + beatnote.Beat)
	return grid.GridKey{Line: gc.Root.Line - interval, Beat: b}
}
