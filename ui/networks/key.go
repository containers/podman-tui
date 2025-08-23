package networks

import (
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// InputHandler returns the handler for this primitive.
func (nets *Networks) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) { //nolint:cyclop
	return nets.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("view: networks event %v received", event)

		if nets.progressDialog.IsDisplay() {
			return
		}

		for _, dialog := range nets.getInnerDialogs() {
			if dialog.HasFocus() {
				if dialogHandler := dialog.InputHandler(); dialogHandler != nil {
					dialogHandler(event, setFocus)
				}
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

			// display sort menu
			if event.Rune() == utils.SortMenuKey.Rune() {
				nets.sortDialog.Display()
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
