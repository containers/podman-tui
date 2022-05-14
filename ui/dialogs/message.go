package dialogs

import (
	"strings"

	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// MessageDialog is a simaple message dialog primitive
type MessageDialog struct {
	*tview.Box
	layout        *tview.Flex
	textview      *tview.TextView
	form          *tview.Form
	display       bool
	message       string
	cancelHandler func()
	selectHandler func()
}

// NewMessageDialog returns new message dialog primitive
func NewMessageDialog(text string) *MessageDialog {
	dialog := &MessageDialog{
		Box:     tview.NewBox(),
		display: false,
		message: text,
	}

	dialog.textview = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)

	bgColor := utils.Styles.MessageDialog.BgColor
	terminalBgColor := utils.Styles.MessageDialog.Terminal.BgColor
	terminalFgColor := utils.Styles.MessageDialog.Terminal.FgColor
	terminalBorderColor := utils.Styles.MessageDialog.Terminal.BorderColor
	buttonBgColor := utils.Styles.ButtonPrimitive.BgColor

	dialog.textview.SetTextColor(terminalFgColor)
	dialog.textview.SetBackgroundColor(terminalBgColor)
	dialog.textview.SetBorderColor(terminalBorderColor)
	dialog.textview.SetBorder(true)

	// textview layout
	tlayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	tlayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	tlayout.AddItem(dialog.textview, 0, 1, true)
	tlayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)

	dialog.form = tview.NewForm().
		AddButton("Enter", nil).
		SetButtonsAlign(tview.AlignRight)

	dialog.form.SetBackgroundColor(bgColor)
	dialog.form.SetButtonBackgroundColor(buttonBgColor)

	dialog.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	dialog.layout.AddItem(tlayout, 0, 1, true)
	dialog.layout.AddItem(dialog.form, DialogFormHeight, 0, true)
	dialog.layout.SetBorder(true)
	dialog.layout.SetBackgroundColor(bgColor)

	return dialog
}

// Display displays this primitive
func (d *MessageDialog) Display() {
	d.display = true
}

// IsDisplay returns true if primitive is shown
func (d *MessageDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive
func (d *MessageDialog) Hide() {
	d.message = ""
	d.textview.SetText("")
	d.display = false
}

// SetTitle sets input dialog title
func (d *MessageDialog) SetTitle(title string) {
	d.layout.SetTitle(strings.ToUpper(title))
}

// SetText sets message dialog text messages
func (d *MessageDialog) SetText(message string) {
	d.message = message
	d.textview.Clear()
	d.textview.SetText(message)
	d.textview.ScrollToBeginning()
}

// TextScrollToEnd scroll downs the text view
func (d *MessageDialog) TextScrollToEnd() {
	d.textview.ScrollToEnd()
}

// Focus is called when this primitive receives focus
func (d *MessageDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(d.form)
}

// HasFocus returns whether or not this primitive has focus
func (d *MessageDialog) HasFocus() bool {
	return d.form.HasFocus()
}

// SetRect set rects for this primitive.
func (d *MessageDialog) SetRect(x, y, width, height int) {
	messageHeight := len(strings.Split(d.message, "\n")) + 1
	messageWidth := getMessageWidth(d.message)

	dWidth := width - (2 * DialogPadding)
	if messageWidth+4 < dWidth {
		dWidth = messageWidth + 4
	}
	if DialogMinWidth < width && dWidth < DialogMinWidth {
		dWidth = DialogMinWidth
	}
	emptySpace := (width - dWidth) / 2
	dX := x + emptySpace

	dHeight := messageHeight + DialogFormHeight + DialogPadding
	if dHeight > height {
		dHeight = height - DialogPadding - 1
	}
	textviewHeight := dHeight - DialogFormHeight - 2
	hs := ((height - dHeight) / 2)
	dY := y + hs

	d.Box.SetRect(dX, dY, dWidth, dHeight)
	//set text view height size
	d.layout.ResizeItem(d.textview, textviewHeight, 0)

}

// Draw draws this primitive onto the screen.
func (d *MessageDialog) Draw(screen tcell.Screen) {

	if !d.display {
		return
	}

	d.Box.DrawForSubclass(screen, d)
	x, y, width, height := d.Box.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// InputHandler returns input handler function for this primitive
func (d *MessageDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("message dialog: event %v received", event)
		if event.Key() == tcell.KeyEsc {
			d.cancelHandler()
			return
		}
		if event.Key() == tcell.KeyEnter {
			d.selectHandler()
			return
		}
		if event.Key() == tcell.KeyTab {
			if formHandler := d.form.InputHandler(); formHandler != nil {
				formHandler(event, setFocus)
				return
			}
		}
		// scroll between message textview
		if textHandler := d.textview.InputHandler(); textHandler != nil {
			textHandler(event, setFocus)
			return
		}
	})
}

// SetSelectedFunc sets form enter button selected function
func (d *MessageDialog) SetSelectedFunc(handler func()) *MessageDialog {
	d.selectHandler = handler
	enterButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	enterButton.SetSelectedFunc(handler)
	return d
}

// SetCancelFunc sets form cancel button selected function
func (d *MessageDialog) SetCancelFunc(handler func()) *MessageDialog {
	d.cancelHandler = handler
	return d
}
