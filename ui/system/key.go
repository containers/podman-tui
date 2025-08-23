package system

import (
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// InputHandler returns the handler for this primitive.
func (sys *System) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) { //nolint:cyclop
	return sys.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("view: system event %v received", event)

		if sys.progressDialog.IsDisplay() {
			return
		}

		for _, dialog := range sys.getInnerDialogs(true) {
			if dialog.HasFocus() {
				if handler := dialog.InputHandler(); handler != nil {
					handler(event, setFocus)
				}
			}
		}

		// table handlers
		if sys.connTable.HasFocus() { //nolint:nestif
			if event.Rune() == utils.CommandMenuKey.Rune() {
				if sys.cmdDialog.GetCommandCount() <= 1 {
					return
				}

				sys.cmdDialog.Display()
				setFocus(sys)

				return
			}

			// display sort menu
			if event.Rune() == utils.SortMenuKey.Rune() {
				sys.sortDialog.Display()
				setFocus(sys)

				return
			}

			if event.Key() == utils.DeleteKey.EventKey() {
				sys.cremove()
				setFocus(sys)

				return
			}

			if tableHandler := sys.connTable.InputHandler(); tableHandler != nil {
				tableHandler(event, setFocus)
			}
		}

		setFocus(sys)
	})
}
