package system

import (
	"github.com/gdamore/tcell/v2"
)

// Draw draws this primitive onto the screen.
func (sys *System) Draw(screen tcell.Screen) { //nolint:cyclop
	sys.Box.DrawForSubclass(screen, sys)

	sysViewX, sysViewY, sysViewW, sysViewH := sys.GetInnerRect()

	sys.connTable.SetRect(sysViewX, sysViewY, sysViewW, sysViewH)
	sys.refresh(sysViewW)
	sys.connTable.Draw(screen)

	x, y, width, height := sys.connTable.GetInnerRect()

	// error dialog
	if sys.errorDialog.IsDisplay() {
		sys.errorDialog.SetRect(x, y, width, height)
		sys.errorDialog.Draw(screen)

		return
	}

	// connection progress dialog
	if sys.connPrgDialog.IsDisplay() {
		sys.connPrgDialog.SetRect(x, y, width, height)
		sys.connPrgDialog.Draw(screen)

		return
	}

	// command dialog
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
		if sys.messageDialog.IsDisplayFullSize() {
			sys.messageDialog.SetRect(sysViewX, sysViewY, sysViewW, sysViewH)
		} else {
			sys.messageDialog.SetRect(x, y, width, height)
		}

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

		return
	}

	// connection create dialog
	if sys.connAddDialog.IsDisplay() {
		sys.connAddDialog.SetRect(x, y, width, height)
		sys.connAddDialog.Draw(screen)

		return
	}

	// event dialog
	if sys.eventDialog.IsDisplay() {
		sys.eventDialog.SetRect(sysViewX, sysViewY, sysViewW, sysViewH)
		sys.eventDialog.Draw(screen)
	}
}
