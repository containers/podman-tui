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
	cntCheckpointDialogMaxWidth     = 73
	cntCheckpointDialogMaxHeight    = 15
	cntCheckpointDialogLabelPadding = 1
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
	cntCheckpointTCPEstablishedFocus
	cntCheckpointWithPreviousFocus
	cntCheckpointFormFocus
)

// ContainerCheckpointDialog implements container checkpoint dialog primitive.
type ContainerCheckpointDialog struct {
	*tview.Box

	layout            *tview.Flex
	containerInfo     *tview.InputField
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
		containerInfo:  tview.NewInputField(),
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

	labelWidth := 14
	chkGroupFirstColLabelWidth := 14
	chkGroupSecondColLabelWidth := 16
	chkGroupThirdColLabelWidth := 16

	// containerInfo
	dialog.containerInfo.SetBackgroundColor(style.DialogBgColor)
	dialog.containerInfo.SetLabel("[::b]" + utils.ContainerIDLabel)
	dialog.containerInfo.SetFieldBackgroundColor(style.DialogBgColor)
	dialog.containerInfo.SetLabelStyle(tcell.StyleDefault.
		Background(style.DialogBorderColor).
		Foreground(style.DialogFgColor))

	// createImage
	dialog.createImage.SetBackgroundColor(style.DialogBgColor)
	dialog.createImage.SetLabel(utils.StringToInputLabel("create image:", labelWidth))
	dialog.createImage.SetFieldStyle(style.InputFieldStyle)
	dialog.createImage.SetLabelStyle(style.InputLabelStyle)

	// export
	dialog.export.SetBackgroundColor(style.DialogBgColor)
	dialog.export.SetLabel(utils.StringToInputLabel("export:", labelWidth))
	dialog.export.SetFieldStyle(style.InputFieldStyle)
	dialog.export.SetLabelStyle(style.InputLabelStyle)

	// printStats
	dialog.printStats.SetLabel("print stats:")
	dialog.printStats.SetLabelWidth(labelWidth)
	dialog.printStats.SetChecked(false)
	dialog.printStats.SetBackgroundColor(style.DialogBgColor)
	dialog.printStats.SetLabelColor(style.DialogFgColor)
	dialog.printStats.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// fileLock
	dialog.fileLock.SetLabel("file lock:")
	dialog.fileLock.SetLabelWidth(labelWidth)
	dialog.fileLock.SetChecked(false)
	dialog.fileLock.SetBackgroundColor(style.DialogBgColor)
	dialog.fileLock.SetLabelColor(style.DialogFgColor)
	dialog.fileLock.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// ignoreRootFS
	ignoreRootFSLabel := fmt.Sprintf("%*s ", chkGroupFirstColLabelWidth, "ignore rootFS:")

	dialog.ignoreRootFS.SetLabel(ignoreRootFSLabel)
	dialog.ignoreRootFS.SetChecked(false)
	dialog.ignoreRootFS.SetBackgroundColor(style.DialogBgColor)
	dialog.ignoreRootFS.SetLabelColor(style.DialogFgColor)
	dialog.ignoreRootFS.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// keep
	keepLabel := fmt.Sprintf("%*s ", chkGroupFirstColLabelWidth, "keep:")

	dialog.keep.SetLabel(keepLabel)
	dialog.keep.SetChecked(false)
	dialog.keep.SetBackgroundColor(style.DialogBgColor)
	dialog.keep.SetLabelColor(style.DialogFgColor)
	dialog.keep.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// tcpEstablished
	tcpEstablishedLabel := fmt.Sprintf("%*s ", chkGroupSecondColLabelWidth, "tcp established:")

	dialog.tcpEstablished.SetLabel(tcpEstablishedLabel)
	dialog.tcpEstablished.SetChecked(false)
	dialog.tcpEstablished.SetBackgroundColor(style.DialogBgColor)
	dialog.tcpEstablished.SetLabelColor(style.DialogFgColor)
	dialog.tcpEstablished.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// leaveRunning
	leaveRunningLabel := fmt.Sprintf("%*s ", chkGroupSecondColLabelWidth, "leave running:")

	dialog.leaveRunning.SetLabel(leaveRunningLabel)
	dialog.leaveRunning.SetChecked(false)
	dialog.leaveRunning.SetBackgroundColor(style.DialogBgColor)
	dialog.leaveRunning.SetLabelColor(style.DialogFgColor)
	dialog.leaveRunning.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// preCheckpoint
	preCheckPointLabel := fmt.Sprintf("%*s ", chkGroupThirdColLabelWidth, "pre checkpoint:")

	dialog.preCheckpoint.SetLabel(preCheckPointLabel)
	dialog.preCheckpoint.SetChecked(false)
	dialog.preCheckpoint.SetBackgroundColor(style.DialogBgColor)
	dialog.preCheckpoint.SetLabelColor(style.DialogFgColor)
	dialog.preCheckpoint.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// withPrevious
	withPreviousLabel := fmt.Sprintf("%*s ", chkGroupThirdColLabelWidth, "with previous:")

	dialog.withPrevious.SetLabel(withPreviousLabel)
	dialog.withPrevious.SetChecked(false)
	dialog.withPrevious.SetBackgroundColor(style.DialogBgColor)
	dialog.withPrevious.SetLabelColor(style.DialogFgColor)
	dialog.withPrevious.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// form
	dialog.form.AddButton(" Cancel ", nil)
	dialog.form.AddButton("Checkpoint", nil)
	dialog.form.SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(style.DialogBgColor)
	dialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// layout
	optionsLayoutRow01 := tview.NewFlex().SetDirection(tview.FlexRow)
	optionsLayoutRow01.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)
	optionsLayoutRow01.AddItem(dialog.containerInfo, 1, 0, true)
	optionsLayoutRow01.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)
	optionsLayoutRow01.AddItem(dialog.createImage, 1, 0, true)
	optionsLayoutRow01.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)
	optionsLayoutRow01.AddItem(dialog.export, 1, 0, true)

	optionsLayoutRow02 := tview.NewFlex().SetDirection(tview.FlexColumn)
	optionsLayoutRow02.AddItem(dialog.fileLock, labelWidth+2, 0, true) //nolint:mnd
	optionsLayoutRow02.AddItem(dialog.ignoreRootFS, 0, 1, true)
	optionsLayoutRow02.AddItem(dialog.tcpEstablished, 0, 1, true)
	optionsLayoutRow02.AddItem(dialog.preCheckpoint, 0, 1, true)

	optionsLayoutRow03 := tview.NewFlex().SetDirection(tview.FlexColumn)
	optionsLayoutRow03.AddItem(dialog.printStats, labelWidth+2, 1, true) //nolint:mnd
	optionsLayoutRow03.AddItem(dialog.keep, 0, 1, true)
	optionsLayoutRow03.AddItem(dialog.leaveRunning, 0, 1, true)
	optionsLayoutRow03.AddItem(dialog.withPrevious, 0, 1, true)

	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	layout.AddItem(optionsLayoutRow01, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)
	layout.AddItem(optionsLayoutRow02, 1, 0, true)
	layout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)
	layout.AddItem(optionsLayoutRow03, 1, 0, true)

	mainOptsLayout := tview.NewFlex().SetDirection(tview.FlexColumn)

	mainOptsLayout.SetBackgroundColor(style.DialogBgColor)
	mainOptsLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)
	mainOptsLayout.AddItem(layout, 0, 1, true)
	mainOptsLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)

	dialog.layout.SetDirection(tview.FlexRow)
	dialog.layout.SetBackgroundColor(style.DialogBgColor)
	dialog.layout.SetBorder(true)
	dialog.layout.SetBorderColor(style.DialogBorderColor)
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
func (d *ContainerCheckpointDialog) HasFocus() bool { //nolint:cyclop
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
func (d *ContainerCheckpointDialog) Focus(delegate func(p tview.Primitive)) { //nolint:cyclop
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
	case cntCheckpointTCPEstablishedFocus:
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
func (d *ContainerCheckpointDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) { //nolint:gocognit,lll,cyclop
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("container checkpoint dialog: event %v received", event)

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
		emptySpace := (width - cntCheckpointDialogMaxWidth) / 2 //nolint:mnd
		x += emptySpace
		width = cntCheckpointDialogMaxWidth
	}

	if height > cntCheckpointDialogMaxHeight {
		emptySpace := (height - cntCheckpointDialogMaxHeight) / 2 //nolint:mnd
		y += emptySpace
		height = cntCheckpointDialogMaxHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive into the screen.
func (d *ContainerCheckpointDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.DrawForSubclass(screen, d)

	x, y, width, height := d.GetInnerRect()

	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetCheckpointFunc sets form checkpoint button selected function.
func (d *ContainerCheckpointDialog) SetCheckpointFunc(handler func()) *ContainerCheckpointDialog {
	d.checkpointHandler = handler
	checkpointButton := d.form.GetButton(d.form.GetButtonCount() - 1)

	checkpointButton.SetSelectedFunc(handler)

	return d
}

// SetCancelFunc sets form cancel button selected function.
func (d *ContainerCheckpointDialog) SetCancelFunc(handler func()) *ContainerCheckpointDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:mnd

	cancelButton.SetSelectedFunc(handler)

	return d
}

// SetContainerInfo sets selected container ID and name information.
func (d *ContainerCheckpointDialog) SetContainerInfo(id string, name string) {
	d.containerID = id
	containerInfo := fmt.Sprintf("%12s (%s)", id, name)
	containerInfo = utils.LabelWidthLeftPadding(containerInfo, cntCheckpointDialogLabelPadding)

	d.containerInfo.SetText(containerInfo)
}

// GetOptions returns checkpoint options.
func (d *ContainerCheckpointDialog) GetCheckpointOptions() containers.CntCheckPointOptions {
	var opts containers.CntCheckPointOptions

	opts.CreateImage = strings.Trim(d.createImage.GetText(), " ")

	exportPath := strings.Trim(d.export.GetText(), " ")
	if strings.Index(exportPath, "~") == 0 {
		exportPath, _ = utils.ResolveHomeDir(exportPath)
	}

	opts.Export = exportPath

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

func (d *ContainerCheckpointDialog) setFocusElement() { //nolint:cyclop
	switch d.focusElement {
	case cntCheckpointImageFocus:
		d.focusElement = cntCheckpointExportFocus
	case cntCheckpointExportFocus:
		d.focusElement = cntCheckpointFileLockFocus
	case cntCheckpointFileLockFocus:
		d.focusElement = cntCheckpointIgnoreRootFsFocus
	case cntCheckpointIgnoreRootFsFocus:
		d.focusElement = cntCheckpointTCPEstablishedFocus
	case cntCheckpointTCPEstablishedFocus:
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
