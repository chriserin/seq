package main

type Undoable interface {
	ApplyUndo(m *model) (overlayKey, gridKey)
}

type UndoStack struct {
	undo Undoable
	redo Undoable
	next *UndoStack
	id   int
}

var NIL_STACK = UndoStack{}

type UndoBeats struct {
	beats  uint8
	partId int
}

func (ub UndoBeats) ApplyUndo(m *model) {
	(*m.definition.parts)[ub.partId].Beats = ub.beats
}

type UndoTempo struct {
	tempo int
}

func (ukl UndoTempo) ApplyUndo(m *model) {
	m.definition.tempo = ukl.tempo
}

type UndoSubdivisions struct {
	subdivisions int
}

func (ukl UndoSubdivisions) ApplyUndo(m *model) {
	m.definition.subdivisions = ukl.subdivisions
}

type UndoGridNote struct {
	overlayKey
	cursorPosition gridKey
	gridNote       GridNote
}

func (ugn UndoGridNote) ApplyUndo(m *model) (overlayKey, gridKey) {
	m.EnsureOverlayWithKey(ugn.overlayKey)
	overlay := m.CurrentPart().Overlays.FindOverlay(ugn.overlayKey)
	overlay.SetNote(ugn.gridNote.gridKey, ugn.gridNote.note)
	return ugn.overlayKey, ugn.gridNote.gridKey
}

type UndoLineGridNotes struct {
	overlayKey
	cursorPosition gridKey
	line           uint8
	gridNotes      []GridNote
}

func (ulgn UndoLineGridNotes) ApplyUndo(m *model) (overlayKey, gridKey) {
	m.EnsureOverlayWithKey(ulgn.overlayKey)
	overlay := m.CurrentPart().Overlays.FindOverlay(ulgn.overlayKey)
	for i := range m.CurrentPart().Beats {
		overlay.RemoveNote(GK(ulgn.line, i))
	}
	for _, gridNote := range ulgn.gridNotes {
		overlay.SetNote(gridNote.gridKey, gridNote.note)
	}
	return ulgn.overlayKey, ulgn.cursorPosition
}

type UndoBounds struct {
	overlayKey
	cursorPosition gridKey
	bounds         Bounds
	gridNotes      []GridNote
}

func (uvs UndoBounds) ApplyUndo(m *model) (overlayKey, gridKey) {
	m.EnsureOverlayWithKey(uvs.overlayKey)
	overlay := m.CurrentPart().Overlays.FindOverlay(uvs.overlayKey)
	for _, k := range uvs.bounds.GridKeys() {
		overlay.RemoveNote(k)
	}
	for _, gridNote := range uvs.gridNotes {
		overlay.SetNote(gridNote.gridKey, gridNote.note)
	}
	return uvs.overlayKey, uvs.cursorPosition
}

type UndoGridNotes struct {
	overlayKey
	gridNotes []GridNote
}

func (ugn UndoGridNotes) ApplyUndo(m *model) (overlayKey, gridKey) {
	m.EnsureOverlayWithKey(ugn.overlayKey)
	overlay := m.CurrentPart().Overlays.FindOverlay(ugn.overlayKey)
	for _, gridNote := range ugn.gridNotes {

		overlay.SetNote(gridNote.gridKey, gridNote.note)
	}
	return ugn.overlayKey, ugn.gridNotes[0].gridKey
}

type UndoToNothing struct {
	overlayKey overlayKey
	location   gridKey
}

func (utn UndoToNothing) ApplyUndo(m *model) (overlayKey, gridKey) {
	overlay := m.CurrentPart().Overlays.FindOverlay(utn.overlayKey)
	overlay.RemoveNote(utn.location)
	return utn.overlayKey, utn.location
}

type UndoLineToNothing struct {
	overlayKey     overlayKey
	cursorPosition gridKey
	line           uint8
}

func (ultn UndoLineToNothing) ApplyUndo(m *model) (overlayKey, gridKey) {
	overlay := m.CurrentPart().Overlays.FindOverlay(ultn.overlayKey)
	for i := range m.CurrentPart().Beats {
		overlay.RemoveNote(GK(ultn.line, i))
	}

	return ultn.overlayKey, ultn.cursorPosition
}

type UndoNewOverlay struct {
	overlayKey     overlayKey
	cursorPosition gridKey
	partId         int
}

func (uno UndoNewOverlay) ApplyUndo(m *model) (overlayKey, gridKey) {
	newOverlay := m.CurrentPart().Overlays.Remove(uno.overlayKey)
	(*m.definition.parts)[uno.partId].Overlays = newOverlay
	return uno.overlayKey, uno.cursorPosition
}
