package networks

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// InputHandler returns the handler for this primitive.
func (nets *Networks) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return nets.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("view: networks event %v received", event.Key())
		if nets.progressDialog.IsDisplay() {
			return
		}
		// error dialog handler
		if nets.errorDialog.HasFocus() {
			if errorDialogHandler := nets.errorDialog.InputHandler(); errorDialogHandler != nil {
				errorDialogHandler(event, setFocus)
			}
		}
		// message dialog handler
		if nets.messageDialog.HasFocus() {
			if messageDialogHandler := nets.messageDialog.InputHandler(); messageDialogHandler != nil {
				messageDialogHandler(event, setFocus)
			}
		}
		// create dialog dialog handler
		if nets.createDialog.HasFocus() {
			if createDialogHandler := nets.createDialog.InputHandler(); createDialogHandler != nil {
				createDialogHandler(event, setFocus)
			}
		}
		// confirm dialog handler
		if nets.confirmDialog.HasFocus() {
			if confirmDialogHandler := nets.confirmDialog.InputHandler(); confirmDialogHandler != nil {
				confirmDialogHandler(event, setFocus)
			}
		}
		// table handlers
		if nets.table.HasFocus() {
			if event.Key() == tcell.KeyCtrlV || event.Key() == tcell.KeyEnter {
				if nets.cmdDialog.GetCommandCount() <= 1 {
					return
				}
				nets.selectedID = nets.getSelectedItem()
				nets.cmdDialog.Display()
			}
			if tableHandler := nets.table.InputHandler(); tableHandler != nil {
				tableHandler(event, setFocus)
			}
		}

		// command dialog handler
		if nets.cmdDialog.HasFocus() {
			if cmdHandler := nets.cmdDialog.InputHandler(); cmdHandler != nil {
				cmdHandler(event, setFocus)
			}
		}

		setFocus(nets)
	})
}
