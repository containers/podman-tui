package images

import (
	"github.com/gdamore/tcell/v2"
)

// Draw draws this primitive onto the screen.
func (img *Images) Draw(screen tcell.Screen) {
	img.DrawForSubclass(screen, img)
	img.SetBorder(false)

	x, y, w, h := img.GetInnerRect()

	img.refresh(w)
	img.table.SetRect(x, y, w, h)
	img.table.SetBorder(true)

	img.table.Draw(screen)

	for _, dialog := range img.getInnerDialogs() {
		if dialog.IsDisplay() {
			dialog.SetRect(x, y, w, h)
			dialog.Draw(screen)

			break
		}
	}

	for _, dialog := range img.getInnerTopDialogs() {
		if dialog.IsDisplay() {
			dialog.SetRect(x, y, w, h)
			dialog.Draw(screen)

			break
		}
	}
}
