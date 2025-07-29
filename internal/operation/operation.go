// Package operation defines the types and values used to represent the operational state of the application.
package operation

type SequencerMode uint8

const (
	SeqModeLine SequencerMode = iota
	SeqModeChord
)

type Focus uint8

// NOTE: Focus is necessary because the selection can be focused for a part operation
// at the same time the arrangement editor is focused.  `FocusOverlayKey` could be switched
// out for a selection indicator, but the pattern is to focus components that have their
// own key event responses.
const (
	FocusGrid Focus = iota
	FocusOverlayKey
	FocusArrangementEditor
)

type Selection uint8

const (
	SelectNothing Selection = iota
	// Definition Change
	SelectTempo
	SelectTempoSubdivision
	SelectSetupChannel
	SelectSetupMessageType
	SelectSetupValue
	SelectAccentDiff
	SelectAccentTarget
	SelectAccentStart

	// Part Change
	SelectBeats

	// Arrangement Change
	SelectPart
	SelectChangePart
	SelectRenamePart
	SelectCycles
	SelectStartBeats
	SelectStartCycles

	// Note Change
	SelectRatchets
	SelectRatchetSpan
	SelectSpecificValue

	// Program Level Operation
	SelectConfirmNew
	SelectConfirmQuit
	SelectConfirmReload
	SelectFileName
	SelectError
)

type PatternMode uint8

const (
	PatternFill PatternMode = iota
	PatternAccent
	PatternGate
	PatternWait
	PatternRatchet
)
