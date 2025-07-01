package main

import (
	"github.com/chriserin/seq/internal/arrangement"
	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/overlays"
)

type Undoable interface {
	ApplyUndo(m *model) Location
}

type Location struct {
	OverlayKey    overlayKey
	GridKey       gridKey
	ApplyLocation bool
}

type UndoStack struct {
	undo Undoable
	redo Undoable
	next *UndoStack
	id   int
}

var NilStack = UndoStack{}

type UndoBeats struct {
	beats     uint8
	ArrCursor arrangement.ArrCursor
}

func (ub UndoBeats) ApplyUndo(m *model) Location {
	m.arrangement.Cursor = ub.ArrCursor
	partID := m.CurrentPartID()
	(*m.definition.parts)[partID].Beats = ub.beats
	return Location{ApplyLocation: false}
}

type UndoTempo struct {
	tempo int
}

func (ukl UndoTempo) ApplyUndo(m *model) Location {
	m.definition.tempo = ukl.tempo
	return Location{ApplyLocation: false}
}

type UndoSubdivisions struct {
	subdivisions int
}

func (ukl UndoSubdivisions) ApplyUndo(m *model) Location {
	m.definition.subdivisions = ukl.subdivisions
	return Location{ApplyLocation: false}
}

type UndoGridNote struct {
	overlayKey
	cursorPosition gridKey
	gridNote       GridNote
	ArrCursor      arrangement.ArrCursor
}

func (ugn UndoGridNote) ApplyUndo(m *model) Location {
	m.arrangement.Cursor = ugn.ArrCursor
	m.EnsureOverlayWithKey(ugn.overlayKey)
	overlay := m.CurrentPart().Overlays.FindOverlay(ugn.overlayKey)
	overlay.SetNote(ugn.gridNote.gridKey, ugn.gridNote.note)
	return Location{ugn.overlayKey, ugn.gridNote.gridKey, true}
}

type UndoLineGridNotes struct {
	overlayKey
	cursorPosition gridKey
	line           uint8
	gridNotes      []GridNote
	ArrCursor      arrangement.ArrCursor
}

func (ulgn UndoLineGridNotes) ApplyUndo(m *model) Location {
	m.arrangement.Cursor = ulgn.ArrCursor
	m.EnsureOverlayWithKey(ulgn.overlayKey)
	overlay := m.CurrentPart().Overlays.FindOverlay(ulgn.overlayKey)
	for i := range m.CurrentPart().Beats {
		overlay.RemoveNote(GK(ulgn.line, i))
	}
	for _, gridNote := range ulgn.gridNotes {
		overlay.SetNote(gridNote.gridKey, gridNote.note)
	}
	return Location{ulgn.overlayKey, ulgn.cursorPosition, true}
}

type UndoBounds struct {
	overlayKey
	cursorPosition gridKey
	bounds         grid.Bounds
	gridNotes      []GridNote
	ArrCursor      arrangement.ArrCursor
}

func (ub UndoBounds) ApplyUndo(m *model) Location {
	m.arrangement.Cursor = ub.ArrCursor
	m.EnsureOverlayWithKey(ub.overlayKey)
	overlay := m.CurrentPart().Overlays.FindOverlay(ub.overlayKey)
	for _, k := range ub.bounds.GridKeys() {
		overlay.RemoveNote(k)
	}
	for _, gridNote := range ub.gridNotes {
		overlay.SetNote(gridNote.gridKey, gridNote.note)
	}
	return Location{ub.overlayKey, ub.cursorPosition, true}
}

type UndoGridNotes struct {
	overlayKey
	gridNotes []GridNote
	ArrCursor arrangement.ArrCursor
}

func (ugn UndoGridNotes) ApplyUndo(m *model) Location {
	m.arrangement.Cursor = ugn.ArrCursor
	m.EnsureOverlayWithKey(ugn.overlayKey)
	overlay := m.CurrentPart().Overlays.FindOverlay(ugn.overlayKey)
	for _, gridNote := range ugn.gridNotes {

		overlay.SetNote(gridNote.gridKey, gridNote.note)
	}
	return Location{ugn.overlayKey, ugn.gridNotes[0].gridKey, true}
}

type UndoToNothing struct {
	overlayKey overlayKey
	location   gridKey
	ArrCursor  arrangement.ArrCursor
}

func (utn UndoToNothing) ApplyUndo(m *model) Location {
	m.arrangement.Cursor = utn.ArrCursor
	overlay := m.CurrentPart().Overlays.FindOverlay(utn.overlayKey)
	overlay.RemoveNote(utn.location)
	return Location{utn.overlayKey, utn.location, true}
}

type UndoLineToNothing struct {
	overlayKey     overlayKey
	cursorPosition gridKey
	line           uint8
	ArrCursor      arrangement.ArrCursor
}

func (ultn UndoLineToNothing) ApplyUndo(m *model) Location {
	m.arrangement.Cursor = ultn.ArrCursor
	overlay := m.CurrentPart().Overlays.FindOverlay(ultn.overlayKey)
	for i := range m.CurrentPart().Beats {
		overlay.RemoveNote(GK(ultn.line, i))
	}

	return Location{ultn.overlayKey, ultn.cursorPosition, true}
}

type UndoNewOverlay struct {
	overlayKey     overlayKey
	cursorPosition gridKey
	ArrCursor      arrangement.ArrCursor
}

func (uno UndoNewOverlay) ApplyUndo(m *model) Location {
	m.arrangement.Cursor = uno.ArrCursor
	currentPartID := m.CurrentPartID()
	newOverlay := m.CurrentPart().Overlays.Remove(uno.overlayKey)
	(*m.definition.parts)[currentPartID].Overlays = newOverlay
	return Location{uno.overlayKey, uno.cursorPosition, true}
}

type UndoArrangement struct {
	arrUndo arrangement.Undoable
}

func (ua UndoArrangement) ApplyUndo(m *model) Location {
	m.arrangement.ApplyArrUndo(ua.arrUndo)
	return Location{ApplyLocation: false}
}

type UndoOverlayDiff struct {
	overlayKey     overlayKey
	cursorPosition gridKey
	ArrCursor      arrangement.ArrCursor
	overlayDiff    overlays.OverlayDiff
}

func (uod UndoOverlayDiff) ApplyUndo(m *model) Location {
	m.arrangement.Cursor = uod.ArrCursor
	overlay := m.CurrentPart().Overlays.FindOverlay(uod.overlayKey)
	if overlay == nil {
		m.currentOverlay = m.CurrentPart().Overlays
	} else {
		m.currentOverlay = overlay
	}
	uod.overlayDiff.Apply(m.currentOverlay)
	if len(uod.overlayDiff.RemovedChords) > 0 {
		m.UnsetActiveChord()
	}
	return Location{uod.overlayKey, uod.cursorPosition, true}
}
