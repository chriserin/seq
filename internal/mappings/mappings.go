package mappings

import (
	"slices"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chriserin/seq/internal/grid"
)

var Keycombo = make([]tea.KeyMsg, 0, 3)
var timer *time.Timer

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
	SelectKeyLine
	PressDownOverlay
	NumberPattern
	RotateRight
	RotateLeft
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
)

type mappingKey [3]string
type registry map[mappingKey]Command

var mappings = registry{
	k(" "):      PlayStop,
	k("+"):      Increase,
	k("-"):      Decrease,
	k("<"):      CursorLineStart,
	k("="):      Increase,
	k(">"):      CursorLineEnd,
	k("?"):      Help,
	k("A"):      AccentIncrease,
	k("B"):      ActionAddLineBounce,
	k("C"):      ClearSeq,
	k("G"):      GateIncrease,
	k("H"):      RotateLeft,
	k("K"):      SelectKeyLine,
	k("L"):      RotateRight,
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
	k("a"):      AccentDecrease,
	k("alt+ "):  PlayLoop,
	k("b"):      ActionAddSkipBeat,
	k("c"):      ClearLine,
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
	k("d"):      TriggerRemove,
	k("e"):      TogglePlayEdit,
	k("enter"):  Enter,
	k("esc"):    Escape,
	k("f"):      TriggerAdd,
	k("g"):      GateDecrease,
	k("h"):      CursorLeft,
	k("j"):      CursorDown,
	k("k"):      CursorUp,
	k("l"):      CursorRight,
	k("m"):      Mute,
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
}

var triggerMappings = registry{
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

func ProcessKey(key tea.KeyMsg, seqtype grid.SequencerType) Mapping {
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
	case grid.SEQTYPE_TRIGGER:
		triggerCommand, triggerExists := triggerMappings[ToMappingKey(Keycombo)]
		if triggerExists {
			command = triggerCommand
			exists = triggerExists
		}
	case grid.SEQTYPE_POLYPHONY:
		chordCommand, chordExists := chordMappings[ToMappingKey(Keycombo)]
		if chordExists {
			command = chordCommand
			exists = chordExists
		}
	}

	if !exists {
		timer = time.AfterFunc(time.Millisecond*750, func() {
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
