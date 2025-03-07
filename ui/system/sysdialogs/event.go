package sysdialogs

import (
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	formFieldHasFocus = 0 + iota
	textviewHasFocus
)

// EventsDialog implements system events view dialog primitive.
type EventsDialog struct {
	*tview.Box
	layout        *tview.Flex
	serviceName   *tview.InputField
	textview      *tview.TextView
	form          *tview.Form
	display       bool
	cancelHandler func()
	focusElement  int
}

// NewEventDialog returns new EventsDialog primitive.
func NewEventDialog() *EventsDialog {
	eventsDialog := EventsDialog{
		Box:         tview.NewBox(),
		serviceName: tview.NewInputField(),
		layout:      tview.NewFlex().SetDirection(tview.FlexRow),
	}

	// service name input field
	serviceNameLabel := "SERVICE NAME:"

	eventsDialog.serviceName.SetBackgroundColor(style.DialogBgColor)
	eventsDialog.serviceName.SetLabel("[::b]" + serviceNameLabel)
	eventsDialog.serviceName.SetLabelWidth(len(serviceNameLabel) + 1)
	eventsDialog.serviceName.SetFieldBackgroundColor(style.DialogBgColor)
	eventsDialog.serviceName.SetLabelStyle(tcell.StyleDefault.
		Background(style.DialogBorderColor).
		Foreground(style.DialogFgColor))

	// text view
	eventsDialog.textview = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)

	eventsDialog.textview.SetTextColor(style.FgColor)
	eventsDialog.textview.SetBackgroundColor(style.BgColor)
	eventsDialog.textview.SetBorderColor(style.DialogSubBoxBorderColor)
	eventsDialog.textview.SetBorder(true)

	// form
	eventsDialog.form = tview.NewForm().
		AddButton("Cancel", nil).
		SetButtonsAlign(tview.AlignRight)

	eventsDialog.form.SetBackgroundColor(style.DialogBgColor)
	eventsDialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// textview layout
	tlayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	tlayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)

	tlayout.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(eventsDialog.serviceName, 1, 0, false).
			AddItem(eventsDialog.textview, 0, 1, true),
			0, 1, true), 0, 1, true)

	tlayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)

	// layout
	eventsDialog.layout.AddItem(tlayout, 0, 1, true)
	eventsDialog.layout.AddItem(eventsDialog.form, dialogs.DialogFormHeight, 0, true)
	eventsDialog.layout.SetBorder(true)
	eventsDialog.layout.SetBackgroundColor(style.DialogBgColor)
	eventsDialog.layout.SetBorderColor(style.DialogBorderColor)

	return &eventsDialog
}

// Display displays this primitive.
func (d *EventsDialog) Display() {
	d.display = true
}

// IsDisplay returns true if primitive is shown.
func (d *EventsDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *EventsDialog) Hide() {
	d.display = false
}

// SetServiceName sets event dialog service (connection) name.
func (d *EventsDialog) SetServiceName(name string) {
	d.layout.SetTitle("SYSTEM EVENTS")
	d.serviceName.SetText(name)
}

// SetText sets message dialog text messages.
func (d *EventsDialog) SetText(message string) {
	d.textview.Clear()
	d.textview.SetText(message)
	d.textview.ScrollToEnd()
}

// Focus is called when this primitive receives focus.
func (d *EventsDialog) Focus(delegate func(p tview.Primitive)) {
	switch d.focusElement {
	// text screen field focus
	case textviewHasFocus:
		d.textview.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == utils.SwitchFocusKey.EventKey() {
				d.focusElement = formFieldHasFocus
				d.Focus(delegate)

				return nil
			}

			return event
		})
		delegate(d.textview)
	// form field focus
	case formFieldHasFocus:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == utils.SwitchFocusKey.EventKey() {
				d.focusElement = textviewHasFocus
				d.Focus(delegate)

				return nil
			}

			if event.Key() == tcell.KeyEnter {
				d.cancelHandler()

				return nil
			}

			return event
		})

		delegate(d.form)
	}
}

// HasFocus returns whether or not this primitive has focus.
func (d *EventsDialog) HasFocus() bool {
	return d.form.HasFocus() || d.textview.HasFocus()
}

// SetRect set rects for this primitive.
func (d *EventsDialog) SetRect(x, y, width, height int) {
	dX := x + 1
	dY := y + 1
	dWidth := width - 2   //nolint:mnd
	dHeight := height - 2 //nolint:mnd

	d.Box.SetRect(dX, dY, dWidth, dHeight)
}

// Draw draws this primitive onto the screen.
func (d *EventsDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.Box.DrawForSubclass(screen, d)
	x, y, width, height := d.Box.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// InputHandler returns input handler function for this primitive.
func (d *EventsDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("events dialog: event %v received", event)

		if event.Key() == utils.CloseDialogKey.EventKey() {
			d.cancelHandler()

			return
		}

		// textview field
		if d.textview.HasFocus() {
			if textviewHandler := d.textview.InputHandler(); textviewHandler != nil {
				textviewHandler(event, setFocus)

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

// SetCancelFunc sets form cancel button selected function.
func (d *EventsDialog) SetCancelFunc(handler func()) *EventsDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	cancelButton.SetSelectedFunc(handler)

	return d
}
