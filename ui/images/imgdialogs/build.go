package imgdialogs

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/containers/buildah/define"
	"github.com/containers/podman-tui/pdcs/images"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	buildDialogMaxWidth = 90
	buildDialogHeight   = 16
)

const (
	buildDialogFormFocus = 0 + iota
	buildDialogCategoriesFocus
	buildDialogCategoryPagesFocus
	buildDialogContainerfilePathFocus
	buildDialogPullPolicyFieldFocus
	buildDialogTagFieldFocus
	buildDialogRegistryFieldFocus
	buildDialogContextDirectoryPathFieldFocus
	buildDialogBuildArgsFieldFocus
	buildDialogLayersFieldFocus
	buildDialogNoCacheFieldfocus
	buildDialogSquashFieldFocus
	buildDialogFormatFieldFocus
	buildDialogLabelsFieldFocus
	buildDialogAnnotationFieldFocus
	buildDialogRemoveCntFieldFocus
	buildDialogForceRemoveCntFieldFocus
	buildDialogSelinuxLabelFieldFocus
	buildDialogApparmorProfileFieldFocus
	buildDialogSeccompProfilePathFieldFocus
	buildDialogNetworkFieldFocus
	buildDialogHTTPProxyFieldFocus
	buildDialogAddHostFieldFocus
	buildDialogDNSServersFieldFocus
	buildDialogDNSOptionsFieldFocus
	buildDialogDNSSearchFieldFocus
	buildDialogAddCapabilityFieldFocus
	buildDialogRemoveCapabilityFieldFocus
	buildDialogCPUPeriodFieldFocus
	buildDialogCPUQuataFieldFocus
	buildDialogCPUSharesFieldFocus
	buildDialogCPUSetCpusFieldFocus
	buildDialogCPUSetMemsFieldFocus
	buildDialogMemoryFieldFocus
	buildDialogMemorySwapFieldFocus
)

const (
	buildDialogBasicInfoPageIndex = 0 + iota
	buildDialogBuildInfoPageIndex
	buildDialogCapabilityPageIndex
	buildDialogCPUMemoryPageIndex
	buildDialogNetworkingPageIndex
	buildDialogSecurityOptsPageIndex
)

// ImageBuildDialog represents image build dialog primitive.
type ImageBuildDialog struct {
	*tview.Box

	layout                  *tview.Flex
	form                    *tview.Form
	categoryLabels          []string
	categories              *tview.TextView
	categoryPages           *tview.Pages
	basicInfoPage           *tview.Flex
	buildInfoPage           *tview.Flex
	securityOptsPage        *tview.Flex
	networkingPage          *tview.Flex
	capabilityPage          *tview.Flex
	cpuMemoryPage           *tview.Flex
	containerFilePath       *tview.InputField
	contextDirectoryPath    *tview.InputField
	tagField                *tview.InputField
	registryField           *tview.InputField
	pullPolicyField         *tview.DropDown
	formatField             *tview.DropDown
	buildArgsField          *tview.InputField
	layersField             *tview.Checkbox
	noCacheField            *tview.Checkbox
	SquashField             *tview.Checkbox
	labelsField             *tview.InputField
	removeCntField          *tview.Checkbox
	forceRemoveCntField     *tview.Checkbox
	annotationsField        *tview.InputField
	selinuxLabelField       *tview.InputField
	apparmorProfileField    *tview.InputField
	seccompProfilePathField *tview.InputField
	networkField            *tview.DropDown
	httpProxyField          *tview.Checkbox
	addHostField            *tview.InputField
	dnsServersField         *tview.InputField
	dnsOptionsField         *tview.InputField
	dnsSearchField          *tview.InputField
	addCapabilityField      *tview.InputField
	removeCapabilityField   *tview.InputField
	cpuPeriodField          *tview.InputField
	cpuQuataField           *tview.InputField
	cpuSharesField          *tview.InputField
	cpuSetCpusField         *tview.InputField
	cpuSetMemsField         *tview.InputField
	memoryField             *tview.InputField
	memorySwapField         *tview.InputField
	display                 bool
	focusElement            int
	activePageIndex         int
	cancelHandler           func()
	buildHandler            func()
}

// NewImageBuildDialog returns new image build dialog primitive.
func NewImageBuildDialog() *ImageBuildDialog { //nolint:maintidx
	buildDialog := &ImageBuildDialog{
		Box:    tview.NewBox(),
		layout: tview.NewFlex().SetDirection(tview.FlexRow),
		form:   tview.NewForm(),
		categoryLabels: []string{
			"Basic Information",
			"Build Settings",
			"Capability",
			"CPU and Memory",
			"Networking",
			"Security Options",
		},
		categories:              tview.NewTextView(),
		categoryPages:           tview.NewPages(),
		basicInfoPage:           tview.NewFlex(),
		buildInfoPage:           tview.NewFlex(),
		securityOptsPage:        tview.NewFlex(),
		networkingPage:          tview.NewFlex(),
		capabilityPage:          tview.NewFlex(),
		cpuMemoryPage:           tview.NewFlex(),
		containerFilePath:       tview.NewInputField(),
		contextDirectoryPath:    tview.NewInputField(),
		pullPolicyField:         tview.NewDropDown(),
		formatField:             tview.NewDropDown(),
		tagField:                tview.NewInputField(),
		registryField:           tview.NewInputField(),
		buildArgsField:          tview.NewInputField(),
		layersField:             tview.NewCheckbox(),
		noCacheField:            tview.NewCheckbox(),
		SquashField:             tview.NewCheckbox(),
		labelsField:             tview.NewInputField(),
		removeCntField:          tview.NewCheckbox(),
		forceRemoveCntField:     tview.NewCheckbox(),
		annotationsField:        tview.NewInputField(),
		selinuxLabelField:       tview.NewInputField(),
		apparmorProfileField:    tview.NewInputField(),
		seccompProfilePathField: tview.NewInputField(),
		networkField:            tview.NewDropDown(),
		httpProxyField:          tview.NewCheckbox(),
		addHostField:            tview.NewInputField(),
		dnsServersField:         tview.NewInputField(),
		dnsOptionsField:         tview.NewInputField(),
		dnsSearchField:          tview.NewInputField(),
		addCapabilityField:      tview.NewInputField(),
		removeCapabilityField:   tview.NewInputField(),
		cpuPeriodField:          tview.NewInputField(),
		cpuQuataField:           tview.NewInputField(),
		cpuSharesField:          tview.NewInputField(),
		cpuSetCpusField:         tview.NewInputField(),
		cpuSetMemsField:         tview.NewInputField(),
		memoryField:             tview.NewInputField(),
		memorySwapField:         tview.NewInputField(),
	}

	bgColor := style.DialogBgColor
	fgColor := style.DialogFgColor
	inputFieldBgColor := style.InputFieldBgColor
	ddUnselectedStyle := style.DropDownUnselected
	ddselectedStyle := style.DropDownSelected

	// categories
	buildDialog.categories.SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)
	buildDialog.categories.SetBackgroundColor(bgColor)
	buildDialog.categories.SetBorder(true)
	buildDialog.categories.SetBorderColor(style.DialogSubBoxBorderColor)

	// basic information setup page
	basicInfoPageLabelWidth := 17
	// context dir path field
	buildDialog.contextDirectoryPath.SetBackgroundColor(bgColor)
	buildDialog.contextDirectoryPath.SetLabel(utils.StringToInputLabel("context dir:", basicInfoPageLabelWidth))
	buildDialog.contextDirectoryPath.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.contextDirectoryPath.SetLabelStyle(style.InputLabelStyle)

	// Containerfile path field
	buildDialog.containerFilePath.SetBackgroundColor(bgColor)
	buildDialog.containerFilePath.SetLabel(utils.StringToInputLabel("container files:", basicInfoPageLabelWidth))
	buildDialog.containerFilePath.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.containerFilePath.SetLabelStyle(style.InputLabelStyle)

	// pull policy dropdown
	buildDialog.pullPolicyField.SetLabel("pull policy:")
	buildDialog.pullPolicyField.SetLabelWidth(basicInfoPageLabelWidth)
	buildDialog.pullPolicyField.SetBackgroundColor(bgColor)
	buildDialog.pullPolicyField.SetLabelColor(fgColor)
	buildDialog.pullPolicyField.SetOptions([]string{
		define.PullIfMissing.String(),
		define.PullAlways.String(),
		define.PullIfNewer.String(),
		define.PullNever.String(),
	},
		nil)
	buildDialog.pullPolicyField.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	buildDialog.pullPolicyField.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.pullPolicyField.SetFocusedStyle(style.DropDownFocused)

	// tag field
	buildDialog.tagField.SetBackgroundColor(bgColor)
	buildDialog.tagField.SetLabel(utils.StringToInputLabel("image tag:", basicInfoPageLabelWidth))
	buildDialog.tagField.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.tagField.SetLabelStyle(style.InputLabelStyle)

	// registry field
	buildDialog.registryField.SetBackgroundColor(bgColor)
	buildDialog.registryField.SetLabel(utils.StringToInputLabel("registry:", basicInfoPageLabelWidth))
	buildDialog.registryField.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.registryField.SetLabelStyle(style.InputLabelStyle)

	// build settings page
	buildSettingFirstColWidth := 15

	buildDialog.buildArgsField.SetBackgroundColor(bgColor)
	buildDialog.buildArgsField.SetLabel(utils.StringToInputLabel("runtime args:", buildSettingFirstColWidth))
	buildDialog.buildArgsField.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.buildArgsField.SetLabelStyle(style.InputLabelStyle)

	// format dropdown
	formatLabel := "output format:"

	buildDialog.formatField.SetLabel(formatLabel)
	buildDialog.formatField.SetTitleAlign(tview.AlignRight)
	buildDialog.formatField.SetLabelWidth(buildSettingFirstColWidth)
	buildDialog.formatField.SetBackgroundColor(bgColor)
	buildDialog.formatField.SetLabelColor(fgColor)
	buildDialog.formatField.SetOptions([]string{
		define.OCI,
		define.DOCKER,
	},
		nil)
	buildDialog.formatField.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	buildDialog.formatField.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.formatField.SetFocusedStyle(style.DropDownFocused)

	// squash
	squashLabel := "squash:"

	buildDialog.SquashField.SetBackgroundColor(bgColor)
	buildDialog.SquashField.SetBorder(false)
	buildDialog.SquashField.SetLabel(squashLabel)
	buildDialog.SquashField.SetLabelColor(fgColor)
	buildDialog.SquashField.SetLabelWidth(len(squashLabel) + 1)
	buildDialog.SquashField.SetFieldBackgroundColor(inputFieldBgColor)

	// layers
	layersLabel := "layers:"

	buildDialog.layersField.SetBackgroundColor(bgColor)
	buildDialog.layersField.SetBorder(false)
	buildDialog.layersField.SetLabel(layersLabel)
	buildDialog.layersField.SetLabelColor(fgColor)
	buildDialog.layersField.SetLabelWidth(len(layersLabel) + 1)
	buildDialog.layersField.SetFieldBackgroundColor(inputFieldBgColor)

	// no-cache
	noCacheLabel := "no cache:"

	buildDialog.noCacheField.SetBackgroundColor(bgColor)
	buildDialog.noCacheField.SetBorder(false)
	buildDialog.noCacheField.SetLabel(noCacheLabel)
	buildDialog.noCacheField.SetLabelColor(fgColor)
	buildDialog.noCacheField.SetLabelWidth(len(noCacheLabel) + 1)
	buildDialog.noCacheField.SetFieldBackgroundColor(inputFieldBgColor)

	// labels
	buildDialog.labelsField.SetBackgroundColor(bgColor)
	buildDialog.labelsField.SetLabel(utils.StringToInputLabel("labels:", buildSettingFirstColWidth))
	buildDialog.labelsField.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.labelsField.SetLabelStyle(style.InputLabelStyle)

	// annotations
	buildDialog.annotationsField.SetBackgroundColor(bgColor)
	buildDialog.annotationsField.SetLabel(utils.StringToInputLabel("annotations:", buildSettingFirstColWidth))
	buildDialog.annotationsField.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.annotationsField.SetLabelStyle(style.InputLabelStyle)

	// force remove field
	buildDialog.removeCntField.SetLabel("remove containers: ")
	buildDialog.removeCntField.SetBackgroundColor(bgColor)
	buildDialog.removeCntField.SetLabelColor(fgColor)
	buildDialog.removeCntField.SetFieldBackgroundColor(inputFieldBgColor)

	buildDialog.forceRemoveCntField.SetLabel("force remove: ")
	buildDialog.forceRemoveCntField.SetLabelWidth(buildSettingFirstColWidth)
	buildDialog.forceRemoveCntField.SetBackgroundColor(bgColor)
	buildDialog.forceRemoveCntField.SetLabelColor(fgColor)
	buildDialog.forceRemoveCntField.SetFieldBackgroundColor(inputFieldBgColor)

	// security options page
	securityOptionsPAgeLabelWidth := 10

	// selinux Label
	buildDialog.selinuxLabelField.SetBackgroundColor(bgColor)
	buildDialog.selinuxLabelField.SetLabel(utils.StringToInputLabel("label:", securityOptionsPAgeLabelWidth))
	buildDialog.selinuxLabelField.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.selinuxLabelField.SetLabelStyle(style.InputLabelStyle)

	// apparmor profile
	buildDialog.apparmorProfileField.SetBackgroundColor(bgColor)
	buildDialog.apparmorProfileField.SetLabel(utils.StringToInputLabel("apparmor:", securityOptionsPAgeLabelWidth))
	buildDialog.apparmorProfileField.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.apparmorProfileField.SetLabelStyle(style.InputLabelStyle)

	// seccomp profile
	buildDialog.seccompProfilePathField.SetBackgroundColor(bgColor)
	buildDialog.seccompProfilePathField.SetLabel(utils.StringToInputLabel("seccomp:", securityOptionsPAgeLabelWidth))
	buildDialog.seccompProfilePathField.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.seccompProfilePathField.SetLabelStyle(style.InputLabelStyle)

	// networking setup page
	networkingPageLabelWidth := 13

	// network dropdown
	buildDialog.networkField.SetLabel("network:")
	buildDialog.networkField.SetLabelWidth(networkingPageLabelWidth)
	buildDialog.networkField.SetBackgroundColor(bgColor)
	buildDialog.networkField.SetLabelColor(fgColor)
	buildDialog.networkField.SetOptions([]string{
		define.NetworkDefault.String(),
		define.NetworkDisabled.String(),
		define.NetworkEnabled.String(),
	},
		nil)
	buildDialog.networkField.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	buildDialog.networkField.SetFocusedStyle(style.DropDownFocused)
	buildDialog.networkField.SetFieldBackgroundColor(inputFieldBgColor)

	// http proxy checkbox
	buildDialog.httpProxyField.SetBackgroundColor(bgColor)
	buildDialog.httpProxyField.SetBorder(false)
	buildDialog.httpProxyField.SetLabel("http proxy:")
	buildDialog.httpProxyField.SetLabelColor(fgColor)
	buildDialog.httpProxyField.SetLabelWidth(networkingPageLabelWidth)
	buildDialog.httpProxyField.SetFieldBackgroundColor(inputFieldBgColor)

	// Add host field
	buildDialog.addHostField.SetBackgroundColor(bgColor)
	buildDialog.addHostField.SetLabel(utils.StringToInputLabel("add host:", networkingPageLabelWidth))
	buildDialog.addHostField.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.addHostField.SetLabelStyle(style.InputLabelStyle)

	// DNS servers field
	buildDialog.dnsServersField.SetBackgroundColor(bgColor)
	buildDialog.dnsServersField.SetLabel(utils.StringToInputLabel("dns servers:", networkingPageLabelWidth))
	buildDialog.dnsServersField.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.dnsServersField.SetLabelStyle(style.InputLabelStyle)

	// DNS options field
	buildDialog.dnsOptionsField.SetBackgroundColor(bgColor)
	buildDialog.dnsOptionsField.SetLabel(utils.StringToInputLabel("dns options:", networkingPageLabelWidth))
	buildDialog.dnsOptionsField.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.dnsOptionsField.SetLabelStyle(style.InputLabelStyle)

	// DNS search field
	buildDialog.dnsSearchField.SetBackgroundColor(bgColor)
	buildDialog.dnsSearchField.SetLabel(utils.StringToInputLabel("dns search:", networkingPageLabelWidth))
	buildDialog.dnsSearchField.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.dnsSearchField.SetLabelStyle(style.InputLabelStyle)

	// capability page
	capabilityPageLabelWidth := 12

	// add capability field
	buildDialog.addCapabilityField.SetBackgroundColor(bgColor)
	buildDialog.addCapabilityField.SetLabel(utils.StringToInputLabel("add cap:", capabilityPageLabelWidth))
	buildDialog.addCapabilityField.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.addCapabilityField.SetLabelStyle(style.InputLabelStyle)

	// remove capability field
	buildDialog.removeCapabilityField.SetBackgroundColor(bgColor)
	buildDialog.removeCapabilityField.SetLabel(utils.StringToInputLabel("remove cap:", capabilityPageLabelWidth))
	buildDialog.removeCapabilityField.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.removeCapabilityField.SetLabelStyle(style.InputLabelStyle)

	// cpu and memory page
	cpuMemoryLabelWidth := 14
	cpuMemoryFieldWidth := 17

	// cpu period field
	buildDialog.cpuPeriodField.SetBackgroundColor(bgColor)
	buildDialog.cpuPeriodField.SetLabel(utils.StringToInputLabel("cpu period:", cpuMemoryLabelWidth))
	buildDialog.cpuPeriodField.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.cpuPeriodField.SetLabelStyle(style.InputLabelStyle)
	buildDialog.cpuPeriodField.SetFieldWidth(cpuMemoryFieldWidth)

	// cpu quota field
	buildDialog.cpuQuataField.SetBackgroundColor(bgColor)
	buildDialog.cpuQuataField.SetLabel(utils.StringToInputLabel("cpu quota:", cpuMemoryLabelWidth))
	buildDialog.cpuQuataField.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.cpuQuataField.SetLabelStyle(style.InputLabelStyle)
	buildDialog.cpuQuataField.SetFieldWidth(cpuMemoryFieldWidth)

	// cpu shares field
	buildDialog.cpuSharesField.SetBackgroundColor(bgColor)
	buildDialog.cpuSharesField.SetLabel(utils.StringToInputLabel("cpu shares:", cpuMemoryLabelWidth))
	buildDialog.cpuSharesField.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.cpuSharesField.SetLabelStyle(style.InputLabelStyle)
	buildDialog.cpuSharesField.SetFieldWidth(cpuMemoryFieldWidth)

	// cpuset cpus field
	buildDialog.cpuSetCpusField.SetBackgroundColor(bgColor)
	buildDialog.cpuSetCpusField.SetLabel(utils.StringToInputLabel("cpu set cpus:", cpuMemoryLabelWidth))
	buildDialog.cpuSetCpusField.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.cpuSetCpusField.SetLabelStyle(style.InputLabelStyle)
	buildDialog.cpuSetCpusField.SetFieldWidth(cpuMemoryFieldWidth)

	// cpuset mems field
	buildDialog.cpuSetMemsField.SetBackgroundColor(bgColor)
	buildDialog.cpuSetMemsField.SetLabel(utils.StringToInputLabel(" cpu set mems:", cpuMemoryLabelWidth+1))
	buildDialog.cpuSetMemsField.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.cpuSetMemsField.SetLabelStyle(style.InputLabelStyle)
	buildDialog.cpuSetMemsField.SetFieldWidth(cpuMemoryFieldWidth)

	// memory field
	buildDialog.memoryField.SetBackgroundColor(bgColor)
	buildDialog.memoryField.SetLabel(utils.StringToInputLabel("memory:", cpuMemoryLabelWidth))
	buildDialog.memoryField.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.memoryField.SetLabelStyle(style.InputLabelStyle)
	buildDialog.memoryField.SetFieldWidth(cpuMemoryFieldWidth)

	// memory swap field
	buildDialog.memorySwapField.SetBackgroundColor(bgColor)
	buildDialog.memorySwapField.SetLabel(utils.StringToInputLabel(" memory swap:", cpuMemoryLabelWidth+1))
	buildDialog.memorySwapField.SetFieldBackgroundColor(inputFieldBgColor)
	buildDialog.memorySwapField.SetLabelStyle(style.InputLabelStyle)
	buildDialog.memorySwapField.SetFieldWidth(cpuMemoryFieldWidth)

	// category pages
	buildDialog.categoryPages.SetBackgroundColor(bgColor)
	buildDialog.categoryPages.SetBorder(true)
	buildDialog.categoryPages.SetBorderColor(style.DialogSubBoxBorderColor)

	// form
	buildDialog.form.SetBackgroundColor(bgColor)
	buildDialog.form.AddButton("Cancel", nil)
	buildDialog.form.AddButton("Build", nil)
	buildDialog.form.SetButtonsAlign(tview.AlignRight)
	buildDialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// layout
	buildDialog.setupLayout()
	buildDialog.layout.SetBackgroundColor(bgColor)
	buildDialog.layout.SetBorder(true)
	buildDialog.layout.SetBorderColor(style.DialogBorderColor)
	buildDialog.layout.SetTitle("PODMAN IMAGE BUILD")
	buildDialog.layout.AddItem(buildDialog.form, dialogs.DialogFormHeight, 0, true)

	return buildDialog
}

// Display displays this primitive.
func (d *ImageBuildDialog) Display() {
	d.focusElement = buildDialogContextDirectoryPathFieldFocus
	d.initData()
	d.display = true
}

// IsDisplay returns true if primitive is shown.
func (d *ImageBuildDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *ImageBuildDialog) Hide() {
	d.display = false
}

// HasFocus returns whether or not this primitive has focus.
func (d *ImageBuildDialog) HasFocus() bool {
	if d.categories.HasFocus() || d.categoryPages.HasFocus() {
		return true
	}

	if d.form.HasFocus() || d.layout.HasFocus() {
		return true
	}

	return d.Box.HasFocus()
}

// Focus is called when this primitive receives focus.
func (d *ImageBuildDialog) Focus(delegate func(p tview.Primitive)) { //nolint:gocyclo,cyclop
	switch d.focusElement {
	// form focus
	case buildDialogFormFocus:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == utils.SwitchFocusKey.Key {
				d.focusElement = buildDialogCategoriesFocus // category text view
				d.Focus(delegate)
				d.form.SetFocus(0)

				return nil
			}

			return event
		})

		delegate(d.form)
	// category text view
	case buildDialogCategoriesFocus:
		d.categories.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == utils.SwitchFocusKey.Key {
				d.focusElement = buildDialogCategoryPagesFocus // category page view
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
	case buildDialogContextDirectoryPathFieldFocus:
		delegate(d.contextDirectoryPath)
	case buildDialogContainerfilePathFocus:
		delegate(d.containerFilePath)
	case buildDialogPullPolicyFieldFocus:
		delegate(d.pullPolicyField)
	case buildDialogTagFieldFocus:
		delegate(d.tagField)
	case buildDialogRegistryFieldFocus:
		delegate(d.registryField)
	// build page
	case buildDialogBuildArgsFieldFocus:
		delegate(d.buildArgsField)
	case buildDialogLayersFieldFocus:
		delegate(d.layersField)
	case buildDialogNoCacheFieldfocus:
		delegate(d.noCacheField)
	case buildDialogFormatFieldFocus:
		delegate(d.formatField)
	case buildDialogSquashFieldFocus:
		delegate(d.SquashField)
	case buildDialogLabelsFieldFocus:
		delegate(d.labelsField)
	case buildDialogAnnotationFieldFocus:
		delegate(d.annotationsField)
	case buildDialogRemoveCntFieldFocus:
		delegate(d.removeCntField)
	case buildDialogForceRemoveCntFieldFocus:
		delegate(d.forceRemoveCntField)
	// security options page
	case buildDialogSelinuxLabelFieldFocus:
		delegate(d.selinuxLabelField)
	case buildDialogApparmorProfileFieldFocus:
		delegate(d.apparmorProfileField)
	case buildDialogSeccompProfilePathFieldFocus:
		delegate(d.seccompProfilePathField)
	// networking page
	case buildDialogNetworkFieldFocus:
		delegate(d.networkField)
	case buildDialogHTTPProxyFieldFocus:
		delegate(d.httpProxyField)
	case buildDialogAddHostFieldFocus:
		delegate(d.addHostField)
	case buildDialogDNSServersFieldFocus:
		delegate(d.dnsServersField)
	case buildDialogDNSOptionsFieldFocus:
		delegate(d.dnsOptionsField)
	case buildDialogDNSSearchFieldFocus:
		delegate(d.dnsSearchField)
	// capability page
	case buildDialogAddCapabilityFieldFocus:
		delegate(d.addCapabilityField)
	case buildDialogRemoveCapabilityFieldFocus:
		delegate(d.removeCapabilityField)
	// cpu and memory page
	case buildDialogCPUPeriodFieldFocus:
		delegate(d.cpuPeriodField)
	case buildDialogCPUQuataFieldFocus:
		delegate(d.cpuQuataField)
	case buildDialogCPUSharesFieldFocus:
		delegate(d.cpuSharesField)
	case buildDialogCPUSetCpusFieldFocus:
		delegate(d.cpuSetCpusField)
	case buildDialogCPUSetMemsFieldFocus:
		delegate(d.cpuSetMemsField)
	case buildDialogMemoryFieldFocus:
		delegate(d.memoryField)
	case buildDialogMemorySwapFieldFocus:
		delegate(d.memorySwapField)
	// category page
	case buildDialogCategoryPagesFocus:
		delegate(d.categoryPages)
	}
}

// InputHandler returns input handler function for this primitive.
func (d *ImageBuildDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) { //nolint:gocognit,gocyclo,lll,cyclop
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("image build dialog: event %v received", event)

		if event.Key() == utils.CloseDialogKey.Key {
			if !d.pullPolicyField.HasFocus() && !d.networkField.HasFocus() && !d.formatField.HasFocus() {
				d.cancelHandler()

				return
			}
		}
		// drop down event
		if d.pullPolicyField.HasFocus() {
			event = utils.ParseKeyEventKey(event)
			if handler := d.pullPolicyField.InputHandler(); handler != nil {
				if event.Key() == utils.SwitchFocusKey.Key {
					d.setBasicInfoPageNextFocus()
				}

				handler(event, setFocus)

				return
			}
		}

		if d.formatField.HasFocus() {
			event = utils.ParseKeyEventKey(event)
			if handler := d.formatField.InputHandler(); handler != nil {
				if event.Key() == utils.SwitchFocusKey.Key {
					d.setBuildSettingsPageNextFocus()
				}

				handler(event, setFocus)

				return
			}
		}

		if d.networkField.HasFocus() {
			event = utils.ParseKeyEventKey(event)
			if handler := d.networkField.InputHandler(); handler != nil {
				if event.Key() == utils.SwitchFocusKey.Key {
					d.setNetworkingPageNextFocus()
				}

				handler(event, setFocus)

				return
			}
		}

		// basic info page
		if d.basicInfoPage.HasFocus() {
			if handler := d.basicInfoPage.InputHandler(); handler != nil {
				if event.Key() == utils.SwitchFocusKey.Key {
					d.setBasicInfoPageNextFocus()
				}

				handler(event, setFocus)

				return
			}
		}

		// build settings page
		if d.buildInfoPage.HasFocus() {
			if handler := d.buildInfoPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setBuildSettingsPageNextFocus()
				}

				handler(event, setFocus)

				return
			}
		}

		// security options page
		if d.securityOptsPage.HasFocus() {
			if handler := d.securityOptsPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setSecurityOptionsPageNextFocus()
				}

				handler(event, setFocus)

				return
			}
		}

		// networking page
		if d.networkingPage.HasFocus() {
			if handler := d.networkingPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setNetworkingPageNextFocus()
				}

				handler(event, setFocus)

				return
			}
		}

		// capability page
		if d.capabilityPage.HasFocus() {
			if handler := d.capabilityPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setCapabilityPageNextFocus()
				}

				handler(event, setFocus)

				return
			}
		}

		// cpu and memory page
		if d.cpuMemoryPage.HasFocus() {
			if handler := d.cpuMemoryPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setCPUMemoryPageNextFocus()
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

		// form
		if d.form.HasFocus() {
			if formHandler := d.form.InputHandler(); formHandler != nil {
				formHandler(event, setFocus)

				return
			}
		}
	})
}

// SetRect set rects for this primitive.
func (d *ImageBuildDialog) SetRect(x, y, width, height int) {
	if width > buildDialogMaxWidth {
		emptySpace := (width - buildDialogMaxWidth) / 2 //nolint:mnd
		x += emptySpace
		width = buildDialogMaxWidth
	}

	if height > buildDialogHeight {
		emptySpace := (height - buildDialogHeight) / 2 //nolint:mnd
		y += emptySpace
		height = buildDialogHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen.
func (d *ImageBuildDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.DrawForSubclass(screen, d)
	x, y, width, height := d.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetCancelFunc sets form cancel button selected function.
func (d *ImageBuildDialog) SetCancelFunc(handler func()) *ImageBuildDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:mnd
	cancelButton.SetSelectedFunc(handler)

	return d
}

// SetBuildFunc sets form build button selected function.
func (d *ImageBuildDialog) SetBuildFunc(handler func()) *ImageBuildDialog {
	d.buildHandler = handler
	buildButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	buildButton.SetSelectedFunc(handler)

	return d
}

// ImageBuildOptions returns image build options.
func (d *ImageBuildDialog) ImageBuildOptions() (images.ImageBuildOptions, error) { //nolint:gocognit,gocyclo,cyclop,maintidx,lll
	var (
		memoryLimit        int64
		memorySwap         int64
		cpuPeriod          uint64
		cpuQuota           int64
		cpuShares          uint64
		cpuSetCpus         string
		cpuSetMems         string
		containerFiles     []string
		dnsServers         []string
		dnsOptions         []string
		dnsSearchDomains   []string
		addHost            []string
		labelOpts          []string
		apparmorProfile    string
		seccompProfilePath string
	)

	// basic info page
	// Containerfiles
	for _, cntFile := range strings.Split(strings.TrimSpace(d.containerFilePath.GetText()), " ") {
		if cntFile != "" {
			cFile, err := utils.ResolveHomeDir(cntFile)
			if err != nil {
				return images.ImageBuildOptions{}, err
			}

			containerFiles = append(containerFiles, cFile)
		}
	}

	opts := images.ImageBuildOptions{
		ContainerFiles: containerFiles,
		BuildOptions:   entities.BuildOptions{},
	}

	dir, err := utils.ResolveHomeDir(strings.TrimSpace(d.contextDirectoryPath.GetText()))
	if err != nil {
		return images.ImageBuildOptions{}, fmt.Errorf("cannot resolve home directory %w", err)
	}

	opts.BuildOptions.ContextDirectory = dir
	opts.BuildOptions.AdditionalTags = append(opts.BuildOptions.AdditionalTags, strings.TrimSpace(d.tagField.GetText()))
	opts.BuildOptions.Registry = strings.TrimSpace(d.registryField.GetText())

	_, pullOption := d.pullPolicyField.GetCurrentOption()
	switch pullOption {
	case "missing":
		opts.BuildOptions.PullPolicy = define.PullIfMissing
	case "always":
		opts.BuildOptions.PullPolicy = define.PullAlways
	case "ifnewer":
		opts.BuildOptions.PullPolicy = define.PullIfNewer
	case "never":
		opts.BuildOptions.PullPolicy = define.PullNever
	}

	_, format := d.formatField.GetCurrentOption()
	switch format {
	case "oci":
		opts.BuildOptions.OutputFormat = define.OCIv1ImageManifest
	case "docker":
		opts.BuildOptions.OutputFormat = define.Dockerv2ImageManifest
	}

	// build settings
	opts.BuildOptions.Squash = d.SquashField.IsChecked()
	opts.BuildOptions.Layers = d.layersField.IsChecked()
	opts.BuildOptions.NoCache = d.noCacheField.IsChecked()
	opts.BuildOptions.RemoveIntermediateCtrs = d.removeCntField.IsChecked()
	opts.BuildOptions.ForceRmIntermediateCtrs = d.forceRemoveCntField.IsChecked()

	labels := strings.TrimSpace(d.labelsField.GetText())
	if labels != "" {
		opts.BuildOptions.Labels = strings.Split(labels, " ")
	}

	annotations := strings.TrimSpace(d.annotationsField.GetText())
	if annotations != "" {
		opts.BuildOptions.Annotations = strings.Split(annotations, " ")
	}

	// capability pages
	addCap := strings.TrimSpace(d.addCapabilityField.GetText())
	if addCap != "" {
		opts.BuildOptions.AddCapabilities = strings.Split(addCap, " ")
	}

	removeCap := strings.TrimSpace(d.removeCapabilityField.GetText())
	if removeCap != "" {
		opts.BuildOptions.DropCapabilities = strings.Split(removeCap, " ")
	}

	// cpu and memory page
	opts.BuildOptions.CommonBuildOpts = &define.CommonBuildOptions{}

	cpuPeriodVal := d.cpuPeriodField.GetText()
	if cpuPeriodVal != "" {
		period, err := strconv.Atoi(cpuPeriodVal)
		if err != nil {
			return images.ImageBuildOptions{}, fmt.Errorf("invalid CPU period value %q %w", cpuPeriodVal, err)
		}

		cpuPeriod = uint64(period) //nolint:gosec
	}

	cpuQuotaVal := d.cpuQuataField.GetText()
	if cpuQuotaVal != "" {
		quota, err := strconv.Atoi(cpuQuotaVal)
		if err != nil {
			return images.ImageBuildOptions{}, fmt.Errorf("invalid CPU quota value %q %w", cpuQuotaVal, err)
		}

		cpuQuota = int64(quota)
	}

	cpuSharesVal := d.cpuSharesField.GetText()
	if cpuSharesVal != "" {
		shares, err := strconv.Atoi(cpuSharesVal)
		if err != nil {
			return images.ImageBuildOptions{}, fmt.Errorf("invalid CPU quota value %q %w", cpuSharesVal, err)
		}

		cpuShares = uint64(shares) //nolint:gosec
	}

	cpuSetCpusVal := d.cpuSetCpusField.GetText()
	if cpuSetCpusVal != "" {
		cpuSetCpus = cpuSetCpusVal
	}

	cpuSetMemsVal := d.cpuSetMemsField.GetText()
	if cpuSetMemsVal != "" {
		cpuSetMems = cpuSetMemsVal
	}

	memoryVal := d.memoryField.GetText()
	if memoryVal != "" {
		memory, err := strconv.Atoi(memoryVal)
		if err != nil {
			return images.ImageBuildOptions{}, fmt.Errorf("invalid memory value %q %w", memoryVal, err)
		}

		memoryLimit = int64(memory)
	}

	memorySwapVal := d.memorySwapField.GetText()
	if memorySwapVal != "" {
		swap, err := strconv.Atoi(memorySwapVal)
		if err != nil {
			return images.ImageBuildOptions{}, fmt.Errorf("invalid memory swap value %q %w", memorySwapVal, err)
		}

		memorySwap = int64(swap)
	}

	// networking page
	// network policy
	_, configureNetwork := d.networkField.GetCurrentOption()
	switch configureNetwork {
	case "NetworkDefault":
		opts.BuildOptions.ConfigureNetwork = define.NetworkDefault
	case "NetworkDisabled":
		opts.BuildOptions.ConfigureNetwork = define.NetworkDisabled
	case "NetworkEnabled":
		opts.BuildOptions.ConfigureNetwork = define.NetworkEnabled
	}
	// add hosts
	hosts := strings.TrimSpace(d.addHostField.GetText())
	if hosts != "" {
		addHost = strings.Split(hosts, " ")
	}

	// dns page
	dnsServersList := strings.TrimSpace(d.dnsServersField.GetText())
	for _, dns := range strings.Split(dnsServersList, " ") {
		if dns != "" {
			dnsServers = append(dnsServers, dns)
		}
	}

	for _, do := range strings.Split(d.dnsOptionsField.GetText(), " ") {
		if do != "" {
			dnsOptions = append(dnsOptions, do)
		}
	}

	for _, ds := range strings.Split(d.dnsSearchField.GetText(), " ") {
		if ds != "" {
			dnsSearchDomains = append(dnsSearchDomains, ds)
		}
	}

	// security options page
	for _, selinuxLabel := range strings.Split(d.selinuxLabelField.GetText(), " ") {
		if selinuxLabel != "" {
			labelOpts = append(labelOpts, selinuxLabel)
		}
	}

	apparmor := strings.TrimSpace(d.apparmorProfileField.GetText())
	if apparmor != "" {
		apparmorProfile = apparmor
	}

	seccomp := strings.TrimSpace(d.seccompProfilePathField.GetText())
	if seccomp != "" {
		seccompProfilePath = seccomp
	}

	commonOpts := &define.CommonBuildOptions{
		AddHost:            addHost,
		HTTPProxy:          d.httpProxyField.IsChecked(),
		CPUPeriod:          cpuPeriod,
		CPUQuota:           cpuQuota,
		CPUSetCPUs:         cpuSetCpus,
		CPUSetMems:         cpuSetMems,
		CPUShares:          cpuShares,
		DNSServers:         dnsServers,
		DNSOptions:         dnsOptions,
		DNSSearch:          dnsSearchDomains,
		Memory:             memoryLimit,
		MemorySwap:         memorySwap,
		LabelOpts:          labelOpts,
		ApparmorProfile:    apparmorProfile,
		SeccompProfilePath: seccompProfilePath,
	}

	opts.BuildOptions.CommonBuildOpts = commonOpts

	return opts, nil
}

func (d *ImageBuildDialog) setupLayout() {
	bgColor := style.DialogBgColor
	// basic info page
	d.basicInfoPage.SetDirection(tview.FlexRow)
	d.basicInfoPage.AddItem(d.contextDirectoryPath, 1, 0, true)
	d.basicInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.basicInfoPage.AddItem(d.containerFilePath, 1, 0, true)
	d.basicInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.basicInfoPage.AddItem(d.pullPolicyField, 1, 0, true)
	d.basicInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.basicInfoPage.AddItem(d.tagField, 1, 0, true)
	d.basicInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.basicInfoPage.AddItem(d.registryField, 1, 0, true)
	d.basicInfoPage.SetBackgroundColor(bgColor)

	// layers setup page
	secondRowLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	secondRowLayout.SetBackgroundColor(bgColor)
	secondRowLayout.AddItem(d.formatField, 0, 2, true) //nolint:mnd
	secondRowLayout.AddItem(d.SquashField, 0, 1, true)
	secondRowLayout.AddItem(d.layersField, 0, 1, true)
	secondRowLayout.AddItem(d.noCacheField, 0, 1, true)

	cntRmRowLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	cntRmRowLayout.SetBackgroundColor(bgColor)
	cntRmRowLayout.AddItem(d.forceRemoveCntField, 0, 1, true)
	cntRmRowLayout.AddItem(d.removeCntField, 0, 2, true) //nolint:mnd

	// build setup page
	d.buildInfoPage.SetDirection(tview.FlexRow)
	d.buildInfoPage.AddItem(d.buildArgsField, 1, 0, true)
	d.buildInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.buildInfoPage.AddItem(secondRowLayout, 1, 0, true)
	d.buildInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.buildInfoPage.AddItem(d.labelsField, 1, 0, true)
	d.buildInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.buildInfoPage.AddItem(d.annotationsField, 1, 0, true)
	d.buildInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.buildInfoPage.AddItem(cntRmRowLayout, 1, 0, true)
	d.buildInfoPage.SetBackgroundColor(bgColor)

	// security options page
	d.securityOptsPage.SetDirection(tview.FlexRow)
	d.securityOptsPage.AddItem(d.selinuxLabelField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.apparmorProfileField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.seccompProfilePathField, 1, 0, true)

	// networking setup page
	netFirstRow := tview.NewFlex().SetDirection(tview.FlexColumn)
	netFirstRow.SetBackgroundColor(bgColor)
	netFirstRow.AddItem(d.networkField, 0, 1, true)
	netFirstRow.AddItem(d.httpProxyField, 0, 1, true)
	d.networkingPage.SetDirection(tview.FlexRow)
	d.networkingPage.AddItem(netFirstRow, 1, 0, true)
	d.networkingPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.networkingPage.AddItem(d.addHostField, 1, 0, true)
	d.networkingPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.networkingPage.AddItem(d.dnsServersField, 1, 0, true)
	d.networkingPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.networkingPage.AddItem(d.dnsOptionsField, 1, 0, true)
	d.networkingPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.networkingPage.AddItem(d.dnsSearchField, 1, 0, true)
	d.networkingPage.SetBackgroundColor(bgColor)

	// capability page
	d.capabilityPage.SetDirection(tview.FlexRow)
	d.capabilityPage.AddItem(d.addCapabilityField, 1, 0, true)
	d.capabilityPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.capabilityPage.AddItem(d.removeCapabilityField, 1, 0, true)

	// cpu and memory page
	cpuSetRow := tview.NewFlex().SetDirection(tview.FlexColumn)
	cpuSetRow.AddItem(d.cpuSetCpusField, 0, 1, true)
	cpuSetRow.AddItem(d.cpuSetMemsField, 0, 1, true)

	// memory and swap
	memSwapRow := tview.NewFlex().SetDirection(tview.FlexColumn)
	memSwapRow.AddItem(d.memoryField, 0, 1, true)
	memSwapRow.AddItem(d.memorySwapField, 0, 1, true)

	d.cpuMemoryPage.SetDirection(tview.FlexRow)
	d.cpuMemoryPage.AddItem(d.cpuPeriodField, 0, 1, true)
	d.cpuMemoryPage.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, true)
	d.cpuMemoryPage.AddItem(d.cpuQuataField, 0, 1, true)
	d.cpuMemoryPage.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, true)
	d.cpuMemoryPage.AddItem(d.cpuSharesField, 0, 1, true)
	d.cpuMemoryPage.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, true)
	d.cpuMemoryPage.AddItem(cpuSetRow, 0, 1, true)
	d.cpuMemoryPage.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, true)
	d.cpuMemoryPage.AddItem(memSwapRow, 0, 1, true)

	// adding category pages
	d.categoryPages.AddPage(d.categoryLabels[buildDialogBasicInfoPageIndex], d.basicInfoPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[buildDialogBuildInfoPageIndex], d.buildInfoPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[buildDialogCapabilityPageIndex], d.capabilityPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[buildDialogCPUMemoryPageIndex], d.cpuMemoryPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[buildDialogNetworkingPageIndex], d.networkingPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[buildDialogSecurityOptsPageIndex], d.securityOptsPage, true, true)

	// add it to layout.
	_, layoutWidth := utils.AlignStringListWidth(d.categoryLabels)
	layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	layout.AddItem(d.categories, layoutWidth+6, 0, true) //nolint:mnd
	layout.AddItem(d.categoryPages, 0, 1, true)
	layout.SetBackgroundColor(bgColor)

	d.layout.AddItem(layout, 0, 1, true)
}

func (d *ImageBuildDialog) initData() {
	d.setActiveCategory(0)

	// basic info page
	d.containerFilePath.SetText("")
	d.contextDirectoryPath.SetText("")
	d.tagField.SetText("")
	d.registryField.SetText("")
	d.pullPolicyField.SetCurrentOption(0)

	// build page
	d.buildArgsField.SetText("")
	d.layersField.SetChecked(false)
	d.noCacheField.SetChecked(false)
	d.formatField.SetCurrentOption(0)
	d.SquashField.SetChecked(false)
	d.labelsField.SetText("")
	d.annotationsField.SetText("")
	d.removeCntField.SetChecked(false)
	d.forceRemoveCntField.SetChecked(false)

	// security options page
	d.selinuxLabelField.SetText("")
	d.apparmorProfileField.SetText("")
	d.seccompProfilePathField.SetText("")

	// networking setting page
	d.networkField.SetCurrentOption(0)
	d.httpProxyField.SetChecked(false)
	d.addHostField.SetText("")
	d.dnsServersField.SetText("")
	d.dnsOptionsField.SetText("")
	d.dnsSearchField.SetText("")

	// capability setting page
	d.addCapabilityField.SetText("")
	d.removeCapabilityField.SetText("")

	// memory and cpu page
	d.cpuPeriodField.SetText("")
	d.cpuQuataField.SetText("")
	d.cpuSharesField.SetText("")
	d.cpuSetCpusField.SetText("")
	d.cpuSetMemsField.SetText("")
	d.memoryField.SetText("")
	d.memorySwapField.SetText("")
}

func (d *ImageBuildDialog) setActiveCategory(index int) {
	bgColor := style.ButtonBgColor
	fgColor := style.DialogFgColor

	ctgFgColor := style.GetColorHex(fgColor)
	ctgBgColor := style.GetColorHex(bgColor)
	d.activePageIndex = index
	d.categories.Clear()

	ctgList := []string{}

	alignedList, _ := utils.AlignStringListWidth(d.categoryLabels)

	for i := range alignedList {
		if i == index {
			ctgList = append(ctgList, fmt.Sprintf("[%s:%s:b]-> %s ", ctgFgColor, ctgBgColor, alignedList[i]))

			continue
		}

		ctgList = append(ctgList, fmt.Sprintf("[-:-:-]   %s ", alignedList[i]))
	}

	d.categories.SetText(strings.Join(ctgList, "\n"))

	// switch the page
	d.categoryPages.SwitchToPage(d.categoryLabels[index])
}

func (d *ImageBuildDialog) nextCategory() {
	activePage := d.activePageIndex
	if d.activePageIndex < len(d.categoryLabels)-1 {
		activePage++
		d.setActiveCategory(activePage)

		return
	}

	d.setActiveCategory(0)
}

func (d *ImageBuildDialog) previousCategory() {
	activePage := d.activePageIndex
	if d.activePageIndex > 0 {
		activePage--
		d.setActiveCategory(activePage)

		return
	}

	d.setActiveCategory(len(d.categoryLabels) - 1)
}

func (d *ImageBuildDialog) setBasicInfoPageNextFocus() {
	if d.contextDirectoryPath.HasFocus() {
		d.focusElement = buildDialogContainerfilePathFocus

		return
	}

	if d.containerFilePath.HasFocus() {
		d.focusElement = buildDialogPullPolicyFieldFocus

		return
	}

	if d.pullPolicyField.HasFocus() {
		d.focusElement = buildDialogTagFieldFocus

		return
	}

	if d.tagField.HasFocus() {
		d.focusElement = buildDialogRegistryFieldFocus

		return
	}

	d.focusElement = buildDialogFormFocus
}

func (d *ImageBuildDialog) setNetworkingPageNextFocus() {
	if d.networkField.HasFocus() {
		d.focusElement = buildDialogHTTPProxyFieldFocus

		return
	}

	if d.httpProxyField.HasFocus() {
		d.focusElement = buildDialogAddHostFieldFocus

		return
	}

	if d.addHostField.HasFocus() {
		d.focusElement = buildDialogDNSServersFieldFocus

		return
	}

	if d.dnsServersField.HasFocus() {
		d.focusElement = buildDialogDNSOptionsFieldFocus

		return
	}

	if d.dnsOptionsField.HasFocus() {
		d.focusElement = buildDialogDNSSearchFieldFocus

		return
	}

	d.focusElement = buildDialogFormFocus
}

func (d *ImageBuildDialog) setBuildSettingsPageNextFocus() {
	if d.buildArgsField.HasFocus() {
		d.focusElement = buildDialogFormatFieldFocus

		return
	}

	if d.formatField.HasFocus() {
		d.focusElement = buildDialogSquashFieldFocus

		return
	}

	if d.SquashField.HasFocus() {
		d.focusElement = buildDialogLayersFieldFocus

		return
	}

	if d.layersField.HasFocus() {
		d.focusElement = buildDialogNoCacheFieldfocus

		return
	}

	if d.noCacheField.HasFocus() {
		d.focusElement = buildDialogLabelsFieldFocus

		return
	}

	if d.labelsField.HasFocus() {
		d.focusElement = buildDialogAnnotationFieldFocus

		return
	}

	if d.annotationsField.HasFocus() {
		d.focusElement = buildDialogForceRemoveCntFieldFocus

		return
	}

	if d.forceRemoveCntField.HasFocus() {
		d.focusElement = buildDialogRemoveCntFieldFocus

		return
	}

	d.focusElement = buildDialogFormFocus
}

func (d *ImageBuildDialog) setCapabilityPageNextFocus() {
	if d.addCapabilityField.HasFocus() {
		d.focusElement = buildDialogRemoveCapabilityFieldFocus
	} else {
		d.focusElement = buildDialogFormFocus
	}
}

func (d *ImageBuildDialog) setCPUMemoryPageNextFocus() {
	if d.cpuPeriodField.HasFocus() {
		d.focusElement = buildDialogCPUQuataFieldFocus

		return
	}

	if d.cpuQuataField.HasFocus() {
		d.focusElement = buildDialogCPUSharesFieldFocus

		return
	}

	if d.cpuSharesField.HasFocus() {
		d.focusElement = buildDialogCPUSetCpusFieldFocus

		return
	}

	if d.cpuSetCpusField.HasFocus() {
		d.focusElement = buildDialogCPUSetMemsFieldFocus

		return
	}

	if d.cpuSetMemsField.HasFocus() {
		d.focusElement = buildDialogMemoryFieldFocus

		return
	}

	if d.memoryField.HasFocus() {
		d.focusElement = buildDialogMemorySwapFieldFocus

		return
	}

	d.focusElement = buildDialogFormFocus
}

func (d *ImageBuildDialog) setSecurityOptionsPageNextFocus() {
	if d.selinuxLabelField.HasFocus() {
		d.focusElement = buildDialogApparmorProfileFieldFocus

		return
	}

	if d.apparmorProfileField.HasFocus() {
		d.focusElement = buildDialogSeccompProfilePathFieldFocus

		return
	}

	d.focusElement = buildDialogFormFocus
}
