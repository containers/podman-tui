package images

import (
	"github.com/gdamore/tcell/v2"
)

// Draw draws this primitive onto the screen.
func (img *Images) Draw(screen tcell.Screen) {
	img.DrawForSubclass(screen, img)
	img.SetBorder(false)

	imagewViewX, imagewViewY, imagewViewW, imagewViewH := img.GetInnerRect()

	img.refresh(imagewViewW)
	img.table.SetRect(imagewViewX, imagewViewY, imagewViewW, imagewViewH)
	img.table.SetBorder(true)

	img.table.Draw(screen)
	x, y, width, height := img.table.GetInnerRect()

	for _, dialog := range img.getInnerDialogs() {
		if dialog.IsDisplay() {
			dialog.SetRect(x, y, width, height)
			dialog.Draw(screen)

			break
		}
	}

	for _, dialog := range img.getInnerTopDialogs() {
		if dialog.IsDisplay() {
			dialog.SetRect(x, y, width, height)
			dialog.Draw(screen)

			break
		}
	}
}
