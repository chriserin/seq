package sequence

import (
	"math"
	"slices"

	"github.com/chriserin/seq/internal/arrangement"
	"github.com/chriserin/seq/internal/config"
	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/operation"
)

type Sequence struct {
	Parts                 *[]arrangement.Part
	Arrangement           *arrangement.Arrangement
	Lines                 []grid.LineDefinition
	Tempo                 int
	Subdivisions          int
	Keyline               uint8
	Accents               PatternAccents
	Instrument            string
	Template              string
	TemplateUIStyle       string
	TemplateSequencerType operation.SequencerMode
}

type PatternAccents struct {
	Data   []config.Accent
	Start  uint8
	End    uint8
	Target AccentTarget
}

type AccentTarget uint8

const (
	AccentTargetNote AccentTarget = iota
	AccentTargetVelocity
)

func (pa *PatternAccents) ReCalc() {
	accents := make([]config.Accent, 9)

	interval := float64(pa.Start-pa.End) / float64(len(pa.Data)-2)

	for i := range pa.Data[1:] {
		calculatedValue := float64(pa.Start) - (interval * float64(i))
		roundedValue := math.Round(calculatedValue)
		a := config.Accent(roundedValue)
		accents[i+1] = a
	}

	pa.Data = accents
}

func (pa *PatternAccents) Equal(other *PatternAccents) bool {
	if pa.Target != other.Target {
		return false
	}
	if pa.Start != other.Start {
		return false
	}
	if pa.End != other.End {
		return false
	}
	if !slices.Equal(pa.Data, other.Data) {
		return false
	}
	return true
}

func InitParts() []arrangement.Part {
	firstPart := arrangement.InitPart("Part 1")
	return []arrangement.Part{firstPart}
}

func InitSequence(template string, instrument string) Sequence {
	gridTemplate, exists := config.GetTemplate(template)
	if !exists {
		gridTemplate = config.GetDefaultTemplate()
	}
	config.LongGates = config.GetGateLengths(32)
	newLines := make([]grid.LineDefinition, len(gridTemplate.Lines))
	copy(newLines, gridTemplate.Lines)

	parts := InitParts()
	return Sequence{
		Parts:                 &parts,
		Arrangement:           InitArrangement(parts),
		Tempo:                 120,
		Keyline:               0,
		Subdivisions:          2,
		Lines:                 newLines,
		Accents:               PatternAccents{End: 15, Data: config.Accents, Start: 120, Target: AccentTargetVelocity},
		Template:              gridTemplate.Name,
		Instrument:            instrument,
		TemplateUIStyle:       gridTemplate.UIStyle,
		TemplateSequencerType: gridTemplate.SequencerType,
	}
}

func InitArrangement(parts []arrangement.Part) *arrangement.Arrangement {
	root := arrangement.InitRoot(parts)

	for i := range parts {
		section := arrangement.InitSongSection(i)

		node := &arrangement.Arrangement{
			Section:    section,
			Iterations: 1,
		}

		root.Nodes = append(root.Nodes, node)
	}

	return root
}
