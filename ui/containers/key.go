package containers

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// InputHandler returns the handler for this primitive.
func (cnt *Containers) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return cnt.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("view: containers event %v received", event.Key())
		if cnt.progressDialog.IsDisplay() {
			return
		}
		// command dialog handler
		if cnt.cmdDialog.HasFocus() {
			if cmdHandler := cnt.cmdDialog.InputHandler(); cmdHandler != nil {
				cmdHandler(event, setFocus)
			}
		}
		// input dialog handler
		if cnt.cmdInputDialog.HasFocus() {
			if cmdInputHandler := cnt.cmdInputDialog.InputHandler(); cmdInputHandler != nil {
				cmdInputHandler(event, setFocus)
			}
		}

		// message dialog handler
		if cnt.messageDialog.HasFocus() {
			if messageDialogHandler := cnt.messageDialog.InputHandler(); messageDialogHandler != nil {
				messageDialogHandler(event, setFocus)
			}
		}
		// create dialog dialog handler
		if cnt.createDialog.HasFocus() {
			if createDialogHandler := cnt.createDialog.InputHandler(); createDialogHandler != nil {
				createDialogHandler(event, setFocus)
			}
		}
		// exec dialog dialog handler
		if cnt.execDialog.HasFocus() {
			if execDialogHandler := cnt.execDialog.InputHandler(); execDialogHandler != nil {
				execDialogHandler(event, setFocus)
			}
		}
		// exec terminal dialog dialog handler
		if cnt.execTerminalDialog.HasFocus() {
			if execTerminalDialogHandler := cnt.execTerminalDialog.InputHandler(); execTerminalDialogHandler != nil {
				execTerminalDialogHandler(event, setFocus)
			}
		}
		// confirm dialog handler
		if cnt.confirmDialog.HasFocus() {
			if confirmDialogHandler := cnt.confirmDialog.InputHandler(); confirmDialogHandler != nil {
				confirmDialogHandler(event, setFocus)
			}
		}
		// table handlers
		if cnt.table.HasFocus() {
			if event.Key() == tcell.KeyCtrlV || event.Key() == tcell.KeyEnter {
				if cnt.cmdDialog.GetCommandCount() <= 1 {
					return
				}
				cnt.selectedID, cnt.selectedName = cnt.getSelectedItem()
				cnt.cmdDialog.Display()
			}
			if tableHandler := cnt.table.InputHandler(); tableHandler != nil {
				tableHandler(event, setFocus)
			}
		}
		// error dialog handler
		if cnt.errorDialog.HasFocus() {
			if errorDialogHandler := cnt.errorDialog.InputHandler(); errorDialogHandler != nil {
				errorDialogHandler(event, setFocus)
			}
		}

		// container top dialog handler
		if cnt.topDialog.HasFocus() {
			if cntTopDialogHandler := cnt.topDialog.InputHandler(); cntTopDialogHandler != nil {
				cntTopDialogHandler(event, setFocus)
			}
		}
		setFocus(cnt)
	})
}
