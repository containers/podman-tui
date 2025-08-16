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

	// error dialog
	if s.errorDialog.IsDisplay() {
		s.errorDialog.SetRect(x, y, width, height)
		s.errorDialog.Draw(screen)

		return
	}

	// progress dialog
	if s.progressDialog.IsDisplay() {
		s.progressDialog.SetRect(x, y, width, height)
		s.progressDialog.Draw(screen)
	}

	// message dialog
	if s.messageDialog.IsDisplay() {
		if s.messageDialog.IsDisplayFullSize() {
			s.messageDialog.SetRect(secretViewX, secretViewY, secretViewW, secretViewH)
		} else {
			s.messageDialog.SetRect(x, y, width, height+1)
		}

		s.messageDialog.Draw(screen)

		return
	}

	// confirm dialog
	if s.confirmDialog.IsDisplay() {
		s.confirmDialog.SetRect(x, y, width, height)
		s.confirmDialog.Draw(screen)

		return
	}

	// command dialog
	if s.cmdDialog.IsDisplay() {
		s.cmdDialog.SetRect(x, y, width, height)
		s.cmdDialog.Draw(screen)

		return
	}

	// create dialog
	if s.createDialog.IsDisplay() {
		s.createDialog.SetRect(x, y, width, height)
		s.createDialog.Draw(screen)
	}
}
