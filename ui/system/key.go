package system

import (
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// InputHandler returns the handler for this primitive.
func (sys *System) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return sys.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("view: system event %v received", event)
		if sys.progressDialog.IsDisplay() {
			return
		}
		// command dialog handler
		if sys.cmdDialog.HasFocus() {
			if cmdHandler := sys.cmdDialog.InputHandler(); cmdHandler != nil {
				cmdHandler(event, setFocus)
			}
		}
		// confirm dialog handler
		if sys.confirmDialog.HasFocus() {
			if confirmDialogHandler := sys.confirmDialog.InputHandler(); confirmDialogHandler != nil {
				confirmDialogHandler(event, setFocus)
			}
		}
		// message dialog handler
		if sys.messageDialog.HasFocus() {
			if messageDialogHandler := sys.messageDialog.InputHandler(); messageDialogHandler != nil {
				messageDialogHandler(event, setFocus)
			}
		}
		// disk usage dialog
		if sys.dfDialog.HasFocus() {
			if dfDialogHandler := sys.dfDialog.InputHandler(); dfDialogHandler != nil {
				dfDialogHandler(event, setFocus)
			}
		}
		// error dialog handler
		if sys.errorDialog.HasFocus() {
			if errorDialogHandler := sys.errorDialog.InputHandler(); errorDialogHandler != nil {
				errorDialogHandler(event, setFocus)
			}
		}
		// textview handlers
		if sys.textview.HasFocus() {
			// workaround to give disk usage dialog focus after first event when dialog is displayed
			// but system page has focus
			if sys.dfDialog.IsDisplay() {
				// disk usage dialog
				if dfDialogHandler := sys.dfDialog.InputHandler(); dfDialogHandler != nil {
					dfDialogHandler(event, setFocus)
				}
			} else {
				if event.Rune() == utils.CommandMenuKey.Rune() {
					if sys.cmdDialog.GetCommandCount() <= 1 {
						return
					}
					sys.cmdDialog.Display()
				} else {
					if textviewHandler := sys.textview.InputHandler(); textviewHandler != nil {
						textviewHandler(event, setFocus)
					}
				}
			}

		}
		setFocus(sys)

	})
}
