package images

import (
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// InputHandler returns the handler for this primitive.
func (img *Images) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) { //nolint:gocognit,cyclop,lll
	return img.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("view: images event %v received", event)

		if img.progressDialog.IsDisplay() {
			return
		}

		for _, dialog := range img.getInnerTopDialogs() {
			if dialog.HasFocus() {
				if handler := dialog.InputHandler(); handler != nil {
					handler(event, setFocus)
				}
			}
		}

		for _, dialog := range img.getInnerDialogs() {
			if dialog.HasFocus() {
				if handler := dialog.InputHandler(); handler != nil {
					handler(event, setFocus)
				}
			}
		}

		// table handlers
		if img.table.HasFocus() { //nolint:nestif
			img.selectedID, img.selectedName = img.getSelectedItem()
			if event.Rune() == utils.CommandMenuKey.Rune() {
				if img.cmdDialog.GetCommandCount() <= 1 {
					return
				}

				img.cmdDialog.Display()
				setFocus(img)

				return
			}

			// display sort menu
			if event.Rune() == utils.SortMenuKey.Rune() {
				img.sortDialog.Display()
				setFocus(img)

				return
			}

			if event.Key() == utils.DeleteKey.EventKey() {
				img.rm()
				setFocus(img)

				return
			}

			if tableHandler := img.table.InputHandler(); tableHandler != nil {
				tableHandler(event, setFocus)
			}
		}

		setFocus(img)
	})
}
