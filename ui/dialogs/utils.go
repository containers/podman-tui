package dialogs

import (
	"strings"
)

const (
	// DialogPadding dialog inner paddign
	DialogPadding = 3

	// DialogFormHeight dialog "Enter"/"Cancel" form height
	DialogFormHeight = 3

	// DialogMinWidth dialog min width
	DialogMinWidth = 40

	// TableHeightOffset table hight offset for border
	TableHeightOffset = 3
)

func getMessageWidth(message string) int {
	var messageWidth int
	for _, msg := range strings.Split(message, "\n") {
		if len(msg) > messageWidth {
			messageWidth = len(msg)
		}
	}
	return messageWidth
}
