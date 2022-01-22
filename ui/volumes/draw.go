package volumes

import (
	"github.com/gdamore/tcell/v2"
)

// Draw draws this primitive onto the screen.
func (vols *Volumes) Draw(screen tcell.Screen) {
	vols.refresh()
	vols.Box.DrawForSubclass(screen, vols)
	vols.Box.SetBorder(false)
	x, y, width, height := vols.GetInnerRect()
	vols.table.SetRect(x, y, width, height)
	vols.table.SetBorder(true)

	vols.table.Draw(screen)
	x, y, width, height = vols.table.GetInnerRect()
	// error dialog
	if vols.errorDialog.IsDisplay() {
		vols.errorDialog.SetRect(x, y, width, height)
		vols.errorDialog.Draw(screen)
		return
	}
	// message dialog
	if vols.messageDialog.IsDisplay() {
		vols.messageDialog.SetRect(x, y, width, height+1)
		vols.messageDialog.Draw(screen)
		return
	}
	// confirm dialog
	if vols.confirmDialog.IsDisplay() {
		vols.confirmDialog.SetRect(x, y, width, height)
		vols.confirmDialog.Draw(screen)
		return
	}
	// command dialog dialog
	if vols.cmdDialog.IsDisplay() {
		vols.cmdDialog.SetRect(x, y, width, height)
		vols.cmdDialog.Draw(screen)
		return
	}
	// create dialog dialog
	if vols.createDialog.IsDisplay() {
		vols.createDialog.SetRect(x, y, width, height)
		vols.createDialog.Draw(screen)
		return
	}
	// progress dialog
	if vols.progressDialog.IsDisplay() {
		vols.progressDialog.SetRect(x, y, width, height)
		vols.progressDialog.Draw(screen)
	}

}
