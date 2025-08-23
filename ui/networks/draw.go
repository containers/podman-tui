package networks

import (
	"github.com/gdamore/tcell/v2"
)

// Draw draws this primitive onto the screen.
func (nets *Networks) Draw(screen tcell.Screen) {
	nets.DrawForSubclass(screen, nets)
	nets.SetBorder(false)

	netViewX, netViewY, netViewW, netViewH := nets.GetInnerRect()

	nets.table.SetRect(netViewX, netViewY, netViewW, netViewH)
	nets.refresh(netViewW)
	nets.table.SetBorder(true)

	nets.table.Draw(screen)

	x, y, width, height := nets.table.GetInnerRect()

	for _, dialog := range nets.getInnerDialogs() {
		if dialog.IsDisplay() {
			dialog.SetRect(x, y, width, height)
			dialog.Draw(screen)

			return
		}
	}
}
