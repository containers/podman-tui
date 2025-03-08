package system

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/system/sysdialogs"
	"github.com/rivo/tview"
)

var ErrConnectionInprogres = errors.New("connection is in progress, need to disconnect")

// System implemnents the system information page primitive.
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
	appFocusHandler          func()
}

type connectionListReport struct {
	mu     sync.Mutex
	report []registry.Connection
}

type sysSelectedItem struct {
	name     string
	status   string
	uri      string
	identity string
}

// NewSystem returns new system page view.
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

	// connection table
	sys.connTable.SetBackgroundColor(style.BgColor)
	sys.connTable.SetBorder(true)
	sys.updateConnTableTitle(0)
	sys.connTable.SetTitleColor(style.FgColor)
	sys.connTable.SetBorderColor(style.BorderColor)
	sys.connTable.SetFixed(1, 1)
	sys.connTable.SetSelectable(true, false)

	for i := range sys.connTableHeaders {
		header := fmt.Sprintf("[::b]%s", strings.ToUpper(sys.connTableHeaders[i])) //nolint:perfsprint
		sys.connTable.SetCell(0, i,
			tview.NewTableCell(header).
				SetExpansion(1).
				SetBackgroundColor(style.PageHeaderBgColor).
				SetTextColor(style.PageHeaderFgColor).
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

	// set command dialog functions.
	sys.cmdDialog.SetSelectedFunc(func() {
		sys.cmdDialog.Hide()
		sys.runCommand(sys.cmdDialog.GetSelectedItem())
	})

	sys.cmdDialog.SetCancelFunc(func() {
		sys.cmdDialog.Hide()
	})

	// set confirm dialogs functions.
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

// SetAppFocusHandler sets application focus handler.
func (sys *System) SetAppFocusHandler(handler func()) {
	sys.appFocusHandler = handler
}

// GetTitle returns primitive title.
func (sys *System) GetTitle() string {
	return sys.title
}

// HasFocus returns whether or not this primitive has focus.
func (sys *System) HasFocus() bool { //nolint:cyclop
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
// displayed on the front screen.
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

// Focus is called when this primitive receives focus.
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

// SetEventMessage appends podman events to textview.
func (sys *System) SetEventMessage(messages []string) {
	msg := strings.Join(messages, "\n")
	sys.eventDialog.SetText(msg)
}

// SetConnectionProgressMessage sets connection progressbar error message.
func (sys *System) SetConnectionProgressMessage(message string) {
	sys.connPrgDialog.SetMessage(message)
}

// SetConnectionProgressDestName sets connection
// progressbar title destination name.
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

func (sys *System) getSelectedItem() *sysSelectedItem {
	selectedItem := sysSelectedItem{}

	if sys.connTable.GetRowCount() <= 1 {
		return &selectedItem
	}

	row, _ := sys.connTable.GetSelection()
	selectedItem.name = sys.connTable.GetCell(row, 0).Text
	selectedItem.status = sys.connTable.GetCell(row, 2).Text   //nolint:mnd
	selectedItem.uri = sys.connTable.GetCell(row, 3).Text      //nolint:mnd
	selectedItem.identity = sys.connTable.GetCell(row, 4).Text //nolint:mnd

	return &selectedItem
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

// SetConnectionListFunc sets list destination function.
func (sys *System) SetConnectionListFunc(list func() []registry.Connection) {
	sys.connectionListFunc = list
}

// SetConnectionSetDefaultFunc sets set destination default function.
func (sys *System) SetConnectionSetDefaultFunc(setDefault func(dest string) error) {
	sys.connectionSetDefaultFunc = setDefault
}

// SetConnectionConnectFunc sets system connect function.
func (sys *System) SetConnectionConnectFunc(connect func(dest registry.Connection)) {
	sys.connectionConnectFunc = connect
}

// SetConnectionDisconnectFunc sets system disconnect function.
func (sys *System) SetConnectionDisconnectFunc(disconnect func()) {
	sys.connectionDisconnectFunc = disconnect
}

// SetConnectionAddFunc sets system add new connection function.
func (sys *System) SetConnectionAddFunc(add func(name string, uri string, identity string) error) {
	sys.connectionAddFunc = add
}

// SetConnectionRemoveFunc sets system remove connection function.
func (sys *System) SetConnectionRemoveFunc(remove func(name string) error) {
	sys.connectionRemoveFunc = remove
}
