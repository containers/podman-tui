package imgdialogs

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/images"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/containers/podman/v4/libpod/define"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	imageSaveDialogMaxWidth  = 70
	imageSaveDialogMaxHeight = 13
)

const (
	imageSaveOutputFocus = 0 + iota
	imageSaveCompressFocus
	imageSaveAcceptUncompressedFocus
	imageSaveFormatFocus
	imageSaveFormFocus
)

// ImageSaveDialog represents image save dialog primitive
type ImageSaveDialog struct {
	*tview.Box
	layout                *tview.Flex
	imageInfo             *tview.InputField
	output                *tview.InputField
	compress              *tview.Checkbox
	format                *tview.DropDown
	ociAcceptUncompressed *tview.Checkbox
	form                  *tview.Form
	display               bool
	saveHandler           func()
	cancelHandler         func()
	focusElement          int
}

// NewImageSaveDialog returns new image save dialog
func NewImageSaveDialog() *ImageSaveDialog {
	dialog := &ImageSaveDialog{
		Box:                   tview.NewBox(),
		layout:                tview.NewFlex(),
		imageInfo:             tview.NewInputField(),
		output:                tview.NewInputField(),
		compress:              tview.NewCheckbox(),
		format:                tview.NewDropDown(),
		ociAcceptUncompressed: tview.NewCheckbox(),
		form:                  tview.NewForm(),
	}

	bgColor := style.DialogBgColor
	fgColor := style.DialogFgColor
	inputFieldBgColor := style.InputFieldBgColor
	ddUnselectedStyle := style.DropDownUnselected
	ddselectedStyle := style.DropDownSelected
	labelWidth := 10

	// image info
	imageInfoLabel := "IMAGE ID:"
	dialog.imageInfo.SetBackgroundColor(style.DialogBgColor)
	dialog.imageInfo.SetLabel("[::b]" + imageInfoLabel)
	dialog.imageInfo.SetLabelWidth(len(imageInfoLabel) + 1)
	dialog.imageInfo.SetFieldBackgroundColor(style.DialogBgColor)
	dialog.imageInfo.SetLabelStyle(tcell.StyleDefault.
		Background(style.DialogBorderColor).
		Foreground(style.DialogFgColor))

	// output
	dialog.output.SetBackgroundColor(bgColor)
	dialog.output.SetLabelColor(fgColor)
	dialog.output.SetLabel("output:")
	dialog.output.SetLabelWidth(labelWidth)
	dialog.output.SetFieldBackgroundColor(inputFieldBgColor)

	// compress
	dialog.compress.SetBackgroundColor(bgColor)
	dialog.compress.SetLabelColor(fgColor)
	dialog.compress.SetLabel("compress:")
	dialog.compress.SetLabelWidth(labelWidth)
	dialog.compress.SetFieldBackgroundColor(inputFieldBgColor)

	// format
	dialog.format.SetBackgroundColor(bgColor)
	dialog.format.SetLabelColor(fgColor)
	dialog.format.SetLabel("format:")
	dialog.format.SetLabelWidth(labelWidth)
	dialog.format.SetOptions([]string{
		define.V2s2Archive,
		define.V2s2ManifestDir,
		define.OCIArchive,
		define.OCIManifestDir},
		nil)
	dialog.format.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	dialog.format.SetCurrentOption(0)
	dialog.format.SetFieldBackgroundColor(inputFieldBgColor)

	// OciAcceptUncompressed
	dialog.ociAcceptUncompressed.SetBackgroundColor(bgColor)
	dialog.ociAcceptUncompressed.SetLabelColor(fgColor)
	dialog.ociAcceptUncompressed.SetLabel("accept uncompressed (OCI images): ")
	dialog.ociAcceptUncompressed.SetFieldBackgroundColor(inputFieldBgColor)

	// form
	dialog.form.AddButton("Cancel", nil)
	dialog.form.AddButton("Save", nil)
	dialog.form.SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(bgColor)
	dialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// layout
	compressRow := tview.NewFlex().SetDirection(tview.FlexColumn)
	compressRow.SetBackgroundColor(bgColor)
	compressRow.AddItem(dialog.compress, 0, 1, true)
	compressRow.AddItem(dialog.ociAcceptUncompressed, 0, 3, true)

	optionsLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(dialog.imageInfo, 0, 1, true)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(dialog.output, 0, 1, true)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(compressRow, 0, 1, true)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(dialog.format, 0, 1, true)

	mainOptsLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	mainOptsLayout.SetBackgroundColor(bgColor)
	mainOptsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	mainOptsLayout.AddItem(optionsLayout, 0, 1, true)
	mainOptsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)

	dialog.layout.SetDirection(tview.FlexRow)
	dialog.layout.SetBackgroundColor(bgColor)
	dialog.layout.SetBorder(true)
	dialog.layout.SetBorderColor(style.DialogBorderColor)
	dialog.layout.SetTitle("PODMAN IMAGE SAVE")
	dialog.layout.AddItem(mainOptsLayout, 0, 1, true)
	dialog.layout.AddItem(dialog.form, dialogs.DialogFormHeight, 0, true)

	return dialog
}

// Display displays this primitive
func (d *ImageSaveDialog) Display() {
	d.display = true
}

// IsDisplay returns true if primitive is shown
func (d *ImageSaveDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive
func (d *ImageSaveDialog) Hide() {
	d.display = false
	d.SetImageInfo("", "")
	d.focusElement = imageSaveOutputFocus
	d.output.SetText("")
	d.compress.SetChecked(false)
	d.ociAcceptUncompressed.SetChecked(false)
	d.format.SetCurrentOption(0)
}

// HasFocus returns whether or not this primitive has focus
func (d *ImageSaveDialog) HasFocus() bool {
	if d.output.HasFocus() || d.compress.HasFocus() {
		return true
	}
	if d.format.HasFocus() || d.ociAcceptUncompressed.HasFocus() {
		return true
	}
	if d.form.HasFocus() || d.layout.HasFocus() {
		return true
	}
	return d.Box.HasFocus()
}

// Focus is called when this primitive receives focus
func (d *ImageSaveDialog) Focus(delegate func(p tview.Primitive)) {
	switch d.focusElement {
	case imageSaveOutputFocus:
		delegate(d.output)
	case imageSaveCompressFocus:
		delegate(d.compress)
	case imageSaveAcceptUncompressedFocus:
		delegate(d.ociAcceptUncompressed)
	case imageSaveFormatFocus:
		delegate(d.format)
	case imageSaveFormFocus:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == utils.SwitchFocusKey.Key {
				d.focusElement = imageSaveOutputFocus
				d.Focus(delegate)
				d.form.SetFocus(0)
				return nil
			}
			return event
		})
		delegate(d.form)
	}
}

// InputHandler returns input handler function for this primitive
func (d *ImageSaveDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("image save dialog: event %v received", event)
		if event.Key() == tcell.KeyEsc {
			if !d.format.HasFocus() {
				d.cancelHandler()
				return
			}
		}
		if event.Key() == utils.SwitchFocusKey.Key {
			d.setFocusElement()
		}
		// drop down event
		if d.format.HasFocus() {
			event = utils.ParseKeyEventKey(event)
			if formatHandler := d.format.InputHandler(); formatHandler != nil {
				formatHandler(event, setFocus)
				return
			}
		}
		if d.output.HasFocus() {
			if outputHandler := d.output.InputHandler(); outputHandler != nil {
				outputHandler(event, setFocus)
				return
			}
		}
		if d.compress.HasFocus() {
			if compressHandler := d.compress.InputHandler(); compressHandler != nil {
				compressHandler(event, setFocus)
				return
			}
		}
		if d.ociAcceptUncompressed.HasFocus() {
			if acceptUncompressedhandler := d.ociAcceptUncompressed.InputHandler(); acceptUncompressedhandler != nil {
				acceptUncompressedhandler(event, setFocus)
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
func (d *ImageSaveDialog) SetRect(x, y, width, height int) {

	if width > imageSaveDialogMaxWidth {
		emptySpace := (width - imageSaveDialogMaxWidth) / 2
		x = x + emptySpace
		width = imageSaveDialogMaxWidth
	}

	if height > imageSaveDialogMaxHeight {
		emptySpace := (height - imageSaveDialogMaxHeight) / 2
		y = y + emptySpace
		height = imageSaveDialogMaxHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen.
func (d *ImageSaveDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}
	d.Box.DrawForSubclass(screen, d)
	x, y, width, height := d.Box.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

func (d *ImageSaveDialog) setFocusElement() {
	switch d.focusElement {
	case imageSaveOutputFocus:
		d.focusElement = imageSaveCompressFocus
	case imageSaveCompressFocus:
		d.focusElement = imageSaveAcceptUncompressedFocus
	case imageSaveAcceptUncompressedFocus:
		d.focusElement = imageSaveFormatFocus
	case imageSaveFormatFocus:
		d.focusElement = imageSaveFormFocus
	}
}

// SetSaveFunc sets form save button selected function
func (d *ImageSaveDialog) SetSaveFunc(handler func()) *ImageSaveDialog {
	d.saveHandler = handler
	saveButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	saveButton.SetSelectedFunc(handler)
	return d
}

// SetCancelFunc sets form cancel button selected function
func (d *ImageSaveDialog) SetCancelFunc(handler func()) *ImageSaveDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2)
	cancelButton.SetSelectedFunc(handler)
	return d
}

// SetImageInfo sets selected image ID and name in save dialog
func (d *ImageSaveDialog) SetImageInfo(id string, name string) {
	nameSplited := strings.Split(name, "/")
	l := len(nameSplited)
	if l > 1 {
		name = nameSplited[l-1]
	}
	imageInfo := fmt.Sprintf("Image ID: %s (%s)", id, name)
	d.imageInfo.SetText(imageInfo)
}

// ImageSaveOptions prepare and returns image save options
func (d *ImageSaveDialog) ImageSaveOptions() (images.ImageSaveOptions, error) {

	opts := images.ImageSaveOptions{
		Compressed:                  d.compress.IsChecked(),
		OciAcceptUncompressedLayers: d.ociAcceptUncompressed.IsChecked(),
	}
	_, format := d.format.GetCurrentOption()
	opts.Format = format

	output := strings.TrimSpace(d.output.GetText())
	if output == "" {
		return opts, fmt.Errorf("empty output name")
	}
	outputPath, err := utils.ResolveHomeDir(output)
	if err != nil {
		return opts, err
	}
	opts.Output = outputPath

	return opts, nil
}
