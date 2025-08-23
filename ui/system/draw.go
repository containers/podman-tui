package system

import (
	"github.com/gdamore/tcell/v2"
)

// Draw draws this primitive onto the screen.
func (sys *System) Draw(screen tcell.Screen) {
	sys.DrawForSubclass(screen, sys)

	sysViewX, sysViewY, sysViewW, sysViewH := sys.GetInnerRect()

	sys.connTable.SetRect(sysViewX, sysViewY, sysViewW, sysViewH)
	sys.refresh(sysViewW)
	sys.connTable.Draw(screen)

	x, y, width, height := sys.connTable.GetInnerRect()

	// message dialog
	if sys.messageDialog.IsDisplay() {
		if sys.messageDialog.IsDisplayFullSize() {
			sys.messageDialog.SetRect(sysViewX, sysViewY, sysViewW, sysViewH)
		} else {
			sys.messageDialog.SetRect(x, y, width, height)
		}

		sys.messageDialog.Draw(screen)

		return
	}

	for _, dialog := range sys.getInnerDialogs(true) {
		if dialog.IsDisplay() {
			dialog.SetRect(x, y, width, height)
			dialog.Draw(screen)

			return
		}
	}
}
