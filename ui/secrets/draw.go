package secrets

import "github.com/gdamore/tcell/v2"

// Draw draws this primitive onto the screen.
func (s *Secrets) Draw(screen tcell.Screen) {
	s.DrawForSubclass(screen, s)
	s.SetBorder(false)

	x, y, w, h := s.GetInnerRect()

	s.table.SetRect(x, y, w, h)
	s.refresh(w)
	s.table.SetBorder(true)
	s.table.Draw(screen)

	for _, dialog := range s.getInnerDialogs() {
		if dialog.IsDisplay() {
			dialog.SetRect(x, y, w, h)
			dialog.Draw(screen)

			return
		}
	}
}
