package dialogs

import (
	"strings"

	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// MessageDialog is a simaple message dialog primitive.
type MessageDialog struct {
	*tview.Box
	layout        *tview.Flex
	infoType      *tview.InputField
	textview      *tview.TextView
	form          *tview.Form
	display       bool
	message       string
	cancelHandler func()
}

type messageInfo int

const (
	// top dialog header label.
	MessageSystemInfo messageInfo = 0 + iota
	MessagePodInfo
	MessageContainerInfo
	MessageVolumeInfo
	MessageImageInfo
	MessageNetworkInfo
)

// NewMessageDialog returns new message dialog primitive.
func NewMessageDialog(text string) *MessageDialog {
	dialog := &MessageDialog{
		Box:      tview.NewBox(),
		infoType: tview.NewInputField(),
		display:  false,
		message:  text,
	}

	dialog.infoType.SetBackgroundColor(style.DialogBgColor)
	dialog.infoType.SetFieldBackgroundColor(style.DialogBgColor)
	dialog.infoType.SetLabelStyle(tcell.StyleDefault.
		Background(style.DialogBorderColor).
		Foreground(style.DialogFgColor))

	dialog.textview = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)

	dialog.textview.SetTextColor(style.FgColor)
	dialog.textview.SetBackgroundColor(style.BgColor)
	dialog.textview.SetBorder(true)
	dialog.textview.SetBorderColor(style.DialogSubBoxBorderColor)

	// textview layout
	tlayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	tlayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)
	tlayout.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(dialog.infoType, 1, 0, false).
		AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false).
		AddItem(dialog.textview, 0, 1, true),
		0, 1, true)
	tlayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)

	dialog.form = tview.NewForm().
		AddButton("Cancel", nil).
		SetButtonsAlign(tview.AlignRight)

	dialog.form.SetBackgroundColor(style.DialogBgColor)
	dialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	dialog.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	dialog.layout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	dialog.layout.AddItem(tlayout, 0, 1, true)
	dialog.layout.AddItem(dialog.form, DialogFormHeight, 0, true)
	dialog.layout.SetBorder(true)
	dialog.layout.SetBorderColor(style.DialogBorderColor)
	dialog.layout.SetBackgroundColor(style.DialogBgColor)

	return dialog
}

// Display displays this primitive.
func (d *MessageDialog) Display() {
	d.display = true
}

// IsDisplay returns true if primitive is shown.
func (d *MessageDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *MessageDialog) Hide() {
	d.message = ""
	d.textview.SetText("")
	d.display = false
}

// SetTitle sets input dialog title.
func (d *MessageDialog) SetTitle(title string) {
	d.layout.SetTitle(strings.ToUpper(title))
}

// SetText sets message dialog text messages.
func (d *MessageDialog) SetText(headerType messageInfo, headerMessage string, message string) {
	msgTypeLabel := ""

	switch headerType {
	case MessageSystemInfo:
		msgTypeLabel = "SERVICE NAME:"
	case MessagePodInfo:
		msgTypeLabel = "POD ID:"
	case MessageContainerInfo:
		msgTypeLabel = "CONTAINER ID:"
	case MessageVolumeInfo:
		msgTypeLabel = "VOLUME NAME:"
	case MessageImageInfo:
		msgTypeLabel = "IMAGE ID:"
	case MessageNetworkInfo:
		msgTypeLabel = "NETWORK ID:"
	}

	if msgTypeLabel != "" {
		d.infoType.SetLabel("[::b]" + msgTypeLabel)
		d.infoType.SetLabelWidth(len(msgTypeLabel) + 1)
		d.infoType.SetText(headerMessage)
	}

	d.message = strings.TrimSpace(message)
	d.textview.Clear()

	if d.message == "" {
		d.textview.SetBorder(false)
		d.textview.SetText("")
	} else {
		d.textview.SetBorder(true)
		d.textview.SetBorderColor(style.DialogSubBoxBorderColor)
		d.textview.SetText(message)
	}

	d.textview.ScrollToBeginning()
}

// TextScrollToEnd scroll downs the text view.
func (d *MessageDialog) TextScrollToEnd() {
	d.textview.ScrollToEnd()
}

// Focus is called when this primitive receives focus.
func (d *MessageDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(d.form)
}

// HasFocus returns whether or not this primitive has focus.
func (d *MessageDialog) HasFocus() bool {
	return d.form.HasFocus()
}

// SetRect set rects for this primitive.
func (d *MessageDialog) SetRect(x, y, width, height int) {
	messageHeight := 0
	if d.message != "" {
		messageHeight = len(strings.Split(d.message, "\n")) + 3 //nolint:gomnd
	}

	messageWidth := getMessageWidth(d.message)

	headerWidth := len(d.infoType.GetText()) + len(d.infoType.GetLabel()) + 4 //nolint:gomnd
	if messageWidth < headerWidth {
		messageWidth = headerWidth
	}

	dWidth := width - (2 * DialogPadding) //nolint:gomnd
	if messageWidth+4 < dWidth {
		dWidth = messageWidth + 4 //nolint:gomnd
	}

	if DialogMinWidth < width && dWidth < DialogMinWidth {
		dWidth = DialogMinWidth
	}

	emptySpace := (width - dWidth) / 2 //nolint:gomnd
	dX := x + emptySpace

	dHeight := messageHeight + DialogFormHeight + DialogPadding + 1
	if dHeight > height {
		dHeight = height - DialogPadding - 1
	}

	textviewHeight := dHeight - DialogFormHeight - 2 //nolint:gomnd
	hs := ((height - dHeight) / 2)                   //nolint:gomnd
	dY := y + hs

	d.Box.SetRect(dX, dY, dWidth, dHeight)

	d.layout.ResizeItem(d.textview, textviewHeight, 0)
}

// Draw draws this primitive onto the screen.
func (d *MessageDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.Box.DrawForSubclass(screen, d)
	x, y, width, height := d.Box.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// InputHandler returns input handler function for this primitive.
func (d *MessageDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("message dialog: event %v received", event)

		if event.Key() == tcell.KeyEsc {
			d.cancelHandler()

			return
		}

		if event.Key() == tcell.KeyEnter {
			d.cancelHandler()

			return
		}

		if event.Key() == tcell.KeyTab {
			if formHandler := d.form.InputHandler(); formHandler != nil {
				formHandler(event, setFocus)

				return
			}
		}

		// scroll between message textview
		if textHandler := d.textview.InputHandler(); textHandler != nil {
			textHandler(event, setFocus)

			return
		}
	})
}

// SetCancelFunc sets form cancel button selected function.
func (d *MessageDialog) SetCancelFunc(handler func()) *MessageDialog {
	d.cancelHandler = handler

	return d
}
