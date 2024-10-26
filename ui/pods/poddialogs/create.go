package poddialogs

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/networks"
	"github.com/containers/podman-tui/pdcs/pods"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
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
	podAddHostFieldFocus
	podNetworkFieldFocus
	podPublishFieldFocus
)

const (
	basicInfoPageIndex = 0 + iota
	dnsSetupPageIndex
	infraSetupPageIndex
	networkingPageIndex
	securityOptsPageIndex
)

// PodCreateDialog implements pod create dialog.
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
	podAddHostField          *tview.InputField
	podNetworkField          *tview.DropDown
	podPublishField          *tview.InputField
	cancelHandler            func()
	createHandler            func()
}

// NewPodCreateDialog returns new pod create dialog primitive PodCreateDialog.
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
			"Security Options",
		},
		activePageIndex:          0,
		display:                  false,
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
		podAddHostField:          tview.NewInputField(),
		podNetworkField:          tview.NewDropDown(),
		podPublishField:          tview.NewInputField(),
	}

	podDialog.categories.SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)
	podDialog.categories.SetBackgroundColor(style.DialogBgColor)
	podDialog.categories.SetBorder(true)
	podDialog.categories.SetBorderColor(style.DialogSubBoxBorderColor)

	// basic information setup page
	basicInfoPageLabelWidth := 12
	// name field
	podDialog.podNameField.SetLabel("name:")
	podDialog.podNameField.SetLabelWidth(basicInfoPageLabelWidth)
	podDialog.podNameField.SetBackgroundColor(style.DialogBgColor)
	podDialog.podNameField.SetLabelColor(style.DialogFgColor)
	podDialog.podNameField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// no hosts check box
	podDialog.podNoHostsCheckBox.SetLabel("no hosts")
	podDialog.podNoHostsCheckBox.SetLabelWidth(basicInfoPageLabelWidth)
	podDialog.podNoHostsCheckBox.SetChecked(false)
	podDialog.podNoHostsCheckBox.SetBackgroundColor(style.DialogBgColor)
	podDialog.podNoHostsCheckBox.SetLabelColor(style.DialogFgColor)
	podDialog.podNoHostsCheckBox.SetFieldBackgroundColor(style.InputFieldBgColor)

	// labels field
	podDialog.podLabelsField.SetLabel("labels:")
	podDialog.podLabelsField.SetLabelWidth(basicInfoPageLabelWidth)
	podDialog.podLabelsField.SetBackgroundColor(style.DialogBgColor)
	podDialog.podLabelsField.SetLabelColor(style.DialogFgColor)
	podDialog.podLabelsField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// security options
	securityOptsPageLabelWidth := 10
	// labels
	podDialog.podSelinuxLabelField.SetLabel("label:")
	podDialog.podSelinuxLabelField.SetLabelWidth(securityOptsPageLabelWidth)
	podDialog.podSelinuxLabelField.SetBackgroundColor(style.DialogBgColor)
	podDialog.podSelinuxLabelField.SetLabelColor(style.DialogFgColor)
	podDialog.podSelinuxLabelField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// apparmor
	podDialog.podApparmorField.SetLabel("apparmor:")
	podDialog.podApparmorField.SetLabelWidth(securityOptsPageLabelWidth)
	podDialog.podApparmorField.SetBackgroundColor(style.DialogBgColor)
	podDialog.podApparmorField.SetLabelColor(style.DialogFgColor)
	podDialog.podApparmorField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// seccomp
	podDialog.podSeccompField.SetLabel("seccomp:")
	podDialog.podSeccompField.SetLabelWidth(securityOptsPageLabelWidth)
	podDialog.podSeccompField.SetBackgroundColor(style.DialogBgColor)
	podDialog.podSeccompField.SetLabelColor(style.DialogFgColor)
	podDialog.podSeccompField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// mask
	podDialog.podMaskField.SetLabel("mask:")
	podDialog.podMaskField.SetLabelWidth(securityOptsPageLabelWidth)
	podDialog.podMaskField.SetBackgroundColor(style.DialogBgColor)
	podDialog.podMaskField.SetLabelColor(style.DialogFgColor)
	podDialog.podMaskField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// unmask
	podDialog.podUnmaskField.SetLabel("unmask:")
	podDialog.podUnmaskField.SetLabelWidth(securityOptsPageLabelWidth)
	podDialog.podUnmaskField.SetBackgroundColor(style.DialogBgColor)
	podDialog.podUnmaskField.SetLabelColor(style.DialogFgColor)
	podDialog.podUnmaskField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// no new privileges
	podDialog.podNoNewPrivField.SetLabel("no new privileges ")
	podDialog.podNoNewPrivField.SetBackgroundColor(style.DialogBgColor)
	podDialog.podNoNewPrivField.SetLabelColor(tcell.ColorWhite)
	podDialog.podNoNewPrivField.SetBackgroundColor(style.DialogBgColor)
	podDialog.podNoNewPrivField.SetLabelColor(style.DialogFgColor)
	podDialog.podNoNewPrivField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// DNS setup page
	dnsPageLabelWidth := 16
	// DNS server field
	podDialog.podDNSServerField.SetLabel("dns servers:")
	podDialog.podDNSServerField.SetLabelWidth(dnsPageLabelWidth)
	podDialog.podDNSServerField.SetBackgroundColor(style.DialogBgColor)
	podDialog.podDNSServerField.SetLabelColor(style.DialogFgColor)
	podDialog.podDNSServerField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// DNS options field
	podDialog.podDNSOptionsField.SetLabel("dns options:")
	podDialog.podDNSOptionsField.SetLabelWidth(dnsPageLabelWidth)
	podDialog.podDNSOptionsField.SetBackgroundColor(style.DialogBgColor)
	podDialog.podDNSOptionsField.SetLabelColor(style.DialogFgColor)
	podDialog.podDNSOptionsField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// DNS search domains field
	podDialog.podDNSSearchDomaindField.SetLabel("search domains:")
	podDialog.podDNSSearchDomaindField.SetLabelWidth(dnsPageLabelWidth)
	podDialog.podDNSSearchDomaindField.SetBackgroundColor(style.DialogBgColor)
	podDialog.podDNSSearchDomaindField.SetLabelColor(style.DialogFgColor)
	podDialog.podDNSSearchDomaindField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// infra page
	infraPageLabelWidth := 15
	// infra check box
	podDialog.podInfraCheckBox.SetLabel("infra")
	podDialog.podInfraCheckBox.SetLabelWidth(infraPageLabelWidth)
	podDialog.podInfraCheckBox.SetChecked(true)
	podDialog.podInfraCheckBox.SetBackgroundColor(style.DialogBgColor)
	podDialog.podInfraCheckBox.SetLabelColor(style.DialogFgColor)
	podDialog.podInfraCheckBox.SetFieldBackgroundColor(style.InputFieldBgColor)

	// infra command field
	podDialog.podInfraCommandField.SetLabel("infra command:")
	podDialog.podInfraCommandField.SetLabelWidth(infraPageLabelWidth)
	podDialog.podInfraCommandField.SetBackgroundColor(style.DialogBgColor)
	podDialog.podInfraCommandField.SetLabelColor(style.DialogFgColor)
	podDialog.podInfraCommandField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// infra image field
	podDialog.podInfraImageField.SetLabel("infra image:")
	podDialog.podInfraImageField.SetText("")
	podDialog.podInfraImageField.SetLabelWidth(infraPageLabelWidth)
	podDialog.podInfraImageField.SetBackgroundColor(style.DialogBgColor)
	podDialog.podInfraImageField.SetLabelColor(style.DialogFgColor)
	podDialog.podInfraImageField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// networking page
	networkingLabelWidth := 17
	// hostname field
	podDialog.podHostnameField.SetLabel("hostname:")
	podDialog.podHostnameField.SetLabelWidth(networkingLabelWidth)
	podDialog.podHostnameField.SetBackgroundColor(style.DialogBgColor)
	podDialog.podHostnameField.SetLabelColor(style.DialogFgColor)
	podDialog.podHostnameField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// ip address field
	podDialog.podIPAddressField.SetLabel("ip address:")
	podDialog.podIPAddressField.SetLabelWidth(networkingLabelWidth)
	podDialog.podIPAddressField.SetBackgroundColor(style.DialogBgColor)
	podDialog.podIPAddressField.SetLabelColor(style.DialogFgColor)
	podDialog.podIPAddressField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// mac address field
	podDialog.podMacAddressField.SetLabel("mac address:")
	podDialog.podMacAddressField.SetLabelWidth(networkingLabelWidth)
	podDialog.podMacAddressField.SetBackgroundColor(style.DialogBgColor)
	podDialog.podMacAddressField.SetLabelColor(style.DialogFgColor)
	podDialog.podMacAddressField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// add host field
	podDialog.podAddHostField.SetLabel("add host:")
	podDialog.podAddHostField.SetLabelWidth(networkingLabelWidth)
	podDialog.podAddHostField.SetBackgroundColor(style.DialogBgColor)
	podDialog.podAddHostField.SetLabelColor(style.DialogFgColor)
	podDialog.podAddHostField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// network field
	podDialog.podNetworkField.SetLabel("network:")
	podDialog.podNetworkField.SetLabelWidth(networkingLabelWidth)
	podDialog.podNetworkField.SetBackgroundColor(style.DialogBgColor)
	podDialog.podNetworkField.SetLabelColor(style.DialogFgColor)
	podDialog.podNetworkField.SetListStyles(style.DropDownUnselected, style.DropDownSelected)
	podDialog.podNetworkField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// publish field
	podDialog.podPublishField.SetLabel("publish:")
	podDialog.podPublishField.SetLabelWidth(networkingLabelWidth)
	podDialog.podPublishField.SetBackgroundColor(style.DialogBgColor)
	podDialog.podPublishField.SetLabelColor(style.DialogFgColor)
	podDialog.podPublishField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// category pages
	podDialog.categoryPages.SetBackgroundColor(style.DialogBgColor)
	podDialog.categoryPages.SetBorder(true)
	podDialog.categoryPages.SetBorderColor(style.DialogSubBoxBorderColor)

	// form
	podDialog.form.SetBackgroundColor(style.DialogBgColor)
	podDialog.form.AddButton("Cancel", nil)
	podDialog.form.AddButton("Create", nil)
	podDialog.form.SetButtonsAlign(tview.AlignRight)
	podDialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	podDialog.setupLayout()
	podDialog.layout.SetBackgroundColor(style.DialogBgColor)
	podDialog.layout.SetBorder(true)
	podDialog.layout.SetBorderColor(style.DialogBorderColor)
	podDialog.layout.SetTitle("PODMAN POD CREATE")
	podDialog.layout.AddItem(podDialog.form, dialogs.DialogFormHeight, 0, true)

	podDialog.setActiveCategory(0)

	podDialog.initCustomInputHanlers()

	return &podDialog
}

func (d *PodCreateDialog) setupLayout() {
	bgColor := style.DialogBgColor

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
	d.networkingPage.AddItem(d.podAddHostField, 1, 0, true)
	d.networkingPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.networkingPage.AddItem(d.podNetworkField, 1, 0, true)
	d.networkingPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.networkingPage.AddItem(d.podPublishField, 1, 0, true)
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
	layout.AddItem(d.categories, layoutWidth+6, 0, true) //nolint:mnd
	layout.AddItem(d.categoryPages, 0, 1, true)
	layout.SetBackgroundColor(bgColor)

	d.layout.AddItem(layout, 0, 1, true)
}

// Display displays this primitive.
func (d *PodCreateDialog) Display() {
	d.display = true
	d.initData()
	d.focusElement = categoryPagesFocus
}

// IsDisplay returns true if primitive is shown.
func (d *PodCreateDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *PodCreateDialog) Hide() {
	d.display = false
}

// HasFocus returns whether or not this primitive has focus.
func (d *PodCreateDialog) HasFocus() bool {
	if d.categories.HasFocus() || d.categoryPages.HasFocus() {
		return true
	}

	return d.Box.HasFocus() || d.form.HasFocus()
}

// dropdownHasFocus returns true if pod create dialog dropdown primitives.
// has focus.
func (d *PodCreateDialog) dropdownHasFocus() bool {
	return d.podNetworkField.HasFocus()
}

// Focus is called when this primitive receives focus.
func (d *PodCreateDialog) Focus(delegate func(p tview.Primitive)) { //nolint:cyclop
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
	case podAddHostFieldFocus:
		delegate(d.podAddHostField)
	case podNetworkFieldFocus:
		delegate(d.podNetworkField)
	case podPublishFieldFocus:
		delegate(d.podPublishField)
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

// InputHandler returns input handler function for this primitive.
func (d *PodCreateDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) { //nolint:gocognit,lll,cyclop
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("pod create dialog: event %v received", event)

		if event.Key() == tcell.KeyEsc && !d.dropdownHasFocus() {
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
func (d *PodCreateDialog) SetRect(x, y, width, height int) {
	if width > podCreateDialogMaxWidth {
		emptySpace := (width - podCreateDialogMaxWidth) / 2 //nolint:mnd
		x += emptySpace
		width = podCreateDialogMaxWidth
	}

	if height > podCreateDialogHeight {
		emptySpace := (height - podCreateDialogHeight) / 2 //nolint:mnd
		y += emptySpace
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

// SetCancelFunc sets form cancel button selected function.
func (d *PodCreateDialog) SetCancelFunc(handler func()) *PodCreateDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:mnd
	cancelButton.SetSelectedFunc(handler)

	return d
}

// SetCreateFunc sets form create button selected function.
func (d *PodCreateDialog) SetCreateFunc(handler func()) *PodCreateDialog {
	d.createHandler = handler
	enterButton := d.form.GetButton(d.form.GetButtonCount() - 1)

	enterButton.SetSelectedFunc(handler)

	return d
}

func (d *PodCreateDialog) setActiveCategory(index int) {
	fgColor := style.DialogFgColor
	bgColor := style.ButtonBgColor
	ctgTextColor := style.GetColorHex(fgColor)
	ctgBgColor := style.GetColorHex(bgColor)

	d.activePageIndex = index

	d.categories.Clear()

	ctgList := []string{}

	alignedList, _ := utils.AlignStringListWidth(d.categoryLabels)

	for i := range len(d.categoryLabels) {
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
		activePage++
		d.setActiveCategory(activePage)

		return
	}

	d.setActiveCategory(0)
}

func (d *PodCreateDialog) previousCategory() {
	activePage := d.activePageIndex
	if d.activePageIndex > 0 {
		activePage--
		d.setActiveCategory(activePage)

		return
	}

	d.setActiveCategory(len(d.categoryLabels) - 1)
}

func (d *PodCreateDialog) initData() {
	// get available networks
	networkOptions := []string{""}
	networkList, _ := networks.List()

	for i := range networkList {
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
	d.podInfraImageField.SetText("")

	d.podHostnameField.SetText("")
	d.podIPAddressField.SetText("")
	d.podMacAddressField.SetText("")
	d.podAddHostField.SetText("")

	d.podNetworkField.SetOptions(networkOptions, nil)
	d.podNetworkField.SetCurrentOption(0)
	d.podPublishField.SetText("")
}

func (d *PodCreateDialog) setBasicInfoPageNextFocus() {
	if d.podNameField.HasFocus() {
		d.focusElement = podNoHostsCheckBoxFocus

		return
	}

	if d.podNoHostsCheckBox.HasFocus() {
		d.focusElement = podLabelsFieldFocus

		return
	}

	d.focusElement = podFormFocus
}

func (d *PodCreateDialog) setSecurityOptsPageNextFocus() {
	if d.podSelinuxLabelField.HasFocus() {
		d.focusElement = podApparmorFieldFocus

		return
	}

	if d.podApparmorField.HasFocus() {
		d.focusElement = podSeccompFieldFocus

		return
	}

	if d.podSeccompField.HasFocus() {
		d.focusElement = podMaskFieldFocus

		return
	}

	if d.podMaskField.HasFocus() {
		d.focusElement = podUnmaskFieldFocus

		return
	}

	if d.podUnmaskField.HasFocus() {
		d.focusElement = podNoNewPrivFieldFocus

		return
	}

	d.focusElement = podFormFocus
}

func (d *PodCreateDialog) setDNSSetupPageNextFocus() {
	if d.podDNSServerField.HasFocus() {
		d.focusElement = podDNSOptionsFieldFocus

		return
	}

	if d.podDNSOptionsField.HasFocus() {
		d.focusElement = podDNSSearchDomaindFieldFocus

		return
	}

	d.focusElement = podFormFocus
}

func (d *PodCreateDialog) setInfraSetupPageNextFocus() {
	if d.podInfraCheckBox.HasFocus() {
		d.focusElement = podInfraCommandFieldFocus

		return
	}

	if d.podInfraCommandField.HasFocus() {
		d.focusElement = podInfraImageFieldFocus

		return
	}

	d.focusElement = podFormFocus
}

func (d *PodCreateDialog) setNetworkingPageNextFocus() {
	if d.podHostnameField.HasFocus() {
		d.focusElement = podIPAddressFieldFocus

		return
	}

	if d.podIPAddressField.HasFocus() {
		d.focusElement = podMacAddressFieldFocus

		return
	}

	if d.podMacAddressField.HasFocus() {
		d.focusElement = podAddHostFieldFocus

		return
	}

	if d.podAddHostField.HasFocus() {
		d.focusElement = podNetworkFieldFocus

		return
	}

	if d.podNetworkField.HasFocus() {
		d.focusElement = podPublishFieldFocus

		return
	}

	d.focusElement = podFormFocus
}

// GetPodSpec returns pod create option spec.
func (d *PodCreateDialog) GetPodSpec() pods.CreateOptions { //nolint:gocognit,cyclop
	var (
		labels           = make(map[string]string)
		dnsServers       []string
		dnsOptions       []string
		dnsSearchDomains []string
		addHost          []string
		network          string
		securityOpts     []string
		publish          []string
	)

	for _, label := range strings.Split(d.podLabelsField.GetText(), " ") {
		if label != "" {
			split := strings.Split(label, "=")
			if len(split) == 2 { //nolint:mnd
				key := split[0]
				value := split[1]

				if key != "" && value != "" {
					labels[key] = value
				}
			}
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

	for _, hadd := range strings.Split(d.podAddHostField.GetText(), " ") {
		if hadd != "" {
			addHost = append(addHost, hadd)
		}
	}

	index, netName := d.podNetworkField.GetCurrentOption()
	if index > 0 {
		network = netName
	}

	for _, p := range strings.Split(d.podPublishField.GetText(), " ") {
		if p != "" {
			publish = append(publish, p)
		}
	}

	// securuty options
	if d.podNoNewPrivField.IsChecked() {
		securityOpts = append(securityOpts, "no-new-privileges")
	}

	apparmor := strings.TrimSpace(d.podApparmorField.GetText())
	if apparmor != "" {
		securityOpts = append(securityOpts, fmt.Sprintf("apparmor=%s", apparmor)) //nolint:perfsprint
	}

	seccomp := strings.TrimSpace(d.podSeccompField.GetText())
	if seccomp != "" {
		securityOpts = append(securityOpts, fmt.Sprintf("seccomp=%s", seccomp)) //nolint:perfsprint
	}

	for _, selinuxLabel := range strings.Split(d.podSelinuxLabelField.GetText(), " ") {
		if selinuxLabel != "" {
			securityOpts = append(securityOpts, fmt.Sprintf("label=%s", selinuxLabel)) //nolint:perfsprint
		}
	}

	mask := strings.TrimSpace(d.podMaskField.GetText())
	if seccomp != "" {
		securityOpts = append(securityOpts, fmt.Sprintf("mask=%s", mask)) //nolint:perfsprint
	}

	unmask := strings.TrimSpace(d.podUnmaskField.GetText())
	if seccomp != "" {
		securityOpts = append(securityOpts, fmt.Sprintf("unmask=%s", unmask)) //nolint:perfsprint
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
		InfraCommand:    d.podInfraCommandField.GetText(),
		Hostname:        d.podHostnameField.GetText(),
		IPAddress:       d.podIPAddressField.GetText(),
		MacAddress:      d.podMacAddressField.GetText(),
		AddHost:         addHost,
		Network:         network,
		SecurityOpts:    securityOpts,
		Publish:         publish,
	}

	return opts
}
