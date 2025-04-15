package arrangement

import (
	"fmt"
	"math"
	"slices"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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

func (ac *ArrCursor) DeleteGroup(depthCursor int) {
	currentNode := (*ac)[depthCursor]
	parentNode := (*ac)[depthCursor-1]

	currentIndex := slices.Index(parentNode.Nodes, currentNode)

	parentNode.Nodes = slices.Replace(parentNode.Nodes, currentIndex, currentIndex+1, currentNode.Nodes...)
	*ac = slices.Delete(*ac, depthCursor, depthCursor+1)
}

func (ac *ArrCursor) DeleteNode() {
	if (*ac)[0].CountEndNodes() <= 1 {
		(*ac).MoveNext()
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
		if !ac.MovePrev() {
			ac.MoveNext()
		}
		parentNode.Nodes = slices.Delete(parentNode.Nodes, currentIndex, currentIndex+1)
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
	SECTION_KEEP_CYCLES
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

type NewPart struct {
	Index     int
	After     bool
	IsPlaying bool
}

func (m *Model) NewPart(index int, after bool, isPlaying bool) {
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

	m.AddPart(after, newNode, isPlaying)
}

type ChangePart struct {
	Index int
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

func (m *Model) AddPart(after bool, newNode *Arrangement, isPlaying bool) {
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

	if !isPlaying {
		newCursor := make(ArrCursor, len(m.Cursor)-1)
		copy(newCursor, m.Cursor[:len(m.Cursor)-1])
		newCursor = append(newCursor, newNode)
		m.Cursor = newCursor
	}
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

func MoveNodeDown(cursor *ArrCursor) {
	if len(*cursor) < 2 {
		return
	}

	currentNode := (*cursor)[len(*cursor)-1]
	parentNode := (*cursor)[len(*cursor)-2]

	currentIndex := slices.Index(parentNode.Nodes, currentNode)

	if currentIndex < len(parentNode.Nodes)-1 && parentNode.Nodes[currentIndex+1].IsEndNode() {
		parentNode.Nodes[currentIndex], parentNode.Nodes[currentIndex+1] = parentNode.Nodes[currentIndex+1], parentNode.Nodes[currentIndex]
		return
	}

	if currentIndex < len(parentNode.Nodes)-1 && parentNode.Nodes[currentIndex+1].IsGroup() {
		newParent := parentNode.Nodes[currentIndex+1]
		newParent.Nodes = slices.Insert(newParent.Nodes, 0, currentNode)
		parentNode.Nodes = slices.Delete(parentNode.Nodes, currentIndex, currentIndex+1)
		*cursor = (*cursor)[:len(*cursor)-1]
		*cursor = append(*cursor, newParent, currentNode)
		return
	}

	if len(*cursor) == 2 && (*cursor)[0] == parentNode {
		// PARENT IS ROOT, CAN'T MOVE DOWN FURTHER
		return
	}

	parentNode.Nodes = parentNode.Nodes[:len(parentNode.Nodes)-1]

	grandparentNode := (*cursor)[len(*cursor)-3]
	parentIndex := slices.Index(grandparentNode.Nodes, parentNode)

	if len(parentNode.Nodes) == 0 {
		grandparentNode.Nodes = slices.Delete(grandparentNode.Nodes, parentIndex, parentIndex+1)
		grandparentNode.Nodes = slices.Insert(grandparentNode.Nodes, parentIndex, currentNode)
	} else {
		grandparentNode.Nodes = slices.Insert(grandparentNode.Nodes, parentIndex+1, currentNode)
	}

	*cursor = (*cursor)[:len(*cursor)-2]   // Remove current node and parent
	*cursor = append(*cursor, currentNode) // Add new path
}

func MoveNodeUp(cursor *ArrCursor) {
	if len(*cursor) < 2 {
		return
	}

	currentNode := (*cursor)[len(*cursor)-1]
	parentNode := (*cursor)[len(*cursor)-2]

	currentIndex := slices.Index(parentNode.Nodes, currentNode)

	if currentIndex == 0 {
		if len(*cursor) > 2 {
			grandparentNode := (*cursor)[len(*cursor)-3]
			parentIndex := slices.Index(grandparentNode.Nodes, parentNode)

			parentNode.Nodes = slices.Delete(parentNode.Nodes, currentIndex, currentIndex+1)
			if len(parentNode.Nodes) == 0 {
				grandparentNode.Nodes = slices.Delete(grandparentNode.Nodes, parentIndex, parentIndex+1)
			}

			grandparentNode.Nodes = slices.Insert(grandparentNode.Nodes, parentIndex, currentNode)

			*cursor = (*cursor)[:len(*cursor)-2]
			*cursor = append(*cursor, currentNode)
		}
		return
	}

	// Swap with previous sibling if both are end nodes
	if parentNode.Nodes[currentIndex-1].IsEndNode() {
		parentNode.Nodes[currentIndex], parentNode.Nodes[currentIndex-1] = parentNode.Nodes[currentIndex-1], parentNode.Nodes[currentIndex]
		return
	}

	// Handle the case where the previous sibling is a group
	if parentNode.Nodes[currentIndex-1].IsGroup() {
		prevGroup := parentNode.Nodes[currentIndex-1]

		// Remove node from current position
		parentNode.Nodes = slices.Delete(parentNode.Nodes, currentIndex, currentIndex+1)

		// Add the node to the end of the previous group
		prevGroup.Nodes = append(prevGroup.Nodes, currentNode)

		// Update cursor to point to the node in its new position
		*cursor = (*cursor)[:len(*cursor)-1]              // Remove current node
		*cursor = append(*cursor, prevGroup, currentNode) // Add path through group
		return
	}
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
	CursorUp     key.Binding
	CursorDown   key.Binding
	CursorLeft   key.Binding
	CursorRight  key.Binding
	Increase     key.Binding
	Decrease     key.Binding
	GroupNodes   key.Binding
	DeleteNode   key.Binding
	MovePartDown key.Binding
	MovePartUp   key.Binding
	Escape       key.Binding
	RenamePart   key.Binding
}

var keys = keymap{
	CursorUp:     Key("Up", "k"),
	CursorDown:   Key("Down", "j"),
	CursorLeft:   Key("Left", "h"),
	CursorRight:  Key("Right", "l"),
	Increase:     Key("Increase", "+", "="),
	Decrease:     Key("Decrease", "-", "_"),
	GroupNodes:   Key("Group", "g"),
	DeleteNode:   Key("Delete", "d"),
	MovePartDown: Key("Move Part Down", "J"),
	MovePartUp:   Key("Move Part Up", "K"),
	RenamePart:   Key("Rename Part", "R"),
	Escape:       Key("Escape", "esc", "enter"),
}

func Key(help string, keyboardKey ...string) key.Binding {
	return key.NewBinding(key.WithKeys(keyboardKey...), key.WithHelp(keyboardKey[0], help))
}

func IsSectionChangeMessage(msg tea.Msg, isEndNode bool) bool {
	if !isEndNode {
		return false
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return key.Matches(msg, keys.Increase, keys.Decrease)
	}
	return false
}

func IsGroupChangeMessage(msg tea.Msg, isGroup bool) bool {
	if !isGroup {
		return false
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return key.Matches(msg, keys.Increase, keys.Decrease)
	}
	return false
}

var changeKeys = []key.Binding{
	keys.GroupNodes,
	keys.DeleteNode,
	keys.MovePartDown,
	keys.MovePartUp,
}

func IsArrChangeMessage(msg tea.Msg) bool {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return key.Matches(msg, changeKeys...)
	case NewPart:
		return true
	case ChangePart:
		return true
	}
	return false
}

type ArrUndo struct {
	Undo Undoable
	Redo Undoable
}

type Undoable interface {
	ApplyUndo(m *Model)
}

type TreeUndo struct {
	undoTree UndoTree
}

func (tu TreeUndo) ApplyUndo(m *Model) {
	m.Root = Convert(tu.undoTree)
}

func Convert(ut UndoTree) *Arrangement {
	if ut.arrRef == nil {
		return nil
	}
	ut.arrRef.Nodes = []*Arrangement{}
	for _, n := range ut.nodes {
		ut.arrRef.Nodes = append(ut.arrRef.Nodes, Convert(n))
	}

	return ut.arrRef
}

type UndoTree struct {
	arrRef *Arrangement
	nodes  []UndoTree
}

type GroupUndo struct {
	arr        *Arrangement
	iterations int
}

func (gu GroupUndo) ApplyUndo(m *Model) {
	(*gu.arr).Iterations = gu.iterations
}

type SectionUndo struct {
	arr     *Arrangement
	section SongSection
}

func (su SectionUndo) ApplyUndo(m *Model) {
	(*su.arr).Section = su.section
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var undo, redo Undoable
	if IsArrChangeMessage(msg) {
		undo = TreeUndo{m.CreateUndoTree()}
	} else if IsSectionChangeMessage(msg, m.Cursor[m.depthCursor].IsEndNode()) {
		undo = SectionUndo{m.Cursor[len(m.Cursor)-1], m.Cursor[len(m.Cursor)-1].Section}
	} else if IsGroupChangeMessage(msg, m.Cursor[m.depthCursor].IsGroup()) {
		undo = GroupUndo{m.Cursor[m.depthCursor], m.Cursor[m.depthCursor].Iterations}
	}

	switch msg := msg.(type) {
	case ChangePart:
		m.ChangePart(msg.Index)
	case NewPart:
		m.NewPart(msg.Index, msg.After, msg.IsPlaying)
	case tea.KeyMsg:
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
			if m.oldCursor.attribute < 3 {
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
				case SECTION_KEEP_CYCLES:
					currentNode.Section.ToggleKeepCycles()
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
				case SECTION_KEEP_CYCLES:
					currentNode.Section.ToggleKeepCycles()
				}
			} else {
				m.Cursor[m.depthCursor].DecreaseIterations()
			}
		case Is(msg, keys.GroupNodes):
			m.GroupNodes()
		case Is(msg, keys.DeleteNode):
			if m.depthCursor == len(m.Cursor)-1 {
				m.Cursor.DeleteNode()
			} else {
				m.Cursor.DeleteGroup(m.depthCursor)
			}
			m.ResetDepth()
		case Is(msg, keys.MovePartDown):
			MoveNodeDown(&m.Cursor)
			m.ResetDepth()
		case Is(msg, keys.MovePartUp):
			MoveNodeUp(&m.Cursor)
			m.ResetDepth()
		case Is(msg, keys.Escape):
			m.Focus = false
			m.ResetDepth()
			return m, func() tea.Msg { return GiveBackFocus{} }
		case Is(msg, keys.RenamePart):
			return m, func() tea.Msg { return RenamePart{} }
		}
	}

	if IsArrChangeMessage(msg) {
		redo = TreeUndo{m.CreateUndoTree()}
		return m, m.CreateUndoCmd(undo, redo)
	} else if IsSectionChangeMessage(msg, m.Cursor[m.depthCursor].IsEndNode()) {
		redo = SectionUndo{m.Cursor[len(m.Cursor)-1], m.Cursor[len(m.Cursor)-1].Section}
		return m, m.CreateUndoCmd(undo, redo)
	} else if IsGroupChangeMessage(msg, m.Cursor[m.depthCursor].IsGroup()) {
		redo = GroupUndo{m.Cursor[m.depthCursor], m.Cursor[m.depthCursor].Iterations}
		return m, m.CreateUndoCmd(undo, redo)
	}

	return m, nil
}

type GiveBackFocus struct{}
type RenamePart struct{}

func Is(msg tea.KeyMsg, k key.Binding) bool {
	return key.Matches(msg, k)
}

type SongSection struct {
	Part        int
	Cycles      int
	StartBeat   int
	StartCycles int
	KeepCycles  bool
	playCycles  int
	infinite    bool
}

func (ss *SongSection) ResetPlayCycles() {
	ss.playCycles = ss.StartCycles
}

func (ss *SongSection) DuringPlayReset() {
	if !ss.KeepCycles {
		ss.playCycles = ss.StartCycles
	}
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
		KeepCycles:  false,
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

func (ss *SongSection) ToggleKeepCycles() {
	ss.KeepCycles = !ss.KeepCycles
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

func (m *Model) ApplyArrUndo(arrUndo Undoable) {
	arrUndo.ApplyUndo(m)
}

func (m Model) CreateUndoTree() UndoTree {
	return CreateUndoTree(m.Root)
}

func (m Model) CreateUndoCmd(undo Undoable, redo Undoable) tea.Cmd {
	return func() tea.Msg {
		return ArrUndo{
			Undo: undo,
			Redo: redo,
		}
	}
}

func CreateUndoTree(arr *Arrangement) UndoTree {
	undoTree := UndoTree{arrRef: arr, nodes: make([]UndoTree, 0)}
	for _, arrRef := range arr.Nodes {
		undoTree.nodes = append(undoTree.nodes, CreateUndoTree(arrRef))
	}
	return undoTree
}
