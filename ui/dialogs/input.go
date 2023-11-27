package dialogs

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	siInputElement     = 1
	siFormElement      = 2
	siDescHeight       = 4
	siDialogWidth      = 60
	siDialogHeight     = 10
	siDialogInputWidth = 57
)

// SimpleInputDialog is an input dialog primitive.
type SimpleInputDialog struct {
	*tview.Box
	height        int
	layout        *tview.Flex
	textview      *tview.TextView
	input         *tview.InputField
	inputWidth    int
	form          *tview.Form
	focusElement  int
	display       bool
	cancelHandler func()
	selectHandler func()
}

// NewSimpleInputDialog returns new input dialog primitive.
func NewSimpleInputDialog(text string) *SimpleInputDialog {
	dialog := &SimpleInputDialog{
		Box:          tview.NewBox(),
		display:      false,
		height:       siDialogHeight,
		focusElement: siInputElement,
	}

	bgColor := style.DialogBgColor
	fgColor := style.DialogFgColor

	dialog.textview = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)

	dialog.textview.SetBackgroundColor(bgColor)

	dialog.form = tview.NewForm().
		AddButton("Cancel", nil).
		AddButton("Enter", nil).
		SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(bgColor)
	dialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	dialog.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	dialog.layout.SetBorder(true)
	dialog.layout.SetBackgroundColor(bgColor)
	dialog.layout.SetBorderColor(style.DialogBorderColor)

	dialog.input = tview.NewInputField()
	dialog.SetInputText(text)
	dialog.input.SetLabelColor(fgColor)
	dialog.input.SetBackgroundColor(bgColor)
	dialog.input.SetFieldBackgroundColor(style.InputFieldBgColor)

	dialog.setLayout(false)

	return dialog
}

// Display displays this primitive.
func (d *SimpleInputDialog) Display() {
	d.focusElement = siInputElement
	d.display = true
}

// IsDisplay returns true if primitive is shown.
func (d *SimpleInputDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *SimpleInputDialog) Hide() {
	d.focusElement = 0
	d.input.SetText("")
	d.display = false
}

func (d *SimpleInputDialog) setLayout(haveDesc bool) {
	d.layout.Clear()

	descHeight := siDescHeight

	if !haveDesc {
		descHeight = 1
		d.height = siDialogHeight - 3 //nolint:gomnd
	} else {
		d.height = siDialogHeight
	}

	d.layout.AddItem(
		tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false).
			AddItem(d.textview, 0, 1, false).
			AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false),
		descHeight, 0, true)

	d.layout.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false).
		AddItem(d.input, 0, 1, false).
		AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false),
		1, 0, true)

	d.layout.AddItem(d.form, DialogFormHeight, 0, true)
	d.layout.SetBorder(true)
	d.layout.SetBorderColor(style.DialogBorderColor)
	d.layout.SetBackgroundColor(style.DialogBgColor)
}

// SetDescription sets dialogs description.
func (d *SimpleInputDialog) SetDescription(text string) {
	d.textview.Clear()

	haveDesc := true

	fmt.Fprintf(d.textview, "\n%s", text)

	if len(text) == 0 {
		haveDesc = false
	}

	d.setLayout(haveDesc)
}

// SetSelectButtonLabel sets form select/enter button name.
func (d *SimpleInputDialog) SetSelectButtonLabel(label string) {
	if len(label) == 0 {
		return
	}

	button := d.form.GetButton(d.form.GetButtonCount() - 1)
	buttonLabel := strings.ToUpper(label[0:1])

	if len(label) > 1 {
		buttonLabel += label[1:]
	}

	button.SetLabel(buttonLabel)
}

// SetTitle sets input dialog title.
func (d *SimpleInputDialog) SetTitle(title string) {
	d.layout.SetTitle(strings.ToUpper(title))
}

// GetInputText returns input dialog input field value.
func (d *SimpleInputDialog) GetInputText() string {
	return d.input.GetText()
}

// SetInputText sets input dialog default value.
func (d *SimpleInputDialog) SetInputText(text string) {
	d.input.SetText(text)
}

// SetLabel sets input fields label message.
func (d *SimpleInputDialog) SetLabel(text string) {
	width := len(text) + 2 //nolint:gomnd
	d.inputWidth = siDialogInputWidth - width

	d.input.SetFieldWidth(d.inputWidth)

	label := fmt.Sprintf("%s: ", text)

	d.input.SetLabel(label)
}

// HasFocus returns whether or not this primitive has focus.
func (d *SimpleInputDialog) HasFocus() bool {
	if d.input.HasFocus() || d.form.HasFocus() {
		return true
	}

	return d.Box.HasFocus()
}

// Focus is called when this primitive receives focus.
func (d *SimpleInputDialog) Focus(delegate func(p tview.Primitive)) {
	inputHandler := func(key tcell.Key) {
		if key == tcell.KeyTab {
			d.focusElement = siFormElement
			d.form.SetFocus(siFormElement)
			d.Focus(delegate)
		}
	}

	switch d.focusElement {
	case siInputElement:
		d.input.SetDoneFunc(inputHandler)
		delegate(d.input)
	case siFormElement:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = siInputElement
				d.Focus(delegate)

				return nil
			}

			return event
		})

		delegate(d.form)
	}
}

// InputHandler returns input handler function for this primitive.
func (d *SimpleInputDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("input dialog: event %v received", event)
		if event.Key() == tcell.KeyEsc {
			d.cancelHandler()

			return
		}

		if event.Key() == tcell.KeyEnter && !d.form.HasFocus() {
			d.selectHandler()

			return
		}

		if d.input.HasFocus() {
			if inputHandler := d.input.InputHandler(); inputHandler != nil {
				inputHandler(event, setFocus)

				return
			}
		}

		if d.form.HasFocus() {
			if formHandler := d.form.InputHandler(); formHandler != nil {
				formHandler(event, setFocus)

				return
			}
		}
	})
}

// SetRect set rects for this primitive.
func (d *SimpleInputDialog) SetRect(x, y, width, height int) {
	ws := (width - siDialogWidth) / 2 //nolint:gomnd
	hs := ((height - d.height) / 2)   //nolint:gomnd
	dy := y + hs
	bWidth := siDialogWidth

	if siDialogWidth > width {
		ws = 0
		bWidth = width - 1
	}

	bHeight := d.height

	if d.height >= height {
		dy = y + 1
		bHeight = height - 1
	}

	d.Box.SetRect(x+ws, dy, bWidth, bHeight)
}

// GetRect returns the current position of the primitive, x, y, width, and
// height.
func (d *SimpleInputDialog) GetRect() (int, int, int, int) {
	return d.Box.GetRect()
}

// Draw draws this primitive onto the screen.
func (d *SimpleInputDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.Box.DrawForSubclass(screen, d)
	x, y, width, height := d.Box.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetSelectedFunc sets form enter button selected function.
func (d *SimpleInputDialog) SetSelectedFunc(handler func()) *SimpleInputDialog {
	d.selectHandler = handler
	enterButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	enterButton.SetSelectedFunc(handler)

	return d
}

// SetCancelFunc sets form cancel button selected function.
func (d *SimpleInputDialog) SetCancelFunc(handler func()) *SimpleInputDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:gomnd
	cancelButton.SetSelectedFunc(handler)

	return d
}
