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
	volumeImportDialogMaxWidth = 80
	volumeImportDialogHeight   = 9
)

const (
	volumeImportFormFieldFocus = 0 + iota
	volumeImportSourceFieldFocus
)

// VolumeImportDialog implements volume input dialog.
type VolumeImportDialog struct {
	*tview.Box

	layout        *tview.Flex
	form          *tview.Form
	source        *tview.InputField
	volume        *tview.InputField
	display       bool
	focusElement  int
	cancelHandler func()
	importHandler func()
}

func NewVolumeImportDialog() *VolumeImportDialog {
	importDialog := VolumeImportDialog{
		Box:     tview.NewBox(),
		layout:  tview.NewFlex().SetDirection(tview.FlexRow),
		form:    tview.NewForm(),
		source:  tview.NewInputField(),
		volume:  tview.NewInputField(),
		display: false,
	}

	bgColor := style.DialogBgColor
	buttonBgColor := style.ButtonBgColor

	// volume
	importDialog.volume.SetBackgroundColor(style.DialogBgColor)
	importDialog.volume.SetFieldBackgroundColor(style.DialogBgColor)
	importDialog.volume.SetLabelStyle(tcell.StyleDefault.
		Background(style.DialogBorderColor).
		Foreground(style.DialogFgColor))

	// source
	label := "source:"

	importDialog.source.SetBackgroundColor(bgColor)
	importDialog.source.SetLabel(utils.StringToInputLabel(label, len(label)+1))
	importDialog.source.SetFieldStyle(style.InputFieldStyle)
	importDialog.source.SetLabelStyle(style.InputLabelStyle)

	// form
	importDialog.form.SetBackgroundColor(bgColor)
	importDialog.form.AddButton("Cancel", nil)
	importDialog.form.AddButton("Import", nil)
	importDialog.form.SetButtonsAlign(tview.AlignRight)
	importDialog.form.SetButtonBackgroundColor(buttonBgColor)

	importDialog.setupLayout()
	importDialog.layout.SetBackgroundColor(bgColor)
	importDialog.layout.SetBorder(true)
	importDialog.layout.SetBorderColor(style.DialogBorderColor)
	importDialog.layout.SetTitle("PODMAN VOLUME IMPORT")

	return &importDialog
}

// Display displays this primitive.
func (d *VolumeImportDialog) Display() {
	d.display = true
	d.initData()
}

// IsDisplay returns true if primitive is shown.
func (d *VolumeImportDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *VolumeImportDialog) Hide() {
	d.display = false
}

// HasFocus returns whether or not this primitive has focus.
func (d *VolumeImportDialog) HasFocus() bool {
	for _, item := range d.getInnerPrimitives() {
		if item.HasFocus() {
			return true
		}
	}

	return d.Box.HasFocus() || d.form.HasFocus()
}

// Focus is called when this primitive receives focus.
func (d *VolumeImportDialog) Focus(delegate func(p tview.Primitive)) {
	switch d.focusElement {
	case volumeImportSourceFieldFocus:
		delegate(d.source)
	// form has focus
	case volumeImportFormFieldFocus:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = volumeImportSourceFieldFocus // category text view
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
func (d *VolumeImportDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return volumeDialogInputHandler(
		"import",
		d.Box,
		d.form,
		d.getInnerPrimitives(),
		d.nextFocus,
		d.cancelHandler,
		d.importHandler,
	)
}

// SetRect set rects for this primitive.
func (d *VolumeImportDialog) SetRect(x, y, width, height int) {
	if width > volumeImportDialogMaxWidth {
		emptySpace := (width - volumeImportDialogMaxWidth) / 2 //nolint:mnd
		x += emptySpace
		width = volumeImportDialogMaxWidth
	}

	if height > volumeImportDialogHeight {
		emptySpace := (height - volumeImportDialogHeight) / 2 //nolint:mnd
		y += emptySpace
		height = volumeImportDialogHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen.
func (d *VolumeImportDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.DrawForSubclass(screen, d)
	x, y, width, height := d.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetCancelFunc sets form cancel button selected function.
func (d *VolumeImportDialog) SetCancelFunc(handler func()) *VolumeImportDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:mnd
	cancelButton.SetSelectedFunc(handler)

	return d
}

// SetImportFunc sets form import button selected function.
func (d *VolumeImportDialog) SetImportFunc(handler func()) *VolumeImportDialog {
	d.importHandler = handler
	importButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	importButton.SetSelectedFunc(handler)

	return d
}

func (d *VolumeImportDialog) VolumeImportSource() string {
	return strings.TrimSpace(d.source.GetText())
}

func (d *VolumeImportDialog) SetVolumeInfo(name string) {
	d.volume.SetLabel("[::b]VOLUME:")
	d.volume.SetText(" " + name)
}

func (d *VolumeImportDialog) initData() {
	d.source.SetText("")
	d.volume.SetText("")
	d.focusElement = volumeImportSourceFieldFocus
}

func (d *VolumeImportDialog) nextFocus() {
	if d.focusElement == volumeImportSourceFieldFocus {
		d.focusElement = volumeImportFormFieldFocus
	}
}

func (d *VolumeImportDialog) getInnerPrimitives() []tview.Primitive {
	return []tview.Primitive{
		d.source,
	}
}

func (d *VolumeImportDialog) setupLayout() {
	bgColor := style.DialogBgColor

	// layouts
	inputFieldLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	inputFieldLayout.SetBackgroundColor(bgColor)
	inputFieldLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	inputFieldLayout.AddItem(d.volume, 1, 0, true)
	inputFieldLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	inputFieldLayout.AddItem(d.source, 1, 0, true)
	inputFieldLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)

	// adding an empty column space to beginning and end of the fields layout
	layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	layout.AddItem(inputFieldLayout, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)

	d.layout.AddItem(layout, 0, 1, true)
	d.layout.AddItem(d.form, dialogs.DialogFormHeight, 0, true)
}
