package main

import (
	"github.com/chriserin/seq/internal/arrangement"
	"github.com/chriserin/seq/internal/operation"
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

var EmptyStack = UndoStack{}

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
	m.focus = operation.FocusArrangementEditor
	m.showArrangementView = true
	m.arrangement.Focus = true
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

type UndoStateDiff struct {
	stateDiff StateDiff
}

func (usd UndoStateDiff) ApplyUndo(m *model) Location {
	usd.stateDiff.Apply(m)

	if usd.stateDiff.AccentsChanged {
		m.selectionIndicator = operation.SelectAccentDiff
	}
	if usd.stateDiff.LinesChanged {
		m.selectionIndicator = operation.SelectSetupChannel
	}
	if usd.stateDiff.SubdivisionsChanged {
		m.selectionIndicator = operation.SelectTempoSubdivision
	}
	if usd.stateDiff.TempoChanged {
		m.selectionIndicator = operation.SelectTempo
	}
	return Location{ApplyLocation: false}
}
