package cntdialogs

import (
	"fmt"
	"strings"

	"github.com/containers/buildah/define"
	"github.com/containers/podman-tui/pdcs/containers"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	cntCommitDialogMaxWidth     = 90
	cntCommitDialogMaxHeight    = 15
	cntCommitDialogLabelPadding = 1
)

const (
	cntCommitImageFocus = 0 + iota
	cntCommitAuthorFocus
	cntCommitChangeFocus
	cntCommitFormatFocus
	cntCommitMessageFocus
	cntCommitPauseFocus
	cntCommitSquashFocus
	cntCommitFormFocus
)

// ContainerCommitDialog represents container commit dialog primitive.
type ContainerCommitDialog struct {
	*tview.Box

	layout        *tview.Flex
	cntInfo       *tview.InputField
	image         *tview.InputField
	author        *tview.InputField
	change        *tview.InputField
	format        *tview.DropDown
	message       *tview.InputField
	pause         *tview.Checkbox
	squash        *tview.Checkbox
	form          *tview.Form
	display       bool
	commitHandler func()
	cancelHandler func()
	focusElement  int
}

// NewContainerCommitDialog returns new container commit dialog primitive.
func NewContainerCommitDialog() *ContainerCommitDialog {
	dialog := &ContainerCommitDialog{
		Box:     tview.NewBox(),
		cntInfo: tview.NewInputField(),
		layout:  tview.NewFlex(),
		image:   tview.NewInputField(),
		author:  tview.NewInputField(),
		change:  tview.NewInputField(),
		format:  tview.NewDropDown(),
		message: tview.NewInputField(),
		pause:   tview.NewCheckbox(),
		squash:  tview.NewCheckbox(),
		form:    tview.NewForm(),
	}

	ddUnselectedStyle := style.DropDownUnselected
	ddselectedStyle := style.DropDownSelected
	labelWidth := 9

	// container info input field
	dialog.cntInfo.SetBackgroundColor(style.DialogBgColor)
	dialog.cntInfo.SetLabel("[::b]" + utils.ContainerIDLabel)
	dialog.cntInfo.SetFieldBackgroundColor(style.DialogBgColor)
	dialog.cntInfo.SetLabelStyle(tcell.StyleDefault.
		Background(style.DialogBorderColor).
		Foreground(style.DialogFgColor))

	// image field
	dialog.image.SetBackgroundColor(style.DialogBgColor)
	dialog.image.SetLabel(utils.StringToInputLabel("image:", labelWidth))
	dialog.image.SetFieldStyle(style.InputFieldStyle)
	dialog.image.SetLabelStyle(style.InputLabelStyle)

	// author field
	dialog.author.SetBackgroundColor(style.DialogBgColor)
	dialog.author.SetLabel(utils.StringToInputLabel("author:", labelWidth))
	dialog.author.SetFieldStyle(style.InputFieldStyle)
	dialog.author.SetLabelStyle(style.InputLabelStyle)

	// change field
	dialog.change.SetBackgroundColor(style.DialogBgColor)
	dialog.change.SetLabel(utils.StringToInputLabel("change:", labelWidth))
	dialog.change.SetFieldStyle(style.InputFieldStyle)
	dialog.change.SetLabelStyle(style.InputLabelStyle)

	// format options dropdown
	dialog.format.SetLabel("format:")
	dialog.format.SetTitleAlign(tview.AlignRight)
	dialog.format.SetLabelColor(style.DialogFgColor)
	dialog.format.SetLabelWidth(labelWidth)
	dialog.format.SetBackgroundColor(style.DialogBgColor)
	dialog.format.SetOptions([]string{
		define.OCI,
		define.DOCKER,
	},
		nil)

	dialog.format.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	dialog.format.SetFocusedStyle(style.DropDownFocused)
	dialog.format.SetFieldStyle(style.InputFieldStyle)

	// commit message field
	dialog.message.SetBackgroundColor(style.DialogBgColor)
	dialog.message.SetLabel(utils.StringToInputLabel("message:", labelWidth))
	dialog.message.SetFieldStyle(style.InputFieldStyle)
	dialog.message.SetLabelStyle(style.InputLabelStyle)

	// pause checkbox
	pauseLabel := "pause container:"

	dialog.pause.SetBackgroundColor(style.DialogBgColor)
	dialog.pause.SetLabelColor(style.DialogFgColor)
	dialog.pause.SetLabel(pauseLabel)
	dialog.pause.SetLabelWidth(len(pauseLabel) + 1)
	dialog.pause.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// squash checkbox
	squashLabel := "squash layers:"

	dialog.squash.SetBackgroundColor(style.DialogBgColor)
	dialog.squash.SetLabelColor(style.DialogFgColor)
	dialog.squash.SetLabel(squashLabel)
	dialog.squash.SetLabelWidth(len(squashLabel) + 1)
	dialog.squash.SetFieldBackgroundColor(style.FieldBackgroundColor)

	// form
	dialog.form.AddButton("Cancel", nil)
	dialog.form.AddButton("Commit", nil)
	dialog.form.SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(style.DialogBgColor)
	dialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// image and author layout row
	iaLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	iaLayout.SetBackgroundColor(style.DialogBgColor)
	iaLayout.AddItem(dialog.image, 0, 1, true)
	iaLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 2, 0, false) //nolint:mnd
	iaLayout.AddItem(dialog.author, 0, 1, true)

	// dropdown and checkbox layout row
	dcLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	dcLayout.SetBackgroundColor(style.DialogBgColor)
	dcLayout.AddItem(dialog.format, 0, 1, true)
	dcLayout.AddItem(dialog.squash, 0, 1, true)
	dcLayout.AddItem(dialog.pause, 0, 1, true)
	dcLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 0, 1, false)

	// inputs layout
	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	layout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 0, 1, false)
	layout.AddItem(dialog.cntInfo, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 0, 1, false)
	layout.AddItem(iaLayout, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 0, 1, false)
	layout.AddItem(dialog.change, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 0, 1, false)
	layout.AddItem(dcLayout, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 0, 1, false)
	layout.AddItem(dialog.message, 0, 1, true)

	inputLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	inputLayout.SetBackgroundColor(style.DialogBgColor)
	inputLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)
	inputLayout.AddItem(layout, 0, 1, true)
	inputLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)

	// main layout
	dialog.layout.SetDirection(tview.FlexRow)
	dialog.layout.SetBackgroundColor(style.DialogBgColor)
	dialog.layout.SetBorder(true)
	dialog.layout.SetBorderColor(style.DialogBorderColor)
	dialog.layout.SetTitle("PODMAN CONTAINER COMMIT")
	dialog.layout.AddItem(inputLayout, 0, 1, true)
	dialog.layout.AddItem(dialog.form, dialogs.DialogFormHeight, 0, true)

	dialog.Hide()

	return dialog
}

// Display displays this primitive.
func (d *ContainerCommitDialog) Display() {
	d.display = true
}

// IsDisplay returns true if primitive is shown.
func (d *ContainerCommitDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *ContainerCommitDialog) Hide() {
	d.display = false
	d.focusElement = cntCommitImageFocus
	d.image.SetText("")
	d.author.SetText("")
	d.change.SetText("")
	d.format.SetCurrentOption(0)
	d.message.SetText("")
	d.pause.SetChecked(false)
	d.squash.SetChecked(false)
	d.SetContainerInfo("", "")
}

// HasFocus returns whether or not this primitive has focus.
func (d *ContainerCommitDialog) HasFocus() bool { //nolint:cyclop
	if d.image.HasFocus() || d.author.HasFocus() {
		return true
	}

	if d.change.HasFocus() || d.format.HasFocus() {
		return true
	}

	if d.pause.HasFocus() || d.squash.HasFocus() {
		return true
	}

	if d.message.HasFocus() || d.form.HasFocus() {
		return true
	}

	if d.layout.HasFocus() || d.Box.HasFocus() {
		return true
	}

	return false
}

// Focus is called when this primitive receives focus.
func (d *ContainerCommitDialog) Focus(delegate func(p tview.Primitive)) {
	switch d.focusElement {
	case cntCommitImageFocus:
		delegate(d.image)
	case cntCommitAuthorFocus:
		delegate(d.author)
	case cntCommitChangeFocus:
		delegate(d.change)
	case cntCommitFormatFocus:
		delegate(d.format)
	case cntCommitMessageFocus:
		delegate(d.message)
	case cntCommitPauseFocus:
		delegate(d.pause)
	case cntCommitSquashFocus:
		delegate(d.squash)
	case cntCommitFormFocus:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == utils.SwitchFocusKey.Key {
				d.focusElement = cntCommitImageFocus

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
func (d *ContainerCommitDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) { //nolint:gocognit,lll,cyclop
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("container commit dialog: event %v received", event)

		if event.Key() == utils.SwitchFocusKey.Key {
			d.setFocusElement()
		}

		// dropdown widgets shall handle events before "Esc" key handler
		if d.format.HasFocus() {
			event = utils.ParseKeyEventKey(event)
			if formatHandler := d.format.InputHandler(); formatHandler != nil {
				formatHandler(event, setFocus)

				return
			}
		}

		if event.Key() == tcell.KeyEsc {
			d.cancelHandler()

			return
		}

		if d.image.HasFocus() {
			if imageHandler := d.image.InputHandler(); imageHandler != nil {
				imageHandler(event, setFocus)

				return
			}
		}

		if d.author.HasFocus() {
			if authorHandler := d.author.InputHandler(); authorHandler != nil {
				authorHandler(event, setFocus)

				return
			}
		}

		if d.change.HasFocus() {
			if changeHandler := d.change.InputHandler(); changeHandler != nil {
				changeHandler(event, setFocus)

				return
			}
		}

		if d.message.HasFocus() {
			if messageHandler := d.message.InputHandler(); messageHandler != nil {
				messageHandler(event, setFocus)

				return
			}
		}

		if d.pause.HasFocus() {
			if pauseHandler := d.pause.InputHandler(); pauseHandler != nil {
				pauseHandler(event, setFocus)

				return
			}
		}

		if d.squash.HasFocus() {
			if squashHandler := d.squash.InputHandler(); squashHandler != nil {
				squashHandler(event, setFocus)

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
func (d *ContainerCommitDialog) SetRect(x, y, width, height int) {
	if width > cntCommitDialogMaxWidth {
		emptySpace := (width - cntCommitDialogMaxWidth) / 2 //nolint:mnd
		x += emptySpace
		width = cntCommitDialogMaxWidth
	}

	if height > cntCommitDialogMaxHeight {
		emptySpace := (height - cntCommitDialogMaxHeight) / 2 //nolint:mnd
		y += emptySpace
		height = cntCommitDialogMaxHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen.
func (d *ContainerCommitDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.DrawForSubclass(screen, d)

	x, y, width, height := d.GetInnerRect()

	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetCommitFunc sets form commit button selected function.
func (d *ContainerCommitDialog) SetCommitFunc(handler func()) *ContainerCommitDialog {
	d.commitHandler = handler
	commitButton := d.form.GetButton(d.form.GetButtonCount() - 1)

	commitButton.SetSelectedFunc(handler)

	return d
}

// SetCancelFunc sets form cancel button selected function.
func (d *ContainerCommitDialog) SetCancelFunc(handler func()) *ContainerCommitDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:mnd

	cancelButton.SetSelectedFunc(handler)

	return d
}

// SetContainerInfo sets selected container ID and name in commit dialog.
func (d *ContainerCommitDialog) SetContainerInfo(id string, name string) {
	containerInfo := fmt.Sprintf("%s (%s)", id, name)
	containerInfo = utils.LabelWidthLeftPadding(containerInfo, cntCommitDialogLabelPadding)

	d.cntInfo.SetText(containerInfo)
}

// GetContainerCommitOptions returns container commit options based on user inputs.
func (d *ContainerCommitDialog) GetContainerCommitOptions() containers.CntCommitOptions {
	var opts containers.CntCommitOptions

	opts.Image = strings.TrimSpace(d.image.GetText())
	opts.Author = strings.TrimSpace(d.author.GetText())
	opts.Changes = strings.Split(d.change.GetText(), " ")
	_, format := d.format.GetCurrentOption()

	switch format {
	case "oci":
		opts.Format = define.OCIv1ImageManifest
	case "docker":
		opts.Format = define.Dockerv2ImageManifest
	}

	opts.Pause = d.pause.IsChecked()
	opts.Squash = d.squash.IsChecked()
	opts.Message = strings.TrimSpace(d.message.GetText())

	return opts
}

func (d *ContainerCommitDialog) setFocusElement() {
	switch d.focusElement {
	case cntCommitImageFocus:
		d.focusElement = cntCommitAuthorFocus
	case cntCommitAuthorFocus:
		d.focusElement = cntCommitChangeFocus
	case cntCommitChangeFocus:
		d.focusElement = cntCommitFormatFocus
	case cntCommitFormatFocus:
		d.focusElement = cntCommitSquashFocus
	case cntCommitSquashFocus:
		d.focusElement = cntCommitPauseFocus
	case cntCommitPauseFocus:
		d.focusElement = cntCommitMessageFocus
	case cntCommitMessageFocus:
		d.focusElement = cntCommitFormFocus
	}
}
