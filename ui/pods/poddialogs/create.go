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
	podCreateDialogMaxWidth = 90
	podCreateDialogHeight   = 19
)

const (
	createPodFormFocus = 0 + iota
	createPodCategoriesFocus
	createPodCategoryPagesFocus
	createPodNameFieldFocus
	createPodNoHostsCheckBoxFocus
	createPodLabelsFieldFocus
	createPodSelinuxLabelFieldFocus
	createPodApparmorFieldFocus
	createPodSeccompFieldFocus
	createPodMaskFieldFocus
	createPodUnmaskFieldFocus
	createPodNoNewPrivFieldFocus
	createPodDNSServerFieldFocus
	createPodDNSOptionsFieldFocus
	createPodDNSSearchDomaindFieldFocus
	createPodInfraCheckBoxFocus
	createPodInfraCommandFieldFocus
	createPodInfraImageFieldFocus
	createPodHostnameFieldFocus
	createPodIPAddressFieldFocus
	createPodMacAddressFieldFocus
	createPodAddHostFieldFocus
	createPodNetworkFieldFocus
	createPodPublishFieldFocus
	createPodMemoryFieldFocus
	createPodMemorySwapFieldFocus
	createPodCPUsFieldFocus
	createPodCPUSharesFieldFocus
	createPodCPUSetCPUsFieldFocus
	createPodCPUSetMemsFieldFocus
	createPodShmSizeFieldFocus
	createPodShmSizeSystemdFieldFocus
	createPodNamespaceShareFieldFocus
	createPodNamespacePidFieldFocus
	createPodNamespaceUserFieldFocus
	createPodNamespaceUtsFieldFocus
	createPodNamespaceUidmapFieldFocus
	createPodNamespaceSubuidNameFieldFocus
	createPodNamespaceGidmapFieldFocus
	createPodNamespaceSubgidNameFieldFocus
)

const (
	createPodBasicInfoPageIndex = 0 + iota
	createPodDNSSetupPageIndex
	createPodInfraSetupPageIndex
	createPodNetworkingPageIndex
	createPodSecurityOptsPageIndex
	createPodResourceSettingsPageIndex
	createPodNamespaceOptionsPageIndex
)

// PodCreateDialog implements pod create dialog.
type PodCreateDialog struct {
	*tview.Box
	layout                      *tview.Flex
	categoryLabels              []string
	categories                  *tview.TextView
	categoryPages               *tview.Pages
	basicInfoPage               *tview.Flex
	securityOptsPage            *tview.Flex
	dnsSetupPage                *tview.Flex
	infraSetupPage              *tview.Flex
	networkingPage              *tview.Flex
	resourcePage                *tview.Flex
	namespacePage               *tview.Flex
	form                        *tview.Form
	display                     bool
	activePageIndex             int
	focusElement                int
	podNameField                *tview.InputField
	podNoHostsCheckBox          *tview.Checkbox
	podLabelsField              *tview.InputField
	podSelinuxLabelField        *tview.InputField
	podApparmorField            *tview.InputField
	podSeccompField             *tview.InputField
	podMaskField                *tview.InputField
	podUnmaskField              *tview.InputField
	podNoNewPrivField           *tview.Checkbox
	podDNSServerField           *tview.InputField
	podDNSOptionsField          *tview.InputField
	podDNSSearchDomaindField    *tview.InputField
	podInfraCheckBox            *tview.Checkbox
	podInfraCommandField        *tview.InputField
	podInfraImageField          *tview.InputField
	podHostnameField            *tview.InputField
	podIPAddressField           *tview.InputField
	podMacAddressField          *tview.InputField
	podAddHostField             *tview.InputField
	podNetworkField             *tview.DropDown
	podPublishField             *tview.InputField
	podMemoryField              *tview.InputField
	podMemorySwapField          *tview.InputField
	podCPUsField                *tview.InputField
	podCPUSharesField           *tview.InputField
	podCPUSetCPUsField          *tview.InputField
	podCPUSetMemsField          *tview.InputField
	podShmSizeField             *tview.InputField
	podShmSizeSystemdField      *tview.InputField
	podNamespaceShareField      *tview.InputField
	podNamespacePidField        *tview.InputField
	podNamespaceUserField       *tview.InputField
	podNamespaceUtsField        *tview.InputField
	podNamespaceUidmapField     *tview.InputField
	podNamespaceSubuidNameField *tview.InputField
	podNamespaceGidmapField     *tview.InputField
	podNamespaceSubgidNameField *tview.InputField
	cancelHandler               func()
	createHandler               func()
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
		resourcePage:     tview.NewFlex(),
		namespacePage:    tview.NewFlex(),
		form:             tview.NewForm(),
		categoryLabels: []string{
			"Basic Information",
			"DNS Setup",
			"Infra Setup",
			"Networking",
			"Security Options",
			"Resource Settings",
			"Namespace Options",
		},
		activePageIndex:             0,
		display:                     false,
		podNameField:                tview.NewInputField(),
		podNoHostsCheckBox:          tview.NewCheckbox(),
		podLabelsField:              tview.NewInputField(),
		podSelinuxLabelField:        tview.NewInputField(),
		podApparmorField:            tview.NewInputField(),
		podSeccompField:             tview.NewInputField(),
		podMaskField:                tview.NewInputField(),
		podUnmaskField:              tview.NewInputField(),
		podNoNewPrivField:           tview.NewCheckbox(),
		podDNSServerField:           tview.NewInputField(),
		podDNSOptionsField:          tview.NewInputField(),
		podDNSSearchDomaindField:    tview.NewInputField(),
		podInfraCheckBox:            tview.NewCheckbox(),
		podInfraCommandField:        tview.NewInputField(),
		podInfraImageField:          tview.NewInputField(),
		podHostnameField:            tview.NewInputField(),
		podIPAddressField:           tview.NewInputField(),
		podMacAddressField:          tview.NewInputField(),
		podAddHostField:             tview.NewInputField(),
		podNetworkField:             tview.NewDropDown(),
		podPublishField:             tview.NewInputField(),
		podMemoryField:              tview.NewInputField(),
		podMemorySwapField:          tview.NewInputField(),
		podCPUsField:                tview.NewInputField(),
		podCPUSharesField:           tview.NewInputField(),
		podCPUSetCPUsField:          tview.NewInputField(),
		podCPUSetMemsField:          tview.NewInputField(),
		podShmSizeField:             tview.NewInputField(),
		podShmSizeSystemdField:      tview.NewInputField(),
		podNamespaceShareField:      tview.NewInputField(),
		podNamespacePidField:        tview.NewInputField(),
		podNamespaceUserField:       tview.NewInputField(),
		podNamespaceUtsField:        tview.NewInputField(),
		podNamespaceUidmapField:     tview.NewInputField(),
		podNamespaceSubuidNameField: tview.NewInputField(),
		podNamespaceGidmapField:     tview.NewInputField(),
		podNamespaceSubgidNameField: tview.NewInputField(),
	}

	podDialog.categories.SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)
	podDialog.categories.SetBackgroundColor(style.DialogBgColor)
	podDialog.categories.SetBorder(true)
	podDialog.categories.SetBorderColor(style.DialogSubBoxBorderColor)

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

	podDialog.initCustomInputHandlers()

	return &podDialog
}

func (d *PodCreateDialog) setupLayout() {
	d.setupBasicInfoUI()
	d.setupDNSSetupUI()
	d.setupInfraSetupUI()
	d.setupNetworkingUI()
	d.setupSecurityOptionsUI()
	d.setupResourceSettingsUI()
	d.setupNamespaceOptionsUI()

	// adding category pages
	d.categoryPages.AddPage(d.categoryLabels[createPodBasicInfoPageIndex], d.basicInfoPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[createPodDNSSetupPageIndex], d.dnsSetupPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[createPodInfraSetupPageIndex], d.infraSetupPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[createPodNetworkingPageIndex], d.networkingPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[createPodSecurityOptsPageIndex], d.securityOptsPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[createPodResourceSettingsPageIndex], d.resourcePage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[createPodNamespaceOptionsPageIndex], d.namespacePage, true, true)

	// add it to layout.
	_, layoutWidth := utils.AlignStringListWidth(d.categoryLabels)
	layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	layout.AddItem(d.categories, layoutWidth+6, 0, true) //nolint:mnd
	layout.AddItem(d.categoryPages, 0, 1, true)
	layout.SetBackgroundColor(style.DialogBgColor)

	d.layout.AddItem(layout, 0, 1, true)
}

func (d *PodCreateDialog) setupBasicInfoUI() {
	// basic information setup page
	basicInfoPageLabelWidth := 12
	// name field
	d.podNameField.SetLabel("name:")
	d.podNameField.SetLabelWidth(basicInfoPageLabelWidth)
	d.podNameField.SetBackgroundColor(style.DialogBgColor)
	d.podNameField.SetLabelColor(style.DialogFgColor)
	d.podNameField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// no hosts check box
	d.podNoHostsCheckBox.SetLabel("no hosts")
	d.podNoHostsCheckBox.SetLabelWidth(basicInfoPageLabelWidth)
	d.podNoHostsCheckBox.SetChecked(false)
	d.podNoHostsCheckBox.SetBackgroundColor(style.DialogBgColor)
	d.podNoHostsCheckBox.SetLabelColor(style.DialogFgColor)
	d.podNoHostsCheckBox.SetFieldBackgroundColor(style.InputFieldBgColor)

	// labels field
	d.podLabelsField.SetLabel("labels:")
	d.podLabelsField.SetLabelWidth(basicInfoPageLabelWidth)
	d.podLabelsField.SetBackgroundColor(style.DialogBgColor)
	d.podLabelsField.SetLabelColor(style.DialogFgColor)
	d.podLabelsField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// basic info page
	d.basicInfoPage.SetDirection(tview.FlexRow)
	d.basicInfoPage.AddItem(d.podNameField, 1, 0, true)
	d.basicInfoPage.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	d.basicInfoPage.AddItem(d.podNoHostsCheckBox, 1, 0, true)
	d.basicInfoPage.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	d.basicInfoPage.AddItem(d.podLabelsField, 1, 0, true)
	d.basicInfoPage.SetBackgroundColor(style.DialogBgColor)
}

func (d *PodCreateDialog) setupDNSSetupUI() {
	// DNS setup page
	dnsPageLabelWidth := 16
	// DNS server field
	d.podDNSServerField.SetLabel("dns servers:")
	d.podDNSServerField.SetLabelWidth(dnsPageLabelWidth)
	d.podDNSServerField.SetBackgroundColor(style.DialogBgColor)
	d.podDNSServerField.SetLabelColor(style.DialogFgColor)
	d.podDNSServerField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// DNS options field
	d.podDNSOptionsField.SetLabel("dns options:")
	d.podDNSOptionsField.SetLabelWidth(dnsPageLabelWidth)
	d.podDNSOptionsField.SetBackgroundColor(style.DialogBgColor)
	d.podDNSOptionsField.SetLabelColor(style.DialogFgColor)
	d.podDNSOptionsField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// DNS search domains field
	d.podDNSSearchDomaindField.SetLabel("search domains:")
	d.podDNSSearchDomaindField.SetLabelWidth(dnsPageLabelWidth)
	d.podDNSSearchDomaindField.SetBackgroundColor(style.DialogBgColor)
	d.podDNSSearchDomaindField.SetLabelColor(style.DialogFgColor)
	d.podDNSSearchDomaindField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// DNS setup page
	d.dnsSetupPage.SetDirection(tview.FlexRow)
	d.dnsSetupPage.AddItem(d.podDNSServerField, 1, 0, true)
	d.dnsSetupPage.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	d.dnsSetupPage.AddItem(d.podDNSOptionsField, 1, 0, true)
	d.dnsSetupPage.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	d.dnsSetupPage.AddItem(d.podDNSSearchDomaindField, 1, 0, true)
	d.dnsSetupPage.SetBackgroundColor(style.DialogBgColor)
}

func (d *PodCreateDialog) setupInfraSetupUI() {
	// infra page
	infraPageLabelWidth := 15
	// infra check box
	d.podInfraCheckBox.SetLabel("infra")
	d.podInfraCheckBox.SetLabelWidth(infraPageLabelWidth)
	d.podInfraCheckBox.SetChecked(true)
	d.podInfraCheckBox.SetBackgroundColor(style.DialogBgColor)
	d.podInfraCheckBox.SetLabelColor(style.DialogFgColor)
	d.podInfraCheckBox.SetFieldBackgroundColor(style.InputFieldBgColor)

	// infra command field
	d.podInfraCommandField.SetLabel("infra command:")
	d.podInfraCommandField.SetLabelWidth(infraPageLabelWidth)
	d.podInfraCommandField.SetBackgroundColor(style.DialogBgColor)
	d.podInfraCommandField.SetLabelColor(style.DialogFgColor)
	d.podInfraCommandField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// infra image field
	d.podInfraImageField.SetLabel("infra image:")
	d.podInfraImageField.SetText("")
	d.podInfraImageField.SetLabelWidth(infraPageLabelWidth)
	d.podInfraImageField.SetBackgroundColor(style.DialogBgColor)
	d.podInfraImageField.SetLabelColor(style.DialogFgColor)
	d.podInfraImageField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// infra page
	d.infraSetupPage.SetDirection(tview.FlexRow)
	d.infraSetupPage.AddItem(d.podInfraCheckBox, 1, 0, true)
	d.infraSetupPage.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	d.infraSetupPage.AddItem(d.podInfraCommandField, 1, 0, true)
	d.infraSetupPage.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	d.infraSetupPage.AddItem(d.podInfraImageField, 1, 0, true)
	d.infraSetupPage.SetBackgroundColor(style.DialogBgColor)
}

func (d *PodCreateDialog) setupNetworkingUI() {
	// networking page
	networkingLabelWidth := 17
	// hostname field
	d.podHostnameField.SetLabel("hostname:")
	d.podHostnameField.SetLabelWidth(networkingLabelWidth)
	d.podHostnameField.SetBackgroundColor(style.DialogBgColor)
	d.podHostnameField.SetLabelColor(style.DialogFgColor)
	d.podHostnameField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// ip address field
	d.podIPAddressField.SetLabel("ip address:")
	d.podIPAddressField.SetLabelWidth(networkingLabelWidth)
	d.podIPAddressField.SetBackgroundColor(style.DialogBgColor)
	d.podIPAddressField.SetLabelColor(style.DialogFgColor)
	d.podIPAddressField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// mac address field
	d.podMacAddressField.SetLabel("mac address:")
	d.podMacAddressField.SetLabelWidth(networkingLabelWidth)
	d.podMacAddressField.SetBackgroundColor(style.DialogBgColor)
	d.podMacAddressField.SetLabelColor(style.DialogFgColor)
	d.podMacAddressField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// add host field
	d.podAddHostField.SetLabel("add host:")
	d.podAddHostField.SetLabelWidth(networkingLabelWidth)
	d.podAddHostField.SetBackgroundColor(style.DialogBgColor)
	d.podAddHostField.SetLabelColor(style.DialogFgColor)
	d.podAddHostField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// network field
	d.podNetworkField.SetLabel("network:")
	d.podNetworkField.SetLabelWidth(networkingLabelWidth)
	d.podNetworkField.SetBackgroundColor(style.DialogBgColor)
	d.podNetworkField.SetLabelColor(style.DialogFgColor)
	d.podNetworkField.SetListStyles(style.DropDownUnselected, style.DropDownSelected)
	d.podNetworkField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// publish field
	d.podPublishField.SetLabel("publish:")
	d.podPublishField.SetLabelWidth(networkingLabelWidth)
	d.podPublishField.SetBackgroundColor(style.DialogBgColor)
	d.podPublishField.SetLabelColor(style.DialogFgColor)
	d.podPublishField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// networking page
	d.networkingPage.SetDirection(tview.FlexRow)
	d.networkingPage.AddItem(d.podHostnameField, 1, 0, true)
	d.networkingPage.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	d.networkingPage.AddItem(d.podIPAddressField, 1, 0, true)
	d.networkingPage.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	d.networkingPage.AddItem(d.podMacAddressField, 1, 0, true)
	d.networkingPage.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	d.networkingPage.AddItem(d.podAddHostField, 1, 0, true)
	d.networkingPage.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	d.networkingPage.AddItem(d.podNetworkField, 1, 0, true)
	d.networkingPage.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	d.networkingPage.AddItem(d.podPublishField, 1, 0, true)
	d.networkingPage.SetBackgroundColor(style.DialogBgColor)
}

func (d *PodCreateDialog) setupSecurityOptionsUI() {
	// security options
	securityOptsPageLabelWidth := 10
	// labels
	d.podSelinuxLabelField.SetLabel("label:")
	d.podSelinuxLabelField.SetLabelWidth(securityOptsPageLabelWidth)
	d.podSelinuxLabelField.SetBackgroundColor(style.DialogBgColor)
	d.podSelinuxLabelField.SetLabelColor(style.DialogFgColor)
	d.podSelinuxLabelField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// apparmor
	d.podApparmorField.SetLabel("apparmor:")
	d.podApparmorField.SetLabelWidth(securityOptsPageLabelWidth)
	d.podApparmorField.SetBackgroundColor(style.DialogBgColor)
	d.podApparmorField.SetLabelColor(style.DialogFgColor)
	d.podApparmorField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// seccomp
	d.podSeccompField.SetLabel("seccomp:")
	d.podSeccompField.SetLabelWidth(securityOptsPageLabelWidth)
	d.podSeccompField.SetBackgroundColor(style.DialogBgColor)
	d.podSeccompField.SetLabelColor(style.DialogFgColor)
	d.podSeccompField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// mask
	d.podMaskField.SetLabel("mask:")
	d.podMaskField.SetLabelWidth(securityOptsPageLabelWidth)
	d.podMaskField.SetBackgroundColor(style.DialogBgColor)
	d.podMaskField.SetLabelColor(style.DialogFgColor)
	d.podMaskField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// unmask
	d.podUnmaskField.SetLabel("unmask:")
	d.podUnmaskField.SetLabelWidth(securityOptsPageLabelWidth)
	d.podUnmaskField.SetBackgroundColor(style.DialogBgColor)
	d.podUnmaskField.SetLabelColor(style.DialogFgColor)
	d.podUnmaskField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// no new privileges
	d.podNoNewPrivField.SetLabel("no new privileges ")
	d.podNoNewPrivField.SetBackgroundColor(style.DialogBgColor)
	d.podNoNewPrivField.SetLabelColor(tcell.ColorWhite)
	d.podNoNewPrivField.SetBackgroundColor(style.DialogBgColor)
	d.podNoNewPrivField.SetLabelColor(style.DialogFgColor)
	d.podNoNewPrivField.SetFieldBackgroundColor(style.InputFieldBgColor)

	// security options page
	d.securityOptsPage.SetDirection(tview.FlexRow)
	d.securityOptsPage.AddItem(d.podSelinuxLabelField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.podApparmorField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.podSeccompField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.podMaskField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.podUnmaskField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.podNoNewPrivField, 1, 0, true)
	d.securityOptsPage.SetBackgroundColor(style.DialogBgColor)
}

func (d *PodCreateDialog) setupResourceSettingsUI() {
	// security options
	labelWidth := 13
	fieldWidth := 15
	getSecColLabel := func(label string) string {
		return fmt.Sprintf("%17s: ", label)
	}

	// memory
	d.podMemoryField.SetLabel("memory:")
	d.podMemoryField.SetLabelWidth(labelWidth)
	d.podMemoryField.SetBackgroundColor(style.DialogBgColor)
	d.podMemoryField.SetLabelColor(style.DialogFgColor)
	d.podMemoryField.SetFieldBackgroundColor(style.InputFieldBgColor)
	d.podMemoryField.SetFieldWidth(fieldWidth)

	// cpus
	d.podCPUsField.SetLabel("cpus:")
	d.podCPUsField.SetLabelWidth(labelWidth)
	d.podCPUsField.SetBackgroundColor(style.DialogBgColor)
	d.podCPUsField.SetLabelColor(style.DialogFgColor)
	d.podCPUsField.SetFieldBackgroundColor(style.InputFieldBgColor)
	d.podCPUsField.SetFieldWidth(fieldWidth)

	// cpuset cpus
	d.podCPUSetCPUsField.SetLabel("cpuset cpus:")
	d.podCPUSetCPUsField.SetLabelWidth(labelWidth)
	d.podCPUSetCPUsField.SetBackgroundColor(style.DialogBgColor)
	d.podCPUSetCPUsField.SetLabelColor(style.DialogFgColor)
	d.podCPUSetCPUsField.SetFieldBackgroundColor(style.InputFieldBgColor)
	d.podCPUSetCPUsField.SetFieldWidth(fieldWidth)

	// shm size
	d.podShmSizeField.SetLabel("shm size:")
	d.podShmSizeField.SetLabelWidth(labelWidth)
	d.podShmSizeField.SetBackgroundColor(style.DialogBgColor)
	d.podShmSizeField.SetLabelColor(style.DialogFgColor)
	d.podShmSizeField.SetFieldBackgroundColor(style.InputFieldBgColor)
	d.podShmSizeField.SetFieldWidth(fieldWidth)

	// memory swap
	d.podMemorySwapField.SetLabel(getSecColLabel("memory swap"))
	d.podMemorySwapField.SetBackgroundColor(style.DialogBgColor)
	d.podMemorySwapField.SetLabelColor(style.DialogFgColor)
	d.podMemorySwapField.SetFieldBackgroundColor(style.InputFieldBgColor)
	d.podMemorySwapField.SetFieldWidth(fieldWidth)

	// cpu shares
	d.podCPUSharesField.SetLabel(getSecColLabel("cpu shares"))
	d.podCPUSharesField.SetBackgroundColor(style.DialogBgColor)
	d.podCPUSharesField.SetLabelColor(style.DialogFgColor)
	d.podCPUSharesField.SetFieldBackgroundColor(style.InputFieldBgColor)
	d.podCPUSharesField.SetFieldWidth(fieldWidth)

	// cpuset mems
	d.podCPUSetMemsField.SetLabel(getSecColLabel("cpuset mems"))
	d.podCPUSetMemsField.SetBackgroundColor(style.DialogBgColor)
	d.podCPUSetMemsField.SetLabelColor(style.DialogFgColor)
	d.podCPUSetMemsField.SetFieldBackgroundColor(style.InputFieldBgColor)
	d.podCPUSetMemsField.SetFieldWidth(fieldWidth)

	// shm size systemd
	d.podShmSizeSystemdField.SetLabel(getSecColLabel("shm size systemd"))
	d.podShmSizeSystemdField.SetBackgroundColor(style.DialogBgColor)
	d.podShmSizeSystemdField.SetLabelColor(style.DialogFgColor)
	d.podShmSizeSystemdField.SetFieldBackgroundColor(style.InputFieldBgColor)
	d.podShmSizeSystemdField.SetFieldWidth(fieldWidth)

	// layout
	row1 := tview.NewFlex().SetDirection(tview.FlexColumn)
	row1.SetBackgroundColor(style.DialogBgColor)
	row1.AddItem(d.podMemoryField, 0, 1, true)
	row1.AddItem(d.podMemorySwapField, 0, 1, true)

	row2 := tview.NewFlex().SetDirection(tview.FlexColumn)
	row2.SetBackgroundColor(style.DialogBgColor)
	row2.AddItem(d.podCPUsField, 0, 1, true)
	row2.AddItem(d.podCPUSharesField, 0, 1, true)

	row3 := tview.NewFlex().SetDirection(tview.FlexColumn)
	row3.SetBackgroundColor(style.DialogBgColor)
	row3.AddItem(d.podCPUSetCPUsField, 0, 1, true)
	row3.AddItem(d.podCPUSetMemsField, 0, 1, true)

	row4 := tview.NewFlex().SetDirection(tview.FlexColumn)
	row4.SetBackgroundColor(style.DialogBgColor)
	row4.AddItem(d.podShmSizeField, 0, 1, true)
	row4.AddItem(d.podShmSizeSystemdField, 0, 1, true)

	// resource page
	d.resourcePage.SetDirection(tview.FlexRow)
	d.resourcePage.AddItem(row1, 1, 0, true)
	d.resourcePage.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	d.resourcePage.AddItem(row2, 1, 0, true)
	d.resourcePage.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	d.resourcePage.AddItem(row3, 1, 0, true)
	d.resourcePage.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	d.resourcePage.AddItem(row4, 1, 0, true)
	d.resourcePage.SetBackgroundColor(style.DialogBgColor)
}

func (d *PodCreateDialog) setupNamespaceOptionsUI() {
	bgColor := style.DialogBgColor
	inputFieldBgColor := style.InputFieldBgColor
	namespacePageLabelWidth := 8

	// share
	d.podNamespaceShareField.SetLabel("share:")
	d.podNamespaceShareField.SetLabelWidth(namespacePageLabelWidth)
	d.podNamespaceShareField.SetBackgroundColor(bgColor)
	d.podNamespaceShareField.SetLabelColor(style.DialogFgColor)
	d.podNamespaceShareField.SetFieldBackgroundColor(inputFieldBgColor)

	// pid
	d.podNamespacePidField.SetLabel("pid:")
	d.podNamespacePidField.SetLabelWidth(namespacePageLabelWidth)
	d.podNamespacePidField.SetBackgroundColor(bgColor)
	d.podNamespacePidField.SetLabelColor(style.DialogFgColor)
	d.podNamespacePidField.SetFieldBackgroundColor(inputFieldBgColor)

	// userns
	d.podNamespaceUserField.SetLabel("userns:")
	d.podNamespaceUserField.SetLabelWidth(namespacePageLabelWidth)
	d.podNamespaceUserField.SetBackgroundColor(bgColor)
	d.podNamespaceUserField.SetLabelColor(style.DialogFgColor)
	d.podNamespaceUserField.SetFieldBackgroundColor(inputFieldBgColor)

	// uts
	d.podNamespaceUtsField.SetLabel("uts:")
	d.podNamespaceUtsField.SetLabelWidth(namespacePageLabelWidth)
	d.podNamespaceUtsField.SetBackgroundColor(bgColor)
	d.podNamespaceUtsField.SetLabelColor(style.DialogFgColor)
	d.podNamespaceUtsField.SetFieldBackgroundColor(inputFieldBgColor)

	// uidmap
	d.podNamespaceUidmapField.SetLabel("uidmap:")
	d.podNamespaceUidmapField.SetLabelWidth(namespacePageLabelWidth)
	d.podNamespaceUidmapField.SetBackgroundColor(bgColor)
	d.podNamespaceUidmapField.SetLabelColor(style.DialogFgColor)
	d.podNamespaceUidmapField.SetFieldBackgroundColor(inputFieldBgColor)

	// subuidname
	d.podNamespaceSubuidNameField.SetLabel("subuidname: ")
	d.podNamespaceSubuidNameField.SetBackgroundColor(bgColor)
	d.podNamespaceSubuidNameField.SetLabelColor(style.DialogFgColor)
	d.podNamespaceSubuidNameField.SetFieldBackgroundColor(inputFieldBgColor)

	// gidmap
	d.podNamespaceGidmapField.SetLabel("gidmap:")
	d.podNamespaceGidmapField.SetLabelWidth(namespacePageLabelWidth)
	d.podNamespaceGidmapField.SetBackgroundColor(bgColor)
	d.podNamespaceGidmapField.SetLabelColor(style.DialogFgColor)
	d.podNamespaceGidmapField.SetFieldBackgroundColor(inputFieldBgColor)

	// subgidname
	d.podNamespaceSubgidNameField.SetLabel("subgidname: ")
	d.podNamespaceSubgidNameField.SetBackgroundColor(bgColor)
	d.podNamespaceSubgidNameField.SetLabelColor(style.DialogFgColor)
	d.podNamespaceSubgidNameField.SetFieldBackgroundColor(inputFieldBgColor)

	// mapRow01Layout
	mapRow01Layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	mapRow01Layout.AddItem(d.podNamespaceUidmapField, 0, 1, true)
	mapRow01Layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	mapRow01Layout.AddItem(d.podNamespaceSubuidNameField, 0, 1, true)
	mapRow01Layout.SetBackgroundColor(bgColor)

	// mapRow02Layout
	mapRow02Layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	mapRow02Layout.AddItem(d.podNamespaceGidmapField, 0, 1, true)
	mapRow02Layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	mapRow02Layout.AddItem(d.podNamespaceSubgidNameField, 0, 1, true)
	mapRow02Layout.SetBackgroundColor(bgColor)

	// namespace options page
	d.namespacePage.SetDirection(tview.FlexRow)
	d.namespacePage.AddItem(d.podNamespaceShareField, 1, 0, true)
	d.namespacePage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.namespacePage.AddItem(d.podNamespacePidField, 1, 0, true)
	d.namespacePage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.namespacePage.AddItem(d.podNamespaceUserField, 1, 0, true)
	d.namespacePage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.namespacePage.AddItem(d.podNamespaceUtsField, 1, 0, true)
	d.namespacePage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.namespacePage.AddItem(mapRow01Layout, 1, 0, true)
	d.namespacePage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.namespacePage.AddItem(mapRow02Layout, 1, 0, true)

	d.namespacePage.SetBackgroundColor(bgColor)
}

// Display displays this primitive.
func (d *PodCreateDialog) Display() {
	d.display = true
	d.initData()
	d.focusElement = createPodCategoryPagesFocus
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
func (d *PodCreateDialog) Focus(delegate func(p tview.Primitive)) { //nolint:cyclop,gocyclo
	switch d.focusElement {
	// form has focus
	case createPodFormFocus:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = createPodCategoriesFocus // category text view
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
	case createPodCategoriesFocus:
		d.categories.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = createPodCategoryPagesFocus // category page view
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
	case createPodNoHostsCheckBoxFocus:
		delegate(d.podNoHostsCheckBox)
	case createPodLabelsFieldFocus:
		delegate(d.podLabelsField)
	// security options page
	case createPodSelinuxLabelFieldFocus:
		delegate(d.podSelinuxLabelField)
	case createPodApparmorFieldFocus:
		delegate(d.podApparmorField)
	case createPodSeccompFieldFocus:
		delegate(d.podSeccompField)
	case createPodMaskFieldFocus:
		delegate(d.podMaskField)
	case createPodUnmaskFieldFocus:
		delegate(d.podUnmaskField)
	case createPodNoNewPrivFieldFocus:
		delegate(d.podNoNewPrivField)
	// dns page
	case createPodDNSOptionsFieldFocus:
		delegate(d.podDNSOptionsField)
	case createPodDNSSearchDomaindFieldFocus:
		delegate(d.podDNSSearchDomaindField)
	// infra page
	case createPodInfraCommandFieldFocus:
		delegate(d.podInfraCommandField)
	case createPodInfraImageFieldFocus:
		delegate(d.podInfraImageField)
	// networking page
	case createPodIPAddressFieldFocus:
		delegate(d.podIPAddressField)
	case createPodMacAddressFieldFocus:
		delegate(d.podMacAddressField)
	case createPodAddHostFieldFocus:
		delegate(d.podAddHostField)
	case createPodNetworkFieldFocus:
		delegate(d.podNetworkField)
	case createPodPublishFieldFocus:
		delegate(d.podPublishField)
	// resource page
	case createPodMemoryFieldFocus:
		delegate(d.podMemoryField)
	case createPodMemorySwapFieldFocus:
		delegate(d.podMemorySwapField)
	case createPodCPUsFieldFocus:
		delegate(d.podCPUsField)
	case createPodCPUSharesFieldFocus:
		delegate(d.podCPUSharesField)
	case createPodCPUSetCPUsFieldFocus:
		delegate(d.podCPUSetCPUsField)
	case createPodCPUSetMemsFieldFocus:
		delegate(d.podCPUSetMemsField)
	case createPodShmSizeFieldFocus:
		delegate(d.podShmSizeField)
	case createPodShmSizeSystemdFieldFocus:
		delegate(d.podShmSizeSystemdField)
	// namespace page
	case createPodNamespaceShareFieldFocus:
		delegate(d.podNamespaceShareField)
	case createPodNamespacePidFieldFocus:
		delegate(d.podNamespacePidField)
	case createPodNamespaceUserFieldFocus:
		delegate(d.podNamespaceUserField)
	case createPodNamespaceUtsFieldFocus:
		delegate(d.podNamespaceUtsField)
	case createPodNamespaceUidmapFieldFocus:
		delegate(d.podNamespaceUidmapField)
	case createPodNamespaceSubuidNameFieldFocus:
		delegate(d.podNamespaceSubuidNameField)
	case createPodNamespaceGidmapFieldFocus:
		delegate(d.podNamespaceGidmapField)
	case createPodNamespaceSubgidNameFieldFocus:
		delegate(d.podNamespaceSubgidNameField)
	// category page
	case createPodCategoryPagesFocus:
		delegate(d.categoryPages)
	}
}

func (d *PodCreateDialog) initCustomInputHandlers() {
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

		if d.resourcePage.HasFocus() {
			if handler := d.resourcePage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setResourcePagePageNextFocus()
				}

				handler(event, setFocus)

				return
			}
		}

		if d.namespacePage.HasFocus() {
			if handler := d.namespacePage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setNamespacePageNextFocus()
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

	d.podMemoryField.SetText("")
	d.podMemorySwapField.SetText("")
	d.podCPUsField.SetText("")
	d.podCPUSharesField.SetText("")
	d.podCPUSetCPUsField.SetText("")
	d.podCPUSetMemsField.SetText("")
	d.podShmSizeField.SetText("")
	d.podShmSizeSystemdField.SetText("")

	// namespace
	d.podNamespaceShareField.SetText("ipc,net,uts")
	d.podNamespacePidField.SetText("")
	d.podNamespaceUserField.SetText("")
	d.podNamespaceUtsField.SetText("")
	d.podNamespaceUidmapField.SetText("")
	d.podNamespaceSubuidNameField.SetText("")
	d.podNamespaceGidmapField.SetText("")
	d.podNamespaceSubgidNameField.SetText("")
}

func (d *PodCreateDialog) setResourcePagePageNextFocus() {
	if d.podMemoryField.HasFocus() {
		d.focusElement = createPodMemorySwapFieldFocus

		return
	}

	if d.podMemorySwapField.HasFocus() {
		d.focusElement = createPodCPUsFieldFocus

		return
	}

	if d.podCPUsField.HasFocus() {
		d.focusElement = createPodCPUSharesFieldFocus

		return
	}

	if d.podCPUSharesField.HasFocus() {
		d.focusElement = createPodCPUSetCPUsFieldFocus

		return
	}

	if d.podCPUSetCPUsField.HasFocus() {
		d.focusElement = createPodCPUSetMemsFieldFocus

		return
	}

	if d.podCPUSetMemsField.HasFocus() {
		d.focusElement = createPodShmSizeFieldFocus

		return
	}

	if d.podShmSizeField.HasFocus() {
		d.focusElement = createPodShmSizeSystemdFieldFocus

		return
	}

	d.focusElement = createPodFormFocus
}

func (d *PodCreateDialog) setBasicInfoPageNextFocus() {
	if d.podNameField.HasFocus() {
		d.focusElement = createPodNoHostsCheckBoxFocus

		return
	}

	if d.podNoHostsCheckBox.HasFocus() {
		d.focusElement = createPodLabelsFieldFocus

		return
	}

	d.focusElement = createPodFormFocus
}

func (d *PodCreateDialog) setSecurityOptsPageNextFocus() {
	if d.podSelinuxLabelField.HasFocus() {
		d.focusElement = createPodApparmorFieldFocus

		return
	}

	if d.podApparmorField.HasFocus() {
		d.focusElement = createPodSeccompFieldFocus

		return
	}

	if d.podSeccompField.HasFocus() {
		d.focusElement = createPodMaskFieldFocus

		return
	}

	if d.podMaskField.HasFocus() {
		d.focusElement = createPodUnmaskFieldFocus

		return
	}

	if d.podUnmaskField.HasFocus() {
		d.focusElement = createPodNoNewPrivFieldFocus

		return
	}

	d.focusElement = createPodFormFocus
}

func (d *PodCreateDialog) setDNSSetupPageNextFocus() {
	if d.podDNSServerField.HasFocus() {
		d.focusElement = createPodDNSOptionsFieldFocus

		return
	}

	if d.podDNSOptionsField.HasFocus() {
		d.focusElement = createPodDNSSearchDomaindFieldFocus

		return
	}

	d.focusElement = createPodFormFocus
}

func (d *PodCreateDialog) setInfraSetupPageNextFocus() {
	if d.podInfraCheckBox.HasFocus() {
		d.focusElement = createPodInfraCommandFieldFocus

		return
	}

	if d.podInfraCommandField.HasFocus() {
		d.focusElement = createPodInfraImageFieldFocus

		return
	}

	d.focusElement = createPodFormFocus
}

func (d *PodCreateDialog) setNetworkingPageNextFocus() {
	if d.podHostnameField.HasFocus() {
		d.focusElement = createPodIPAddressFieldFocus

		return
	}

	if d.podIPAddressField.HasFocus() {
		d.focusElement = createPodMacAddressFieldFocus

		return
	}

	if d.podMacAddressField.HasFocus() {
		d.focusElement = createPodAddHostFieldFocus

		return
	}

	if d.podAddHostField.HasFocus() {
		d.focusElement = createPodNetworkFieldFocus

		return
	}

	if d.podNetworkField.HasFocus() {
		d.focusElement = createPodPublishFieldFocus

		return
	}

	d.focusElement = createPodFormFocus
}

func (d *PodCreateDialog) setNamespacePageNextFocus() {
	if d.podNamespaceShareField.HasFocus() {
		d.focusElement = createPodNamespacePidFieldFocus

		return
	}

	if d.podNamespacePidField.HasFocus() {
		d.focusElement = createPodNamespaceUserFieldFocus

		return
	}

	if d.podNamespaceUserField.HasFocus() {
		d.focusElement = createPodNamespaceUtsFieldFocus

		return
	}

	if d.podNamespaceUtsField.HasFocus() {
		d.focusElement = createPodNamespaceUidmapFieldFocus

		return
	}

	if d.podNamespaceUidmapField.HasFocus() {
		d.focusElement = createPodNamespaceSubuidNameFieldFocus

		return
	}

	if d.podNamespaceSubuidNameField.HasFocus() {
		d.focusElement = createPodNamespaceGidmapFieldFocus

		return
	}

	if d.podNamespaceGidmapField.HasFocus() {
		d.focusElement = createPodNamespaceSubgidNameFieldFocus

		return
	}

	d.focusElement = createPodFormFocus
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
		namespaceShare   []string
	)

	for _, nsshare := range strings.Split(d.podNamespaceShareField.GetText(), " ") {
		if nsshare != "" {
			namespaceShare = append(namespaceShare, nsshare)
		}
	}

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
		Name:                strings.TrimSpace(d.podNameField.GetText()),
		NoHost:              d.podNoHostsCheckBox.IsChecked(),
		Labels:              labels,
		DNSServer:           dnsServers,
		DNSOptions:          dnsOptions,
		DNSSearchDomain:     dnsSearchDomains,
		Infra:               d.podInfraCheckBox.IsChecked(),
		InfraImage:          strings.TrimSpace(d.podInfraImageField.GetText()),
		InfraCommand:        strings.TrimSpace(d.podInfraCommandField.GetText()),
		Hostname:            strings.TrimSpace(d.podHostnameField.GetText()),
		IPAddress:           strings.TrimSpace(d.podIPAddressField.GetText()),
		MacAddress:          strings.TrimSpace(d.podMacAddressField.GetText()),
		AddHost:             addHost,
		Network:             network,
		SecurityOpts:        securityOpts,
		Publish:             publish,
		Memory:              strings.TrimSpace(d.podMemoryField.GetText()),
		MemorySwap:          strings.TrimSpace(d.podMemorySwapField.GetText()),
		CPUs:                strings.TrimSpace(d.podCPUsField.GetText()),
		CPUShares:           strings.TrimSpace(d.podCPUSharesField.GetText()),
		CPUSetCPUs:          strings.TrimSpace(d.podCPUSetCPUsField.GetText()),
		CPUSetMems:          strings.TrimSpace(d.podCPUSetMemsField.GetText()),
		ShmSize:             strings.TrimSpace(d.podShmSizeField.GetText()),
		ShmSizeSystemd:      strings.TrimSpace(d.podShmSizeSystemdField.GetText()),
		NamespaceShare:      namespaceShare,
		NamespacePid:        strings.TrimSpace(d.podNamespacePidField.GetText()),
		NamespaceUser:       strings.TrimSpace(d.podNamespaceUserField.GetText()),
		NamespaceUts:        strings.TrimSpace(d.podNamespaceUtsField.GetText()),
		NamespaceUidmap:     strings.TrimSpace(d.podNamespaceUidmapField.GetText()),
		NamespaceSubuidName: strings.TrimSpace(d.podNamespaceSubuidNameField.GetText()),
		NamespaceGidmap:     strings.TrimSpace(d.podNamespaceGidmapField.GetText()),
		NamespaceSubgidName: strings.TrimSpace(d.podNamespaceSubgidNameField.GetText()),
	}

	return opts
}
