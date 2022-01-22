package connection

import "github.com/gdamore/tcell/v2"

// Draw draws this primitive onto the screen.
func (conn *Connection) Draw(screen tcell.Screen) {
	conn.Box.DrawForSubclass(screen, conn)
	x, y, width, height := conn.Box.GetInnerRect()
	emptyWidth := (width - maxWidth) / 2
	emptyHeight := (height - maxHeight) / 2
	if width > maxWidth {
		width = maxWidth
		x = x + emptyWidth
	}
	if height > maxHeight {
		height = maxHeight
		y = y + emptyHeight
	}
	if conn.textview.GetText(true) == "" {
		height = height - 5
	}
	conn.layout.SetRect(x, y, width, height)
	conn.progressDialog.Pulse()
	conn.layout.Draw(screen)
}
