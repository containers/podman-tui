package tvxwidgets

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

const (
	// gauge cell
	prgCell = "▉"
	// form height
	dialogFormHeight = 3
)

// getColorName returns convert tcell color to its name
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
