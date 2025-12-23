package voldialogs

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

func volumeDialogInputHandler(
	name string,
	d *tview.Box,
	form *tview.Form,
	innerPrimitives []tview.Primitive,
	nextFocusHandler func(),
	cancelHandler func(),
	doneHandler func(),
) func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("volume %s dialog: event %v received", name, event)

		if event.Key() == tcell.KeyEsc {
			cancelHandler()

			return
		}

		if form.HasFocus() { //nolint:nestif
			if formHandler := form.InputHandler(); formHandler != nil {
				if event.Key() == tcell.KeyEnter {
					enterButton := form.GetButton(form.GetButtonCount() - 1)
					if enterButton.HasFocus() {
						doneHandler()
					}
				}

				formHandler(event, setFocus)

				return
			}
		}

		if event.Key() == tcell.KeyTab {
			nextFocusHandler()
			d.Focus(setFocus)

			return
		}

		for _, item := range innerPrimitives {
			if item.HasFocus() {
				if handler := item.InputHandler(); handler != nil {
					handler(event, setFocus)

					return
				}
			}
		}
	})
}
