package netdialogs

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/networks"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	netConnectDialogMaxWidth  = 60
	netConnectDialogMaxHeight = 17
	labelWidth                = 13
)

const (
	netConnectContainerFocus = 0 + iota
	netConnectAliasesFocus
	netConnectAliasesIPv4Focus
	netConnectAliasesIPv6Focus
	netConnectMacAddrFocus
	netConnectFormFocus
)

// NetworkConnectDialog implements network connect dialog primitive.
type NetworkConnectDialog struct {
	*tview.Box
	layout         *tview.Flex
	network        *tview.InputField
	container      *tview.DropDown
	aliases        *tview.InputField
	ipv4           *tview.InputField
	ipv6           *tview.InputField
	macAddr        *tview.InputField
	form           *tview.Form
	display        bool
	focusElement   int
	networkName    string
	connectHandler func()
	cancelHandler  func()
}

// NewNetworkConnectDialog returns a new network connect dialog primitive.
func NewNetworkConnectDialog() *NetworkConnectDialog {
	dialog := &NetworkConnectDialog{
		Box:       tview.NewBox(),
		layout:    tview.NewFlex(),
		network:   tview.NewInputField(),
		container: tview.NewDropDown(),
		aliases:   tview.NewInputField(),
		ipv4:      tview.NewInputField(),
		ipv6:      tview.NewInputField(),
		macAddr:   tview.NewInputField(),
		form:      tview.NewForm(),
	}

	bgColor := style.DialogBgColor
	fgColor := style.DialogFgColor
	inputFieldBgColor := style.InputFieldBgColor
	ddUnselectedStyle := style.DropDownUnselected
	ddselectedStyle := style.DropDownSelected

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

	// aliases input field
	dialog.aliases.SetBackgroundColor(bgColor)
	dialog.aliases.SetLabelColor(fgColor)
	dialog.aliases.SetLabel("alias:")
	dialog.aliases.SetLabelWidth(labelWidth)
	dialog.aliases.SetFieldBackgroundColor(inputFieldBgColor)

	// ipv4 input field
	dialog.ipv4.SetBackgroundColor(bgColor)
	dialog.ipv4.SetLabelColor(fgColor)
	dialog.ipv4.SetLabel("ipv4:")
	dialog.ipv4.SetLabelWidth(labelWidth)
	dialog.ipv4.SetFieldBackgroundColor(inputFieldBgColor)

	// ipv6 input field
	dialog.ipv6.SetBackgroundColor(bgColor)
	dialog.ipv6.SetLabelColor(fgColor)
	dialog.ipv6.SetLabel("ipv6:")
	dialog.ipv6.SetLabelWidth(labelWidth)
	dialog.ipv6.SetFieldBackgroundColor(inputFieldBgColor)

	// mac address input field
	dialog.macAddr.SetBackgroundColor(bgColor)
	dialog.macAddr.SetLabelColor(fgColor)
	dialog.macAddr.SetLabel("mac address:")
	dialog.macAddr.SetLabelWidth(labelWidth)
	dialog.macAddr.SetFieldBackgroundColor(inputFieldBgColor)

	// form
	dialog.form.AddButton("Cancel", nil)
	dialog.form.AddButton("Connect", nil)
	dialog.form.SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(bgColor)
	dialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// layout
	optionsLayout := tview.NewFlex().SetDirection(tview.FlexRow)

	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(dialog.network, 1, 0, true)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(dialog.container, 1, 0, true)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(dialog.aliases, 1, 0, true)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(dialog.ipv4, 1, 0, true)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(dialog.ipv6, 1, 0, true)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(dialog.macAddr, 1, 0, true)

	mainOptsLayout := tview.NewFlex().SetDirection(tview.FlexColumn)

	mainOptsLayout.SetBackgroundColor(bgColor)
	mainOptsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	mainOptsLayout.AddItem(optionsLayout, 0, 1, true)
	mainOptsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)

	dialog.layout.SetDirection(tview.FlexRow)
	dialog.layout.SetBackgroundColor(bgColor)
	dialog.layout.SetBorder(true)
	dialog.layout.SetBorderColor(style.DialogBorderColor)
	dialog.layout.SetTitle("PODMAN NETWORK CONNECT")
	dialog.layout.AddItem(mainOptsLayout, 0, 1, true)
	dialog.layout.AddItem(dialog.form, dialogs.DialogFormHeight, 0, true)

	return dialog
}

// Display displays this primitive.
func (d *NetworkConnectDialog) Display() {
	d.display = true
}

// IsDisplay returns true if this primitive is shown.
func (d *NetworkConnectDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *NetworkConnectDialog) Hide() {
	d.display = false
	d.focusElement = netConnectContainerFocus
	d.networkName = ""

	d.aliases.SetText("")
	d.ipv4.SetText("")
	d.ipv6.SetText("")
	d.macAddr.SetText("")
	d.SetNetworkInfo("", "")
	d.container.SetCurrentOption(0)
}

// HasFocus returns whether or not this primitive has focus.
func (d *NetworkConnectDialog) HasFocus() bool {
	if d.container.HasFocus() || d.aliases.HasFocus() {
		return true
	}

	if d.ipv4.HasFocus() || d.ipv6.HasFocus() {
		return true
	}

	if d.layout.HasFocus() || d.form.HasFocus() {
		return true
	}

	return d.Box.HasFocus()
}

// Focus is called when this primitive receives focus.
func (d *NetworkConnectDialog) Focus(delegate func(p tview.Primitive)) {
	switch d.focusElement {
	case netConnectContainerFocus:
		delegate(d.container)
	case netConnectAliasesFocus:
		delegate(d.aliases)
	case netConnectAliasesIPv4Focus:
		delegate(d.ipv4)
	case netConnectAliasesIPv6Focus:
		delegate(d.ipv6)
	case netConnectMacAddrFocus:
		delegate(d.macAddr)
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
func (d *NetworkConnectDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) { //nolint:gocognit,lll,cyclop
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("network connect dialog: event %v received", event)

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

		if d.aliases.HasFocus() {
			if aliasesHandler := d.aliases.InputHandler(); aliasesHandler != nil {
				aliasesHandler(event, setFocus)

				return
			}
		}

		if d.ipv4.HasFocus() {
			if ipv4Handler := d.ipv4.InputHandler(); ipv4Handler != nil {
				ipv4Handler(event, setFocus)

				return
			}
		}

		if d.ipv6.HasFocus() {
			if ipv6Handler := d.ipv6.InputHandler(); ipv6Handler != nil {
				ipv6Handler(event, setFocus)

				return
			}
		}

		if d.macAddr.HasFocus() {
			if macAddrHandler := d.macAddr.InputHandler(); macAddrHandler != nil {
				macAddrHandler(event, setFocus)

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
func (d *NetworkConnectDialog) SetRect(x, y, width, height int) {
	if width > netConnectDialogMaxWidth {
		emptySpace := (width - netConnectDialogMaxWidth) / 2 //nolint:gomnd
		x += emptySpace
		width = netConnectDialogMaxWidth
	}

	if height > netConnectDialogMaxHeight {
		emptySpace := (height - netConnectDialogMaxHeight) / 2 //nolint:gomnd
		y += emptySpace
		height = netConnectDialogMaxHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive into the screen.
func (d *NetworkConnectDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.Box.DrawForSubclass(screen, d)

	x, y, width, height := d.Box.GetInnerRect()

	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetConnectFunc sets form connect button selected function.
func (d *NetworkConnectDialog) SetConnectFunc(handler func()) *NetworkConnectDialog {
	d.connectHandler = handler
	connectButton := d.form.GetButton(d.form.GetButtonCount() - 1)

	connectButton.SetSelectedFunc(handler)

	return d
}

// SetCancelFunc sets form cancel button selected function.
func (d *NetworkConnectDialog) SetCancelFunc(handler func()) *NetworkConnectDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:gomnd

	cancelButton.SetSelectedFunc(handler)

	return d
}

func (d *NetworkConnectDialog) setFocusElement() {
	switch d.focusElement {
	case netConnectContainerFocus:
		d.focusElement = netConnectAliasesFocus
	case netConnectAliasesFocus:
		d.focusElement = netConnectAliasesIPv4Focus
	case netConnectAliasesIPv4Focus:
		d.focusElement = netConnectAliasesIPv6Focus
	case netConnectAliasesIPv6Focus:
		d.focusElement = netConnectMacAddrFocus
	case netConnectMacAddrFocus:
		d.focusElement = netConnectFormFocus
	}
}

// SetNetworkInfo sets selected network name in connect dialog.
func (d *NetworkConnectDialog) SetNetworkInfo(id, name string) {
	d.networkName = name
	network := fmt.Sprintf("%12s (%s)", id, name)

	d.network.SetText(network)
}

// SetContainers sets container drop down list content.
func (d *NetworkConnectDialog) SetContainers(cntList []entities.ListContainer) {
	containers := make([]string, 0)

	for _, cnt := range cntList {
		container := fmt.Sprintf("%s (%s)", cnt.ID[0:12], cnt.Names[0])
		containers = append(containers, container)
	}

	d.container.SetOptions(containers, nil)
}

// GetConnectOptions returns network connect options.
func (d *NetworkConnectDialog) GetConnectOptions() networks.NetworkConnect {
	var connectOptions networks.NetworkConnect

	_, selectedCnt := d.container.GetCurrentOption()
	container := strings.Split(selectedCnt, " ")[0]
	connectOptions.Container = container
	connectOptions.Network = d.networkName
	connectOptions.IPv4 = d.ipv4.GetText()
	connectOptions.IPv6 = d.ipv6.GetText()
	connectOptions.MacAddress = d.macAddr.GetText()
	connectOptions.Aliases = strings.Split(d.aliases.GetText(), " ")

	return connectOptions
}
