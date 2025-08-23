package secrets

import "github.com/gdamore/tcell/v2"

// Draw draws this primitive onto the screen.
func (s *Secrets) Draw(screen tcell.Screen) {
	s.DrawForSubclass(screen, s)
	s.SetBorder(false)

	secretViewX, secretViewY, secretViewW, secretViewH := s.GetInnerRect()

	s.table.SetRect(secretViewX, secretViewY, secretViewW, secretViewH)
	s.refresh(secretViewW)
	s.table.SetBorder(true)
	s.table.Draw(screen)

	x, y, width, height := s.table.GetInnerRect()
	for _, dialog := range s.getInnerDialogs() {
		if dialog.IsDisplay() {
			dialog.SetRect(x, y, width, height)
			dialog.Draw(screen)

			return
		}
	}
}
