package tvxwidgets

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

const (
	// gauge cell.
	prgCell = "▉"
	// form height.
	dialogFormHeight = 3
	// gauge warning percentage.
	gaugeWarnPc = 60.00
	// gauge critical percentage.
	gaugeCritPc = 85.00
	// gauge min percentage.
	gaugeMinPc = 0.00
	// gauge max percentage.
	gaugeMaxPc = 100
	// dialog padding.
	dialogPadding = 2
	// empty space parts.
	emptySpaceParts = 2
)

// getColorName returns convert tcell color to its name.
func getColorName(color tcell.Color) string {
	for name, c := range tcell.ColorNames {
		if c == color {
			return name
		}
	}

	return ""
}

// getMessageWidth returns width size for dialogs based on messages.
func getMessageWidth(message string) int {
	var messageWidth int
	for _, msg := range strings.Split(message, "\n") {
		if len(msg) > messageWidth {
			messageWidth = len(msg)
		}
	}

	return messageWidth
}
