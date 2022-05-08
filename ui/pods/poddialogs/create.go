package poddialogs

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/networks"
	"github.com/containers/podman-tui/pdcs/pods"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	podCreateDialogMaxWidth = 80
	podCreateDialogHeight   = 19
)

const (
	podFormFocus = 0 + iota
	categoriesFocus
	categoryPagesFocus
	podNameFieldFocus
	podNoHostsCheckBoxFocus
	podLabelsFieldFocus
	podSelinuxLabelFieldFocus
	podApparmorFieldFocus
	podSeccompFieldFocus
	podMaskFieldFocus
	podUnmaskFieldFocus
	podNoNewPrivFieldFocus
	podDNSServerFieldFocus
	podDNSOptionsFieldFocus
	podDNSSearchDomaindFieldFocus
	podInfraCheckBoxFocus
	podInfraCommandFieldFocus
	podInfraImageFieldFocus
	podHostnameFieldFocus
	podIPAddressFieldFocus
	podMacAddressFieldFocus
	podHostToIPMapFieldFocus
	podNetworkFieldFocus
)

const (
	basicInfoPageIndex = 0 + iota
	dnsSetupPageIndex
	infraSetupPageIndex
	networkingPageIndex
	securityOptsPageIndex
)

// PodCreateDialog implements pod create dialog
type PodCreateDialog struct {
	*tview.Box
	layout                   *tview.Flex
	categoryLabels           []string
	categories               *tview.TextView
	categoryPages            *tview.Pages
	basicInfoPage            *tview.Flex
	securityOptsPage         *tview.Flex
	dnsSetupPage             *tview.Flex
	infraSetupPage           *tview.Flex
	networkingPage           *tview.Flex
	form                     *tview.Form
	display                  bool
	activePageIndex          int
	focusElement             int
	defaultInfraImage        string
	podNameField             *tview.InputField
	podNoHostsCheckBox       *tview.Checkbox
	podLabelsField           *tview.InputField
	podSelinuxLabelField     *tview.InputField
	podApparmorField         *tview.InputField
	podSeccompField          *tview.InputField
	podMaskField             *tview.InputField
	podUnmaskField           *tview.InputField
	podNoNewPrivField        *tview.Checkbox
	podDNSServerField        *tview.InputField
	podDNSOptionsField       *tview.InputField
	podDNSSearchDomaindField *tview.InputField
	podInfraCheckBox         *tview.Checkbox
	podInfraCommandField     *tview.InputField
	podInfraImageField       *tview.InputField
	podHostnameField         *tview.InputField
	podIPAddressField        *tview.InputField
	podMacAddressField       *tview.InputField
	podHostToIPMapField      *tview.InputField
	podNetworkField          *tview.DropDown
	//podNetworkAliasesField   *tview.InputField
	cancelHandler func()
	createHandler func()
}

// NewPodCreateDialog returns new pod create dialog primitive PodCreateDialog
func NewPodCreateDialog() *PodCreateDialog {
	podDialog := PodCreateDialog{
		Box:              tview.NewBox(),
		layout:           tview.NewFlex().SetDirection(tview.FlexRow),
		categories:       tview.NewTextView(),
		categoryPages:    tview.NewPages(),
		basicInfoPage:    tview.NewFlex(),
		securityOptsPage: tview.NewFlex(),
		dnsSetupPage:     tview.NewFlex(),
		infraSetupPage:   tview.NewFlex(),
		networkingPage:   tview.NewFlex(),
		form:             tview.NewForm(),
		categoryLabels: []string{
			"Basic Information",
			"DNS Setup",
			"Infra Setup",
			"Networking",
			"Security Options"},
		activePageIndex:          0,
		display:                  false,
		defaultInfraImage:        pods.DefaultPodInfraImage(),
		podNameField:             tview.NewInputField(),
		podNoHostsCheckBox:       tview.NewCheckbox(),
		podLabelsField:           tview.NewInputField(),
		podSelinuxLabelField:     tview.NewInputField(),
		podApparmorField:         tview.NewInputField(),
		podSeccompField:          tview.NewInputField(),
		podMaskField:             tview.NewInputField(),
		podUnmaskField:           tview.NewInputField(),
		podNoNewPrivField:        tview.NewCheckbox(),
		podDNSServerField:        tview.NewInputField(),
		podDNSOptionsField:       tview.NewInputField(),
		podDNSSearchDomaindField: tview.NewInputField(),
		podInfraCheckBox:         tview.NewCheckbox(),
		podInfraCommandField:     tview.NewInputField(),
		podInfraImageField:       tview.NewInputField(),
		podHostnameField:         tview.NewInputField(),
		podIPAddressField:        tview.NewInputField(),
		podMacAddressField:       tview.NewInputField(),
		podHostToIPMapField:      tview.NewInputField(),
		podNetworkField:          tview.NewDropDown(),
		//podNetworkAliasesField:   tview.NewInputField(),
	}

	bgColor := utils.Styles.PodCreateDialog.BgColor
	fgColor := utils.Styles.PodCreateDialog.FgColor
	inputFieldBgColor := utils.Styles.InputFieldPrimitive.BgColor

	podDialog.categories.SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)
	podDialog.categories.SetBackgroundColor(bgColor)
	podDialog.categories.SetBorder(true)

	// basic information setup page
	basicInfoPageLabelWidth := 12
	// name field
	podDialog.podNameField.SetLabel("name:")
	podDialog.podNameField.SetLabelWidth(basicInfoPageLabelWidth)
	podDialog.podNameField.SetBackgroundColor(bgColor)
	podDialog.podNameField.SetLabelColor(fgColor)
	podDialog.podNameField.SetFieldBackgroundColor(inputFieldBgColor)

	// no hosts check box
	podDialog.podNoHostsCheckBox.SetLabel("no hosts")
	podDialog.podNoHostsCheckBox.SetLabelWidth(basicInfoPageLabelWidth)
	podDialog.podNoHostsCheckBox.SetChecked(false)
	podDialog.podNoHostsCheckBox.SetBackgroundColor(bgColor)
	podDialog.podNoHostsCheckBox.SetLabelColor(fgColor)
	podDialog.podNoHostsCheckBox.SetFieldBackgroundColor(inputFieldBgColor)

	// labels field
	podDialog.podLabelsField.SetLabel("labels:")
	podDialog.podLabelsField.SetLabelWidth(basicInfoPageLabelWidth)
	podDialog.podLabelsField.SetBackgroundColor(bgColor)
	podDialog.podLabelsField.SetLabelColor(fgColor)
	podDialog.podLabelsField.SetFieldBackgroundColor(inputFieldBgColor)

	// security options
	securityOptsPageLabelWidth := 10
	// labels
	podDialog.podSelinuxLabelField.SetLabel("Label:")
	podDialog.podSelinuxLabelField.SetLabelWidth(securityOptsPageLabelWidth)
	podDialog.podSelinuxLabelField.SetBackgroundColor(bgColor)
	podDialog.podSelinuxLabelField.SetLabelColor(fgColor)
	podDialog.podSelinuxLabelField.SetFieldBackgroundColor(inputFieldBgColor)

	// apparmor
	podDialog.podApparmorField.SetLabel("Apparmor:")
	podDialog.podApparmorField.SetLabelWidth(securityOptsPageLabelWidth)
	podDialog.podApparmorField.SetBackgroundColor(bgColor)
	podDialog.podApparmorField.SetLabelColor(fgColor)
	podDialog.podApparmorField.SetFieldBackgroundColor(inputFieldBgColor)

	// seccomp
	podDialog.podSeccompField.SetLabel("Seccomp:")
	podDialog.podSeccompField.SetLabelWidth(securityOptsPageLabelWidth)
	podDialog.podSeccompField.SetBackgroundColor(bgColor)
	podDialog.podSeccompField.SetLabelColor(fgColor)
	podDialog.podSeccompField.SetFieldBackgroundColor(inputFieldBgColor)

	// mask
	podDialog.podMaskField.SetLabel("Mask:")
	podDialog.podMaskField.SetLabelWidth(securityOptsPageLabelWidth)
	podDialog.podMaskField.SetBackgroundColor(bgColor)
	podDialog.podMaskField.SetLabelColor(fgColor)
	podDialog.podMaskField.SetFieldBackgroundColor(inputFieldBgColor)

	// unmask
	podDialog.podUnmaskField.SetLabel("Unmask:")
	podDialog.podUnmaskField.SetLabelWidth(securityOptsPageLabelWidth)
	podDialog.podUnmaskField.SetBackgroundColor(bgColor)
	podDialog.podUnmaskField.SetLabelColor(fgColor)
	podDialog.podUnmaskField.SetFieldBackgroundColor(inputFieldBgColor)

	// no new privileges
	podDialog.podNoNewPrivField.SetLabel("No new privileges ")
	podDialog.podNoNewPrivField.SetBackgroundColor(bgColor)
	podDialog.podNoNewPrivField.SetLabelColor(tcell.ColorWhite)

	// DNS setup page
	dnsPageLabelWidth := 16
	// DNS server field
	podDialog.podDNSServerField.SetLabel("DNS servers:")
	podDialog.podDNSServerField.SetLabelWidth(dnsPageLabelWidth)
	podDialog.podDNSServerField.SetBackgroundColor(bgColor)
	podDialog.podDNSServerField.SetLabelColor(fgColor)
	podDialog.podDNSServerField.SetFieldBackgroundColor(inputFieldBgColor)

	// DNS options field
	podDialog.podDNSOptionsField.SetLabel("DNS options:")
	podDialog.podDNSOptionsField.SetLabelWidth(dnsPageLabelWidth)
	podDialog.podDNSOptionsField.SetBackgroundColor(bgColor)
	podDialog.podDNSOptionsField.SetLabelColor(fgColor)
	podDialog.podDNSOptionsField.SetFieldBackgroundColor(inputFieldBgColor)

	// DNS search domains field
	podDialog.podDNSSearchDomaindField.SetLabel("search domains:")
	podDialog.podDNSSearchDomaindField.SetLabelWidth(dnsPageLabelWidth)
	podDialog.podDNSSearchDomaindField.SetBackgroundColor(bgColor)
	podDialog.podDNSSearchDomaindField.SetLabelColor(fgColor)
	podDialog.podDNSSearchDomaindField.SetFieldBackgroundColor(inputFieldBgColor)

	// infra page
	infraPageLabelWidth := 15
	// infra check box
	podDialog.podInfraCheckBox.SetLabel("infra")
	podDialog.podInfraCheckBox.SetLabelWidth(infraPageLabelWidth)
	podDialog.podInfraCheckBox.SetChecked(true)
	podDialog.podInfraCheckBox.SetBackgroundColor(bgColor)
	podDialog.podInfraCheckBox.SetLabelColor(fgColor)
	podDialog.podInfraCheckBox.SetFieldBackgroundColor(inputFieldBgColor)

	// infra command field
	podDialog.podInfraCommandField.SetLabel("infra command:")
	podDialog.podInfraCommandField.SetLabelWidth(infraPageLabelWidth)
	podDialog.podInfraCommandField.SetBackgroundColor(bgColor)
	podDialog.podInfraCommandField.SetLabelColor(fgColor)
	podDialog.podInfraCommandField.SetFieldBackgroundColor(inputFieldBgColor)

	// infra image field
	podDialog.podInfraImageField.SetLabel("infra image:")
	podDialog.podInfraImageField.SetText(podDialog.defaultInfraImage)
	podDialog.podInfraImageField.SetLabelWidth(infraPageLabelWidth)
	podDialog.podInfraImageField.SetBackgroundColor(bgColor)
	podDialog.podInfraImageField.SetLabelColor(fgColor)
	podDialog.podInfraImageField.SetFieldBackgroundColor(inputFieldBgColor)

	// networking page
	networkingLabelWidth := 17
	// hostname field
	podDialog.podHostnameField.SetLabel("hostname:")
	podDialog.podHostnameField.SetLabelWidth(networkingLabelWidth)
	podDialog.podHostnameField.SetBackgroundColor(bgColor)
	podDialog.podHostnameField.SetLabelColor(fgColor)
	podDialog.podHostnameField.SetFieldBackgroundColor(inputFieldBgColor)

	// ip address field
	podDialog.podIPAddressField.SetLabel("ip address:")
	podDialog.podIPAddressField.SetLabelWidth(networkingLabelWidth)
	podDialog.podIPAddressField.SetBackgroundColor(bgColor)
	podDialog.podIPAddressField.SetLabelColor(fgColor)
	podDialog.podIPAddressField.SetFieldBackgroundColor(inputFieldBgColor)

	// mac address field
	podDialog.podMacAddressField.SetLabel("mac address:")
	podDialog.podMacAddressField.SetLabelWidth(networkingLabelWidth)
	podDialog.podMacAddressField.SetBackgroundColor(bgColor)
	podDialog.podMacAddressField.SetLabelColor(fgColor)
	podDialog.podMacAddressField.SetFieldBackgroundColor(inputFieldBgColor)

	// host-to-ip map field
	podDialog.podHostToIPMapField.SetLabel("host-to-ip:")
	podDialog.podHostToIPMapField.SetLabelWidth(networkingLabelWidth)
	podDialog.podHostToIPMapField.SetBackgroundColor(bgColor)
	podDialog.podHostToIPMapField.SetLabelColor(fgColor)
	podDialog.podHostToIPMapField.SetFieldBackgroundColor(inputFieldBgColor)

	// network field
	podDialog.podNetworkField.SetLabel("network:")
	podDialog.podNetworkField.SetLabelWidth(networkingLabelWidth)
	podDialog.podNetworkField.SetBackgroundColor(bgColor)
	podDialog.podNetworkField.SetLabelColor(fgColor)
	ddUnselectedStyle := utils.Styles.DropdownStyle.Unselected
	ddselectedStyle := utils.Styles.DropdownStyle.Selected
	podDialog.podNetworkField.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	podDialog.podNetworkField.SetFieldBackgroundColor(inputFieldBgColor)

	// category pages
	podDialog.categoryPages.SetBackgroundColor(bgColor)
	podDialog.categoryPages.SetBorder(true)

	// form
	podDialog.form.SetBackgroundColor(bgColor)
	podDialog.form.AddButton("Cancel", nil)
	podDialog.form.AddButton("Create", nil)
	podDialog.form.SetButtonsAlign(tview.AlignRight)
	podDialog.form.SetButtonBackgroundColor(utils.Styles.ButtonPrimitive.BgColor)

	podDialog.layout.AddItem(tview.NewBox().SetBackgroundColor(bgColor), 1, 0, true)
	podDialog.setupLayout()
	podDialog.layout.SetBackgroundColor(bgColor)
	podDialog.layout.SetBorder(true)
	podDialog.layout.SetTitle("PODMAN POD CREATE")
	podDialog.layout.AddItem(podDialog.form, dialogs.DialogFormHeight, 0, true)

	podDialog.setActiveCategory(0)

	podDialog.initCustomInputHanlers()
	return &podDialog
}

func (d *PodCreateDialog) setupLayout() {
	bgColor := utils.Styles.PodCreateDialog.BgColor

	// basic info page
	d.basicInfoPage.SetDirection(tview.FlexRow)
	d.basicInfoPage.AddItem(d.podNameField, 1, 0, true)
	d.basicInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.basicInfoPage.AddItem(d.podNoHostsCheckBox, 1, 0, true)
	d.basicInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.basicInfoPage.AddItem(d.podLabelsField, 1, 0, true)
	d.basicInfoPage.SetBackgroundColor(bgColor)

	// security options page
	d.securityOptsPage.SetDirection(tview.FlexRow)
	d.securityOptsPage.AddItem(d.podSelinuxLabelField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.podApparmorField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.podSeccompField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.podMaskField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.podUnmaskField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.podNoNewPrivField, 1, 0, true)
	d.securityOptsPage.SetBackgroundColor(bgColor)

	// DNS setup page
	d.dnsSetupPage.SetDirection(tview.FlexRow)
	d.dnsSetupPage.AddItem(d.podDNSServerField, 1, 0, true)
	d.dnsSetupPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.dnsSetupPage.AddItem(d.podDNSOptionsField, 1, 0, true)
	d.dnsSetupPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.dnsSetupPage.AddItem(d.podDNSSearchDomaindField, 1, 0, true)
	d.dnsSetupPage.SetBackgroundColor(bgColor)

	// infra page
	d.infraSetupPage.SetDirection(tview.FlexRow)
	d.infraSetupPage.AddItem(d.podInfraCheckBox, 1, 0, true)
	d.infraSetupPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.infraSetupPage.AddItem(d.podInfraCommandField, 1, 0, true)
	d.infraSetupPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.infraSetupPage.AddItem(d.podInfraImageField, 1, 0, true)
	d.infraSetupPage.SetBackgroundColor(bgColor)

	// networking page
	d.networkingPage.SetDirection(tview.FlexRow)
	d.networkingPage.AddItem(d.podHostnameField, 1, 0, true)
	d.networkingPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.networkingPage.AddItem(d.podIPAddressField, 1, 0, true)
	d.networkingPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.networkingPage.AddItem(d.podMacAddressField, 1, 0, true)
	d.networkingPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.networkingPage.AddItem(d.podHostToIPMapField, 1, 0, true)
	d.networkingPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.networkingPage.AddItem(d.podNetworkField, 1, 0, true)
	d.networkingPage.SetBackgroundColor(bgColor)

	// adding category pages
	d.categoryPages.AddPage(d.categoryLabels[basicInfoPageIndex], d.basicInfoPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[dnsSetupPageIndex], d.dnsSetupPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[infraSetupPageIndex], d.infraSetupPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[networkingPageIndex], d.networkingPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[securityOptsPageIndex], d.securityOptsPage, true, true)

	// add it to layout.
	_, layoutWidth := utils.AlignStringListWidth(d.categoryLabels)
	layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	layout.AddItem(d.categories, layoutWidth+6, 0, true)
	layout.AddItem(d.categoryPages, 0, 1, true)
	layout.SetBackgroundColor(bgColor)

	d.layout.AddItem(layout, 0, 1, true)

}

// Display displays this primitive
func (d *PodCreateDialog) Display() {
	d.display = true
	d.initData()
	d.focusElement = categoryPagesFocus
}

// IsDisplay returns true if primitive is shown
func (d *PodCreateDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive
func (d *PodCreateDialog) Hide() {
	d.display = false
}

// HasFocus returns whether or not this primitive has focus
func (d *PodCreateDialog) HasFocus() bool {
	if d.categories.HasFocus() || d.categoryPages.HasFocus() {
		return true
	}

	return d.Box.HasFocus() || d.form.HasFocus()
}

// Focus is called when this primitive receives focus
func (d *PodCreateDialog) Focus(delegate func(p tview.Primitive)) {
	switch d.focusElement {
	// form has focus
	case podFormFocus:
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
			// scroll between categories
			event = utils.ParseKeyEventKey(event)
			if event.Key() == tcell.KeyDown {
				d.nextCategory()
			}
			if event.Key() == tcell.KeyUp {
				d.previousCategory()
			}
			return nil
		})
		delegate(d.categories)
	// basic info page
	case podNoHostsCheckBoxFocus:
		delegate(d.podNoHostsCheckBox)
	case podLabelsFieldFocus:
		delegate(d.podLabelsField)
	// security options page
	case podSelinuxLabelFieldFocus:
		delegate(d.podSelinuxLabelField)
	case podApparmorFieldFocus:
		delegate(d.podApparmorField)
	case podSeccompFieldFocus:
		delegate(d.podSeccompField)
	case podMaskFieldFocus:
		delegate(d.podMaskField)
	case podUnmaskFieldFocus:
		delegate(d.podUnmaskField)
	case podNoNewPrivFieldFocus:
		delegate(d.podNoNewPrivField)
	// dns page
	case podDNSOptionsFieldFocus:
		delegate(d.podDNSOptionsField)
	case podDNSSearchDomaindFieldFocus:
		delegate(d.podDNSSearchDomaindField)
	// infra page
	case podInfraCommandFieldFocus:
		delegate(d.podInfraCommandField)
	case podInfraImageFieldFocus:
		delegate(d.podInfraImageField)
	// networking page
	case podIPAddressFieldFocus:
		delegate(d.podIPAddressField)
	case podMacAddressFieldFocus:
		delegate(d.podMacAddressField)
	case podHostToIPMapFieldFocus:
		delegate(d.podHostToIPMapField)
	case podNetworkFieldFocus:
		delegate(d.podNetworkField)
	//case podNetworkAliasesFieldFocus:
	//	delegate(d.podNetworkAliasesField)
	// category page
	case categoryPagesFocus:
		delegate(d.categoryPages)
	}

}

func (d *PodCreateDialog) initCustomInputHanlers() {
	// newtwork dropdown
	d.podNetworkField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = utils.ParseKeyEventKey(event)
		return event
	})
}

// InputHandler returns input handler function for this primitive
func (d *PodCreateDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("pod create dialog: event %v received", event)
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
		if d.dnsSetupPage.HasFocus() {
			if handler := d.dnsSetupPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setDNSSetupPageNextFocus()
				}
				handler(event, setFocus)
				return
			}
		}
		if d.infraSetupPage.HasFocus() {
			if handler := d.infraSetupPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setInfraSetupPageNextFocus()
				}
				handler(event, setFocus)
				return
			}
		}
		if d.networkingPage.HasFocus() {
			if handler := d.networkingPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setNetworkingPageNextFocus()
				}
				handler(event, setFocus)
				return
			}
		}
		if d.securityOptsPage.HasFocus() {
			if handler := d.securityOptsPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setSecurityOptsPageNextFocus()
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
func (d *PodCreateDialog) SetRect(x, y, width, height int) {

	if width > podCreateDialogMaxWidth {
		emptySpace := (width - podCreateDialogMaxWidth) / 2
		x = x + emptySpace
		width = podCreateDialogMaxWidth
	}

	if height > podCreateDialogHeight {
		emptySpace := (height - podCreateDialogHeight) / 2
		y = y + emptySpace
		height = podCreateDialogHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen.
func (d *PodCreateDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}
	d.Box.DrawForSubclass(screen, d)
	x, y, width, height := d.Box.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetCancelFunc sets form cancel button selected function
func (d *PodCreateDialog) SetCancelFunc(handler func()) *PodCreateDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2)
	cancelButton.SetSelectedFunc(handler)
	return d
}

// SetCreateFunc sets form create button selected function
func (d *PodCreateDialog) SetCreateFunc(handler func()) *PodCreateDialog {
	d.createHandler = handler
	enterButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	enterButton.SetSelectedFunc(handler)
	return d
}

func (d *PodCreateDialog) setActiveCategory(index int) {
	fgColor := utils.Styles.PodCreateDialog.FgColor
	bgColor := utils.Styles.ButtonPrimitive.BgColor
	ctgTextColor := utils.GetColorName(fgColor)
	ctgBgColor := utils.GetColorName(bgColor)

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

func (d *PodCreateDialog) nextCategory() {
	activePage := d.activePageIndex
	if d.activePageIndex < len(d.categoryLabels)-1 {
		activePage = activePage + 1
		d.setActiveCategory(activePage)
		return
	}
	d.setActiveCategory(0)
}

func (d *PodCreateDialog) previousCategory() {
	activePage := d.activePageIndex
	if d.activePageIndex > 0 {
		activePage = activePage - 1
		d.setActiveCategory(activePage)
		return
	}
	d.setActiveCategory(len(d.categoryLabels) - 1)
}

func (d *PodCreateDialog) initData() {
	// get available networks
	networkOptions := []string{""}
	networkList, _ := networks.List()
	for i := 0; i < len(networkList); i++ {
		networkOptions = append(networkOptions, networkList[i][1])
	}

	d.setActiveCategory(0)
	d.podNameField.SetText("")
	d.podNoHostsCheckBox.SetChecked(false)
	d.podLabelsField.SetText("")

	d.podSelinuxLabelField.SetText("")
	d.podApparmorField.SetText("")
	d.podSeccompField.SetText("")
	d.podMaskField.SetText("")
	d.podUnmaskField.SetText("")
	d.podNoNewPrivField.SetChecked(false)

	d.podDNSServerField.SetText("")
	d.podDNSOptionsField.SetText("")
	d.podDNSSearchDomaindField.SetText("")

	d.podInfraCheckBox.SetChecked(true)
	d.podInfraCommandField.SetText("")
	d.podInfraImageField.SetText(d.defaultInfraImage)

	d.podHostnameField.SetText("")
	d.podIPAddressField.SetText("")
	d.podMacAddressField.SetText("")
	d.podHostToIPMapField.SetText("")

	d.podNetworkField.SetOptions(networkOptions, nil)
	d.podNetworkField.SetCurrentOption(0)
	//d.podNetworkAliasesField.SetText("")

}

func (d *PodCreateDialog) setBasicInfoPageNextFocus() {
	if d.podNameField.HasFocus() {
		d.focusElement = podNoHostsCheckBoxFocus
	} else if d.podNoHostsCheckBox.HasFocus() {
		d.focusElement = podLabelsFieldFocus
	} else {
		d.focusElement = podFormFocus
	}
}

func (d *PodCreateDialog) setSecurityOptsPageNextFocus() {
	if d.podSelinuxLabelField.HasFocus() {
		d.focusElement = podApparmorFieldFocus
	} else if d.podApparmorField.HasFocus() {
		d.focusElement = podSeccompFieldFocus
	} else if d.podSeccompField.HasFocus() {
		d.focusElement = podMaskFieldFocus
	} else if d.podMaskField.HasFocus() {
		d.focusElement = podUnmaskFieldFocus
	} else if d.podUnmaskField.HasFocus() {
		d.focusElement = podNoNewPrivFieldFocus
	} else {
		d.focusElement = podFormFocus
	}
}

func (d *PodCreateDialog) setDNSSetupPageNextFocus() {
	if d.podDNSServerField.HasFocus() {
		d.focusElement = podDNSOptionsFieldFocus
	} else if d.podDNSOptionsField.HasFocus() {
		d.focusElement = podDNSSearchDomaindFieldFocus
	} else {
		d.focusElement = podFormFocus
	}
}

func (d *PodCreateDialog) setInfraSetupPageNextFocus() {
	if d.podInfraCheckBox.HasFocus() {
		d.focusElement = podInfraCommandFieldFocus
	} else if d.podInfraCommandField.HasFocus() {
		d.focusElement = podInfraImageFieldFocus
	} else {
		d.focusElement = podFormFocus
	}
}

func (d *PodCreateDialog) setNetworkingPageNextFocus() {
	if d.podHostnameField.HasFocus() {
		d.focusElement = podIPAddressFieldFocus
	} else if d.podIPAddressField.HasFocus() {
		d.focusElement = podMacAddressFieldFocus
	} else if d.podMacAddressField.HasFocus() {
		d.focusElement = podHostToIPMapFieldFocus
	} else if d.podHostToIPMapField.HasFocus() {
		d.focusElement = podNetworkFieldFocus
		//} else if d.podNetworkField.HasFocus() {
		//	d.focusElement = podNetworkAliasesFieldFocus
	} else {
		d.focusElement = podFormFocus
	}
}

// GetPodSpec returns pod create option spec
func (d *PodCreateDialog) GetPodSpec() pods.CreateOptions {

	var (
		labels           = make(map[string]string)
		dnsServers       []string
		dnsOptions       []string
		dnsSearchDomains []string
		hostAdd          []string
		infraCommand     []string
		network          string
		securityOpts     []string
	)
	for _, label := range strings.Split(d.podLabelsField.GetText(), " ") {
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

	for _, icmd := range strings.Split(d.podInfraCommandField.GetText(), " ") {
		if icmd != "" {
			infraCommand = append(infraCommand, icmd)
		}
	}

	for _, dns := range strings.Split(d.podDNSServerField.GetText(), " ") {
		if dns != "" {
			dnsServers = append(dnsServers, dns)
		}
	}
	for _, do := range strings.Split(d.podDNSOptionsField.GetText(), " ") {
		if do != "" {
			dnsOptions = append(dnsOptions, do)
		}
	}
	for _, ds := range strings.Split(d.podDNSSearchDomaindField.GetText(), " ") {
		if ds != "" {
			dnsSearchDomains = append(dnsSearchDomains, ds)
		}
	}

	for _, hadd := range strings.Split(d.podHostToIPMapField.GetText(), " ") {
		if hadd != "" {
			hostAdd = append(hostAdd, hadd)
		}
	}

	index, netName := d.podNetworkField.GetCurrentOption()
	if index > 0 {
		network = netName
	}

	// securuty options
	if d.podNoNewPrivField.IsChecked() {
		securityOpts = append(securityOpts, "no-new-privileges")
	}
	apparmor := strings.TrimSpace(d.podApparmorField.GetText())
	if apparmor != "" {
		securityOpts = append(securityOpts, fmt.Sprintf("apparmor=%s", apparmor))
	}
	seccomp := strings.TrimSpace(d.podSeccompField.GetText())
	if seccomp != "" {
		securityOpts = append(securityOpts, fmt.Sprintf("seccomp=%s", seccomp))
	}
	for _, selinuxLabel := range strings.Split(d.podSelinuxLabelField.GetText(), " ") {
		if selinuxLabel != "" {
			securityOpts = append(securityOpts, fmt.Sprintf("label=%s", selinuxLabel))
		}
	}
	mask := strings.TrimSpace(d.podMaskField.GetText())
	if seccomp != "" {
		securityOpts = append(securityOpts, fmt.Sprintf("mask=%s", mask))
	}
	unmask := strings.TrimSpace(d.podUnmaskField.GetText())
	if seccomp != "" {
		securityOpts = append(securityOpts, fmt.Sprintf("unmask=%s", unmask))
	}

	opts := pods.CreateOptions{
		Name:            d.podNameField.GetText(),
		NoHost:          d.podNoHostsCheckBox.IsChecked(),
		Labels:          labels,
		DNSServer:       dnsServers,
		DNSOptions:      dnsOptions,
		DNSSearchDomain: dnsSearchDomains,
		Infra:           d.podInfraCheckBox.IsChecked(),
		InfraImage:      d.podInfraImageField.GetText(),
		InfraCommand:    infraCommand,
		Hostname:        d.podHostnameField.GetText(),
		IPAddress:       d.podIPAddressField.GetText(),
		MacAddress:      d.podMacAddressField.GetText(),
		HostToIP:        hostAdd,
		Network:         network,
		SecurityOpts:    securityOpts,
	}
	return opts
}
