package dialogs

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// MessageDialog is a simaple messsage dialog primitive
type MessageDialog struct {
	*tview.Box
	layout        *tview.Flex
	textview      *tview.TextView
	form          *tview.Form
	display       bool
	message       string
	x             int
	y             int
	width         int
	height        int
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

	bgColor := utils.Styles.CommandDialog.BgColor
	dialog.textview.SetTextColor(tcell.ColorBlack)
	dialog.textview.SetBackgroundColor(bgColor)
	dialog.textview.SetBorderColor(utils.Styles.CommandDialog.HeaderRow.BgColor)
	dialog.textview.SetBorder(true)

	dialog.form = tview.NewForm().
		AddButton("Enter", nil).
		SetButtonsAlign(tview.AlignRight)

	dialog.form.SetBackgroundColor(bgColor)

	dialog.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	dialog.layout.AddItem(dialog.textview, len(text), 0, true)
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
	d.setRect()
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

func (d *MessageDialog) setRect() {
	maxHeight := d.height
	maxWidth := d.width
	messageHeight := len(strings.Split(d.message, "\n"))
	messageWidth := getMessageWidth(d.message)

	if d.width > DialogMinWidth {
		if messageWidth < DialogMinWidth {
			d.width = DialogMinWidth + 2
		}
	}

	if maxWidth > d.width {
		emptyWidth := (maxWidth - d.width) / 2
		d.x = d.x + emptyWidth + DialogPadding
	}

	layoutHeight := messageHeight + 2
	for _, line := range strings.Split(d.message, "\n") {
		if len(line) > maxWidth-2 {
			layoutHeight = layoutHeight + 1
		}
	}

	if maxHeight > layoutHeight+DialogFormHeight {
		d.height = layoutHeight + DialogFormHeight + 2
	} else {
		d.height = maxHeight
		layoutHeight = d.height - DialogFormHeight - 2
	}

	if maxHeight > d.height {
		emptyHeight := (maxHeight - d.height) / 2
		d.y = d.y + emptyHeight - DialogPadding

	}

	d.layout.Clear()

	bgColor := utils.Styles.CommandDialog.BgColor
	d.layout.AddItem(d.textview, layoutHeight, 0, true)
	d.layout.AddItem(d.form, DialogFormHeight, 0, true)
	d.layout.SetBorder(true)
	d.layout.SetBackgroundColor(bgColor)

	d.Box.SetRect(d.x, d.y, d.width, d.height)
}

// SetRect set rects for this primitive.
func (d *MessageDialog) SetRect(x, y, width, height int) {
	d.x = x + DialogPadding
	d.y = y + DialogPadding
	d.width = width - (2 * DialogPadding)
	d.height = height - (2 * DialogPadding)
	d.setRect()
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
		log.Debug().Msgf("message dialog: event %v received", event.Key())
		if event.Key() == tcell.KeyEsc {
			d.cancelHandler()
			return
		}
		if event.Key() == tcell.KeyEnter {
			d.selectHandler()
			return
		}
		if event.Key() == tcell.KeyLeft || event.Key() == tcell.KeyRight || event.Key() == tcell.KeyTab {
			if formHandler := d.form.InputHandler(); formHandler != nil {
				formHandler(event, setFocus)
				return
			}
		}
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
