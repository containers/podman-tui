package pods

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// InputHandler returns the handler for this primitive.
func (pods *Pods) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return pods.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("view: pods event %v received", event.Key())
		if pods.progressDialog.IsDisplay() {
			return
		}
		// error dialog handler
		if pods.errorDialog.HasFocus() {
			if errorDialogHandler := pods.errorDialog.InputHandler(); errorDialogHandler != nil {
				errorDialogHandler(event, setFocus)
			}
		}
		// message dialog handler
		if pods.messageDialog.HasFocus() {
			if messageDialogHandler := pods.messageDialog.InputHandler(); messageDialogHandler != nil {
				messageDialogHandler(event, setFocus)
			}
		}
		// create dialog dialog handler
		if pods.createDialog.HasFocus() {
			if createDialogHandler := pods.createDialog.InputHandler(); createDialogHandler != nil {
				createDialogHandler(event, setFocus)
			}
		}
		// confirm dialog handler
		if pods.confirmDialog.HasFocus() {
			if confirmDialogHandler := pods.confirmDialog.InputHandler(); confirmDialogHandler != nil {
				confirmDialogHandler(event, setFocus)
			}
		}
		// table handlers
		if pods.table.HasFocus() {
			if event.Key() == tcell.KeyCtrlV || event.Key() == tcell.KeyEnter {
				if pods.cmdDialog.GetCommandCount() <= 1 {
					return
				}
				pods.selectedID = pods.getSelectedItem()
				pods.cmdDialog.Display()
			}
			if tableHandler := pods.table.InputHandler(); tableHandler != nil {
				tableHandler(event, setFocus)
			}
		}

		// command dialog handler
		if pods.cmdDialog.HasFocus() {
			if cmdHandler := pods.cmdDialog.InputHandler(); cmdHandler != nil {
				cmdHandler(event, setFocus)
			}
		}

		// top dialog handler
		if pods.topDialog.HasFocus() {
			if topDialogHandler := pods.topDialog.InputHandler(); topDialogHandler != nil {
				topDialogHandler(event, setFocus)
			}
		}
		// container stats dialog handler
		if pods.statsDialog.HasFocus() {
			if podStatsDialogHandler := pods.statsDialog.InputHandler(); podStatsDialogHandler != nil {
				podStatsDialogHandler(event, setFocus)
			}
		}
		setFocus(pods)
	})
}
