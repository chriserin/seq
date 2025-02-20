package config

import (
	"fmt"

	"github.com/aarzilli/golua/lua"
	"github.com/charmbracelet/lipgloss"
	"github.com/chriserin/seq/internal/grid"
)

type Config struct {
	Accents     []Accent
	C1          int
	LineActions map[grid.Action]lineaction
}

type Accent struct {
	Shape rune
	Color lipgloss.Color
	Value uint8
}

var Accents = []Accent{
	{' ', "#000000", 0},
	{'✤', "#ed3902", 120},
	{'⎈', "#f564a9", 105},
	{'⚙', "#f8730e", 90},
	{'⊚', "#fcc05c", 75},
	{'✦', "#5cdffb", 60},
	{'❖', "#1e89ef", 45},
	{'✥', "#164de5", 30},
	{'❄', "#0246a7", 15},
}

const C1 = 36

type lineaction struct {
	Shape rune
	Color lipgloss.Color
}

var Lineactions = map[grid.Action]lineaction{
	grid.ACTION_NOTHING:        {' ', "#000000"},
	grid.ACTION_LINE_RESET:     {'↔', "#cf142b"},
	grid.ACTION_LINE_REVERSE:   {'←', "#f8730e"},
	grid.ACTION_LINE_SKIP_BEAT: {'⇒', "#a9e5bb"},
	grid.ACTION_RESET:          {'⇚', "#fcf6b1"},
	grid.ACTION_LINE_BOUNCE:    {'↨', "#fcf6b1"},
	grid.ACTION_LINE_DELAY:     {'ℤ', "#cc4bc2"},
}

type ratchetDiacritical string

var Ratchets = []ratchetDiacritical{
	"",
	"\u0307",
	"\u030A",
	"\u030B",
	"\u030C",
	"\u0312",
	"\u0313",
	"\u0344",
}

type Gate struct {
	Shape string
	Value uint16
}

var Gates = []Gate{
	{"", 20},
	{"\u032A", 40},
	{"\u032B", 80},
	{"\u032C", 160},
	{"\u032D", 240},
	{"\u032E", 320},
	{"\u032F", 480},
	{"\u0330", 640},
}

type Wait int16

var WaitPercentages = []Wait{
	0,
	8,
	16,
	24,
	32,
	40,
	48,
	54,
}

type ControlChange struct {
	Value      uint8
	UpperLimit uint8
	Name       string
}

var StandardCCs = []ControlChange{
	{0, 127, "Bank Select"},
	{1, 127, "Modulation Wheel or Lever"},
	{2, 127, "Breath Controller"},
	{4, 127, "Foot Controller"},
	{5, 127, "Portamento Time"},
	{6, 127, "Data Entry MSB"},
	{7, 127, "Channel Volume"},
	{8, 127, "Balance"},
	{10, 127, "Pan"},
	{11, 127, "Expression Controller"},
	{12, 127, "Effect Control 1"},
	{13, 127, "Effect Control 2"},
	{16, 127, "General Purpose Controller 1"},
	{17, 127, "General Purpose Controller 2"},
	{18, 127, "General Purpose Controller 3"},
	{19, 127, "General Purpose Controller 4"},
}

func FindCC(value uint8, instrumentName string) ControlChange {
	instrument := GetInstrument(instrumentName)
	for _, cc := range instrument.CCs {
		if cc.Value == value {
			return cc
		}
	}
	for _, cc := range StandardCCs {
		if cc.Value == value {
			return cc
		}
	}
	return StandardCCs[0]
}

type Template struct {
	Name    string
	Lines   []grid.LineDefinition
	UIStyle string
}

var templates []Template

func GetTemplate(name string) Template {
	for _, template := range templates {
		if template.Name == name {
			return template
		}
	}
	return Template{}
}

type Instrument struct {
	Name string
	CCs  []ControlChange
}

var instruments []Instrument

func GetInstrument(name string) Instrument {
	for _, instrument := range instruments {
		if instrument.Name == name {
			return instrument
		}
	}
	return Instrument{}
}

func ProcessConfig(luafilepath string) {
	L := lua.NewState()
	defer L.Close()

	L.OpenPackage()
	L.OpenLibs()
	L.RegisterLibrary("seq", seqFunctions)

	err := L.DoFile(luafilepath)
	if err != nil {
		fmt.Println("Do File error!!")
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Number of templates %d\n", len(templates))
	}
}

// Lua Function
func addInstrument(L *lua.State) int {
	if L.IsTable(1) {
		L.GetField(1, "name")
		name := L.ToString(2)
		instrument := Instrument{Name: name}

		L.Pop(1)
		L.GetField(1, "controlchanges")
		if L.IsTable(2) {

			for i := 1; true; i++ {
				L.PushInteger(int64(i))
				L.GetTable(2)
				if L.IsTable(3) {
					cc := ControlChange{}
					for i := range 3 {
						L.PushInteger(int64(i + 1))
						L.GetTable(3)
						switch i + 1 {
						case 1:
							value := L.ToNumber(4)
							cc.Value = uint8(value)
						case 2:
							upperLimit := L.ToNumber(4)
							cc.UpperLimit = uint8(upperLimit)
						case 3:
							name := L.ToString(4)
							cc.Name = name
						}
						L.Pop(1)
					}
					instrument.CCs = append(instrument.CCs, cc)
				} else {
					break
				}
				L.Pop(1)
			}
		}
		instruments = append(instruments, instrument)
	} else {
		panic("Instrument not formatted correctly")
		// Communicate Error
	}
	return 0
}

// Lua Function
func addTemplate(L *lua.State) int {
	if L.IsTable(1) {
		L.GetField(1, "name")
		name := L.ToString(2)
		L.Pop(1)
		L.GetField(1, "uistyle")
		uistyle := L.ToString(2)
		if uistyle == "" {
			uistyle = "plain"
		}
		L.Pop(1)

		template := Template{Name: name, UIStyle: uistyle}

		L.GetField(1, "lines")
		if L.IsTable(2) {

			for i := 1; true; i++ {
				L.PushInteger(int64(i))
				L.GetTable(2)
				if L.IsTable(3) {
					ld := grid.LineDefinition{}
					for i := range 3 {
						L.PushInteger(int64(i + 1))
						L.GetTable(3)
						switch i + 1 {
						case 1:
							channel := L.ToNumber(4)
							ld.Channel = uint8(channel)
						case 2:
							messageType := L.ToString(4)
							switch messageType {
							case "NOTE":
								ld.MsgType = grid.MESSAGE_TYPE_NOTE
							case "CC":
								ld.MsgType = grid.MESSAGE_TYPE_CC
							}
						case 3:
							note := L.ToNumber(4)
							ld.Note = uint8(note)
						}
						L.Pop(1)
					}
					template.Lines = append(template.Lines, ld)
				} else {
					break
				}
				L.Pop(1)
			}
		}
		templates = append(templates, template)
	} else {
		panic("Template not formatted correctly")
		// Communicate Error
	}
	return 0
}

type LuaFn = lua.LuaGoFunction

var seqFunctions = map[string]LuaFn{
	"addtemplate":   addTemplate,
	"addinstrument": addInstrument,
}
