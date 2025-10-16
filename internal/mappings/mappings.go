// Package mappings provides keyboard input handling and command mapping functionality
// for the sequencer application. It manages key combinations, processes user input,
// and maps keyboard shortcuts to sequencer commands based on the current mode
// (trigger, polyphony, pattern mode, chord mode, etc.).
package mappings

import (
	"maps"
	"slices"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chriserin/seq/internal/operation"
)

var Keycombo = make([]tea.KeyMsg, 0, 3)
var timer *time.Timer
var mutex = sync.Mutex{}

func KeycomboView() string {
	var buf strings.Builder
	for _, msg := range Keycombo {
		buf.WriteString(msg.String())
	}
	return buf.String()
}

type Command int
type Mapping struct {
	Command   Command
	LastValue string
}

const (
	HoldingKeys Command = iota
	Quit
	CursorUp
	CursorDown
	CursorLeft
	CursorRight
	CursorLineStart
	CursorLineEnd
	CursorLastLine
	CursorFirstLine
	Escape
	PlayStop
	PlayPart
	PlayLoop
	PlayOverlayLoop
	PlayRecord
	PlayAlong
	TempoInputSwitch
	OverlayInputSwitch
	ModifyKeyInputSwitch
	SetupInputSwitch
	AccentInputSwitch
	RatchetInputSwitch
	BeatInputSwitch
	CyclesInputSwitch
	ToggleArrangementView
	Increase
	Decrease
	Enter
	ToggleGateMode
	ToggleGateNoteMode
	ToggleWaitMode
	ToggleWaitNoteMode
	ToggleAccentMode
	ToggleAccentNoteMode
	ToggleRatchetMode
	ToggleRatchetNoteMode
	NextOverlay
	PrevOverlay
	Save
	SaveAs
	Undo
	Redo
	New
	ToggleVisualMode
	ToggleVisualLineMode
	TogglePlayEdit
	ToggleHideLines
	NewLine
	NewSectionAfter
	NewSectionBefore
	ChangePart
	NextSection
	PrevSection
	NextTheme
	PrevTheme
	Yank
	Mute
	MuteAll
	UnMuteAll
	Solo
	NoteAdd
	NoteRemove
	OverlayNoteRemove
	AccentIncrease
	AccentDecrease
	GateIncrease
	GateDecrease
	GateBigIncrease
	GateBigDecrease
	WaitIncrease
	WaitDecrease
	ClearLine
	ClearOverlay
	ClearAllOverlays
	RemoveOverlay
	RatchetIncrease
	RatchetDecrease
	ActionAddLineReset
	ActionAddLineResetAll
	ActionAddLineReverse
	ActionAddSkipBeat
	ActionAddSkipBeatAll
	ActionAddLineBounce
	ActionAddLineBounceAll
	ActionAddLineDelay
	ActionAddSpecificValue
	SelectKeyLine
	OverlayStackToggle
	NumberPattern
	RotateRight
	RotateLeft
	RotateUp
	RotateDown
	Paste
	MajorTriad
	MinorTriad
	AugmentedTriad
	DiminishedTriad
	MinorSecond
	MajorSecond
	MinorThird
	MajorThird
	PerfectFourth
	AugFifth
	DimFifth
	PerfectFifth
	MajorSixth
	MinorSeventh
	MajorSeventh
	Octave
	MinorNinth
	MajorNinth
	IncreaseInversions
	DecreaseInversions
	ToggleChordMode
	ToggleMonoMode
	NextArpeggio
	PrevArpeggio
	NextDouble
	PrevDouble
	OmitRoot
	OmitSecond
	OmitThird
	OmitFourth
	OmitFifth
	OmitSixth
	OmitSeventh
	OmitOctave
	OmitNinth
	RemoveChord
	ConvertToNotes
	ReloadFile
	OverlayKeyMessage
	ArrKeyMessage
	TextInputMessage
	ConfirmOverlayKey
	ConfirmRenamePart
	ConfirmFileName
	ConfirmSelectPart
	ConfirmChangePart
	ConfirmConfirmNew
	ConfirmConfirmReload
	ConfirmConfirmQuit
	MidiPanic
	IncreaseAllChannels
	DecreaseAllChannels
	IncreaseAllNote
	DecreaseAllNote
	ToggleTransmitting
	PurposePanic
)

type mappingKey [3]string
type registry map[OperationKey]Command

// KeysForCommand Function that gets keys for mappings by looking at the looping
// through the registry and returning the keys for the given command
// If the command is not found, it returns an empty string slice.
// This is used to get the keys for a given command in the mappings.
func KeysForCommand(command Command) []string {
	var keys []string
	for key, cmd := range allCommands() {
		if cmd == command {
			for _, k := range key.key {
				if k != "" {
					keys = append(keys, k)
				}
			}
			break
		}
	}
	return keys
}

func allCommands() registry {
	// Combine all mappings into a single registry
	all := make(registry)
	for _, m := range []registry{mappings} {
		maps.Copy(all, m)
	}
	return all
}

type OperationKey struct {
	focus       operation.Focus
	selection   operation.Selection
	mode        operation.SequencerMode
	patternMode operation.PatternMode
	key         mappingKey
}

var mappings = registry{
	OperationKey{key: k("b", "p")}:                                                                                               MidiPanic,
	OperationKey{focus: operation.FocusAny, key: k(" ")}:                                                                         PlayStop,
	OperationKey{focus: operation.FocusAny, key: k("'", " ")}:                                                                    PlayOverlayLoop,
	OperationKey{focus: operation.FocusAny, key: k(":", " ")}:                                                                    PlayRecord,
	OperationKey{focus: operation.FocusAny, key: k(";", " ")}:                                                                    PlayAlong,
	OperationKey{focus: operation.FocusGrid, key: k("+")}:                                                                        Increase,
	OperationKey{focus: operation.FocusGrid, key: k("=")}:                                                                        Increase,
	OperationKey{focus: operation.FocusGrid, key: k("-")}:                                                                        Decrease,
	OperationKey{selection: operation.SelectTempo, key: k("+")}:                                                                  Increase,
	OperationKey{selection: operation.SelectTempo, key: k("=")}:                                                                  Increase,
	OperationKey{selection: operation.SelectTempo, key: k("-")}:                                                                  Decrease,
	OperationKey{selection: operation.SelectTempoSubdivision, key: k("+")}:                                                       Increase,
	OperationKey{selection: operation.SelectTempoSubdivision, key: k("=")}:                                                       Increase,
	OperationKey{selection: operation.SelectTempoSubdivision, key: k("-")}:                                                       Decrease,
	OperationKey{selection: operation.SelectPart, key: k("+")}:                                                                   Increase,
	OperationKey{selection: operation.SelectPart, key: k("=")}:                                                                   Increase,
	OperationKey{selection: operation.SelectPart, key: k("-")}:                                                                   Decrease,
	OperationKey{selection: operation.SelectChangePart, key: k("+")}:                                                             Increase,
	OperationKey{selection: operation.SelectChangePart, key: k("=")}:                                                             Increase,
	OperationKey{selection: operation.SelectChangePart, key: k("-")}:                                                             Decrease,
	OperationKey{focus: operation.FocusGrid, key: k("<")}:                                                                        CursorLineStart,
	OperationKey{focus: operation.FocusGrid, key: k(">")}:                                                                        CursorLineEnd,
	OperationKey{focus: operation.FocusGrid, key: k("b", "l")}:                                                                   CursorLastLine,
	OperationKey{focus: operation.FocusGrid, key: k("b", "f")}:                                                                   CursorFirstLine,
	OperationKey{focus: operation.FocusGrid, key: k("b", "h")}:                                                                   ToggleHideLines,
	OperationKey{focus: operation.FocusGrid, key: k("b", "t")}:                                                                   ToggleTransmitting,
	OperationKey{focus: operation.FocusGrid, key: k("A")}:                                                                        AccentIncrease,
	OperationKey{focus: operation.FocusGrid, key: k("C")}:                                                                        ClearOverlay,
	OperationKey{focus: operation.FocusGrid, key: k("b", "C")}:                                                                   ClearAllOverlays,
	OperationKey{focus: operation.FocusGrid, key: k("D")}:                                                                        RemoveOverlay,
	OperationKey{focus: operation.FocusGrid, key: k("G")}:                                                                        GateIncrease,
	OperationKey{focus: operation.FocusGrid, key: k("E")}:                                                                        GateBigIncrease,
	OperationKey{focus: operation.FocusGrid, key: k("J")}:                                                                        RotateDown,
	OperationKey{focus: operation.FocusGrid, key: k("K")}:                                                                        RotateUp,
	OperationKey{focus: operation.FocusGrid, key: k("H")}:                                                                        RotateLeft,
	OperationKey{focus: operation.FocusGrid, key: k("L")}:                                                                        RotateRight,
	OperationKey{focus: operation.FocusGrid, key: k("Y")}:                                                                        SelectKeyLine,
	OperationKey{focus: operation.FocusGrid, key: k("M")}:                                                                        Solo,
	OperationKey{focus: operation.FocusGrid, key: k("R")}:                                                                        RatchetIncrease,
	OperationKey{focus: operation.FocusAny, key: k("U")}:                                                                         Redo,
	OperationKey{focus: operation.FocusGrid, key: k("W")}:                                                                        WaitIncrease,
	OperationKey{focus: operation.FocusGrid, key: k("[", "c")}:                                                                   PrevTheme,
	OperationKey{focus: operation.FocusAny, key: k("[", "s")}:                                                                    PrevSection,
	OperationKey{focus: operation.FocusGrid, key: k("]", "c")}:                                                                   NextTheme,
	OperationKey{focus: operation.FocusAny, key: k("]", "s")}:                                                                    NextSection,
	OperationKey{focus: operation.FocusGrid, key: k("g")}:                                                                        GateDecrease,
	OperationKey{focus: operation.FocusGrid, key: k("e")}:                                                                        GateBigDecrease,
	OperationKey{focus: operation.FocusAny, key: k("alt+ ")}:                                                                     PlayLoop,
	OperationKey{focus: operation.FocusAny, key: k("ctrl+@")}:                                                                    PlayPart,
	OperationKey{focus: operation.FocusAny, key: k("ctrl+]")}:                                                                    NewSectionAfter,
	OperationKey{focus: operation.FocusGrid, key: k("ctrl+b")}:                                                                   BeatInputSwitch,
	OperationKey{focus: operation.FocusGrid, key: k("ctrl+k")}:                                                                   CyclesInputSwitch,
	OperationKey{focus: operation.FocusAny, key: k("ctrl+c")}:                                                                    ChangePart,
	OperationKey{focus: operation.FocusGrid, key: k("ctrl+e")}:                                                                   AccentInputSwitch,
	OperationKey{focus: operation.FocusGrid, key: k("ctrl+a")}:                                                                   ToggleArrangementView,
	OperationKey{focus: operation.FocusArrangementEditor, key: k("ctrl+a")}:                                                      ToggleArrangementView,
	OperationKey{focus: operation.FocusAny, key: k("ctrl+l")}:                                                                    NewLine,
	OperationKey{focus: operation.FocusAny, key: k("ctrl+n")}:                                                                    New,
	OperationKey{focus: operation.FocusGrid, key: k("ctrl+o")}:                                                                   OverlayInputSwitch,
	OperationKey{focus: operation.FocusGrid, key: k("ctrl+x")}:                                                                   ModifyKeyInputSwitch,
	OperationKey{focus: operation.FocusOverlayKey, key: k("ctrl+o")}:                                                             OverlayInputSwitch,
	OperationKey{focus: operation.FocusAny, key: k("ctrl+p")}:                                                                    NewSectionBefore,
	OperationKey{focus: operation.FocusGrid, key: k("ctrl+d")}:                                                                   SetupInputSwitch,
	OperationKey{focus: operation.FocusAny, key: k("ctrl+t")}:                                                                    TempoInputSwitch,
	OperationKey{focus: operation.FocusAny, key: k("ctrl+u")}:                                                                    OverlayStackToggle,
	OperationKey{focus: operation.FocusAny, key: k("ctrl+s")}:                                                                    Save,
	OperationKey{focus: operation.FocusAny, key: k("ctrl+w")}:                                                                    SaveAs,
	OperationKey{focus: operation.FocusGrid, key: k("ctrl+y")}:                                                                   RatchetInputSwitch,
	OperationKey{focus: operation.FocusGrid, key: k("a")}:                                                                        AccentDecrease,
	OperationKey{focus: operation.FocusGrid, key: k("c")}:                                                                        ClearLine,
	OperationKey{focus: operation.FocusGrid, key: k("d")}:                                                                        NoteRemove,
	OperationKey{focus: operation.FocusGrid, key: k("b", "e")}:                                                                   TogglePlayEdit,
	OperationKey{focus: operation.FocusGrid, key: k("f")}:                                                                        NoteAdd,
	OperationKey{focus: operation.FocusGrid, key: k("b", "r")}:                                                                   ReloadFile,
	OperationKey{focus: operation.FocusGrid, key: k("b", "v")}:                                                                   ActionAddSpecificValue,
	OperationKey{focus: operation.FocusGrid, key: k("h")}:                                                                        CursorLeft,
	OperationKey{focus: operation.FocusGrid, key: k("j")}:                                                                        CursorDown,
	OperationKey{focus: operation.FocusGrid, key: k("k")}:                                                                        CursorUp,
	OperationKey{focus: operation.FocusGrid, key: k("l")}:                                                                        CursorRight,
	OperationKey{focus: operation.FocusGrid, key: k("m")}:                                                                        Mute,
	OperationKey{focus: operation.FocusGrid, key: k("b", "m")}:                                                                   MuteAll,
	OperationKey{focus: operation.FocusGrid, key: k("b", "M")}:                                                                   UnMuteAll,
	OperationKey{focus: operation.FocusGrid, key: k("o")}:                                                                        ToggleChordMode,
	OperationKey{focus: operation.FocusGrid, key: k("O")}:                                                                        ToggleMonoMode,
	OperationKey{focus: operation.FocusGrid, key: k("n", "a")}:                                                                   ToggleAccentMode,
	OperationKey{focus: operation.FocusGrid, key: k("n", "A")}:                                                                   ToggleAccentNoteMode,
	OperationKey{focus: operation.FocusGrid, key: k("n", "w")}:                                                                   ToggleWaitMode,
	OperationKey{focus: operation.FocusGrid, key: k("n", "W")}:                                                                   ToggleWaitNoteMode,
	OperationKey{focus: operation.FocusGrid, key: k("n", "g")}:                                                                   ToggleGateMode,
	OperationKey{focus: operation.FocusGrid, key: k("n", "G")}:                                                                   ToggleGateNoteMode,
	OperationKey{focus: operation.FocusGrid, key: k("n", "r")}:                                                                   ToggleRatchetMode,
	OperationKey{focus: operation.FocusGrid, key: k("n", "R")}:                                                                   ToggleRatchetNoteMode,
	OperationKey{focus: operation.FocusGrid, key: k("n", "P")}:                                                                   PurposePanic,
	OperationKey{focus: operation.FocusGrid, key: k("p")}:                                                                        Paste,
	OperationKey{focus: operation.FocusAny, key: k("q")}:                                                                         Quit,
	OperationKey{focus: operation.FocusGrid, key: k("r")}:                                                                        RatchetDecrease,
	OperationKey{focus: operation.FocusGrid, key: k("s", "s")}:                                                                   ActionAddLineReset,
	OperationKey{focus: operation.FocusGrid, key: k("s", "S")}:                                                                   ActionAddLineResetAll,
	OperationKey{focus: operation.FocusGrid, key: k("s", "b")}:                                                                   ActionAddLineBounce,
	OperationKey{focus: operation.FocusGrid, key: k("s", "B")}:                                                                   ActionAddLineBounceAll,
	OperationKey{focus: operation.FocusGrid, key: k("s", "k")}:                                                                   ActionAddSkipBeat,
	OperationKey{focus: operation.FocusGrid, key: k("s", "K")}:                                                                   ActionAddSkipBeatAll,
	OperationKey{focus: operation.FocusGrid, key: k("s", "r")}:                                                                   ActionAddLineReverse,
	OperationKey{focus: operation.FocusGrid, key: k("s", "z")}:                                                                   ActionAddLineDelay,
	OperationKey{focus: operation.FocusAny, key: k("u")}:                                                                         Undo,
	OperationKey{focus: operation.FocusGrid, key: k("v")}:                                                                        ToggleVisualMode,
	OperationKey{focus: operation.FocusGrid, key: k("V")}:                                                                        ToggleVisualLineMode,
	OperationKey{focus: operation.FocusGrid, key: k("w")}:                                                                        WaitDecrease,
	OperationKey{focus: operation.FocusGrid, key: k("x")}:                                                                        OverlayNoteRemove,
	OperationKey{focus: operation.FocusGrid, key: k("y")}:                                                                        Yank,
	OperationKey{focus: operation.FocusGrid, key: k("{")}:                                                                        NextOverlay,
	OperationKey{focus: operation.FocusGrid, key: k("}")}:                                                                        PrevOverlay,
	OperationKey{focus: operation.FocusGrid, key: k("enter")}:                                                                    Enter,
	OperationKey{focus: operation.FocusArrangementEditor, key: k("enter")}:                                                       Enter,
	OperationKey{focus: operation.FocusOverlayKey, key: k("enter")}:                                                              ConfirmOverlayKey,
	OperationKey{selection: operation.SelectRenamePart, key: k("enter")}:                                                         ConfirmRenamePart,
	OperationKey{selection: operation.SelectFileName, key: k("enter")}:                                                           ConfirmFileName,
	OperationKey{selection: operation.SelectPart, key: k("enter")}:                                                               ConfirmSelectPart,
	OperationKey{selection: operation.SelectChangePart, key: k("enter")}:                                                         ConfirmChangePart,
	OperationKey{selection: operation.SelectConfirmNew, key: k("enter")}:                                                         ConfirmConfirmNew,
	OperationKey{selection: operation.SelectConfirmReload, key: k("enter")}:                                                      ConfirmConfirmReload,
	OperationKey{selection: operation.SelectConfirmQuit, key: k("enter")}:                                                        ConfirmConfirmQuit,
	OperationKey{selection: operation.SelectFileName, key: k("esc")}:                                                             Escape,
	OperationKey{selection: operation.SelectSetupChannel, key: k("J")}:                                                           DecreaseAllChannels,
	OperationKey{selection: operation.SelectSetupChannel, key: k("K")}:                                                           IncreaseAllChannels,
	OperationKey{selection: operation.SelectSetupValue, key: k("J")}:                                                             DecreaseAllNote,
	OperationKey{selection: operation.SelectSetupValue, key: k("K")}:                                                             IncreaseAllNote,
	OperationKey{key: k("esc")}:                                                                                                  Escape,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeLine, key: k("!")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeLine, key: k("@")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeLine, key: k("#")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeLine, key: k("$")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeLine, key: k("%")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeLine, key: k("^")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeLine, key: k("&")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeLine, key: k("*")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeLine, key: k("(")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeLine, key: k("1")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeLine, key: k("2")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeLine, key: k("3")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeLine, key: k("4")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeLine, key: k("5")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeLine, key: k("6")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeLine, key: k("7")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeLine, key: k("8")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeLine, key: k("9")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeMono, key: k("!")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeMono, key: k("@")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeMono, key: k("#")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeMono, key: k("$")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeMono, key: k("%")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeMono, key: k("^")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeMono, key: k("&")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeMono, key: k("*")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeMono, key: k("(")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeMono, key: k("1")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeMono, key: k("2")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeMono, key: k("3")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeMono, key: k("4")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeMono, key: k("5")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeMono, key: k("6")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeMono, key: k("7")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeMono, key: k("8")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeMono, key: k("9")}:                                           NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternAccent, mode: operation.SeqModeChord, key: k("!")}:    NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternAccent, mode: operation.SeqModeChord, key: k("@")}:    NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternAccent, mode: operation.SeqModeChord, key: k("#")}:    NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternAccent, mode: operation.SeqModeChord, key: k("$")}:    NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternAccent, mode: operation.SeqModeChord, key: k("%")}:    NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternAccent, mode: operation.SeqModeChord, key: k("^")}:    NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternAccent, mode: operation.SeqModeChord, key: k("&")}:    NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternAccent, mode: operation.SeqModeChord, key: k("*")}:    NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternAccent, mode: operation.SeqModeChord, key: k("(")}:    NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternAccent, mode: operation.SeqModeChord, key: k("1")}:    NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternAccent, mode: operation.SeqModeChord, key: k("2")}:    NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternAccent, mode: operation.SeqModeChord, key: k("3")}:    NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternAccent, mode: operation.SeqModeChord, key: k("4")}:    NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternAccent, mode: operation.SeqModeChord, key: k("5")}:    NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternAccent, mode: operation.SeqModeChord, key: k("6")}:    NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternAccent, mode: operation.SeqModeChord, key: k("7")}:    NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternAccent, mode: operation.SeqModeChord, key: k("8")}:    NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternAccent, mode: operation.SeqModeChord, key: k("9")}:    NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternGate, mode: operation.SeqModeChord, key: k("!")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternGate, mode: operation.SeqModeChord, key: k("@")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternGate, mode: operation.SeqModeChord, key: k("#")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternGate, mode: operation.SeqModeChord, key: k("$")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternGate, mode: operation.SeqModeChord, key: k("%")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternGate, mode: operation.SeqModeChord, key: k("^")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternGate, mode: operation.SeqModeChord, key: k("&")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternGate, mode: operation.SeqModeChord, key: k("*")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternGate, mode: operation.SeqModeChord, key: k("(")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternGate, mode: operation.SeqModeChord, key: k("1")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternGate, mode: operation.SeqModeChord, key: k("2")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternGate, mode: operation.SeqModeChord, key: k("3")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternGate, mode: operation.SeqModeChord, key: k("4")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternGate, mode: operation.SeqModeChord, key: k("5")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternGate, mode: operation.SeqModeChord, key: k("6")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternGate, mode: operation.SeqModeChord, key: k("7")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternGate, mode: operation.SeqModeChord, key: k("8")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternGate, mode: operation.SeqModeChord, key: k("9")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternWait, mode: operation.SeqModeChord, key: k("!")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternWait, mode: operation.SeqModeChord, key: k("@")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternWait, mode: operation.SeqModeChord, key: k("#")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternWait, mode: operation.SeqModeChord, key: k("$")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternWait, mode: operation.SeqModeChord, key: k("%")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternWait, mode: operation.SeqModeChord, key: k("^")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternWait, mode: operation.SeqModeChord, key: k("&")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternWait, mode: operation.SeqModeChord, key: k("*")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternWait, mode: operation.SeqModeChord, key: k("(")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternWait, mode: operation.SeqModeChord, key: k("1")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternWait, mode: operation.SeqModeChord, key: k("2")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternWait, mode: operation.SeqModeChord, key: k("3")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternWait, mode: operation.SeqModeChord, key: k("4")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternWait, mode: operation.SeqModeChord, key: k("5")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternWait, mode: operation.SeqModeChord, key: k("6")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternWait, mode: operation.SeqModeChord, key: k("7")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternWait, mode: operation.SeqModeChord, key: k("8")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternWait, mode: operation.SeqModeChord, key: k("9")}:      NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternRatchet, mode: operation.SeqModeChord, key: k("!")}:   NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternRatchet, mode: operation.SeqModeChord, key: k("@")}:   NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternRatchet, mode: operation.SeqModeChord, key: k("#")}:   NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternRatchet, mode: operation.SeqModeChord, key: k("$")}:   NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternRatchet, mode: operation.SeqModeChord, key: k("%")}:   NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternRatchet, mode: operation.SeqModeChord, key: k("^")}:   NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternRatchet, mode: operation.SeqModeChord, key: k("&")}:   NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternRatchet, mode: operation.SeqModeChord, key: k("*")}:   NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternRatchet, mode: operation.SeqModeChord, key: k("(")}:   NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternRatchet, mode: operation.SeqModeChord, key: k("1")}:   NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternRatchet, mode: operation.SeqModeChord, key: k("2")}:   NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternRatchet, mode: operation.SeqModeChord, key: k("3")}:   NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternRatchet, mode: operation.SeqModeChord, key: k("4")}:   NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternRatchet, mode: operation.SeqModeChord, key: k("5")}:   NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternRatchet, mode: operation.SeqModeChord, key: k("6")}:   NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternRatchet, mode: operation.SeqModeChord, key: k("7")}:   NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternRatchet, mode: operation.SeqModeChord, key: k("8")}:   NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternRatchet, mode: operation.SeqModeChord, key: k("9")}:   NumberPattern,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("t", "M")}: MajorTriad,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("t", "m")}: MinorTriad,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("t", "d")}: DiminishedTriad,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("t", "a")}: AugmentedTriad,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("7", "m")}: MinorSeventh,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("7", "M")}: MajorSeventh,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("5", "a")}: AugFifth,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("5", "d")}: DimFifth,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("5", "p")}: PerfectFifth,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("2", "m")}: MinorSecond,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("2", "M")}: MajorSecond,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("3", "m")}: MinorThird,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("3", "M")}: MajorThird,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("4", "p")}: PerfectFourth,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("6", "M")}: MajorSixth,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("8", "p")}: Octave,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("9", "m")}: MinorNinth,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("9", "M")}: MajorNinth,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("[", "i")}: DecreaseInversions,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("]", "i")}: IncreaseInversions,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("1", "o")}: OmitRoot,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("2", "o")}: OmitSecond,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("3", "o")}: OmitThird,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("4", "o")}: OmitFourth,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("5", "o")}: OmitFifth,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("6", "o")}: OmitSixth,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("7", "o")}: OmitSeventh,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("8", "o")}: OmitOctave,
	OperationKey{focus: operation.FocusGrid, patternMode: operation.PatternFill, mode: operation.SeqModeChord, key: k("9", "o")}: OmitNinth,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeChord, key: k("X")}:                                          RemoveChord,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeChord, key: k("]", "p")}:                                     NextArpeggio,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeChord, key: k("[", "p")}:                                     PrevArpeggio,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeChord, key: k("]", "d")}:                                     NextDouble,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeChord, key: k("[", "d")}:                                     PrevDouble,
	OperationKey{focus: operation.FocusGrid, mode: operation.SeqModeChord, key: k("n", "n")}:                                     ConvertToNotes,
	OperationKey{focus: operation.FocusOverlayKey, key: [3]string{}}:                                                             OverlayKeyMessage,
	OperationKey{selection: operation.SelectRenamePart, key: [3]string{}}:                                                        TextInputMessage,
	OperationKey{selection: operation.SelectFileName, key: [3]string{}}:                                                          TextInputMessage,
	OperationKey{focus: operation.FocusArrangementEditor, selection: operation.SelectFileName, key: [3]string{}}:                 TextInputMessage,
	OperationKey{focus: operation.FocusArrangementEditor, key: [3]string{}}:                                                      ArrKeyMessage,
	OperationKey{focus: operation.FocusArrangementEditor, key: k("'")}:                                                           HoldingKeys,
}

func k(x ...string) [3]string {
	if len(x) <= 3 {
		combo := [3]string{}
		copy(combo[:], x)
		return combo
	} else {
		panic("Can't have key combos longer than 3")
	}
}

var holdKeysTime = time.Millisecond * 500

func ResetKeycombo() {
	mutex.Lock()
	defer mutex.Unlock()
	if timer != nil {
		timer.Stop()
	}
	Keycombo = make([]tea.KeyMsg, 0, 3)
}

func ProcessKey(key tea.KeyMsg, focus operation.Focus, selection operation.Selection, seqtype operation.SequencerMode, patternMode operation.PatternMode) Mapping {
	mutex.Lock()
	defer mutex.Unlock()
	if len(Keycombo) < 3 {
		Keycombo = append(Keycombo, key)
	} else {
		Keycombo = slices.Delete(Keycombo, 0, 1)
		Keycombo = append(Keycombo, key)
	}

	if timer != nil {
		timer.Stop()
	}

	var mk mappingKey
	for i, msg := range Keycombo {
		mk[i] = msg.String()
	}

	operationKeys := []OperationKey{
		// Route enter and esc to the textInput mappings
		ToMappingKey(mk, operation.FocusAny, selection, operation.SeqModeAny, operation.PatternAny),
		// Route any text keys to the textInput
		ToMappingKey([3]string{}, operation.FocusAny, selection, operation.SeqModeAny, operation.PatternAny),
		// For mappings good in any focus, like Play
		ToMappingKey(mk, operation.FocusAny, operation.SelectAny, operation.SeqModeAny, operation.PatternAny),
		// Handle focus specific commands handled by ui
		ToMappingKey(mk, focus, operation.SelectAny, operation.SeqModeAny, operation.PatternAny),
		// Route any text keys to other focuses
		ToMappingKey([3]string{}, focus, operation.SelectAny, operation.SeqModeAny, operation.PatternAny),
		// Route selection specific mappings to those selections.
		ToMappingKey(mk, focus, selection, operation.SeqModeAny, operation.PatternAny),
		ToMappingKey(mk, focus, operation.SelectAny, seqtype, patternMode),
		ToMappingKey(mk, focus, operation.SelectAny, seqtype, operation.PatternAny),
	}

	var command Command
	var exists bool

	for _, opKey := range operationKeys {
		command, exists = mappings[opKey]
		if exists {
			break
		}
	}

	if !exists {
		timer = time.AfterFunc(holdKeysTime, func() {
			mutex.Lock()
			defer mutex.Unlock()
			Keycombo = make([]tea.KeyMsg, 0, 3)
		})
	}

	if exists && HoldingKeys == command {
		return Mapping{HoldingKeys, key.String()}
	} else if exists {
		Keycombo = make([]tea.KeyMsg, 0, 3)
		return Mapping{command, key.String()}
	} else {
		return Mapping{HoldingKeys, key.String()}
	}
}

func ToMappingKey(mk [3]string, focus operation.Focus, selection operation.Selection, seqMode operation.SequencerMode, patternMode operation.PatternMode) OperationKey {

	return OperationKey{
		key:         mk,
		focus:       focus,
		selection:   selection,
		mode:        seqMode,
		patternMode: patternMode,
	}
}
