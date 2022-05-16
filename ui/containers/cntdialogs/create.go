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
	"github.com/containers/podman-tui/ui/utils"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	containerCreateDialogMaxWidth = 80
	containerCreateDialogHeight   = 19
)

const (
	createContainerFormFocus = 0 + iota
	createCategoriesFocus
	createCategoryPagesFocus
	createContainerNameFieldFocus
	createContainerImageFieldFocus
	createcontainerPodFieldFocis
	createContainerLabelsFieldFocus
	createContainerRemoveFieldFocus
	createContainerSelinuxLabelFieldFocus
	createContainerApprarmorFieldFocus
	createContainerSeccompFeildFocus
	createContainerMaskFieldFocus
	createContainerUnmaskFieldFocus
	createContainerNoNewPrivFieldFocus
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
)

const (
	basicInfoPageIndex = 0 + iota
	dnsPageIndex
	networkingPageIndex
	portPageIndex
	securityOptsPageIndex
	volumePageIndex
)

// ContainerCreateDialog implements container create dialog
type ContainerCreateDialog struct {
	*tview.Box
	layout                       *tview.Flex
	categoryLabels               []string
	categories                   *tview.TextView
	categoryPages                *tview.Pages
	basicInfoPage                *tview.Flex
	securityOptsPage             *tview.Flex
	portPage                     *tview.Flex
	networkingPage               *tview.Flex
	dnsPage                      *tview.Flex
	volumePage                   *tview.Flex
	form                         *tview.Form
	display                      bool
	activePageIndex              int
	focusElement                 int
	imageList                    []images.ImageListReporter
	podList                      []*entities.ListPodsReport
	containerNameField           *tview.InputField
	containerImageField          *tview.DropDown
	containerPodField            *tview.DropDown
	containerLabelsField         *tview.InputField
	containerRemoveField         *tview.Checkbox
	containerSelinuxLabelField   *tview.InputField
	containerApparmorField       *tview.InputField
	containerSeccompField        *tview.InputField
	containerMaskField           *tview.InputField
	containerUnmaskField         *tview.InputField
	containerNoNewPrivField      *tview.Checkbox
	containerPortExposeField     *tview.InputField
	containerPortPublishField    *tview.InputField
	ContainerPortPublishAllField *tview.Checkbox
	containerHostnameField       *tview.InputField
	containerIPAddrField         *tview.InputField
	containerMacAddrField        *tview.InputField
	containerNetworkField        *tview.DropDown
	containerDNSServersField     *tview.InputField
	containerDNSOptionsField     *tview.InputField
	containerDNSSearchField      *tview.InputField
	containerVolumeField         *tview.DropDown
	containerImageVolumeField    *tview.DropDown
	cancelHandler                func()
	createHandler                func()
}

// NewContainerCreateDialog returns new container create dialog primitive ContainerCreateDialog
func NewContainerCreateDialog() *ContainerCreateDialog {
	containerDialog := ContainerCreateDialog{
		Box:              tview.NewBox(),
		layout:           tview.NewFlex().SetDirection(tview.FlexRow),
		categories:       tview.NewTextView(),
		categoryPages:    tview.NewPages(),
		basicInfoPage:    tview.NewFlex(),
		securityOptsPage: tview.NewFlex(),
		networkingPage:   tview.NewFlex(),
		dnsPage:          tview.NewFlex(),
		portPage:         tview.NewFlex(),
		volumePage:       tview.NewFlex(),
		form:             tview.NewForm(),
		categoryLabels: []string{
			"Basic Information",
			"DNS Settings",
			"Network Settings",
			"Ports Settings",
			"Security Options",
			"Volumes Settings"},
		activePageIndex:              0,
		display:                      false,
		containerNameField:           tview.NewInputField(),
		containerImageField:          tview.NewDropDown(),
		containerPodField:            tview.NewDropDown(),
		containerLabelsField:         tview.NewInputField(),
		containerRemoveField:         tview.NewCheckbox(),
		containerSelinuxLabelField:   tview.NewInputField(),
		containerApparmorField:       tview.NewInputField(),
		containerSeccompField:        tview.NewInputField(),
		containerMaskField:           tview.NewInputField(),
		containerUnmaskField:         tview.NewInputField(),
		containerNoNewPrivField:      tview.NewCheckbox(),
		containerPortExposeField:     tview.NewInputField(),
		containerPortPublishField:    tview.NewInputField(),
		ContainerPortPublishAllField: tview.NewCheckbox(),
		containerHostnameField:       tview.NewInputField(),
		containerIPAddrField:         tview.NewInputField(),
		containerMacAddrField:        tview.NewInputField(),
		containerNetworkField:        tview.NewDropDown(),
		containerDNSServersField:     tview.NewInputField(),
		containerDNSOptionsField:     tview.NewInputField(),
		containerDNSSearchField:      tview.NewInputField(),
		containerVolumeField:         tview.NewDropDown(),
		containerImageVolumeField:    tview.NewDropDown(),
	}

	bgColor := utils.Styles.ContainerCreateDialog.BgColor
	ddUnselectedStyle := utils.Styles.DropdownStyle.Unselected
	ddselectedStyle := utils.Styles.DropdownStyle.Selected
	inputFieldBgColor := utils.Styles.InputFieldPrimitive.BgColor

	containerDialog.categories.SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)
	containerDialog.categories.SetBackgroundColor(bgColor)
	containerDialog.categories.SetBorder(true)

	// basic information setup page
	basicInfoPageLabelWidth := 14
	// name field
	containerDialog.containerNameField.SetLabel("name:")
	containerDialog.containerNameField.SetLabelWidth(basicInfoPageLabelWidth)
	containerDialog.containerNameField.SetBackgroundColor(bgColor)
	containerDialog.containerNameField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerNameField.SetFieldBackgroundColor(inputFieldBgColor)

	// image field
	containerDialog.containerImageField.SetLabel("select image:")
	containerDialog.containerImageField.SetLabelWidth(basicInfoPageLabelWidth)
	containerDialog.containerImageField.SetBackgroundColor(bgColor)
	containerDialog.containerImageField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerImageField.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	containerDialog.containerImageField.SetFieldBackgroundColor(inputFieldBgColor)

	// pod field
	containerDialog.containerPodField.SetLabel("select pod:")
	containerDialog.containerPodField.SetLabelWidth(basicInfoPageLabelWidth)
	containerDialog.containerPodField.SetBackgroundColor(bgColor)
	containerDialog.containerPodField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerPodField.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	containerDialog.containerPodField.SetFieldBackgroundColor(inputFieldBgColor)

	// labels field
	containerDialog.containerLabelsField.SetLabel("labels:")
	containerDialog.containerLabelsField.SetLabelWidth(basicInfoPageLabelWidth)
	containerDialog.containerLabelsField.SetBackgroundColor(bgColor)
	containerDialog.containerLabelsField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerLabelsField.SetFieldBackgroundColor(inputFieldBgColor)

	// remove field
	containerDialog.containerRemoveField.SetLabel("remove container after exit ")
	//containerDialog.containerRemoveField.SetLabelWidth(basicInfoPageLabelWidth)
	containerDialog.containerRemoveField.SetBackgroundColor(bgColor)
	containerDialog.containerRemoveField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerRemoveField.SetChecked(true)
	containerDialog.containerRemoveField.SetFieldBackgroundColor(inputFieldBgColor)

	// security options page
	securityOptsLabelWidth := 10
	// selinux label
	containerDialog.containerSelinuxLabelField.SetLabel("Label:")
	containerDialog.containerSelinuxLabelField.SetLabelWidth(securityOptsLabelWidth)
	containerDialog.containerSelinuxLabelField.SetBackgroundColor(bgColor)
	containerDialog.containerSelinuxLabelField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerSelinuxLabelField.SetFieldBackgroundColor(inputFieldBgColor)

	// apparmor
	containerDialog.containerApparmorField.SetLabel("Apparmor:")
	containerDialog.containerApparmorField.SetLabelWidth(securityOptsLabelWidth)
	containerDialog.containerApparmorField.SetBackgroundColor(bgColor)
	containerDialog.containerApparmorField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerApparmorField.SetFieldBackgroundColor(inputFieldBgColor)

	// seccomp
	containerDialog.containerSeccompField.SetLabel("Seccomp:")
	containerDialog.containerSeccompField.SetLabelWidth(securityOptsLabelWidth)
	containerDialog.containerSeccompField.SetBackgroundColor(bgColor)
	containerDialog.containerSeccompField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerSeccompField.SetFieldBackgroundColor(inputFieldBgColor)

	// mask
	containerDialog.containerMaskField.SetLabel("Mask:")
	containerDialog.containerMaskField.SetLabelWidth(securityOptsLabelWidth)
	containerDialog.containerMaskField.SetBackgroundColor(bgColor)
	containerDialog.containerMaskField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerMaskField.SetFieldBackgroundColor(inputFieldBgColor)

	// unmask
	containerDialog.containerUnmaskField.SetLabel("Unmask:")
	containerDialog.containerUnmaskField.SetLabelWidth(securityOptsLabelWidth)
	containerDialog.containerUnmaskField.SetBackgroundColor(bgColor)
	containerDialog.containerUnmaskField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerUnmaskField.SetFieldBackgroundColor(inputFieldBgColor)

	// no-new-privileges
	containerDialog.containerNoNewPrivField.SetLabel("No new privileges ")
	containerDialog.containerNoNewPrivField.SetBackgroundColor(bgColor)
	containerDialog.containerNoNewPrivField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerNoNewPrivField.SetBackgroundColor(bgColor)
	containerDialog.containerNoNewPrivField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerNoNewPrivField.SetChecked(false)
	containerDialog.containerNoNewPrivField.SetFieldBackgroundColor(inputFieldBgColor)

	// networking setup page
	networkingPageLabelWidth := 13
	// hostname field
	containerDialog.containerHostnameField.SetLabel("hostname:")
	containerDialog.containerHostnameField.SetLabelWidth(networkingPageLabelWidth)
	containerDialog.containerHostnameField.SetBackgroundColor(bgColor)
	containerDialog.containerHostnameField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerHostnameField.SetFieldBackgroundColor(inputFieldBgColor)

	// IP field
	containerDialog.containerIPAddrField.SetLabel("IP address:")
	containerDialog.containerIPAddrField.SetLabelWidth(networkingPageLabelWidth)
	containerDialog.containerIPAddrField.SetBackgroundColor(bgColor)
	containerDialog.containerIPAddrField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerIPAddrField.SetFieldBackgroundColor(inputFieldBgColor)

	// mac field
	containerDialog.containerMacAddrField.SetLabel("MAC address:")
	containerDialog.containerMacAddrField.SetLabelWidth(networkingPageLabelWidth)
	containerDialog.containerMacAddrField.SetBackgroundColor(bgColor)
	containerDialog.containerMacAddrField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerMacAddrField.SetFieldBackgroundColor(inputFieldBgColor)

	// network field
	containerDialog.containerNetworkField.SetLabel("network:")
	containerDialog.containerNetworkField.SetLabelWidth(networkingPageLabelWidth)
	containerDialog.containerNetworkField.SetBackgroundColor(bgColor)
	containerDialog.containerNetworkField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerNetworkField.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	containerDialog.containerNetworkField.SetFieldBackgroundColor(inputFieldBgColor)

	// ports setup page
	portPageLabelWidth := 15
	// publish field
	containerDialog.containerPortPublishField.SetLabel("publish ports:")
	containerDialog.containerPortPublishField.SetLabelWidth(portPageLabelWidth)
	containerDialog.containerPortPublishField.SetBackgroundColor(bgColor)
	containerDialog.containerPortPublishField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerPortPublishField.SetFieldBackgroundColor(inputFieldBgColor)

	// expose field
	containerDialog.containerPortExposeField.SetLabel("expose ports:")
	containerDialog.containerPortExposeField.SetLabelWidth(portPageLabelWidth)
	containerDialog.containerPortExposeField.SetBackgroundColor(bgColor)
	containerDialog.containerPortExposeField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerPortExposeField.SetFieldBackgroundColor(inputFieldBgColor)

	// publish all field
	containerDialog.ContainerPortPublishAllField.SetLabel("publish all ")
	containerDialog.ContainerPortPublishAllField.SetLabelWidth(portPageLabelWidth)
	containerDialog.ContainerPortPublishAllField.SetBackgroundColor(bgColor)
	containerDialog.ContainerPortPublishAllField.SetLabelColor(tcell.ColorWhite)
	containerDialog.ContainerPortPublishAllField.SetChecked(false)
	containerDialog.ContainerPortPublishAllField.SetFieldBackgroundColor(inputFieldBgColor)

	// dns setup page
	dnsPageLabelWidth := 13
	// hostname field
	containerDialog.containerDNSServersField.SetLabel("DNS servers:")
	containerDialog.containerDNSServersField.SetLabelWidth(dnsPageLabelWidth)
	containerDialog.containerDNSServersField.SetBackgroundColor(bgColor)
	containerDialog.containerDNSServersField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerDNSServersField.SetFieldBackgroundColor(inputFieldBgColor)

	// IP field
	containerDialog.containerDNSOptionsField.SetLabel("DNS options:")
	containerDialog.containerDNSOptionsField.SetLabelWidth(dnsPageLabelWidth)
	containerDialog.containerDNSOptionsField.SetBackgroundColor(bgColor)
	containerDialog.containerDNSOptionsField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerDNSOptionsField.SetFieldBackgroundColor(inputFieldBgColor)

	// mac field
	containerDialog.containerDNSSearchField.SetLabel("DNS search:")
	containerDialog.containerDNSSearchField.SetLabelWidth(dnsPageLabelWidth)
	containerDialog.containerDNSSearchField.SetBackgroundColor(bgColor)
	containerDialog.containerDNSSearchField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerDNSSearchField.SetFieldBackgroundColor(inputFieldBgColor)

	// volume setup page
	volumePageLabelWidth := 14
	// volume
	containerDialog.containerVolumeField.SetLabel("Volume:")
	containerDialog.containerVolumeField.SetLabelWidth(volumePageLabelWidth)
	containerDialog.containerVolumeField.SetBackgroundColor(bgColor)
	containerDialog.containerVolumeField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerVolumeField.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	containerDialog.containerVolumeField.SetFieldBackgroundColor(inputFieldBgColor)

	// image volume
	containerDialog.containerImageVolumeField.SetLabel("Image volume:")
	containerDialog.containerImageVolumeField.SetLabelWidth(volumePageLabelWidth)
	containerDialog.containerImageVolumeField.SetBackgroundColor(bgColor)
	containerDialog.containerImageVolumeField.SetLabelColor(tcell.ColorWhite)
	containerDialog.containerImageVolumeField.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	containerDialog.containerImageVolumeField.SetFieldBackgroundColor(inputFieldBgColor)

	// category pages
	containerDialog.categoryPages.SetBackgroundColor(bgColor)
	containerDialog.categoryPages.SetBorder(true)

	// form
	containerDialog.form.SetBackgroundColor(bgColor)
	containerDialog.form.AddButton("Cancel", nil)
	containerDialog.form.AddButton("Create", nil)
	containerDialog.form.SetButtonsAlign(tview.AlignRight)
	containerDialog.form.SetButtonBackgroundColor(utils.Styles.ButtonPrimitive.BgColor)

	containerDialog.layout.AddItem(tview.NewBox().SetBackgroundColor(bgColor), 1, 0, true)
	containerDialog.setupLayout()
	containerDialog.layout.SetBackgroundColor(bgColor)
	containerDialog.layout.SetBorder(true)
	containerDialog.layout.SetTitle("PODMAN CONTAINER CREATE")
	containerDialog.layout.AddItem(containerDialog.form, dialogs.DialogFormHeight, 0, true)

	containerDialog.setActiveCategory(0)

	containerDialog.initCustomInputHanlers()
	return &containerDialog
}

func (d *ContainerCreateDialog) setupLayout() {
	bgColor := utils.Styles.ContainerCreateDialog.BgColor

	// basic info page
	d.basicInfoPage.SetDirection(tview.FlexRow)
	d.basicInfoPage.AddItem(d.containerNameField, 1, 0, true)
	d.basicInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.basicInfoPage.AddItem(d.containerImageField, 1, 0, true)
	d.basicInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.basicInfoPage.AddItem(d.containerPodField, 1, 0, true)
	d.basicInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.basicInfoPage.AddItem(d.containerLabelsField, 1, 0, true)
	d.basicInfoPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.basicInfoPage.AddItem(d.containerRemoveField, 1, 0, true)
	d.basicInfoPage.SetBackgroundColor(bgColor)

	// security options page
	d.securityOptsPage.SetDirection(tview.FlexRow)
	d.securityOptsPage.AddItem(d.containerSelinuxLabelField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.containerApparmorField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.containerSeccompField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.containerMaskField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.containerUnmaskField, 1, 0, true)
	d.securityOptsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.securityOptsPage.AddItem(d.containerNoNewPrivField, 1, 0, true)

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

	// port settings page
	d.portPage.SetDirection(tview.FlexRow)
	d.portPage.AddItem(d.containerPortPublishField, 1, 0, true)
	d.portPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.portPage.AddItem(d.ContainerPortPublishAllField, 1, 0, true)
	d.portPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.portPage.AddItem(d.containerPortExposeField, 1, 0, true)
	d.portPage.SetBackgroundColor(bgColor)

	// dns settings page
	d.dnsPage.SetDirection(tview.FlexRow)
	d.dnsPage.AddItem(d.containerDNSServersField, 1, 0, true)
	d.dnsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.dnsPage.AddItem(d.containerDNSOptionsField, 1, 0, true)
	d.dnsPage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.dnsPage.AddItem(d.containerDNSSearchField, 1, 0, true)
	d.dnsPage.SetBackgroundColor(bgColor)

	// volume settings page
	d.volumePage.SetDirection(tview.FlexRow)
	d.volumePage.AddItem(d.containerVolumeField, 1, 0, true)
	d.volumePage.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.volumePage.AddItem(d.containerImageVolumeField, 1, 0, true)
	d.volumePage.SetBackgroundColor(bgColor)

	// adding category pages
	d.categoryPages.AddPage(d.categoryLabels[basicInfoPageIndex], d.basicInfoPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[dnsPageIndex], d.dnsPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[networkingPageIndex], d.networkingPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[portPageIndex], d.portPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[securityOptsPageIndex], d.securityOptsPage, true, true)
	d.categoryPages.AddPage(d.categoryLabels[volumePageIndex], d.volumePage, true, true)

	// add it to layout.
	_, layoutWidth := utils.AlignStringListWidth(d.categoryLabels)
	layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	layout.AddItem(d.categories, layoutWidth+6, 0, true)
	layout.AddItem(d.categoryPages, 0, 1, true)
	layout.SetBackgroundColor(bgColor)

	d.layout.AddItem(layout, 0, 1, true)

}

// Display displays this primitive
func (d *ContainerCreateDialog) Display() {
	d.display = true
	d.initData()
	d.focusElement = createCategoryPagesFocus
}

// IsDisplay returns true if primitive is shown
func (d *ContainerCreateDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive
func (d *ContainerCreateDialog) Hide() {
	d.display = false
}

// HasFocus returns whether or not this primitive has focus
func (d *ContainerCreateDialog) HasFocus() bool {
	if d.categories.HasFocus() || d.categoryPages.HasFocus() {
		return true
	}

	return d.Box.HasFocus() || d.form.HasFocus()
}

// Focus is called when this primitive receives focus
func (d *ContainerCreateDialog) Focus(delegate func(p tview.Primitive)) {
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
				//d.pullSelectHandler()
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
	// basic info page
	case createContainerNameFieldFocus:
		delegate(d.containerNameField)
	case createContainerImageFieldFocus:
		delegate(d.containerImageField)
	case createcontainerPodFieldFocis:
		delegate(d.containerPodField)
	case createContainerLabelsFieldFocus:
		delegate(d.containerLabelsField)
	case createContainerRemoveFieldFocus:
		delegate(d.containerRemoveField)
	// security options page
	case createContainerSelinuxLabelFieldFocus:
		delegate(d.containerSelinuxLabelField)
	case createContainerApprarmorFieldFocus:
		delegate(d.containerApparmorField)
	case createContainerSeccompFeildFocus:
		delegate(d.containerSeccompField)
	case createContainerMaskFieldFocus:
		delegate(d.containerMaskField)
	case createContainerUnmaskFieldFocus:
		delegate(d.containerUnmaskField)
	case createContainerNoNewPrivFieldFocus:
		delegate(d.containerNoNewPrivField)
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
	// container volume dropdown
	d.containerVolumeField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = utils.ParseKeyEventKey(event)
		return event
	})
	// container image volume dropdown
	d.containerImageVolumeField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = utils.ParseKeyEventKey(event)
		return event
	})
}

// InputHandler returns input handler function for this primitive
func (d *ContainerCreateDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("container create dialog: event %v received", event)
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
func (d *ContainerCreateDialog) SetRect(x, y, width, height int) {

	if width > containerCreateDialogMaxWidth {
		emptySpace := (width - containerCreateDialogMaxWidth) / 2
		x = x + emptySpace
		width = containerCreateDialogMaxWidth
	}

	if height > containerCreateDialogHeight {
		emptySpace := (height - containerCreateDialogHeight) / 2
		y = y + emptySpace
		height = containerCreateDialogHeight
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

// SetCancelFunc sets form cancel button selected function
func (d *ContainerCreateDialog) SetCancelFunc(handler func()) *ContainerCreateDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2)
	cancelButton.SetSelectedFunc(handler)
	return d
}

// SetCreateFunc sets form create button selected function
func (d *ContainerCreateDialog) SetCreateFunc(handler func()) *ContainerCreateDialog {
	d.createHandler = handler
	enterButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	enterButton.SetSelectedFunc(handler)
	return d
}

func (d *ContainerCreateDialog) setActiveCategory(index int) {
	fgColor := utils.Styles.ContainerCreateDialog.FgColor
	bgColor := int(utils.Styles.ButtonPrimitive.BgColor)
	ctgTextColor := utils.GetColorName(fgColor)
	ctgBgColor := utils.GetColorName(tcell.Color(bgColor))

	d.activePageIndex = index
	d.categories.Clear()
	var ctgList []string
	alignedList, _ := utils.AlignStringListWidth(d.categoryLabels)
	for i := 0; i < len(alignedList); i++ {
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
		activePage = activePage + 1
		d.setActiveCategory(activePage)
		return
	}
	d.setActiveCategory(0)
}

func (d *ContainerCreateDialog) previousCategory() {
	activePage := d.activePageIndex
	if d.activePageIndex > 0 {
		activePage = activePage - 1
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
	for i := 0; i < len(d.imageList); i++ {
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
	for i := 0; i < len(podList); i++ {
		podOptions = append(podOptions, podList[i].Name)
	}

	// get available networks
	networkOptions := []string{""}
	networkList, _ := networks.List()
	for i := 0; i < len(networkList); i++ {
		networkOptions = append(networkOptions, networkList[i][1])
	}

	// get available volumes
	imageVolumeOptions := []string{"", "ignore", "tmpfs", "anonymous"}
	volumeOptions := []string{""}
	volList, _ := volumes.List()
	for i := 0; i < len(volList); i++ {
		volumeOptions = append(volumeOptions, volList[i].Name)
	}

	d.setActiveCategory(0)
	d.containerNameField.SetText("")
	d.containerImageField.SetOptions(imgOptions, nil)
	d.containerImageField.SetCurrentOption(0)
	d.containerPodField.SetOptions(podOptions, nil)
	d.containerPodField.SetCurrentOption(0)
	d.containerLabelsField.SetText("")
	d.containerRemoveField.SetChecked(false)
	d.containerSelinuxLabelField.SetText("")
	d.containerApparmorField.SetText("")
	d.containerSeccompField.SetText("")
	d.containerMaskField.SetText("")
	d.containerUnmaskField.SetText("")
	d.containerNoNewPrivField.SetChecked(false)
	d.containerPortPublishField.SetText("")
	d.ContainerPortPublishAllField.SetChecked(false)
	d.containerPortExposeField.SetText("")
	d.containerHostnameField.SetText("")
	d.containerIPAddrField.SetText("")
	d.containerMacAddrField.SetText("")
	d.containerNetworkField.SetOptions(networkOptions, nil)
	d.containerNetworkField.SetCurrentOption(0)
	d.containerDNSServersField.SetText("")
	d.containerDNSSearchField.SetText("")
	d.containerDNSOptionsField.SetText("")
	d.containerVolumeField.SetOptions(volumeOptions, nil)
	d.containerVolumeField.SetCurrentOption(0)
	d.containerImageVolumeField.SetOptions(imageVolumeOptions, nil)
	d.containerImageVolumeField.SetCurrentOption(0)
}

func (d *ContainerCreateDialog) setPortPageNextFocus() {
	if d.containerPortPublishField.HasFocus() {
		d.focusElement = createContainerPortPublishAllFieldFocus
	} else if d.ContainerPortPublishAllField.HasFocus() {
		d.focusElement = createContainerPortExposeFieldFocus
	} else {
		d.focusElement = createContainerFormFocus
	}
}

func (d *ContainerCreateDialog) setBasicInfoPageNextFocus() {
	if d.containerNameField.HasFocus() {
		d.focusElement = createContainerImageFieldFocus
	} else if d.containerImageField.HasFocus() {
		d.focusElement = createcontainerPodFieldFocis
	} else if d.containerPodField.HasFocus() {
		d.focusElement = createContainerLabelsFieldFocus
	} else if d.containerLabelsField.HasFocus() {
		d.focusElement = createContainerRemoveFieldFocus
	} else {
		d.focusElement = createContainerFormFocus
	}
}

func (d *ContainerCreateDialog) setSecurityOptionsPageNextFocus() {
	if d.containerSelinuxLabelField.HasFocus() {
		d.focusElement = createContainerApprarmorFieldFocus
	} else if d.containerApparmorField.HasFocus() {
		d.focusElement = createContainerSeccompFeildFocus
	} else if d.containerSeccompField.HasFocus() {
		d.focusElement = createContainerMaskFieldFocus
	} else if d.containerMaskField.HasFocus() {
		d.focusElement = createContainerUnmaskFieldFocus
	} else if d.containerUnmaskField.HasFocus() {
		d.focusElement = createContainerNoNewPrivFieldFocus
	} else {
		d.focusElement = createContainerFormFocus
	}
}

func (d *ContainerCreateDialog) setNetworkSettingsPageNextFocus() {
	if d.containerHostnameField.HasFocus() {
		d.focusElement = createContainerIPAddrFieldFocus
	} else if d.containerIPAddrField.HasFocus() {
		d.focusElement = createContainerMacAddrFieldFocus
	} else if d.containerMacAddrField.HasFocus() {
		d.focusElement = createContainerNetworkFieldFocus
	} else {
		d.focusElement = createContainerFormFocus
	}
}

func (d *ContainerCreateDialog) setDNSSettingsPageNextFocus() {
	if d.containerDNSServersField.HasFocus() {
		d.focusElement = createContainerDNSOptionsFieldFocus
	} else if d.containerDNSOptionsField.HasFocus() {
		d.focusElement = createContainerDNSSearchFieldFocus
	} else {
		d.focusElement = createContainerFormFocus
	}
}

func (d *ContainerCreateDialog) setVolumeSettingsPageNextFocus() {
	if d.containerVolumeField.HasFocus() {
		d.focusElement = createContainerImageVolumeFieldFocus
	} else {
		d.focusElement = createContainerFormFocus
	}
}

// ContainerCreateOptions returns new network options
func (d *ContainerCreateDialog) ContainerCreateOptions() containers.CreateOptions {
	var (
		labels           = make(map[string]string)
		imageID          string
		podID            string
		dnsServers       []string
		dnsOptions       []string
		dnsSearchDomains []string
		publish          []string
		expose           []string
		volume           string
		imageVolume      string
		selinuxOpts      []string
		apparmorProfile  string
		seccompProfile   string
		maskPaths        []string
		unmaskPaths      []string
	)
	for _, label := range strings.Split(d.containerLabelsField.GetText(), " ") {
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
	_, volume = d.containerVolumeField.GetCurrentOption()
	_, imageVolume = d.containerImageVolumeField.GetCurrentOption()

	// security options
	for _, selinuxLabel := range strings.Split(d.containerSelinuxLabelField.GetText(), " ") {
		if selinuxLabel != "" {
			selinuxOpts = append(selinuxOpts, selinuxLabel)
		}
	}
	apparmor := strings.TrimSpace(d.containerApparmorField.GetText())
	if apparmor != "" {
		apparmorProfile = apparmor
	}
	for _, maskPath := range strings.Split(d.containerMaskField.GetText(), ":") {
		if maskPath != "" {
			maskPaths = append(maskPaths, maskPath)
		}
	}
	for _, unmaskPath := range strings.Split(d.containerUnmaskField.GetText(), ":") {
		if unmaskPath != "" {
			unmaskPaths = append(unmaskPaths, unmaskPath)
		}
	}

	_, network := d.containerNetworkField.GetCurrentOption()
	opts := containers.CreateOptions{
		Name:            d.containerNameField.GetText(),
		Image:           imageID,
		Pod:             podID,
		Labels:          labels,
		Remove:          d.containerRemoveField.IsChecked(),
		Hostname:        d.containerHostnameField.GetText(),
		MacAddress:      d.containerMacAddrField.GetText(),
		IPAddress:       d.containerIPAddrField.GetText(),
		Network:         network,
		Publish:         publish,
		Expose:          expose,
		PublishAll:      d.ContainerPortPublishAllField.IsChecked(),
		DNSServer:       dnsServers,
		DNSOptions:      dnsOptions,
		DNSSearchDomain: dnsSearchDomains,
		Volume:          volume,
		ImageVolume:     imageVolume,
		SelinuxOpts:     selinuxOpts,
		ApparmorProfile: apparmorProfile,
		Seccomp:         seccompProfile,
		NoNewPriv:       d.containerNoNewPrivField.IsChecked(),
		Mask:            maskPaths,
		Unmask:          unmaskPaths,
	}
	return opts
}
