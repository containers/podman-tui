package cntdialogs

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/containers"
	"github.com/containers/podman-tui/pdcs/images"
	"github.com/containers/podman-tui/pdcs/networks"
	"github.com/containers/podman-tui/pdcs/pods"
	"github.com/containers/podman-tui/pdcs/volumes"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	containerCreateDialogMaxWidth     = 100
	containerCreateOnlyDialogHeight   = 20
	containerCreateAndRunDialogHeight = 22
)

const (
	ContainerCreateOnlyDialogMode = 0 + iota
	ContainerCreateAndRunDialogMode
)

const (
	createContainerFormFocus = 0 + iota
	createCategoriesFocus
	createCategoryPagesFocus
	createContainerNameFieldFocus
	createContainerCommandFieldFocus
	createContainerImageFieldFocus
	createcontainerPodFieldFocus
	createContainerLabelsFieldFocus
	createContainerRemoveFieldFocus
	createContainerPrivilegedFieldFocus
	createContainerTimeoutFieldFocus
	createContainerInteractiveFieldFocus
	createContainerDetachFieldFocus
	createContainerTtyFieldFocus
	createContainerSecretFieldFocus
	createContainerEnvHostFieldFocus
	createContainerEnvVarsFieldFocus
	createContainerEnvFileFieldFocus
	createContainerEnvMergeFieldFocus
	createContainerWorkDirFieldFocus
	createContainerUmaskFieldFocus
	createContainerUnsetEnvFieldFocus
	createContainerUnsetEnvAllFieldFocus
	createContainerUserFieldFocus
	createContainerHostUsersFieldFocus
	createContainerPasswdEntryFieldFocus
	createContainerGroupEntryFieldFocus
	createcontainerSecLabelFieldFocus
	createContainerApprarmorFieldFocus
	createContainerSeccompFeildFocus
	createcontainerSecMaskFieldFocus
	createcontainerSecUnmaskFieldFocus
	createcontainerSecNoNewPrivFieldFocus
	createContainerPortExposeFieldFocus
	createContainerPortPublishFieldFocus
	createContainerPortPublishAllFieldFocus
	createContainerHostnameFieldFocus
	createContainerIPAddrFieldFocus
	createContainerMacAddrFieldFocus
	createContainerNetworkFieldFocus
	createContainerDNSServersFieldFocus
	createContainerDNSOptionsFieldFocus
	createContainerDNSSearchFieldFocus
	createContainerImageVolumeFieldFocus
	createContainerVolumeFieldFocus
	createContainerMountFieldFocus
	createContainerHealthCmdFieldFocus
	createContainerHealthStartupCmdFieldFocus
	createContainerHealthOnFailureFieldFocus
	createContainerHealthIntervalFieldFocus
	createContainerHealthStartupIntervalFieldFocus
	createContainerHealthTimeoutFieldFocus
	createContainerHealthStartupTimeoutFieldFocus
	createContainerHealthRetriesFieldFocus
	createContainerHealthStartupRetriesFieldFocus
	createContainerHealthStartPeriodFieldFocus
	createContainerHealthStartupSuccessFieldFocus
	createContainerHealthLogDestFocus
	createContainerHealthMaxLogCountFocus
	createContainerHealthMaxLogSizeFocus
	createContainerMemoryFieldFocus
	createContainerMemoryReservatoinFieldFocus
	createContainerMemorySwapFieldFocus
	createcontainerMemorySwappinessFieldFocus
	createContainerCPUsFieldFocus
	createContainerCPUSharesFieldFocus
	createContainerCPUPeriodFieldFocus
	createContainerCPURtPeriodFieldFocus
	createContainerCPUQuotaFieldFocus
	createContainerCPURtRuntimeFeildFocus
	createContainerCPUSetCPUsFieldFocus
	createContainerCPUSetMemsFieldFocus
	createContainerShmSizeFieldFocus
	createContainerShmSizeSystemdFieldFocus
	createContainerNamespaceCgroupFieldFocus
	createContainerNamespaceIpcFieldFocus
	createContainerNamespacePidFieldFocus
	createContainerNamespaceUserFieldFocus
	createContainerNamespaceUtsFieldFocus
	createContainerNamespaceUidmapFieldFocus
	createContainerNamespaceGidmapFieldFocus
	createContainerNamespaceSubuidNameFieldFocus
	createContainerNamespaceSubgidNameFieldFocus
)

const (
	createContainerInfoPageIndex = 0 + iota
	createContainerEnvironmentPageIndex
	createContainerUserGroupsPageIndex
	createContainerDNSPageIndex
	createContainerHealthPageIndex
	createContainerNetworkingPageIndex
	createContainerPortPageIndex
	createContainerSecurityOptsPageIndex
	createContainerVolumePageIndex
	createContainerResourcePageIndex
	createContainerNamespacePageIndex
)

type ContainerCreateDialogMode int

// ContainerCreateDialog implements container create dialog.
type ContainerCreateDialog struct {
	*tview.Box
	mode                                ContainerCreateDialogMode
	layout                              *tview.Flex
	categoryLabels                      []string
	categories                          *tview.TextView
	categoryPages                       *tview.Pages
	containerInfoPage                   *tview.Flex
	environmentPage                     *tview.Flex
	userGroupsPage                      *tview.Flex
	securityOptsPage                    *tview.Flex
	portPage                            *tview.Flex
	networkingPage                      *tview.Flex
	dnsPage                             *tview.Flex
	volumePage                          *tview.Flex
	healthPage                          *tview.Flex
	resourcePage                        *tview.Flex
	namespacePage                       *tview.Flex
	form                                *tview.Form
	display                             bool
	activePageIndex                     int
	focusElement                        int
	imageList                           []images.ImageListReporter
	podList                             []*entities.ListPodsReport
	containerNameField                  *tview.InputField
	containerCommandField               *tview.InputField
	containerImageField                 *tview.DropDown
	containerPodField                   *tview.DropDown
	containerLabelsField                *tview.InputField
	containerRemoveField                *tview.Checkbox
	containerPrivilegedField            *tview.Checkbox
	containerTimeoutField               *tview.InputField
	containerInteractiveField           *tview.Checkbox
	containerTtyField                   *tview.Checkbox
	containerDetachField                *tview.Checkbox
	containerSecretField                *tview.InputField
	containerWorkDirField               *tview.InputField
	containerEnvHostField               *tview.Checkbox
	containerEnvVarsField               *tview.InputField
	containerEnvFileField               *tview.InputField
	containerEnvMergeField              *tview.InputField
	containerUmaskField                 *tview.InputField
	containerUnsetEnvField              *tview.InputField
	containerUnsetEnvAllField           *tview.Checkbox
	containerUserField                  *tview.InputField
	containerHostUsersField             *tview.InputField
	containerPasswdEntryField           *tview.InputField
	containerGroupEntryField            *tview.InputField
	containerSecLabelField              *tview.InputField
	containerSecApparmorField           *tview.InputField
	containerSeccompField               *tview.InputField
	containerSecMaskField               *tview.InputField
	containerSecUnmaskField             *tview.InputField
	containerSecNoNewPrivField          *tview.Checkbox
	containerPortExposeField            *tview.InputField
	containerPortPublishField           *tview.InputField
	ContainerPortPublishAllField        *tview.Checkbox
	containerHostnameField              *tview.InputField
	containerIPAddrField                *tview.InputField
	containerMacAddrField               *tview.InputField
	containerNetworkField               *tview.DropDown
	containerDNSServersField            *tview.InputField
	containerDNSOptionsField            *tview.InputField
	containerDNSSearchField             *tview.InputField
	containerHealthCmdField             *tview.InputField
	containerHealthIntervalField        *tview.InputField
	containerHealthOnFailureField       *tview.DropDown
	containerHealthRetriesField         *tview.InputField
	containerHealthStartPeriodField     *tview.InputField
	containerHealthTimeoutField         *tview.InputField
	containerHealthStartupCmdField      *tview.InputField
	containerHealthStartupIntervalField *tview.InputField
	containerHealthStartupRetriesField  *tview.InputField
	containerHealthStartupSuccessField  *tview.InputField
	containerHealthStartupTimeoutField  *tview.InputField
	containerHealthLogDestField         *tview.InputField
	containerHealthMaxLogCountField     *tview.InputField
	containerHealthMaxLogSizeField      *tview.InputField
	containerVolumeField                *tview.InputField
	containerImageVolumeField           *tview.DropDown
	containerMountField                 *tview.InputField
	containerMemoryField                *tview.InputField
	containerMemoryReservationField     *tview.InputField
	containerMemorySwapField            *tview.InputField
	containerMemorySwappinessField      *tview.InputField
	containerCPUsField                  *tview.InputField
	containerCPUPeriodField             *tview.InputField
	containerCPUQuotaField              *tview.InputField
	containerCPURtPeriodField           *tview.InputField
	containerCPURtRuntimeField          *tview.InputField
	containerCPUSharesField             *tview.InputField
	containerCPUSetCPUsField            *tview.InputField
	containerCPUSetMemsField            *tview.InputField
	containerShmSizeField               *tview.InputField
	containerShmSizeSystemdField        *tview.InputField
	containerNamespaceCgroupField       *tview.InputField
	containerNamespaceIpcField          *tview.InputField
	containerNamespacePidField          *tview.InputField
	containerNamespaceUserField         *tview.InputField
	containerNamespaceUtsField          *tview.InputField
	containerNamespaceUidmapField       *tview.InputField
	containerNamespaceGidmapField       *tview.InputField
	containerNamespaceSubuidNameField   *tview.InputField
	containerNamespaceSubgidNameField   *tview.InputField
	cancelHandler                       func()
	enterHandler                        func()
}

// NewContainerCreateDialog returns new container create dialog primitive ContainerCreateDialog.
func NewContainerCreateDialog(mode ContainerCreateDialogMode) *ContainerCreateDialog {
	containerDialog := ContainerCreateDialog{
		Box:               tview.NewBox(),
		mode:              mode,
		layout:            tview.NewFlex().SetDirection(tview.FlexRow),
		categories:        tview.NewTextView(),
		categoryPages:     tview.NewPages(),
		containerInfoPage: tview.NewFlex(),
		environmentPage:   tview.NewFlex(),
		userGroupsPage:    tview.NewFlex(),
		securityOptsPage:  tview.NewFlex(),
		networkingPage:    tview.NewFlex(),
		dnsPage:           tview.NewFlex(),
		portPage:          tview.NewFlex(),
		volumePage:        tview.NewFlex(),
		healthPage:        tview.NewFlex(),
		resourcePage:      tview.NewFlex(),
		namespacePage:     tview.NewFlex(),
		form:              tview.NewForm(),
		categoryLabels: []string{
			"Container",
			"Environment",
			"User and groups",
			"DNS Settings",
			"Health check",
			"Network Settings",
			"Ports Settings",
			"Security Options",
			"Volumes Settings",
			"Resource Settings",
			"Namespace Options",
		},
		activePageIndex:                     0,
		display:                             false,
		containerNameField:                  tview.NewInputField(),
		containerCommandField:               tview.NewInputField(),
		containerImageField:                 tview.NewDropDown(),
		containerPodField:                   tview.NewDropDown(),
		containerLabelsField:                tview.NewInputField(),
		containerRemoveField:                tview.NewCheckbox(),
		containerPrivilegedField:            tview.NewCheckbox(),
		containerTimeoutField:               tview.NewInputField(),
		containerInteractiveField:           tview.NewCheckbox(),
		containerTtyField:                   tview.NewCheckbox(),
		containerDetachField:                tview.NewCheckbox(),
		containerSecretField:                tview.NewInputField(),
		containerWorkDirField:               tview.NewInputField(),
		containerEnvHostField:               tview.NewCheckbox(),
		containerEnvVarsField:               tview.NewInputField(),
		containerEnvFileField:               tview.NewInputField(),
		containerEnvMergeField:              tview.NewInputField(),
		containerUmaskField:                 tview.NewInputField(),
		containerUnsetEnvField:              tview.NewInputField(),
		containerUnsetEnvAllField:           tview.NewCheckbox(),
		containerUserField:                  tview.NewInputField(),
		containerHostUsersField:             tview.NewInputField(),
		containerPasswdEntryField:           tview.NewInputField(),
		containerGroupEntryField:            tview.NewInputField(),
		containerSecLabelField:              tview.NewInputField(),
		containerSecApparmorField:           tview.NewInputField(),
		containerSeccompField:               tview.NewInputField(),
		containerSecMaskField:               tview.NewInputField(),
		containerSecUnmaskField:             tview.NewInputField(),
		containerSecNoNewPrivField:          tview.NewCheckbox(),
		containerPortExposeField:            tview.NewInputField(),
		containerPortPublishField:           tview.NewInputField(),
		ContainerPortPublishAllField:        tview.NewCheckbox(),
		containerHostnameField:              tview.NewInputField(),
		containerIPAddrField:                tview.NewInputField(),
		containerMacAddrField:               tview.NewInputField(),
		containerNetworkField:               tview.NewDropDown(),
		containerDNSServersField:            tview.NewInputField(),
		containerDNSOptionsField:            tview.NewInputField(),
		containerDNSSearchField:             tview.NewInputField(),
		containerVolumeField:                tview.NewInputField(),
		containerImageVolumeField:           tview.NewDropDown(),
		containerMountField:                 tview.NewInputField(),
		containerHealthCmdField:             tview.NewInputField(),
		containerHealthIntervalField:        tview.NewInputField(),
		containerHealthOnFailureField:       tview.NewDropDown(),
		containerHealthRetriesField:         tview.NewInputField(),
		containerHealthStartPeriodField:     tview.NewInputField(),
		containerHealthTimeoutField:         tview.NewInputField(),
		containerHealthStartupCmdField:      tview.NewInputField(),
		containerHealthStartupIntervalField: tview.NewInputField(),
		containerHealthStartupRetriesField:  tview.NewInputField(),
		containerHealthStartupSuccessField:  tview.NewInputField(),
		containerHealthStartupTimeoutField:  tview.NewInputField(),
		containerHealthLogDestField:         tview.NewInputField(),
		containerHealthMaxLogCountField:     tview.NewInputField(),
		containerHealthMaxLogSizeField:      tview.NewInputField(),
		containerMemoryField:                tview.NewInputField(),
		containerMemoryReservationField:     tview.NewInputField(),
		containerMemorySwapField:            tview.NewInputField(),
		containerMemorySwappinessField:      tview.NewInputField(),
		containerCPUsField:                  tview.NewInputField(),
		containerCPUPeriodField:             tview.NewInputField(),
		containerCPUQuotaField:              tview.NewInputField(),
		containerCPURtPeriodField:           tview.NewInputField(),
		containerCPURtRuntimeField:          tview.NewInputField(),
		containerCPUSharesField:             tview.NewInputField(),
		containerCPUSetCPUsField:            tview.NewInputField(),
		containerCPUSetMemsField:            tview.NewInputField(),
		containerShmSizeField:               tview.NewInputField(),
		containerShmSizeSystemdField:        tview.NewInputField(),
		containerNamespaceCgroupField:       tview.NewInputField(),
		containerNamespacePidField:          tview.NewInputField(),
		containerNamespaceIpcField:          tview.NewInputField(),
		containerNamespaceUserField:         tview.NewInputField(),
		containerNamespaceUtsField:          tview.NewInputField(),
		containerNamespaceUidmapField:       tview.NewInputField(),
		containerNamespaceGidmapField:       tview.NewInputField(),
		containerNamespaceSubuidNameField:   tview.NewInputField(),
		containerNamespaceSubgidNameField:   tview.NewInputField(),
	}

	containerDialog.setupLayout()
	containerDialog.setActiveCategory(0)
	containerDialog.initCustomInputHanlers()

	return &containerDialog
}

func (d *ContainerCreateDialog) setupLayout() {
	bgColor := style.DialogBgColor

	d.categories.SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)
	d.categories.SetBackgroundColor(bgColor)
	d.categories.SetBorder(true)
	d.categories.SetBorderColor(style.DialogSubBoxBorderColor)

	// category pages
	d.categoryPages.SetBackgroundColor(bgColor)
	d.categoryPages.SetBorder(true)
	d.categoryPages.SetBorderColor(style.DialogSubBoxBorderColor)

	d.setupContainerInfoPageUI()
	d.setupEnvironmentPageUI()
	d.setupUserGroupsPageUI()
	d.setupDNSPageUI()
	d.setupHealthPageUI()
	d.setupResourcePageUI()
	d.setupNetworkPageUI()
	d.setupPortsPageUI()
	d.setupSecurityPageUI()
	d.setupVolumePageUI()
	d.setupNamespacePageUI()

	// form
	d.form.SetBackgroundColor(bgColor)
	d.form.AddButton("Cancel", nil)

	if d.mode == ContainerCreateOnlyDialogMode {
		d.form.AddButton("Create", nil)
	} else {
		d.form.AddButton("Run", nil)
	}

	d.form.SetButtonsAlign(tview.AlignRight)
	d.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// adding category pages
	d.categoryPages.AddPage(d.categoryLabels[createContainerInfoPageIndex], d.containerInfoPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[createContainerEnvironmentPageIndex], d.environmentPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[createContainerUserGroupsPageIndex], d.userGroupsPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[createContainerDNSPageIndex], d.dnsPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[createContainerHealthPageIndex], d.healthPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[createContainerNetworkingPageIndex], d.networkingPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[createContainerPortPageIndex], d.portPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[createContainerSecurityOptsPageIndex], d.securityOptsPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[createContainerVolumePageIndex], d.volumePage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[createContainerResourcePageIndex], d.resourcePage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[createContainerNamespacePageIndex], d.namespacePage, true, true)

	// add it to layout.
	d.layout.SetBackgroundColor(bgColor)
	d.layout.SetBorder(true)
	d.layout.SetBorderColor(style.DialogBorderColor)

	if d.mode == ContainerCreateOnlyDialogMode {
		d.layout.SetTitle("PODMAN CONTAINER CREATE")
	} else {
		d.layout.SetTitle("PODMAN CONTAINER RUN")
	}

	_, layoutWidth := utils.AlignStringListWidth(d.categoryLabels)
	layout := tview.NewFlex().SetDirection(tview.FlexColumn)

	layout.AddItem(d.categories, layoutWidth+6, 0, true) //nolint:mnd
	layout.AddItem(d.categoryPages, 0, 1, true)
	layout.SetBackgroundColor(bgColor)
	d.layout.AddItem(layout, 0, 1, true)

	d.layout.AddItem(d.form, dialogs.DialogFormHeight, 0, true)
}

func (d *ContainerCreateDialog) setupContainerInfoPageUI() {
	bgColor := style.DialogBgColor
	ddUnselectedStyle := style.DropDownUnselected
	ddselectedStyle := style.DropDownSelected
	inputFieldBgColor := style.InputFieldBgColor
	cntInfoPageLabelWidth := 12

	// name field
	d.containerNameField.SetLabel("name:")
	d.containerNameField.SetLabelWidth(cntInfoPageLabelWidth)
	d.containerNameField.SetBackgroundColor(bgColor)
	d.containerNameField.SetLabelColor(style.DialogFgColor)
	d.containerNameField.SetFieldBackgroundColor(inputFieldBgColor)

	// command field
	d.containerCommandField.SetLabel("command:")
	d.containerCommandField.SetLabelWidth(cntInfoPageLabelWidth)
	d.containerCommandField.SetBackgroundColor(bgColor)
	d.containerCommandField.SetLabelColor(style.DialogFgColor)
	d.containerCommandField.SetFieldBackgroundColor(inputFieldBgColor)

	// image field
	d.containerImageField.SetLabel("image:")
	d.containerImageField.SetLabelWidth(cntInfoPageLabelWidth)
	d.containerImageField.SetBackgroundColor(bgColor)
	d.containerImageField.SetLabelColor(style.DialogFgColor)
	d.containerImageField.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	d.containerImageField.SetFieldBackgroundColor(inputFieldBgColor)

	// pod field
	d.containerPodField.SetLabel("pod:")
	d.containerPodField.SetLabelWidth(cntInfoPageLabelWidth)
	d.containerPodField.SetBackgroundColor(bgColor)
	d.containerPodField.SetLabelColor(style.DialogFgColor)
	d.containerPodField.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	d.containerPodField.SetFieldBackgroundColor(inputFieldBgColor)

	// labels field
	d.containerLabelsField.SetLabel("labels:")
	d.containerLabelsField.SetLabelWidth(cntInfoPageLabelWidth)
	d.containerLabelsField.SetBackgroundColor(bgColor)
	d.containerLabelsField.SetLabelColor(style.DialogFgColor)
	d.containerLabelsField.SetFieldBackgroundColor(inputFieldBgColor)

	// privileged
	d.containerPrivilegedField.SetLabel("privileged:")
	d.containerPrivilegedField.SetLabelWidth(cntInfoPageLabelWidth)
	d.containerPrivilegedField.SetBackgroundColor(bgColor)
	d.containerPrivilegedField.SetLabelColor(style.DialogFgColor)
	d.containerPrivilegedField.SetFieldBackgroundColor(inputFieldBgColor)

	// timeout field
	timeoutLabel := "timeout:"

	d.containerTimeoutField.SetLabel(timeoutLabel)
	d.containerTimeoutField.SetLabelWidth(len(timeoutLabel) + 1)
	d.containerTimeoutField.SetBackgroundColor(bgColor)
	d.containerTimeoutField.SetLabelColor(style.DialogFgColor)
	d.containerTimeoutField.SetFieldBackgroundColor(inputFieldBgColor)

	// interactive
	interactiveLabel := "interactive:"
	d.containerInteractiveField.SetLabel(interactiveLabel)
	d.containerInteractiveField.SetLabelWidth(len(interactiveLabel) + 1)
	d.containerInteractiveField.SetBackgroundColor(bgColor)
	d.containerInteractiveField.SetLabelColor(style.DialogFgColor)
	d.containerInteractiveField.SetFieldBackgroundColor(inputFieldBgColor)

	// detach
	d.containerDetachField.SetLabel("detach:")
	d.containerDetachField.SetLabelWidth(cntInfoPageLabelWidth)
	d.containerDetachField.SetBackgroundColor(bgColor)
	d.containerDetachField.SetLabelColor(style.DialogFgColor)
	d.containerDetachField.SetFieldBackgroundColor(inputFieldBgColor)

	// tty
	ttyLabel := fmt.Sprintf("%7s:", "tty")
	d.containerTtyField.SetLabel(ttyLabel)
	d.containerTtyField.SetLabelWidth(len(timeoutLabel) + 1)
	d.containerTtyField.SetBackgroundColor(bgColor)
	d.containerTtyField.SetLabelColor(style.DialogFgColor)
	d.containerTtyField.SetFieldBackgroundColor(inputFieldBgColor)

	// remove field
	removeLabel := fmt.Sprintf("%11s:", "remove")
	d.containerRemoveField.SetLabel(removeLabel)
	d.containerRemoveField.SetLabelWidth(len(interactiveLabel) + 1)
	d.containerRemoveField.SetBackgroundColor(bgColor)
	d.containerRemoveField.SetLabelColor(style.DialogFgColor)
	d.containerRemoveField.SetFieldBackgroundColor(inputFieldBgColor)

	// secrets
	d.containerSecretField.SetLabel("secret:")
	d.containerSecretField.SetLabelWidth(cntInfoPageLabelWidth)
	d.containerSecretField.SetBackgroundColor(bgColor)
	d.containerSecretField.SetLabelColor(style.DialogFgColor)
	d.containerSecretField.SetFieldBackgroundColor(inputFieldBgColor)

	// layout
	labelPaddings := 4
	checkBoxLayout1 := tview.NewFlex().SetDirection(tview.FlexColumn)

	checkBoxLayout1.SetBackgroundColor(bgColor)
	checkBoxLayout1.AddItem(d.containerPrivilegedField, cntInfoPageLabelWidth+labelPaddings, 0, false)
	checkBoxLayout1.AddItem(d.containerRemoveField, 0, 1, false)
	checkBoxLayout1.AddItem(d.containerTimeoutField, 0, 1, false)
	checkBoxLayout1.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, true)

	checkBoxLayout2 := tview.NewFlex().SetDirection(tview.FlexColumn)

	checkBoxLayout2.SetBackgroundColor(bgColor)
	checkBoxLayout2.AddItem(d.containerDetachField, cntInfoPageLabelWidth+labelPaddings, 0, false)
	checkBoxLayout2.AddItem(d.containerInteractiveField, 0, 1, false)
	checkBoxLayout2.AddItem(d.containerTtyField, 0, 1, false)
	checkBoxLayout2.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, true)

	d.containerInfoPage.SetDirection(tview.FlexRow)
	d.containerInfoPage.AddItem(d.containerNameField, 1, 0, true)
	d.containerInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.containerInfoPage.AddItem(d.containerCommandField, 1, 0, true)
	d.containerInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.containerInfoPage.AddItem(d.containerImageField, 1, 0, true)
	d.containerInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.containerInfoPage.AddItem(d.containerPodField, 1, 0, true)
	d.containerInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.containerInfoPage.AddItem(d.containerLabelsField, 1, 0, true)
	d.containerInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.containerInfoPage.AddItem(checkBoxLayout1, 1, 0, true)
	d.containerInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)

	if d.mode == ContainerCreateAndRunDialogMode {
		d.containerInfoPage.AddItem(checkBoxLayout2, 1, 0, true)
		d.containerInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	}

	d.containerInfoPage.AddItem(d.containerSecretField, 1, 0, true)
	d.containerInfoPage.SetBackgroundColor(bgColor)
}

func (d *ContainerCreateDialog) setupEnvironmentPageUI() {
	bgColor := style.DialogBgColor
	inputFieldBgColor := style.InputFieldBgColor
	envPageLabelWidth := 12

	// environment host
	d.containerEnvHostField.SetLabel("env host:")
	d.containerEnvHostField.SetLabelWidth(envPageLabelWidth)
	d.containerEnvHostField.SetBackgroundColor(bgColor)
	d.containerEnvHostField.SetLabelColor(style.DialogFgColor)
	d.containerEnvHostField.SetFieldBackgroundColor(inputFieldBgColor)

	// unset all
	unsetEnvAllLabel := "unsetenv all"
	d.containerUnsetEnvAllField.SetLabel(unsetEnvAllLabel)
	d.containerUnsetEnvAllField.SetLabelWidth(len(unsetEnvAllLabel) + 1)
	d.containerUnsetEnvAllField.SetBackgroundColor(bgColor)
	d.containerUnsetEnvAllField.SetLabelColor(style.DialogFgColor)
	d.containerUnsetEnvAllField.SetFieldBackgroundColor(inputFieldBgColor)

	// environment variables
	d.containerEnvVarsField.SetLabel("env vars:")
	d.containerEnvVarsField.SetLabelWidth(envPageLabelWidth)
	d.containerEnvVarsField.SetBackgroundColor(bgColor)
	d.containerEnvVarsField.SetLabelColor(style.DialogFgColor)
	d.containerEnvVarsField.SetFieldBackgroundColor(inputFieldBgColor)

	// environment variables file
	d.containerEnvFileField.SetLabel("env file:")
	d.containerEnvFileField.SetLabelWidth(envPageLabelWidth)
	d.containerEnvFileField.SetBackgroundColor(bgColor)
	d.containerEnvFileField.SetLabelColor(style.DialogFgColor)
	d.containerEnvFileField.SetFieldBackgroundColor(inputFieldBgColor)

	// environment merge
	d.containerEnvMergeField.SetLabel("env merge:")
	d.containerEnvMergeField.SetLabelWidth(envPageLabelWidth)
	d.containerEnvMergeField.SetBackgroundColor(bgColor)
	d.containerEnvMergeField.SetLabelColor(style.DialogFgColor)
	d.containerEnvMergeField.SetFieldBackgroundColor(inputFieldBgColor)

	// environment unset variables
	d.containerUnsetEnvField.SetLabel("unset env:")
	d.containerUnsetEnvField.SetLabelWidth(envPageLabelWidth)
	d.containerUnsetEnvField.SetBackgroundColor(bgColor)
	d.containerUnsetEnvField.SetLabelColor(style.DialogFgColor)
	d.containerUnsetEnvField.SetFieldBackgroundColor(inputFieldBgColor)

	// working directory
	d.containerWorkDirField.SetLabel("work dir:")
	d.containerWorkDirField.SetLabelWidth(envPageLabelWidth)
	d.containerWorkDirField.SetBackgroundColor(bgColor)
	d.containerWorkDirField.SetLabelColor(style.DialogFgColor)
	d.containerWorkDirField.SetFieldBackgroundColor(inputFieldBgColor)

	// umask
	umaskLabel := "umask:"
	d.containerUmaskField.SetLabel(umaskLabel)
	d.containerUmaskField.SetLabelWidth(len(umaskLabel) + 1)
	d.containerUmaskField.SetBackgroundColor(bgColor)
	d.containerUmaskField.SetLabelColor(style.DialogFgColor)
	d.containerUmaskField.SetFieldBackgroundColor(inputFieldBgColor)

	// layout
	labelPaddings := 4
	checkBoxLayout := tview.NewFlex().SetDirection(tview.FlexColumn)

	checkBoxLayout.SetBackgroundColor(bgColor)
	checkBoxLayout.AddItem(d.containerEnvHostField, envPageLabelWidth+labelPaddings, 0, false)
	checkBoxLayout.AddItem(d.containerUnsetEnvAllField, len(unsetEnvAllLabel)+labelPaddings, 0, false)
	checkBoxLayout.AddItem(d.containerUmaskField, 0, 1, false)
	checkBoxLayout.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, true)

	d.environmentPage.SetDirection(tview.FlexRow)
	d.environmentPage.AddItem(d.containerWorkDirField, 1, 0, true)
	d.environmentPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.environmentPage.AddItem(d.containerEnvVarsField, 1, 0, true)
	d.environmentPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.environmentPage.AddItem(d.containerEnvFileField, 1, 0, true)
	d.environmentPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.environmentPage.AddItem(d.containerEnvMergeField, 1, 0, true)
	d.environmentPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.environmentPage.AddItem(d.containerUnsetEnvField, 1, 0, true)
	d.environmentPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.environmentPage.AddItem(checkBoxLayout, 1, 0, true)
	d.environmentPage.SetBackgroundColor(bgColor)
}

func (d *ContainerCreateDialog) setupUserGroupsPageUI() {
	bgColor := style.DialogBgColor
	inputFieldBgColor := style.InputFieldBgColor
	userGroupLabelWidth := 14
	userFieldWidth := 30

	// user
	d.containerUserField.SetLabel("user:")
	d.containerUserField.SetLabelWidth(userGroupLabelWidth)
	d.containerUserField.SetBackgroundColor(bgColor)
	d.containerUserField.SetLabelColor(style.DialogFgColor)
	d.containerUserField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerUserField.SetFieldWidth(userFieldWidth)

	// host users
	d.containerHostUsersField.SetLabel("host user:")
	d.containerHostUsersField.SetLabelWidth(userGroupLabelWidth)
	d.containerHostUsersField.SetBackgroundColor(bgColor)
	d.containerHostUsersField.SetLabelColor(style.DialogFgColor)
	d.containerHostUsersField.SetFieldBackgroundColor(inputFieldBgColor)

	// passwd entry
	d.containerPasswdEntryField.SetLabel("passwd entry:")
	d.containerPasswdEntryField.SetLabelWidth(userGroupLabelWidth)
	d.containerPasswdEntryField.SetBackgroundColor(bgColor)
	d.containerPasswdEntryField.SetLabelColor(style.DialogFgColor)
	d.containerPasswdEntryField.SetFieldBackgroundColor(inputFieldBgColor)

	// group entry
	d.containerGroupEntryField.SetLabel("group entry:")
	d.containerGroupEntryField.SetLabelWidth(userGroupLabelWidth)
	d.containerGroupEntryField.SetBackgroundColor(bgColor)
	d.containerGroupEntryField.SetLabelColor(style.DialogFgColor)
	d.containerGroupEntryField.SetFieldBackgroundColor(inputFieldBgColor)

	d.userGroupsPage.SetDirection(tview.FlexRow)
	d.userGroupsPage.AddItem(d.containerUserField, 1, 0, true)
	d.userGroupsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.userGroupsPage.AddItem(d.containerHostUsersField, 1, 0, true)
	d.userGroupsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.userGroupsPage.AddItem(d.containerPasswdEntryField, 1, 0, true)
	d.userGroupsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.userGroupsPage.AddItem(d.containerGroupEntryField, 1, 0, true)
	d.userGroupsPage.SetBackgroundColor(bgColor)
}

func (d *ContainerCreateDialog) setupDNSPageUI() {
	bgColor := style.DialogBgColor
	inputFieldBgColor := style.InputFieldBgColor
	dnsPageLabelWidth := 13

	// hostname field
	d.containerDNSServersField.SetLabel("dns servers:")
	d.containerDNSServersField.SetLabelWidth(dnsPageLabelWidth)
	d.containerDNSServersField.SetBackgroundColor(bgColor)
	d.containerDNSServersField.SetLabelColor(style.DialogFgColor)
	d.containerDNSServersField.SetFieldBackgroundColor(inputFieldBgColor)

	// IP field
	d.containerDNSOptionsField.SetLabel("dns options:")
	d.containerDNSOptionsField.SetLabelWidth(dnsPageLabelWidth)
	d.containerDNSOptionsField.SetBackgroundColor(bgColor)
	d.containerDNSOptionsField.SetLabelColor(style.DialogFgColor)
	d.containerDNSOptionsField.SetFieldBackgroundColor(inputFieldBgColor)

	// mac field
	d.containerDNSSearchField.SetLabel("dns search:")
	d.containerDNSSearchField.SetLabelWidth(dnsPageLabelWidth)
	d.containerDNSSearchField.SetBackgroundColor(bgColor)
	d.containerDNSSearchField.SetLabelColor(style.DialogFgColor)
	d.containerDNSSearchField.SetFieldBackgroundColor(inputFieldBgColor)

	d.dnsPage.SetDirection(tview.FlexRow)
	d.dnsPage.AddItem(d.containerDNSServersField, 1, 0, true)
	d.dnsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.dnsPage.AddItem(d.containerDNSOptionsField, 1, 0, true)
	d.dnsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.dnsPage.AddItem(d.containerDNSSearchField, 1, 0, true)
	d.dnsPage.SetBackgroundColor(bgColor)
}

func (d *ContainerCreateDialog) setupHealthPageUI() {
	bgColor := style.DialogBgColor
	inputFieldBgColor := style.InputFieldBgColor
	ddUnselectedStyle := style.DropDownUnselected
	ddselectedStyle := style.DropDownSelected
	healthPageLabelWidth := 13
	healthPageSecColLabelWidth := 18
	healthPageMultiRowFieldWidth := 7

	// health cmd
	d.containerHealthCmdField.SetLabel("Command:")
	d.containerHealthCmdField.SetLabelWidth(healthPageLabelWidth)
	d.containerHealthCmdField.SetBackgroundColor(bgColor)
	d.containerHealthCmdField.SetLabelColor(style.DialogFgColor)
	d.containerHealthCmdField.SetFieldBackgroundColor(inputFieldBgColor)

	// startup cmd
	d.containerHealthStartupCmdField.SetLabel("Startup cmd:")
	d.containerHealthStartupCmdField.SetLabelWidth(healthPageLabelWidth)
	d.containerHealthStartupCmdField.SetBackgroundColor(bgColor)
	d.containerHealthStartupCmdField.SetLabelColor(style.DialogFgColor)
	d.containerHealthStartupCmdField.SetFieldBackgroundColor(inputFieldBgColor)

	d.containerHealthLogDestField.SetLabel("Log dest:")
	d.containerHealthLogDestField.SetLabelWidth(healthPageLabelWidth)
	d.containerHealthLogDestField.SetBackgroundColor(bgColor)
	d.containerHealthLogDestField.SetLabelColor(style.DialogFgColor)
	d.containerHealthLogDestField.SetFieldBackgroundColor(inputFieldBgColor)

	// multi primitive row01
	// max log size
	d.containerHealthMaxLogSizeField.SetLabel("Max log size:")
	d.containerHealthMaxLogSizeField.SetLabelWidth(healthPageLabelWidth)
	d.containerHealthMaxLogSizeField.SetBackgroundColor(bgColor)
	d.containerHealthMaxLogSizeField.SetLabelColor(style.DialogFgColor)
	d.containerHealthMaxLogSizeField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerHealthMaxLogSizeField.SetFieldWidth(healthPageMultiRowFieldWidth)

	// max log count
	d.containerHealthMaxLogCountField.SetLabel("Max log count:")
	d.containerHealthMaxLogCountField.SetLabelWidth(healthPageSecColLabelWidth)
	d.containerHealthMaxLogCountField.SetBackgroundColor(bgColor)
	d.containerHealthMaxLogCountField.SetLabelColor(style.DialogFgColor)
	d.containerHealthMaxLogCountField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerHealthMaxLogCountField.SetFieldWidth(healthPageMultiRowFieldWidth)

	// on-failure
	onfailureLabel := fmt.Sprintf("%17s: ", "On failure")

	d.containerHealthOnFailureField.SetOptions([]string{"none", "kill", "restart", "stop"}, nil)
	d.containerHealthOnFailureField.SetLabel(onfailureLabel)
	d.containerHealthOnFailureField.SetBackgroundColor(bgColor)
	d.containerHealthOnFailureField.SetLabelColor(style.DialogFgColor)
	d.containerHealthOnFailureField.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	d.containerHealthOnFailureField.SetFieldBackgroundColor(inputFieldBgColor)

	multiItemRow01 := tview.NewFlex().SetDirection(tview.FlexColumn)
	multiItemRow01.AddItem(d.containerHealthMaxLogSizeField, 0, 1, true)
	multiItemRow01.AddItem(d.containerHealthMaxLogCountField, 0, 1, true)
	multiItemRow01.AddItem(d.containerHealthOnFailureField, 0, 1, true)
	multiItemRow01.SetBackgroundColor(bgColor)

	// multi primitive row02
	// interval
	d.containerHealthIntervalField.SetLabel("Interval:")
	d.containerHealthIntervalField.SetLabelWidth(healthPageLabelWidth)
	d.containerHealthIntervalField.SetBackgroundColor(bgColor)
	d.containerHealthIntervalField.SetLabelColor(style.DialogFgColor)
	d.containerHealthIntervalField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerHealthIntervalField.SetFieldWidth(healthPageMultiRowFieldWidth)

	// startup interval
	d.containerHealthStartupIntervalField.SetLabel("Startup interval:")
	d.containerHealthStartupIntervalField.SetLabelWidth(healthPageSecColLabelWidth)
	d.containerHealthStartupIntervalField.SetBackgroundColor(bgColor)
	d.containerHealthStartupIntervalField.SetLabelColor(style.DialogFgColor)
	d.containerHealthStartupIntervalField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerHealthStartupIntervalField.SetFieldWidth(healthPageMultiRowFieldWidth)

	// start period
	startPeroidLabel := fmt.Sprintf("%17s: ", "Start period")
	d.containerHealthStartPeriodField.SetLabel(startPeroidLabel)
	d.containerHealthStartPeriodField.SetBackgroundColor(bgColor)
	d.containerHealthStartPeriodField.SetLabelColor(style.DialogFgColor)
	d.containerHealthStartPeriodField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerHealthStartPeriodField.SetFieldWidth(healthPageMultiRowFieldWidth)

	multiItemRow02 := tview.NewFlex().SetDirection(tview.FlexColumn)
	multiItemRow02.AddItem(d.containerHealthIntervalField, 0, 1, true)
	multiItemRow02.AddItem(d.containerHealthStartupIntervalField, 0, 1, true)
	multiItemRow02.AddItem(d.containerHealthStartPeriodField, 0, 1, true)
	multiItemRow02.SetBackgroundColor(bgColor)

	// multi primitive row03
	// retires
	d.containerHealthRetriesField.SetLabel("Retries:")
	d.containerHealthRetriesField.SetLabelWidth(healthPageLabelWidth)
	d.containerHealthRetriesField.SetBackgroundColor(bgColor)
	d.containerHealthRetriesField.SetLabelColor(style.DialogFgColor)
	d.containerHealthRetriesField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerHealthRetriesField.SetFieldWidth(healthPageMultiRowFieldWidth)

	// startup retries
	d.containerHealthStartupRetriesField.SetLabel("Startup retries:")
	d.containerHealthStartupRetriesField.SetLabelWidth(healthPageSecColLabelWidth)
	d.containerHealthStartupRetriesField.SetBackgroundColor(bgColor)
	d.containerHealthStartupRetriesField.SetLabelColor(style.DialogFgColor)
	d.containerHealthStartupRetriesField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerHealthStartupRetriesField.SetFieldWidth(healthPageMultiRowFieldWidth)

	// startup success
	startupSuccessLabel := fmt.Sprintf("%17s: ", "Startup success")
	d.containerHealthStartupSuccessField.SetLabel(startupSuccessLabel)
	d.containerHealthStartupSuccessField.SetBackgroundColor(bgColor)
	d.containerHealthStartupSuccessField.SetLabelColor(style.DialogFgColor)
	d.containerHealthStartupSuccessField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerHealthStartupSuccessField.SetFieldWidth(healthPageMultiRowFieldWidth)

	multiItemRow03 := tview.NewFlex().SetDirection(tview.FlexColumn)
	multiItemRow03.AddItem(d.containerHealthRetriesField, 0, 1, true)
	multiItemRow03.AddItem(d.containerHealthStartupRetriesField, 0, 1, true)
	multiItemRow03.AddItem(d.containerHealthStartupSuccessField, 0, 1, true)
	multiItemRow03.SetBackgroundColor(bgColor)

	// multi primitive row04
	// timeout
	d.containerHealthTimeoutField.SetLabel("Timeout:")
	d.containerHealthTimeoutField.SetLabelWidth(healthPageLabelWidth)
	d.containerHealthTimeoutField.SetBackgroundColor(bgColor)
	d.containerHealthTimeoutField.SetLabelColor(style.DialogFgColor)
	d.containerHealthTimeoutField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerHealthTimeoutField.SetFieldWidth(healthPageMultiRowFieldWidth)
	// startup timeout
	d.containerHealthStartupTimeoutField.SetLabel("Startup timeout:")
	d.containerHealthStartupTimeoutField.SetLabelWidth(healthPageSecColLabelWidth)
	d.containerHealthStartupTimeoutField.SetBackgroundColor(bgColor)
	d.containerHealthStartupTimeoutField.SetLabelColor(style.DialogFgColor)
	d.containerHealthStartupTimeoutField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerHealthStartupTimeoutField.SetFieldWidth(healthPageMultiRowFieldWidth)

	multiItemRow04 := tview.NewFlex().SetDirection(tview.FlexColumn)
	multiItemRow04.AddItem(d.containerHealthTimeoutField, 0, 1, true)
	multiItemRow04.AddItem(d.containerHealthStartupTimeoutField, 0, 1, true)
	multiItemRow04.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, true)
	multiItemRow04.SetBackgroundColor(bgColor)

	// health page
	d.healthPage.SetDirection(tview.FlexRow)
	d.healthPage.AddItem(d.containerHealthCmdField, 1, 0, true)
	d.healthPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.healthPage.AddItem(d.containerHealthStartupCmdField, 1, 0, true)
	d.healthPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.healthPage.AddItem(d.containerHealthLogDestField, 1, 0, true)
	d.healthPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.healthPage.AddItem(multiItemRow01, 1, 0, true)
	d.healthPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.healthPage.AddItem(multiItemRow02, 1, 0, true)
	d.healthPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.healthPage.AddItem(multiItemRow03, 1, 0, true)
	d.healthPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.healthPage.AddItem(multiItemRow04, 1, 0, true)
	d.healthPage.SetBackgroundColor(bgColor)
}

func (d *ContainerCreateDialog) setupNetworkPageUI() {
	bgColor := style.DialogBgColor
	ddUnselectedStyle := style.DropDownUnselected
	ddselectedStyle := style.DropDownSelected
	inputFieldBgColor := style.InputFieldBgColor
	networkingPageLabelWidth := 13

	// hostname field
	d.containerHostnameField.SetLabel("hostname:")
	d.containerHostnameField.SetLabelWidth(networkingPageLabelWidth)
	d.containerHostnameField.SetBackgroundColor(bgColor)
	d.containerHostnameField.SetLabelColor(style.DialogFgColor)
	d.containerHostnameField.SetFieldBackgroundColor(inputFieldBgColor)

	// IP field
	d.containerIPAddrField.SetLabel("ip address:")
	d.containerIPAddrField.SetLabelWidth(networkingPageLabelWidth)
	d.containerIPAddrField.SetBackgroundColor(bgColor)
	d.containerIPAddrField.SetLabelColor(style.DialogFgColor)
	d.containerIPAddrField.SetFieldBackgroundColor(inputFieldBgColor)

	// mac field
	d.containerMacAddrField.SetLabel("mac address:")
	d.containerMacAddrField.SetLabelWidth(networkingPageLabelWidth)
	d.containerMacAddrField.SetBackgroundColor(bgColor)
	d.containerMacAddrField.SetLabelColor(style.DialogFgColor)
	d.containerMacAddrField.SetFieldBackgroundColor(inputFieldBgColor)

	// network field
	d.containerNetworkField.SetLabel("network:")
	d.containerNetworkField.SetLabelWidth(networkingPageLabelWidth)
	d.containerNetworkField.SetBackgroundColor(bgColor)
	d.containerNetworkField.SetLabelColor(style.DialogFgColor)
	d.containerNetworkField.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	d.containerNetworkField.SetFieldBackgroundColor(inputFieldBgColor)

	// network settings page
	d.networkingPage.SetDirection(tview.FlexRow)
	d.networkingPage.AddItem(d.containerHostnameField, 1, 0, true)
	d.networkingPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.networkingPage.AddItem(d.containerIPAddrField, 1, 0, true)
	d.networkingPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.networkingPage.AddItem(d.containerMacAddrField, 1, 0, true)
	d.networkingPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.networkingPage.AddItem(d.containerNetworkField, 1, 0, true)
	d.networkingPage.SetBackgroundColor(bgColor)
}

func (d *ContainerCreateDialog) setupPortsPageUI() {
	bgColor := style.DialogBgColor
	inputFieldBgColor := style.InputFieldBgColor
	portPageLabelWidth := 15

	inputFieldItems := []struct {
		label  string
		widget *tview.InputField
	}{
		{label: "publish ports:", widget: d.containerPortPublishField},
		{label: "expose ports:", widget: d.containerPortExposeField},
	}

	for _, inputField := range inputFieldItems {
		inputField.widget.SetLabel(inputField.label)
		inputField.widget.SetLabelWidth(portPageLabelWidth)
		inputField.widget.SetBackgroundColor(bgColor)
		inputField.widget.SetLabelColor(style.DialogFgColor)
		inputField.widget.SetFieldBackgroundColor(inputFieldBgColor)
	}

	// publish all field
	d.ContainerPortPublishAllField.SetLabel("publish all ")
	d.ContainerPortPublishAllField.SetLabelWidth(portPageLabelWidth)
	d.ContainerPortPublishAllField.SetBackgroundColor(bgColor)
	d.ContainerPortPublishAllField.SetLabelColor(style.DialogFgColor)
	d.ContainerPortPublishAllField.SetFieldBackgroundColor(inputFieldBgColor)

	d.portPage.SetDirection(tview.FlexRow)
	d.portPage.AddItem(d.containerPortPublishField, 1, 0, true)
	d.portPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.portPage.AddItem(d.ContainerPortPublishAllField, 1, 0, true)
	d.portPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.portPage.AddItem(d.containerPortExposeField, 1, 0, true)
	d.portPage.SetBackgroundColor(bgColor)
}

func (d *ContainerCreateDialog) setupSecurityPageUI() {
	bgColor := style.DialogBgColor
	inputFieldBgColor := style.InputFieldBgColor
	securityOptsLabelWidth := 10

	// selinux label
	d.containerSecLabelField.SetLabel("label:")
	d.containerSecLabelField.SetLabelWidth(securityOptsLabelWidth)
	d.containerSecLabelField.SetBackgroundColor(bgColor)
	d.containerSecLabelField.SetLabelColor(style.DialogFgColor)
	d.containerSecLabelField.SetFieldBackgroundColor(inputFieldBgColor)

	// apparmor
	d.containerSecApparmorField.SetLabel("apparmor:")
	d.containerSecApparmorField.SetLabelWidth(securityOptsLabelWidth)
	d.containerSecApparmorField.SetBackgroundColor(bgColor)
	d.containerSecApparmorField.SetLabelColor(style.DialogFgColor)
	d.containerSecApparmorField.SetFieldBackgroundColor(inputFieldBgColor)

	// seccomp
	d.containerSeccompField.SetLabel("seccomp:")
	d.containerSeccompField.SetLabelWidth(securityOptsLabelWidth)
	d.containerSeccompField.SetBackgroundColor(bgColor)
	d.containerSeccompField.SetLabelColor(style.DialogFgColor)
	d.containerSeccompField.SetFieldBackgroundColor(inputFieldBgColor)

	// mask
	d.containerSecMaskField.SetLabel("mask:")
	d.containerSecMaskField.SetLabelWidth(securityOptsLabelWidth)
	d.containerSecMaskField.SetBackgroundColor(bgColor)
	d.containerSecMaskField.SetLabelColor(style.DialogFgColor)
	d.containerSecMaskField.SetFieldBackgroundColor(inputFieldBgColor)

	// unmask
	d.containerSecUnmaskField.SetLabel("unmask:")
	d.containerSecUnmaskField.SetLabelWidth(securityOptsLabelWidth)
	d.containerSecUnmaskField.SetBackgroundColor(bgColor)
	d.containerSecUnmaskField.SetLabelColor(style.DialogFgColor)
	d.containerSecUnmaskField.SetFieldBackgroundColor(inputFieldBgColor)

	// no-new-privileges
	d.containerSecNoNewPrivField.SetLabel("no new privileges ")
	d.containerSecNoNewPrivField.SetBackgroundColor(bgColor)
	d.containerSecNoNewPrivField.SetLabelColor(style.DialogFgColor)
	d.containerSecNoNewPrivField.SetBackgroundColor(bgColor)
	d.containerSecNoNewPrivField.SetLabelColor(style.DialogFgColor)
	d.containerSecNoNewPrivField.SetFieldBackgroundColor(inputFieldBgColor)

	// security options page
	d.securityOptsPage.SetDirection(tview.FlexRow)
	d.securityOptsPage.AddItem(d.containerSecLabelField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.containerSecApparmorField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.containerSeccompField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.containerSecMaskField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.containerSecUnmaskField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.containerSecNoNewPrivField, 1, 0, true)
}

func (d *ContainerCreateDialog) setupVolumePageUI() {
	bgColor := style.DialogBgColor
	ddUnselectedStyle := style.DropDownUnselected
	ddselectedStyle := style.DropDownSelected
	inputFieldBgColor := style.InputFieldBgColor
	volumePageLabelWidth := 14

	// volume
	d.containerVolumeField.SetLabel("volume:")
	d.containerVolumeField.SetLabelWidth(volumePageLabelWidth)
	d.containerVolumeField.SetBackgroundColor(bgColor)
	d.containerVolumeField.SetLabelColor(style.DialogFgColor)
	d.containerVolumeField.SetFieldBackgroundColor(inputFieldBgColor)

	// image volume
	d.containerImageVolumeField.SetLabel("image volume:")
	d.containerImageVolumeField.SetLabelWidth(volumePageLabelWidth)
	d.containerImageVolumeField.SetBackgroundColor(bgColor)
	d.containerImageVolumeField.SetLabelColor(style.DialogFgColor)
	d.containerImageVolumeField.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	d.containerImageVolumeField.SetFieldBackgroundColor(inputFieldBgColor)

	// mounts
	d.containerMountField.SetLabel("mount:")
	d.containerMountField.SetLabelWidth(volumePageLabelWidth)
	d.containerMountField.SetBackgroundColor(bgColor)
	d.containerMountField.SetLabelColor(style.DialogFgColor)
	d.containerMountField.SetFieldBackgroundColor(inputFieldBgColor)

	// volume settings page
	d.volumePage.SetDirection(tview.FlexRow)
	d.volumePage.AddItem(d.containerVolumeField, 1, 0, true)
	d.volumePage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.volumePage.AddItem(d.containerImageVolumeField, 1, 0, true)
	d.volumePage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.volumePage.AddItem(d.containerMountField, 1, 0, true)
	d.volumePage.SetBackgroundColor(bgColor)
}

func (d *ContainerCreateDialog) setupResourcePageUI() {
	bgColor := style.DialogBgColor
	inputFieldBgColor := style.InputFieldBgColor
	resourcePageLabelWidth := 13
	inputFieldWidth := 18

	getSecondColLabel := func(label string) string {
		return fmt.Sprintf("%18s:", label)
	}

	// memory
	d.containerMemoryField.SetLabel("memory:")
	d.containerMemoryField.SetLabelWidth(resourcePageLabelWidth)
	d.containerMemoryField.SetBackgroundColor(bgColor)
	d.containerMemoryField.SetLabelColor(style.DialogFgColor)
	d.containerMemoryField.SetFieldWidth(inputFieldWidth)
	d.containerMemoryField.SetFieldBackgroundColor(inputFieldBgColor)

	// memory reservation
	memResLabel := "memory reservation:"
	d.containerMemoryReservationField.SetLabel(memResLabel)
	d.containerMemoryReservationField.SetLabelWidth(len(memResLabel) + 1)
	d.containerMemoryReservationField.SetBackgroundColor(bgColor)
	d.containerMemoryReservationField.SetLabelColor(style.DialogFgColor)
	d.containerMemoryReservationField.SetFieldBackgroundColor(inputFieldBgColor)

	// memory swap
	d.containerMemorySwapField.SetLabel("memory swap:")
	d.containerMemorySwapField.SetLabelWidth(resourcePageLabelWidth)
	d.containerMemorySwapField.SetBackgroundColor(bgColor)
	d.containerMemorySwapField.SetLabelColor(style.DialogFgColor)
	d.containerMemorySwapField.SetFieldWidth(inputFieldWidth)
	d.containerMemorySwapField.SetFieldBackgroundColor(inputFieldBgColor)

	// memory swappiness
	d.containerMemorySwappinessField.SetLabel(" memory swappiness:")
	d.containerMemorySwappinessField.SetLabelWidth(len(memResLabel) + 1)
	d.containerMemorySwappinessField.SetBackgroundColor(bgColor)
	d.containerMemorySwappinessField.SetLabelColor(style.DialogFgColor)
	d.containerMemorySwappinessField.SetFieldBackgroundColor(inputFieldBgColor)

	// memRow1
	memRow1Layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	memRow1Layout.AddItem(d.containerMemoryField, 0, 1, true)
	memRow1Layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	memRow1Layout.AddItem(d.containerMemoryReservationField, 0, 1, true)
	memRow1Layout.SetBackgroundColor(bgColor)

	// memRow2
	memRow2Layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	memRow2Layout.AddItem(d.containerMemorySwapField, 0, 1, true)
	memRow2Layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	memRow2Layout.AddItem(d.containerMemorySwappinessField, 0, 1, true)
	memRow2Layout.SetBackgroundColor(bgColor)

	// cpus
	d.containerCPUsField.SetLabel("cpus:")
	d.containerCPUsField.SetLabelWidth(resourcePageLabelWidth)
	d.containerCPUsField.SetBackgroundColor(bgColor)
	d.containerCPUsField.SetLabelColor(style.DialogFgColor)
	d.containerCPUsField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerCPUsField.SetFieldWidth(inputFieldWidth)

	// cpu shares
	d.containerCPUSharesField.SetLabel(getSecondColLabel("cpu shares"))
	d.containerCPUSharesField.SetLabelWidth(len(memResLabel) + 1)
	d.containerCPUSharesField.SetBackgroundColor(bgColor)
	d.containerCPUSharesField.SetLabelColor(style.DialogFgColor)
	d.containerCPUSharesField.SetFieldBackgroundColor(inputFieldBgColor)

	// cpuRow1
	cpuRow1Layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	cpuRow1Layout.AddItem(d.containerCPUsField, 0, 1, true)
	cpuRow1Layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	cpuRow1Layout.AddItem(d.containerCPUSharesField, 0, 1, true)
	cpuRow1Layout.SetBackgroundColor(bgColor)

	// cpus period
	d.containerCPUPeriodField.SetLabel("cpu period:")
	d.containerCPUPeriodField.SetLabelWidth(resourcePageLabelWidth)
	d.containerCPUPeriodField.SetBackgroundColor(bgColor)
	d.containerCPUPeriodField.SetLabelColor(style.DialogFgColor)
	d.containerCPUPeriodField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerCPUPeriodField.SetFieldWidth(inputFieldWidth)

	// cpu rt period
	d.containerCPURtPeriodField.SetLabel(getSecondColLabel("cpu rt period"))
	d.containerCPURtPeriodField.SetLabelWidth(len(memResLabel) + 1)
	d.containerCPURtPeriodField.SetBackgroundColor(bgColor)
	d.containerCPURtPeriodField.SetLabelColor(style.DialogFgColor)
	d.containerCPURtPeriodField.SetFieldBackgroundColor(inputFieldBgColor)

	// cpuRow2
	cpuRow2Layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	cpuRow2Layout.AddItem(d.containerCPUPeriodField, 0, 1, true)
	cpuRow2Layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	cpuRow2Layout.AddItem(d.containerCPURtPeriodField, 0, 1, true)
	cpuRow2Layout.SetBackgroundColor(bgColor)

	// cpus quota
	d.containerCPUQuotaField.SetLabel("cpu quota:")
	d.containerCPUQuotaField.SetLabelWidth(resourcePageLabelWidth)
	d.containerCPUQuotaField.SetBackgroundColor(bgColor)
	d.containerCPUQuotaField.SetLabelColor(style.DialogFgColor)
	d.containerCPUQuotaField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerCPUQuotaField.SetFieldWidth(inputFieldWidth)

	// cpu rt runtime
	d.containerCPURtRuntimeField.SetLabel(getSecondColLabel("cpu rt runtime"))
	d.containerCPURtRuntimeField.SetLabelWidth(len(memResLabel) + 1)
	d.containerCPURtRuntimeField.SetBackgroundColor(bgColor)
	d.containerCPURtRuntimeField.SetLabelColor(style.DialogFgColor)
	d.containerCPURtRuntimeField.SetFieldBackgroundColor(inputFieldBgColor)

	// cpuRow3
	cpuRow3Layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	cpuRow3Layout.AddItem(d.containerCPUQuotaField, 0, 1, true)
	cpuRow3Layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	cpuRow3Layout.AddItem(d.containerCPURtRuntimeField, 0, 1, true)
	cpuRow3Layout.SetBackgroundColor(bgColor)

	// cpuset cpus
	d.containerCPUSetCPUsField.SetLabel("cpuset cpus:")
	d.containerCPUSetCPUsField.SetLabelWidth(resourcePageLabelWidth)
	d.containerCPUSetCPUsField.SetBackgroundColor(bgColor)
	d.containerCPUSetCPUsField.SetLabelColor(style.DialogFgColor)
	d.containerCPUSetCPUsField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerCPUSetCPUsField.SetFieldWidth(inputFieldWidth)

	// cpuset mems
	d.containerCPUSetMemsField.SetLabel(getSecondColLabel("cpuset mems"))
	d.containerCPUSetMemsField.SetLabelWidth(len(memResLabel) + 1)
	d.containerCPUSetMemsField.SetBackgroundColor(bgColor)
	d.containerCPUSetMemsField.SetLabelColor(style.DialogFgColor)
	d.containerCPUSetMemsField.SetFieldBackgroundColor(inputFieldBgColor)

	// cpuRow4
	cpuRow4Layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	cpuRow4Layout.AddItem(d.containerCPUSetCPUsField, 0, 1, true)
	cpuRow4Layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	cpuRow4Layout.AddItem(d.containerCPUSetMemsField, 0, 1, true)
	cpuRow4Layout.SetBackgroundColor(bgColor)

	// shm size
	d.containerShmSizeField.SetLabel("shm size:")
	d.containerShmSizeField.SetLabelWidth(resourcePageLabelWidth)
	d.containerShmSizeField.SetBackgroundColor(bgColor)
	d.containerShmSizeField.SetLabelColor(style.DialogFgColor)
	d.containerShmSizeField.SetFieldBackgroundColor(inputFieldBgColor)
	d.containerShmSizeField.SetFieldWidth(inputFieldWidth)

	// shm size systemd
	d.containerShmSizeSystemdField.SetLabel(getSecondColLabel("shm size systemd"))
	d.containerShmSizeSystemdField.SetLabelWidth(len(memResLabel) + 1)
	d.containerShmSizeSystemdField.SetBackgroundColor(bgColor)
	d.containerShmSizeSystemdField.SetLabelColor(style.DialogFgColor)
	d.containerShmSizeSystemdField.SetFieldBackgroundColor(inputFieldBgColor)

	// shmRow1
	shmRow1Layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	shmRow1Layout.AddItem(d.containerShmSizeField, 0, 1, true)
	shmRow1Layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	shmRow1Layout.AddItem(d.containerShmSizeSystemdField, 0, 1, true)
	shmRow1Layout.SetBackgroundColor(bgColor)

	// resource settings page
	d.resourcePage.SetDirection(tview.FlexRow)
	d.resourcePage.AddItem(memRow1Layout, 1, 0, true)
	d.resourcePage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.resourcePage.AddItem(memRow2Layout, 1, 0, true)
	d.resourcePage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.resourcePage.AddItem(cpuRow1Layout, 1, 0, true)
	d.resourcePage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.resourcePage.AddItem(cpuRow2Layout, 1, 0, true)
	d.resourcePage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.resourcePage.AddItem(cpuRow3Layout, 1, 0, true)
	d.resourcePage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.resourcePage.AddItem(cpuRow4Layout, 1, 0, true)
	d.resourcePage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.resourcePage.AddItem(shmRow1Layout, 1, 0, true)
	d.resourcePage.SetBackgroundColor(bgColor)
}

func (d *ContainerCreateDialog) setupNamespacePageUI() {
	bgColor := style.DialogBgColor
	inputFieldBgColor := style.InputFieldBgColor
	namespacePageLabelWidth := 10

	// cgroupns
	d.containerNamespaceCgroupField.SetLabel("cgroupns:")
	d.containerNamespaceCgroupField.SetLabelWidth(namespacePageLabelWidth)
	d.containerNamespaceCgroupField.SetBackgroundColor(bgColor)
	d.containerNamespaceCgroupField.SetLabelColor(style.DialogFgColor)
	d.containerNamespaceCgroupField.SetFieldBackgroundColor(inputFieldBgColor)

	// ipc
	d.containerNamespaceIpcField.SetLabel("ipc:")
	d.containerNamespaceIpcField.SetLabelWidth(namespacePageLabelWidth)
	d.containerNamespaceIpcField.SetBackgroundColor(bgColor)
	d.containerNamespaceIpcField.SetLabelColor(style.DialogFgColor)
	d.containerNamespaceIpcField.SetFieldBackgroundColor(inputFieldBgColor)

	// pid
	d.containerNamespacePidField.SetLabel("pid:")
	d.containerNamespacePidField.SetLabelWidth(namespacePageLabelWidth)
	d.containerNamespacePidField.SetBackgroundColor(bgColor)
	d.containerNamespacePidField.SetLabelColor(style.DialogFgColor)
	d.containerNamespacePidField.SetFieldBackgroundColor(inputFieldBgColor)

	// userns
	d.containerNamespaceUserField.SetLabel("userns:")
	d.containerNamespaceUserField.SetLabelWidth(namespacePageLabelWidth)
	d.containerNamespaceUserField.SetBackgroundColor(bgColor)
	d.containerNamespaceUserField.SetLabelColor(style.DialogFgColor)
	d.containerNamespaceUserField.SetFieldBackgroundColor(inputFieldBgColor)

	// uts
	d.containerNamespaceUtsField.SetLabel("uts:")
	d.containerNamespaceUtsField.SetLabelWidth(namespacePageLabelWidth)
	d.containerNamespaceUtsField.SetBackgroundColor(bgColor)
	d.containerNamespaceUtsField.SetLabelColor(style.DialogFgColor)
	d.containerNamespaceUtsField.SetFieldBackgroundColor(inputFieldBgColor)

	// uidmap
	d.containerNamespaceUidmapField.SetLabel("uidmap:")
	d.containerNamespaceUidmapField.SetLabelWidth(namespacePageLabelWidth)
	d.containerNamespaceUidmapField.SetBackgroundColor(bgColor)
	d.containerNamespaceUidmapField.SetLabelColor(style.DialogFgColor)
	d.containerNamespaceUidmapField.SetFieldBackgroundColor(inputFieldBgColor)

	// subuidname
	d.containerNamespaceSubuidNameField.SetLabel("subuidname: ")
	d.containerNamespaceSubuidNameField.SetBackgroundColor(bgColor)
	d.containerNamespaceSubuidNameField.SetLabelColor(style.DialogFgColor)
	d.containerNamespaceSubuidNameField.SetFieldBackgroundColor(inputFieldBgColor)

	// gidmap
	d.containerNamespaceGidmapField.SetLabel("gidmap:")
	d.containerNamespaceGidmapField.SetLabelWidth(namespacePageLabelWidth)
	d.containerNamespaceGidmapField.SetBackgroundColor(bgColor)
	d.containerNamespaceGidmapField.SetLabelColor(style.DialogFgColor)
	d.containerNamespaceGidmapField.SetFieldBackgroundColor(inputFieldBgColor)

	// subgidname
	d.containerNamespaceSubgidNameField.SetLabel("subgidname: ")
	d.containerNamespaceSubgidNameField.SetBackgroundColor(bgColor)
	d.containerNamespaceSubgidNameField.SetLabelColor(style.DialogFgColor)
	d.containerNamespaceSubgidNameField.SetFieldBackgroundColor(inputFieldBgColor)

	// mapRow01Layout
	mapRow01Layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	mapRow01Layout.AddItem(d.containerNamespaceUidmapField, 0, 1, true)
	mapRow01Layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	mapRow01Layout.AddItem(d.containerNamespaceSubuidNameField, 0, 1, true)
	mapRow01Layout.SetBackgroundColor(bgColor)

	// mapRow02Layout
	mapRow02Layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	mapRow02Layout.AddItem(d.containerNamespaceGidmapField, 0, 1, true)
	mapRow02Layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	mapRow02Layout.AddItem(d.containerNamespaceSubgidNameField, 0, 1, true)
	mapRow02Layout.SetBackgroundColor(bgColor)

	// namespace options page
	d.namespacePage.SetDirection(tview.FlexRow)
	d.namespacePage.AddItem(d.containerNamespaceCgroupField, 1, 0, true)
	d.namespacePage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.namespacePage.AddItem(d.containerNamespaceIpcField, 1, 0, true)
	d.namespacePage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.namespacePage.AddItem(d.containerNamespacePidField, 1, 0, true)
	d.namespacePage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.namespacePage.AddItem(d.containerNamespaceUserField, 1, 0, true)
	d.namespacePage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.namespacePage.AddItem(d.containerNamespaceUtsField, 1, 0, true)
	d.namespacePage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.namespacePage.AddItem(mapRow01Layout, 1, 0, true)
	d.namespacePage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.namespacePage.AddItem(mapRow02Layout, 1, 0, true)

	d.resourcePage.SetBackgroundColor(bgColor)
}

// Display displays this primitive.
func (d *ContainerCreateDialog) Display() {
	d.display = true
	d.initData()
	d.focusElement = createCategoryPagesFocus
}

// IsDisplay returns true if primitive is shown.
func (d *ContainerCreateDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *ContainerCreateDialog) Hide() {
	d.display = false
}

// HasFocus returns whether or not this primitive has focus.
func (d *ContainerCreateDialog) HasFocus() bool {
	if d.categories.HasFocus() || d.categoryPages.HasFocus() {
		return true
	}

	return d.Box.HasFocus() || d.form.HasFocus()
}

// dropdownHasFocus returns true if container create dialog dropdown primitives
// has focus.
func (d *ContainerCreateDialog) dropdownHasFocus() bool {
	if d.containerImageField.HasFocus() || d.containerPodField.HasFocus() {
		return true
	}

	if d.containerNetworkField.HasFocus() || d.containerImageVolumeField.HasFocus() {
		return true
	}

	return d.containerHealthOnFailureField.HasFocus()
}

// Focus is called when this primitive receives focus.
func (d *ContainerCreateDialog) Focus(delegate func(p tview.Primitive)) { //nolint:gocyclo,cyclop,maintidx
	switch d.focusElement {
	// form has focus
	case createContainerFormFocus:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = createCategoriesFocus // category text view

				d.Focus(delegate)
				d.form.SetFocus(0)

				return nil
			}

			if event.Key() == tcell.KeyEnter {
				return nil
			}

			return event
		})

		delegate(d.form)
	// category text view
	case createCategoriesFocus:
		d.categories.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = createCategoryPagesFocus // category page view
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
	// container info page
	case createContainerNameFieldFocus:
		delegate(d.containerNameField)
	case createContainerCommandFieldFocus:
		delegate(d.containerCommandField)
	case createContainerImageFieldFocus:
		delegate(d.containerImageField)
	case createcontainerPodFieldFocus:
		delegate(d.containerPodField)
	case createContainerLabelsFieldFocus:
		delegate(d.containerLabelsField)
	case createContainerRemoveFieldFocus:
		delegate(d.containerRemoveField)
	case createContainerPrivilegedFieldFocus:
		delegate(d.containerPrivilegedField)
	case createContainerTimeoutFieldFocus:
		delegate(d.containerTimeoutField)
	case createContainerInteractiveFieldFocus:
		delegate(d.containerInteractiveField)
	case createContainerTtyFieldFocus:
		delegate(d.containerTtyField)
	case createContainerDetachFieldFocus:
		delegate(d.containerDetachField)
	case createContainerSecretFieldFocus:
		delegate(d.containerSecretField)
	// environment options page
	case createContainerWorkDirFieldFocus:
		delegate(d.containerWorkDirField)
	case createContainerEnvVarsFieldFocus:
		delegate(d.containerEnvVarsField)
	case createContainerEnvFileFieldFocus:
		delegate(d.containerEnvFileField)
	case createContainerEnvMergeFieldFocus:
		delegate(d.containerEnvMergeField)
	case createContainerUnsetEnvFieldFocus:
		delegate(d.containerUnsetEnvField)
	case createContainerEnvHostFieldFocus:
		delegate(d.containerEnvHostField)
	case createContainerUnsetEnvAllFieldFocus:
		delegate(d.containerUnsetEnvAllField)
	case createContainerUmaskFieldFocus:
		delegate(d.containerUmaskField)
	// user and groups page
	case createContainerUserFieldFocus:
		delegate(d.containerUserField)
	case createContainerHostUsersFieldFocus:
		delegate(d.containerHostUsersField)
	case createContainerPasswdEntryFieldFocus:
		delegate(d.containerPasswdEntryField)
	case createContainerGroupEntryFieldFocus:
		delegate(d.containerGroupEntryField)
	// security options page
	case createcontainerSecLabelFieldFocus:
		delegate(d.containerSecLabelField)
	case createContainerApprarmorFieldFocus:
		delegate(d.containerSecApparmorField)
	case createContainerSeccompFeildFocus:
		delegate(d.containerSeccompField)
	case createcontainerSecMaskFieldFocus:
		delegate(d.containerSecMaskField)
	case createcontainerSecUnmaskFieldFocus:
		delegate(d.containerSecUnmaskField)
	case createcontainerSecNoNewPrivFieldFocus:
		delegate(d.containerSecNoNewPrivField)
	// networking page
	case createContainerHostnameFieldFocus:
		delegate(d.containerHostnameField)
	case createContainerIPAddrFieldFocus:
		delegate(d.containerIPAddrField)
	case createContainerMacAddrFieldFocus:
		delegate(d.containerMacAddrField)
	case createContainerNetworkFieldFocus:
		delegate(d.containerNetworkField)
	// port page
	// networking page
	case createContainerPortPublishFieldFocus:
		delegate(d.containerPortPublishField)
	case createContainerPortPublishAllFieldFocus:
		delegate(d.ContainerPortPublishAllField)
	case createContainerPortExposeFieldFocus:
		delegate(d.containerPortExposeField)
	// dns page
	case createContainerDNSServersFieldFocus:
		delegate(d.containerDNSServersField)
	case createContainerDNSOptionsFieldFocus:
		delegate(d.containerDNSOptionsField)
	case createContainerDNSSearchFieldFocus:
		delegate(d.containerDNSSearchField)
	// volume page
	case createContainerVolumeFieldFocus:
		delegate(d.containerVolumeField)
	case createContainerImageVolumeFieldFocus:
		delegate(d.containerImageVolumeField)
	case createContainerMountFieldFocus:
		delegate(d.containerMountField)
	// health page
	case createContainerHealthCmdFieldFocus:
		delegate(d.containerHealthCmdField)
	case createContainerHealthStartupCmdFieldFocus:
		delegate(d.containerHealthStartupCmdField)
	case createContainerHealthOnFailureFieldFocus:
		delegate(d.containerHealthOnFailureField)
	case createContainerHealthIntervalFieldFocus:
		delegate(d.containerHealthIntervalField)
	case createContainerHealthStartupIntervalFieldFocus:
		delegate(d.containerHealthStartupIntervalField)
	case createContainerHealthTimeoutFieldFocus:
		delegate(d.containerHealthTimeoutField)
	case createContainerHealthStartupTimeoutFieldFocus:
		delegate(d.containerHealthStartupTimeoutField)
	case createContainerHealthRetriesFieldFocus:
		delegate(d.containerHealthRetriesField)
	case createContainerHealthStartupRetriesFieldFocus:
		delegate(d.containerHealthStartupRetriesField)
	case createContainerHealthStartPeriodFieldFocus:
		delegate(d.containerHealthStartPeriodField)
	case createContainerHealthStartupSuccessFieldFocus:
		delegate(d.containerHealthStartupSuccessField)
	case createContainerHealthLogDestFocus:
		delegate(d.containerHealthLogDestField)
	case createContainerHealthMaxLogCountFocus:
		delegate(d.containerHealthMaxLogCountField)
	case createContainerHealthMaxLogSizeFocus:
		delegate(d.containerHealthMaxLogSizeField)
	// resource page
	case createContainerMemoryFieldFocus:
		delegate(d.containerMemoryField)
	case createContainerMemoryReservatoinFieldFocus:
		delegate(d.containerMemoryReservationField)
	case createContainerMemorySwapFieldFocus:
		delegate(d.containerMemorySwapField)
	case createcontainerMemorySwappinessFieldFocus:
		delegate(d.containerMemorySwappinessField)
	case createContainerCPUsFieldFocus:
		delegate(d.containerCPUsField)
	case createContainerCPUSharesFieldFocus:
		delegate(d.containerCPUSharesField)
	case createContainerCPUPeriodFieldFocus:
		delegate(d.containerCPUPeriodField)
	case createContainerCPURtPeriodFieldFocus:
		delegate(d.containerCPURtPeriodField)
	case createContainerCPUQuotaFieldFocus:
		delegate(d.containerCPUQuotaField)
	case createContainerCPURtRuntimeFeildFocus:
		delegate(d.containerCPURtRuntimeField)
	case createContainerCPUSetCPUsFieldFocus:
		delegate(d.containerCPUSetCPUsField)
	case createContainerCPUSetMemsFieldFocus:
		delegate(d.containerCPUSetMemsField)
	case createContainerShmSizeFieldFocus:
		delegate(d.containerShmSizeField)
	case createContainerShmSizeSystemdFieldFocus:
		delegate(d.containerShmSizeSystemdField)
	// namespace page
	case createContainerNamespaceCgroupFieldFocus:
		delegate(d.containerNamespaceCgroupField)
	case createContainerNamespaceIpcFieldFocus:
		delegate(d.containerNamespaceIpcField)
	case createContainerNamespacePidFieldFocus:
		delegate(d.containerNamespacePidField)
	case createContainerNamespaceUserFieldFocus:
		delegate(d.containerNamespaceUserField)
	case createContainerNamespaceUtsFieldFocus:
		delegate(d.containerNamespaceUtsField)
	case createContainerNamespaceUidmapFieldFocus:
		delegate(d.containerNamespaceUidmapField)
	case createContainerNamespaceSubuidNameFieldFocus:
		delegate(d.containerNamespaceSubuidNameField)
	case createContainerNamespaceGidmapFieldFocus:
		delegate(d.containerNamespaceGidmapField)
	case createContainerNamespaceSubgidNameFieldFocus:
		delegate(d.containerNamespaceSubgidNameField)
	// category page
	case createCategoryPagesFocus:
		delegate(d.categoryPages)
	}
}

func (d *ContainerCreateDialog) initCustomInputHanlers() {
	// pod name dropdown
	d.containerPodField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = utils.ParseKeyEventKey(event)

		return event
	})

	// container image volume
	d.containerImageField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = utils.ParseKeyEventKey(event)

		return event
	})

	// container network dropdown
	d.containerNetworkField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = utils.ParseKeyEventKey(event)

		return event
	})

	// container image volume dropdown
	d.containerImageVolumeField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = utils.ParseKeyEventKey(event)

		return event
	})
}

// InputHandler returns input handler function for this primitive.
func (d *ContainerCreateDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) { //nolint:gocognit,lll,cyclop,gocyclo
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("container create dialog: event %v received", event)

		if event.Key() == tcell.KeyEsc && !d.dropdownHasFocus() {
			d.cancelHandler()

			return
		}

		if d.containerInfoPage.HasFocus() {
			if handler := d.containerInfoPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setContainerInfoPageNextFocus()
				}

				handler(event, setFocus)

				return
			}
		}

		if d.environmentPage.HasFocus() {
			if handler := d.environmentPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setEnvironmentPageNextFocus()
				}

				handler(event, setFocus)

				return
			}
		}

		if d.userGroupsPage.HasFocus() {
			if handler := d.userGroupsPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setUserGroupsPageNextFocus()
				}

				handler(event, setFocus)

				return
			}
		}

		if d.dnsPage.HasFocus() {
			if handler := d.dnsPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setDNSSettingsPageNextFocus()
				}

				handler(event, setFocus)

				return
			}
		}

		if d.networkingPage.HasFocus() {
			if handler := d.networkingPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setNetworkSettingsPageNextFocus()
				}

				handler(event, setFocus)

				return
			}
		}

		if d.portPage.HasFocus() {
			if handler := d.portPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setPortPageNextFocus()
				}

				handler(event, setFocus)

				return
			}
		}

		if d.securityOptsPage.HasFocus() {
			if handler := d.securityOptsPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setSecurityOptionsPageNextFocus()
				}

				handler(event, setFocus)

				return
			}
		}

		if d.volumePage.HasFocus() {
			if handler := d.volumePage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setVolumeSettingsPageNextFocus()
				}

				handler(event, setFocus)

				return
			}
		}

		if d.healthPage.HasFocus() {
			if handler := d.healthPage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setHealthSettingsPageNextFocus()
				}

				handler(event, setFocus)

				return
			}
		}

		if d.resourcePage.HasFocus() {
			if handler := d.resourcePage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setResourceSettingsPageNextFocus()
				}

				handler(event, setFocus)

				return
			}
		}

		if d.namespacePage.HasFocus() {
			if handler := d.namespacePage.InputHandler(); handler != nil {
				if event.Key() == tcell.KeyTab {
					d.setNamespaceOptionsPageNextFocus()
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
						d.enterHandler()
					}
				}

				formHandler(event, setFocus)

				return
			}
		}
	})
}

// SetRect set rects for this primitive.
func (d *ContainerCreateDialog) SetRect(x, y, width, height int) {
	if width > containerCreateDialogMaxWidth {
		emptySpace := (width - containerCreateDialogMaxWidth) / 2 //nolint:mnd
		x += emptySpace
		width = containerCreateDialogMaxWidth
	}

	maxAllowedHeight := containerCreateOnlyDialogHeight
	if d.mode == ContainerCreateAndRunDialogMode {
		maxAllowedHeight = containerCreateAndRunDialogHeight
	}

	if height > maxAllowedHeight {
		emptySpace := (height - maxAllowedHeight) / 2 //nolint:mnd
		y += emptySpace
		height = maxAllowedHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen.
func (d *ContainerCreateDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.Box.DrawForSubclass(screen, d)

	x, y, width, height := d.Box.GetInnerRect()

	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetCancelFunc sets form cancel button selected function.
func (d *ContainerCreateDialog) SetCancelFunc(handler func()) *ContainerCreateDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:mnd

	cancelButton.SetSelectedFunc(handler)

	return d
}

// SetHandlerFunc sets form create or run button selected function.
func (d *ContainerCreateDialog) SetHandlerFunc(handler func()) *ContainerCreateDialog {
	d.enterHandler = handler
	enterButton := d.form.GetButton(d.form.GetButtonCount() - 1)

	enterButton.SetSelectedFunc(handler)

	return d
}

func (d *ContainerCreateDialog) setActiveCategory(index int) {
	fgColor := style.DialogFgColor
	bgColor := style.ButtonBgColor
	ctgTextColor := style.GetColorHex(fgColor)
	ctgBgColor := style.GetColorHex(bgColor)

	d.activePageIndex = index

	d.categories.Clear()

	ctgList := []string{}

	alignedList, _ := utils.AlignStringListWidth(d.categoryLabels)

	for i := range alignedList {
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

func (d *ContainerCreateDialog) nextCategory() {
	activePage := d.activePageIndex
	if d.activePageIndex < len(d.categoryLabels)-1 {
		activePage++

		d.setActiveCategory(activePage)

		return
	}

	d.setActiveCategory(0)
}

func (d *ContainerCreateDialog) previousCategory() {
	activePage := d.activePageIndex
	if d.activePageIndex > 0 {
		activePage--

		d.setActiveCategory(activePage)

		return
	}

	d.setActiveCategory(len(d.categoryLabels) - 1)
}

func (d *ContainerCreateDialog) initData() {
	// get available images
	imgList, _ := images.List()
	d.imageList = imgList
	imgOptions := []string{""}

	for i := range d.imageList {
		if d.imageList[i].ID == "<none>" {
			imgOptions = append(imgOptions, d.imageList[i].ID)

			continue
		}

		imgname := d.imageList[i].Repository + ":" + d.imageList[i].Tag
		imgOptions = append(imgOptions, imgname)
	}

	// get available pods
	podOptions := []string{""}
	podList, _ := pods.List()
	d.podList = podList

	for i := range podList {
		podOptions = append(podOptions, podList[i].Name)
	}

	// get available networks
	networkOptions := []string{""}
	networkList, _ := networks.List()

	for i := range networkList {
		networkOptions = append(networkOptions, networkList[i][1])
	}

	// get available volumes
	imageVolumeOptions := []string{"", "ignore", "tmpfs", "anonymous"}
	volumeOptions := []string{""}
	volList, _ := volumes.List()

	for i := range volList {
		volumeOptions = append(volumeOptions, volList[i].Name)
	}

	d.setActiveCategory(0)
	// container category
	d.containerNameField.SetText("")
	d.containerCommandField.SetText("")
	d.containerImageField.SetOptions(imgOptions, nil)
	d.containerImageField.SetCurrentOption(0)
	d.containerPodField.SetOptions(podOptions, nil)
	d.containerPodField.SetCurrentOption(0)
	d.containerLabelsField.SetText("")
	d.containerRemoveField.SetChecked(false)
	d.containerPrivilegedField.SetChecked(false)
	d.containerTimeoutField.SetText("")
	d.containerSecretField.SetText("")
	d.containerInteractiveField.SetChecked(false)
	d.containerTtyField.SetChecked(false)
	d.containerDetachField.SetChecked(false)

	// environment category
	d.containerWorkDirField.SetText("")
	d.containerEnvVarsField.SetText("")
	d.containerEnvFileField.SetText("")
	d.containerEnvMergeField.SetText("")
	d.containerUnsetEnvField.SetText("")
	d.containerEnvHostField.SetChecked(false)
	d.containerUnsetEnvAllField.SetChecked(false)
	d.containerUmaskField.SetText("")

	// user and groups category
	d.containerUserField.SetText("")
	d.containerHostUsersField.SetText("")
	d.containerPasswdEntryField.SetText("")
	d.containerGroupEntryField.SetText("")

	// dns settings category
	d.containerDNSServersField.SetText("")
	d.containerDNSSearchField.SetText("")
	d.containerDNSOptionsField.SetText("")

	// health options category
	d.containerHealthCmdField.SetText("")
	d.containerHealthStartupCmdField.SetText("")
	d.containerHealthOnFailureField.SetCurrentOption(0)
	d.containerHealthIntervalField.SetText("30s")
	d.containerHealthStartupIntervalField.SetText("30s")
	d.containerHealthTimeoutField.SetText("30s")
	d.containerHealthStartupTimeoutField.SetText("30s")
	d.containerHealthRetriesField.SetText("3")
	d.containerHealthStartupRetriesField.SetText("")
	d.containerHealthStartPeriodField.SetText("0s")
	d.containerHealthStartupSuccessField.SetText("")
	d.containerHealthLogDestField.SetText("local")
	d.containerHealthMaxLogCountField.SetText("5")
	d.containerHealthMaxLogSizeField.SetText("500")

	// network settings category
	d.containerHostnameField.SetText("")
	d.containerIPAddrField.SetText("")
	d.containerMacAddrField.SetText("")
	d.containerNetworkField.SetOptions(networkOptions, nil)
	d.containerNetworkField.SetCurrentOption(0)

	// ports settings category
	d.containerPortPublishField.SetText("")
	d.ContainerPortPublishAllField.SetChecked(false)
	d.containerPortExposeField.SetText("")

	// security options category
	d.containerSecLabelField.SetText("")
	d.containerSecApparmorField.SetText("")
	d.containerSeccompField.SetText("")
	d.containerSecMaskField.SetText("")
	d.containerSecUnmaskField.SetText("")
	d.containerSecNoNewPrivField.SetChecked(false)

	// volumes options category
	d.containerVolumeField.SetText("")
	d.containerMountField.SetText("")
	d.containerImageVolumeField.SetOptions(imageVolumeOptions, nil)
	d.containerImageVolumeField.SetCurrentOption(0)

	// resource settings category
	d.containerMemoryField.SetText("")
	d.containerMemoryReservationField.SetText("")
	d.containerMemorySwapField.SetText("")
	d.containerMemorySwappinessField.SetText("")
	d.containerCPUsField.SetText("")
	d.containerCPUSharesField.SetText("")
	d.containerCPUPeriodField.SetText("")
	d.containerCPURtPeriodField.SetText("")
	d.containerCPUQuotaField.SetText("")
	d.containerCPURtRuntimeField.SetText("")
	d.containerCPUSetCPUsField.SetText("")
	d.containerCPUSetMemsField.SetText("")
	d.containerShmSizeField.SetText("")
	d.containerShmSizeSystemdField.SetText("")

	// namespace options category
	d.containerNamespaceCgroupField.SetText("")
	d.containerNamespaceIpcField.SetText("")
	d.containerNamespacePidField.SetText("")
	d.containerNamespaceUserField.SetText("")
	d.containerNamespaceUtsField.SetText("")
	d.containerNamespaceUidmapField.SetText("")
	d.containerNamespaceGidmapField.SetText("")
	d.containerNamespaceSubuidNameField.SetText("")
	d.containerNamespaceSubgidNameField.SetText("")
}

func (d *ContainerCreateDialog) setPortPageNextFocus() {
	if d.containerPortPublishField.HasFocus() {
		d.focusElement = createContainerPortPublishAllFieldFocus

		return
	}

	if d.ContainerPortPublishAllField.HasFocus() {
		d.focusElement = createContainerPortExposeFieldFocus

		return
	}

	d.focusElement = createContainerFormFocus
}

func (d *ContainerCreateDialog) setNamespaceOptionsPageNextFocus() {
	if d.containerNamespaceCgroupField.HasFocus() {
		d.focusElement = createContainerNamespaceIpcFieldFocus

		return
	}

	if d.containerNamespaceIpcField.HasFocus() {
		d.focusElement = createContainerNamespacePidFieldFocus

		return
	}

	if d.containerNamespacePidField.HasFocus() {
		d.focusElement = createContainerNamespaceUserFieldFocus

		return
	}

	if d.containerNamespaceUserField.HasFocus() {
		d.focusElement = createContainerNamespaceUtsFieldFocus

		return
	}

	if d.containerNamespaceUtsField.HasFocus() {
		d.focusElement = createContainerNamespaceUidmapFieldFocus

		return
	}

	if d.containerNamespaceUidmapField.HasFocus() {
		d.focusElement = createContainerNamespaceSubuidNameFieldFocus

		return
	}

	if d.containerNamespaceSubuidNameField.HasFocus() {
		d.focusElement = createContainerNamespaceGidmapFieldFocus

		return
	}

	if d.containerNamespaceGidmapField.HasFocus() {
		d.focusElement = createContainerNamespaceSubgidNameFieldFocus

		return
	}

	d.focusElement = createContainerFormFocus
}

func (d *ContainerCreateDialog) setContainerInfoPageNextFocus() { //nolint:cyclop
	if d.containerNameField.HasFocus() {
		d.focusElement = createContainerCommandFieldFocus

		return
	}

	if d.containerCommandField.HasFocus() {
		d.focusElement = createContainerImageFieldFocus

		return
	}

	if d.containerImageField.HasFocus() {
		d.focusElement = createcontainerPodFieldFocus

		return
	}

	if d.containerPodField.HasFocus() {
		d.focusElement = createContainerLabelsFieldFocus

		return
	}

	if d.containerLabelsField.HasFocus() {
		d.focusElement = createContainerPrivilegedFieldFocus

		return
	}

	if d.containerPrivilegedField.HasFocus() {
		d.focusElement = createContainerRemoveFieldFocus

		return
	}

	if d.containerRemoveField.HasFocus() {
		d.focusElement = createContainerTimeoutFieldFocus

		return
	}

	if d.containerTimeoutField.HasFocus() {
		if d.mode == ContainerCreateOnlyDialogMode {
			d.focusElement = createContainerSecretFieldFocus

			return
		}

		d.focusElement = createContainerDetachFieldFocus

		return
	}

	if d.containerDetachField.HasFocus() {
		d.focusElement = createContainerInteractiveFieldFocus

		return
	}

	if d.containerInteractiveField.HasFocus() {
		d.focusElement = createContainerTtyFieldFocus

		return
	}

	if d.containerTtyField.HasFocus() {
		d.focusElement = createContainerSecretFieldFocus

		return
	}

	d.focusElement = createContainerFormFocus
}

func (d *ContainerCreateDialog) setEnvironmentPageNextFocus() {
	if d.containerWorkDirField.HasFocus() {
		d.focusElement = createContainerEnvVarsFieldFocus

		return
	}

	if d.containerEnvVarsField.HasFocus() {
		d.focusElement = createContainerEnvFileFieldFocus

		return
	}

	if d.containerEnvFileField.HasFocus() {
		d.focusElement = createContainerEnvMergeFieldFocus

		return
	}

	if d.containerEnvMergeField.HasFocus() {
		d.focusElement = createContainerUnsetEnvFieldFocus

		return
	}

	if d.containerUnsetEnvField.HasFocus() {
		d.focusElement = createContainerEnvHostFieldFocus

		return
	}

	if d.containerEnvHostField.HasFocus() {
		d.focusElement = createContainerUnsetEnvAllFieldFocus

		return
	}

	if d.containerUnsetEnvAllField.HasFocus() {
		d.focusElement = createContainerUmaskFieldFocus

		return
	}

	d.focusElement = createContainerFormFocus
}

func (d *ContainerCreateDialog) setUserGroupsPageNextFocus() {
	if d.containerUserField.HasFocus() {
		d.focusElement = createContainerHostUsersFieldFocus

		return
	}

	if d.containerHostUsersField.HasFocus() {
		d.focusElement = createContainerPasswdEntryFieldFocus

		return
	}

	if d.containerPasswdEntryField.HasFocus() {
		d.focusElement = createContainerGroupEntryFieldFocus

		return
	}

	d.focusElement = createContainerFormFocus
}

func (d *ContainerCreateDialog) setSecurityOptionsPageNextFocus() {
	if d.containerSecLabelField.HasFocus() {
		d.focusElement = createContainerApprarmorFieldFocus

		return
	}

	if d.containerSecApparmorField.HasFocus() {
		d.focusElement = createContainerSeccompFeildFocus

		return
	}

	if d.containerSeccompField.HasFocus() {
		d.focusElement = createcontainerSecMaskFieldFocus

		return
	}

	if d.containerSecMaskField.HasFocus() {
		d.focusElement = createcontainerSecUnmaskFieldFocus

		return
	}

	if d.containerSecUnmaskField.HasFocus() {
		d.focusElement = createcontainerSecNoNewPrivFieldFocus

		return
	}

	d.focusElement = createContainerFormFocus
}

func (d *ContainerCreateDialog) setNetworkSettingsPageNextFocus() {
	if d.containerHostnameField.HasFocus() {
		d.focusElement = createContainerIPAddrFieldFocus

		return
	}

	if d.containerIPAddrField.HasFocus() {
		d.focusElement = createContainerMacAddrFieldFocus

		return
	}

	if d.containerMacAddrField.HasFocus() {
		d.focusElement = createContainerNetworkFieldFocus

		return
	}

	d.focusElement = createContainerFormFocus
}

func (d *ContainerCreateDialog) setDNSSettingsPageNextFocus() {
	if d.containerDNSServersField.HasFocus() {
		d.focusElement = createContainerDNSOptionsFieldFocus

		return
	}

	if d.containerDNSOptionsField.HasFocus() {
		d.focusElement = createContainerDNSSearchFieldFocus

		return
	}

	d.focusElement = createContainerFormFocus
}

func (d *ContainerCreateDialog) setResourceSettingsPageNextFocus() { //nolint:cyclop,dupl
	if d.containerMemoryField.HasFocus() {
		d.focusElement = createContainerMemoryReservatoinFieldFocus

		return
	}

	if d.containerMemoryReservationField.HasFocus() {
		d.focusElement = createContainerMemorySwapFieldFocus

		return
	}

	if d.containerMemorySwapField.HasFocus() {
		d.focusElement = createcontainerMemorySwappinessFieldFocus

		return
	}

	if d.containerMemorySwappinessField.HasFocus() {
		d.focusElement = createContainerCPUsFieldFocus

		return
	}

	if d.containerCPUsField.HasFocus() {
		d.focusElement = createContainerCPUSharesFieldFocus

		return
	}

	if d.containerCPUSharesField.HasFocus() {
		d.focusElement = createContainerCPUPeriodFieldFocus

		return
	}

	if d.containerCPUPeriodField.HasFocus() {
		d.focusElement = createContainerCPURtPeriodFieldFocus

		return
	}

	if d.containerCPURtPeriodField.HasFocus() {
		d.focusElement = createContainerCPUQuotaFieldFocus

		return
	}

	if d.containerCPUQuotaField.HasFocus() {
		d.focusElement = createContainerCPURtRuntimeFeildFocus

		return
	}

	if d.containerCPURtRuntimeField.HasFocus() {
		d.focusElement = createContainerCPUSetCPUsFieldFocus

		return
	}

	if d.containerCPUSetCPUsField.HasFocus() {
		d.focusElement = createContainerCPUSetMemsFieldFocus

		return
	}

	if d.containerCPUSetMemsField.HasFocus() {
		d.focusElement = createContainerShmSizeFieldFocus

		return
	}

	if d.containerShmSizeField.HasFocus() {
		d.focusElement = createContainerShmSizeSystemdFieldFocus

		return
	}

	d.focusElement = createContainerFormFocus
}

func (d *ContainerCreateDialog) setHealthSettingsPageNextFocus() { //nolint:cyclop,dupl
	if d.containerHealthCmdField.HasFocus() {
		d.focusElement = createContainerHealthStartupCmdFieldFocus

		return
	}

	if d.containerHealthStartupCmdField.HasFocus() {
		d.focusElement = createContainerHealthLogDestFocus

		return
	}

	if d.containerHealthLogDestField.HasFocus() {
		d.focusElement = createContainerHealthMaxLogSizeFocus

		return
	}

	if d.containerHealthMaxLogSizeField.HasFocus() {
		d.focusElement = createContainerHealthMaxLogCountFocus

		return
	}

	if d.containerHealthMaxLogCountField.HasFocus() {
		d.focusElement = createContainerHealthOnFailureFieldFocus

		return
	}

	if d.containerHealthOnFailureField.HasFocus() {
		d.focusElement = createContainerHealthIntervalFieldFocus

		return
	}

	if d.containerHealthIntervalField.HasFocus() {
		d.focusElement = createContainerHealthStartupIntervalFieldFocus

		return
	}

	if d.containerHealthStartupIntervalField.HasFocus() {
		d.focusElement = createContainerHealthStartPeriodFieldFocus

		return
	}

	if d.containerHealthStartPeriodField.HasFocus() {
		d.focusElement = createContainerHealthRetriesFieldFocus

		return
	}

	if d.containerHealthRetriesField.HasFocus() {
		d.focusElement = createContainerHealthStartupRetriesFieldFocus

		return
	}

	if d.containerHealthStartupRetriesField.HasFocus() {
		d.focusElement = createContainerHealthStartupSuccessFieldFocus

		return
	}

	if d.containerHealthStartupSuccessField.HasFocus() {
		d.focusElement = createContainerHealthTimeoutFieldFocus

		return
	}

	if d.containerHealthTimeoutField.HasFocus() {
		d.focusElement = createContainerHealthStartupTimeoutFieldFocus

		return
	}

	d.focusElement = createContainerFormFocus
}

func (d *ContainerCreateDialog) setVolumeSettingsPageNextFocus() {
	if d.containerVolumeField.HasFocus() {
		d.focusElement = createContainerImageVolumeFieldFocus

		return
	}

	if d.containerImageVolumeField.HasFocus() {
		d.focusElement = createContainerMountFieldFocus

		return
	}

	d.focusElement = createContainerFormFocus
}

// ContainerCreateOptions returns new network options.
func (d *ContainerCreateDialog) ContainerCreateOptions() containers.CreateOptions { //nolint:cyclop,gocognit,gocyclo,maintidx,lll
	var (
		labels           []string
		imageID          string
		podID            string
		dnsServers       []string
		dnsOptions       []string
		dnsSearchDomains []string
		publish          []string
		expose           []string
		imageVolume      string
		selinuxOpts      []string
		envVars          []string
		envFile          []string
		envMerge         []string
		unsetEnv         []string
		hostUsers        []string
		secret           []string
	)

	for _, label := range strings.Split(d.containerLabelsField.GetText(), " ") {
		if label != "" {
			labels = append(labels, label)
		}
	}

	selectedImageIndex, _ := d.containerImageField.GetCurrentOption()
	if len(d.imageList) > 0 && selectedImageIndex > 0 {
		imageID = d.imageList[selectedImageIndex-1].ID
	}

	selectedPodIndex, _ := d.containerPodField.GetCurrentOption()
	if len(d.podList) > 0 && selectedPodIndex > 0 {
		podID = d.podList[selectedPodIndex-1].Id
	}

	// ports
	for _, p := range strings.Split(d.containerPortPublishField.GetText(), " ") {
		if p != "" {
			publish = append(publish, p)
		}
	}

	for _, e := range strings.Split(d.containerPortExposeField.GetText(), " ") {
		if e != "" {
			expose = append(expose, e)
		}
	}

	// DNS setting
	for _, dns := range strings.Split(d.containerDNSServersField.GetText(), " ") {
		if dns != "" {
			dnsServers = append(dnsServers, dns)
		}
	}

	for _, do := range strings.Split(d.containerDNSOptionsField.GetText(), " ") {
		if do != "" {
			dnsOptions = append(dnsOptions, do)
		}
	}

	for _, ds := range strings.Split(d.containerDNSSearchField.GetText(), " ") {
		if ds != "" {
			dnsSearchDomains = append(dnsSearchDomains, ds)
		}
	}

	_, imageVolume = d.containerImageVolumeField.GetCurrentOption()

	// security options
	for _, selinuxLabel := range strings.Split(d.containerSecLabelField.GetText(), " ") {
		if selinuxLabel != "" {
			selinuxOpts = append(selinuxOpts, selinuxLabel)
		}
	}

	// health check
	_, healthOnFailure := d.containerHealthOnFailureField.GetCurrentOption()

	// env vars
	for _, evar := range strings.Split(d.containerEnvVarsField.GetText(), " ") {
		if evar != "" {
			envVars = append(envVars, evar)
		}
	}

	// env file
	for _, efile := range strings.Split(d.containerEnvFileField.GetText(), " ") {
		if efile != "" {
			envFile = append(envFile, efile)
		}
	}

	// env merge
	for _, emerge := range strings.Split(d.containerEnvMergeField.GetText(), " ") {
		if emerge != "" {
			envMerge = append(envMerge, emerge)
		}
	}

	// unset env
	for _, eunset := range strings.Split(d.containerUnsetEnvField.GetText(), " ") {
		if eunset != "" {
			unsetEnv = append(unsetEnv, eunset)
		}
	}

	// host users
	for _, huser := range strings.Split(d.containerHostUsersField.GetText(), " ") {
		if huser != "" {
			hostUsers = append(hostUsers, huser)
		}
	}

	// secret
	for _, sec := range strings.Split(d.containerSecretField.GetText(), " ") {
		if sec != "" {
			secret = append(secret, sec)
		}
	}

	_, network := d.containerNetworkField.GetCurrentOption()
	opts := containers.CreateOptions{
		Name:                  strings.TrimSpace(d.containerNameField.GetText()),
		Command:               strings.TrimSpace(d.containerCommandField.GetText()),
		Image:                 imageID,
		Pod:                   podID,
		Labels:                labels,
		Remove:                d.containerRemoveField.IsChecked(),
		Privileged:            d.containerPrivilegedField.IsChecked(),
		Timeout:               strings.TrimSpace(d.containerTimeoutField.GetText()),
		TTY:                   d.containerTtyField.IsChecked(),
		Detach:                d.containerDetachField.IsChecked(),
		Interactive:           d.containerInteractiveField.IsChecked(),
		Secret:                secret,
		WorkDir:               strings.TrimSpace(d.containerWorkDirField.GetText()),
		EnvVars:               envVars,
		EnvFile:               envFile,
		EnvMerge:              envMerge,
		UnsetEnv:              unsetEnv,
		EnvHost:               d.containerEnvHostField.IsChecked(),
		UnsetEnvAll:           d.containerUnsetEnvAllField.IsChecked(),
		Umask:                 strings.TrimSpace(d.containerUmaskField.GetText()),
		User:                  strings.TrimSpace(d.containerUserField.GetText()),
		HostUsers:             hostUsers,
		PasswdEntry:           strings.TrimSpace(d.containerPasswdEntryField.GetText()),
		GroupEntry:            strings.TrimSpace(d.containerGroupEntryField.GetText()),
		Hostname:              strings.TrimSpace(d.containerHostnameField.GetText()),
		MacAddress:            strings.TrimSpace(d.containerMacAddrField.GetText()),
		IPAddress:             strings.TrimSpace(d.containerIPAddrField.GetText()),
		Network:               network,
		Publish:               publish,
		Expose:                expose,
		PublishAll:            d.ContainerPortPublishAllField.IsChecked(),
		DNSServer:             dnsServers,
		DNSOptions:            dnsOptions,
		DNSSearchDomain:       dnsSearchDomains,
		Volume:                strings.TrimSpace(d.containerVolumeField.GetText()),
		ImageVolume:           imageVolume,
		Mount:                 strings.TrimSpace(d.containerMountField.GetText()),
		SelinuxOpts:           selinuxOpts,
		ApparmorProfile:       strings.TrimSpace(d.containerSecApparmorField.GetText()),
		Seccomp:               strings.TrimSpace(d.containerSeccompField.GetText()),
		SecNoNewPriv:          d.containerSecNoNewPrivField.IsChecked(),
		SecMask:               strings.TrimSpace(d.containerSecMaskField.GetText()),
		SecUnmask:             strings.TrimSpace(d.containerSecUnmaskField.GetText()),
		HealthCmd:             strings.TrimSpace(d.containerHealthCmdField.GetText()),
		HealthInterval:        strings.TrimSpace(d.containerHealthIntervalField.GetText()),
		HealthRetries:         strings.TrimSpace(d.containerHealthRetriesField.GetText()),
		HealthStartPeroid:     strings.TrimSpace(d.containerHealthStartPeriodField.GetText()),
		HealthTimeout:         strings.TrimSpace(d.containerHealthTimeoutField.GetText()),
		HealthOnFailure:       strings.TrimSpace(healthOnFailure),
		HealthStartupCmd:      strings.TrimSpace(d.containerHealthStartupCmdField.GetText()),
		HealthStartupInterval: strings.TrimSpace(d.containerHealthStartupIntervalField.GetText()),
		HealthStartupRetries:  strings.TrimSpace(d.containerHealthStartupRetriesField.GetText()),
		HealthStartupSuccess:  strings.TrimSpace(d.containerHealthStartupSuccessField.GetText()),
		HealthStartupTimeout:  strings.TrimSpace(d.containerHealthStartupTimeoutField.GetText()),
		HealthLogDestination:  strings.TrimSpace(d.containerHealthLogDestField.GetText()),
		HealthMaxLogSize:      strings.TrimSpace(d.containerHealthMaxLogSizeField.GetText()),
		HealthMaxLogCount:     strings.TrimSpace(d.containerHealthMaxLogCountField.GetText()),
		Memory:                strings.TrimSpace(d.containerMemoryField.GetText()),
		MemoryReservation:     strings.TrimSpace(d.containerMemoryReservationField.GetText()),
		MemorySwap:            strings.TrimSpace(d.containerMemorySwapField.GetText()),
		MemorySwappiness:      strings.TrimSpace(d.containerMemorySwappinessField.GetText()),
		CPUs:                  strings.TrimSpace(d.containerCPUsField.GetText()),
		CPUShares:             strings.TrimSpace(d.containerCPUSharesField.GetText()),
		CPUPeriod:             strings.TrimSpace(d.containerCPUPeriodField.GetText()),
		CPURtPeriod:           strings.TrimSpace(d.containerCPURtPeriodField.GetText()),
		CPUQuota:              strings.TrimSpace(d.containerCPUQuotaField.GetText()),
		CPURtRuntime:          strings.TrimSpace(d.containerCPURtRuntimeField.GetText()),
		CPUSetCPUs:            strings.TrimSpace(d.containerCPUSetCPUsField.GetText()),
		CPUSetMems:            strings.TrimSpace(d.containerCPUSetMemsField.GetText()),
		SHMSize:               strings.TrimSpace(d.containerShmSizeField.GetText()),
		SHMSizeSystemd:        strings.TrimSpace(d.containerShmSizeSystemdField.GetText()),
		NamespaceCgroup:       strings.TrimSpace(d.containerNamespaceCgroupField.GetText()),
		NamespaceIpc:          strings.TrimSpace(d.containerNamespaceIpcField.GetText()),
		NamespacePid:          strings.TrimSpace(d.containerNamespacePidField.GetText()),
		NamespaceUser:         strings.TrimSpace(d.containerNamespaceUserField.GetText()),
		NamespaceUts:          strings.TrimSpace(d.containerNamespaceUtsField.GetText()),
		NamespaceUidmap:       strings.TrimSpace(d.containerNamespaceUidmapField.GetText()),
		NamespaceSubuidName:   strings.TrimSpace(d.containerNamespaceSubuidNameField.GetText()),
		NamespaceGidmap:       strings.TrimSpace(d.containerNamespaceGidmapField.GetText()),
		NamespaceSubgidName:   strings.TrimSpace(d.containerNamespaceSubgidNameField.GetText()),
	}

	return opts
}
