package secrets

import (
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// InputHandler returns the handler for this primitive.
func (s *Secrets) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) { //nolint:gocognit,cyclop,lll
	return s.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("view: secrets event %v received", event)

		if s.progressDialog.IsDisplay() {
			return
		}

		// error dialog handler
		if s.errorDialog.HasFocus() {
			if errorDialogHandler := s.errorDialog.InputHandler(); errorDialogHandler != nil {
				errorDialogHandler(event, setFocus)
			}
		}

		// message dialog handler
		if s.messageDialog.HasFocus() {
			if messageDialogHandler := s.messageDialog.InputHandler(); messageDialogHandler != nil {
				messageDialogHandler(event, setFocus)
			}
		}

		// confirm dialog handler
		if s.confirmDialog.HasFocus() {
			if confirmDialogHandler := s.confirmDialog.InputHandler(); confirmDialogHandler != nil {
				confirmDialogHandler(event, setFocus)
			}
		}

		// command dialog handler
		if s.cmdDialog.HasFocus() {
			if cmdHandler := s.cmdDialog.InputHandler(); cmdHandler != nil {
				cmdHandler(event, setFocus)
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
