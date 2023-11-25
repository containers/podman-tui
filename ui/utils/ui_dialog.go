package utils

import "github.com/rivo/tview"

type UIDialog interface {
	tview.Primitive
	IsDisplay() bool
	Hide()
}
