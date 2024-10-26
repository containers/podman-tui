package sysdialogs

import (
	"strings"

	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	connCreateDialogMaxWidth  int = 100
	connCreateDialogMaxHeight int = 11
)

const (
	connNameFieldFocus = 0 + iota
	connURIFieldFocus
	connIdentityFieldFocus
	connFormFocus
)

// AddConnectionDialog implements new connection create dialog.
type AddConnectionDialog struct {
	*tview.Box
	layout               *tview.Flex
	connNameField        *tview.InputField
	connURIField         *tview.InputField
	identityField        *tview.InputField
	form                 *tview.Form
	focusElement         int
	display              bool
	cancelHandler        func()
	addConnectionHandler func()
}

// NewAddConnectionDialog returns a new connection create dialog primitive.
func NewAddConnectionDialog() *AddConnectionDialog {
	connDialog := AddConnectionDialog{
		Box:     tview.NewBox().SetBorder(false),
		layout:  tview.NewFlex().SetDirection(tview.FlexRow),
		display: false,
	}

	labelWidth := 10
	// connection name
	connDialog.connNameField = tview.NewInputField()
	connDialog.connNameField.SetLabel("Name:")
	connDialog.connNameField.SetLabelWidth(labelWidth)
	connDialog.connNameField.SetBackgroundColor(style.DialogBgColor)
	connDialog.connNameField.SetLabelColor(style.DialogFgColor)
	connDialog.connNameField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// connection URI
	connDialog.connURIField = tview.NewInputField()
	connDialog.connURIField.SetLabel("URI:")
	connDialog.connURIField.SetLabelWidth(labelWidth)
	connDialog.connURIField.SetBackgroundColor(style.DialogBgColor)
	connDialog.connURIField.SetLabelColor(style.DialogFgColor)
	connDialog.connURIField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// identity
	connDialog.identityField = tview.NewInputField()
	connDialog.identityField.SetLabel("Identity:")
	connDialog.identityField.SetLabelWidth(labelWidth)
	connDialog.identityField.SetBackgroundColor(style.DialogBgColor)
	connDialog.identityField.SetLabelColor(style.DialogFgColor)
	connDialog.identityField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// form
	connDialog.form = tview.NewForm().
		AddButton("Cancel", nil).
		AddButton(" Add ", nil).
		SetButtonsAlign(tview.AlignRight)
	connDialog.form.SetBackgroundColor(style.DialogBgColor)
	connDialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// layouts
	inputFieldLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	inputFieldLayout.SetBackgroundColor(style.DialogBgColor)
	inputFieldLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	inputFieldLayout.AddItem(connDialog.connNameField, 1, 0, true)
	inputFieldLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	inputFieldLayout.AddItem(connDialog.connURIField, 1, 0, true)
	inputFieldLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	inputFieldLayout.AddItem(connDialog.identityField, 1, 0, true)
	// adding an empty column space to beginning and end of the fields layout
	layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	layout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	layout.AddItem(inputFieldLayout, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)

	connDialog.layout.SetBorder(true)
	connDialog.layout.SetBorderColor(style.DialogBorderColor)
	connDialog.layout.SetTitle("ADD NEW SYSTEM CONNECTION")
	connDialog.layout.SetBackgroundColor(style.DialogBgColor)
	connDialog.layout.AddItem(layout, 0, 1, true)
	connDialog.layout.AddItem(connDialog.form, dialogs.DialogFormHeight, 0, true)

	// returns the command primitive
	return &connDialog
}

// Display displays this primitive.
func (addDialog *AddConnectionDialog) Display() {
	addDialog.focusElement = 0
	addDialog.connNameField.SetText("")
	addDialog.connURIField.SetText("")
	addDialog.identityField.SetText("")
	addDialog.display = true
}

// IsDisplay returns true if primitive is shown.
func (addDialog *AddConnectionDialog) IsDisplay() bool {
	return addDialog.display
}

// Hide stops displaying this primitive.
func (addDialog *AddConnectionDialog) Hide() {
	addDialog.display = false
}

// HasFocus returns whether or not this primitive has focus.
func (addDialog *AddConnectionDialog) HasFocus() bool {
	if addDialog.connNameField.HasFocus() || addDialog.connURIField.HasFocus() {
		return true
	}

	if addDialog.identityField.HasFocus() || addDialog.layout.HasFocus() {
		return true
	}

	if addDialog.layout.HasFocus() || addDialog.form.HasFocus() {
		return true
	}

	return addDialog.Box.HasFocus()
}

// Focus is called when this primitive receives focus.
func (addDialog *AddConnectionDialog) Focus(delegate func(p tview.Primitive)) {
	switch addDialog.focusElement {
	case connNameFieldFocus:
		addDialog.connNameField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == utils.SwitchFocusKey.EventKey() {
				addDialog.focusElement = connURIFieldFocus
				addDialog.Focus(delegate)

				return nil
			}

			return event
		})

		delegate(addDialog.connNameField)
	case connURIFieldFocus:
		addDialog.connURIField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == utils.SwitchFocusKey.EventKey() {
				addDialog.focusElement = connIdentityFieldFocus
				addDialog.Focus(delegate)

				return nil
			}

			return event
		})

		delegate(addDialog.connURIField)
	case connIdentityFieldFocus:
		addDialog.identityField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == utils.SwitchFocusKey.EventKey() {
				addDialog.focusElement = connFormFocus
				addDialog.Focus(delegate)

				return nil
			}

			return event
		})

		delegate(addDialog.identityField)
	case connFormFocus:
		button := addDialog.form.GetButton(addDialog.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == utils.SwitchFocusKey.EventKey() {
				addDialog.focusElement = connNameFieldFocus
				addDialog.Focus(delegate)
				addDialog.form.SetFocus(0)

				return nil
			}

			return event
		})

		delegate(addDialog.form)
	}
}

// InputHandler returns input handler function for this primitive.
func (addDialog *AddConnectionDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return addDialog.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("connection create dialog: event %v received", event)

		if event.Key() == utils.CloseDialogKey.EventKey() {
			addDialog.cancelHandler()

			return
		}
		// connection name field
		if addDialog.connNameField.HasFocus() {
			if inputHandler := addDialog.connNameField.InputHandler(); inputHandler != nil {
				inputHandler(event, setFocus)

				return
			}
		}
		// connection URI field
		if addDialog.connURIField.HasFocus() {
			if inputHandler := addDialog.connURIField.InputHandler(); inputHandler != nil {
				inputHandler(event, setFocus)

				return
			}
		}
		// identity field handler
		if addDialog.identityField.HasFocus() {
			if inputHandler := addDialog.identityField.InputHandler(); inputHandler != nil {
				inputHandler(event, setFocus)

				return
			}
		}
		// form handler
		if addDialog.form.HasFocus() {
			if formHandler := addDialog.form.InputHandler(); formHandler != nil {
				formHandler(event, setFocus)

				return
			}
		}
	})
}

// SetAddFunc sets form add button selected function.
func (addDialog *AddConnectionDialog) SetAddFunc(handler func()) *AddConnectionDialog {
	addDialog.addConnectionHandler = handler
	addButton := addDialog.form.GetButton(addDialog.form.GetButtonCount() - 1)
	addButton.SetSelectedFunc(handler)

	return addDialog
}

// SetCancelFunc sets form cancel button selected function.
func (addDialog *AddConnectionDialog) SetCancelFunc(handler func()) *AddConnectionDialog {
	addDialog.cancelHandler = handler
	cancelButton := addDialog.form.GetButton(addDialog.form.GetButtonCount() - 2) //nolint:mnd
	cancelButton.SetSelectedFunc(handler)

	return addDialog
}

// SetRect set rects for this primitive.
func (addDialog *AddConnectionDialog) SetRect(x, y, width, height int) {
	dWidth := width
	if width > connCreateDialogMaxWidth {
		dWidth = connCreateDialogMaxWidth
	}

	dBWidth := dWidth - (2 * dialogs.DialogPadding) //nolint:mnd

	widthEmptySpace := (width - dWidth) / 2 //nolint:mnd

	x = x + widthEmptySpace + dialogs.DialogPadding

	dHeight := height
	if height > connCreateDialogMaxHeight {
		dHeight = connCreateDialogMaxHeight
	}

	heightEmptySpace := (height - dHeight) / 2 //nolint:mnd
	y += heightEmptySpace

	addDialog.Box.SetRect(x, y, dBWidth, dHeight)

	x, y, width, height = addDialog.Box.GetInnerRect()

	addDialog.layout.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen.
func (addDialog *AddConnectionDialog) Draw(screen tcell.Screen) {
	if !addDialog.display {
		return
	}

	addDialog.Box.DrawForSubclass(screen, addDialog)
	addDialog.layout.Draw(screen)
}

// GetItems returns new connection name, uri and identity.
func (addDialog *AddConnectionDialog) GetItems() (string, string, string) {
	var (
		name     string
		uri      string
		identity string
	)

	name = addDialog.connNameField.GetText()
	name = strings.TrimSpace(name)

	uri = addDialog.connURIField.GetText()
	uri = strings.TrimSpace(uri)

	identity = addDialog.identityField.GetText()
	identity = strings.TrimSpace(identity)

	return name, uri, identity
}
