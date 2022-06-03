package dialogs

import (
	"strings"

	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// ConfirmDialog is a simple confirmation dialog primitive
type ConfirmDialog struct {
	*tview.Box
	layout        *tview.Flex
	textview      *tview.TextView
	form          *tview.Form
	x             int
	y             int
	width         int
	height        int
	message       string
	display       bool
	cancelHandler func()
	selectHandler func()
}

// NewConfirmDialog returns new confirm dialog primitive
func NewConfirmDialog() *ConfirmDialog {
	dialog := &ConfirmDialog{
		Box:     tview.NewBox(),
		display: false,
	}

	bgColor := utils.Styles.ConfirmDialog.BgColor
	fgColor := utils.Styles.ConfirmDialog.FgColor

	dialog.textview = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)

	dialog.textview.SetBackgroundColor(bgColor)
	dialog.textview.SetTextColor(fgColor)

	dialog.form = tview.NewForm().
		AddButton("Cancel", nil).
		AddButton("  OK  ", nil).
		SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(bgColor)
	dialog.form.SetButtonBackgroundColor(utils.Styles.ConfirmDialog.ButtonColor)

	dialog.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	dialog.layout.SetBorder(true)
	dialog.layout.SetBorderColor(fgColor)
	dialog.layout.SetBackgroundColor(bgColor)

	return dialog

}

// Display displays this primitive
func (d *ConfirmDialog) Display() {
	d.form.SetFocus(1)
	d.display = true
}

// IsDisplay returns true if primitive is shown
func (d *ConfirmDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive
func (d *ConfirmDialog) Hide() {
	d.textview.SetText("")
	d.message = ""
	d.display = false
}

// SetTitle sets dialog title
func (d *ConfirmDialog) SetTitle(title string) {
	d.layout.SetTitle(strings.ToUpper(title))
	fgColor := utils.Styles.ConfirmDialog.FgColor
	d.layout.SetTitleColor(fgColor)
}

// SetText sets dialog title
func (d *ConfirmDialog) SetText(message string) {
	d.message = message
	d.textview.Clear()
	msg := "\n" + message
	d.textview.SetText(msg)
	d.textview.ScrollToBeginning()
	d.setRect()
}

// Focus is called when this primitive receives focus
func (d *ConfirmDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(d.form)
}

// HasFocus returns whether or not this primitive has focus
func (d *ConfirmDialog) HasFocus() bool {
	return d.form.HasFocus()
}

// SetRect set rects for this primitive.
func (d *ConfirmDialog) SetRect(x, y, width, height int) {
	d.x = x + DialogPadding
	d.y = y + DialogPadding
	d.width = width - (2 * DialogPadding)
	d.height = height - (2 * DialogPadding)
	d.setRect()
}

func (d *ConfirmDialog) setRect() {
	maxHeight := d.height
	maxWidth := d.width
	messageHeight := len(strings.Split(d.message, "\n"))
	messageWidth := getMessageWidth(d.message)

	layoutHeight := messageHeight + 2

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

	if d.width > DialogMinWidth {
		if messageWidth < DialogMinWidth {
			d.width = DialogMinWidth + 2
		} else if messageWidth < d.width {
			d.width = messageWidth + 2
		}
	}

	if maxWidth > d.width {
		emptyWidth := (maxWidth - d.width) / 2
		d.x = d.x + emptyWidth + DialogPadding
	}

	d.layout.Clear()
	d.layout.AddItem(d.textview, layoutHeight, 0, true)
	d.layout.AddItem(d.form, DialogFormHeight, 0, true)

	d.Box.SetRect(d.x, d.y, d.width, d.height)
}

// Draw draws this primitive onto the screen.
func (d *ConfirmDialog) Draw(screen tcell.Screen) {
	fgColor := utils.Styles.ConfirmDialog.FgColor
	bgColor := utils.Styles.ConfirmDialog.BgColor
	if !d.display {
		return
	}
	d.Box.DrawForSubclass(screen, d)
	x, y, width, height := d.Box.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.SetBorder(true)
	d.layout.SetBackgroundColor(bgColor)
	d.layout.SetBorderColor(fgColor)
	d.layout.Draw(screen)
}

// InputHandler returns input handler function for this primitive
func (d *ConfirmDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("confirm dialog: event %v received", event)
		if event.Key() == tcell.KeyEsc {
			d.cancelHandler()
			return
		}
		if formHandler := d.form.InputHandler(); formHandler != nil {
			formHandler(event, setFocus)
			return
		}
	})
}

// SetCancelFunc sets form cancel button selected function
func (d *ConfirmDialog) SetCancelFunc(handler func()) *ConfirmDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2)
	cancelButton.SetSelectedFunc(handler)
	return d
}

// SetSelectedFunc sets form select button selected function
func (d *ConfirmDialog) SetSelectedFunc(handler func()) *ConfirmDialog {
	d.selectHandler = handler
	enterButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	enterButton.SetSelectedFunc(handler)
	return d
}
