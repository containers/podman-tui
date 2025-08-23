package secrets

import (
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// InputHandler returns the handler for this primitive.
func (s *Secrets) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) { //nolint:cyclop
	return s.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("view: secrets event %v received", event)

		if s.progressDialog.IsDisplay() {
			return
		}

		for _, dialog := range s.getInnerDialogs() {
			if dialog.HasFocus() {
				if dialogHandler := dialog.InputHandler(); dialogHandler != nil {
					dialogHandler(event, setFocus)
				}
			}
		}

		// table handlers
		if s.table.HasFocus() { //nolint:nestif
			if event.Rune() == utils.CommandMenuKey.Rune() {
				if s.cmdDialog.GetCommandCount() <= 1 {
					return
				}

				s.cmdDialog.Display()
				setFocus(s)

				return
			}

			// display sort menu
			if event.Rune() == utils.SortMenuKey.Rune() {
				s.sortDialog.Display()
				setFocus(s)

				return
			}

			if event.Key() == utils.DeleteKey.EventKey() {
				s.rm()
				setFocus(s)

				return
			}

			if tableHandler := s.table.InputHandler(); tableHandler != nil {
				tableHandler(event, setFocus)
			}
		}

		setFocus(s)
	})
}
