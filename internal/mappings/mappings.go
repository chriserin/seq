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
	Help
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
	TempoInputSwitch
	OverlayInputSwitch
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
	ToggleWaitMode
	ToggleAccentMode
	ToggleRatchetMode
	NextOverlay
	PrevOverlay
	Save
	Undo
	Redo
	New
	ToggleVisualMode
	TogglePlayEdit
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
	RatchetIncrease
	RatchetDecrease
	ActionAddLineReset
	ActionAddLineResetAll
	ActionAddLineReverse
	// ActionAddLineReverseAll
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
	key         mappingKey
	focus       operation.Focus
	selection   operation.Selection
	mode        operation.SequencerMode
	patternMode operation.PatternMode
}

var mappings = registry{
	OperationKey{key: k(" ")}:                                    PlayStop,
	OperationKey{key: k("'", " ")}:                               PlayOverlayLoop,
	OperationKey{key: k(":", " ")}:                               PlayRecord,
	OperationKey{key: k("+")}:                                    Increase,
	OperationKey{key: k("-")}:                                    Decrease,
	OperationKey{key: k("<")}:                                    CursorLineStart,
	OperationKey{key: k("=")}:                                    Increase,
	OperationKey{key: k(">")}:                                    CursorLineEnd,
	OperationKey{key: k("b", "l")}:                               CursorLastLine,
	OperationKey{key: k("b", "f")}:                               CursorFirstLine,
	OperationKey{key: k("?")}:                                    Help,
	OperationKey{key: k("A")}:                                    AccentIncrease,
	OperationKey{key: k("C")}:                                    ClearOverlay,
	OperationKey{key: k("G")}:                                    GateIncrease,
	OperationKey{key: k("E")}:                                    GateBigIncrease,
	OperationKey{key: k("J")}:                                    RotateDown,
	OperationKey{key: k("K")}:                                    RotateUp,
	OperationKey{key: k("H")}:                                    RotateLeft,
	OperationKey{key: k("L")}:                                    RotateRight,
	OperationKey{key: k("Y")}:                                    SelectKeyLine,
	OperationKey{key: k("M")}:                                    Solo,
	OperationKey{key: k("R")}:                                    RatchetIncrease,
	OperationKey{key: k("U")}:                                    Redo,
	OperationKey{key: k("W")}:                                    WaitIncrease,
	OperationKey{key: k("[", "c")}:                               PrevTheme,
	OperationKey{key: k("[", "s")}:                               PrevSection,
	OperationKey{key: k("]", "c")}:                               NextTheme,
	OperationKey{key: k("]", "s")}:                               NextSection,
	OperationKey{key: k("g")}:                                    GateDecrease,
	OperationKey{key: k("e")}:                                    GateBigDecrease,
	OperationKey{key: k("alt+ ")}:                                PlayLoop,
	OperationKey{key: k("ctrl+@")}:                               PlayPart,
	OperationKey{key: k("ctrl+]")}:                               NewSectionAfter,
	OperationKey{key: k("ctrl+b")}:                               BeatInputSwitch,
	OperationKey{key: k("ctrl+k")}:                               CyclesInputSwitch,
	OperationKey{key: k("ctrl+c")}:                               ChangePart,
	OperationKey{key: k("ctrl+e")}:                               AccentInputSwitch,
	OperationKey{key: k("ctrl+a")}:                               ToggleArrangementView,
	OperationKey{key: k("ctrl+l")}:                               NewLine,
	OperationKey{key: k("ctrl+n")}:                               New,
	OperationKey{key: k("ctrl+o")}:                               OverlayInputSwitch,
	OperationKey{key: k("ctrl+p")}:                               NewSectionBefore,
	OperationKey{key: k("ctrl+d")}:                               SetupInputSwitch,
	OperationKey{key: k("ctrl+t")}:                               TempoInputSwitch,
	OperationKey{key: k("ctrl+u")}:                               OverlayStackToggle,
	OperationKey{key: k("ctrl+s")}:                               Save,
	OperationKey{key: k("ctrl+y")}:                               RatchetInputSwitch,
	OperationKey{key: k("a")}:                                    AccentDecrease,
	OperationKey{key: k("c")}:                                    ClearLine,
	OperationKey{key: k("d")}:                                    NoteRemove,
	OperationKey{key: k("b", "e")}:                               TogglePlayEdit,
	OperationKey{key: k("f")}:                                    NoteAdd,
	OperationKey{key: k("b", "r")}:                               ReloadFile,
	OperationKey{key: k("b", "v")}:                               ActionAddSpecificValue,
	OperationKey{key: k("h")}:                                    CursorLeft,
	OperationKey{key: k("j")}:                                    CursorDown,
	OperationKey{key: k("k")}:                                    CursorUp,
	OperationKey{key: k("l")}:                                    CursorRight,
	OperationKey{key: k("m")}:                                    Mute,
	OperationKey{key: k("o")}:                                    ToggleChordMode,
	OperationKey{key: k("n", "a")}:                               ToggleAccentMode,
	OperationKey{key: k("n", "w")}:                               ToggleWaitMode,
	OperationKey{key: k("n", "g")}:                               ToggleGateMode,
	OperationKey{key: k("n", "r")}:                               ToggleRatchetMode,
	OperationKey{key: k("p")}:                                    Paste,
	OperationKey{key: k("q")}:                                    Quit,
	OperationKey{key: k("r")}:                                    RatchetDecrease,
	OperationKey{key: k("s", "s")}:                               ActionAddLineReset,
	OperationKey{key: k("s", "S")}:                               ActionAddLineResetAll,
	OperationKey{key: k("s", "b")}:                               ActionAddLineBounce,
	OperationKey{key: k("s", "B")}:                               ActionAddLineBounceAll,
	OperationKey{key: k("s", "k")}:                               ActionAddSkipBeat,
	OperationKey{key: k("s", "K")}:                               ActionAddSkipBeatAll,
	OperationKey{key: k("s", "r")}:                               ActionAddLineReverse,
	OperationKey{key: k("s", "z")}:                               ActionAddLineDelay,
	OperationKey{key: k("u")}:                                    Undo,
	OperationKey{key: k("v")}:                                    ToggleVisualMode,
	OperationKey{key: k("w")}:                                    WaitDecrease,
	OperationKey{key: k("x")}:                                    OverlayNoteRemove,
	OperationKey{key: k("y")}:                                    Yank,
	OperationKey{key: k("{")}:                                    NextOverlay,
	OperationKey{key: k("}")}:                                    PrevOverlay,
	OperationKey{key: k("enter")}:                                Enter,
	OperationKey{key: k("esc")}:                                  Escape,
	OperationKey{key: k("!")}:                                    NumberPattern,
	OperationKey{key: k("@")}:                                    NumberPattern,
	OperationKey{key: k("#")}:                                    NumberPattern,
	OperationKey{key: k("$")}:                                    NumberPattern,
	OperationKey{key: k("%")}:                                    NumberPattern,
	OperationKey{key: k("^")}:                                    NumberPattern,
	OperationKey{key: k("&")}:                                    NumberPattern,
	OperationKey{key: k("*")}:                                    NumberPattern,
	OperationKey{key: k("(")}:                                    NumberPattern,
	OperationKey{key: k("1")}:                                    NumberPattern,
	OperationKey{key: k("2")}:                                    NumberPattern,
	OperationKey{key: k("3")}:                                    NumberPattern,
	OperationKey{key: k("4")}:                                    NumberPattern,
	OperationKey{key: k("5")}:                                    NumberPattern,
	OperationKey{key: k("6")}:                                    NumberPattern,
	OperationKey{key: k("7")}:                                    NumberPattern,
	OperationKey{key: k("8")}:                                    NumberPattern,
	OperationKey{key: k("9")}:                                    NumberPattern,
	OperationKey{mode: operation.SeqModeChord, key: k("t", "M")}: MajorTriad,
	OperationKey{mode: operation.SeqModeChord, key: k("t", "m")}: MinorTriad,
	OperationKey{mode: operation.SeqModeChord, key: k("t", "d")}: DiminishedTriad,
	OperationKey{mode: operation.SeqModeChord, key: k("t", "a")}: AugmentedTriad,
	OperationKey{mode: operation.SeqModeChord, key: k("7", "m")}: MinorSeventh,
	OperationKey{mode: operation.SeqModeChord, key: k("7", "M")}: MajorSeventh,
	OperationKey{mode: operation.SeqModeChord, key: k("5", "a")}: AugFifth,
	OperationKey{mode: operation.SeqModeChord, key: k("5", "d")}: DimFifth,
	OperationKey{mode: operation.SeqModeChord, key: k("5", "p")}: PerfectFifth,
	OperationKey{mode: operation.SeqModeChord, key: k("2", "m")}: MinorSecond,
	OperationKey{mode: operation.SeqModeChord, key: k("2", "M")}: MajorSecond,
	OperationKey{mode: operation.SeqModeChord, key: k("3", "m")}: MinorThird,
	OperationKey{mode: operation.SeqModeChord, key: k("3", "M")}: MajorThird,
	OperationKey{mode: operation.SeqModeChord, key: k("4", "p")}: PerfectFourth,
	OperationKey{mode: operation.SeqModeChord, key: k("6", "M")}: MajorSixth,
	OperationKey{mode: operation.SeqModeChord, key: k("8", "p")}: Octave,
	OperationKey{mode: operation.SeqModeChord, key: k("9", "m")}: MinorNinth,
	OperationKey{mode: operation.SeqModeChord, key: k("9", "M")}: MajorNinth,
	OperationKey{mode: operation.SeqModeChord, key: k("[", "i")}: DecreaseInversions,
	OperationKey{mode: operation.SeqModeChord, key: k("]", "i")}: IncreaseInversions,
	OperationKey{mode: operation.SeqModeChord, key: k("1", "o")}: OmitRoot,
	OperationKey{mode: operation.SeqModeChord, key: k("2", "o")}: OmitSecond,
	OperationKey{mode: operation.SeqModeChord, key: k("3", "o")}: OmitThird,
	OperationKey{mode: operation.SeqModeChord, key: k("4", "o")}: OmitFourth,
	OperationKey{mode: operation.SeqModeChord, key: k("5", "o")}: OmitFifth,
	OperationKey{mode: operation.SeqModeChord, key: k("6", "o")}: OmitSixth,
	OperationKey{mode: operation.SeqModeChord, key: k("7", "o")}: OmitSeventh,
	OperationKey{mode: operation.SeqModeChord, key: k("8", "o")}: OmitOctave,
	OperationKey{mode: operation.SeqModeChord, key: k("9", "o")}: OmitNinth,
	OperationKey{mode: operation.SeqModeChord, key: k("D")}:      RemoveChord,
	OperationKey{mode: operation.SeqModeChord, key: k("]", "p")}: NextArpeggio,
	OperationKey{mode: operation.SeqModeChord, key: k("[", "p")}: PrevArpeggio,
	OperationKey{mode: operation.SeqModeChord, key: k("]", "d")}: NextDouble,
	OperationKey{mode: operation.SeqModeChord, key: k("[", "d")}: PrevDouble,
	OperationKey{mode: operation.SeqModeChord, key: k("n", "n")}: ConvertToNotes,
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

func ProcessKey(key tea.KeyMsg, seqtype operation.SequencerMode, patternMode operation.PatternMode) Mapping {
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
		ToMappingKey(mk, seqtype, patternMode),
		ToMappingKey(mk, seqtype, operation.PatternAny),
		ToMappingKey(mk, operation.SeqModeAny, operation.PatternAny),
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

	if exists {
		Keycombo = make([]tea.KeyMsg, 0, 3)
		return Mapping{command, key.String()}
	} else {
		return Mapping{HoldingKeys, key.String()}
	}
}

func ToMappingKey(mk [3]string, seqMode operation.SequencerMode, patternMode operation.PatternMode) OperationKey {

	return OperationKey{
		key:         mk,
		focus:       operation.FocusAny,
		selection:   operation.SelectAny,
		mode:        seqMode,
		patternMode: patternMode,
	}
}
