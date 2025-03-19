package arrangement

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestArrangementStructure tests basic initialization and structure of the arrangement
func TestArrangementStructure(t *testing.T) {
	// Create a basic arrangement tree
	root := &Arrangement{
		Iterations: 1,
		Nodes:      make([]*Arrangement, 0),
	}

	// Add some sections
	section1 := SongSection{Part: 0, Cycles: 1, StartBeat: 0, StartCycles: 1}
	section2 := SongSection{Part: 1, Cycles: 2, StartBeat: 4, StartCycles: 2}
	section3 := SongSection{Part: 2, Cycles: 1, StartBeat: 0, StartCycles: 1}

	// Create nodes
	node1 := &Arrangement{Section: section1, Iterations: 1}
	node2 := &Arrangement{Section: section2, Iterations: 1}
	node3 := &Arrangement{Section: section3, Iterations: 1}

	// Add nodes to root
	root.Nodes = append(root.Nodes, node1, node2)

	// Test tree structure
	assert.Equal(t, 2, len(root.Nodes), "Root should have 2 child nodes")
	assert.Equal(t, 0, len(node1.Nodes), "Node1 should have no children")
	assert.Equal(t, 0, len(node2.Nodes), "Node2 should have no children")

	// Test IsEndNode
	assert.False(t, root.IsEndNode(), "Root should not be an end node")
	assert.True(t, node1.IsEndNode(), "Node1 should be an end node")
	assert.True(t, node2.IsEndNode(), "Node2 should be an end node")

	// Create a group with two nodes
	groupNode := &Arrangement{
		Iterations: 2,
		Nodes:      []*Arrangement{node2, node3},
	}

	// Update root
	root.Nodes = []*Arrangement{node1, groupNode}

	// Test the updated structure
	assert.Equal(t, 2, len(root.Nodes), "Root should have 2 child nodes")
	assert.False(t, groupNode.IsEndNode(), "Group node should not be an end node")
	assert.Equal(t, 2, len(groupNode.Nodes), "Group node should have 2 child nodes")
	assert.Equal(t, 2, groupNode.Iterations, "Group node should have 2 iterations")
}

// TestArrCursor tests the cursor navigation through the arrangement tree
func TestArrCursor(t *testing.T) {
	// Create a basic arrangement tree
	root := &Arrangement{
		Iterations: 1,
		Nodes:      make([]*Arrangement, 0),
	}

	// Create some sections and nodes
	section1 := SongSection{Part: 0, Cycles: 1, StartBeat: 0, StartCycles: 1}
	section2 := SongSection{Part: 1, Cycles: 2, StartBeat: 4, StartCycles: 2}
	section3 := SongSection{Part: 2, Cycles: 1, StartBeat: 0, StartCycles: 1}

	node1 := &Arrangement{Section: section1, Iterations: 1}
	node2 := &Arrangement{Section: section2, Iterations: 1}
	node3 := &Arrangement{Section: section3, Iterations: 1}

	// Add nodes to root in a specific order
	root.Nodes = append(root.Nodes, node1, node2, node3)

	// Initialize cursor at the root
	cursor := ArrCursor{root}

	// Test GetCurrentNode when cursor only has root
	currentNode := cursor.GetCurrentNode()
	assert.Equal(t, root, currentNode, "Current node should be the root")

	// Move to the first child
	cursor = append(cursor, root.Nodes[0])
	currentNode = cursor.GetCurrentNode()
	assert.Equal(t, node1, currentNode, "Current node should be node1")

	// Test IncreaseIterations and DecreaseIterations on a group node
	groupNode := &Arrangement{
		Iterations: 1,
		Nodes:      []*Arrangement{node2, node3},
	}

	groupNode.IncreaseIterations()
	assert.Equal(t, 2, groupNode.Iterations, "Group node should have 2 iterations")

	groupNode.DecreaseIterations()
	assert.Equal(t, 1, groupNode.Iterations, "Group node should have 1 iteration")

	groupNode.DecreaseIterations()
	assert.Equal(t, 1, groupNode.Iterations, "Group node iterations should not go below 1")
}

// TestMoveNext tests the MoveNext functionality of the ArrCursor
func TestMoveNext(t *testing.T) {
	// Create a basic arrangement tree with a more complex structure
	root := &Arrangement{
		Iterations: 1,
		Nodes:      make([]*Arrangement, 0),
	}

	// Create simple nodes
	nodeA := &Arrangement{
		Section:    SongSection{Part: 0, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	nodeB := &Arrangement{
		Section:    SongSection{Part: 1, Cycles: 2, StartBeat: 4, StartCycles: 2},
		Iterations: 1,
	}

	// Create a group node with children
	group1 := &Arrangement{
		Iterations: 2,
		Nodes:      make([]*Arrangement, 0),
	}

	nodeC := &Arrangement{
		Section:    SongSection{Part: 2, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	nodeD := &Arrangement{
		Section:    SongSection{Part: 3, Cycles: 3, StartBeat: 2, StartCycles: 0},
		Iterations: 1,
	}

	// Add children to group1
	group1.Nodes = append(group1.Nodes, nodeC, nodeD)

	// Add nodes to root: nodeA, group1, nodeB
	root.Nodes = append(root.Nodes, nodeA, group1, nodeB)

	// Test 1: Starting at root, navigate through the arrangement
	cursor := ArrCursor{root}

	assert.True(t, cursor.MoveNext())
	assert.Equal(t, nodeA, cursor.GetCurrentNode())
	assert.True(t, cursor.MoveNext())
	assert.Equal(t, nodeC, cursor.GetCurrentNode())
	assert.True(t, cursor.MoveNext())
	assert.Equal(t, nodeD, cursor.GetCurrentNode())
	assert.True(t, cursor.MoveNext())
	assert.Equal(t, nodeB, cursor.GetCurrentNode())
	assert.False(t, cursor.MoveNext())
}

// TestMovePrev tests the MovePrev functionality of the ArrCursor
func TestMovePrev(t *testing.T) {
	// Create a basic arrangement tree with a more complex structure
	root := &Arrangement{
		Iterations: 1,
		Nodes:      make([]*Arrangement, 0),
	}

	// Create simple nodes
	nodeA := &Arrangement{
		Section:    SongSection{Part: 0, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	nodeB := &Arrangement{
		Section:    SongSection{Part: 1, Cycles: 2, StartBeat: 4, StartCycles: 2},
		Iterations: 1,
	}

	// Create a group node with children
	group1 := &Arrangement{
		Iterations: 2,
		Nodes:      make([]*Arrangement, 0),
	}

	nodeC := &Arrangement{
		Section:    SongSection{Part: 2, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	nodeD := &Arrangement{
		Section:    SongSection{Part: 3, Cycles: 3, StartBeat: 2, StartCycles: 0},
		Iterations: 1,
	}

	//Add children to group1
	group1.Nodes = append(group1.Nodes, nodeC, nodeD)

	// Add nodes to root: nodeA, group1, nodeB
	root.Nodes = append(root.Nodes, nodeA, group1, nodeB)

	// Test 1: Starting at nodeB (last child of root), navigate backwards
	cursor := ArrCursor{root, nodeB}

	assert.Equal(t, nodeB, cursor.GetCurrentNode())
	assert.True(t, cursor.MovePrev())
	assert.Equal(t, nodeD, cursor.GetCurrentNode())
	assert.True(t, cursor.MovePrev())
	assert.Equal(t, nodeC, cursor.GetCurrentNode())
	assert.True(t, cursor.MovePrev())
	assert.Equal(t, nodeA, cursor.GetCurrentNode())
}

// TestGroupNodes tests the grouping function for arrangement nodes
func TestGroupNodes(t *testing.T) {
	// Create a basic arrangement tree
	root := &Arrangement{
		Iterations: 1,
		Nodes:      make([]*Arrangement, 0),
	}

	// Create some sections and nodes
	section1 := SongSection{Part: 0, Cycles: 1, StartBeat: 0, StartCycles: 1}
	section2 := SongSection{Part: 1, Cycles: 2, StartBeat: 4, StartCycles: 2}
	section3 := SongSection{Part: 2, Cycles: 1, StartBeat: 0, StartCycles: 1}

	node1 := &Arrangement{Section: section1, Iterations: 1}
	node2 := &Arrangement{Section: section2, Iterations: 1}
	node3 := &Arrangement{Section: section3, Iterations: 1}

	// Add nodes to root in a specific order
	root.Nodes = append(root.Nodes, node1, node2, node3)

	// Group nodes 1 and 2
	GroupNodes(root, 0, 1)

	// Check the structure
	assert.Equal(t, 2, len(root.Nodes), "Root should have 2 nodes after grouping")

	// Get the group node (should be the last one)
	groupNode := root.Nodes[0]
	assert.Equal(t, 2, len(groupNode.Nodes), "Group node should have 2 children")
	assert.Equal(t, node1, groupNode.Nodes[0], "First child of group should be node1")
	assert.Equal(t, node2, groupNode.Nodes[1], "Second child of group should be node2")

	// Group with invalid indices
	originalNodes := len(root.Nodes)
	GroupNodes(root, -1, 1)
	assert.Equal(t, originalNodes, len(root.Nodes), "Invalid grouping should not change structure")

	GroupNodes(root, 0, 10)
	assert.Equal(t, originalNodes, len(root.Nodes), "Invalid grouping should not change structure")
}

// TestSongSectionMethods tests the SongSection helper methods
func TestSongSectionMethods(t *testing.T) {
	section := SongSection{
		Part:        0,
		Cycles:      2,
		StartBeat:   1,
		StartCycles: 1,
	}

	// Test increase methods
	section.IncreaseStartBeats()
	assert.Equal(t, 2, section.StartBeat)

	section.IncreaseStartCycles()
	assert.Equal(t, 2, section.StartCycles)

	section.IncreaseCycles()
	assert.Equal(t, 3, section.Cycles)

	// Test decrease methods
	section.DecreaseStartBeats()
	assert.Equal(t, 1, section.StartBeat)

	section.DecreaseStartCycles()
	assert.Equal(t, 1, section.StartCycles)

	section.DecreaseCycles()
	assert.Equal(t, 2, section.Cycles)

	// Test decrease bounds
	section.DecreaseStartBeats()
	section.DecreaseStartBeats()
	assert.Equal(t, 0, section.StartBeat, "StartBeat should not go below 0")

	section.DecreaseStartCycles()
	section.DecreaseStartCycles()
	assert.Equal(t, 0, section.StartCycles, "StartCycles should not go below 0")

	section.DecreaseCycles()
	section.DecreaseCycles()
	section.DecreaseCycles()
	assert.Equal(t, 0, section.Cycles, "Cycles should not go below 0")

	// Test increase bounds
	for i := 0; i < 130; i++ {
		section.IncreaseStartBeats()
		section.IncreaseStartCycles()
		section.IncreaseCycles()
	}

	assert.Equal(t, 127, section.StartBeat, "StartBeat should not exceed 127")
	assert.Equal(t, 127, section.StartCycles, "StartCycles should not exceed 127")
	assert.Equal(t, 127, section.Cycles, "Cycles should not exceed 127")
}

// TestModelSetup tests the Model's initialization
func TestModelSetup(t *testing.T) {
	// Create parts
	parts := []Part{
		{Name: "Part 1"},
		{Name: "Part 2"},
	}

	// Create a test arrangement
	arr := &Arrangement{
		Iterations: 1,
		Nodes:      make([]*Arrangement, 0),
	}

	// Add some nodes
	node1 := &Arrangement{
		Section:    SongSection{Part: 0, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	node2 := &Arrangement{
		Section:    SongSection{Part: 1, Cycles: 2, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	arr.Nodes = append(arr.Nodes, node1, node2)

	// Create model with the arrangement
	partPtr := &parts
	model := InitModel(arr, partPtr)

	// Test the initialization
	assert.Equal(t, arr, model.root, "Model root should be the arrangement")
	// Check that the cursor has length at least 1
	assert.True(t, len(model.Cursor) >= 1, "Cursor should have at least one element")
	assert.Equal(t, arr, model.Cursor[0], "First cursor element should be root")
}

// TestMatches tests the cursor matching functionality
func TestMatches(t *testing.T) {
	// Create a basic arrangement
	root := &Arrangement{
		Iterations: 1,
		Nodes:      make([]*Arrangement, 0),
	}

	nodeA := &Arrangement{
		Section:    SongSection{Part: 0, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	// Test with a normal cursor
	cursor := ArrCursor{root, nodeA}
	assert.True(t, cursor.Matches(nodeA), "Cursor should match nodeA")
	assert.False(t, cursor.Matches(root), "Cursor should not match root")

	// Test with a different node
	nodeB := &Arrangement{
		Section:    SongSection{Part: 1, Cycles: 2, StartBeat: 4, StartCycles: 2},
		Iterations: 1,
	}
	assert.False(t, cursor.Matches(nodeB), "Cursor should not match unrelated node")

	// Test edge cases
	// Empty cursor
	emptyCursor := ArrCursor{}
	assert.False(t, emptyCursor.Matches(nodeA), "Empty cursor should not match any node")

	// Nil node
	assert.False(t, cursor.Matches(nil), "Cursor should not match nil node")

	// Root-only cursor
	rootCursor := ArrCursor{root}
	assert.True(t, rootCursor.Matches(root), "Root cursor should match root")
	assert.False(t, rootCursor.Matches(nodeA), "Root cursor should not match non-root node")
}

// TestModelView tests the rendering of the arrangement Model
func TestModelView(t *testing.T) {
	// Create parts
	parts := []Part{
		{Name: "Part 1"},
		{Name: "Part 2"},
	}

	// Create a test arrangement
	arr := &Arrangement{
		Iterations: 1,
		Nodes:      make([]*Arrangement, 0),
	}

	// Add some nodes
	node1 := &Arrangement{
		Section:    SongSection{Part: 0, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	// Create a group node
	group := &Arrangement{
		Iterations: 2,
		Nodes: []*Arrangement{
			{
				Section:    SongSection{Part: 1, Cycles: 2, StartBeat: 0, StartCycles: 1},
				Iterations: 1,
			},
		},
	}

	arr.Nodes = append(arr.Nodes, node1, group)

	// Create model with the arrangement
	partPtr := &parts
	model := InitModel(arr, partPtr)

	// Set focus to true to test rendering
	model.Focus = true

	// Test view rendering
	output := model.View(0)

	// Basic validation of output
	assert.Contains(t, output, "Part 1", "Output should contain Part 1")
	assert.Contains(t, output, "[Group]", "Output should contain a group")
	assert.Contains(t, output, "2", "Output should contain iteration count")

	// Test that indentation is working
	var buf strings.Builder
	model.renderNode(&buf, group, 1, 0)
	groupOutput := buf.String()
	assert.Contains(t, groupOutput, "  ", "Group rendering should include indentation")
}

func TestIsFirstChild(t *testing.T) {
	root := &Arrangement{
		Iterations: 1,
		Nodes:      make([]*Arrangement, 0),
	}

	nodeA := &Arrangement{
		Section:    SongSection{Part: 0, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	nodeB := &Arrangement{
		Section:    SongSection{Part: 1, Cycles: 2, StartBeat: 4, StartCycles: 2},
		Iterations: 1,
	}

	group1 := &Arrangement{
		Iterations: 2,
		Nodes:      make([]*Arrangement, 0),
	}

	nodeC := &Arrangement{
		Section:    SongSection{Part: 2, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	nodeD := &Arrangement{
		Section:    SongSection{Part: 3, Cycles: 3, StartBeat: 2, StartCycles: 0},
		Iterations: 1,
	}

	group1.Nodes = append(group1.Nodes, nodeC, nodeD)

	root.Nodes = append(root.Nodes, nodeA, group1, nodeB)

	cursor := ArrCursor{root, nodeA}
	assert.True(t, cursor.IsFirstSibling())

	cursor = ArrCursor{root, nodeB}
	assert.False(t, cursor.IsFirstSibling())

	cursor = ArrCursor{root, group1, nodeC}
	assert.True(t, cursor.IsFirstSibling())

	cursor = ArrCursor{root, group1, nodeD}
	assert.False(t, cursor.IsFirstSibling())
}

func TestIsLastSibling(t *testing.T) {
	// Setup test arrangement tree
	root := &Arrangement{
		Iterations: 1,
		Nodes:      make([]*Arrangement, 0),
	}

	nodeA := &Arrangement{
		Section:    SongSection{Part: 0, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	nodeB := &Arrangement{
		Section:    SongSection{Part: 1, Cycles: 2, StartBeat: 4, StartCycles: 2},
		Iterations: 1,
	}

	group1 := &Arrangement{
		Iterations: 2,
		Nodes:      make([]*Arrangement, 0),
	}

	nodeC := &Arrangement{
		Section:    SongSection{Part: 2, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	nodeD := &Arrangement{
		Section:    SongSection{Part: 3, Cycles: 3, StartBeat: 2, StartCycles: 0},
		Iterations: 1,
	}

	group1.Nodes = append(group1.Nodes, nodeC, nodeD)
	root.Nodes = append(root.Nodes, nodeA, group1, nodeB)

	// Define test cases
	tests := []struct {
		name     string
		cursor   ArrCursor
		expected bool
	}{
		{
			name:     "nodeA is not the last sibling of root",
			cursor:   ArrCursor{root, nodeA},
			expected: false,
		},
		{
			name:     "nodeB is the last sibling of root",
			cursor:   ArrCursor{root, nodeB},
			expected: true,
		},
		{
			name:     "nodeC is not the last sibling of group1",
			cursor:   ArrCursor{root, group1, nodeC},
			expected: false,
		},
		{
			name:     "nodeD is the last sibling of group1",
			cursor:   ArrCursor{root, group1, nodeD},
			expected: true,
		},
		{
			name:     "Single item cursor (root only) should return false",
			cursor:   ArrCursor{root},
			expected: false,
		},
		{
			name:     "Empty cursor should return false",
			cursor:   ArrCursor{},
			expected: false,
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.cursor.IsLastSibling()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsRoot(t *testing.T) {
	root := &Arrangement{
		Iterations: 1,
		Nodes:      make([]*Arrangement, 0),
	}

	nodeA := &Arrangement{
		Section:    SongSection{Part: 0, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	// Test cases
	tests := []struct {
		name     string
		cursor   ArrCursor
		expected bool
	}{
		{
			name:     "Root only cursor should return true",
			cursor:   ArrCursor{root},
			expected: true,
		},
		{
			name:     "Multi-node cursor should return false",
			cursor:   ArrCursor{root, nodeA},
			expected: false,
		},
		{
			name:     "Empty cursor should return false",
			cursor:   ArrCursor{},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.cursor.IsRoot()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestHasParentIterations(t *testing.T) {
	// Setup arrangement with playingIterations
	// Remember: only group nodes have iterations, end nodes (with Sections) don't

	// Root node (group)
	root := &Arrangement{
		Iterations:        2,
		playingIterations: 1, // Has iterations left
		Nodes:             make([]*Arrangement, 0),
	}

	// End node with section
	nodeA := &Arrangement{
		Section: SongSection{Part: 0, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Nodes:   nil, // End node
	}

	// Group node with iterations
	groupWithIter := &Arrangement{
		Iterations:        2,
		playingIterations: 2, // Has iterations left
		Nodes:             make([]*Arrangement, 0),
	}

	// Group node without iterations left
	groupNoIter := &Arrangement{
		Iterations:        2,
		playingIterations: 0, // No iterations left
		Nodes:             make([]*Arrangement, 0),
	}

	// End node in groups
	nodeB := &Arrangement{
		Section: SongSection{Part: 1, Cycles: 2, StartBeat: 4, StartCycles: 2},
		Nodes:   nil, // End node
	}

	nodeC := &Arrangement{
		Section: SongSection{Part: 2, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Nodes:   nil, // End node
	}

	// Set up the tree
	groupWithIter.Nodes = append(groupWithIter.Nodes, nodeB)
	groupNoIter.Nodes = append(groupNoIter.Nodes, nodeC)
	root.Nodes = append(root.Nodes, nodeA, groupWithIter, groupNoIter)

	// Test cases
	tests := []struct {
		name     string
		cursor   ArrCursor
		expected bool
	}{
		{
			name:     "Root only cursor has no parent",
			cursor:   ArrCursor{root},
			expected: false,
		},
		{
			name:     "Parent with playingIterations > 0",
			cursor:   ArrCursor{root, groupWithIter, nodeB},
			expected: true,
		},
		{
			name:     "Parent with playingIterations = 0",
			cursor:   ArrCursor{root, groupNoIter, nodeC},
			expected: false,
		},
		{
			name:     "End node directly under root",
			cursor:   ArrCursor{root, nodeA},
			expected: true, // root has playingIterations > 0
		},
		{
			name:     "Empty cursor has no parent",
			cursor:   ArrCursor{},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.cursor.HasParentIterations()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMoveToFirstSibling(t *testing.T) {
	root := &Arrangement{
		Iterations: 1,
		Nodes:      make([]*Arrangement, 0),
	}

	nodeA := &Arrangement{
		Section:    SongSection{Part: 0, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	nodeB := &Arrangement{
		Section:    SongSection{Part: 1, Cycles: 2, StartBeat: 4, StartCycles: 2},
		Iterations: 1,
	}

	nodeC := &Arrangement{
		Section:    SongSection{Part: 2, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	group1 := &Arrangement{
		Iterations: 1,
		Nodes:      []*Arrangement{nodeB, nodeC},
	}

	root.Nodes = append(root.Nodes, nodeA, group1)

	cursor := ArrCursor{root, group1, nodeC}

	cursorCopy := make(ArrCursor, len(cursor))
	copy(cursorCopy, cursor)

	cursorCopy.MoveToFirstSibling()

	assert.Equal(t, root, cursorCopy[0], "Root should remain the same")
	assert.Equal(t, group1, cursorCopy[1], "Group should remain the same")
	assert.Equal(t, nodeB, cursorCopy[2], "Should now point to the first sibling")
}

func TestUp(t *testing.T) {
	root := &Arrangement{
		Iterations: 1,
		Nodes:      make([]*Arrangement, 0),
	}

	nodeA := &Arrangement{
		Section:    SongSection{Part: 0, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations: 1,
	}

	group1 := &Arrangement{
		Iterations: 2,
		Nodes:      make([]*Arrangement, 0),
	}

	nodeB := &Arrangement{
		Section:    SongSection{Part: 1, Cycles: 2, StartBeat: 4, StartCycles: 2},
		Iterations: 1,
	}

	// Set up the tree
	group1.Nodes = append(group1.Nodes, nodeB)
	root.Nodes = append(root.Nodes, nodeA, group1)

	// Test cases
	tests := []struct {
		name     string
		cursor   ArrCursor
		expected ArrCursor
	}{
		{
			name:     "Move up from nodeB to group1",
			cursor:   ArrCursor{root, group1, nodeB},
			expected: ArrCursor{root, group1},
		},
		{
			name:     "Move up from group1 to root",
			cursor:   ArrCursor{root, group1},
			expected: ArrCursor{root},
		},
		{
			name:     "Move up from root (should remain root)",
			cursor:   ArrCursor{root},
			expected: ArrCursor{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a copy for testing
			cursorCopy := make(ArrCursor, len(tc.cursor))
			copy(cursorCopy, tc.cursor)

			cursorCopy.Up()
			assert.Equal(t, tc.expected, cursorCopy)
		})
	}
}

func TestResetIterations(t *testing.T) {
	root := &Arrangement{
		Iterations:        2,
		playingIterations: 0, // Needs reset
		Nodes:             make([]*Arrangement, 0),
	}

	nodeA := &Arrangement{
		Section:           SongSection{Part: 0, Cycles: 1, StartBeat: 0, StartCycles: 1},
		Iterations:        3,
		playingIterations: 0, // Needs reset
	}

	group1 := &Arrangement{
		Iterations:        1,
		playingIterations: 1, // Already set, should not change
	}

	root.Nodes = append(root.Nodes, nodeA, group1)

	cursor := ArrCursor{root, nodeA}

	cursor.ResetIterations()

	assert.Equal(t, 2, root.playingIterations, "Root playingIterations should be reset to Iterations value")
	assert.Equal(t, 3, nodeA.playingIterations, "NodeA playingIterations should be reset to Iterations value")
	assert.Equal(t, 1, group1.playingIterations, "Group1 playingIterations should remain unchanged")
}

// TestDeleteNodeComplex tests the DeleteNode function in more complex scenarios
func TestDeleteNodeComplex(t *testing.T) {
	// Test case 1: Delete a node from a nested structure
	t.Run("delete node from nested structure", func(t *testing.T) {
		// Create a more complex arrangement tree
		root := &Arrangement{
			Iterations: 1,
			Nodes:      make([]*Arrangement, 0),
		}

		outerGroup := &Arrangement{
			Iterations: 2,
			Nodes:      make([]*Arrangement, 0),
		}

		innerGroup := &Arrangement{
			Iterations: 3,
			Nodes:      make([]*Arrangement, 0),
		}

		nodeA := &Arrangement{
			Section: SongSection{Part: 0, Cycles: 1, StartBeat: 0, StartCycles: 1},
		}

		nodeB := &Arrangement{
			Section: SongSection{Part: 1, Cycles: 2, StartBeat: 1, StartCycles: 0},
		}

		nodeC := &Arrangement{
			Section: SongSection{Part: 2, Cycles: 3, StartBeat: 2, StartCycles: 1},
		}

		innerGroup.Nodes = append(innerGroup.Nodes, nodeA, nodeB)
		outerGroup.Nodes = append(outerGroup.Nodes, innerGroup, nodeC)
		root.Nodes = append(root.Nodes, outerGroup)

		cursor := ArrCursor{root, outerGroup, innerGroup, nodeB}

		cursor.DeleteNode()

		assert.Equal(t, 4, len(cursor), "Cursor should point to previous node")
		assert.Equal(t, nodeA, cursor[3], "Cursor should point to previous node")

		assert.Equal(t, 1, len(innerGroup.Nodes), "InnerGroup should have 1 node left")
		assert.Equal(t, nodeA, innerGroup.Nodes[0], "NodeA should be the only node left in innerGroup")
	})

	t.Run("delete last node in a group", func(t *testing.T) {
		root := &Arrangement{
			Iterations: 1,
			Nodes:      make([]*Arrangement, 0),
		}

		firstNode := &Arrangement{
			Section: SongSection{Part: 1, Cycles: 1, StartBeat: 0, StartCycles: 1},
		}

		group := &Arrangement{
			Iterations: 2,
			Nodes:      make([]*Arrangement, 0),
		}

		lastNode := &Arrangement{
			Section: SongSection{Part: 0, Cycles: 1, StartBeat: 0, StartCycles: 1},
		}

		group.Nodes = append(group.Nodes, lastNode)
		root.Nodes = append(root.Nodes, firstNode, group)

		cursor := ArrCursor{root, group, lastNode}

		cursor.DeleteNode()

		assert.Equal(t, 2, len(cursor), "Cursor should move to previous node")
		assert.Equal(t, firstNode, cursor[1], "Cursor should now point to previous node")

		assert.Equal(t, 1, len(root.Nodes), "root should only have the node")
	})

	// Test case 3: Delete middle node in a flat structure
	t.Run("delete middle node in flat structure", func(t *testing.T) {
		root := &Arrangement{
			Iterations: 1,
			Nodes:      make([]*Arrangement, 0),
		}

		node1 := &Arrangement{
			Section: SongSection{Part: 0, Cycles: 1},
		}
		node2 := &Arrangement{
			Section: SongSection{Part: 1, Cycles: 2},
		}
		node3 := &Arrangement{
			Section: SongSection{Part: 2, Cycles: 3},
		}

		root.Nodes = append(root.Nodes, node1, node2, node3)

		cursor := ArrCursor{root, node2}
		cursor.DeleteNode()

		assert.Equal(t, 2, len(root.Nodes), "Root should have 2 nodes after deletion")
		assert.Equal(t, node1, root.Nodes[0], "First node should remain")
		assert.Equal(t, node3, root.Nodes[1], "Third node should now be second")

		assert.Equal(t, 2, len(cursor), "Cursor should move to first node")
		assert.Equal(t, root, cursor[0], "Cursor should move to first node")
	})
}
