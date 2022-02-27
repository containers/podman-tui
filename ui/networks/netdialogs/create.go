package netdialogs

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/networks"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	networkCreateDialogMaxWidth = 80
	networkCreateDialogHeight   = 19
)

const (
	formFocus = 0 + iota
	categoriesFocus
	categoryPagesFocus
	networkNameFieldFocus
	networkLabelFieldFocus
	networkInternalCheckBoxFocus
	networkMacvlanFieldFocus
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

// NetworkCreateDialog implements network create dialog
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
	networkMacvlanField       *tview.InputField
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

// NewNetworkCreateDialog returns new network create dialog primitive NetworkCreateDialog
func NewNetworkCreateDialog() *NetworkCreateDialog {
	podDialog := NetworkCreateDialog{
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
		networkMacvlanField:       tview.NewInputField(),
		networkDriverField:        tview.NewInputField(),
		networkDriverOptionsField: tview.NewInputField(),
		networkIpv6CheckBox:       tview.NewCheckbox(),
		networkGatewayField:       tview.NewInputField(),
		networkIPRangeField:       tview.NewInputField(),
		networkSubnetField:        tview.NewInputField(),
		networkDisableDNSCheckBox: tview.NewCheckbox(),
	}

	bgColor := utils.Styles.ImageHistoryDialog.BgColor

	podDialog.categories.SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)
	podDialog.categories.SetBackgroundColor(bgColor)
	podDialog.categories.SetBorder(true)

	// basic information setup page
	basicInfoPageLabelWidth := 12
	// name field
	podDialog.networkNameField.SetLabel("name:")
	podDialog.networkNameField.SetLabelWidth(basicInfoPageLabelWidth)
	podDialog.networkNameField.SetBackgroundColor(bgColor)
	podDialog.networkNameField.SetLabelColor(tcell.ColorWhite)
	// labels field
	podDialog.networkLabelsField.SetLabel("labels:")
	podDialog.networkLabelsField.SetLabelWidth(basicInfoPageLabelWidth)
	podDialog.networkLabelsField.SetBackgroundColor(bgColor)
	podDialog.networkLabelsField.SetLabelColor(tcell.ColorWhite)
	// internal check box
	podDialog.networkInternalCheckBox.SetLabel("internal")
	podDialog.networkInternalCheckBox.SetLabelWidth(basicInfoPageLabelWidth)
	podDialog.networkInternalCheckBox.SetChecked(false)
	podDialog.networkInternalCheckBox.SetBackgroundColor(bgColor)
	podDialog.networkInternalCheckBox.SetLabelColor(tcell.ColorWhite)
	// macvlan
	podDialog.networkMacvlanField.SetLabel("macvlan:")
	podDialog.networkMacvlanField.SetLabelWidth(basicInfoPageLabelWidth)
	podDialog.networkMacvlanField.SetBackgroundColor(bgColor)
	podDialog.networkMacvlanField.SetLabelColor(tcell.ColorWhite)
	// drivers
	podDialog.networkDriverField.SetLabel("drivers:")
	podDialog.networkDriverField.SetLabelWidth(basicInfoPageLabelWidth)
	podDialog.networkDriverField.SetBackgroundColor(bgColor)
	podDialog.networkDriverField.SetLabelColor(tcell.ColorWhite)
	// drivers options
	podDialog.networkDriverOptionsField.SetLabel("options:")
	podDialog.networkDriverOptionsField.SetLabelWidth(basicInfoPageLabelWidth)
	podDialog.networkDriverOptionsField.SetBackgroundColor(bgColor)
	podDialog.networkDriverOptionsField.SetLabelColor(tcell.ColorWhite)

	// ip settings page
	ipSettingsPageLabelWidth := 12
	// ipv6 check box
	podDialog.networkIpv6CheckBox.SetLabel("IPv6")
	podDialog.networkIpv6CheckBox.SetLabelWidth(ipSettingsPageLabelWidth)
	podDialog.networkIpv6CheckBox.SetChecked(false)
	podDialog.networkIpv6CheckBox.SetBackgroundColor(bgColor)
	podDialog.networkIpv6CheckBox.SetLabelColor(tcell.ColorWhite)

	// gateway
	podDialog.networkGatewayField.SetLabel("gateway:")
	podDialog.networkGatewayField.SetLabelWidth(basicInfoPageLabelWidth)
	podDialog.networkGatewayField.SetBackgroundColor(bgColor)
	podDialog.networkGatewayField.SetLabelColor(tcell.ColorWhite)

	// ip range
	podDialog.networkIPRangeField.SetLabel("IP range:")
	podDialog.networkIPRangeField.SetLabelWidth(basicInfoPageLabelWidth)
	podDialog.networkIPRangeField.SetBackgroundColor(bgColor)
	podDialog.networkIPRangeField.SetLabelColor(tcell.ColorWhite)

	// subnet
	podDialog.networkSubnetField.SetLabel("subnet:")
	podDialog.networkSubnetField.SetLabelWidth(basicInfoPageLabelWidth)
	podDialog.networkSubnetField.SetBackgroundColor(bgColor)
	podDialog.networkSubnetField.SetLabelColor(tcell.ColorWhite)
	// dns check box
	podDialog.networkDisableDNSCheckBox.SetLabel("disable DNS")
	podDialog.networkDisableDNSCheckBox.SetLabelWidth(basicInfoPageLabelWidth)
	podDialog.networkDisableDNSCheckBox.SetChecked(false)
	podDialog.networkDisableDNSCheckBox.SetBackgroundColor(bgColor)
	podDialog.networkDisableDNSCheckBox.SetLabelColor(tcell.ColorWhite)

	// category pages
	podDialog.categoryPages.SetBackgroundColor(bgColor)
	podDialog.categoryPages.SetBorder(true)

	// form
	podDialog.form.SetBackgroundColor(bgColor)
	podDialog.form.AddButton("Cancel", nil)
	podDialog.form.AddButton("Create", nil)
	podDialog.form.SetButtonsAlign(tview.AlignRight)

	podDialog.layout.AddItem(tview.NewBox().SetBackgroundColor(bgColor), 1, 0, true)
	podDialog.setupLayout()
	podDialog.layout.SetBackgroundColor(bgColor)
	podDialog.layout.SetBorder(true)
	podDialog.layout.SetTitle("PODMAN NETWORK CREATE")
	podDialog.layout.AddItem(podDialog.form, dialogs.DialogFormHeight, 0, true)

	podDialog.setActiveCategory(0)
	return &podDialog
}

func (d *NetworkCreateDialog) setupLayout() {
	bgColor := utils.Styles.ImageHistoryDialog.BgColor

	emptySpace := func() *tview.Box {
		box := tview.NewBox()
		box.SetBackgroundColor(bgColor)
		return box
	}

	// basic info page
	d.basicInfoPage.SetDirection(tview.FlexRow)
	d.basicInfoPage.AddItem(d.networkNameField, 1, 0, true)
	d.basicInfoPage.AddItem(emptySpace(), 1, 0, true)
	d.basicInfoPage.AddItem(d.networkLabelsField, 1, 0, true)
	d.basicInfoPage.AddItem(emptySpace(), 1, 0, true)
	d.basicInfoPage.AddItem(d.networkInternalCheckBox, 1, 0, true)
	d.basicInfoPage.AddItem(emptySpace(), 1, 0, true)
	d.basicInfoPage.AddItem(d.networkMacvlanField, 1, 0, true)
	d.basicInfoPage.AddItem(emptySpace(), 1, 0, true)
	d.basicInfoPage.AddItem(d.networkDriverField, 1, 0, true)
	d.basicInfoPage.AddItem(emptySpace(), 1, 0, true)
	d.basicInfoPage.AddItem(d.networkDriverOptionsField, 1, 0, true)
	d.basicInfoPage.SetBackgroundColor(bgColor)

	// ip settings page
	d.ipSettingsPage.SetDirection(tview.FlexRow)
	d.ipSettingsPage.AddItem(d.networkIpv6CheckBox, 1, 0, true)
	d.ipSettingsPage.AddItem(emptySpace(), 1, 0, true)
	d.ipSettingsPage.AddItem(d.networkGatewayField, 1, 0, true)
	d.ipSettingsPage.AddItem(emptySpace(), 1, 0, true)
	d.ipSettingsPage.AddItem(d.networkIPRangeField, 1, 0, true)
	d.ipSettingsPage.AddItem(emptySpace(), 1, 0, true)
	d.ipSettingsPage.AddItem(d.networkSubnetField, 1, 0, true)
	d.ipSettingsPage.AddItem(emptySpace(), 1, 0, true)
	d.ipSettingsPage.AddItem(d.networkDisableDNSCheckBox, 1, 0, true)
	d.ipSettingsPage.SetBackgroundColor(bgColor)

	// adding category pages
	d.categoryPages.AddPage(d.categoryLabels[basicInfoPageIndex], d.basicInfoPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[ipSettingsPageIndex], d.ipSettingsPage, true, true)

	// add it to layout.
	_, layoutWidth := utils.AlignStringListWidth(d.categoryLabels)
	layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	layout.AddItem(d.categories, layoutWidth+6, 0, true)
	layout.AddItem(d.categoryPages, 0, 1, true)
	layout.SetBackgroundColor(bgColor)

	d.layout.AddItem(layout, 0, 1, true)

}

// Display displays this primitive
func (d *NetworkCreateDialog) Display() {
	d.display = true
	d.initData()
	d.focusElement = categoryPagesFocus
}

// IsDisplay returns true if primitive is shown
func (d *NetworkCreateDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive
func (d *NetworkCreateDialog) Hide() {
	d.display = false
}

// HasFocus returns whether or not this primitive has focus
func (d *NetworkCreateDialog) HasFocus() bool {
	if d.categories.HasFocus() || d.categoryPages.HasFocus() {
		return true
	}

	return d.Box.HasFocus() || d.form.HasFocus()
}

// Focus is called when this primitive receives focus
func (d *NetworkCreateDialog) Focus(delegate func(p tview.Primitive)) {
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
				//d.pullSelectHandler()
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
	case networkMacvlanFieldFocus:
		delegate(d.networkMacvlanField)
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

// InputHandler returns input handler function for this primitive
func (d *NetworkCreateDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("network create dialog: event %v received", event.Key())
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
		if d.form.HasFocus() {
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
		emptySpace := (width - networkCreateDialogMaxWidth) / 2
		x = x + emptySpace
		width = networkCreateDialogMaxWidth
	}

	if height > networkCreateDialogHeight {
		emptySpace := (height - networkCreateDialogHeight) / 2
		y = y + emptySpace
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

// SetCancelFunc sets form cancel button selected function
func (d *NetworkCreateDialog) SetCancelFunc(handler func()) *NetworkCreateDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2)
	cancelButton.SetSelectedFunc(handler)
	return d
}

// SetCreateFunc sets form create button selected function
func (d *NetworkCreateDialog) SetCreateFunc(handler func()) *NetworkCreateDialog {
	d.createHandler = handler
	enterButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	enterButton.SetSelectedFunc(handler)
	return d
}

func (d *NetworkCreateDialog) setActiveCategory(index int) {
	d.activePageIndex = index
	d.categories.Clear()
	var ctgList []string
	alignedList, _ := utils.AlignStringListWidth(d.categoryLabels)
	for i := 0; i < len(d.categoryLabels); i++ {
		if i == index {
			ctgList = append(ctgList, fmt.Sprintf("[white:blue:b]-> %s ", alignedList[i]))
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
		activePage = activePage + 1
		d.setActiveCategory(activePage)
		return
	}
	d.setActiveCategory(0)
}

func (d *NetworkCreateDialog) previousCategory() {
	activePage := d.activePageIndex
	if d.activePageIndex > 0 {
		activePage = activePage - 1
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
	d.networkMacvlanField.SetText("")
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
	} else if d.networkLabelsField.HasFocus() {
		d.focusElement = networkInternalCheckBoxFocus
	} else if d.networkInternalCheckBox.HasFocus() {
		d.focusElement = networkMacvlanFieldFocus
	} else if d.networkMacvlanField.HasFocus() {
		d.focusElement = networkDriverFieldFocus
	} else if d.networkDriverField.HasFocus() {
		d.focusElement = networkDriverOptionsFieldFocus
	} else {
		d.focusElement = formFocus
	}
}

func (d *NetworkCreateDialog) setIPSettingsPageNextFocus() {
	if d.networkIpv6CheckBox.HasFocus() {
		d.focusElement = networkGatewatFieldFocus
	} else if d.networkGatewayField.HasFocus() {
		d.focusElement = networkIPRangeFieldFocus
	} else if d.networkIPRangeField.HasFocus() {
		d.focusElement = networkSubnetFieldFocus
	} else if d.networkSubnetField.HasFocus() {
		d.focusElement = networkDisableDNSCheckBoxFocus
	} else {
		d.focusElement = formFocus
	}
}

// NetworkCreateOptions returns new network options
func (d *NetworkCreateDialog) NetworkCreateOptions() networks.CreateOptions {
	var (
		labels  = make(map[string]string)
		options = make(map[string]string)
	)
	for _, label := range strings.Split(d.networkLabelsField.GetText(), " ") {
		if label != "" {
			split := strings.Split(label, "=")
			if len(split) == 2 {
				key := split[0]
				value := split[1]
				if key != "" && value != "" {
					labels[key] = value
				}
			}
		}
	}
	for _, option := range strings.Split(d.networkLabelsField.GetText(), " ") {
		if option != "" {
			split := strings.Split(option, "=")
			if len(split) == 2 {
				key := split[0]
				value := split[1]
				if key != "" && value != "" {
					options[key] = value
				}
			}
		}
	}
	opts := networks.CreateOptions{
		Name:           d.networkNameField.GetText(),
		Labels:         labels,
		Internal:       d.networkInternalCheckBox.IsChecked(),
		Macvlan:        d.networkMacvlanField.GetText(),
		Drivers:        d.networkDriverField.GetText(),
		DriversOptions: options,
		IPv6:           d.networkIpv6CheckBox.IsChecked(),
		Subnet:         d.networkSubnetField.GetText(),
		IPRange:        d.networkIPRangeField.GetText(),
		DisableDNS:     d.networkDisableDNSCheckBox.IsChecked(),
	}
	return opts
}
