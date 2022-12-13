package imgdialogs

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/images"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	imagePushDialogMaxWidth  = 90
	imagePushDialogMaxHeight = 15
)

const (
	imagePushDesitnationFocus = 0 + iota
	imagePushCompressFocus
	imagePushFormatFocus
	imagePushSkipTLSVerifyFocus
	imagePushUsernameFocus
	imagePushPasswordFocus
	imagePushAuthFileFocus
	imagePushFormFocus
)

// ImagePushDialog represents image push dialog primitive
type ImagePushDialog struct {
	*tview.Box
	layout        *tview.Flex
	imageInfo     *tview.InputField
	destination   *tview.InputField
	compress      *tview.Checkbox
	format        *tview.DropDown
	skipTLSVerify *tview.Checkbox
	authFile      *tview.InputField
	username      *tview.InputField
	password      *tview.InputField
	form          *tview.Form
	display       bool
	pushHandler   func()
	cancelHandler func()
	focusElement  int
}

// NewImagePushDialog returns a new image push dialog primitive
func NewImagePushDialog() *ImagePushDialog {
	dialog := &ImagePushDialog{
		Box:           tview.NewBox(),
		layout:        tview.NewFlex(),
		imageInfo:     tview.NewInputField(),
		destination:   tview.NewInputField(),
		compress:      tview.NewCheckbox(),
		format:        tview.NewDropDown(),
		skipTLSVerify: tview.NewCheckbox(),
		authFile:      tview.NewInputField(),
		username:      tview.NewInputField(),
		password:      tview.NewInputField(),
		form:          tview.NewForm(),
	}

	bgColor := style.DialogBgColor
	fgColor := style.DialogFgColor
	inputFieldBgColor := style.InputFieldBgColor
	ddUnselectedStyle := style.DropDownUnselected
	ddselectedStyle := style.DropDownSelected
	labelWidth := 13

	// image info field
	dialog.imageInfo.SetBackgroundColor(style.DialogBgColor)
	dialog.imageInfo.SetLabel("[::b]IMAGE ID:")
	dialog.imageInfo.SetLabelWidth(labelWidth)
	dialog.imageInfo.SetFieldBackgroundColor(style.DialogBgColor)
	dialog.imageInfo.SetLabelStyle(tcell.StyleDefault.
		Background(style.DialogBorderColor).
		Foreground(style.DialogFgColor))

	// destination input field
	dialog.destination.SetBackgroundColor(bgColor)
	dialog.destination.SetLabelColor(fgColor)
	dialog.destination.SetLabel("destination:")
	dialog.destination.SetLabelWidth(labelWidth)
	dialog.destination.SetFieldBackgroundColor(inputFieldBgColor)

	// compress checkbox
	dialog.compress.SetBackgroundColor(bgColor)
	dialog.compress.SetLabelColor(fgColor)
	dialog.compress.SetLabel("compress:")
	dialog.compress.SetLabelWidth(labelWidth)
	dialog.compress.SetFieldBackgroundColor(inputFieldBgColor)

	// format dropdown
	formatLabel := "format:"
	dialog.format.SetLabel(formatLabel)
	dialog.format.SetTitleAlign(tview.AlignRight)
	dialog.format.SetLabelColor(fgColor)
	dialog.format.SetLabelWidth(len(formatLabel) + 1)
	dialog.format.SetBackgroundColor(bgColor)
	dialog.format.SetOptions([]string{
		"oci",
		"v2v2",
		"v2v1"},
		nil)
	dialog.format.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	dialog.format.SetFieldBackgroundColor(inputFieldBgColor)

	// skipTLSVerify checkbox
	skipTLSVerifyLabel := "skip tls verify:"
	dialog.skipTLSVerify.SetBackgroundColor(bgColor)
	dialog.skipTLSVerify.SetLabelColor(fgColor)
	dialog.skipTLSVerify.SetLabel(skipTLSVerifyLabel)
	dialog.skipTLSVerify.SetLabelWidth(len(skipTLSVerifyLabel) + 1)
	dialog.skipTLSVerify.SetFieldBackgroundColor(inputFieldBgColor)

	// authfile input field
	dialog.authFile.SetBackgroundColor(bgColor)
	dialog.authFile.SetLabelColor(fgColor)
	dialog.authFile.SetLabel("authfile:")
	dialog.authFile.SetLabelWidth(labelWidth)
	dialog.authFile.SetFieldBackgroundColor(inputFieldBgColor)

	// username input field
	dialog.username.SetBackgroundColor(bgColor)
	dialog.username.SetLabelColor(fgColor)
	dialog.username.SetLabel("username:")
	dialog.username.SetLabelWidth(labelWidth)
	dialog.username.SetFieldBackgroundColor(inputFieldBgColor)

	// password input field
	passwordLabel := "password:"
	dialog.password.SetBackgroundColor(bgColor)
	dialog.password.SetLabelColor(fgColor)
	dialog.password.SetLabel(passwordLabel)
	dialog.password.SetLabelWidth(len(passwordLabel) + 1)
	dialog.password.SetFieldBackgroundColor(inputFieldBgColor)
	dialog.password.SetMaskCharacter('*')

	// form
	dialog.form.AddButton("Cancel", nil)
	dialog.form.AddButton("Push", nil)
	dialog.form.SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(bgColor)
	dialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// layout
	// dropdowns and checkbox row layour
	dcLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	dcLayout.AddItem(dialog.compress, labelWidth+1, 1, true)
	dcLayout.AddItem(utils.EmptyBoxSpace(bgColor), 2, 0, false)
	dcLayout.AddItem(dialog.format, len(formatLabel)+5, 0, true)
	dcLayout.AddItem(utils.EmptyBoxSpace(bgColor), 2, 0, false)
	dcLayout.AddItem(dialog.skipTLSVerify, 0, 1, true)

	// username and password row layout
	userPassLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	userPassLayout.AddItem(dialog.username, 0, 1, true)
	userPassLayout.AddItem(utils.EmptyBoxSpace(bgColor), 3, 0, false)
	userPassLayout.AddItem(dialog.password, 0, 1, true)

	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, false)
	layout.AddItem(dialog.imageInfo, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, false)
	layout.AddItem(dialog.destination, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, false)
	layout.AddItem(dcLayout, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, false)
	layout.AddItem(userPassLayout, 0, 1, true)
	layout.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, false)
	layout.AddItem(dialog.authFile, 0, 1, true)

	inputLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	inputLayout.SetBackgroundColor(bgColor)
	inputLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	inputLayout.AddItem(layout, 0, 1, true)
	inputLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)

	dialog.layout.SetDirection(tview.FlexRow)
	dialog.layout.SetBackgroundColor(bgColor)
	dialog.layout.SetBorder(true)
	dialog.layout.SetBorderColor(style.DialogBorderColor)
	dialog.layout.SetTitle("PODMAN IMAGE PUSH")
	dialog.layout.AddItem(inputLayout, 0, 1, true)
	dialog.layout.AddItem(dialog.form, dialogs.DialogFormHeight, 0, true)

	dialog.Hide()
	return dialog
}

// Display displays this primitive
func (d *ImagePushDialog) Display() {
	d.display = true
}

// IsDisplay returns true if primitive is shown
func (d *ImagePushDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive
func (d *ImagePushDialog) Hide() {
	d.display = false
	d.focusElement = imagePushDesitnationFocus
	d.destination.SetText("")
	d.compress.SetChecked(false)
	d.format.SetCurrentOption(0)
	d.skipTLSVerify.SetChecked(false)
	d.authFile.SetText("")
	d.username.SetText("")
	d.password.SetText("")
}

// HasFocus returns whether or not this primitive has focus
func (d *ImagePushDialog) HasFocus() bool {
	if d.destination.HasFocus() || d.compress.HasFocus() {
		return true
	}
	if d.format.HasFocus() || d.skipTLSVerify.HasFocus() {
		return true
	}
	if d.username.HasFocus() || d.password.HasFocus() {
		return true
	}
	if d.authFile.HasFocus() || d.form.HasFocus() {
		return true
	}
	if d.layout.HasFocus() || d.Box.HasFocus() {
		return true
	}
	return d.Box.HasFocus()
}

// dropdownHasFocus returns true if image push dialog dropdown primitives
// has focus
func (d *ImagePushDialog) dropdownHasFocus() bool {
	return d.format.HasFocus()
}

// Focus is called when this primitive receives focus
func (d *ImagePushDialog) Focus(delegate func(p tview.Primitive)) {
	switch d.focusElement {
	case imagePushDesitnationFocus:
		delegate(d.destination)
	case imagePushCompressFocus:
		delegate(d.compress)
	case imagePushFormatFocus:
		delegate(d.format)
	case imagePushSkipTLSVerifyFocus:
		delegate(d.skipTLSVerify)
	case imagePushAuthFileFocus:
		delegate(d.authFile)
	case imagePushUsernameFocus:
		delegate(d.username)
	case imagePushPasswordFocus:
		delegate(d.password)
	case imagePushFormFocus:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == utils.SwitchFocusKey.Key {
				d.focusElement = imagePushDesitnationFocus
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
func (d *ImagePushDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("image push dialog: event %v received", event)
		if event.Key() == utils.SwitchFocusKey.Key {
			d.setFocusElement()
		}
		if event.Key() == tcell.KeyEsc && !d.dropdownHasFocus() {
			d.cancelHandler()
			return
		}
		if d.destination.HasFocus() {
			if destinationHandler := d.destination.InputHandler(); destinationHandler != nil {
				destinationHandler(event, setFocus)
				return
			}
		}
		if d.compress.HasFocus() {
			if compressHandler := d.compress.InputHandler(); compressHandler != nil {
				compressHandler(event, setFocus)
				return
			}
		}
		if d.format.HasFocus() {
			if formatHandler := d.format.InputHandler(); formatHandler != nil {
				event = utils.ParseKeyEventKey(event)
				formatHandler(event, setFocus)
				return
			}
		}
		if d.skipTLSVerify.HasFocus() {
			if skipTLSVerifyHandler := d.skipTLSVerify.InputHandler(); skipTLSVerifyHandler != nil {
				skipTLSVerifyHandler(event, setFocus)
				return
			}
		}
		if d.authFile.HasFocus() {
			if authFileHandler := d.authFile.InputHandler(); authFileHandler != nil {
				authFileHandler(event, setFocus)
				return
			}
		}
		if d.username.HasFocus() {
			if usernameHandler := d.username.InputHandler(); usernameHandler != nil {
				usernameHandler(event, setFocus)
				return
			}
		}
		if d.password.HasFocus() {
			if passwordHandler := d.password.InputHandler(); passwordHandler != nil {
				passwordHandler(event, setFocus)
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

func (d *ImagePushDialog) setFocusElement() {
	switch d.focusElement {
	case imagePushDesitnationFocus:
		d.focusElement = imagePushCompressFocus
	case imagePushCompressFocus:
		d.focusElement = imagePushFormatFocus
	case imagePushFormatFocus:
		d.focusElement = imagePushSkipTLSVerifyFocus
	case imagePushSkipTLSVerifyFocus:
		d.focusElement = imagePushUsernameFocus
	case imagePushUsernameFocus:
		d.focusElement = imagePushPasswordFocus
	case imagePushPasswordFocus:
		d.focusElement = imagePushAuthFileFocus
	case imagePushAuthFileFocus:
		d.focusElement = imagePushFormFocus
	}
}

// SetRect set rects for this primitive.
func (d *ImagePushDialog) SetRect(x, y, width, height int) {

	if width > imagePushDialogMaxWidth {
		emptySpace := (width - imagePushDialogMaxWidth) / 2
		x = x + emptySpace
		width = imagePushDialogMaxWidth
	}

	if height > imagePushDialogMaxHeight {
		emptySpace := (height - imagePushDialogMaxHeight) / 2
		y = y + emptySpace
		height = imagePushDialogMaxHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen.
func (d *ImagePushDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}
	d.Box.DrawForSubclass(screen, d)
	x, y, width, height := d.Box.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetPushFunc sets form push button selected function
func (d *ImagePushDialog) SetPushFunc(handler func()) *ImagePushDialog {
	d.pushHandler = handler
	pushButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	pushButton.SetSelectedFunc(handler)
	return d
}

// SetCancelFunc sets form cancel button selected function
func (d *ImagePushDialog) SetCancelFunc(handler func()) *ImagePushDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2)
	cancelButton.SetSelectedFunc(handler)
	return d
}

// SetImageInfo sets selected image ID and name in push dialog
func (d *ImagePushDialog) SetImageInfo(id string, name string) {
	containerInfo := fmt.Sprintf("%12s (%s)", id, name)
	d.imageInfo.SetText(containerInfo)
}

// GetImagePushOptions returns image push options based on user inputs
func (d *ImagePushDialog) GetImagePushOptions() images.ImagePushOptions {
	var opts images.ImagePushOptions

	opts.Destination = strings.TrimSpace(d.destination.GetText())
	_, format := d.format.GetCurrentOption()
	format = strings.TrimSpace(format)
	opts.Format = format
	opts.Compress = d.compress.IsChecked()
	opts.SkipTLSVerify = d.skipTLSVerify.IsChecked()
	opts.Username = strings.TrimSpace(d.username.GetText())
	opts.Password = strings.TrimSpace(d.password.GetText())
	opts.AuthFile = strings.TrimSpace(d.authFile.GetText())

	return opts
}
