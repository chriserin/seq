package playstate

import (
	"fmt"
	"strings"

	"github.com/chriserin/seq/internal/arrangement"
)

func View(playState PlayState, cursor arrangement.ArrCursor) string {

	var buf strings.Builder

	buf.WriteString(" ▶ ")
	for i, arr := range cursor {
		if i == 0 {
			if arr == playState.LoopedArrangement {
				buf.WriteString("∞ ⬩ ")
			}
			continue
		} else if i != 1 {
			buf.WriteString(" ⬩ ")
		}
		ArrView(&buf, playState, arr)
	}
	buf.WriteString("\n")
	return buf.String()
}

func ArrView(buf *strings.Builder, playState PlayState, a *arrangement.Arrangement) {
	if a == playState.LoopedArrangement {
		fmt.Fprintf(buf, "∞/%d", a.Iterations)
	} else {
		if a.IsGroup() {
			fmt.Fprintf(buf, "%d/%d", (*playState.Iterations)[a]+1, a.Iterations)
		} else {
			fmt.Fprintf(buf, "%d/%d", (*playState.Iterations)[a], a.Section.Cycles)
		}
	}
}
