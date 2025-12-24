package voldialogs

import (
	"strings"

	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	volumeExportDialogMaxWidth = 80
	volumeExportDialogHeight   = 8
)

const (
	volumeExportFormFieldFocus = 0 + iota
	volumeExportOutputFieldFocus
)

// VolumeExportDialog implements volume export dialog.
type VolumeExportDialog struct {
	*tview.Box

	layout        *tview.Flex
	form          *tview.Form
	output        *tview.InputField
	display       bool
	focusElement  int
	cancelHandler func()
	exportHandler func()
}

func NewVolumeExportDialog() *VolumeExportDialog {
	exportDialog := VolumeExportDialog{
		Box:     tview.NewBox(),
		layout:  tview.NewFlex().SetDirection(tview.FlexRow),
		form:    tview.NewForm(),
		output:  tview.NewInputField(),
		display: false,
	}

	bgColor := style.DialogBgColor
	buttonBgColor := style.ButtonBgColor

	// output
	label := "output:"

	exportDialog.output.SetBackgroundColor(bgColor)
	exportDialog.output.SetLabel(utils.StringToInputLabel(label, len(label)+1))
	exportDialog.output.SetFieldStyle(style.InputFieldStyle)
	exportDialog.output.SetLabelStyle(style.InputLabelStyle)

	// form
	exportDialog.form.SetBackgroundColor(bgColor)
	exportDialog.form.AddButton("Cancel", nil)
	exportDialog.form.AddButton("Export", nil)
	exportDialog.form.SetButtonsAlign(tview.AlignRight)
	exportDialog.form.SetButtonBackgroundColor(buttonBgColor)

	exportDialog.setupLayout()
	exportDialog.layout.SetBackgroundColor(bgColor)
	exportDialog.layout.SetBorder(true)
	exportDialog.layout.SetBorderColor(style.DialogBorderColor)
	exportDialog.layout.SetTitle("PODMAN VOLUME EXPORT")

	return &exportDialog
}

// Display displays this primitive.
func (d *VolumeExportDialog) Display() {
	d.display = true
	d.initData()
}

// IsDisplay returns true if primitive is shown.
func (d *VolumeExportDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *VolumeExportDialog) Hide() {
	d.display = false
}

// HasFocus returns whether or not this primitive has focus.
func (d *VolumeExportDialog) HasFocus() bool {
	for _, item := range d.getInnerPrimitives() {
		if item.HasFocus() {
			return true
		}
	}

	return d.Box.HasFocus() || d.form.HasFocus()
}

// Focus is called when this primitive receives focus.
func (d *VolumeExportDialog) Focus(delegate func(p tview.Primitive)) {
	switch d.focusElement {
	case volumeExportOutputFieldFocus:
		delegate(d.output)
	// form has focus
	case volumeExportFormFieldFocus:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = volumeExportOutputFieldFocus // category text view
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
	}
}

// InputHandler returns input handler function for this primitive.
func (d *VolumeExportDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return volumeDialogInputHandler(
		"export",
		d.Box,
		d.form,
		d.getInnerPrimitives(),
		d.nextFocus,
		d.cancelHandler,
		d.exportHandler,
	)
}

// SetRect set rects for this primitive.
func (d *VolumeExportDialog) SetRect(x, y, width, height int) {
	if width > volumeExportDialogMaxWidth {
		emptySpace := (width - volumeExportDialogMaxWidth) / 2 //nolint:mnd
		x += emptySpace
		width = volumeExportDialogMaxWidth
	}

	if height > volumeExportDialogHeight {
		emptySpace := (height - volumeExportDialogHeight) / 2 //nolint:mnd
		y += emptySpace
		height = volumeExportDialogHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen.
func (d *VolumeExportDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.DrawForSubclass(screen, d)
	x, y, width, height := d.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetCancelFunc sets form cancel button selected function.
func (d *VolumeExportDialog) SetCancelFunc(handler func()) *VolumeExportDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:mnd
	cancelButton.SetSelectedFunc(handler)

	return d
}

// SetExportFunc sets form export button selected function.
func (d *VolumeExportDialog) SetExportFunc(handler func()) *VolumeExportDialog {
	d.exportHandler = handler
	enterButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	enterButton.SetSelectedFunc(handler)

	return d
}

func (d *VolumeExportDialog) VolumeExportOutput() string {
	return strings.TrimSpace(d.output.GetText())
}

func (d *VolumeExportDialog) initData() {
	d.output.SetText("")
	d.focusElement = volumeExportOutputFieldFocus
}

func (d *VolumeExportDialog) nextFocus() {
	if d.focusElement == volumeExportOutputFieldFocus {
		d.focusElement = volumeExportFormFieldFocus
	}
}

func (d *VolumeExportDialog) getInnerPrimitives() []tview.Primitive {
	return []tview.Primitive{
		d.output,
	}
}

func (d *VolumeExportDialog) setupLayout() {
	bgColor := style.DialogBgColor

	// layouts
	inputFieldLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	inputFieldLayout.SetBackgroundColor(bgColor)
	inputFieldLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	inputFieldLayout.AddItem(d.output, 1, 0, true)
	inputFieldLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)

	// adding an empty column space to beginning and end of the fields layout
	layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	layout.AddItem(inputFieldLayout, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)

	d.layout.AddItem(layout, 0, 1, true)
	d.layout.AddItem(d.form, dialogs.DialogFormHeight, 0, true)
}
