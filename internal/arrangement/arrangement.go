package arrangement

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chriserin/seq/internal/colors"
	"github.com/chriserin/seq/internal/overlays"
)

// New tree-based arrangement structure
type Arrangement struct {
	Section    SongSection
	Nodes      []*Arrangement
	Iterations int
}

// ArrCursor represents the path through the tree to the current node
type ArrCursor []*Arrangement

// IsEndNode checks if an arrangement is an end node (no children)
func (a *Arrangement) IsEndNode() bool {
	return len(a.Nodes) == 0
}

// GetCurrentNode returns the current node based on the cursor path
func (ac ArrCursor) GetCurrentNode() *Arrangement {
	if len(ac) == 0 {
		return nil
	}
	return ac[len(ac)-1]
}

// MoveNext moves to the next node in the arrangement
func (ac *ArrCursor) MoveNext() bool {
	if len(*ac) == 0 {
		return false
	}

	// If we only have the root node, we can't move next
	if len(*ac) == 1 {
		return false
	}

	currentNode := (*ac)[len(*ac)-1]
	parentNode := (*ac)[len(*ac)-2]

	// Find current node's index in parent's Nodes
	var currentIndex int
	for i, node := range parentNode.Nodes {
		if node == currentNode {
			currentIndex = i
			break
		}
	}

	// If there's a next sibling, move to it
	if currentIndex+1 < len(parentNode.Nodes) {
		// Remove current node from cursor
		*ac = (*ac)[:len(*ac)-1]
		// Add next sibling
		*ac = append(*ac, parentNode.Nodes[currentIndex+1])
		return true
	}

	// Otherwise, try to move up and then to next sibling
	*ac = (*ac)[:len(*ac)-1] // Move up
	
	// If we're at the root after moving up, we can't go further
	if len(*ac) <= 1 {
		return false
	}
	
	return ac.MoveNext() // Try again
}

// MovePrev moves to the previous node in the arrangement
func (ac *ArrCursor) MovePrev() bool {
	if len(*ac) == 0 {
		return false
	}

	// If we only have the root node, we can't move prev
	if len(*ac) == 1 {
		return false
	}

	currentNode := (*ac)[len(*ac)-1]
	parentNode := (*ac)[len(*ac)-2]

	// Find current node's index in parent's Nodes
	var currentIndex int
	for i, node := range parentNode.Nodes {
		if node == currentNode {
			currentIndex = i
			break
		}
	}

	// If there's a previous sibling, move to it
	if currentIndex > 0 {
		// Remove current node from cursor
		*ac = (*ac)[:len(*ac)-1]
		// Add previous sibling
		*ac = append(*ac, parentNode.Nodes[currentIndex-1])
		return true
	}

	// Otherwise, we're the first child, so move up to parent
	*ac = (*ac)[:len(*ac)-1] // Move up
	return false             // No more previous nodes
}

func (ac *ArrCursor) Matches(node *Arrangement) bool {
	if len(*ac) == 0 || node == nil {
		return false
	}
	return (*ac)[len(*ac)-1] == node
}

// GroupNodes groups two end nodes together
func GroupNodes(parent *Arrangement, index1, index2 int) {
	if index1 < 0 || index2 < 0 || index1 >= len(parent.Nodes) || index2 >= len(parent.Nodes) {
		return
	}

	node1 := parent.Nodes[index1]
	node2 := parent.Nodes[index2]

	// Create a new parent node
	newParent := &Arrangement{
		Nodes:      []*Arrangement{node1, node2},
		Iterations: 1,
	}

	// Remove original nodes
	if index1 < index2 {
		parent.Nodes = append(parent.Nodes[:index1], parent.Nodes[index1+1:]...)
		parent.Nodes = append(parent.Nodes[:index2-1], parent.Nodes[index2:]...)
	} else {
		parent.Nodes = append(parent.Nodes[:index2], parent.Nodes[index2+1:]...)
		parent.Nodes = append(parent.Nodes[:index1-1], parent.Nodes[index1:]...)
	}

	// Add new parent
	parent.Nodes = append(parent.Nodes, newParent)
}

// DeleteNode removes the current node and restructures the tree
func (ac *ArrCursor) DeleteNode() {
	if len(*ac) < 2 {
		return // Can't delete root
	}

	currentNode := (*ac)[len(*ac)-1]
	parentNode := (*ac)[len(*ac)-2]

	// Find current node's index in parent's Nodes
	var currentIndex int
	for i, node := range parentNode.Nodes {
		if node == currentNode {
			currentIndex = i
			break
		}
	}

	// Remove current node from parent
	parentNode.Nodes = append(parentNode.Nodes[:currentIndex], parentNode.Nodes[currentIndex+1:]...)

	// Move cursor up one level
	*ac = (*ac)[:len(*ac)-1]
}

// IncreaseIterations increases the iterations count of current node
func (ac ArrCursor) IncreaseIterations() {
	if len(ac) == 0 {
		return
	}
	currentNode := ac[len(ac)-1]
	if currentNode.Iterations < 128 {
		currentNode.Iterations++
	}
}

// DecreaseIterations decreases the iterations count of current node
func (ac ArrCursor) DecreaseIterations() {
	if len(ac) == 0 {
		return
	}
	currentNode := ac[len(ac)-1]
	if currentNode.Iterations > 1 {
		currentNode.Iterations--
	}
}

type cursor struct {
	attribute SectionAttribute
}

func (c cursor) Matches(attribute SectionAttribute) bool {
	return c.attribute == attribute
}

type SectionAttribute int

const (
	SECTION_START_BEAT SectionAttribute = iota
	SECTION_START_CYCLE
	SECTION_CYCLES
)

type Model struct {
	Focus     bool
	Cursor    ArrCursor
	oldCursor cursor
	root      *Arrangement
	parts     *[]Part
}

func InitModel(arrangement *Arrangement, parts *[]Part) Model {
	var root *Arrangement

	if arrangement != nil {
		// Use the provided arrangement tree
		root = arrangement
	} else {
		root = &Arrangement{
			Iterations: 1,
			Nodes:      make([]*Arrangement, 0),
			Section:    SongSection{0, 1, 0, 1},
		}
	}

	// Initialize cursor to point to the root and first node if available
	cursor := ArrCursor{root}

	return Model{
		root:   root,
		Cursor: cursor,
		parts:  parts,
	}
}

type keymap struct {
	CursorUp    key.Binding
	CursorDown  key.Binding
	CursorLeft  key.Binding
	CursorRight key.Binding
	Increase    key.Binding
	Decrease    key.Binding
	GroupNodes  key.Binding
	DeleteNode  key.Binding
	Escape      key.Binding
}

var keys = keymap{
	CursorUp:    Key("Up", "k"),
	CursorDown:  Key("Down", "j"),
	CursorLeft:  Key("Left", "h"),
	CursorRight: Key("Right", "l"),
	Increase:    Key("Increase", "+", "="),
	Decrease:    Key("Decrease", "-", "_"),
	GroupNodes:  Key("Group", "g"),
	DeleteNode:  Key("Delete", "d"),
	Escape:      Key("Escape", "esc", "enter"),
}

func Key(help string, keyboardKey ...string) key.Binding {
	return key.NewBinding(key.WithKeys(keyboardKey...), key.WithHelp(keyboardKey[0], help))
}

func (m Model) Update(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case Is(msg, keys.CursorDown):
		m.Cursor.MoveNext()
	case Is(msg, keys.CursorUp):
		m.Cursor.MovePrev()
	case Is(msg, keys.CursorLeft):
		if m.oldCursor.attribute > 0 {
			m.oldCursor.attribute--
		}
	case Is(msg, keys.CursorRight):
		if m.oldCursor.attribute < 2 {
			m.oldCursor.attribute++
		}
	case Is(msg, keys.Increase):
		currentNode := m.Cursor.GetCurrentNode()
		if currentNode != nil {
			if currentNode.IsEndNode() {
				// For end nodes, modify section properties
				switch m.oldCursor.attribute {
				case SECTION_START_BEAT:
					currentNode.Section.IncreaseStartBeats()
				case SECTION_START_CYCLE:
					currentNode.Section.IncreaseStartCycles()
				case SECTION_CYCLES:
					currentNode.Section.IncreaseCycles()
				}
			} else {
				// For parent nodes, increase iterations
				m.Cursor.IncreaseIterations()
			}
		}
	case Is(msg, keys.Decrease):
		currentNode := m.Cursor.GetCurrentNode()
		if currentNode != nil {
			if currentNode.IsEndNode() {
				// For end nodes, modify section properties
				switch m.oldCursor.attribute {
				case SECTION_START_BEAT:
					currentNode.Section.DecreaseStartBeats()
				case SECTION_START_CYCLE:
					currentNode.Section.DecreaseStartCycles()
				case SECTION_CYCLES:
					currentNode.Section.DecreaseCycles()
				}
			} else {
				// For parent nodes, decrease iterations
				m.Cursor.DecreaseIterations()
			}
		}
	case Is(msg, keys.GroupNodes):
		// Group current node with next sibling if possible
		if len(m.Cursor) >= 2 {
			currentNode := m.Cursor[len(m.Cursor)-1]
			parentNode := m.Cursor[len(m.Cursor)-2]

			// Find current node's index in parent
			var currentIndex int
			for i, node := range parentNode.Nodes {
				if node == currentNode {
					currentIndex = i
					break
				}
			}

			// Group with next node if possible
			if currentIndex+1 < len(parentNode.Nodes) {
				GroupNodes(parentNode, currentIndex, currentIndex+1)
			}
		}
	case Is(msg, keys.DeleteNode):
		m.Cursor.DeleteNode()
	case Is(msg, keys.Escape):
		m.Focus = false
		return m, func() tea.Msg { return GiveBackFocus{} }
	}
	return m, nil
}

type GiveBackFocus struct{}

func Is(msg tea.KeyMsg, k key.Binding) bool {
	return key.Matches(msg, k)
}

type SongSection struct {
	Part        int
	Cycles      int
	StartBeat   int
	StartCycles int
}

func (ss *SongSection) IncreaseStartBeats() {
	newStartBeats := ss.StartBeat + 1
	if newStartBeats < 128 {
		ss.StartBeat = newStartBeats
	}
}

func (ss *SongSection) IncreaseStartCycles() {
	newStartCycle := ss.StartCycles + 1
	if newStartCycle < 128 {
		ss.StartCycles = newStartCycle
	}
}

func (ss *SongSection) IncreaseCycles() {
	newCycle := ss.Cycles + 1
	if newCycle < 128 {
		ss.Cycles = newCycle
	}
}

func (ss *SongSection) DecreaseStartBeats() {
	newStartBeats := ss.StartBeat - 1
	if newStartBeats >= 0 {
		ss.StartBeat = newStartBeats
	}
}

func (ss *SongSection) DecreaseStartCycles() {
	newStartCycle := ss.StartCycles - 1
	if newStartCycle >= 0 {
		ss.StartCycles = newStartCycle
	}
}

func (ss *SongSection) DecreaseCycles() {
	newCycle := ss.Cycles - 1
	if newCycle >= 0 {
		ss.Cycles = newCycle
	}
}

type Part struct {
	Overlays *overlays.Overlay
	Beats    uint8
	Name     string
}

func (m Model) View(currentSongSection int) string {
	var buf strings.Builder
	buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, "Section"))
	buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, "Start Beat"))
	buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, "Start Cycle"))
	buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, "Cycles"))
	buf.WriteString("\n")

	// Recursively render the arrangement tree
	m.renderNode(&buf, m.root, 0, currentSongSection, "")

	return buf.String()
}

// GetCursorPath returns the current cursor path
func (m Model) GetCursorPath() ArrCursor {
	return m.Cursor
}

// Recursively render a node and its children
func (m Model) renderNode(buf *strings.Builder, node *Arrangement, depth int, currentSongSection int, prefix string) {
	if node == nil {
		return
	}

	// For non-end nodes (groups), show iterations
	if !node.IsEndNode() {
		indentation := strings.Repeat("  ", depth)
		nodeName := fmt.Sprintf("%s%s[Group]", indentation, prefix)
		buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, nodeName))
		buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, ""))
		buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, ""))

		// Check if this node is selected
		isSelected := false
		if len(m.Cursor) > 0 && m.Cursor[len(m.Cursor)-1] == node {
			isSelected = true
		}

		// Display iterations
		iterations := fmt.Sprintf("%d", node.Iterations)
		if isSelected && m.Focus {
			buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, colors.SelectedColor.Render(iterations)))
		} else {
			buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, colors.NumberColor.Render(iterations)))
		}
		buf.WriteString("\n")

		// Render child nodes
		for i, childNode := range node.Nodes {
			childPrefix := fmt.Sprintf("%d) ", i+1)
			m.renderNode(buf, childNode, depth+1, currentSongSection, childPrefix)
		}
	} else {
		// For end nodes (song sections), show detailed information
		songSection := node.Section
		indentation := strings.Repeat("  ", depth)
		section := fmt.Sprintf("%s%s%s", indentation, prefix, (*m.parts)[songSection.Part].GetName())

		isSelected := false
		if len(m.Cursor) > 0 && m.Cursor[len(m.Cursor)-1] == node {
			isSelected = true
		}

		// Check if this is the currently playing section
		var sectionOutput string
		if m.Cursor.Matches(node) {
			sectionOutput = fmt.Sprintf("%s%s", colors.CurrentlyPlayingDot, section)
		} else {
			sectionOutput = section
		}

		buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, sectionOutput))

		// Display start beat
		startBeat := songSection.StartBeat
		if isSelected && m.Focus && m.oldCursor.attribute == SECTION_START_BEAT {
			buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, colors.SelectedColor.Render(fmt.Sprintf("%d", startBeat))))
		} else {
			buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, colors.NumberColor.Render(fmt.Sprintf("%d", startBeat))))
		}

		// Display start cycle
		startCycle := songSection.StartCycles
		if isSelected && m.Focus && m.oldCursor.attribute == SECTION_START_CYCLE {
			buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, colors.SelectedColor.Render(fmt.Sprintf("%d", startCycle))))
		} else {
			buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, colors.NumberColor.Render(fmt.Sprintf("%d", startCycle))))
		}

		// Display cycles
		cycles := songSection.Cycles
		if isSelected && m.Focus && m.oldCursor.attribute == SECTION_CYCLES {
			buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, colors.SelectedColor.Render(fmt.Sprintf("%d", cycles))))
		} else {
			buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, colors.NumberColor.Render(fmt.Sprintf("%d", cycles))))
		}

		buf.WriteString("\n")
	}
}

func (p Part) GetName() string {
	return p.Name
}
