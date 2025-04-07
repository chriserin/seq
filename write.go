package main

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/charmbracelet/log"

	"github.com/chriserin/seq/internal/arrangement"
	"github.com/chriserin/seq/internal/grid"
	"github.com/chriserin/seq/internal/overlays"
)

// Write saves the parts and arrangement attributes of the model's definition struct to a file
// using a custom format that is easy to diff with tools like git diff
func Write(m *model, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		log.Error("Failed to create file", "filename", filename, "error", err)
		return err
	}
	defer f.Close()

	if m.definition.parts == nil || len(*m.definition.parts) == 0 {
		log.Warn("No parts to write", "filename", filename)
		return nil
	}

	// Write parts first
	if err := writeParts(f, *m.definition.parts); err != nil {
		return err
	}

	// Write arrangement next
	if m.definition.arrangement != nil {
		if err := writeArrangement(f, m.definition.arrangement); err != nil {
			return err
		}
	}

	return nil
}

// writeParts writes all parts to the provided writer
func writeParts(w io.Writer, parts []arrangement.Part) error {
	fmt.Fprintln(w, "--------------------------- PARTS ---------------------------")

	for i, part := range parts {
		partName := part.Name
		if partName == "" {
			partName = fmt.Sprintf("Part%d", i+1)
		}

		separator := fmt.Sprintf("------------------------ PART %s ------------------------", partName)
		fmt.Fprintln(w, separator)
		fmt.Fprintf(w, "Name: %s\n", part.Name)
		fmt.Fprintf(w, "Beats: %d\n", part.Beats)

		if err := writeOverlays(w, part.Overlays); err != nil {
			return err
		}

		fmt.Fprintln(w, "") // Extra newline for better readability between parts
	}

	return nil
}

// writeArrangement writes the arrangement tree structure recursively
func writeArrangement(w io.Writer, arr *arrangement.Arrangement) error {
	if arr == nil {
		return nil
	}

	fmt.Fprintln(w, "------------------------ ARRANGEMENT ------------------------")
	
	// Write arrangement tree recursively using depth-first traversal
	return writeArrangementNode(w, arr, 0, -1) // Pass -1 as childIndex to indicate root node
}

// writeArrangementNode writes a single arrangement node and its children
func writeArrangementNode(w io.Writer, node *arrangement.Arrangement, depth int, childIndex ...int) error {
	// Create indentation based on depth
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	// Write node identifier
	if depth == 0 {
		fmt.Fprintln(w, "------------------------ ROOT NODE ------------------------")
	} else {
		isGroup := len(node.Nodes) > 0
		nodeType := "PART"
		if isGroup {
			nodeType = "GROUP"
		}
		
		// Include child index in the separator if provided
		indexStr := ""
		if len(childIndex) > 0 && childIndex[0] >= 0 {
			indexStr = fmt.Sprintf(" #%d", childIndex[0]+1)
		}
		
		fmt.Fprintf(w, "%s------------------------ %s%s NODE ------------------------\n", 
			indent, nodeType, indexStr)
	}

	// Write node properties
	fmt.Fprintf(w, "%sIterations: %d\n", indent, node.Iterations)
	
	// If it's an end node (contains a section), write section data
	if len(node.Nodes) == 0 {
		fmt.Fprintf(w, "%sPart: %d\n", indent, node.Section.Part)
		fmt.Fprintf(w, "%sCycles: %d\n", indent, node.Section.Cycles)
		fmt.Fprintf(w, "%sStartBeat: %d\n", indent, node.Section.StartBeat)
		fmt.Fprintf(w, "%sStartCycles: %d\n", indent, node.Section.StartCycles)
		fmt.Fprintf(w, "%sKeepCycles: %t\n", indent, node.Section.KeepCycles)
	}

	// Write child nodes recursively
	if len(node.Nodes) > 0 {
		fmt.Fprintf(w, "%sChildren: %d\n", indent, len(node.Nodes))
		fmt.Fprintf(w, "%s------------------------ CHILDREN ------------------------\n", indent)
		
		for i, child := range node.Nodes {
			// Add empty line between children for readability
			if i > 0 {
				fmt.Fprintf(w, "%s\n", indent)
			}
			
			// Pass the child index to the recursive call
			if err := writeArrangementNode(w, child, depth+1, i); err != nil {
				return err
			}
		}
	}

	return nil
}

// writeOverlays writes the overlay tree structure recursively
func writeOverlays(w io.Writer, overlay *overlays.Overlay) error {
	if overlay == nil {
		return nil
	}

	// Write the current overlay
	fmt.Fprintln(w, "----------------------- OVERLAY -------------------------")

	// Write overlay key properties
	fmt.Fprintf(w, "Shift: %d\n", overlay.Key.Shift)
	fmt.Fprintf(w, "Interval: %d\n", overlay.Key.Interval)
	fmt.Fprintf(w, "Width: %d\n", overlay.Key.Width)
	fmt.Fprintf(w, "StartCycle: %d\n", overlay.Key.StartCycle)

	// Write notes in a formatted way
	fmt.Fprintln(w, "------------------------ NOTES --------------------------")
	if len(overlay.Notes) > 0 {
		// Get all GridKeys and sort them for stable output
		gridKeys := make([]grid.GridKey, 0, len(overlay.Notes))
		for k := range overlay.Notes {
			gridKeys = append(gridKeys, k)
		}

		// Sort GridKeys by Line then Beat for consistent output
		sort.Slice(gridKeys, func(i, j int) bool {
			if gridKeys[i].Line != gridKeys[j].Line {
				return gridKeys[i].Line < gridKeys[j].Line
			}
			return gridKeys[i].Beat < gridKeys[j].Beat
		})

		for _, k := range gridKeys {
			note := overlay.Notes[k]
			fmt.Fprintf(w, "GridKey(%d,%d): AccentIndex=%d, Ratchets={Hits:%d,Length:%d,Span:%d}, Action=%d, GateIndex=%d, WaitIndex=%d\n",
				k.Line, k.Beat,
				note.AccentIndex, note.Ratchets.Hits, note.Ratchets.Length, note.Ratchets.Span,
				note.Action, note.GateIndex, note.WaitIndex)
		}
	} else {
		fmt.Fprintln(w, "(empty)")
	}

	// Recursively process the overlay below this one
	if overlay.Below != nil {
		writeOverlays(w, overlay.Below)
	}

	return nil
}
