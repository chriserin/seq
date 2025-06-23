package arrangement

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	themes "github.com/chriserin/seq/internal/themes"
)

func (m Model) View() string {
	var buf strings.Builder

	// Create stylish header
	header := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.PlaceHorizontal(24, lipgloss.Left, themes.AppTitleStyle.Render("Section "), lipgloss.WithWhitespaceChars("─"), lipgloss.WithWhitespaceForeground(themes.ArrangementSelectedLineColor)),
		lipgloss.PlaceHorizontal(12, lipgloss.Right, themes.AppTitleStyle.Render("Start Beat"), lipgloss.WithWhitespaceChars("─"), lipgloss.WithWhitespaceForeground(themes.ArrangementSelectedLineColor)),
		lipgloss.PlaceHorizontal(12, lipgloss.Right, themes.AppTitleStyle.Render("⟳ Start"), lipgloss.WithWhitespaceChars("─"), lipgloss.WithWhitespaceForeground(themes.ArrangementSelectedLineColor)),
		lipgloss.PlaceHorizontal(12, lipgloss.Right, themes.AppTitleStyle.Render("⟳ Amount"), lipgloss.WithWhitespaceChars("─"), lipgloss.WithWhitespaceForeground(themes.ArrangementSelectedLineColor)),
		lipgloss.PlaceHorizontal(12, lipgloss.Right, themes.AppTitleStyle.Render("⟳ Keep"), lipgloss.WithWhitespaceChars("─"), lipgloss.WithWhitespaceForeground(themes.ArrangementSelectedLineColor)),
	)
	buf.WriteString(header)
	buf.WriteString("\n")

	m.renderNode(&buf, m.Root, 0, false)

	return buf.String()
}

// Style definitions for node rendering

// Recursively render a node and its children
func (m Model) renderNode(buf *strings.Builder, node *Arrangement, depth int, isLast bool) {
	if node == nil {
		return
	}

	// For non-end nodes (groups), show iterations
	if node.IsGroup() && depth > 0 {
		var indent, nodeName string
		if depth > 1 {
			indent = strings.Repeat("│ ", max(0, depth-2)) + "├─"
			indentation := themes.IndentStyle.Render(indent)
			nodeName = fmt.Sprintf("%s %s", indentation, themes.GroupStyle.Render("Group "))
		} else {
			nodeName = themes.GroupStyle.Render("Group")
		}

		isSelected := depth == m.depthCursor && slices.Contains(m.Cursor, node)

		var options []lipgloss.WhitespaceOption
		if isSelected {
			options = []lipgloss.WhitespaceOption{lipgloss.WithWhitespaceChars("─"), lipgloss.WithWhitespaceForeground(themes.ArrangementSelectedLineColor)}
		} else {
			options = []lipgloss.WhitespaceOption{}
		}

		row := lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.PlaceHorizontal(22, lipgloss.Left, nodeName, options...),
			lipgloss.PlaceHorizontal(12, lipgloss.Right, "", options...),
			lipgloss.PlaceHorizontal(12, lipgloss.Right, "", options...),
		)

		// Display iterations
		iterations := fmt.Sprintf("%d", node.Iterations)
		iterationsText := ""
		if isSelected && m.Focus {
			iterationsText = themes.SelectedStyle.MarginLeft(1).Render(iterations)
		} else {
			iterationsText = themes.NumberStyle.Render(iterations)
		}

		row = lipgloss.JoinHorizontal(lipgloss.Top,
			row,
			lipgloss.PlaceHorizontal(12, lipgloss.Right, iterationsText, options...),
			lipgloss.PlaceHorizontal(12, lipgloss.Right, "", options...),
		)

		buf.WriteString(themes.NodeRowStyle.Render(row))
		buf.WriteString("\n")

		// Render child nodes
		for i, childNode := range node.Nodes {
			m.renderNode(buf, childNode, depth+1, len(node.Nodes)-1 == i)
		}
	} else if node.IsEndNode() {
		songSection := node.Section

		// Create fancy indentation with tree-like structure
		var indentation, indentChar, section string
		sectionName := (*m.parts)[songSection.Part].GetName()
		if depth > 1 {
			if isLast {
				indentChar = "└─"
			} else {
				indentChar = "├─"
			}
			indentChars := strings.Repeat("│ ", depth-2) + indentChar
			indentation = themes.IndentStyle.Render(indentChars)
			section = fmt.Sprintf("%s %s", indentation, themes.AppTitleStyle.Bold(false).Render(sectionName))
		} else {
			section = themes.AppTitleStyle.Bold(false).Render(sectionName)
		}

		isSelected := len(m.Cursor)-1 == m.depthCursor &&
			m.Cursor[len(m.Cursor)-1] == node

		var options []lipgloss.WhitespaceOption
		if isSelected {
			options = []lipgloss.WhitespaceOption{lipgloss.WithWhitespaceChars("─"), lipgloss.WithWhitespaceForeground(themes.ArrangementSelectedLineColor)}
		} else {
			options = []lipgloss.WhitespaceOption{}
		}

		// Check if this is the currently playing section
		var sectionOutput string
		if m.Cursor.Matches(node) {
			sectionOutput = fmt.Sprintf("%s %s", section, themes.CurrentlyPlayingSymbol)
		} else {
			sectionOutput = section
		}

		// Start building the row
		row := lipgloss.PlaceHorizontal(22, lipgloss.Left, sectionOutput, options...)

		selectedStyle := themes.SelectedStyle.MarginLeft(1)

		// Handle start beat
		startBeat := songSection.StartBeat
		startBeatText := ""
		if isSelected && m.Focus && m.oldCursor.attribute == SECTION_START_BEAT {
			startBeatText = selectedStyle.Render(fmt.Sprintf("%d", startBeat))
		} else {
			startBeatText = themes.NumberStyle.Render(fmt.Sprintf("%d", startBeat))
		}
		row = lipgloss.JoinHorizontal(lipgloss.Top, row,
			lipgloss.PlaceHorizontal(12, lipgloss.Right, startBeatText, options...))

		// Handle start cycle
		startCycle := songSection.StartCycles
		startCycleText := ""
		if isSelected && m.Focus && m.oldCursor.attribute == SECTION_START_CYCLE {
			startCycleText = selectedStyle.Render(fmt.Sprintf("%d", startCycle))
		} else {
			startCycleText = themes.NumberStyle.Render(fmt.Sprintf("%d", startCycle))
		}
		row = lipgloss.JoinHorizontal(lipgloss.Top, row,
			lipgloss.PlaceHorizontal(12, lipgloss.Right, startCycleText, options...))

		// Handle cycles
		var cyclesString string
		if songSection.infinite {
			cyclesString = "∞"
		} else {
			cyclesString = fmt.Sprintf("%d", songSection.Cycles)
		}

		cyclesText := ""
		if isSelected && m.Focus && m.oldCursor.attribute == SECTION_CYCLES {
			cyclesText = selectedStyle.Render(cyclesString)
		} else {
			cyclesText = themes.NumberStyle.Render(cyclesString)
		}
		row = lipgloss.JoinHorizontal(lipgloss.Top, row,
			lipgloss.PlaceHorizontal(12, lipgloss.Right, cyclesText, options...))

		// Handle keep cycles
		var keepCycles string
		if songSection.KeepCycles {
			keepCycles = "✔"
		} else {
			keepCycles = "-"
		}
		keepText := ""
		if isSelected && m.Focus && m.oldCursor.attribute == SECTION_KEEP_CYCLES {
			keepText = selectedStyle.Render(keepCycles)
		} else {
			keepText = themes.NumberStyle.Render(keepCycles)
		}
		row = lipgloss.JoinHorizontal(lipgloss.Top, row,
			lipgloss.PlaceHorizontal(12, lipgloss.Right, keepText, options...))

		buf.WriteString(themes.NodeRowStyle.Render(row))
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
	buf.WriteString(" ▶ ")
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
		if arr.IsInfinite() {
			fmt.Fprintf(buf, "∞/%d", arr.Iterations)
		} else {
			fmt.Fprintf(buf, "%d/%d", arr.playingIterations, arr.Iterations)
		}
	} else {
		if arr.Section.infinite {
			fmt.Fprintf(buf, "%d/∞", cycles)
		} else {
			fmt.Fprintf(buf, "%d/%d", cycles, arr.Section.Cycles)
		}
	}
}
