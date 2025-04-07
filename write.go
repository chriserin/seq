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

// Write saves the parts attribute of the model's definition struct to a file
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

	return writeParts(f, *m.definition.parts)
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
