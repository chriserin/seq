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
	"github.com/chriserin/seq/internal/grid"
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
type registry map[mappingKey]Command

// KeysForCommand Function that gets keys for mappings by looking at the looping
// through the registry and returning the keys for the given command
// If the command is not found, it returns an empty string slice.
// This is used to get the keys for a given command in the mappings.
func KeysForCommand(command Command) []string {
	var keys []string
	for key, cmd := range allCommands() {
		if cmd == command {
			for _, k := range key {
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
	for _, m := range []registry{mappings, patternModeMappings, chordMappings} {
		maps.Copy(all, m)
	}
	return all
}

var mappings = registry{
	k(" "):      PlayStop,
	k("'", " "): PlayOverlayLoop,
	k(":", " "): PlayRecord,
	k("+"):      Increase,
	k("-"):      Decrease,
	k("<"):      CursorLineStart,
	k("="):      Increase,
	k(">"):      CursorLineEnd,
	k("b", "l"): CursorLastLine,
	k("b", "f"): CursorFirstLine,
	k("?"):      Help,
	k("A"):      AccentIncrease,
	k("C"):      ClearOverlay,
	k("G"):      GateIncrease,
	k("E"):      GateBigIncrease,
	k("J"):      RotateDown,
	k("K"):      RotateUp,
	k("H"):      RotateLeft,
	k("L"):      RotateRight,
	k("Y"):      SelectKeyLine,
	k("M"):      Solo,
	k("R"):      RatchetIncrease,
	k("U"):      Redo,
	k("W"):      WaitIncrease,
	k("[", "c"): PrevTheme,
	k("[", "s"): PrevSection,
	k("]", "c"): NextTheme,
	k("]", "s"): NextSection,
	k("g"):      GateDecrease,
	k("e"):      GateBigDecrease,
	k("alt+ "):  PlayLoop,
	k("ctrl+@"): PlayPart,
	k("ctrl+]"): NewSectionAfter,
	k("ctrl+b"): BeatInputSwitch,
	k("ctrl+k"): CyclesInputSwitch,
	k("ctrl+c"): ChangePart,
	k("ctrl+e"): AccentInputSwitch,
	k("ctrl+a"): ToggleArrangementView,
	k("ctrl+l"): NewLine,
	k("ctrl+n"): New,
	k("ctrl+o"): OverlayInputSwitch,
	k("ctrl+p"): NewSectionBefore,
	k("ctrl+d"): SetupInputSwitch,
	k("ctrl+t"): TempoInputSwitch,
	k("ctrl+u"): OverlayStackToggle,
	k("ctrl+s"): Save,
	k("ctrl+y"): RatchetInputSwitch,
	k("a"):      AccentDecrease,
	k("c"):      ClearLine,
	k("d"):      NoteRemove,
	k("b", "e"): TogglePlayEdit,
	k("f"):      NoteAdd,
	k("b", "r"): ReloadFile,
	k("b", "v"): ActionAddSpecificValue,
	k("h"):      CursorLeft,
	k("j"):      CursorDown,
	k("k"):      CursorUp,
	k("l"):      CursorRight,
	k("m"):      Mute,
	k("o"):      ToggleChordMode,
	k("n", "a"): ToggleAccentMode,
	k("n", "w"): ToggleWaitMode,
	k("n", "g"): ToggleGateMode,
	k("n", "r"): ToggleRatchetMode,
	k("p"):      Paste,
	k("q"):      Quit,
	k("r"):      RatchetDecrease,
	k("s", "s"): ActionAddLineReset,
	k("s", "S"): ActionAddLineResetAll,
	k("s", "b"): ActionAddLineBounce,
	k("s", "B"): ActionAddLineBounceAll,
	k("s", "k"): ActionAddSkipBeat,
	k("s", "K"): ActionAddSkipBeatAll,
	k("s", "r"): ActionAddLineReverse,
	// k("s", "R"): ActionAddLineReverseAll,
	k("s", "z"): ActionAddLineDelay,
	k("u"):      Undo,
	k("v"):      ToggleVisualMode,
	k("w"):      WaitDecrease,
	k("x"):      OverlayNoteRemove,
	k("y"):      Yank,
	k("{"):      NextOverlay,
	k("}"):      PrevOverlay,
	k("enter"):  Enter,
	k("esc"):    Escape,
}

var patternModeMappings = registry{
	k("!"): NumberPattern,
	k("@"): NumberPattern,
	k("#"): NumberPattern,
	k("$"): NumberPattern,
	k("%"): NumberPattern,
	k("^"): NumberPattern,
	k("&"): NumberPattern,
	k("*"): NumberPattern,
	k("("): NumberPattern,
	k("1"): NumberPattern,
	k("2"): NumberPattern,
	k("3"): NumberPattern,
	k("4"): NumberPattern,
	k("5"): NumberPattern,
	k("6"): NumberPattern,
	k("7"): NumberPattern,
	k("8"): NumberPattern,
	k("9"): NumberPattern,
}

var chordMappings = registry{
	k("t", "M"): MajorTriad,
	k("t", "m"): MinorTriad,
	k("t", "d"): DiminishedTriad,
	k("t", "a"): AugmentedTriad,
	k("7", "m"): MinorSeventh,
	k("7", "M"): MajorSeventh,
	k("5", "a"): AugFifth,
	k("5", "d"): DimFifth,
	k("5", "p"): PerfectFifth,
	k("2", "m"): MinorSecond,
	k("2", "M"): MajorSecond,
	k("3", "m"): MinorThird,
	k("3", "M"): MajorThird,
	k("4", "p"): PerfectFourth,
	k("6", "M"): MajorSixth,
	k("8", "p"): Octave,
	k("9", "m"): MinorNinth,
	k("9", "M"): MajorNinth,
	k("[", "i"): DecreaseInversions,
	k("]", "i"): IncreaseInversions,
	k("1", "o"): OmitRoot,
	k("2", "o"): OmitSecond,
	k("3", "o"): OmitThird,
	k("4", "o"): OmitFourth,
	k("5", "o"): OmitFifth,
	k("6", "o"): OmitSixth,
	k("7", "o"): OmitSeventh,
	k("8", "o"): OmitOctave,
	k("9", "o"): OmitNinth,
	k("D"):      RemoveChord,
	k("]", "p"): NextArpeggio,
	k("[", "p"): PrevArpeggio,
	k("]", "d"): NextDouble,
	k("[", "d"): PrevDouble,
	k("n", "n"): ConvertToNotes,
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

func ProcessKey(key tea.KeyMsg, seqtype grid.SequencerType, patternMode bool) Mapping {
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

	command, exists := mappings[ToMappingKey(Keycombo)]
	switch seqtype {
	case grid.SeqtypeTrigger:
		triggerCommand, triggerExists := patternModeMappings[ToMappingKey(Keycombo)]
		if triggerExists {
			command = triggerCommand
			exists = triggerExists
		}
	case grid.SeqtypePolyphony:
		var chordCommand Command
		var chordExists bool

		if patternMode {
			chordCommand, chordExists = patternModeMappings[ToMappingKey(Keycombo)]
		} else {
			chordCommand, chordExists = chordMappings[ToMappingKey(Keycombo)]
		}

		if chordExists {
			command = chordCommand
			exists = chordExists
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

func ToMappingKey(keyCombo []tea.KeyMsg) mappingKey {
	var mk mappingKey
	for i, msg := range keyCombo {
		mk[i] = msg.String()
	}
	return mk
}
