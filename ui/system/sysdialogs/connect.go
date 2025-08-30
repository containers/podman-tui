package sysdialogs

import (
	"fmt"

	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	maxWidth  = 60
	maxHeight = 12
)

// ConnectDialog implements the Connection progress dialog primitive.
type ConnectDialog struct {
	*tview.Box

	layout         *tview.Flex
	textview       *tview.TextView
	progressDialog *tvxwidgets.ActivityModeGauge
	display        bool
	cancelButton   *tview.Button
}

// NewConnectDialog returns connection progress dialog.
func NewConnectDialog() *ConnectDialog {
	conn := &ConnectDialog{
		Box:            tview.NewBox(),
		layout:         tview.NewFlex().SetDirection(tview.FlexRow),
		textview:       tview.NewTextView(),
		progressDialog: tvxwidgets.NewActivityModeGauge(),
		cancelButton:   tview.NewButton(" Cancel "),
	}

	// connect dialog box
	conn.SetBorder(false)

	// progress bar
	conn.progressDialog.SetBackgroundColor(style.DialogBgColor)
	conn.progressDialog.SetPgBgColor(style.PrgBarColor)

	// connection message text view
	conn.textview.SetBackgroundColor(style.DialogBgColor)
	conn.textview.SetTextColor(style.DialogFgColor)

	// cancel button and layout
	conn.cancelButton.SetStyle(tcell.StyleDefault.Background(style.ButtonBgColor))
	conn.cancelButton.SetLabelColorActivated(style.ButtonBgColor)

	cancelLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	cancelLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 0, 1, false)
	cancelLayout.AddItem(conn.cancelButton, 10, 0, true) //nolint:mnd
	cancelLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)
	cancelLayout.SetBackgroundColor(style.DialogBgColor)

	// connection progress layout
	conn.layout.SetBorder(true)
	conn.layout.SetBorderColor(style.DialogBorderColor)
	conn.layout.SetTitleColor(style.DialogFgColor)

	conn.layout.SetBackgroundColor(style.DialogBgColor)
	conn.layout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)
	conn.layout.AddItem(conn.progressDialog, 1, 0, false)
	conn.layout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)
	conn.layout.AddItem(conn.textview, 0, 0, false)
	conn.layout.AddItem(cancelLayout, 1, 0, false)
	conn.layout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)

	conn.display = false

	return conn
}

// Display displays this primitive.
func (d *ConnectDialog) Display() {
	d.display = true
}

// IsDisplay returns true if primitive is shown.
func (d *ConnectDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *ConnectDialog) Hide() {
	d.display = false

	d.SetDestinationName("")
	d.reset()
}

// SetMessage sets connection page error message.
func (d *ConnectDialog) SetMessage(message string) {
	if message == "" {
		d.layout.ResizeItem(d.textview, 0, 0)
	} else {
		d.layout.ResizeItem(d.textview, 0, 1)
	}

	d.textview.SetText(message)
}

// HasFocus returns whether or not this primitive has focus.
func (d *ConnectDialog) HasFocus() bool {
	return d.Box.HasFocus() || d.cancelButton.HasFocus()
}

// Focus is called when this primitive receives focus.
func (d *ConnectDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(d.cancelButton)
}

// SetRect sets a new position of the primitive.
func (d *ConnectDialog) SetRect(x int, y int, width int, height int) {
	emptyWidth := (width - maxWidth) / 2    //nolint:mnd
	emptyHeight := (height - maxHeight) / 2 //nolint:mnd

	if width > maxWidth {
		width = maxWidth
		x += emptyWidth
	}

	if height > maxHeight {
		height = maxHeight
		y += emptyHeight
	}

	if d.textview.GetText(true) == "" {
		height -= 5
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen.
func (d *ConnectDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.DrawForSubclass(screen, d)
	x, y, width, height := d.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.progressDialog.Pulse()
	d.layout.Draw(screen)
}

// InputHandler returns input handler function for this primitive.
func (d *ConnectDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("connection progress dialog: event %v received", event)

		if cancelButtonHandler := d.cancelButton.InputHandler(); cancelButtonHandler != nil {
			cancelButtonHandler(event, setFocus)

			return
		}
	})
}

// SetDestinationName sets progress bar title destination name.
func (d *ConnectDialog) SetDestinationName(name string) {
	title := fmt.Sprintf("connecting to %s", name) //nolint:perfsprint

	d.layout.SetTitle(title)
}

// SetCancelFunc sets progress bar cancel button function.
func (d *ConnectDialog) SetCancelFunc(cancel func()) {
	d.cancelButton.SetSelectedFunc(cancel)
}

func (d *ConnectDialog) reset() {
	d.progressDialog.Reset()
	d.textview.SetText("")
}
