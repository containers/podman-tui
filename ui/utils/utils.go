package utils

import (
	"strings"
	"time"

	"github.com/containers/storage/pkg/unshare"
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
	// RefreshInterval application refresh interval
	RefreshInterval = 1000 * time.Millisecond
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

// ResolveHomeDir converts a path referencing the home directory via "~"
// to an absolute path
func ResolveHomeDir(path string) (string, error) {
	// check if the path references the home dir to avoid work
	// don't use strings.HasPrefix(path, "~") as this doesn't match "~" alone
	// use strings.HasPrefix(...) to not match "something/~/something"
	if !(path == "~" || strings.HasPrefix(path, "~/")) {
		// path does not reference home dir -> Nothing to do
		return path, nil
	}

	// only get HomeDir when necessary
	home, err := unshare.HomeDir()
	if err != nil {
		return "", err
	}

	// replace the first "~" (start of path) with the HomeDir to resolve "~"
	return strings.Replace(path, "~", home, 1), nil
}
