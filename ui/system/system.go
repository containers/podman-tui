package system

import (
	"fmt"
	"strings"
	"sync"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/system/sysdialogs"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/rivo/tview"
)

// System implemnents the system information page primitive
type System struct {
	*tview.Box
	title                    string
	connTable                *tview.Table
	connTableHeaders         []string
	connectionList           connectionListReport
	cmdDialog                *dialogs.CommandDialog
	confirmDialog            *dialogs.ConfirmDialog
	messageDialog            *dialogs.MessageDialog
	progressDialog           *dialogs.ProgressDialog
	errorDialog              *dialogs.ErrorDialog
	eventDialog              *sysdialogs.EventsDialog
	dfDialog                 *sysdialogs.DfDialog
	connPrgDialog            *sysdialogs.ConnectDialog
	connAddDialog            *sysdialogs.AddConnectionDialog
	confirmData              string
	connectionListFunc       func() []registry.Connection
	connectionAddFunc        func(string, string, string) error
	connectionRemoveFunc     func(string) error
	connectionSetDefaultFunc func(string) error
	connectionConnectFunc    func(registry.Connection)
	connectionDisconnectFunc func()
}

type connectionListReport struct {
	mu     sync.Mutex
	report []registry.Connection
}

// NewSystem returns new system page view
func NewSystem() *System {
	sys := &System{
		Box:              tview.NewBox(),
		title:            "system",
		connTable:        tview.NewTable(),
		connTableHeaders: []string{"name", "default", "status", "uri", "identity"},
		confirmDialog:    dialogs.NewConfirmDialog(),
		progressDialog:   dialogs.NewProgressDialog(),
		errorDialog:      dialogs.NewErrorDialog(),
		messageDialog:    dialogs.NewMessageDialog(""),
		eventDialog:      sysdialogs.NewEventDialog(),
		dfDialog:         sysdialogs.NewDfDialog(),
		connPrgDialog:    sysdialogs.NewConnectDialog(),
		connAddDialog:    sysdialogs.NewAddConnectionDialog(),
	}
	fgColor := utils.Styles.PageTable.FgColor
	bgColor := utils.Styles.PageTable.BgColor

	// connection table
	sys.connTable.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	sys.connTable.SetBorder(true)
	sys.updateConnTableTitle(0)
	sys.connTable.SetTitleColor(fgColor)
	sys.connTable.SetBorderColor(bgColor)
	sys.connTable.SetFixed(1, 1)
	sys.connTable.SetSelectable(true, false)

	fgColor = utils.Styles.PageTable.HeaderRow.FgColor
	bgColor = utils.Styles.PageTable.HeaderRow.BgColor
	for i := 0; i < len(sys.connTableHeaders); i++ {
		header := fmt.Sprintf("[::b]%s", strings.ToUpper(sys.connTableHeaders[i]))
		sys.connTable.SetCell(0, i,
			tview.NewTableCell(header).
				SetExpansion(1).
				SetBackgroundColor(bgColor).
				SetTextColor(fgColor).
				SetAlign(tview.AlignLeft).
				SetSelectable(false))
	}

	sys.cmdDialog = dialogs.NewCommandDialog([][]string{
		{"add connection", "record destination for the Podman TUI service"},
		{"connect", "connect to selected destination"},
		{"disconnect", "disconnect from connected destination"},
		{"disk usage", "display destination podman related disk usage"},
		{"events", "display destination system events"},
		{"info", "display destination podman system information"},
		{"prune", "remove all unused pod, container, image and volume data"},
		{"remove connection", "delete named destination for the Podman TUI"},
		{"set default", "set selected destination as a default service"},
	})

	// set command dialog functions
	sys.cmdDialog.SetSelectedFunc(func() {
		sys.cmdDialog.Hide()
		sys.runCommand(sys.cmdDialog.GetSelectedItem())
	})
	sys.cmdDialog.SetCancelFunc(func() {
		sys.cmdDialog.Hide()
	})
	// set confirm dialogs functions
	sys.confirmDialog.SetSelectedFunc(func() {
		sys.confirmDialog.Hide()
		switch sys.confirmData {
		case "prune":
			sys.prune()
		case "remove_conn":
			sys.remove()
		}
	})
	sys.confirmDialog.SetCancelFunc(func() {
		sys.confirmDialog.Hide()
	})

	// set message dialog functions
	sys.messageDialog.SetCancelFunc(func() {
		sys.messageDialog.Hide()
	})

	// set event dialog functions
	sys.eventDialog.SetCancelFunc(func() {
		sys.eventDialog.Hide()
	})

	// set disk usage function
	sys.dfDialog.SetCancelFunc(func() {
		sys.dfDialog.Hide()
	})

	// set connection progress bar cancel function
	sys.connPrgDialog.SetCancelFunc(func() {
		sys.connPrgDialog.Hide()
		registry.UnsetConnection()
		sys.eventDialog.SetText("")
		sys.UpdateConnectionsData()
	})
	// set connection create dialog functions
	sys.connAddDialog.SetCancelFunc(sys.connAddDialog.Hide)
	sys.connAddDialog.SetAddFunc(func() {
		sys.addConnection()
	})
	return sys
}

// GetTitle returns primitive title
func (sys *System) GetTitle() string {
	return sys.title
}

// HasFocus returns whether or not this primitive has focus
func (sys *System) HasFocus() bool {
	if sys.cmdDialog.HasFocus() || sys.confirmDialog.HasFocus() {
		return true
	}

	if sys.progressDialog.HasFocus() || sys.errorDialog.HasFocus() {
		return true
	}
	if sys.eventDialog.HasFocus() || sys.dfDialog.HasFocus() {
		return true
	}
	if sys.messageDialog.HasFocus() || sys.connTable.HasFocus() {
		return true
	}
	if sys.connPrgDialog.HasFocus() || sys.connAddDialog.HasFocus() {
		return true
	}
	return sys.Box.HasFocus()
}

// SubDialogHasFocus returns true if there is an active dialog
// displayed on the front screen
func (sys *System) SubDialogHasFocus() bool {
	if sys.cmdDialog.HasFocus() || sys.confirmDialog.HasFocus() {
		return true
	}
	if sys.progressDialog.HasFocus() || sys.errorDialog.HasFocus() {
		return true
	}
	if sys.dfDialog.HasFocus() || sys.messageDialog.HasFocus() {
		return true
	}
	if sys.eventDialog.HasFocus() || sys.connAddDialog.HasFocus() {
		return true
	}
	return false
}

// Focus is called when this primitive receives focus
func (sys *System) Focus(delegate func(p tview.Primitive)) {
	// error dialog
	if sys.errorDialog.IsDisplay() {
		delegate(sys.errorDialog)
		return
	}
	// message dialog
	if sys.messageDialog.IsDisplay() {
		delegate(sys.messageDialog)
		return
	}
	// command dialog
	if sys.cmdDialog.IsDisplay() {
		delegate(sys.cmdDialog)
		return
	}
	// confirm dialog
	if sys.confirmDialog.IsDisplay() {
		delegate(sys.confirmDialog)
		return
	}
	// disk usage dialog
	if sys.dfDialog.IsDisplay() {
		delegate(sys.dfDialog)
		return
	}
	// connection progress dialog
	if sys.connPrgDialog.IsDisplay() {
		delegate(sys.connPrgDialog)
		return
	}
	// event dialog
	if sys.eventDialog.IsDisplay() {
		delegate(sys.eventDialog)
		return
	}
	// connection create dialog
	if sys.connAddDialog.IsDisplay() {
		delegate(sys.connAddDialog)
		return
	}
	delegate(sys.connTable)
}

// SetEventMessage appends podman events to textview
func (sys *System) SetEventMessage(messages []string) {
	msg := strings.Join(messages, "\n")
	sys.eventDialog.SetText(msg)
}

// SetConnectionProgressMessage sets connection progressbar error message
func (sys *System) SetConnectionProgressMessage(message string) {
	sys.connPrgDialog.SetMessage(message)
}

// SetConnectionProgressDestName sets connection
// progressbar title destination name
func (sys *System) SetConnectionProgressDestName(name string) {
	sys.connPrgDialog.SetDestinationName(name)
}

// ConnectionProgressDisplay displays or hide the connection progress dialog.
func (sys *System) ConnectionProgressDisplay(display bool) {
	if display {
		sys.hideAllDialogs()
		sys.connPrgDialog.Display()
		return
	}
	sys.connPrgDialog.Hide()
}

func (sys *System) getSelectedItem() (string, string, string, string) {
	var (
		name     string
		uri      string
		status   string
		identity string
	)
	if sys.connTable.GetRowCount() <= 1 {
		return name, status, uri, identity
	}
	row, _ := sys.connTable.GetSelection()
	name = sys.connTable.GetCell(row, 0).Text
	status = sys.connTable.GetCell(row, 2).Text
	uri = sys.connTable.GetCell(row, 3).Text
	identity = sys.connTable.GetCell(row, 4).Text
	return name, status, uri, identity
}

func (sys *System) hideAllDialogs() {
	sys.errorDialog.Hide()
	sys.cmdDialog.Hide()
	sys.confirmDialog.Hide()
	sys.messageDialog.Hide()
	sys.dfDialog.Hide()
	sys.progressDialog.Hide()
	sys.eventDialog.Hide()
	sys.connAddDialog.Hide()
}

// SetConnectionListFunc sets list destination function
func (sys *System) SetConnectionListFunc(list func() []registry.Connection) {
	sys.connectionListFunc = list
}

// SetConnectionSetDefaultFunc sets set destination default function
func (sys *System) SetConnectionSetDefaultFunc(setDefault func(dest string) error) {
	sys.connectionSetDefaultFunc = setDefault
}

// SetConnectionConnectFunc sets system connect function
func (sys *System) SetConnectionConnectFunc(connect func(dest registry.Connection)) {
	sys.connectionConnectFunc = connect
}

// SetConnectionDisconnectFunc sets system disconnect function
func (sys *System) SetConnectionDisconnectFunc(disconnect func()) {
	sys.connectionDisconnectFunc = disconnect
}

// SetConnectionAddFunc sets system add new connection function
func (sys *System) SetConnectionAddFunc(add func(name string, uri string, identity string) error) {
	sys.connectionAddFunc = add
}

// SetConnectionRemoveFunc sets system remove connection function
func (sys *System) SetConnectionRemoveFunc(remove func(name string) error) {
	sys.connectionRemoveFunc = remove
}
