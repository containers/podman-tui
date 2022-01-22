package dialogs

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// ErrorDialog is an error dialog primitive
type ErrorDialog struct {
	*tview.Box
	modal   *tview.Modal
	display bool
}

// NewErrorDialog returns new error dialog primitive
func NewErrorDialog() *ErrorDialog {
	dialog := ErrorDialog{
		Box:     tview.NewBox(),
		modal:   tview.NewModal().SetBackgroundColor(tcell.ColorRed).AddButtons([]string{"OK"}),
		display: false,
	}
	dialog.modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		dialog.Hide()
	})

	return &dialog
}

// Display displays this primitive
func (d *ErrorDialog) Display() {
	d.display = true
}

// IsDisplay returns true if primitive is shown
func (d *ErrorDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive
func (d *ErrorDialog) Hide() {
	d.SetText("")
	d.display = false
}

// SetText sets error dialog message
func (d *ErrorDialog) SetText(message string) {
	d.modal.SetText(message)
}

// HasFocus returns whether or not this primitive has focus
func (d *ErrorDialog) HasFocus() bool {
	return d.modal.HasFocus()
}

// Focus is called when this primitive receives focus
func (d *ErrorDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(d.modal)
}

// InputHandler returns input handler function for this primitive
func (d *ErrorDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("error dialog: event %v received", event.Key())
		if modalHandler := d.modal.InputHandler(); modalHandler != nil {
			modalHandler(event, setFocus)
			return
		}
	})
}

// SetRect set rects for this primitive.
func (d *ErrorDialog) SetRect(x, y, width, height int) {

	d.Box.SetRect(x, y, width, height)
}

// GetRect returns the current position of the primitive, x, y, width, and
// height.
func (d *ErrorDialog) GetRect() (int, int, int, int) {
	return d.Box.GetRect()
}

// Draw draws this primitive onto the screen.
func (d *ErrorDialog) Draw(screen tcell.Screen) {
	d.modal.Draw(screen)
}

// SetDoneFunc sets modal done function
func (d *ErrorDialog) SetDoneFunc(handler func()) *ErrorDialog {
	d.modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		handler()
	})
	return d
}
