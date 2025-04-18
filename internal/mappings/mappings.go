package mappings

import (
	"slices"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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
)

type mappingKey [3]string
type registry map[mappingKey]Command

var mappings = registry{
	k("q"):      Quit,
	k("q"):      Quit,
	k("?"):      Help,
	k("k"):      CursorUp,
	k("j"):      CursorDown,
	k("h"):      CursorLeft,
	k("l"):      CursorRight,
	k("<"):      CursorLineStart,
	k(">"):      CursorLineEnd,
	k("esc"):    Escape,
	k(" "):      PlayStop,
	k("ctrl+@"): PlayPart,
	k("alt+ "):  PlayLoop,
	k("+"):      Increase,
	k("="):      Increase,
	k("-"):      Decrease,
	k("enter"):  Enter,
	k("ctrl+t"): TempoInputSwitch,
	k("ctrl+o"): OverlayInputSwitch,
	k("ctrl+s"): SetupInputSwitch,
	k("ctrl+e"): AccentInputSwitch,
	k("ctrl+y"): RatchetInputSwitch,
	k("ctrl+b"): BeatsInputSwitch,
	k("ctrl+x"): ArrangementInputSwitch,
	k("ctrl+f"): ToggleArrangementView,
	k("ctrl+r"): ToggleRatchetMode,
	k("ctrl+g"): ToggleGateMode,
	k("ctrl+w"): ToggleWaitMode,
	k("ctrl+a"): ToggleAccentMode,
	k("{"):      NextOverlay,
	k("}"):      PrevOverlay,
	k("ctrl+v"): Save,
	k("u"):      Undo,
	k("U"):      Redo,
	k("v"):      ToggleVisualMode,
	k("e"):      TogglePlayEdit,
	k("ctrl+n"): New,
	k("ctrl+l"): NewLine,
	k("ctrl+]"): NewSectionAfter,
	k("ctrl+p"): NewSectionBefore,
	k("ctrl+c"): ChangePart,
	k("]", "s"): NextSection,
	k("[", "s"): PrevSection,
	k("]", "c"): NextTheme,
	k("[", "c"): PrevTheme,
	k("y"):      Yank,
	k("m"):      Mute,
	k("M"):      Solo,
	k("f"):      TriggerAdd,
	k("d"):      TriggerRemove,
	k("A"):      AccentIncrease,
	k("a"):      AccentDecrease,
	k("G"):      GateIncrease,
	k("g"):      GateDecrease,
	k("W"):      WaitIncrease,
	k("w"):      WaitDecrease,
	k("x"):      OverlayTriggerRemove,
	k("c"):      ClearLine,
	k("C"):      ClearSeq,
	k("R"):      RatchetIncrease,
	k("r"):      RatchetDecrease,
	k("s"):      ActionAddLineReset,
	k("S"):      ActionAddLineReverse,
	k("B"):      ActionAddLineBounce,
	k("z"):      ActionAddLineDelay,
	k("b"):      ActionAddSkipBeat,
	k("T"):      ActionAddReset,
	k("K"):      SelectKeyLine,
	k("ctrl+u"): PressDownOverlay,
	k("1"):      NumberPattern,
	k("2"):      NumberPattern,
	k("3"):      NumberPattern,
	k("4"):      NumberPattern,
	k("5"):      NumberPattern,
	k("6"):      NumberPattern,
	k("7"):      NumberPattern,
	k("8"):      NumberPattern,
	k("9"):      NumberPattern,
	k("!"):      NumberPattern,
	k("@"):      NumberPattern,
	k("#"):      NumberPattern,
	k("$"):      NumberPattern,
	k("%"):      NumberPattern,
	k("^"):      NumberPattern,
	k("&"):      NumberPattern,
	k("*"):      NumberPattern,
	k("("):      NumberPattern,
	k("L"):      RotateRight,
	k("H"):      RotateLeft,
	k("p"):      Paste,
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

func ProcessKey(key tea.KeyMsg) Mapping {
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
