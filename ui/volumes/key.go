package volumes

import (
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// InputHandler returns the handler for this primitive.
func (vols *Volumes) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return vols.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("view: volumes event %v received", event)

		if vols.progressDialog.IsDisplay() {
			return
		}

		for _, dialog := range vols.getInnerDialogs() {
			if dialog.HasFocus() {
				if handler := dialog.InputHandler(); handler != nil {
					handler(event, setFocus)
				}
			}
		}

		// table handlers
		if vols.table.HasFocus() {
			vols.processTableInputHandler(event, setFocus)
		}

		setFocus(vols)
	})
}

func (vols *Volumes) processTableInputHandler(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	if event.Rune() == utils.CommandMenuKey.Rune() {
		if vols.cmdDialog.GetCommandCount() <= 1 {
			return
		}

		vols.cmdDialog.Display()

		return
	}

	// display sort menu
	if event.Rune() == utils.SortMenuKey.Rune() {
		vols.sortDialog.Display()
		setFocus(vols)

		return
	}

	if event.Key() == utils.DeleteKey.EventKey() {
		vols.removePrep()

		return
	}

	if tableHandler := vols.table.InputHandler(); tableHandler != nil {
		tableHandler(event, setFocus)
	}
}
