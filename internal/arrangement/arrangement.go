package arrangement

import (
	"fmt"
	"strings"

	"github.com/chriserin/seq/internal/colors"
	"github.com/chriserin/seq/internal/overlays"
)

type Model struct {
	focus       bool
	position    int
	arrangement *[]SongSection
	parts       *[]Part
}

func InitModel(arrangement *[]SongSection, parts *[]Part) Model {
	return Model{
		arrangement: arrangement,
		parts:       parts,
	}
}

type SongSection struct {
	Part        int
	Cycles      int
	StartBeat   int
	StartCycles int
}

type Part struct {
	Overlays *overlays.Overlay
	Beats    uint8
	Name     string
}

func (m Model) View(currentSongSection int) string {
	var buf strings.Builder
	for i, songSection := range *m.arrangement {
		section := fmt.Sprintf("%d) %s", i+1, (*m.parts)[songSection.Part].GetName())
		var sectionOutput string
		if currentSongSection == i {
			sectionOutput = colors.SelectedColor.Render(section)
		} else {
			sectionOutput = section
		}
		buf.WriteString(sectionOutput)
		buf.WriteString("\n")
	}
	return buf.String()
}

func (p Part) GetName() string {
	return p.Name
}
