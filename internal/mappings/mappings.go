// Package mappings provides keyboard input handling and command mapping functionality
// for the sequencer application. It manages key combinations, processes user input,
// and maps keyboard shortcuts to sequencer commands based on the current mode
// (trigger, polyphony, pattern mode, chord mode, etc.).
package mappings

import (
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
	BeatsInputSwitch
	ArrangementInputSwitch
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
	TriggerAdd
	TriggerRemove
	AccentIncrease
	AccentDecrease
	GateIncrease
	GateDecrease
	GateBigIncrease
	GateBigDecrease
	WaitIncrease
	WaitDecrease
	OverlayTriggerRemove
	ClearLine
	ClearSeq
	RatchetIncrease
	RatchetDecrease
	ActionAddLineReset
	ActionAddLineReverse
	ActionAddSkipBeat
	ActionAddReset
	ActionAddLineBounce
	ActionAddLineDelay
	ActionAddSpecificValue
	SelectKeyLine
	PressDownOverlay
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
	MinorSeventh
	MajorSeventh
	AugFifth
	DimFifth
	PerfectFifth
	IncreaseInversions
	DecreaseInversions
	ToggleChordMode
	PatternMode
	NextArppegio
	PrevArppegio
	NextDouble
	PrevDouble
	OmitRoot
	OmitSecond
	OmitThird
	OmitFourth
	OmitFifth
	OmitSixth
	OmitSeventh
	OmitNinth
	RemoveChord
	ConvertToNotes
	ReloadFile
)

type mappingKey [3]string
type registry map[mappingKey]Command

var mappings = registry{
	k(" "):      PlayStop,
	k("'", " "): PlayOverlayLoop,
	k(":", " "): PlayRecord,
	k("+"):      Increase,
	k("-"):      Decrease,
	k("<"):      CursorLineStart,
	k("="):      Increase,
	k(">"):      CursorLineEnd,
	k("?"):      Help,
	k("A"):      AccentIncrease,
	k("B"):      ActionAddLineBounce,
	k("C"):      ClearSeq,
	k("]", "g"): GateIncrease,
	k("]", "e"): GateBigIncrease,
	k("J"):      RotateDown,
	k("K"):      RotateUp,
	k("H"):      RotateLeft,
	k("L"):      RotateRight,
	k("Y"):      SelectKeyLine,
	k("M"):      Solo,
	k("R"):      RatchetIncrease,
	k("S"):      ActionAddLineReverse,
	k("T"):      ActionAddReset,
	k("U"):      Redo,
	k("W"):      WaitIncrease,
	k("[", "c"): PrevTheme,
	k("[", "s"): PrevSection,
	k("]", "c"): NextTheme,
	k("]", "s"): NextSection,
	k("[", "g"): GateDecrease,
	k("[", "e"): GateBigDecrease,
	k("alt+ "):  PlayLoop,
	k("ctrl+@"): PlayPart,
	k("ctrl+]"): NewSectionAfter,
	k("ctrl+a"): ToggleAccentMode,
	k("ctrl+b"): BeatsInputSwitch,
	k("ctrl+c"): ChangePart,
	k("ctrl+e"): AccentInputSwitch,
	k("ctrl+f"): ToggleArrangementView,
	k("ctrl+g"): ToggleGateMode,
	k("ctrl+l"): NewLine,
	k("ctrl+n"): New,
	k("ctrl+o"): OverlayInputSwitch,
	k("ctrl+p"): NewSectionBefore,
	k("ctrl+r"): ToggleRatchetMode,
	k("ctrl+s"): SetupInputSwitch,
	k("ctrl+t"): TempoInputSwitch,
	k("ctrl+u"): PressDownOverlay,
	k("ctrl+v"): Save,
	k("ctrl+w"): ToggleWaitMode,
	k("ctrl+x"): ArrangementInputSwitch,
	k("ctrl+y"): RatchetInputSwitch,
	k("a"):      AccentDecrease,
	k("b"):      ActionAddSkipBeat,
	k("c"):      ClearLine,
	k("d"):      TriggerRemove,
	k("e"):      TogglePlayEdit,
	k("f"):      TriggerAdd,
	k("g", "r"): ReloadFile,
	k("g", "v"): ActionAddSpecificValue,
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
	k("s"):      ActionAddLineReset,
	k("u"):      Undo,
	k("v"):      ToggleVisualMode,
	k("w"):      WaitDecrease,
	k("x"):      OverlayTriggerRemove,
	k("y"):      Yank,
	k("z"):      ActionAddLineDelay,
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
	k("[", "i"): DecreaseInversions,
	k("]", "i"): IncreaseInversions,
	k("1", "o"): OmitRoot,
	k("2", "o"): OmitSecond,
	k("3", "o"): OmitThird,
	k("4", "o"): OmitFourth,
	k("5", "o"): OmitFifth,
	k("6", "o"): OmitSixth,
	k("7", "o"): OmitSeventh,
	k("9", "o"): OmitNinth,
	k("D"):      RemoveChord,
	k("]", "p"): NextArppegio,
	k("[", "p"): PrevArppegio,
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
		timer = time.AfterFunc(time.Millisecond*750, func() {
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
