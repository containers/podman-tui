package cntdialogs

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/containers"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	cntRestoreContainersFocus = 0 + iota
	cntRestorePodFocus
	cntRestoreNameFocus
	cntRestorePublishPortsFocus
	cntRestoreImportArchiveFocus
	cntRestoreKeepFocus
	cntRestoreIgnoreStaticIPFocus
	cntRestoreIgnoreStaticMACFocus
	cntRestoreFileLocksFocus
	cntRestorePrintStatsFocus
	cntRestoreTCPEstablishedFocus
	cntRestoreIgnroeVolumesFocus
	cntRestoreIgnoreRootFsFocus
	cntRestoreFormFocus
)

const (
	cntRestoreDialogLabelWidth            = 14
	cntRestoreDialogPadding               = 2
	cntRestoreDialogChkGroupColTwoWidth   = 18
	cntRestoreDialogChkGroupColThreeWidth = 20
	cntRestoreDialogChkGroupColFourWidth  = 16
	cntRestoreDialogMaxWidth              = cntRestoreDialogLabelWidth +
		cntRestoreDialogChkGroupColTwoWidth +
		cntRestoreDialogChkGroupColThreeWidth +
		cntRestoreDialogChkGroupColFourWidth + (6 * cntRestoreDialogPadding) //nolint:gomnd
	cntRestoreDialogMaxHeight        = 17
	cntRestoreDialogSingleFieldWidth = cntRestoreDialogMaxWidth -
		cntRestoreDialogLabelWidth - (2 * cntRestoreDialogPadding) //nolint:gomnd
)

// ContainerRestoreDialog implements container restore dialog primitive.
type ContainerRestoreDialog struct {
	*tview.Box
	layout          *tview.Flex
	containers      *tview.DropDown
	pods            *tview.DropDown
	name            *tview.InputField
	publishPorts    *tview.InputField
	importArchive   *tview.InputField
	ignoreRootFS    *tview.Checkbox
	ignoreVolumes   *tview.Checkbox
	ignoreStaticIP  *tview.Checkbox
	ignoreStaticMAC *tview.Checkbox
	keep            *tview.Checkbox
	tcpEstablished  *tview.Checkbox
	fileLocks       *tview.Checkbox
	printStats      *tview.Checkbox
	form            *tview.Form
	display         bool
	focusElement    int
	restoreHandler  func()
	cancelHandler   func()
}

// NewContainerRestoreDialog returns new container dialog primitive.
func NewContainerRestoreDialog() *ContainerRestoreDialog {
	dialog := &ContainerRestoreDialog{
		Box:             tview.NewBox(),
		layout:          tview.NewFlex(),
		containers:      tview.NewDropDown(),
		pods:            tview.NewDropDown(),
		name:            tview.NewInputField(),
		publishPorts:    tview.NewInputField(),
		importArchive:   tview.NewInputField(),
		keep:            tview.NewCheckbox(),
		ignoreStaticIP:  tview.NewCheckbox(),
		ignoreStaticMAC: tview.NewCheckbox(),
		fileLocks:       tview.NewCheckbox(),
		printStats:      tview.NewCheckbox(),
		ignoreRootFS:    tview.NewCheckbox(),
		ignoreVolumes:   tview.NewCheckbox(),
		tcpEstablished:  tview.NewCheckbox(),
		form:            tview.NewForm(),
	}

	fgColor := style.DialogFgColor
	bgColor := style.DialogBgColor
	ddUnselectedStyle := style.DropDownUnselected
	ddselectedStyle := style.DropDownSelected
	inputFieldBgColor := style.InputFieldBgColor

	// containers
	containersLabel := fmt.Sprintf("[:#%x:b]CONTAINER ID:[:-:-]", style.DialogBorderColor.Hex())

	dialog.containers.SetLabel(containersLabel)
	dialog.containers.SetLabelWidth(cntRestoreDialogLabelWidth)
	dialog.containers.SetFieldWidth(cntRestoreDialogSingleFieldWidth)
	dialog.containers.SetBackgroundColor(bgColor)
	dialog.containers.SetLabelColor(fgColor)
	dialog.containers.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	dialog.containers.SetFieldBackgroundColor(inputFieldBgColor)
	dialog.SetContainers(nil)
	dialog.containers.SetCurrentOption(0)

	// pod
	dialog.pods.SetLabel("pod:")
	dialog.pods.SetLabelWidth(cntRestoreDialogLabelWidth)
	dialog.pods.SetFieldWidth(cntRestoreDialogSingleFieldWidth)
	dialog.pods.SetBackgroundColor(bgColor)
	dialog.pods.SetLabelColor(fgColor)
	dialog.pods.SetLabelColor(fgColor)
	dialog.pods.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	dialog.pods.SetFieldBackgroundColor(inputFieldBgColor)
	dialog.SetPods(nil)
	dialog.pods.SetCurrentOption(0)

	// name
	dialog.name.SetLabel("name:")
	dialog.name.SetLabelWidth(cntRestoreDialogLabelWidth)
	dialog.name.SetBackgroundColor(bgColor)
	dialog.name.SetLabelColor(fgColor)
	dialog.name.SetFieldBackgroundColor(inputFieldBgColor)

	// Publish ports
	publishLabel := "publish:"
	publishLabelWidth := len(publishLabel) + cntRestoreDialogPadding
	publishPortsLabel := fmt.Sprintf("%*s ",
		publishLabelWidth, publishLabel)

	dialog.publishPorts.SetLabel(publishPortsLabel)
	dialog.publishPorts.SetBackgroundColor(bgColor)
	dialog.publishPorts.SetLabelColor(fgColor)
	dialog.publishPorts.SetFieldBackgroundColor(inputFieldBgColor)

	// Import
	dialog.importArchive.SetBackgroundColor(bgColor)
	dialog.importArchive.SetLabelColor(fgColor)
	dialog.importArchive.SetLabel("import:")
	dialog.importArchive.SetLabelWidth(cntRestoreDialogLabelWidth)
	dialog.importArchive.SetFieldBackgroundColor(inputFieldBgColor)

	// keep
	dialog.keep.SetLabel("keep:")
	dialog.keep.SetLabelWidth(cntRestoreDialogLabelWidth)
	dialog.keep.SetChecked(false)
	dialog.keep.SetBackgroundColor(bgColor)
	dialog.keep.SetLabelColor(fgColor)
	dialog.keep.SetFieldBackgroundColor(inputFieldBgColor)

	// ignoreStaticIP
	ignoreStaticIPLabel := fmt.Sprintf("%*s ",
		cntRestoreDialogChkGroupColTwoWidth, "ignore static IP:")

	dialog.ignoreStaticIP.SetLabel(ignoreStaticIPLabel)
	dialog.ignoreStaticIP.SetChecked(false)
	dialog.ignoreStaticIP.SetBackgroundColor(bgColor)
	dialog.ignoreStaticIP.SetLabelColor(fgColor)
	dialog.ignoreStaticIP.SetFieldBackgroundColor(inputFieldBgColor)

	// ignoreStaticMAC
	ignoreStaticMACLabel := fmt.Sprintf("%*s ",
		cntRestoreDialogChkGroupColThreeWidth, "ignore static MAC:")

	dialog.ignoreStaticMAC.SetLabel(ignoreStaticMACLabel)
	dialog.ignoreStaticMAC.SetChecked(false)
	dialog.ignoreStaticMAC.SetBackgroundColor(bgColor)
	dialog.ignoreStaticMAC.SetLabelColor(fgColor)
	dialog.ignoreStaticMAC.SetFieldBackgroundColor(inputFieldBgColor)

	// fileLocks
	fileLocksLabel := fmt.Sprintf("%*s ",
		cntRestoreDialogChkGroupColFourWidth, "file locks:")

	dialog.fileLocks.SetLabel(fileLocksLabel)
	dialog.fileLocks.SetChecked(false)
	dialog.fileLocks.SetBackgroundColor(bgColor)
	dialog.fileLocks.SetLabelColor(fgColor)
	dialog.fileLocks.SetFieldBackgroundColor(inputFieldBgColor)

	// printStats
	dialog.printStats.SetLabel("print Stats: ")
	dialog.printStats.SetLabelWidth(cntRestoreDialogLabelWidth)
	dialog.printStats.SetChecked(false)
	dialog.printStats.SetBackgroundColor(bgColor)
	dialog.printStats.SetLabelColor(fgColor)
	dialog.printStats.SetFieldBackgroundColor(inputFieldBgColor)

	// tcpEstablished
	tcpEstablishedLabel := fmt.Sprintf("%*s ",
		cntRestoreDialogChkGroupColTwoWidth, "tcp established:")

	dialog.tcpEstablished.SetLabel(tcpEstablishedLabel)
	dialog.tcpEstablished.SetChecked(false)
	dialog.tcpEstablished.SetBackgroundColor(bgColor)
	dialog.tcpEstablished.SetLabelColor(fgColor)
	dialog.tcpEstablished.SetFieldBackgroundColor(inputFieldBgColor)

	// ignoreVolumes
	ignoreVolumesLabel := fmt.Sprintf("%*s ",
		cntRestoreDialogChkGroupColThreeWidth, "ignore volumes:")

	dialog.ignoreVolumes.SetLabel(ignoreVolumesLabel)
	dialog.ignoreVolumes.SetChecked(false)
	dialog.ignoreVolumes.SetBackgroundColor(bgColor)
	dialog.ignoreVolumes.SetLabelColor(fgColor)
	dialog.ignoreVolumes.SetFieldBackgroundColor(inputFieldBgColor)

	// ignoreRootFS
	ignoreRootFSLabel := fmt.Sprintf("%*s ",
		cntRestoreDialogChkGroupColFourWidth, "ignore rootfs:")

	dialog.ignoreRootFS.SetLabel(ignoreRootFSLabel)
	dialog.ignoreRootFS.SetChecked(false)
	dialog.ignoreRootFS.SetBackgroundColor(bgColor)
	dialog.ignoreRootFS.SetLabelColor(fgColor)
	dialog.ignoreRootFS.SetFieldBackgroundColor(inputFieldBgColor)

	// form
	dialog.form.AddButton("Cancel", nil)
	dialog.form.AddButton("Restore", nil)
	dialog.form.SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(bgColor)
	dialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// layout
	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	// layout row #one
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	layout.AddItem(dialog.containers, 0, 1, true)
	// layout row #two
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	layout.AddItem(dialog.pods, 0, 1, true)

	// layout row #three
	row := tview.NewFlex().SetDirection(tview.FlexColumn)
	row.AddItem(dialog.name, 0, 1, true)
	row.AddItem(dialog.publishPorts, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	layout.AddItem(row, 0, 1, true)

	// layout row #four
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	layout.AddItem(dialog.importArchive, 0, 1, true)

	// layout row #five
	row = tview.NewFlex().SetDirection(tview.FlexColumn)

	row.AddItem(dialog.keep,
		cntRestoreDialogLabelWidth+cntRestoreDialogPadding,
		0, true)
	row.AddItem(dialog.ignoreStaticIP,
		cntRestoreDialogChkGroupColTwoWidth+cntRestoreDialogPadding,
		0, true)
	row.AddItem(dialog.ignoreStaticMAC,
		cntRestoreDialogChkGroupColThreeWidth+cntRestoreDialogPadding,
		0, true)
	row.AddItem(dialog.fileLocks,
		cntRestoreDialogChkGroupColFourWidth+cntRestoreDialogPadding,
		0, true)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	layout.AddItem(row, 0, 1, true)

	// layout row #six
	row = tview.NewFlex().SetDirection(tview.FlexColumn)

	row.AddItem(dialog.printStats,
		cntRestoreDialogLabelWidth+cntRestoreDialogPadding,
		0, true)
	row.AddItem(dialog.tcpEstablished,
		cntRestoreDialogChkGroupColTwoWidth+cntRestoreDialogPadding,
		0, true)
	row.AddItem(dialog.ignoreVolumes,
		cntRestoreDialogChkGroupColThreeWidth+cntRestoreDialogPadding,
		0, true)
	row.AddItem(dialog.ignoreRootFS,
		cntRestoreDialogChkGroupColFourWidth+cntRestoreDialogPadding,
		0, true)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	layout.AddItem(row, 0, 1, true)

	mainOptsLayout := tview.NewFlex().SetDirection(tview.FlexColumn)

	mainOptsLayout.SetBackgroundColor(bgColor)
	mainOptsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	mainOptsLayout.AddItem(layout, 0, 1, true)
	mainOptsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)

	dialog.layout.SetDirection(tview.FlexRow)
	dialog.layout.SetBackgroundColor(bgColor)
	dialog.layout.SetBorder(true)
	dialog.layout.SetBorderColor(style.DialogBorderColor)
	dialog.layout.SetTitle("PODMAN CONTAINER RESTORE")
	dialog.layout.AddItem(mainOptsLayout, 0, 1, true)
	dialog.layout.AddItem(dialog.form, dialogs.DialogFormHeight, 0, true)

	return dialog
}

// Display displays this primitive.
func (d *ContainerRestoreDialog) Display() {
	d.display = true
	d.focusElement = cntRestoreContainersFocus

	d.containers.SetCurrentOption(0)
	d.pods.SetCurrentOption(0)
	d.name.SetText("")
	d.publishPorts.SetText("")
	d.importArchive.SetText("")
	d.keep.SetChecked(false)
	d.ignoreStaticIP.SetChecked(false)
	d.ignoreStaticMAC.SetChecked(false)
	d.fileLocks.SetChecked(false)
	d.printStats.SetChecked(false)
	d.tcpEstablished.SetChecked(false)
	d.ignoreVolumes.SetChecked(false)
	d.ignoreRootFS.SetChecked(false)
}

// IsDisplay returns true if this primitive is shown.
func (d *ContainerRestoreDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *ContainerRestoreDialog) Hide() {
	d.display = false
}

// HasFocus returns whether or not this primitive has focus.
func (d *ContainerRestoreDialog) HasFocus() bool {
	for _, primitive := range d.getInnerPrimitives() {
		if primitive.HasFocus() {
			return true
		}
	}

	if d.layout.HasFocus() || d.Box.HasFocus() {
		return true
	}

	return false
}

// Focus is called when this primitive receives focus.
func (d *ContainerRestoreDialog) Focus(delegate func(p tview.Primitive)) {
	// all priviteves that can accept inputs
	if d.focusElement != cntRestoreFormFocus {
		primitives := d.getInnerPrimitives()
		delegate(primitives[d.focusElement])

		return
	}

	button := d.form.GetButton(d.form.GetButtonCount() - 1)

	button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == utils.SwitchFocusKey.Key {
			d.focusElement = cntRestoreContainersFocus

			d.Focus(delegate)
			d.form.SetFocus(0)

			return nil
		}

		return event
	})

	delegate(d.form)
}

// InputHandler returns input handler function for this primitive.
func (d *ContainerRestoreDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) { //nolint:cyclop,lll
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("container restore dialog: event %v received", event)

		if event.Key() == utils.CloseDialogKey.Key {
			if !d.containers.HasFocus() && !d.pods.HasFocus() {
				d.cancelHandler()
			}
		}

		if event.Key() == utils.SwitchFocusKey.Key {
			if d.focusElement != cntRestoreFormFocus {
				d.setFocusElement()
			}
		}

		// all priviteves that can accept inputs
		for _, primitive := range d.getInnerPrimitives() {
			if primitive.HasFocus() {
				if d.containers.HasFocus() || d.pods.HasFocus() {
					event = utils.ParseKeyEventKey(event)
				}

				if handler := primitive.InputHandler(); handler != nil {
					handler(event, setFocus)

					return
				}
			}
		}
	})
}

// SetRect set rects for this primitive.
func (d *ContainerRestoreDialog) SetRect(x, y, width, height int) {
	if width > cntRestoreDialogMaxWidth {
		emptySpace := (width - cntRestoreDialogMaxWidth) / 2 //nolint:gomnd
		x += emptySpace
		width = cntRestoreDialogMaxWidth
	}

	if height > cntRestoreDialogMaxHeight {
		emptySpace := (height - cntRestoreDialogMaxHeight) / 2 //nolint:gomnd
		y += emptySpace
		height = cntRestoreDialogMaxHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive into the screen.
func (d *ContainerRestoreDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.Box.DrawForSubclass(screen, d)

	x, y, width, height := d.Box.GetInnerRect()

	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetRestoreFunc sets form restore button selected function.
func (d *ContainerRestoreDialog) SetRestoreFunc(handler func()) *ContainerRestoreDialog {
	d.restoreHandler = handler
	restoreButton := d.form.GetButton(d.form.GetButtonCount() - 1)

	restoreButton.SetSelectedFunc(handler)

	return d
}

// SetCancelFunc sets form cancel button selected function.
func (d *ContainerRestoreDialog) SetCancelFunc(handler func()) *ContainerRestoreDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:gomnd

	cancelButton.SetSelectedFunc(handler)

	return d
}

func (d *ContainerRestoreDialog) setFocusElement() {
	if d.focusElement < cntRestoreFormFocus {
		d.focusElement++

		return
	}

	d.focusElement = cntRestoreContainersFocus
}

func (d *ContainerRestoreDialog) getInnerPrimitives() []tview.Primitive {
	// the item sort is important to be same as focus element number
	return []tview.Primitive{
		d.containers,
		d.pods,
		d.name,
		d.publishPorts,
		d.importArchive,
		d.keep,
		d.ignoreStaticIP,
		d.ignoreStaticMAC,
		d.fileLocks,
		d.printStats,
		d.tcpEstablished,
		d.ignoreVolumes,
		d.ignoreRootFS,
		d.form,
	}
}

// SetContainers sets containers dropdown options.
func (d *ContainerRestoreDialog) SetContainers(cnts [][]string) {
	emptyOptions := fmt.Sprintf("%*s", cntRestoreDialogSingleFieldWidth, " ")
	cntOptions := []string{emptyOptions}

	for i := 0; i < len(cnts); i++ {
		cntInfo := fmt.Sprintf("%s (%s)", utils.GetIDWithLimit(cnts[i][0]), cnts[i][1])
		cntInfoOption := fmt.Sprintf("%-*s", cntRestoreDialogSingleFieldWidth, cntInfo)
		cntOptions = append(cntOptions, cntInfoOption)
	}

	d.containers.SetOptions(cntOptions, nil)
}

// SetPods sets pods dropdown options.
func (d *ContainerRestoreDialog) SetPods(pods [][]string) {
	emptyOptions := fmt.Sprintf("%*s", cntRestoreDialogSingleFieldWidth, " ")
	podOptions := []string{emptyOptions}

	for i := 0; i < len(pods); i++ {
		podInfo := fmt.Sprintf("%s (%s)", utils.GetIDWithLimit(pods[i][0]), pods[i][1])
		podInfoOption := fmt.Sprintf("%-*s", cntRestoreDialogSingleFieldWidth, podInfo)
		podOptions = append(podOptions, podInfoOption)
	}

	d.pods.SetOptions(podOptions, nil)
}

func (d *ContainerRestoreDialog) GetRestoreOptions() containers.CntRestoreOptions {
	var opts containers.CntRestoreOptions

	_, cntInfoString := d.containers.GetCurrentOption()
	if strings.TrimSpace(cntInfoString) != "" {
		opts.ContainerID = strings.Split(cntInfoString, " ")[0]
	}

	_, podInfoString := d.pods.GetCurrentOption()
	if strings.TrimSpace(podInfoString) != "" {
		opts.PodID = strings.Split(podInfoString, " ")[0]
	}

	opts.Name = strings.TrimSpace(d.name.GetText())

	publishPortsList := strings.TrimSpace(d.publishPorts.GetText())
	opts.Publish = strings.Split(publishPortsList, " ")

	importArchive := strings.TrimSpace(d.importArchive.GetText())
	if strings.Index(importArchive, "~") == 0 {
		importArchive, _ = utils.ResolveHomeDir(importArchive)
	}

	opts.Import = importArchive

	opts.Keep = d.keep.IsChecked()
	opts.IgnoreStaticIP = d.ignoreStaticIP.IsChecked()
	opts.IgnoreStaticMAC = d.ignoreStaticMAC.IsChecked()
	opts.FileLocks = d.fileLocks.IsChecked()
	opts.PrintStats = d.printStats.IsChecked()
	opts.TCPEstablished = d.tcpEstablished.IsChecked()
	opts.IgnoreVolumes = d.ignoreVolumes.IsChecked()
	opts.IgnoreRootfs = d.ignoreRootFS.IsChecked()

	return opts
}
