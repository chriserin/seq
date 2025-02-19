package config

import (
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
	{1, 127, "Mod Wheel"},
	{4, 120, "Foot Controller"},
	{5, 120, "Portamento"},
}

var CCs = []ControlChange{
	{3, 120, "OSC A FREQUENCY"},
	{7, 120, "MASTER VOLUME"},
	{9, 120, "OSC B FREQUENCY"},
	{14, 127, "OSC B FINE TUNE"},
	{15, 1, "OSC A SAW ON/FF"},
	{20, 1, "OSC A SQUARE ON/OFF"},
	{21, 120, "OSC A PULSE WIDTH"},
	{22, 120, "OSC B PULSE WIDTH"},
	{23, 1, "OSC SYNC ON/OFF"},
	{24, 1, "OSC B LOW FREQ ON/OFF"},
	{25, 1, "OSC B KEYBOARD ON/OFF"},
	{26, 120, "GLIDE RATE"},
	{27, 120, "OSC A LEVEL"},
	{28, 120, "OSC B LEVEL"},
	{29, 120, "NOISE LEVEL"},
	{30, 1, "OSC B SAW ON/OFF"},
	{31, 120, "RESONANCE"},
	{35, 2, "FILTER KEYBOARD TRACK OFF/HALF/FULL"},
	{41, 1, "FILTER REV SELECT"},
	{46, 120, "LFO FREQUENCY"},
	{47, 120, "LFO INITIAL AMOUNT"},
	{52, 1, "OSC B TRI ON/OFF"},
	{53, 120, "LFO SOURCE MIX"},
	{54, 1, "LFO FREQ A ON/OFF"},
	{55, 1, "LFO FREQ B ON/OFF"},
	{56, 1, "LFO FREQ PW A ON/OFF"},
	{57, 1, "LFO FREQ PW B ON/OFF"},
	{58, 1, "LFO FILTER ON/OFF"},
	{59, 127, "POLY MOD FILT ENV AMOUNT"},
	{60, 120, "POLY MOD OSC B AMOUNT"},
	{61, 1, "POLY MOD FREQ A ON/OFF"},
	{62, 1, "POLY MOD PW ON/OFF"},
	{63, 1, "POLY MOD FILTER ON/OFF"},
	{70, 11, "PITCH WHEEL RANGE"},
	{71, 3, "RETRIGGER AND UNISON ASSIGN"},
	{73, 120, "CUTOFF"},
	{74, 127, "BRIGHTNESS"},
	{85, 127, "VINTAGE"},
	{86, 1, "PRESSURE FILTER"},
	{87, 1, "PRESSURE LFO"},
	{89, 120, "ENVELOPE FILTER AMOUNT"},
	{90, 1, "ENVELOPE FILTER VELOCITY ON/OFF"},
	{102, 1, "ENVELOPE VCA VELOCITY ON/OFF"},
	{103, 120, "ATTACK FILTER"},
	{104, 120, "ATTACK VCA"},
	{105, 120, "DECAY FILTER"},
	{106, 120, "DECAY VCA"},
	{107, 120, "SUSTAIN FILTER"},
	{108, 120, "SUSTAIN VCA"},
	{109, 120, "RELEASE FILTER"},
	{110, 120, "RELEASE VCA"},
	{111, 1, "RELEASE ON/OFF"},
	{112, 1, "UNISON ON/OFF"},
	{113, 10, "UNISON VOICE COUNT"},
	{114, 7, "UNISON DETUNE"},
	{116, 1, "OSC B SQUARE ON/OFF"},
	{117, 1, "LFO SAW ON/OFF"},
	{118, 1, "LFO TRI ON/OFF"},
	{119, 1, "LFO SQUARE ON/OFF"},
}

func FindCC(value uint8) ControlChange {
	for _, cc := range CCs {
		if cc.Value == value {
			return cc
		}
	}
	for _, cc := range StandardCCs {
		if cc.Value == value {
			return cc
		}
	}
	return CCs[0]
}
