package utils

import "github.com/rivo/tview"

type UiDialog interface {
	tview.Primitive
	IsDisplay() bool
	Hide()
}
