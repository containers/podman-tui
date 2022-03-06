package voldialogs

import (
	"strings"

	"github.com/containers/podman-tui/pdcs/volumes"
	"github.com/containers/podman-tui/ui/dialogs"
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

// VolumeCreateDialog implements volume create dialog
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

// NewVolumeCreateDialog returns new pod create dialog primitive VolumeCreateDialog
func NewVolumeCreateDialog() *VolumeCreateDialog {
	podDialog := VolumeCreateDialog{
		Box:                      tview.NewBox(),
		layout:                   tview.NewFlex().SetDirection(tview.FlexRow),
		form:                     tview.NewForm(),
		display:                  false,
		volumeNameField:          tview.NewInputField(),
		volumeLabelField:         tview.NewInputField(),
		volumeDriverField:        tview.NewInputField(),
		volumeDriverOptionsField: tview.NewInputField(),
	}

	bgColor := utils.Styles.ImageHistoryDialog.BgColor

	// basic information setup page
	basicInfoPageLabelWidth := 9
	// name field
	podDialog.volumeNameField.SetLabel("name:")
	podDialog.volumeNameField.SetLabelWidth(basicInfoPageLabelWidth)
	podDialog.volumeNameField.SetBackgroundColor(bgColor)
	podDialog.volumeNameField.SetLabelColor(tcell.ColorWhite)
	// labels field
	podDialog.volumeLabelField.SetLabel("labels:")
	podDialog.volumeLabelField.SetLabelWidth(basicInfoPageLabelWidth)
	podDialog.volumeLabelField.SetBackgroundColor(bgColor)
	podDialog.volumeLabelField.SetLabelColor(tcell.ColorWhite)
	// drivers
	podDialog.volumeDriverField.SetLabel("drivers:")
	podDialog.volumeDriverField.SetLabelWidth(basicInfoPageLabelWidth)
	podDialog.volumeDriverField.SetBackgroundColor(bgColor)
	podDialog.volumeDriverField.SetLabelColor(tcell.ColorWhite)
	// drivers options
	podDialog.volumeDriverOptionsField.SetLabel("options:")
	podDialog.volumeDriverOptionsField.SetLabelWidth(basicInfoPageLabelWidth)
	podDialog.volumeDriverOptionsField.SetBackgroundColor(bgColor)
	podDialog.volumeDriverOptionsField.SetLabelColor(tcell.ColorWhite)

	// form
	podDialog.form.SetBackgroundColor(bgColor)
	podDialog.form.AddButton("Cancel", nil)
	podDialog.form.AddButton("Create", nil)
	podDialog.form.SetButtonsAlign(tview.AlignRight)

	podDialog.setupLayout()
	podDialog.layout.SetBackgroundColor(bgColor)
	podDialog.layout.SetBorder(true)
	podDialog.layout.SetTitle("PODMAN VOLUME CREATE")

	return &podDialog
}

func (d *VolumeCreateDialog) setupLayout() {
	bgColor := utils.Styles.ImageHistoryDialog.BgColor

	// basic info page
	d.layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.layout.AddItem(d.volumeNameField, 1, 0, true)
	d.layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.layout.AddItem(d.volumeLabelField, 1, 0, true)
	d.layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.layout.AddItem(d.volumeDriverField, 1, 0, true)
	d.layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	d.layout.AddItem(d.volumeDriverOptionsField, 1, 0, true)
	d.layout.AddItem(d.form, dialogs.DialogFormHeight, 0, true)

}

// Display displays this primitive
func (d *VolumeCreateDialog) Display() {
	d.display = true
	d.focusElement = volumeNameFieldFocus
	d.initData()
}

// IsDisplay returns true if primitive is shown
func (d *VolumeCreateDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive
func (d *VolumeCreateDialog) Hide() {
	d.display = false
}

// HasFocus returns whether or not this primitive has focus
func (d *VolumeCreateDialog) HasFocus() bool {
	if d.volumeNameField.HasFocus() || d.volumeLabelField.HasFocus() {
		return true
	}
	if d.volumeDriverField.HasFocus() || d.volumeDriverOptionsField.HasFocus() {
		return true
	}
	return d.Box.HasFocus() || d.form.HasFocus()
}

// Focus is called when this primitive receives focus
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
				//d.pullSelectHandler()
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

// InputHandler returns input handler function for this primitive
func (d *VolumeCreateDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("volume create dialog: event %v received", event)
		if event.Key() == tcell.KeyEsc {
			d.cancelHandler()
			return
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
		if event.Key() == tcell.KeyTab {
			d.nextFocus()
			d.Focus(setFocus)
			return
		}
		if d.volumeNameField.HasFocus() {
			if handler := d.volumeNameField.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
		if d.volumeLabelField.HasFocus() {
			if handler := d.volumeLabelField.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
		if d.volumeDriverField.HasFocus() {
			if handler := d.volumeDriverField.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
		if d.volumeDriverOptionsField.HasFocus() {
			if handler := d.volumeDriverOptionsField.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
	})
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

// SetRect set rects for this primitive.
func (d *VolumeCreateDialog) SetRect(x, y, width, height int) {

	if width > volumeCreateDialogMaxWidth {
		emptySpace := (width - volumeCreateDialogMaxWidth) / 2
		x = x + emptySpace
		width = volumeCreateDialogMaxWidth
	}

	if height > volumeCreateDialogHeight {
		emptySpace := (height - volumeCreateDialogHeight) / 2
		y = y + emptySpace
		height = volumeCreateDialogHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen.
func (d *VolumeCreateDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}
	d.Box.DrawForSubclass(screen, d)
	x, y, width, height := d.Box.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetCancelFunc sets form cancel button selected function
func (d *VolumeCreateDialog) SetCancelFunc(handler func()) *VolumeCreateDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2)
	cancelButton.SetSelectedFunc(handler)
	return d
}

// SetCreateFunc sets form create button selected function
func (d *VolumeCreateDialog) SetCreateFunc(handler func()) *VolumeCreateDialog {
	d.createHandler = handler
	enterButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	enterButton.SetSelectedFunc(handler)
	return d
}

func (d *VolumeCreateDialog) initData() {
	d.volumeNameField.SetText("")
	d.volumeLabelField.SetText("")
	d.volumeDriverField.SetText("")
	d.volumeDriverOptionsField.SetText("")

}

// VolumeCreateOptions returns new volume options
func (d *VolumeCreateDialog) VolumeCreateOptions() volumes.CreateOptions {
	var (
		labels  = make(map[string]string)
		options = make(map[string]string)
	)
	for _, label := range strings.Split(d.volumeLabelField.GetText(), " ") {
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
	for _, option := range strings.Split(d.volumeDriverOptionsField.GetText(), " ") {
		if option != "" {
			split := strings.Split(option, "=")
			if len(split) == 2 {
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
