package arrangement

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	colors "github.com/chriserin/seq/internal/colors"
)

func (m Model) View() string {
	var buf strings.Builder
	buf.WriteString(lipgloss.PlaceHorizontal(18, lipgloss.Left, "  Section"))
	buf.WriteString(lipgloss.PlaceHorizontal(11, lipgloss.Right, "Start Beat"))
	buf.WriteString(lipgloss.PlaceHorizontal(11, lipgloss.Right, "Start Cycle"))
	buf.WriteString(lipgloss.PlaceHorizontal(11, lipgloss.Right, "Cycles"))
	buf.WriteString(lipgloss.PlaceHorizontal(11, lipgloss.Right, "Keep"))
	buf.WriteString("\n")

	m.renderNode(&buf, m.Root, 0)

	return buf.String()
}

// Recursively render a node and its children
func (m Model) renderNode(buf *strings.Builder, node *Arrangement, depth int) {
	if node == nil {
		return
	}

	// For non-end nodes (groups), show iterations
	if node.IsGroup() && depth > 0 {
		indentation := strings.Repeat("  ", depth)
		nodeName := fmt.Sprintf("%s[Group]", indentation)
		buf.WriteString(lipgloss.PlaceHorizontal(18, lipgloss.Left, nodeName))
		buf.WriteString(lipgloss.PlaceHorizontal(11, lipgloss.Right, ""))
		buf.WriteString(lipgloss.PlaceHorizontal(11, lipgloss.Right, ""))

		isSelected := depth == m.depthCursor && slices.Contains(m.Cursor, node)

		// Display iterations
		iterations := fmt.Sprintf("%d", node.Iterations)
		if isSelected && m.Focus {
			buf.WriteString(lipgloss.PlaceHorizontal(11, lipgloss.Right, colors.SelectedColor.Render(iterations)))
		} else {
			buf.WriteString(lipgloss.PlaceHorizontal(11, lipgloss.Right, colors.NumberColor.Render(iterations)))
		}
		buf.WriteString("\n")

		// Render child nodes
		for _, childNode := range node.Nodes {
			m.renderNode(buf, childNode, depth+1)
		}
	} else if node.IsEndNode() {
		songSection := node.Section
		indentation := strings.Repeat("  ", depth)
		section := fmt.Sprintf("%s%s", indentation, (*m.parts)[songSection.Part].GetName())

		isSelected := len(m.Cursor)-1 == m.depthCursor &&
			m.Cursor[len(m.Cursor)-1] == node

		// Check if this is the currently playing section
		var sectionOutput string
		if m.Cursor.Matches(node) {
			sectionOutput = fmt.Sprintf("%s%s", section, colors.CurrentlyPlayingDot)
		} else {
			sectionOutput = section
		}

		buf.WriteString(lipgloss.PlaceHorizontal(18, lipgloss.Left, sectionOutput))

		startBeat := songSection.StartBeat
		if isSelected && m.Focus && m.oldCursor.attribute == SECTION_START_BEAT {
			buf.WriteString(lipgloss.PlaceHorizontal(11, lipgloss.Right, colors.SelectedColor.Render(fmt.Sprintf("%d", startBeat))))
		} else {
			buf.WriteString(lipgloss.PlaceHorizontal(11, lipgloss.Right, colors.NumberColor.Render(fmt.Sprintf("%d", startBeat))))
		}

		startCycle := songSection.StartCycles
		if isSelected && m.Focus && m.oldCursor.attribute == SECTION_START_CYCLE {
			buf.WriteString(lipgloss.PlaceHorizontal(11, lipgloss.Right, colors.SelectedColor.Render(fmt.Sprintf("%d", startCycle))))
		} else {
			buf.WriteString(lipgloss.PlaceHorizontal(11, lipgloss.Right, colors.NumberColor.Render(fmt.Sprintf("%d", startCycle))))
		}

		var cyclesString string
		if songSection.infinite {
			cyclesString = "∞"
		} else {
			cyclesString = fmt.Sprintf("%d", songSection.Cycles)
		}
		if isSelected && m.Focus && m.oldCursor.attribute == SECTION_CYCLES {
			buf.WriteString(lipgloss.PlaceHorizontal(11, lipgloss.Right, colors.SelectedColor.Render(cyclesString)))
		} else {
			buf.WriteString(lipgloss.PlaceHorizontal(11, lipgloss.Right, colors.NumberColor.Render(cyclesString)))
		}

		keepCycles := songSection.KeepCycles
		if isSelected && m.Focus && m.oldCursor.attribute == SECTION_KEEP_CYCLES {
			buf.WriteString(lipgloss.PlaceHorizontal(11, lipgloss.Right, colors.SelectedColor.Render(fmt.Sprintf("%v", keepCycles))))
		} else {
			buf.WriteString(lipgloss.PlaceHorizontal(11, lipgloss.Right, colors.NumberColor.Render(fmt.Sprintf("%v", keepCycles))))
		}

		buf.WriteString("\n")
	} else {
		for _, childNode := range node.Nodes {
			m.renderNode(buf, childNode, depth+1)
		}
	}
}

func (p Part) GetName() string {
	return p.Name
}

func (ac ArrCursor) PlayStateView(cycles int) string {
	var buf strings.Builder
	buf.WriteString("    ▶ ")
	buf.WriteString("")
	for i, arr := range ac {
		if i == 0 {
			if arr.playingIterations == math.MaxInt64 {
				buf.WriteString("∞ ⬩ ")
			}
			continue
		} else if i != 1 {
			buf.WriteString(" ⬩ ")
		}
		arr.PlayStateView(&buf, cycles)
	}
	buf.WriteString("\n")
	return buf.String()
}

func (arr Arrangement) PlayStateView(buf *strings.Builder, cycles int) {
	if arr.IsGroup() {
		fmt.Fprintf(buf, "%d/%d", arr.playingIterations, arr.Iterations)
	} else {
		if arr.Section.infinite {
			fmt.Fprintf(buf, "%d/∞", cycles)
		} else {
			fmt.Fprintf(buf, "%d/%d", cycles, arr.Section.Cycles)
		}
	}
}
