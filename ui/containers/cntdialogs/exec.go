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
	// maxheight = button height + total input widgets row + 11
	execDialogMaxHeight  = dialogs.DialogFormHeight + 7 + 9
	execDialogMaxWidth   = 80
	execDialogLabelWidth = 13
)

const (
	execCommandFieldFocus = 0 + iota
	execInteractiveFieldFocus
	execTtyFieldFocus
	execPrivilegedFieldFocus
	execWorkingDirFieldFocus
	execEnvVariablesFieldFocus
	execEnvFileFieldFocus
	execUserFieldFocus
	execDetachFieldFocus
	execFormFieldFocus
)

// ContainerExecDialog represents container exec dialog primitive
type ContainerExecDialog struct {
	*tview.Box
	layout        *tview.Flex
	label         *tview.TextView
	command       *tview.InputField
	interactive   *tview.Checkbox
	tty           *tview.Checkbox
	privileged    *tview.Checkbox
	detach        *tview.Checkbox
	workingDir    *tview.InputField
	envVariables  *tview.InputField
	envFile       *tview.InputField
	user          *tview.InputField
	form          *tview.Form
	display       bool
	containerID   string
	focusElement  int
	execHandler   func()
	cancelHandler func()
}

// NewContainerExecDialog returns new container exec dialog
func NewContainerExecDialog() *ContainerExecDialog {
	dialog := &ContainerExecDialog{
		Box:          tview.NewBox(),
		label:        tview.NewTextView(),
		command:      tview.NewInputField(),
		tty:          tview.NewCheckbox(),
		interactive:  tview.NewCheckbox(),
		privileged:   tview.NewCheckbox(),
		detach:       tview.NewCheckbox(),
		workingDir:   tview.NewInputField(),
		envVariables: tview.NewInputField(),
		envFile:      tview.NewInputField(),
		user:         tview.NewInputField(),
		display:      false,
	}
	bgColor := utils.Styles.ContainerExecDialog.BgColor
	fgColor := utils.Styles.ContainerExecDialog.FgColor
	inputFieldBgColor := utils.Styles.InputFieldPrimitive.BgColor

	// label (container ID and Name)
	dialog.label.SetDynamicColors(true)
	dialog.label.SetBackgroundColor(bgColor)
	dialog.label.SetBorder(false)
	dialog.SetContainerID("", "")

	// command
	dialog.command.SetBackgroundColor(bgColor)
	dialog.command.SetBorder(false)
	dialog.command.SetLabel("command:")
	dialog.command.SetLabelColor(fgColor)
	dialog.command.SetLabelWidth(execDialogLabelWidth)
	dialog.command.SetFieldBackgroundColor(inputFieldBgColor)

	// interactive
	dialog.interactive.SetBackgroundColor(bgColor)
	dialog.interactive.SetBorder(false)
	dialog.interactive.SetLabel("interactive:")
	dialog.interactive.SetLabelColor(fgColor)
	dialog.interactive.SetLabelWidth(execDialogLabelWidth)
	dialog.interactive.SetFieldBackgroundColor(inputFieldBgColor)

	// tty
	tLabel := "tty:"
	dialog.tty.SetBackgroundColor(bgColor)
	dialog.tty.SetBorder(false)
	dialog.tty.SetLabel(tLabel)
	dialog.tty.SetLabelColor(fgColor)
	dialog.tty.SetLabelWidth(len(tLabel) + 1)
	dialog.tty.SetFieldBackgroundColor(inputFieldBgColor)

	// privileged
	pLabel := "privileged:"
	dialog.privileged.SetBackgroundColor(bgColor)
	dialog.privileged.SetBorder(false)
	dialog.privileged.SetLabel(pLabel)
	dialog.privileged.SetLabelColor(fgColor)
	dialog.privileged.SetLabelWidth(len(pLabel) + 1)
	dialog.privileged.SetFieldBackgroundColor(inputFieldBgColor)

	// detach
	dLabel := "detach:"
	dialog.detach.SetBackgroundColor(bgColor)
	dialog.detach.SetBorder(false)
	dialog.detach.SetLabel(dLabel)
	dialog.detach.SetLabelColor(fgColor)
	dialog.detach.SetLabelWidth(len(dLabel) + 1)
	dialog.detach.SetFieldBackgroundColor(inputFieldBgColor)

	// working dir
	dialog.workingDir.SetBackgroundColor(bgColor)
	dialog.workingDir.SetBorder(false)
	dialog.workingDir.SetLabel("working dir:")
	dialog.workingDir.SetLabelColor(fgColor)
	dialog.workingDir.SetLabelWidth(execDialogLabelWidth)
	dialog.workingDir.SetFieldBackgroundColor(inputFieldBgColor)

	// env variables
	dialog.envVariables.SetBackgroundColor(bgColor)
	dialog.envVariables.SetBorder(false)
	dialog.envVariables.SetLabel("env vars:")
	dialog.envVariables.SetLabelColor(fgColor)
	dialog.envVariables.SetLabelWidth(execDialogLabelWidth)
	dialog.envVariables.SetFieldBackgroundColor(inputFieldBgColor)

	// env file
	dialog.envFile.SetBackgroundColor(bgColor)
	dialog.envFile.SetBorder(false)
	dialog.envFile.SetLabel("env file:")
	dialog.envFile.SetLabelColor(fgColor)
	dialog.envFile.SetLabelWidth(execDialogLabelWidth)
	dialog.envFile.SetFieldBackgroundColor(inputFieldBgColor)

	// user
	dialog.user.SetBackgroundColor(bgColor)
	dialog.user.SetBorder(false)
	dialog.user.SetLabel("user: ")
	dialog.user.SetLabelColor(fgColor)
	dialog.user.SetLabelWidth(execDialogLabelWidth)
	dialog.user.SetFieldBackgroundColor(inputFieldBgColor)

	// form fields
	dialog.form = tview.NewForm().
		AddButton("Cancel", nil).
		AddButton("Execute", nil).
		SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(bgColor)
	dialog.form.SetButtonBackgroundColor(utils.Styles.ButtonPrimitive.BgColor)

	// main dialog layout
	dialog.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	dialog.layout.SetBorder(true)
	dialog.layout.SetBackgroundColor(bgColor)
	dialog.layout.SetTitle("PODMAN CONTAINER EXEC")

	mLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	// label
	mLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	mLayout.AddItem(dialog.label, 1, 0, true)
	// command
	mLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	mLayout.AddItem(dialog.command, 1, 0, true)
	// interactive, tty, privileged and detach
	checkBoxWidth := execDialogLabelWidth + 4
	labelPaddings := 5
	checkBoxLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	checkBoxLayout.SetBackgroundColor(bgColor)
	checkBoxLayout.AddItem(dialog.interactive, checkBoxWidth, 0, false)
	checkBoxLayout.AddItem(dialog.tty, len(tLabel)+labelPaddings, 0, false)
	checkBoxLayout.AddItem(dialog.privileged, len(pLabel)+labelPaddings, 0, false)
	checkBoxLayout.AddItem(dialog.detach, len(dLabel)+labelPaddings, 0, false)
	checkBoxLayout.AddItem(utils.EmptyBoxSpace(bgColor), 0, 1, true)
	mLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	mLayout.AddItem(checkBoxLayout, 1, 0, true)
	// user
	mLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	mLayout.AddItem(dialog.user, 1, 0, true)
	// working dir
	mLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	mLayout.AddItem(dialog.workingDir, 1, 0, true)
	// env variables
	mLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	mLayout.AddItem(dialog.envVariables, 1, 0, true)
	// env file
	mLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	mLayout.AddItem(dialog.envFile, 1, 0, true)

	// main layout
	mainLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	mainLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	mainLayout.AddItem(mLayout, 0, 1, true)
	mainLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)

	dialog.layout.AddItem(mainLayout, 0, 1, true)
	dialog.layout.AddItem(dialog.form, dialogs.DialogFormHeight, 0, true)
	return dialog
}

// Display displays this primitive
func (d *ContainerExecDialog) Display() {
	d.focusElement = execCommandFieldFocus
	d.command.SetText("")
	d.tty.SetChecked(true)
	d.interactive.SetChecked(false)
	d.privileged.SetChecked(false)
	d.detach.SetChecked(false)
	d.workingDir.SetText("")
	d.envVariables.SetText("")
	d.envFile.SetText("")
	d.user.SetText("")
	d.display = true
}

// IsDisplay returns true if primitive is shown
func (d *ContainerExecDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive
func (d *ContainerExecDialog) Hide() {
	d.SetContainerID("", "")
	d.display = false

}

// HasFocus returns whether or not this primitive has focus
func (d *ContainerExecDialog) HasFocus() bool {
	if d.command.HasFocus() || d.tty.HasFocus() {
		return true
	}
	if d.interactive.HasFocus() || d.privileged.HasFocus() {
		return true
	}
	if d.workingDir.HasFocus() || d.envVariables.HasFocus() {
		return true
	}
	if d.envFile.HasFocus() || d.user.HasFocus() {
		return true
	}
	if d.detach.HasFocus() || d.form.HasFocus() {
		return true
	}
	return d.Box.HasFocus() || d.layout.HasFocus()
}

// Focus is called when this primitive receives focus
func (d *ContainerExecDialog) Focus(delegate func(p tview.Primitive)) {
	switch d.focusElement {
	// command field focus
	case execCommandFieldFocus:
		d.command.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = execInteractiveFieldFocus
				d.Focus(delegate)
				return nil
			}
			if event.Key() == tcell.KeyEnter {
				d.execHandler()
				return nil
			}
			return event
		})
		delegate(d.command)
		return
	// interactive field focus
	case execInteractiveFieldFocus:
		d.interactive.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = execTtyFieldFocus
				d.Focus(delegate)
				return nil
			}
			return event
		})
		delegate(d.interactive)
		return
	// tty field focus
	case execTtyFieldFocus:
		d.tty.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = execPrivilegedFieldFocus
				d.Focus(delegate)
				return nil
			}
			return event
		})
		delegate(d.tty)
		return

	// privileged field focus
	case execPrivilegedFieldFocus:
		d.privileged.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = execDetachFieldFocus
				d.Focus(delegate)
				return nil
			}
			return event
		})
		delegate(d.privileged)
		return
	// detach field focus
	case execDetachFieldFocus:
		d.detach.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = execUserFieldFocus
				d.Focus(delegate)
				return nil
			}
			return event
		})
		delegate(d.detach)
		return
	// user field focus
	case execUserFieldFocus:
		d.user.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = execWorkingDirFieldFocus
				d.Focus(delegate)
				return nil
			}
			return event
		})
		delegate(d.user)
		return
	// working directory field focus
	case execWorkingDirFieldFocus:
		d.workingDir.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = execEnvVariablesFieldFocus
				d.Focus(delegate)
				return nil
			}
			return event
		})
		delegate(d.workingDir)
		return
	// env variable field focus
	case execEnvVariablesFieldFocus:
		d.envVariables.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = execEnvFileFieldFocus
				d.Focus(delegate)
				return nil
			}
			return event
		})
		delegate(d.envVariables)
		return
	// env file field focus
	case execEnvFileFieldFocus:
		d.envFile.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = execFormFieldFocus
				d.Focus(delegate)
				return nil
			}
			return event
		})
		delegate(d.envFile)
		return
	// form field focus
	case execFormFieldFocus:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = execCommandFieldFocus
				d.Focus(delegate)
				d.form.SetFocus(execCommandFieldFocus)
				return nil
			}
			if event.Key() == tcell.KeyEnter {
				d.execHandler()
				return nil
			}
			return event
		})
		delegate(d.form)
	}

}

// InputHandler returns input handler function for this primitive
func (d *ContainerExecDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("container exec dialog: event %v received", event)
		if event.Key() == tcell.KeyEsc {
			d.cancelHandler()
			return
		}
		// command field
		if d.command.HasFocus() {
			if commandHandler := d.command.InputHandler(); commandHandler != nil {
				commandHandler(event, setFocus)
				return
			}

		}
		// interactive field
		if d.interactive.HasFocus() {
			if interactiveHandler := d.interactive.InputHandler(); interactiveHandler != nil {
				interactiveHandler(event, setFocus)
				return
			}

		}
		// privileged field
		if d.privileged.HasFocus() {
			if privilegedHandler := d.privileged.InputHandler(); privilegedHandler != nil {
				privilegedHandler(event, setFocus)
				return
			}

		}
		// tty field
		if d.tty.HasFocus() {
			if ttyHandler := d.tty.InputHandler(); ttyHandler != nil {
				ttyHandler(event, setFocus)
				return
			}

		}
		// detach field
		if d.detach.HasFocus() {
			if detachHandler := d.detach.InputHandler(); detachHandler != nil {
				detachHandler(event, setFocus)
				return
			}

		}
		// working directory field
		if d.workingDir.HasFocus() {
			if workingDirHandler := d.workingDir.InputHandler(); workingDirHandler != nil {
				workingDirHandler(event, setFocus)
				return
			}

		}
		// env variables field
		if d.envVariables.HasFocus() {
			if envVariablesHandler := d.envVariables.InputHandler(); envVariablesHandler != nil {
				envVariablesHandler(event, setFocus)
				return
			}

		}
		// env file field
		if d.envFile.HasFocus() {
			if envFileHandler := d.envFile.InputHandler(); envFileHandler != nil {
				envFileHandler(event, setFocus)
				return
			}

		}
		// user field
		if d.user.HasFocus() {
			if userHandler := d.user.InputHandler(); userHandler != nil {
				userHandler(event, setFocus)
				return
			}

		}
		// form primitive
		if d.form.HasFocus() {
			if formHandler := d.form.InputHandler(); formHandler != nil {
				formHandler(event, setFocus)
				return
			}
		}
	})
}

// SetRect set rects for this primitive.
func (d *ContainerExecDialog) SetRect(x, y, width, height int) {

	dWidth := width
	dX := x
	if width > execDialogMaxWidth {
		wEmptySpace := width - execDialogMaxWidth
		if wEmptySpace > 0 {
			dX = x + (wEmptySpace / 2)
		}
		dWidth = execDialogMaxWidth
	}

	dHeight := height
	dY := y
	if height > execDialogMaxHeight {
		hEmptySpace := height - execDialogMaxHeight
		if hEmptySpace > 0 {
			dY = y + (hEmptySpace / 2)
		}
		dHeight = execDialogMaxHeight

	}
	d.Box.SetRect(dX, dY, dWidth, dHeight)
}

// Draw draws this primitive onto the screen.
func (d *ContainerExecDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}
	d.Box.DrawForSubclass(screen, d)
	x, y, width, height := d.Box.GetInnerRect()
	d.layout.SetRect(x, y, width, height)

	d.layout.Draw(screen)
}

// SetCancelFunc sets form cancel button selected function
func (d *ContainerExecDialog) SetCancelFunc(handler func()) *ContainerExecDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2)
	cancelButton.SetSelectedFunc(handler)
	return d
}

// SetExecFunc sets form execute button selected function
func (d *ContainerExecDialog) SetExecFunc(handler func()) *ContainerExecDialog {
	d.execHandler = handler
	execButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	execButton.SetSelectedFunc(handler)
	return d
}

// SetContainerID sets container ID label
func (d *ContainerExecDialog) SetContainerID(id string, name string) {
	d.containerID = id
	label := fmt.Sprintf("[white::]container:   [-::]%s", id)
	if name != "" {
		label = fmt.Sprintf("%s (%s)", label, name)
	}
	d.label.SetText(label)
}

// ContainerExecOptions returns new container exec options
func (d *ContainerExecDialog) ContainerExecOptions() containers.ExecOption {
	execOptions := containers.ExecOption{}

	cmdString := strings.TrimSpace(d.command.GetText())
	execOptions.Cmd = strings.Split(cmdString, " ")

	execOptions.Tty = d.tty.IsChecked()
	execOptions.Interactive = d.interactive.IsChecked()
	execOptions.Detach = d.detach.IsChecked()
	execOptions.Privileged = d.privileged.IsChecked()
	execOptions.WorkDir = strings.TrimSpace(d.workingDir.GetText())

	varString := strings.TrimSpace(d.envVariables.GetText())
	if varString != "" {
		execOptions.EnvVariables = strings.Split(varString, " ")
	}

	envFileString := strings.TrimSpace(d.envFile.GetText())
	if envFileString != "" {
		execOptions.EnvFile = strings.Split(envFileString, " ")
	}

	execOptions.User = strings.TrimSpace(d.user.GetText())

	return execOptions
}
