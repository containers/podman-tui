package volumes

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// InputHandler returns the handler for this primitive.
func (vols *Volumes) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return vols.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("view: volumes event %v received", event.Key())
		if vols.progressDialog.IsDisplay() {
			return
		}
		// error dialog handler
		if vols.errorDialog.HasFocus() {
			if errorDialogHandler := vols.errorDialog.InputHandler(); errorDialogHandler != nil {
				errorDialogHandler(event, setFocus)
			}
		}
		// message dialog handler
		if vols.messageDialog.HasFocus() {
			if messageDialogHandler := vols.messageDialog.InputHandler(); messageDialogHandler != nil {
				messageDialogHandler(event, setFocus)
			}
		}
		// create dialog dialog handler
		if vols.createDialog.HasFocus() {
			if createDialogHandler := vols.createDialog.InputHandler(); createDialogHandler != nil {
				createDialogHandler(event, setFocus)
			}
		}
		// confirm dialog handler
		if vols.confirmDialog.HasFocus() {
			if confirmDialogHandler := vols.confirmDialog.InputHandler(); confirmDialogHandler != nil {
				confirmDialogHandler(event, setFocus)
			}
		}
		// table handlers
		if vols.table.HasFocus() {
			if event.Key() == tcell.KeyCtrlV || event.Key() == tcell.KeyEnter {
				if vols.cmdDialog.GetCommandCount() <= 1 {
					return
				}
				vols.selectedID = vols.getSelectedItem()
				vols.cmdDialog.Display()
			}
			if tableHandler := vols.table.InputHandler(); tableHandler != nil {
				tableHandler(event, setFocus)
			}
		}

		// command dialog handler
		if vols.cmdDialog.HasFocus() {
			if cmdHandler := vols.cmdDialog.InputHandler(); cmdHandler != nil {
				cmdHandler(event, setFocus)
			}
		}

		setFocus(vols)
	})
}
