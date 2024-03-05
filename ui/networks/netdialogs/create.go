package netdialogs

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/networks"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	networkCreateDialogMaxWidth = 80
	networkCreateDialogHeight   = 17
)

const (
	formFocus = 0 + iota
	categoriesFocus
	categoryPagesFocus
	networkNameFieldFocus
	networkLabelFieldFocus
	networkInternalCheckBoxFocus
	networkDriverFieldFocus
	networkDriverOptionsFieldFocus
	networkIPv6CheckBoxFocus
	networkGatewatFieldFocus
	networkIPRangeFieldFocus
	networkSubnetFieldFocus
	networkDisableDNSCheckBoxFocus
)

const (
	basicInfoPageIndex = 0 + iota
	ipSettingsPageIndex
)

// NetworkCreateDialog implements network create dialog.
type NetworkCreateDialog struct {
	*tview.Box
	layout                    *tview.Flex
	categoryLabels            []string
	categories                *tview.TextView
	categoryPages             *tview.Pages
	basicInfoPage             *tview.Flex
	ipSettingsPage            *tview.Flex
	form                      *tview.Form
	display                   bool
	activePageIndex           int
	focusElement              int
	networkNameField          *tview.InputField
	networkLabelsField        *tview.InputField
	networkInternalCheckBox   *tview.Checkbox
	networkDriverField        *tview.InputField
	networkDriverOptionsField *tview.InputField
	networkIpv6CheckBox       *tview.Checkbox
	networkGatewayField       *tview.InputField
	networkIPRangeField       *tview.InputField
	networkSubnetField        *tview.InputField
	networkDisableDNSCheckBox *tview.Checkbox
	cancelHandler             func()
	createHandler             func()
}

// NewNetworkCreateDialog returns new network create dialog primitive NetworkCreateDialog.
func NewNetworkCreateDialog() *NetworkCreateDialog {
	netDialog := NetworkCreateDialog{
		Box:                       tview.NewBox(),
		layout:                    tview.NewFlex().SetDirection(tview.FlexRow),
		categories:                tview.NewTextView(),
		categoryPages:             tview.NewPages(),
		basicInfoPage:             tview.NewFlex(),
		ipSettingsPage:            tview.NewFlex(),
		form:                      tview.NewForm(),
		categoryLabels:            []string{"Basic Information", "IP Settings"},
		activePageIndex:           0,
		display:                   false,
		networkNameField:          tview.NewInputField(),
		networkLabelsField:        tview.NewInputField(),
		networkInternalCheckBox:   tview.NewCheckbox(),
		networkDriverField:        tview.NewInputField(),
		networkDriverOptionsField: tview.NewInputField(),
		networkIpv6CheckBox:       tview.NewCheckbox(),
		networkGatewayField:       tview.NewInputField(),
		networkIPRangeField:       tview.NewInputField(),
		networkSubnetField:        tview.NewInputField(),
		networkDisableDNSCheckBox: tview.NewCheckbox(),
	}

	bgColor := style.DialogBgColor
	fgColor := style.DialogFgColor
	inputFieldBgColor := style.InputFieldBgColor
	buttonBgColor := style.ButtonBgColor

	netDialog.categories.SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)
	netDialog.categories.SetBackgroundColor(bgColor)
	netDialog.categories.SetBorder(true)
	netDialog.categories.SetBorderColor(style.DialogSubBoxBorderColor)

	// basic information setup page
	basicInfoPageLabelWidth := 12
	// name field
	netDialog.networkNameField.SetLabel("name:")
	netDialog.networkNameField.SetLabelWidth(basicInfoPageLabelWidth)
	netDialog.networkNameField.SetBackgroundColor(bgColor)
	netDialog.networkNameField.SetLabelColor(fgColor)
	netDialog.networkNameField.SetFieldBackgroundColor(inputFieldBgColor)

	// labels field
	netDialog.networkLabelsField.SetLabel("labels:")
	netDialog.networkLabelsField.SetLabelWidth(basicInfoPageLabelWidth)
	netDialog.networkLabelsField.SetBackgroundColor(bgColor)
	netDialog.networkLabelsField.SetLabelColor(fgColor)
	netDialog.networkLabelsField.SetFieldBackgroundColor(inputFieldBgColor)

	// internal check box
	netDialog.networkInternalCheckBox.SetLabel("internal")
	netDialog.networkInternalCheckBox.SetLabelWidth(basicInfoPageLabelWidth)
	netDialog.networkInternalCheckBox.SetChecked(false)
	netDialog.networkInternalCheckBox.SetBackgroundColor(bgColor)
	netDialog.networkInternalCheckBox.SetLabelColor(fgColor)
	netDialog.networkInternalCheckBox.SetFieldBackgroundColor(inputFieldBgColor)

	// drivers
	netDialog.networkDriverField.SetLabel("drivers:")
	netDialog.networkDriverField.SetLabelWidth(basicInfoPageLabelWidth)
	netDialog.networkDriverField.SetBackgroundColor(bgColor)
	netDialog.networkDriverField.SetLabelColor(fgColor)
	netDialog.networkDriverField.SetFieldBackgroundColor(inputFieldBgColor)

	// drivers options
	netDialog.networkDriverOptionsField.SetLabel("options:")
	netDialog.networkDriverOptionsField.SetLabelWidth(basicInfoPageLabelWidth)
	netDialog.networkDriverOptionsField.SetBackgroundColor(bgColor)
	netDialog.networkDriverOptionsField.SetLabelColor(fgColor)
	netDialog.networkDriverOptionsField.SetFieldBackgroundColor(inputFieldBgColor)

	// ip settings page
	ipSettingsPageLabelWidth := 12
	// ipv6 check box
	netDialog.networkIpv6CheckBox.SetLabel("ipv6")
	netDialog.networkIpv6CheckBox.SetLabelWidth(ipSettingsPageLabelWidth)
	netDialog.networkIpv6CheckBox.SetChecked(false)
	netDialog.networkIpv6CheckBox.SetBackgroundColor(bgColor)
	netDialog.networkIpv6CheckBox.SetLabelColor(tcell.ColorWhite)
	netDialog.networkIpv6CheckBox.SetFieldBackgroundColor(inputFieldBgColor)

	// gateway
	netDialog.networkGatewayField.SetLabel("gateway:")
	netDialog.networkGatewayField.SetLabelWidth(basicInfoPageLabelWidth)
	netDialog.networkGatewayField.SetBackgroundColor(bgColor)
	netDialog.networkGatewayField.SetLabelColor(tcell.ColorWhite)
	netDialog.networkGatewayField.SetFieldBackgroundColor(inputFieldBgColor)

	// ip range
	netDialog.networkIPRangeField.SetLabel("ip range:")
	netDialog.networkIPRangeField.SetLabelWidth(basicInfoPageLabelWidth)
	netDialog.networkIPRangeField.SetBackgroundColor(bgColor)
	netDialog.networkIPRangeField.SetLabelColor(tcell.ColorWhite)
	netDialog.networkIPRangeField.SetFieldBackgroundColor(inputFieldBgColor)

	// subnet
	netDialog.networkSubnetField.SetLabel("subnet:")
	netDialog.networkSubnetField.SetLabelWidth(basicInfoPageLabelWidth)
	netDialog.networkSubnetField.SetBackgroundColor(bgColor)
	netDialog.networkSubnetField.SetLabelColor(tcell.ColorWhite)
	netDialog.networkSubnetField.SetFieldBackgroundColor(inputFieldBgColor)

	// dns check box
	netDialog.networkDisableDNSCheckBox.SetLabel("disable DNS")
	netDialog.networkDisableDNSCheckBox.SetLabelWidth(basicInfoPageLabelWidth)
	netDialog.networkDisableDNSCheckBox.SetChecked(false)
	netDialog.networkDisableDNSCheckBox.SetBackgroundColor(bgColor)
	netDialog.networkDisableDNSCheckBox.SetLabelColor(tcell.ColorWhite)
	netDialog.networkDisableDNSCheckBox.SetFieldBackgroundColor(inputFieldBgColor)

	// category pages
	netDialog.categoryPages.SetBackgroundColor(bgColor)
	netDialog.categoryPages.SetBorder(true)
	netDialog.categoryPages.SetBorderColor(style.DialogSubBoxBorderColor)

	// form
	netDialog.form.SetBackgroundColor(bgColor)
	netDialog.form.AddButton("Cancel", nil)
	netDialog.form.AddButton("Create", nil)
	netDialog.form.SetButtonsAlign(tview.AlignRight)
	netDialog.form.SetButtonBackgroundColor(buttonBgColor)

	netDialog.setupLayout()
	netDialog.layout.SetBackgroundColor(bgColor)
	netDialog.layout.SetBorder(true)
	netDialog.layout.SetBorderColor(style.DialogBorderColor)
	netDialog.layout.SetTitle("PODMAN NETWORK CREATE")
	netDialog.layout.AddItem(netDialog.form, dialogs.DialogFormHeight, 0, true)

	netDialog.setActiveCategory(0)

	return &netDialog
}

func (d *NetworkCreateDialog) setupLayout() {
	bgColor := style.DialogBgColor

	// basic info page
	d.basicInfoPage.SetDirection(tview.FlexRow)
	d.basicInfoPage.AddItem(d.networkNameField, 1, 0, true)
	d.basicInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.basicInfoPage.AddItem(d.networkLabelsField, 1, 0, true)
	d.basicInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.basicInfoPage.AddItem(d.networkInternalCheckBox, 1, 0, true)
	d.basicInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.basicInfoPage.AddItem(d.networkDriverField, 1, 0, true)
	d.basicInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.basicInfoPage.AddItem(d.networkDriverOptionsField, 1, 0, true)
	d.basicInfoPage.SetBackgroundColor(bgColor)

	// ip settings page
	d.ipSettingsPage.SetDirection(tview.FlexRow)
	d.ipSettingsPage.AddItem(d.networkIpv6CheckBox, 1, 0, true)
	d.ipSettingsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.ipSettingsPage.AddItem(d.networkGatewayField, 1, 0, true)
	d.ipSettingsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.ipSettingsPage.AddItem(d.networkIPRangeField, 1, 0, true)
	d.ipSettingsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.ipSettingsPage.AddItem(d.networkSubnetField, 1, 0, true)
	d.ipSettingsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.ipSettingsPage.AddItem(d.networkDisableDNSCheckBox, 1, 0, true)
	d.ipSettingsPage.SetBackgroundColor(bgColor)

	// adding category pages
	d.categoryPages.AddPage(d.categoryLabels[basicInfoPageIndex], d.basicInfoPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[ipSettingsPageIndex], d.ipSettingsPage, true, true)

	// add it to layout.
	_, layoutWidth := utils.AlignStringListWidth(d.categoryLabels)
	layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	layout.AddItem(d.categories, layoutWidth+6, 0, true) //nolint:gomnd
	layout.AddItem(d.categoryPages, 0, 1, true)
	layout.SetBackgroundColor(bgColor)

	d.layout.AddItem(layout, 0, 1, true)
}

// Display displays this primitive.
func (d *NetworkCreateDialog) Display() {
	d.display = true

	d.initData()

	d.focusElement = categoryPagesFocus
}

// IsDisplay returns true if primitive is shown.
func (d *NetworkCreateDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *NetworkCreateDialog) Hide() {
	d.display = false
}

// HasFocus returns whether or not this primitive has focus.
func (d *NetworkCreateDialog) HasFocus() bool {
	if d.categories.HasFocus() || d.categoryPages.HasFocus() {
		return true
	}

	return d.Box.HasFocus() || d.form.HasFocus()
}

// Focus is called when this primitive receives focus.
func (d *NetworkCreateDialog) Focus(delegate func(p tview.Primitive)) { //nolint:cyclop
	switch d.focusElement {
	// form has focus
	case formFocus:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = categoriesFocus // category text view
				d.Focus(delegate)
				d.form.SetFocus(0)

				return nil
			}

			if event.Key() == tcell.KeyEnter {
				// d.pullSelectHandler()
				return nil
			}

			return event
		})

		delegate(d.form)

	// category text view
	case categoriesFocus:
		d.categories.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = categoryPagesFocus // category page view
				d.Focus(delegate)

				return nil
			}

			// scroll between categories
			event = utils.ParseKeyEventKey(event)
			if event.Key() == tcell.KeyDown {
				d.nextCategory()
			}

			if event.Key() == tcell.KeyUp {
				d.previousCategory()
			}

			return event
		})
		delegate(d.categories)
	// basic info page
	case networkNameFieldFocus:
		delegate(d.networkNameField)
	case networkLabelFieldFocus:
		delegate(d.networkLabelsField)
	case networkInternalCheckBoxFocus:
		delegate(d.networkInternalCheckBox)
	case networkDriverFieldFocus:
		delegate(d.networkDriverField)
	case networkDriverOptionsFieldFocus:
		delegate(d.networkDriverOptionsField)
	// ip settings page
	case networkIPv6CheckBoxFocus:
		delegate(d.networkIpv6CheckBox)
	case networkGatewatFieldFocus:
		delegate(d.networkGatewayField)
	case networkIPRangeFieldFocus:
		delegate(d.networkIPRangeField)
	case networkSubnetFieldFocus:
		delegate(d.networkSubnetField)
	case networkDisableDNSCheckBoxFocus:
		delegate(d.networkDisableDNSCheckBox)
	// category page
	case categoryPagesFocus:
		delegate(d.categoryPages)
	}
}

// InputHandler returns input handler function for this primitive.
func (d *NetworkCreateDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) { //nolint:gocognit,lll,cyclop
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("network create dialog: event %v received", event)

		if event.Key() == tcell.KeyEsc {
			d.cancelHandler()

			return
		}

		if d.basicInfoPage.HasFocus() {
			if handler := d.basicInfoPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setBasicInfoPageNextFocus()
				}

				handler(event, setFocus)

				return
			}
		}

		if d.ipSettingsPage.HasFocus() {
			if handler := d.ipSettingsPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setIPSettingsPageNextFocus()
				}

				handler(event, setFocus)

				return
			}
		}

		if d.categories.HasFocus() {
			if categroryHandler := d.categories.InputHandler(); categroryHandler != nil {
				categroryHandler(event, setFocus)

				return
			}
		}

		if d.form.HasFocus() { //nolint:nestif
			if formHandler := d.form.InputHandler(); formHandler != nil {
				if event.Key() == tcell.KeyEnter {
					enterButton := d.form.GetButton(d.form.GetButtonCount() - 1)
					if enterButton.HasFocus() {
						d.createHandler()
					}
				}

				formHandler(event, setFocus)

				return
			}
		}
	})
}

// SetRect set rects for this primitive.
func (d *NetworkCreateDialog) SetRect(x, y, width, height int) {
	if width > networkCreateDialogMaxWidth {
		emptySpace := (width - networkCreateDialogMaxWidth) / 2 //nolint:gomnd
		x += emptySpace
		width = networkCreateDialogMaxWidth
	}

	if height > networkCreateDialogHeight {
		emptySpace := (height - networkCreateDialogHeight) / 2 //nolint:gomnd
		y += emptySpace
		height = networkCreateDialogHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen.
func (d *NetworkCreateDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.Box.DrawForSubclass(screen, d)
	x, y, width, height := d.Box.GetInnerRect()

	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetCancelFunc sets form cancel button selected function.
func (d *NetworkCreateDialog) SetCancelFunc(handler func()) *NetworkCreateDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:gomnd

	cancelButton.SetSelectedFunc(handler)

	return d
}

// SetCreateFunc sets form create button selected function.
func (d *NetworkCreateDialog) SetCreateFunc(handler func()) *NetworkCreateDialog {
	d.createHandler = handler
	enterButton := d.form.GetButton(d.form.GetButtonCount() - 1)

	enterButton.SetSelectedFunc(handler)

	return d
}

func (d *NetworkCreateDialog) setActiveCategory(index int) {
	fgColor := style.DialogFgColor
	bgColor := style.ButtonBgColor
	ctgTextColor := style.GetColorHex(fgColor)
	ctgBgColor := style.GetColorHex(bgColor)

	d.activePageIndex = index

	d.categories.Clear()

	var ctgList []string

	alignedList, _ := utils.AlignStringListWidth(d.categoryLabels)

	for i := 0; i < len(d.categoryLabels); i++ {
		if i == index {
			ctgList = append(ctgList, fmt.Sprintf("[%s:%s:b]-> %s ", ctgTextColor, ctgBgColor, alignedList[i]))

			continue
		}

		ctgList = append(ctgList, fmt.Sprintf("[-:-:-]   %s ", alignedList[i]))
	}

	d.categories.SetText(strings.Join(ctgList, "\n"))

	// switch the page
	d.categoryPages.SwitchToPage(d.categoryLabels[index])
}

func (d *NetworkCreateDialog) nextCategory() {
	activePage := d.activePageIndex
	if d.activePageIndex < len(d.categoryLabels)-1 {
		activePage++
		d.setActiveCategory(activePage)

		return
	}

	d.setActiveCategory(0)
}

func (d *NetworkCreateDialog) previousCategory() {
	activePage := d.activePageIndex
	if d.activePageIndex > 0 {
		activePage--
		d.setActiveCategory(activePage)

		return
	}

	d.setActiveCategory(len(d.categoryLabels) - 1)
}

func (d *NetworkCreateDialog) initData() {
	d.setActiveCategory(0)
	d.networkNameField.SetText("")
	d.networkLabelsField.SetText("")
	d.networkInternalCheckBox.SetChecked(false)
	d.networkDriverField.SetText(networks.DefaultNetworkDriver())
	d.networkDriverOptionsField.SetText("")
	d.networkIpv6CheckBox.SetChecked(false)
	d.networkGatewayField.SetText("")
	d.networkIPRangeField.SetText("")
	d.networkSubnetField.SetText("")
	d.networkDisableDNSCheckBox.SetChecked(false)
}

func (d *NetworkCreateDialog) setBasicInfoPageNextFocus() {
	if d.networkNameField.HasFocus() {
		d.focusElement = networkLabelFieldFocus

		return
	}

	if d.networkLabelsField.HasFocus() {
		d.focusElement = networkInternalCheckBoxFocus

		return
	}

	if d.networkInternalCheckBox.HasFocus() {
		d.focusElement = networkDriverFieldFocus

		return
	}

	if d.networkDriverField.HasFocus() {
		d.focusElement = networkDriverOptionsFieldFocus

		return
	}

	d.focusElement = formFocus
}

func (d *NetworkCreateDialog) setIPSettingsPageNextFocus() {
	if d.networkIpv6CheckBox.HasFocus() {
		d.focusElement = networkGatewatFieldFocus

		return
	}

	if d.networkGatewayField.HasFocus() {
		d.focusElement = networkIPRangeFieldFocus

		return
	}

	if d.networkIPRangeField.HasFocus() {
		d.focusElement = networkSubnetFieldFocus

		return
	}

	if d.networkSubnetField.HasFocus() {
		d.focusElement = networkDisableDNSCheckBoxFocus

		return
	}

	d.focusElement = formFocus
}

// NetworkCreateOptions returns new network options.
func (d *NetworkCreateDialog) NetworkCreateOptions() networks.CreateOptions { //nolint:cyclop
	var (
		labels   = make(map[string]string)
		options  = make(map[string]string)
		subnets  []string
		gateways []string
		ipranges []string
	)

	for _, label := range strings.Split(d.networkLabelsField.GetText(), " ") {
		if label != "" {
			split := strings.Split(label, "=")
			if len(split) == 2 { //nolint:gomnd
				key := split[0]
				value := split[1]

				if key != "" && value != "" {
					labels[key] = value
				}
			}
		}
	}

	for _, option := range strings.Split(d.networkDriverOptionsField.GetText(), " ") {
		if option != "" {
			split := strings.Split(option, "=")
			if len(split) == 2 { //nolint:gomnd
				key := split[0]
				value := split[1]

				if key != "" && value != "" {
					options[key] = value
				}
			}
		}
	}

	if strings.Trim(d.networkGatewayField.GetText(), " ") != "" {
		gateways = strings.Split(d.networkGatewayField.GetText(), " ")
	}

	if strings.Trim(d.networkSubnetField.GetText(), " ") != "" {
		subnets = strings.Split(d.networkSubnetField.GetText(), " ")
	}

	if strings.Trim(d.networkIPRangeField.GetText(), " ") != "" {
		ipranges = strings.Split(d.networkIPRangeField.GetText(), " ")
	}

	opts := networks.CreateOptions{
		Name:           d.networkNameField.GetText(),
		Labels:         labels,
		Internal:       d.networkInternalCheckBox.IsChecked(),
		Drivers:        d.networkDriverField.GetText(),
		DriversOptions: options,
		IPv6:           d.networkIpv6CheckBox.IsChecked(),
		Gateways:       gateways,
		Subnets:        subnets,
		IPRanges:       ipranges,
		DisableDNS:     d.networkDisableDNSCheckBox.IsChecked(),
	}

	return opts
}
