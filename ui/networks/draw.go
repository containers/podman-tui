package networks

import (
	"github.com/gdamore/tcell/v2"
)

// Draw draws this primitive onto the screen.
func (nets *Networks) Draw(screen tcell.Screen) {
	nets.DrawForSubclass(screen, nets)
	nets.SetBorder(false)

	x, y, w, h := nets.GetInnerRect()

	nets.table.SetRect(x, y, w, h)
	nets.refresh(w)
	nets.table.SetBorder(true)

	nets.table.Draw(screen)

	for _, dialog := range nets.getInnerDialogs() {
		if dialog.IsDisplay() {
			dialog.SetRect(x, y, w, h)
			dialog.Draw(screen)

			return
		}
	}
}
