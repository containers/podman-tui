package containers

import (
	"github.com/gdamore/tcell/v2"
)

// Draw draws this primitive onto the screen.
func (cnt *Containers) Draw(screen tcell.Screen) { //nolint:cyclop
	cnt.DrawForSubclass(screen, cnt)
	cnt.SetBorder(false)

	cntViewX, cntViewY, cntViewW, cntViewH := cnt.GetInnerRect()

	cnt.refresh(cntViewW)
	cnt.table.SetRect(cntViewX, cntViewY, cntViewW, cntViewH)
	cnt.table.SetBorder(true)
	cnt.table.Draw(screen)

	x, y, width, height := cnt.table.GetInnerRect()

	// error dialog
	if cnt.errorDialog.IsDisplay() {
		cnt.errorDialog.SetRect(x, y, width, height)
		cnt.errorDialog.Draw(screen)

		return
	}

	// command dialog
	if cnt.cmdDialog.IsDisplay() {
		cnt.cmdDialog.SetRect(x, y, width, height)
		cnt.cmdDialog.Draw(screen)

		return
	}

	// command input dialog
	if cnt.cmdInputDialog.IsDisplay() {
		cnt.cmdInputDialog.SetRect(x, y, width, height)
		cnt.cmdInputDialog.Draw(screen)

		return
	}

	// create dialog
	if cnt.createDialog.IsDisplay() {
		cnt.createDialog.SetRect(x, y, width, height)
		cnt.createDialog.Draw(screen)

		return
	}

	// run dialog
	if cnt.runDialog.IsDisplay() {
		cnt.runDialog.SetRect(x, y, width, height)
		cnt.runDialog.Draw(screen)

		return
	}

	// message dialog
	if cnt.messageDialog.IsDisplay() {
		if cnt.messageDialog.IsDisplayFullSize() {
			cnt.messageDialog.SetRect(cntViewX, cntViewY, cntViewW, cntViewH)
		} else {
			cnt.messageDialog.SetRect(x, y, width, height+1)
		}

		cnt.messageDialog.Draw(screen)

		return
	}

	// confirm dialog
	if cnt.confirmDialog.IsDisplay() {
		cnt.confirmDialog.SetRect(x, y, width, height)
		cnt.confirmDialog.Draw(screen)

		return
	}

	// progress dialog
	if cnt.progressDialog.IsDisplay() {
		cnt.progressDialog.SetRect(x, y, width, height)
		cnt.progressDialog.Draw(screen)
	}

	// top dialog
	if cnt.topDialog.IsDisplay() {
		cnt.topDialog.SetRect(x, y, width, height)
		cnt.topDialog.Draw(screen)

		return
	}

	// exec dialog
	if cnt.execDialog.IsDisplay() {
		cnt.execDialog.SetRect(x, y, width, height)
		cnt.execDialog.Draw(screen)

		return
	}

	// stats dialogs
	if cnt.statsDialog.IsDisplay() {
		cnt.statsDialog.SetRect(x, y, width, height)
		cnt.statsDialog.Draw(screen)

		return
	}

	// commit dialog
	if cnt.commitDialog.IsDisplay() {
		cnt.commitDialog.SetRect(x, y, width, height)
		cnt.commitDialog.Draw(screen)

		return
	}

	// checkpoint dialog
	if cnt.checkpointDialog.IsDisplay() {
		cnt.checkpointDialog.SetRect(x, y, width, height)
		cnt.checkpointDialog.Draw(screen)

		return
	}

	// restore dialog
	if cnt.restoreDialog.IsDisplay() {
		cnt.restoreDialog.SetRect(x, y, width, height)
		cnt.restoreDialog.Draw(screen)

		return
	}

	// terminal dialog
	if cnt.terminalDialog.IsDisplay() {
		cnt.terminalDialog.SetRect(cntViewX, cntViewY, cntViewW, cntViewH)
		cnt.terminalDialog.Draw(screen)

		return
	}
}
