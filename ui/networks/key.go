package networks

import (
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// InputHandler returns the handler for this primitive.
func (nets *Networks) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) { //nolint:gocognit,lll,cyclop
	return nets.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("view: networks event %v received", event)
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

		// create dialog handler
		if nets.createDialog.HasFocus() {
			if createDialogHandler := nets.createDialog.InputHandler(); createDialogHandler != nil {
				createDialogHandler(event, setFocus)
			}
		}

		// connect dialog handler
		if nets.connectDialog.HasFocus() {
			if connectDialogHandler := nets.connectDialog.InputHandler(); connectDialogHandler != nil {
				connectDialogHandler(event, setFocus)
			}
		}

		// disconnect dialog handler
		if nets.disconnectDialog.HasFocus() {
			if disconnectDialogHandler := nets.disconnectDialog.InputHandler(); disconnectDialogHandler != nil {
				disconnectDialogHandler(event, setFocus)
			}
		}

		// confirm dialog handler
		if nets.confirmDialog.HasFocus() {
			if confirmDialogHandler := nets.confirmDialog.InputHandler(); confirmDialogHandler != nil {
				confirmDialogHandler(event, setFocus)
			}
		}

		// command dialog handler
		if nets.cmdDialog.HasFocus() {
			if cmdHandler := nets.cmdDialog.InputHandler(); cmdHandler != nil {
				cmdHandler(event, setFocus)
			}
		}

		// table handlers
		if nets.table.HasFocus() { //nolint:nestif
			nets.selectedID, _ = nets.getSelectedItem()
			if event.Rune() == utils.CommandMenuKey.Rune() {
				if nets.cmdDialog.GetCommandCount() <= 1 {
					return
				}

				nets.cmdDialog.Display()
				setFocus(nets)

				return
			}
			if event.Key() == utils.DeleteKey.EventKey() {
				nets.rm()
				setFocus(nets)

				return
			}

			if tableHandler := nets.table.InputHandler(); tableHandler != nil {
				tableHandler(event, setFocus)
			}
		}

		setFocus(nets)
	})
}
