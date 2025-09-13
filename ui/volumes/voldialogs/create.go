package voldialogs

import (
	"strings"

	"github.com/containers/podman-tui/pdcs/volumes"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	volumeCreateDialogMaxWidth = 60
	volumeCreateDialogHeight   = 13
)

const (
	formFocus = 0 + iota
	volumeNameFieldFocus
	volumeLabelsFieldFocus
	volumeDriverNameFocus
	volumeDriverOptionsFocus
)

// VolumeCreateDialog implements volume create dialog.
type VolumeCreateDialog struct {
	*tview.Box

	layout                   *tview.Flex
	form                     *tview.Form
	display                  bool
	focusElement             int
	volumeNameField          *tview.InputField
	volumeLabelField         *tview.InputField
	volumeDriverField        *tview.InputField
	volumeDriverOptionsField *tview.InputField
	cancelHandler            func()
	createHandler            func()
}

// NewVolumeCreateDialog returns new pod create dialog primitive VolumeCreateDialog.
func NewVolumeCreateDialog() *VolumeCreateDialog {
	volDialog := VolumeCreateDialog{
		Box:                      tview.NewBox(),
		layout:                   tview.NewFlex().SetDirection(tview.FlexRow),
		form:                     tview.NewForm(),
		display:                  false,
		volumeNameField:          tview.NewInputField(),
		volumeLabelField:         tview.NewInputField(),
		volumeDriverField:        tview.NewInputField(),
		volumeDriverOptionsField: tview.NewInputField(),
	}

	bgColor := style.DialogBgColor
	buttonBgColor := style.ButtonBgColor

	// basic information setup page
	basicInfoPageLabelWidth := 9
	// name field
	volDialog.volumeNameField.SetBackgroundColor(bgColor)
	volDialog.volumeNameField.SetLabel(utils.StringToInputLabel("name:", basicInfoPageLabelWidth))
	volDialog.volumeNameField.SetFieldStyle(style.InputFieldStyle)
	volDialog.volumeNameField.SetLabelStyle(style.InputLabelStyle)

	// labels field
	volDialog.volumeLabelField.SetBackgroundColor(bgColor)
	volDialog.volumeLabelField.SetLabel(utils.StringToInputLabel("labels:", basicInfoPageLabelWidth))
	volDialog.volumeLabelField.SetFieldStyle(style.InputFieldStyle)
	volDialog.volumeLabelField.SetLabelStyle(style.InputLabelStyle)

	// drivers
	volDialog.volumeDriverField.SetBackgroundColor(bgColor)
	volDialog.volumeDriverField.SetLabel(utils.StringToInputLabel("drivers:", basicInfoPageLabelWidth))
	volDialog.volumeDriverField.SetFieldStyle(style.InputFieldStyle)
	volDialog.volumeDriverField.SetLabelStyle(style.InputLabelStyle)

	// drivers options
	volDialog.volumeDriverOptionsField.SetBackgroundColor(bgColor)
	volDialog.volumeDriverOptionsField.SetLabel(utils.StringToInputLabel("options:", basicInfoPageLabelWidth))
	volDialog.volumeDriverOptionsField.SetFieldStyle(style.InputFieldStyle)
	volDialog.volumeDriverOptionsField.SetLabelStyle(style.InputLabelStyle)

	// form
	volDialog.form.SetBackgroundColor(bgColor)
	volDialog.form.AddButton("Cancel", nil)
	volDialog.form.AddButton("Create", nil)
	volDialog.form.SetButtonsAlign(tview.AlignRight)
	volDialog.form.SetButtonBackgroundColor(buttonBgColor)

	volDialog.setupLayout()
	volDialog.layout.SetBackgroundColor(bgColor)
	volDialog.layout.SetBorder(true)
	volDialog.layout.SetBorderColor(style.DialogBorderColor)
	volDialog.layout.SetTitle("PODMAN VOLUME CREATE")

	return &volDialog
}

// Display displays this primitive.
func (d *VolumeCreateDialog) Display() {
	d.display = true
	d.focusElement = volumeNameFieldFocus
	d.initData()
}

// IsDisplay returns true if primitive is shown.
func (d *VolumeCreateDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *VolumeCreateDialog) Hide() {
	d.display = false
}

// HasFocus returns whether or not this primitive has focus.
func (d *VolumeCreateDialog) HasFocus() bool {
	for _, item := range d.getInnerPrimitives() {
		if item.HasFocus() {
			return true
		}
	}

	return d.Box.HasFocus() || d.form.HasFocus()
}

// Focus is called when this primitive receives focus.
func (d *VolumeCreateDialog) Focus(delegate func(p tview.Primitive)) {
	switch d.focusElement {
	// form has focus
	case formFocus:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = volumeNameFieldFocus // category text view
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
	// basic info page
	case volumeNameFieldFocus:
		delegate(d.volumeNameField)
	case volumeLabelsFieldFocus:
		delegate(d.volumeLabelField)
	case volumeDriverNameFocus:
		delegate(d.volumeDriverField)
	case volumeDriverOptionsFocus:
		delegate(d.volumeDriverOptionsField)
	}
}

// InputHandler returns input handler function for this primitive.
func (d *VolumeCreateDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("volume create dialog: event %v received", event)

		if event.Key() == tcell.KeyEsc {
			d.cancelHandler()

			return
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

		if event.Key() == tcell.KeyTab {
			d.nextFocus()
			d.Focus(setFocus)

			return
		}

		for _, item := range d.getInnerPrimitives() {
			if item.HasFocus() {
				if handler := item.InputHandler(); handler != nil {
					handler(event, setFocus)

					return
				}
			}
		}
	})
}

// SetRect set rects for this primitive.
func (d *VolumeCreateDialog) SetRect(x, y, width, height int) {
	if width > volumeCreateDialogMaxWidth {
		emptySpace := (width - volumeCreateDialogMaxWidth) / 2 //nolint:mnd
		x += emptySpace
		width = volumeCreateDialogMaxWidth
	}

	if height > volumeCreateDialogHeight {
		emptySpace := (height - volumeCreateDialogHeight) / 2 //nolint:mnd
		y += emptySpace
		height = volumeCreateDialogHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen.
func (d *VolumeCreateDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.DrawForSubclass(screen, d)
	x, y, width, height := d.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetCancelFunc sets form cancel button selected function.
func (d *VolumeCreateDialog) SetCancelFunc(handler func()) *VolumeCreateDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:mnd
	cancelButton.SetSelectedFunc(handler)

	return d
}

// SetCreateFunc sets form create button selected function.
func (d *VolumeCreateDialog) SetCreateFunc(handler func()) *VolumeCreateDialog {
	d.createHandler = handler
	enterButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	enterButton.SetSelectedFunc(handler)

	return d
}

// VolumeCreateOptions returns new volume options.
func (d *VolumeCreateDialog) VolumeCreateOptions() volumes.CreateOptions { //nolint:cyclop
	var (
		labels  = make(map[string]string)
		options = make(map[string]string)
	)

	for _, label := range strings.Split(d.volumeLabelField.GetText(), " ") {
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

	for _, option := range strings.Split(d.volumeDriverOptionsField.GetText(), " ") {
		if option != "" {
			split := strings.Split(option, "=")
			if len(split) == 2 { //nolint:mnd
				key := split[0]
				value := split[1]

				if key != "" && value != "" {
					options[key] = value
				}
			}
		}
	}

	opts := volumes.CreateOptions{
		Name:          d.volumeNameField.GetText(),
		Labels:        labels,
		Driver:        d.volumeDriverField.GetText(),
		DriverOptions: options,
	}

	return opts
}

func (d *VolumeCreateDialog) initData() {
	d.volumeNameField.SetText("")
	d.volumeLabelField.SetText("")
	d.volumeDriverField.SetText("")
	d.volumeDriverOptionsField.SetText("")
}

func (d *VolumeCreateDialog) getInnerPrimitives() []tview.Primitive {
	return []tview.Primitive{
		d.volumeNameField,
		d.volumeLabelField,
		d.volumeDriverField,
		d.volumeDriverOptionsField,
	}
}

func (d *VolumeCreateDialog) nextFocus() {
	switch d.focusElement {
	case volumeNameFieldFocus:
		d.focusElement = volumeLabelsFieldFocus
	case volumeLabelsFieldFocus:
		d.focusElement = volumeDriverNameFocus
	case volumeDriverNameFocus:
		d.focusElement = volumeDriverOptionsFocus
	case volumeDriverOptionsFocus:
		d.focusElement = formFocus
	}
}

func (d *VolumeCreateDialog) setupLayout() {
	bgColor := style.DialogBgColor

	// layouts
	inputFieldLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	inputFieldLayout.SetBackgroundColor(bgColor)
	inputFieldLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	inputFieldLayout.AddItem(d.volumeNameField, 1, 0, true)
	inputFieldLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	inputFieldLayout.AddItem(d.volumeLabelField, 1, 0, true)
	inputFieldLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	inputFieldLayout.AddItem(d.volumeDriverField, 1, 0, true)
	inputFieldLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	inputFieldLayout.AddItem(d.volumeDriverOptionsField, 1, 0, true)

	// adding an empty column space to beginning and end of the fields layout
	layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	layout.AddItem(inputFieldLayout, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)

	d.layout.AddItem(layout, 0, 1, true)
	d.layout.AddItem(d.form, dialogs.DialogFormHeight, 0, true)
}
