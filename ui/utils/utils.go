package utils

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	// IDLength max ID length to display
	IDLength = 12
	// HeavyGreenCheckMark unicode
	HeavyGreenCheckMark = "\u2705"
	// HeavyRedCrossMark unicode
	HeavyRedCrossMark = "\u274C"
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

// EmptyBoxSpace returns simple Box without border with bgColor as background
func EmptyBoxSpace(bgColor tcell.Color) *tview.Box {
	box := tview.NewBox()
	box.SetBackgroundColor(bgColor)
	box.SetBorder(false)
	return box
}
