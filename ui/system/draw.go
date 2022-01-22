package system

import (
	"github.com/gdamore/tcell/v2"
)

// Draw draws this primitive onto the screen.
func (sys *System) Draw(screen tcell.Screen) {
	sys.Box.DrawForSubclass(screen, sys)
	x, y, width, height := sys.GetInnerRect()
	sys.textview.SetRect(x, y, width, height)
	sys.textview.Draw(screen)

	x, y, width, height = sys.textview.GetInnerRect()
	// error dialog
	if sys.errorDialog.IsDisplay() {
		sys.errorDialog.SetRect(x, y, width, height)
		sys.errorDialog.Draw(screen)
		return
	}
	// command dialog dialog
	if sys.cmdDialog.IsDisplay() {
		sys.cmdDialog.SetRect(x, y, width, height)
		sys.cmdDialog.Draw(screen)
		return
	}
	// confirm dialog
	if sys.confirmDialog.IsDisplay() {
		sys.confirmDialog.SetRect(x, y, width, height)
		sys.confirmDialog.Draw(screen)
		return
	}
	// message dialog
	if sys.messageDialog.IsDisplay() {
		sys.messageDialog.SetRect(x, y-2, width, height+4)
		sys.messageDialog.Draw(screen)
		return
	}
	// disk usage dialog
	if sys.dfDialog.IsDisplay() {
		sys.dfDialog.SetRect(x, y, width, height)
		sys.dfDialog.Draw(screen)
		return
	}
	// progress dialog
	if sys.progressDialog.IsDisplay() {
		sys.progressDialog.SetRect(x, y, width, height)
		sys.progressDialog.Draw(screen)
	}
}
