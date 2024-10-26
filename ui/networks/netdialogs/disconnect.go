package netdialogs

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	netDisconnectDialogMaxWidth  = 60
	netDisconnectDialogMaxHeight = 9
)

const (
	netDisconnectContainerFocus = 0 + iota
	netDisconnectFormFocus
)

// NetworkDisconnectDialog implements network disconnect dialog primitive.
type NetworkDisconnectDialog struct {
	*tview.Box
	layout            *tview.Flex
	network           *tview.InputField
	container         *tview.DropDown
	form              *tview.Form
	display           bool
	networkName       string
	focusElement      int
	disconnectHandler func()
	cancelHandler     func()
}

func NewNetworkDisconnectDialog() *NetworkDisconnectDialog {
	dialog := &NetworkDisconnectDialog{
		Box:       tview.NewBox(),
		layout:    tview.NewFlex(),
		network:   tview.NewInputField(),
		container: tview.NewDropDown(),
		form:      tview.NewForm(),
	}

	bgColor := style.DialogBgColor
	fgColor := style.DialogFgColor
	inputFieldBgColor := style.InputFieldBgColor
	ddUnselectedStyle := style.DropDownUnselected
	ddselectedStyle := style.DropDownSelected
	labelWidth := 12

	// network input field
	dialog.network.SetBackgroundColor(style.DialogBgColor)
	dialog.network.SetLabel("[::b]NETWORK ID:")
	dialog.network.SetLabelWidth(labelWidth)
	dialog.network.SetFieldBackgroundColor(style.DialogBgColor)
	dialog.network.SetLabelStyle(tcell.StyleDefault.
		Background(style.DialogBorderColor).
		Foreground(style.DialogFgColor))

	// container drop down
	dialog.container.SetBackgroundColor(bgColor)
	dialog.container.SetLabelColor(fgColor)
	dialog.container.SetLabel("container:")
	dialog.container.SetLabelWidth(labelWidth)
	dialog.container.SetOptions([]string{""}, nil)
	dialog.container.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	dialog.container.SetCurrentOption(0)
	dialog.container.SetFieldWidth(netConnectDialogMaxWidth)
	dialog.container.SetFieldBackgroundColor(inputFieldBgColor)

	// form
	dialog.form.AddButton(" Cancel ", nil)
	dialog.form.AddButton("Disconnect", nil)
	dialog.form.SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(bgColor)
	dialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// layout
	optionsLayout := tview.NewFlex().SetDirection(tview.FlexRow)

	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(dialog.network, 1, 0, true)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(dialog.container, 1, 0, true)

	mainOptsLayout := tview.NewFlex().SetDirection(tview.FlexColumn)

	mainOptsLayout.SetBackgroundColor(bgColor)
	mainOptsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	mainOptsLayout.AddItem(optionsLayout, 0, 1, true)
	mainOptsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)

	dialog.layout.SetDirection(tview.FlexRow)
	dialog.layout.SetBackgroundColor(bgColor)
	dialog.layout.SetBorder(true)
	dialog.layout.SetBorderColor(style.DialogBorderColor)
	dialog.layout.SetTitle("PODMAN NETWORK DISCONNECT")
	dialog.layout.AddItem(mainOptsLayout, 0, 1, true)
	dialog.layout.AddItem(dialog.form, dialogs.DialogFormHeight, 0, true)

	return dialog
}

// Display displays this primitive.
func (d *NetworkDisconnectDialog) Display() {
	d.display = true
}

// IsDisplay returns true if this primitive is shown.
func (d *NetworkDisconnectDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *NetworkDisconnectDialog) Hide() {
	d.display = false
	d.focusElement = netConnectContainerFocus
	d.networkName = ""

	d.SetNetworkInfo("", "")
	d.container.SetCurrentOption(0)
}

// HasFocus returns whether or not this primitive has focus.
func (d *NetworkDisconnectDialog) HasFocus() bool {
	if d.container.HasFocus() || d.layout.HasFocus() {
		return true
	}

	if d.Box.HasFocus() || d.form.HasFocus() {
		return true
	}

	return false
}

// Focus is called when this primitive receives focus.
func (d *NetworkDisconnectDialog) Focus(delegate func(p tview.Primitive)) {
	switch d.focusElement {
	case netConnectContainerFocus:
		delegate(d.container)
	case netConnectFormFocus:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)

		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == utils.SwitchFocusKey.Key {
				d.focusElement = netConnectContainerFocus

				d.Focus(delegate)
				d.form.SetFocus(0)

				return nil
			}

			return event
		})

		delegate(d.form)
	}
}

// InputHandler returns input handler function for this primitive.
func (d *NetworkDisconnectDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("network disconnect dialog: event %v received", event)

		if event.Key() == utils.CloseDialogKey.Key {
			if !d.container.HasFocus() {
				d.cancelHandler()

				return
			}
		}

		event = utils.ParseKeyEventKey(event)
		if event.Key() == utils.SwitchFocusKey.Key {
			d.setFocusElement()
		}

		// dropdown events
		if d.container.HasFocus() {
			if containerHandler := d.container.InputHandler(); containerHandler != nil {
				containerHandler(event, setFocus)

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
func (d *NetworkDisconnectDialog) SetRect(x, y, width, height int) {
	if width > netDisconnectDialogMaxWidth {
		emptySpace := (width - netDisconnectDialogMaxWidth) / 2 //nolint:mnd
		x += emptySpace
		width = netDisconnectDialogMaxWidth
	}

	if height > netDisconnectDialogMaxHeight {
		emptySpace := (height - netDisconnectDialogMaxHeight) / 2 //nolint:mnd
		y += emptySpace
		height = netDisconnectDialogMaxHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive into the screen.
func (d *NetworkDisconnectDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.Box.DrawForSubclass(screen, d)

	x, y, width, height := d.Box.GetInnerRect()

	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetDisconnectFunc sets form disconnect button selected function.
func (d *NetworkDisconnectDialog) SetDisconnectFunc(handler func()) *NetworkDisconnectDialog {
	d.disconnectHandler = handler
	connectButton := d.form.GetButton(d.form.GetButtonCount() - 1)

	connectButton.SetSelectedFunc(handler)

	return d
}

// SetCancelFunc sets form cancel button selected function.
func (d *NetworkDisconnectDialog) SetCancelFunc(handler func()) *NetworkDisconnectDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:mnd

	cancelButton.SetSelectedFunc(handler)

	return d
}

// SetContainers sets container drop down list content.
func (d *NetworkDisconnectDialog) SetContainers(cntList []entities.ListContainer) {
	containers := make([]string, 0)

	for _, cnt := range cntList {
		container := fmt.Sprintf("%s (%s)", cnt.ID[0:12], cnt.Names[0])
		containers = append(containers, container)
	}

	d.container.SetOptions(containers, nil)
}

func (d *NetworkDisconnectDialog) setFocusElement() {
	if d.focusElement == netDisconnectContainerFocus {
		d.focusElement = netConnectFormFocus
	}
}

// SetNetworkInfo sets selected network name in disconnect dialog.
func (d *NetworkDisconnectDialog) SetNetworkInfo(id string, name string) {
	d.networkName = name
	network := fmt.Sprintf("%12s (%s)", id, name)

	d.network.SetText(network)
}

// GetDisconnectOptions returns network disconnect options.
func (d *NetworkDisconnectDialog) GetDisconnectOptions() (string, string) {
	_, selectedCnt := d.container.GetCurrentOption()
	container := strings.Split(selectedCnt, " ")[0]

	return d.networkName, container
}
