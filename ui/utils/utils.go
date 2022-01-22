package utils

import (
	"github.com/gdamore/tcell/v2"
)

const (
	// IDLength max ID length to display
	IDLength = 12
)

// GetColorName returns convert tcell color to its name
func GetColorName(color tcell.Color) string {
	for name, c := range tcell.ColorNames {
		if c == color {
			return name
		}
	}
	return ""
}

// AlignStringListWidth returns max string len in the list.
func AlignStringListWidth(list []string) ([]string, int) {
	var (
		max         = 0
		alignedList []string
	)
	for _, item := range list {
		if len(item) > max {
			max = len(item)
		}
	}
	for _, item := range list {
		if len(item) < max {
			need := max - len(item)
			for i := 0; i < need; i++ {
				item = item + " "
			}
		}
		alignedList = append(alignedList, item)
	}
	return alignedList, max
}
