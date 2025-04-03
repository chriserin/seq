package arrangement

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	colors "github.com/chriserin/seq/internal/colors"
)

var (
	// Styles for header
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			PaddingRight(2)

	// Title style
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			MarginLeft(1).
			MarginRight(5).
			Width(19)
)

func (m Model) View() string {
	var buf strings.Builder

	// Create stylish header
	header := lipgloss.JoinHorizontal(lipgloss.Top,
		titleStyle.Render("♫ Section"),
		lipgloss.PlaceHorizontal(12, lipgloss.Right, headerStyle.Render("Start ♪")),
		lipgloss.PlaceHorizontal(12, lipgloss.Right, headerStyle.Render("Start ⟳")),
		lipgloss.PlaceHorizontal(12, lipgloss.Right, headerStyle.Render("Cycles")),
		lipgloss.PlaceHorizontal(12, lipgloss.Right, headerStyle.Render("Keep")),
	)

	buf.WriteString(header)
	buf.WriteString("\n")

	m.renderNode(&buf, m.Root, 0, false)

	return buf.String()
}

// Style definitions for node rendering
var (
	groupStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F25D94")).
			Bold(true)

	indentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5AFFFF"))

	nodeRowStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			MarginBottom(0)

	sectionNameStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(true)
)

// Recursively render a node and its children
func (m Model) renderNode(buf *strings.Builder, node *Arrangement, depth int, isLast bool) {
	if node == nil {
		return
	}

	// For non-end nodes (groups), show iterations
	if node.IsGroup() && depth > 0 {
		var indent, nodeName string
		if depth > 1 {
			indent = strings.Repeat("┃ ", max(0, depth-2)) + "┣━"
			indentation := indentStyle.Render(indent)
			nodeName = fmt.Sprintf("%s %s", indentation, groupStyle.Render("Group"))
		} else {
			nodeName = groupStyle.Render("Group")
		}

		row := lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.PlaceHorizontal(22, lipgloss.Left, nodeName),
			lipgloss.PlaceHorizontal(12, lipgloss.Right, ""),
			lipgloss.PlaceHorizontal(12, lipgloss.Right, ""),
		)

		isSelected := depth == m.depthCursor && slices.Contains(m.Cursor, node)

		// Display iterations
		iterations := fmt.Sprintf("%d", node.Iterations)
		iterationsText := ""
		if isSelected && m.Focus {
			iterationsText = colors.SelectedColor.Render(iterations)
		} else {
			iterationsText = colors.NumberColor.Render(iterations)
		}

		row = lipgloss.JoinHorizontal(lipgloss.Top,
			row,
			lipgloss.PlaceHorizontal(12, lipgloss.Right, iterationsText),
		)

		buf.WriteString(nodeRowStyle.Render(row))
		buf.WriteString("\n")

		// Render child nodes
		for i, childNode := range node.Nodes {
			if i == len(node.Nodes)-1 && depth > 0 {
				// Draw a connection line for the last child
				indentStyle = indentStyle.Foreground(lipgloss.Color("#5AFFFF"))
			}
			m.renderNode(buf, childNode, depth+1, len(node.Nodes)-1 == i)
		}
	} else if node.IsEndNode() {
		songSection := node.Section

		// Create fancy indentation with tree-like structure
		var indentation, indentChar, section string
		sectionName := (*m.parts)[songSection.Part].GetName()
		if depth > 1 {
			if isLast {
				indentChar = "┗━"
			} else {
				indentChar = "┣━"
			}
			indentChars := strings.Repeat("┃ ", depth-2) + indentChar
			indentation = indentStyle.Render(indentChars)
			section = fmt.Sprintf("%s %s", indentation, sectionNameStyle.Render(sectionName))
		} else {
			section = sectionNameStyle.Render(sectionName)
		}

		isSelected := len(m.Cursor)-1 == m.depthCursor &&
			m.Cursor[len(m.Cursor)-1] == node

		// Check if this is the currently playing section
		var sectionOutput string
		if m.Cursor.Matches(node) {
			sectionOutput = fmt.Sprintf("%s %s", section, colors.CurrentlyPlayingDot)
		} else {
			sectionOutput = section
		}

		// Start building the row
		row := lipgloss.PlaceHorizontal(22, lipgloss.Left, sectionOutput)

		// Handle start beat
		startBeat := songSection.StartBeat
		startBeatText := ""
		if isSelected && m.Focus && m.oldCursor.attribute == SECTION_START_BEAT {
			startBeatText = colors.SelectedColor.Render(fmt.Sprintf("♪ %d", startBeat))
		} else {
			startBeatText = colors.NumberColor.Render(fmt.Sprintf("♪ %d", startBeat))
		}
		row = lipgloss.JoinHorizontal(lipgloss.Top, row,
			lipgloss.PlaceHorizontal(12, lipgloss.Right, startBeatText))

		// Handle start cycle
		startCycle := songSection.StartCycles
		startCycleText := ""
		if isSelected && m.Focus && m.oldCursor.attribute == SECTION_START_CYCLE {
			startCycleText = colors.SelectedColor.Render(fmt.Sprintf("↺ %d", startCycle))
		} else {
			startCycleText = colors.NumberColor.Render(fmt.Sprintf("↺ %d", startCycle))
		}
		row = lipgloss.JoinHorizontal(lipgloss.Top, row,
			lipgloss.PlaceHorizontal(12, lipgloss.Right, startCycleText))

		// Handle cycles
		var cyclesString string
		if songSection.infinite {
			cyclesString = "∞"
		} else {
			cyclesString = fmt.Sprintf("%d", songSection.Cycles)
		}

		cyclesText := ""
		if isSelected && m.Focus && m.oldCursor.attribute == SECTION_CYCLES {
			cyclesText = colors.SelectedColor.Render(fmt.Sprintf("⟳ %s", cyclesString))
		} else {
			cyclesText = colors.NumberColor.Render(fmt.Sprintf("⟳ %s", cyclesString))
		}
		row = lipgloss.JoinHorizontal(lipgloss.Top, row,
			lipgloss.PlaceHorizontal(12, lipgloss.Right, cyclesText))

		// Handle keep cycles
		var keepCycles string
		if songSection.KeepCycles {
			keepCycles = "-"
		} else {
			keepCycles = "✔"
		}
		keepText := ""
		if isSelected && m.Focus && m.oldCursor.attribute == SECTION_KEEP_CYCLES {
			keepText = colors.SelectedColor.Render(fmt.Sprintf("%s", keepCycles))
		} else {
			keepText = colors.NumberColor.Render(fmt.Sprintf("%s", keepCycles))
		}
		row = lipgloss.JoinHorizontal(lipgloss.Top, row,
			lipgloss.PlaceHorizontal(12, lipgloss.Right, keepText))

		buf.WriteString(nodeRowStyle.Render(row))
		buf.WriteString("\n")
	} else {
		for i, childNode := range node.Nodes {
			m.renderNode(buf, childNode, depth+1, len(node.Nodes)-1 == i)
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
