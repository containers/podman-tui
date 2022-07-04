package tvxwidgets

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// represents dialog type.
const (
	InfoDialog = 0 + iota
	ErrorDailog
)

// MessageDialog represents message dialog primitive.
type MessageDialog struct {
	*tview.Box
	// layout message dialog layout
	layout *tview.Flex
	// message view
	textview *tview.TextView
	// dialog form buttons
	form *tview.Form
	// message dialog X
	x int
	// message dialog Y
	y int
	// message dialog width
	width int
	// message dialog heights
	height int
	// dialog type info and error
	// type will change the default background color for the dialog
	messageType int
	// background color
	bgColor tcell.Color
	// message dialog text message to display.
	message string
	// callback for whwen user clicked on the the button or presses "enter" or "esc"
	doneHandler func()
}

// NewMessageDialog returns a new message dialog primitive.
func NewMessageDialog(dtype int) *MessageDialog {
	dialog := &MessageDialog{
		Box:         tview.NewBox(),
		messageType: dtype,
		bgColor:     tcell.ColorSteelBlue,
	}

	dialog.textview = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)

	dialog.form = tview.NewForm().
		AddButton("Enter", nil).
		SetButtonsAlign(tview.AlignRight)

	dialog.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	dialog.layout.AddItem(dialog.textview, 0, 0, true)
	dialog.layout.AddItem(dialog.form, dialogFormHeight, 0, true)
	dialog.layout.SetBorder(true)

	dialog.setColor()

	return dialog
}

// SetBorder sets dialogs border - no effect always true.
func (d *MessageDialog) SetBorder(status bool) {}

// SetType sets dialog type to info or error.
func (d *MessageDialog) SetType(dtype int) {
	if dtype >= 0 && dtype <= 2 {
		d.messageType = dtype
		d.setColor()
	}
}

// SetTitle sets title for this primitive.
func (d *MessageDialog) SetTitle(title string) {
	d.layout.SetTitle(title)
}

// SetBackgroundColor sets dialog background color.
func (d *MessageDialog) SetBackgroundColor(color tcell.Color) {
	d.bgColor = color
	d.setColor()
}

// SetMessage sets the dialog message to display.
func (d *MessageDialog) SetMessage(message string) {
	d.message = "\n" + message
	d.textview.Clear()
	d.textview.SetText(d.message)
	d.textview.ScrollToBeginning()
	d.setRect()
}

// Focus is called when this primitive receives focus.
func (d *MessageDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(d.form)
}

// HasFocus returns whether or not this primitive has focus.
func (d *MessageDialog) HasFocus() bool {
	return d.form.HasFocus()
}

// SetRect sets rect for this primitive.
func (d *MessageDialog) SetRect(x, y, width, height int) {
	d.x = x
	d.y = y
	d.width = width
	d.height = height
	d.setRect()
}

// SetTextColor sets dialog's message text color.
func (d *MessageDialog) SetTextColor(color tcell.Color) {
	d.textview.SetTextColor(color)
}

// Draw draws this primitive onto the screen.
func (d *MessageDialog) Draw(screen tcell.Screen) {
	d.Box.DrawForSubclass(screen, d)
	x, y, width, height := d.Box.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// InputHandler returns input handler function for this primitive.
func (d *MessageDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if event.Key() == tcell.KeyDown || event.Key() == tcell.KeyUp || event.Key() == tcell.KeyPgDn || event.Key() == tcell.KeyPgUp { // nolint:lll
			if textHandler := d.textview.InputHandler(); textHandler != nil {
				textHandler(event, setFocus)

				return
			}
		}
		if formHandler := d.form.InputHandler(); formHandler != nil {
			formHandler(event, setFocus)

			return
		}
	})
}

// MouseHandler returns the mouse handler for this primitive.
func (d *MessageDialog) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) { // nolint:lll
	return d.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) { // nolint:lll,nonamedreturns
		// Pass mouse events on to the form.
		consumed, capture = d.form.MouseHandler()(action, event, setFocus)
		if !consumed && action == tview.MouseLeftClick && d.InRect(event.Position()) {
			setFocus(d)
			consumed = true
		}

		return consumed, capture
	})
}

// SetDoneFunc sets callback function for when user clicked on
// the the button or presses "enter" or "esc".
func (d *MessageDialog) SetDoneFunc(handler func()) *MessageDialog {
	d.doneHandler = handler
	enterButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	enterButton.SetSelectedFunc(handler)

	return d
}

func (d *MessageDialog) setColor() {
	var bgColor tcell.Color

	switch d.messageType {
	case InfoDialog:
		bgColor = d.bgColor
	case ErrorDailog:
		bgColor = tcell.ColorOrangeRed
	}

	d.form.SetBackgroundColor(bgColor)
	d.textview.SetBackgroundColor(bgColor)
	d.layout.SetBackgroundColor(bgColor)
}

func (d *MessageDialog) setRect() {
	maxHeight := d.height
	maxWidth := d.width // nolint:ifshort
	messageHeight := len(strings.Split(d.message, "\n"))
	messageWidth := getMessageWidth(d.message)

	layoutHeight := messageHeight

	if maxHeight > layoutHeight+dialogFormHeight {
		d.height = layoutHeight + dialogFormHeight + dialogPadding
	} else {
		d.height = maxHeight
		layoutHeight = d.height - dialogFormHeight - dialogPadding
	}

	if maxHeight > d.height {
		emptyHeight := (maxHeight - d.height) / emptySpaceParts
		d.y += emptyHeight
	}

	if d.width > messageWidth {
		d.width = messageWidth + dialogPadding
	}

	if maxWidth > d.width {
		emptyWidth := (maxWidth - d.width) / emptySpaceParts
		d.x += emptyWidth
	}

	d.layout.Clear()

	d.layout.AddItem(d.textview, layoutHeight, 0, true)
	d.layout.AddItem(d.form, dialogFormHeight, 0, true)

	d.Box.SetRect(d.x, d.y, d.width, d.height)
}
