package images

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// InputHandler returns the handler for this primitive.
func (img *Images) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return img.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("view: images event %v received", event.Key())
		if img.progressDialog.IsDisplay() {
			return
		}
		// error dialog handler
		if img.errorDialog.HasFocus() || img.errorDialog.IsDisplay() {
			if errorDialogHandler := img.errorDialog.InputHandler(); errorDialogHandler != nil {
				errorDialogHandler(event, setFocus)
			}
		}
		// message dialog handler
		if img.messageDialog.HasFocus() || img.messageDialog.IsDisplay() {
			if messageDialogHandler := img.messageDialog.InputHandler(); messageDialogHandler != nil {
				messageDialogHandler(event, setFocus)
			}
		}

		// table handlers
		if img.table.HasFocus() {
			if event.Key() == tcell.KeyCtrlV || event.Key() == tcell.KeyEnter {
				if img.cmdDialog.GetCommandCount() <= 1 {
					return
				}
				img.selectedID, img.selectedName = img.getSelectedItem()
				img.cmdDialog.Display()
			}
			if tableHandler := img.table.InputHandler(); tableHandler != nil {
				tableHandler(event, setFocus)
			}
		}

		// command dialog handler
		if img.cmdDialog.HasFocus() {
			if cmdHandler := img.cmdDialog.InputHandler(); cmdHandler != nil {
				cmdHandler(event, setFocus)
			}
		}
		// input dialog handler
		if img.cmdInputDialog.HasFocus() {
			if cmdInputHandler := img.cmdInputDialog.InputHandler(); cmdInputHandler != nil {
				cmdInputHandler(event, setFocus)
			}
		}

		// confirm dialog handler
		if img.confirmDialog.HasFocus() {
			if confirmDialogHandler := img.confirmDialog.InputHandler(); confirmDialogHandler != nil {
				confirmDialogHandler(event, setFocus)
			}
		}

		// search dialog handler
		if img.searchDialog.HasFocus() {
			if searchDialogHandler := img.searchDialog.InputHandler(); searchDialogHandler != nil {
				searchDialogHandler(event, setFocus)
			}
		}

		// history dialog handler
		if img.historyDialog.HasFocus() {
			if historyDialogHandler := img.historyDialog.InputHandler(); historyDialogHandler != nil {
				historyDialogHandler(event, setFocus)
			}
		}

		setFocus(img)
	})
}
