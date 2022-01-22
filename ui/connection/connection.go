package connection

import (
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
)

const (
	maxWidth  = 60
	maxHeight = 10
)

// Connection implements the Connection page primitive
type Connection struct {
	*tview.Box
	title          string
	layout         *tview.Flex
	textview       *tview.TextView
	progressDialog *tvxwidgets.ActivityModeGauge
}

// NewConnection returns containers page view
func NewConnection() *Connection {
	conn := &Connection{
		Box:            tview.NewBox(),
		title:          "connection",
		layout:         tview.NewFlex().SetDirection(tview.FlexRow),
		textview:       tview.NewTextView(),
		progressDialog: tvxwidgets.NewActivityModeGauge(),
	}

	conn.progressDialog.SetPgBgColor(tcell.ColorOrange)

	fgColor := utils.Styles.PageTable.FgColor
	bgColor := utils.Styles.PageTable.BgColor

	conn.Box.SetBorderColor(bgColor)
	conn.Box.SetTitleColor(fgColor)
	conn.Box.SetBorder(true)
	conn.layout.SetBorder(true)
	conn.layout.SetTitle("connecting to podman")
	conn.layout.SetBackgroundColor(tcell.ColorOrangeRed)
	conn.progressDialog.SetBackgroundColor(tcell.ColorOrangeRed)
	conn.textview.SetBackgroundColor(tcell.ColorOrangeRed)

	conn.layout.AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorOrangeRed), 1, 1, false)
	conn.layout.AddItem(conn.progressDialog, 1, 1, false)
	conn.layout.AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorOrangeRed), 1, 1, false)
	conn.layout.AddItem(conn.textview, 0, 0, false)

	return conn
}

// GetTitle returns primitive title
func (conn *Connection) GetTitle() string {
	return conn.title
}

// SetErrorMessage sets connection page error message
func (conn *Connection) SetErrorMessage(message string) {
	if message == "" {
		conn.layout.ResizeItem(conn.textview, 0, 0)
	} else {
		conn.layout.ResizeItem(conn.textview, 0, 1)
	}
	conn.textview.SetText(message)
}

// HasFocus returns whether or not this primitive has focus
func (conn *Connection) HasFocus() bool {
	return conn.Box.HasFocus()
}

// Focus is called when this primitive receives focus
func (conn *Connection) Focus(delegate func(p tview.Primitive)) {
	delegate(conn.Box)
}

// Reset resets progress bar and text view
func (conn *Connection) Reset() {
	conn.progressDialog.Reset()
	conn.textview.SetText("")
}
