package images

import (
	"github.com/gdamore/tcell/v2"
)

// Draw draws this primitive onto the screen.
func (img *Images) Draw(screen tcell.Screen) { //nolint:cyclop
	img.refresh()
	img.Box.DrawForSubclass(screen, img)
	img.Box.SetBorder(false)

	x, y, width, height := img.GetInnerRect()

	img.table.SetRect(x, y, width, height)
	img.table.SetBorder(true)

	img.table.Draw(screen)
	x, y, width, height = img.table.GetInnerRect()

	// error dialog
	if img.errorDialog.IsDisplay() {
		img.errorDialog.SetRect(x, y, width, height)
		img.errorDialog.Draw(screen)

		return
	}

	// command dialog
	if img.cmdDialog.IsDisplay() {
		img.cmdDialog.SetRect(x, y, width, height)
		img.cmdDialog.Draw(screen)

		return
	}

	// command input dialog
	if img.cmdInputDialog.IsDisplay() {
		img.cmdInputDialog.SetRect(x, y, width, height)
		img.cmdInputDialog.Draw(screen)

		return
	}

	// message dialog
	if img.messageDialog.IsDisplay() {
		img.messageDialog.SetRect(x, y, width, height+1)
		img.messageDialog.Draw(screen)

		return
	}

	// confirm dialog
	if img.confirmDialog.IsDisplay() {
		img.confirmDialog.SetRect(x, y, width, height)
		img.confirmDialog.Draw(screen)

		return
	}

	// search dialog
	if img.searchDialog.IsDisplay() {
		img.searchDialog.SetRect(x, y, width, height)
		img.searchDialog.Draw(screen)
	}

	// progress dialog
	if img.progressDialog.IsDisplay() {
		img.progressDialog.SetRect(x, y, width, height)
		img.progressDialog.Draw(screen)
	}

	// history dialog
	if img.historyDialog.IsDisplay() {
		img.historyDialog.SetRect(x, y, width, height)
		img.historyDialog.Draw(screen)

		return
	}

	// build dialog
	if img.buildDialog.IsDisplay() {
		img.buildDialog.SetRect(x, y, width, height)
		img.buildDialog.Draw(screen)

		return
	}

	// build progress dialog
	if img.buildPrgDialog.IsDisplay() {
		img.buildPrgDialog.SetRect(x, y, width, height)
		img.buildPrgDialog.Draw(screen)

		return
	}

	// save dialog
	if img.saveDialog.IsDisplay() {
		img.saveDialog.SetRect(x, y, width, height)
		img.saveDialog.Draw(screen)

		return
	}

	// import dialog
	if img.importDialog.IsDisplay() {
		img.importDialog.SetRect(x, y, width, height)
		img.importDialog.Draw(screen)

		return
	}

	// push dialog
	if img.pushDialog.IsDisplay() {
		img.pushDialog.SetRect(x, y, width, height)
		img.pushDialog.Draw(screen)

		return
	}
}
