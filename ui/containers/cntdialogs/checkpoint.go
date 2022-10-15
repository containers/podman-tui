package cntdialogs

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/containers"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	cntCheckpointDialogMaxWidth  = 73
	cntCheckpointDialogMaxHeight = 15
)

const (
	cntCheckpointImageFocus = 0 + iota
	cntCheckpointExportFocus
	cntCheckpointFileLockFocus
	cntCheckpointIgnoreRootFsFocus
	cntCheckpointKeepFocus
	cntCheckpointLeaveRunningFocus
	cntCheckpointPreCheckpointFocus
	cntCheckpointPrintStatsFocus
	cntCheckpointTcpEstablishedFocus
	cntCheckpointWithPreviousFocus
	cntCheckpointFormFocus
)

// ContainerCheckpointDialog implements container checkpint dialog primitive.
type ContainerCheckpointDialog struct {
	*tview.Box
	layout            *tview.Flex
	containerInfo     *tview.TextView
	createImage       *tview.InputField
	export            *tview.InputField
	fileLock          *tview.Checkbox
	ignoreRootFS      *tview.Checkbox
	keep              *tview.Checkbox
	leaveRunning      *tview.Checkbox
	preCheckpoint     *tview.Checkbox
	printStats        *tview.Checkbox
	tcpEstablished    *tview.Checkbox
	withPrevious      *tview.Checkbox
	form              *tview.Form
	display           bool
	focusElement      int
	containerID       string
	checkpointHandler func()
	cancelHandler     func()
}

// NewContainerCheckpointDialog returns new container checkpoint dialog primitive.
func NewContainerCheckpointDialog() *ContainerCheckpointDialog {
	dialog := &ContainerCheckpointDialog{
		Box:            tview.NewBox(),
		layout:         tview.NewFlex(),
		containerInfo:  tview.NewTextView(),
		createImage:    tview.NewInputField(),
		export:         tview.NewInputField(),
		fileLock:       tview.NewCheckbox(),
		ignoreRootFS:   tview.NewCheckbox(),
		keep:           tview.NewCheckbox(),
		leaveRunning:   tview.NewCheckbox(),
		preCheckpoint:  tview.NewCheckbox(),
		printStats:     tview.NewCheckbox(),
		tcpEstablished: tview.NewCheckbox(),
		withPrevious:   tview.NewCheckbox(),
		form:           tview.NewForm(),
	}

	bgColor := utils.Styles.ContainerCheckpointDialog.BgColor
	fgColor := utils.Styles.ContainerCheckpointDialog.FgColor
	inputFieldBgColor := utils.Styles.InputFieldPrimitive.BgColor
	labelWidth := 14
	chkGroupFirstColLabelWidth := 14
	chkGroupSecondColLabelWidth := 16
	chkGroupThirdColLabelWidth := 16

	// containerInfo
	dialog.containerInfo.SetBackgroundColor(bgColor)
	dialog.containerInfo.SetTextColor(fgColor)
	dialog.containerInfo.SetDynamicColors(true)

	// createImage
	dialog.createImage.SetBackgroundColor(bgColor)
	dialog.createImage.SetLabelColor(fgColor)
	dialog.createImage.SetLabel("Create image:")
	dialog.createImage.SetLabelWidth(labelWidth)
	dialog.createImage.SetFieldBackgroundColor(inputFieldBgColor)

	// export
	dialog.export.SetBackgroundColor(bgColor)
	dialog.export.SetLabelColor(fgColor)
	dialog.export.SetLabel("Export:")
	dialog.export.SetLabelWidth(labelWidth)
	dialog.export.SetFieldBackgroundColor(inputFieldBgColor)

	// printStats
	dialog.printStats.SetLabel("Print stats:")
	dialog.printStats.SetLabelWidth(labelWidth)
	dialog.printStats.SetChecked(false)
	dialog.printStats.SetBackgroundColor(bgColor)
	dialog.printStats.SetLabelColor(fgColor)
	dialog.printStats.SetFieldBackgroundColor(inputFieldBgColor)

	// fileLock
	dialog.fileLock.SetLabel("File lock:")
	dialog.fileLock.SetLabelWidth(labelWidth)
	dialog.fileLock.SetChecked(false)
	dialog.fileLock.SetBackgroundColor(bgColor)
	dialog.fileLock.SetLabelColor(fgColor)
	dialog.fileLock.SetFieldBackgroundColor(inputFieldBgColor)

	// ignoreRootFS
	ignoreRootFSLabel := fmt.Sprintf("%*s ", chkGroupFirstColLabelWidth, "Ignore rootFS:")
	dialog.ignoreRootFS.SetLabel(ignoreRootFSLabel)
	dialog.ignoreRootFS.SetChecked(false)
	dialog.ignoreRootFS.SetBackgroundColor(bgColor)
	dialog.ignoreRootFS.SetLabelColor(fgColor)
	dialog.ignoreRootFS.SetFieldBackgroundColor(inputFieldBgColor)

	// keep
	keepLabel := fmt.Sprintf("%*s ", chkGroupFirstColLabelWidth, "Keep:")
	dialog.keep.SetLabel(keepLabel)
	dialog.keep.SetChecked(false)
	dialog.keep.SetBackgroundColor(bgColor)
	dialog.keep.SetLabelColor(fgColor)
	dialog.keep.SetFieldBackgroundColor(inputFieldBgColor)

	// tcpEstablished
	tcpEstablishedLabel := fmt.Sprintf("%*s ", chkGroupSecondColLabelWidth, "TCP established:")
	dialog.tcpEstablished.SetLabel(tcpEstablishedLabel)
	dialog.tcpEstablished.SetChecked(false)
	dialog.tcpEstablished.SetBackgroundColor(bgColor)
	dialog.tcpEstablished.SetLabelColor(fgColor)
	dialog.tcpEstablished.SetFieldBackgroundColor(inputFieldBgColor)

	// leaveRunning
	leaveRunningLabel := fmt.Sprintf("%*s ", chkGroupSecondColLabelWidth, "Leave running:")
	dialog.leaveRunning.SetLabel(leaveRunningLabel)
	dialog.leaveRunning.SetChecked(false)
	dialog.leaveRunning.SetBackgroundColor(bgColor)
	dialog.leaveRunning.SetLabelColor(fgColor)
	dialog.leaveRunning.SetFieldBackgroundColor(inputFieldBgColor)

	// preCheckpoint
	preCheckPointLabel := fmt.Sprintf("%*s ", chkGroupThirdColLabelWidth, "Pre checkpoint:")
	dialog.preCheckpoint.SetLabel(preCheckPointLabel)
	dialog.preCheckpoint.SetChecked(false)
	dialog.preCheckpoint.SetBackgroundColor(bgColor)
	dialog.preCheckpoint.SetLabelColor(fgColor)
	dialog.preCheckpoint.SetFieldBackgroundColor(inputFieldBgColor)

	// withPrevious
	withPreviousLabel := fmt.Sprintf("%*s ", chkGroupThirdColLabelWidth, "With previous:")
	dialog.withPrevious.SetLabel(withPreviousLabel)
	dialog.withPrevious.SetChecked(false)
	dialog.withPrevious.SetBackgroundColor(bgColor)
	dialog.withPrevious.SetLabelColor(fgColor)
	dialog.withPrevious.SetFieldBackgroundColor(inputFieldBgColor)

	// form
	dialog.form.AddButton(" Cancel ", nil)
	dialog.form.AddButton("Checkpoint", nil)
	dialog.form.SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(bgColor)
	dialog.form.SetButtonBackgroundColor(utils.Styles.ButtonPrimitive.BgColor)

	// layout
	optionsLayoutRow01 := tview.NewFlex().SetDirection(tview.FlexRow)
	optionsLayoutRow01.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayoutRow01.AddItem(dialog.containerInfo, 1, 0, true)
	optionsLayoutRow01.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayoutRow01.AddItem(dialog.createImage, 1, 0, true)
	optionsLayoutRow01.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayoutRow01.AddItem(dialog.export, 1, 0, true)

	optionsLayoutRow02 := tview.NewFlex().SetDirection(tview.FlexColumn)
	optionsLayoutRow02.AddItem(dialog.fileLock, labelWidth+2, 0, true)
	optionsLayoutRow02.AddItem(dialog.ignoreRootFS, 0, 1, true)
	optionsLayoutRow02.AddItem(dialog.tcpEstablished, 0, 1, true)
	optionsLayoutRow02.AddItem(dialog.preCheckpoint, 0, 1, true)

	optionsLayoutRow03 := tview.NewFlex().SetDirection(tview.FlexColumn)
	optionsLayoutRow03.AddItem(dialog.printStats, labelWidth+2, 1, true)
	optionsLayoutRow03.AddItem(dialog.keep, 0, 1, true)
	optionsLayoutRow03.AddItem(dialog.leaveRunning, 0, 1, true)
	optionsLayoutRow03.AddItem(dialog.withPrevious, 0, 1, true)

	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	layout.AddItem(optionsLayoutRow01, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	layout.AddItem(optionsLayoutRow02, 1, 0, true)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	layout.AddItem(optionsLayoutRow03, 1, 0, true)

	mainOptsLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	mainOptsLayout.SetBackgroundColor(bgColor)
	mainOptsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	mainOptsLayout.AddItem(layout, 0, 1, true)
	mainOptsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)

	dialog.layout.SetDirection(tview.FlexRow)
	dialog.layout.SetBackgroundColor(bgColor)
	dialog.layout.SetBorder(true)
	dialog.layout.SetTitle("PODMAN CONTAINER CHECKPOINT")
	dialog.layout.AddItem(mainOptsLayout, 0, 1, true)
	dialog.layout.AddItem(dialog.form, dialogs.DialogFormHeight, 0, true)

	return dialog
}

// Display displays this primitive.
func (d *ContainerCheckpointDialog) Display() {
	d.display = true

	d.createImage.SetText("")
	d.export.SetText("")
	d.fileLock.SetChecked(false)
	d.ignoreRootFS.SetChecked(false)
	d.keep.SetChecked(false)
	d.leaveRunning.SetChecked(false)
	d.preCheckpoint.SetChecked(false)
	d.printStats.SetChecked(false)
	d.tcpEstablished.SetChecked(false)
	d.withPrevious.SetChecked(false)

}

// IsDisplay returns true if this primitive is shown.
func (d *ContainerCheckpointDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *ContainerCheckpointDialog) Hide() {
	d.display = false
	d.focusElement = cntCheckpointImageFocus
}

// HasFocus returns whether or not this primitive has focus.
func (d *ContainerCheckpointDialog) HasFocus() bool {
	if d.createImage.HasFocus() || d.export.HasFocus() {
		return true
	}

	if d.withPrevious.HasFocus() || d.fileLock.HasFocus() {
		return true
	}

	if d.ignoreRootFS.HasFocus() || d.keep.HasFocus() {
		return true
	}

	if d.leaveRunning.HasFocus() || d.preCheckpoint.HasFocus() {
		return true
	}

	if d.printStats.HasFocus() || d.tcpEstablished.HasFocus() {
		return true
	}

	if d.form.HasFocus() || d.layout.HasFocus() {
		return true
	}

	return d.Box.HasFocus()
}

// Focus is called when this primitive receives focus.
func (d *ContainerCheckpointDialog) Focus(delegate func(p tview.Primitive)) {
	switch d.focusElement {
	case cntCheckpointImageFocus:
		delegate(d.createImage)
	case cntCheckpointExportFocus:
		delegate(d.export)
	case cntCheckpointFileLockFocus:
		delegate(d.fileLock)
	case cntCheckpointIgnoreRootFsFocus:
		delegate(d.ignoreRootFS)
	case cntCheckpointKeepFocus:
		delegate(d.keep)
	case cntCheckpointLeaveRunningFocus:
		delegate(d.leaveRunning)
	case cntCheckpointPreCheckpointFocus:
		delegate(d.preCheckpoint)
	case cntCheckpointPrintStatsFocus:
		delegate(d.printStats)
	case cntCheckpointTcpEstablishedFocus:
		delegate(d.tcpEstablished)
	case cntCheckpointWithPreviousFocus:
		delegate(d.withPrevious)
	case cntCheckpointFormFocus:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)

		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == utils.SwitchFocusKey.Key {
				d.focusElement = cntCheckpointImageFocus

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
func (d *ContainerCheckpointDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("network connect dialog: event %v received", event)

		if event.Key() == utils.CloseDialogKey.Key {
			d.cancelHandler()
		}

		if event.Key() == utils.SwitchFocusKey.Key {
			d.setFocusElement()
		}

		// events
		if d.createImage.HasFocus() {
			if createImageHandler := d.createImage.InputHandler(); createImageHandler != nil {
				createImageHandler(event, setFocus)

				return
			}
		}

		if d.export.HasFocus() {
			if exportHandler := d.export.InputHandler(); exportHandler != nil {
				exportHandler(event, setFocus)

				return
			}
		}

		if d.fileLock.HasFocus() {
			if fileLockHandler := d.fileLock.InputHandler(); fileLockHandler != nil {
				fileLockHandler(event, setFocus)

				return
			}
		}

		if d.ignoreRootFS.HasFocus() {
			if ignoreRootFSHandler := d.ignoreRootFS.InputHandler(); ignoreRootFSHandler != nil {
				ignoreRootFSHandler(event, setFocus)

				return
			}
		}

		if d.tcpEstablished.HasFocus() {
			if tcpEstablishedHandler := d.tcpEstablished.InputHandler(); tcpEstablishedHandler != nil {
				tcpEstablishedHandler(event, setFocus)

				return
			}
		}

		if d.keep.HasFocus() {
			if keephedHandler := d.keep.InputHandler(); keephedHandler != nil {
				keephedHandler(event, setFocus)

				return
			}
		}

		if d.printStats.HasFocus() {
			if printStatsHandler := d.printStats.InputHandler(); printStatsHandler != nil {
				printStatsHandler(event, setFocus)

				return
			}
		}

		if d.preCheckpoint.HasFocus() {
			if preCheckpointHandler := d.preCheckpoint.InputHandler(); preCheckpointHandler != nil {
				preCheckpointHandler(event, setFocus)

				return
			}
		}

		if d.leaveRunning.HasFocus() {
			if leaveRunningHandler := d.leaveRunning.InputHandler(); leaveRunningHandler != nil {
				leaveRunningHandler(event, setFocus)

				return
			}
		}

		if d.withPrevious.HasFocus() {
			if withPreviousHandler := d.withPrevious.InputHandler(); withPreviousHandler != nil {
				withPreviousHandler(event, setFocus)

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
func (d *ContainerCheckpointDialog) SetRect(x, y, width, height int) {

	if width > cntCheckpointDialogMaxWidth {
		emptySpace := (width - cntCheckpointDialogMaxWidth) / 2
		x = x + emptySpace
		width = cntCheckpointDialogMaxWidth
	}

	if height > cntCheckpointDialogMaxHeight {
		emptySpace := (height - cntCheckpointDialogMaxHeight) / 2
		y = y + emptySpace
		height = cntCheckpointDialogMaxHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive into the screen.
func (d *ContainerCheckpointDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.Box.DrawForSubclass(screen, d)

	x, y, width, height := d.Box.GetInnerRect()

	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetCheckpointFunc sets form checkpoint button selected function.
func (d *ContainerCheckpointDialog) SetCheckpointFunc(handler func()) *ContainerCheckpointDialog {
	d.checkpointHandler = handler
	connectButton := d.form.GetButton(d.form.GetButtonCount() - 1)

	connectButton.SetSelectedFunc(handler)

	return d
}

// SetCancelFunc sets form cancel button selected function.
func (d *ContainerCheckpointDialog) SetCancelFunc(handler func()) *ContainerCheckpointDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2)

	cancelButton.SetSelectedFunc(handler)

	return d
}

func (d *ContainerCheckpointDialog) setFocusElement() {
	switch d.focusElement {
	case cntCheckpointImageFocus:
		d.focusElement = cntCheckpointExportFocus
	case cntCheckpointExportFocus:
		d.focusElement = cntCheckpointFileLockFocus
	case cntCheckpointFileLockFocus:
		d.focusElement = cntCheckpointIgnoreRootFsFocus
	case cntCheckpointIgnoreRootFsFocus:
		d.focusElement = cntCheckpointTcpEstablishedFocus
	case cntCheckpointTcpEstablishedFocus:
		d.focusElement = cntCheckpointPreCheckpointFocus
	case cntCheckpointPreCheckpointFocus:
		d.focusElement = cntCheckpointPrintStatsFocus
	case cntCheckpointPrintStatsFocus:
		d.focusElement = cntCheckpointKeepFocus
	case cntCheckpointKeepFocus:
		d.focusElement = cntCheckpointLeaveRunningFocus
	case cntCheckpointLeaveRunningFocus:
		d.focusElement = cntCheckpointWithPreviousFocus
	case cntCheckpointWithPreviousFocus:
		d.focusElement = cntCheckpointFormFocus
	}
}

// SetContainerInfo sets selected container ID and name information.
func (d *ContainerCheckpointDialog) SetContainerInfo(id string, name string) {
	d.containerID = id
	containerInfo := fmt.Sprintf("Container: %s (%s)", id, name)

	d.containerInfo.SetText(containerInfo)
}

// GetOptions returns checkpoint options
func (d *ContainerCheckpointDialog) GetCheckpointOptions() containers.CntCheckPointOptions {
	var opts containers.CntCheckPointOptions

	opts.CreateImage = strings.Trim(d.createImage.GetText(), " ")
	opts.Export = strings.Trim(d.export.GetText(), " ")
	opts.FileLocks = d.fileLock.IsChecked()
	opts.IgnoreRootFs = d.ignoreRootFS.IsChecked()
	opts.TCPEstablished = d.tcpEstablished.IsChecked()
	opts.Keep = d.keep.IsChecked()
	opts.PrintStats = d.printStats.IsChecked()
	opts.PreCheckpoint = d.preCheckpoint.IsChecked()
	opts.LeaveRunning = d.leaveRunning.IsChecked()
	opts.WithPrevious = d.withPrevious.IsChecked()

	return opts
}
