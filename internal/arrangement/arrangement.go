package arrangement

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chriserin/seq/internal/colors"
	"github.com/chriserin/seq/internal/overlaykey"
	"github.com/chriserin/seq/internal/overlays"
)

// New tree-based arrangement structure
type Arrangement struct {
	Section           SongSection
	Nodes             []*Arrangement
	Iterations        int
	playingIterations int
}

func (a *Arrangement) Reset() {
	a.playingIterations = a.Iterations
}

func (a *Arrangement) DrawDown() {
	if a.playingIterations == math.MaxInt64 {
		return
	} else if a.playingIterations > 0 {
		a.playingIterations--
	}
}

func (a *Arrangement) PlayingIterations() int {
	return a.playingIterations
}

func (a *Arrangement) SetInfinite() {
	a.playingIterations = math.MaxInt64
}

func (a *Arrangement) ResetIterations() {
	if len(a.Nodes) > 0 {
		a.playingIterations = 0
		for _, n := range a.Nodes {
			n.ResetIterations()
		}
	}
}

func (a *Arrangement) ResetCycles() {
	if len(a.Nodes) == 0 {
		a.Section.ResetCycles()
	} else {
		for _, n := range a.Nodes {
			n.ResetCycles()
		}
	}
}

func (a *Arrangement) ResetAllPlayCycles() {
	if len(a.Nodes) == 0 {
		a.Section.ResetPlayCycles()
	} else {
		for _, n := range a.Nodes {
			n.ResetAllPlayCycles()
		}
	}
}

// CountEndNodes recursively counts the total number of end nodes in an arrangement
func (a *Arrangement) CountEndNodes() int {
	if a.IsEndNode() {
		return 1
	}

	count := 0
	for _, node := range a.Nodes {
		count += node.CountEndNodes()
	}
	return count
}

type ArrCursor []*Arrangement

// IsEndNode checks if an arrangement is an end node (no children)
func (a *Arrangement) IsEndNode() bool {
	return len(a.Nodes) == 0
}

func (a *Arrangement) IsGroup() bool {
	return len(a.Nodes) != 0
}

func (a *ArrCursor) IsRoot() bool {
	return len(*a) == 1
}

func (a *ArrCursor) HasParentIterations() bool {
	cursorLength := len(*a)
	if cursorLength < 2 {
		return false
	}

	parent := (*a)[cursorLength-2]
	return parent.playingIterations > 0
}

func (a *ArrCursor) GetParentNode() *Arrangement {
	cursorLength := len(*a)
	if cursorLength < 2 {
		return (*a)[0] // Root node or empty cursor
	}

	parent := (*a)[cursorLength-2]
	return parent
}

func (ac *ArrCursor) MoveToFirstSibling() {
	var workingCursor ArrCursor = make([]*Arrangement, len(*ac))
	copy(workingCursor, *ac)
	workingCursor = workingCursor[:len(workingCursor)-1]
	MoveToFirstChild(ac, &workingCursor)
}

func (ac *ArrCursor) Up() {
	*ac = (*ac)[:len(*ac)-1]
}

func (ac *ArrCursor) ResetIterations() {
	for i := range *ac {
		if (*ac)[i].playingIterations == 0 {
			(*ac)[i].playingIterations = (*ac)[i].Iterations
		}
	}
}

// GetCurrentNode returns the current node based on the cursor path
func (ac ArrCursor) GetCurrentNode() *Arrangement {
	if len(ac) == 0 {
		return nil
	}
	return ac[len(ac)-1]
}

func (ac ArrCursor) GetNextSiblingNode() *Arrangement {
	node := ac[len(ac)-1]
	parentNode := ac[len(ac)-2]

	currentIndex := slices.Index(parentNode.Nodes, node)

	return parentNode.Nodes[currentIndex+1]
}

// MoveNext moves to the next node in the arrangement
func (ac *ArrCursor) MoveNext() bool {
	var workingCursor ArrCursor = make([]*Arrangement, len(*ac))
	copy(workingCursor, *ac)
	if (*ac).GetCurrentNode().IsGroup() {
		return MoveToFirstChild(ac, &workingCursor)
	}
	return MoveToNextEndNode(ac, &workingCursor)
}

func (ac *ArrCursor) MoveToSibling() bool {
	var workingCursor ArrCursor = make([]*Arrangement, len(*ac))
	copy(workingCursor, *ac)
	return MoveToSibling(ac, &workingCursor)
}

func (ac *ArrCursor) MovePrev() bool {
	var workingCursor ArrCursor = make([]*Arrangement, len(*ac))
	copy(workingCursor, *ac)
	if len(*ac) <= 1 {
		return MoveToLastChild(ac, &workingCursor)
	}
	return MoveToPrevEndNode(ac, &workingCursor)
}

func (ac *ArrCursor) Matches(node *Arrangement) bool {
	if len(*ac) == 0 || node == nil {
		return false
	}
	return (*ac)[len(*ac)-1] == node
}

func (ac *ArrCursor) DeleteNode() {
	if (*ac)[0].CountEndNodes() <= 1 {
		return
	}

	currentNode := (*ac)[len(*ac)-1]
	parentNode := (*ac)[len(*ac)-2]

	currentIndex := slices.Index(ac.GetParentNode().Nodes, currentNode)

	if len(parentNode.Nodes) == 1 {
		parentNode.Nodes = []*Arrangement{}
		newCursor := make(ArrCursor, len(*ac))
		copy(newCursor, (*ac))
		newCursor.MovePrev()
		ac.Up()
		ac.DeleteNode()
		*ac = newCursor
	} else {
		ac.MovePrev()
		parentNode.Nodes = append(parentNode.Nodes[:currentIndex], parentNode.Nodes[currentIndex+1:]...)
	}
}

func MoveToNextEndNode(currentCursor *ArrCursor, workingCursor *ArrCursor) bool {
	cursorLength := len(*workingCursor)
	if cursorLength == 0 {
		return false
	}

	var scopeCursor ArrCursor = make([]*Arrangement, len(*workingCursor))
	copy(scopeCursor, *workingCursor)

	if MoveToSibling(currentCursor, workingCursor) {
		return true
	} else {
		scopeCursor = scopeCursor[:len(scopeCursor)-1]
		return MoveToNextEndNode(currentCursor, &scopeCursor)
	}
}

func MoveToSibling(currentCursor *ArrCursor, workingCursor *ArrCursor) bool {
	cursorLength := len(*workingCursor)
	if cursorLength == 0 {
		return false
	}
	leaf := (*workingCursor)[cursorLength-1]
	if cursorLength >= 2 {
		parent := (*workingCursor)[cursorLength-2]
		index := slices.Index(parent.Nodes, leaf)
		if index+1 < len(parent.Nodes) {
			*workingCursor = (*workingCursor)[:len(*workingCursor)-1]
			*workingCursor = append(*workingCursor, parent.Nodes[index+1])
			return MoveToFirstChild(currentCursor, workingCursor)
		} else {
			return false
		}
	} else {
		return false
	}
}

func MoveToFirstChild(currentCursor *ArrCursor, workingCursor *ArrCursor) bool {
	cursorLength := len(*workingCursor)
	leaf := (*workingCursor)[cursorLength-1]
	if leaf.IsEndNode() {
		*currentCursor = *workingCursor
		return true
	} else if len(leaf.Nodes) > 0 {
		*workingCursor = append(*workingCursor, leaf.Nodes[0])
		return MoveToFirstChild(currentCursor, workingCursor)
	} else {
		panic("Malformed arrangement tree")
	}
}

func MoveToPrevEndNode(currentCursor *ArrCursor, workingCursor *ArrCursor) bool {
	cursorLength := len(*workingCursor)
	if cursorLength == 0 {
		return false
	}

	var scopeCursor ArrCursor = make([]*Arrangement, len(*workingCursor))
	copy(scopeCursor, *workingCursor)

	if MoveToPrevSibling(currentCursor, workingCursor) {
		return true
	} else {
		scopeCursor = scopeCursor[:len(scopeCursor)-1]
		return MoveToPrevEndNode(currentCursor, &scopeCursor)
	}
}

func MoveToPrevSibling(currentCursor *ArrCursor, workingCursor *ArrCursor) bool {
	cursorLength := len(*workingCursor)
	if cursorLength == 0 {
		return false
	}
	leaf := (*workingCursor)[cursorLength-1]
	if cursorLength >= 2 {
		parent := (*workingCursor)[cursorLength-2]
		index := slices.Index(parent.Nodes, leaf)
		if index-1 >= 0 {
			*workingCursor = (*workingCursor)[:len(*workingCursor)-1]
			*workingCursor = append(*workingCursor, parent.Nodes[index-1])
			return MoveToLastChild(currentCursor, workingCursor)
		} else {
			return false
		}
	} else {
		return false
	}
}

func MoveToLastChild(currentCursor *ArrCursor, workingCursor *ArrCursor) bool {
	cursorLength := len(*workingCursor)
	leaf := (*workingCursor)[cursorLength-1]
	if leaf.IsEndNode() {
		*currentCursor = *workingCursor
		return true
	} else if len(leaf.Nodes) > 0 {
		*workingCursor = append(*workingCursor, leaf.Nodes[len(leaf.Nodes)-1])
		return MoveToLastChild(currentCursor, workingCursor)
	} else {
		panic("Malformed arrangement tree")
	}
}

func (m Model) CurrentNodeCursor(currentCursor ArrCursor) ArrCursor {
	if m.depthCursor == len(currentCursor)-1 {
		currentNode := currentCursor[len(currentCursor)-1]
		currentNode.Section.infinite = true
		partGroup := &Arrangement{
			Nodes:      []*Arrangement{currentNode},
			Iterations: 1,
		}
		cursor := ArrCursor{partGroup, currentNode}
		return cursor
	} else {
		group := currentCursor[m.depthCursor]
		group.SetInfinite()
		cursor := ArrCursor{group}
		cursor.MoveNext()
		return cursor
	}
}

// GroupNodes groups two end nodes together
func GroupNodes(parent *Arrangement, index1, index2 int) {
	if index1 < 0 || index2 < 0 || index1 >= len(parent.Nodes) || index2 >= len(parent.Nodes) {
		return
	}

	var newParent *Arrangement

	if index1 < index2 {
		node1 := parent.Nodes[index1]
		node2 := parent.Nodes[index2]

		newParent = &Arrangement{
			Nodes:      []*Arrangement{node1, node2},
			Iterations: 1,
		}

		parent.Nodes = append(parent.Nodes[:index1], parent.Nodes[index1+1:]...)
		parent.Nodes = append(parent.Nodes[:index2-1], parent.Nodes[index2:]...)
	} else if index1 == index2 {
		node1 := parent.Nodes[index1]

		newParent = &Arrangement{
			Nodes:      []*Arrangement{node1},
			Iterations: 1,
		}
		parent.Nodes = append(parent.Nodes[:index1], parent.Nodes[index1+1:]...)
	}

	parent.Nodes = slices.Insert(parent.Nodes, index1, newParent)
}

func (arr *Arrangement) IncreaseIterations() {
	if arr.Iterations < 128 {
		arr.Iterations++
	}
}

func (arr *Arrangement) DecreaseIterations() {
	if arr.Iterations > 1 {
		arr.Iterations--
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
	Focus       bool
	SavedCursor ArrCursor
	Cursor      ArrCursor
	oldCursor   cursor
	Root        *Arrangement
	parts       *[]Part
	depthCursor int
}

func (m *Model) ResetDepth() {
	m.depthCursor = len(m.Cursor) - 1
}

func (m *Model) GroupNodes() {
	if len(m.Cursor) >= 2 {
		currentNode := m.Cursor[len(m.Cursor)-1]
		parentNode := m.Cursor[len(m.Cursor)-2]

		currentIndex := slices.Index(m.Cursor.GetParentNode().Nodes, currentNode)

		m.Cursor.MovePrev()
		if currentIndex+1 < len(parentNode.Nodes) {
			GroupNodes(parentNode, currentIndex, currentIndex+1)
		} else {
			GroupNodes(parentNode, currentIndex, currentIndex)
		}
		m.Cursor.MoveNext()
		m.ResetDepth()
	}
}

func (m *Model) NewPart(index int, after bool) {
	partId := index
	if index < 0 {
		partId = len(*m.parts)
		*m.parts = append(*m.parts, InitPart(fmt.Sprintf("Part %d", partId+1)))
	}

	section := InitSongSection(partId)
	newNode := &Arrangement{
		Section:    section,
		Iterations: 1,
	}

	m.AddPart(after, newNode)
}

func (m *Model) ChangePart(index int) {
	partId := index
	if index < 0 {
		partId = len(*m.parts)
		*m.parts = append(*m.parts, InitPart(fmt.Sprintf("Part %d", partId+1)))
	}
	currentNode := m.Cursor[m.depthCursor]
	currentNode.Section.Part = partId
}

func (m *Model) AddPart(after bool, newNode *Arrangement) {
	currentNode := m.Cursor[m.depthCursor]
	parentNode := m.Cursor[m.depthCursor-1]

	currentIndex := slices.Index(parentNode.Nodes, currentNode)
	if after {
		currentIndex++
	}

	if currentIndex > len(parentNode.Nodes) {
		parentNode.Nodes = append(parentNode.Nodes, newNode)
	} else {
		parentNode.Nodes = slices.Insert(parentNode.Nodes, currentIndex, newNode)
	}

	// Update cursor to point to new node
	newCursor := make(ArrCursor, len(m.Cursor)-1)
	copy(newCursor, m.Cursor[:len(m.Cursor)-1])
	newCursor = append(newCursor, newNode)
	m.Cursor = newCursor
}

func (cursor ArrCursor) IsFirstSibling() bool {
	if len(cursor) < 2 {
		return false
	}
	leaf := cursor[len(cursor)-1]
	parent := cursor[len(cursor)-2]
	return slices.Index(parent.Nodes, leaf) == 0
}

func (a *ArrCursor) IsLastSibling() bool {
	cursorLength := len(*a)
	if cursorLength < 2 {
		return false // Root node or empty cursor
	}

	leaf := (*a)[cursorLength-1]
	parent := (*a)[cursorLength-2]

	index := slices.Index(parent.Nodes, leaf)
	return index == len(parent.Nodes)-1
}

func InitModel(arrangement *Arrangement, parts *[]Part) Model {
	var root *Arrangement
	var cursor ArrCursor

	if arrangement != nil {
		// Use the provided arrangement tree
		root = arrangement
		cursor = ArrCursor{root}
		cursor.MoveNext()
	} else {
		root = &Arrangement{
			Iterations: 1,
			Nodes:      make([]*Arrangement, 0),
			Section:    InitSongSection(0),
		}
		cursor = ArrCursor{root}
	}

	// Initialize cursor to point to the root and first node if available

	return Model{
		Root:        root,
		Cursor:      cursor,
		parts:       parts,
		depthCursor: len(cursor) - 1,
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
		if !m.Cursor.IsLastSibling() && m.Cursor.GetNextSiblingNode().IsGroup() {
			m.ResetDepth()
			m.Cursor.MoveNext()
		} else if m.depthCursor < len(m.Cursor)-1 {
			m.depthCursor++
		} else {
			m.Cursor.MoveNext()
			m.ResetDepth()
		}
	case Is(msg, keys.CursorUp):
		if m.Cursor[:m.depthCursor+1].IsFirstSibling() {
			if m.depthCursor > 1 {
				m.depthCursor--
			}
		} else {
			m.Cursor.MovePrev()
			m.ResetDepth()
		}
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
		if m.depthCursor == len(m.Cursor)-1 {
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
			m.Cursor[m.depthCursor].IncreaseIterations()
		}
	case Is(msg, keys.Decrease):
		currentNode := m.Cursor.GetCurrentNode()
		if m.depthCursor+1 == len(m.Cursor) {
			switch m.oldCursor.attribute {
			case SECTION_START_BEAT:
				currentNode.Section.DecreaseStartBeats()
			case SECTION_START_CYCLE:
				currentNode.Section.DecreaseStartCycles()
			case SECTION_CYCLES:
				currentNode.Section.DecreaseCycles()
			}
		} else {
			m.Cursor[m.depthCursor].DecreaseIterations()
		}
	case Is(msg, keys.GroupNodes):
		m.GroupNodes()
	case Is(msg, keys.DeleteNode):
		m.Cursor.DeleteNode()
		m.ResetDepth()
	case Is(msg, keys.Escape):
		m.Focus = false
		m.ResetDepth()
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
	playCycles  int
	infinite    bool
}

func (ss *SongSection) ResetPlayCycles() {
	ss.playCycles = ss.StartCycles
}

func (ss *SongSection) DuringPlayReset() {
	ss.playCycles = ss.StartCycles
}

func (ss SongSection) PlayCycles() int {
	return ss.playCycles
}

func (ss *SongSection) ResetCycles() {
	ss.infinite = false
}

func (ss *SongSection) IncrementPlayCycles() {
	ss.playCycles++
}

func (ss *SongSection) IsDone() bool {
	return !ss.infinite &&
		ss.Cycles+ss.StartCycles <= ss.playCycles
}

func InitSongSection(part int) SongSection {
	return SongSection{
		Part:        part,
		Cycles:      1,
		StartBeat:   0,
		StartCycles: 1,
		infinite:    false,
	}
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

func InitPart(name string) Part {
	return Part{Overlays: overlays.InitOverlay(overlaykey.ROOT, nil), Beats: 32, Name: name}
}

func (m Model) View() string {
	var buf strings.Builder
	buf.WriteString(lipgloss.PlaceHorizontal(20, lipgloss.Left, "  Section"))
	buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, "Start Beat"))
	buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, "Start Cycle"))
	buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, "Cycles"))
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
		buf.WriteString(lipgloss.PlaceHorizontal(20, lipgloss.Left, nodeName))
		buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, ""))
		buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, ""))

		isSelected := depth == m.depthCursor && slices.Contains(m.Cursor, node)

		// Display iterations
		iterations := fmt.Sprintf("%d", node.Iterations)
		if isSelected && m.Focus {
			buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, colors.SelectedColor.Render(iterations)))
		} else {
			buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, colors.NumberColor.Render(iterations)))
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

		buf.WriteString(lipgloss.PlaceHorizontal(20, lipgloss.Left, sectionOutput))

		startBeat := songSection.StartBeat
		if isSelected && m.Focus && m.oldCursor.attribute == SECTION_START_BEAT {
			buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, colors.SelectedColor.Render(fmt.Sprintf("%d", startBeat))))
		} else {
			buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, colors.NumberColor.Render(fmt.Sprintf("%d", startBeat))))
		}

		startCycle := songSection.StartCycles
		if isSelected && m.Focus && m.oldCursor.attribute == SECTION_START_CYCLE {
			buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, colors.SelectedColor.Render(fmt.Sprintf("%d", startCycle))))
		} else {
			buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, colors.NumberColor.Render(fmt.Sprintf("%d", startCycle))))
		}

		var cyclesString string
		if songSection.infinite {
			cyclesString = "∞"
		} else {
			cyclesString = fmt.Sprintf("%d", songSection.Cycles)
		}
		if isSelected && m.Focus && m.oldCursor.attribute == SECTION_CYCLES {
			buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, colors.SelectedColor.Render(cyclesString)))
		} else {
			buf.WriteString(lipgloss.PlaceHorizontal(15, lipgloss.Right, colors.NumberColor.Render(cyclesString)))
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
