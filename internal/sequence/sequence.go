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

	for i, a := range pa.Data[1:] {
		calculatedValue := float64(pa.Start) - (interval * float64(i))
		roundedValue := math.Round(calculatedValue)
		a.Value = uint8(roundedValue)
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
