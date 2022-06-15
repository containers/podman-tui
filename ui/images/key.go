package images

import (
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// InputHandler returns the handler for this primitive.
func (img *Images) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return img.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("view: images event %v received", event)
		if img.progressDialog.IsDisplay() {
			return
		}
		// error dialog handler
		if img.errorDialog.HasFocus() || img.errorDialog.IsDisplay() {
			if errorDialogHandler := img.errorDialog.InputHandler(); errorDialogHandler != nil {
				errorDialogHandler(event, setFocus)
				setFocus(img.errorDialog)
			}
		}
		// message dialog handler
		if img.messageDialog.HasFocus() || img.messageDialog.IsDisplay() {
			if messageDialogHandler := img.messageDialog.InputHandler(); messageDialogHandler != nil {
				messageDialogHandler(event, setFocus)
				setFocus(img.messageDialog)
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

		// build dialog handler
		if img.buildDialog.HasFocus() {
			if buildDialogHandler := img.buildDialog.InputHandler(); buildDialogHandler != nil {
				buildDialogHandler(event, setFocus)
			}
		}

		// build progress dialog handler
		if img.buildPrgDialog.HasFocus() {
			if buildPrgDialogHandler := img.buildPrgDialog.InputHandler(); buildPrgDialogHandler != nil {
				buildPrgDialogHandler(event, setFocus)
			}
		}

		// save dialog handler
		if img.saveDialog.HasFocus() {
			if saveDialogHandler := img.saveDialog.InputHandler(); saveDialogHandler != nil {
				saveDialogHandler(event, setFocus)
			}
		}

		// import dialog handler
		if img.importDialog.HasFocus() {
			if importDialogHandler := img.importDialog.InputHandler(); importDialogHandler != nil {
				importDialogHandler(event, setFocus)
			}
		}

		// push dialog handler
		if img.pushDialog.HasFocus() {
			if pushDialogHandler := img.pushDialog.InputHandler(); pushDialogHandler != nil {
				pushDialogHandler(event, setFocus)
			}
		}

		// table handlers
		if img.table.HasFocus() {
			img.selectedID, img.selectedName = img.getSelectedItem()
			if event.Rune() == utils.CommandMenuKey.Rune() {
				if img.cmdDialog.GetCommandCount() <= 1 {
					return
				}
				img.cmdDialog.Display()
			} else if event.Key() == utils.DeleteKey.EventKey() {
				img.rm()
			} else {
				if tableHandler := img.table.InputHandler(); tableHandler != nil {
					tableHandler(event, setFocus)
				}
			}
		}

		setFocus(img)
	})
}
