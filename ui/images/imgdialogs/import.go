package imgdialogs

import (
	"errors"
	"strings"

	"github.com/containers/podman-tui/pdcs/images"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/hashicorp/go-multierror"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	imageImportDialogMaxWidth  = 70
	imageImportDialogMaxHeight = 13
)

var errImportEmptySource = errors.New("empty source value for the image tarball")

const (
	imageImportPathFocus = 0 + iota
	imageImportCommitMessageFocus
	imageImportChangeFocus
	imageImportReferenceFocus
	imageImportFormFocus
)

// ImageImportDialog represents image import dialog primitive.
type ImageImportDialog struct {
	*tview.Box
	layout        *tview.Flex
	path          *tview.InputField
	change        *tview.InputField
	commitMessage *tview.InputField
	reference     *tview.InputField
	form          *tview.Form
	display       bool
	importHandler func()
	cancelHandler func()
	focusElement  int
}

// NewImageImportDialog returns new image import dialog.
func NewImageImportDialog() *ImageImportDialog {
	dialog := &ImageImportDialog{
		Box:           tview.NewBox(),
		layout:        tview.NewFlex(),
		path:          tview.NewInputField(),
		change:        tview.NewInputField(),
		reference:     tview.NewInputField(),
		commitMessage: tview.NewInputField(),
		form:          tview.NewForm(),
	}

	bgColor := style.DialogBgColor
	fgColor := style.DialogFgColor
	inputFieldBgColor := style.InputFieldBgColor
	labelWidth := 11

	// path field
	dialog.path.SetBackgroundColor(bgColor)
	dialog.path.SetLabelColor(fgColor)
	dialog.path.SetLabel("source:")
	dialog.path.SetLabelWidth(labelWidth)
	dialog.path.SetFieldBackgroundColor(inputFieldBgColor)

	// change field
	dialog.change.SetBackgroundColor(bgColor)
	dialog.change.SetLabelColor(fgColor)
	dialog.change.SetLabel("change:")
	dialog.change.SetLabelWidth(labelWidth)
	dialog.change.SetFieldBackgroundColor(inputFieldBgColor)

	// commit field
	dialog.commitMessage.SetBackgroundColor(bgColor)
	dialog.commitMessage.SetLabelColor(fgColor)
	dialog.commitMessage.SetLabel("message:")
	dialog.commitMessage.SetLabelWidth(labelWidth)
	dialog.commitMessage.SetFieldBackgroundColor(inputFieldBgColor)

	// reference field
	dialog.reference.SetBackgroundColor(bgColor)
	dialog.reference.SetLabelColor(fgColor)
	dialog.reference.SetLabel("reference:")
	dialog.reference.SetLabelWidth(labelWidth)
	dialog.reference.SetFieldBackgroundColor(inputFieldBgColor)

	// form
	dialog.form.AddButton("Cancel", nil)
	dialog.form.AddButton("Import", nil)
	dialog.form.SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(bgColor)
	dialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// layout
	optionsLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(dialog.path, 1, 0, true)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(dialog.change, 1, 0, true)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(dialog.commitMessage, 1, 0, true)
	optionsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	optionsLayout.AddItem(dialog.reference, 1, 0, true)

	mainOptsLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	mainOptsLayout.SetBackgroundColor(bgColor)
	mainOptsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	mainOptsLayout.AddItem(optionsLayout, 0, 1, true)
	mainOptsLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)

	dialog.layout.SetDirection(tview.FlexRow)
	dialog.layout.SetBackgroundColor(bgColor)
	dialog.layout.SetBorder(true)
	dialog.layout.SetBorderColor(style.DialogBorderColor)
	dialog.layout.SetTitle("PODMAN IMAGE IMPORT")
	dialog.layout.AddItem(mainOptsLayout, 0, 1, true)
	dialog.layout.AddItem(dialog.form, dialogs.DialogFormHeight, 0, true)

	return dialog
}

// Display displays this primitive.
func (d *ImageImportDialog) Display() {
	d.display = true
}

// IsDisplay returns true if primitive is shown.
func (d *ImageImportDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *ImageImportDialog) Hide() {
	d.display = false
	d.focusElement = imageImportPathFocus

	d.path.SetText("")
	d.change.SetText("")
	d.commitMessage.SetText("")
	d.reference.SetText("")
}

// HasFocus returns whether or not this primitive has focus.
func (d *ImageImportDialog) HasFocus() bool {
	if d.path.HasFocus() || d.commitMessage.HasFocus() {
		return true
	}

	if d.form.HasFocus() || d.reference.HasFocus() {
		return true
	}

	if d.change.HasFocus() || d.layout.HasFocus() {
		return true
	}

	return d.Box.HasFocus()
}

// Focus is called when this primitive receives focus.
func (d *ImageImportDialog) Focus(delegate func(p tview.Primitive)) {
	switch d.focusElement {
	case imageImportPathFocus:
		delegate(d.path)
	case imageImportChangeFocus:
		delegate(d.change)
	case imageImportCommitMessageFocus:
		delegate(d.commitMessage)
	case imageImportReferenceFocus:
		delegate(d.reference)
	case imageImportFormFocus:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == utils.SwitchFocusKey.Key {
				d.focusElement = imageImportPathFocus
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
func (d *ImageImportDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) { //nolint:cyclop,lll
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("image import dialog: event %v received", event)

		if event.Key() == tcell.KeyEsc {
			d.cancelHandler()

			return
		}

		if event.Key() == utils.SwitchFocusKey.Key {
			d.setFocusElement()
		}

		if d.path.HasFocus() {
			if pathHandler := d.path.InputHandler(); pathHandler != nil {
				pathHandler(event, setFocus)

				return
			}
		}

		if d.change.HasFocus() {
			if changeHandler := d.change.InputHandler(); changeHandler != nil {
				changeHandler(event, setFocus)

				return
			}
		}

		if d.commitMessage.HasFocus() {
			if commitHandler := d.commitMessage.InputHandler(); commitHandler != nil {
				commitHandler(event, setFocus)

				return
			}
		}

		if d.reference.HasFocus() {
			if referenceHandler := d.reference.InputHandler(); referenceHandler != nil {
				referenceHandler(event, setFocus)

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
func (d *ImageImportDialog) SetRect(x, y, width, height int) {
	if width > imageImportDialogMaxWidth {
		emptySpace := (width - imageImportDialogMaxWidth) / 2 //nolint:mnd
		x += emptySpace
		width = imageImportDialogMaxWidth
	}

	if height > imageImportDialogMaxHeight {
		emptySpace := (height - imageImportDialogMaxHeight) / 2 //nolint:mnd
		y += emptySpace
		height = imageImportDialogMaxHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen.
func (d *ImageImportDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.Box.DrawForSubclass(screen, d)
	x, y, width, height := d.Box.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

func (d *ImageImportDialog) setFocusElement() {
	switch d.focusElement {
	case imageImportPathFocus:
		d.focusElement = imageImportChangeFocus
	case imageImportChangeFocus:
		d.focusElement = imageImportCommitMessageFocus
	case imageImportCommitMessageFocus:
		d.focusElement = imageImportReferenceFocus
	case imageImportReferenceFocus:
		d.focusElement = imageImportFormFocus
	}
}

// SetImportFunc sets form import button selected function.
func (d *ImageImportDialog) SetImportFunc(handler func()) *ImageImportDialog {
	d.importHandler = handler
	importButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	importButton.SetSelectedFunc(handler)

	return d
}

// SetCancelFunc sets form cancel button selected function.
func (d *ImageImportDialog) SetCancelFunc(handler func()) *ImageImportDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:mnd
	cancelButton.SetSelectedFunc(handler)

	return d
}

// ImageImportOptions return image import options.
func (d *ImageImportDialog) ImageImportOptions() (images.ImageImportOptions, error) {
	var (
		path      string
		change    []string
		commit    string
		reference string
	)

	commit = strings.TrimSpace(d.commitMessage.GetText())
	reference = strings.TrimSpace(d.reference.GetText())
	change = strings.Split(d.change.GetText(), " ")

	opts := images.ImageImportOptions{
		Change:    change,
		Message:   commit,
		Reference: reference,
	}

	path = strings.TrimSpace(d.path.GetText())
	if path == "" {
		return opts, errImportEmptySource
	}

	path, err := utils.ResolveHomeDir(path)
	if err != nil {
		return opts, err
	}

	errFileName := utils.ValidateFileName(path)
	errURL := utils.ValidURL(path)

	if errURL == nil {
		opts.URL = true
	}

	if errFileName != nil && errURL != nil {
		return opts, multierror.Append(errFileName, errURL)
	}

	opts.Source = path

	return opts, nil
}
