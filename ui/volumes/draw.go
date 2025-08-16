package volumes

import (
	"github.com/gdamore/tcell/v2"
)

// Draw draws this primitive onto the screen.
func (vols *Volumes) Draw(screen tcell.Screen) {
	vols.DrawForSubclass(screen, vols)
	vols.SetBorder(false)

	x, y, width, height := vols.GetInnerRect()

	vols.refresh(width)
	vols.table.SetRect(x, y, width, height)
	vols.table.SetBorder(true)

	vols.table.Draw(screen)

	x, y, width, height = vols.table.GetInnerRect()

	for _, dialog := range vols.getInnerDialogs() {
		if dialog.IsDisplay() {
			dialog.SetRect(x, y, width, height)
			dialog.Draw(screen)

			break
		}
	}
}
