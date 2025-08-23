package networks

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/containers/common/libnetwork/types"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/networks/netdialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/rivo/tview"
)

const (
	viewNetworkNameColIndex = 0 + iota
	viewNetworkVersionColIndex
	viewNetworkPluginColIndex
)

var (
	errNoNetworkRemove     = errors.New("there is no network to remove")
	errNoNetworkInspect    = errors.New("there is no network to display inspect")
	errNoNetworkDisconnect = errors.New("there is no network to disconnect")
	errNoNetworkConnect    = errors.New("there is no network to connect")
)

// Networks implemnents the Networks page primitive.
type Networks struct {
	*tview.Box

	title            string
	headers          []string
	table            *tview.Table
	errorDialog      *dialogs.ErrorDialog
	progressDialog   *dialogs.ProgressDialog
	confirmDialog    *dialogs.ConfirmDialog
	cmdDialog        *dialogs.CommandDialog
	messageDialog    *dialogs.MessageDialog
	sortDialog       *dialogs.SortDialog
	createDialog     *netdialogs.NetworkCreateDialog
	connectDialog    *netdialogs.NetworkConnectDialog
	disconnectDialog *netdialogs.NetworkDisconnectDialog
	networkList      networkListReport
	selectedID       string
	confirmData      string
	appFocusHandler  func()
}

type networkListReport struct {
	mu        sync.Mutex
	report    []types.Network
	sortBy    string
	ascending bool
}

// NewNetworks returns nets page view.
func NewNetworks() *Networks {
	nets := &Networks{
		Box:              tview.NewBox(),
		title:            "networks",
		headers:          []string{"id", "name", "driver"},
		errorDialog:      dialogs.NewErrorDialog(),
		progressDialog:   dialogs.NewProgressDialog(),
		confirmDialog:    dialogs.NewConfirmDialog(),
		messageDialog:    dialogs.NewMessageDialog(""),
		sortDialog:       dialogs.NewSortDialog([]string{"name", "driver"}, 0),
		createDialog:     netdialogs.NewNetworkCreateDialog(),
		connectDialog:    netdialogs.NewNetworkConnectDialog(),
		disconnectDialog: netdialogs.NewNetworkDisconnectDialog(),
		networkList:      networkListReport{sortBy: "name", ascending: true},
	}

	nets.cmdDialog = dialogs.NewCommandDialog([][]string{
		{"connect", "connect a container to a network"},
		{"create", "create a Podman CNI network"},
		{"disconnect", "disconnect a container from a network"},
		{"inspect", "displays the raw CNI network configuration"},
		{"prune", "remove all unused networks"},
		// {"reload", "reload the network for containers"},
		{"rm", "remove a CNI networks"},
	})

	nets.table = tview.NewTable()
	nets.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(nets.title)))
	nets.table.SetBorderColor(style.BorderColor)
	nets.table.SetBackgroundColor(style.BgColor)
	nets.table.SetTitleColor(style.FgColor)
	nets.table.SetBorder(true)

	for i := range nets.headers {
		nets.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[::b]%s", strings.ToUpper(nets.headers[i]))). //nolint:perfsprint
													SetExpansion(1).
													SetBackgroundColor(style.PageHeaderBgColor).
													SetTextColor(style.PageHeaderFgColor).
													SetAlign(tview.AlignLeft).
													SetSelectable(false))
	}

	nets.table.SetFixed(1, 1)
	nets.table.SetSelectable(true, false)

	// set command dialog functions
	nets.cmdDialog.SetSelectedFunc(func() {
		nets.cmdDialog.Hide()
		nets.runCommand(nets.cmdDialog.GetSelectedItem())
	})
	nets.cmdDialog.SetCancelFunc(func() {
		nets.cmdDialog.Hide()
	})

	// set message dialog functions
	nets.messageDialog.SetCancelFunc(func() {
		nets.messageDialog.Hide()
	})

	// set confirm dialogs functions
	nets.confirmDialog.SetSelectedFunc(func() {
		nets.confirmDialog.Hide()

		switch nets.confirmData {
		case utils.PruneCommandLabel:
			nets.prune()
		case "rm":
			nets.remove()
		}
	})

	nets.confirmDialog.SetCancelFunc(func() {
		nets.confirmDialog.Hide()
	})

	// set create dialog functions
	nets.createDialog.SetCancelFunc(func() {
		nets.createDialog.Hide()
	})

	nets.createDialog.SetCreateFunc(func() {
		nets.createDialog.Hide()
		nets.create()
	})

	// set connect dialog functions
	nets.connectDialog.SetCancelFunc(nets.connectDialog.Hide)
	nets.connectDialog.SetConnectFunc(nets.connect)

	// set disconnect dialog functions
	nets.disconnectDialog.SetCancelFunc(nets.disconnectDialog.Hide)
	nets.disconnectDialog.SetDisconnectFunc(nets.disconnect)

	// set sort dialog functions
	nets.sortDialog.SetCancelFunc(nets.sortDialog.Hide)
	nets.sortDialog.SetSelectFunc(nets.SortView)

	return nets
}

// SetAppFocusHandler sets application focus handler.
func (nets *Networks) SetAppFocusHandler(handler func()) {
	nets.appFocusHandler = handler
}

// GetTitle returns primitive title.
func (nets *Networks) GetTitle() string {
	return nets.title
}

// HasFocus returns whether or not this primitive has focus.
func (nets *Networks) HasFocus() bool {
	if nets.SubDialogHasFocus() {
		return true
	}

	if nets.table.HasFocus() || nets.Box.HasFocus() {
		return true
	}

	return false
}

// SubDialogHasFocus returns whether or not sub dialog primitive has focus.
func (nets *Networks) SubDialogHasFocus() bool {
	for _, dialog := range nets.getInnerDialogs() {
		if dialog.HasFocus() {
			return true
		}
	}

	return false
}

// Focus is called when this primitive receives focus.
func (nets *Networks) Focus(delegate func(p tview.Primitive)) {
	if nets.errorDialog.IsDisplay() {
		delegate(nets.errorDialog)

		return
	}

	if nets.confirmDialog.IsDisplay() {
		delegate(nets.confirmDialog)

		return
	}

	for _, dialog := range nets.getInnerDialogs() {
		if dialog.IsDisplay() {
			delegate(dialog)

			return
		}
	}

	delegate(nets.table)
}

// HideAllDialogs hides all sub dialogs.
func (nets *Networks) HideAllDialogs() {
	for _, dialog := range nets.getInnerDialogs() {
		if dialog.IsDisplay() {
			dialog.Hide()
		}
	}
}

func (nets *Networks) getSelectedItem() (string, string) {
	if nets.table.GetRowCount() <= 1 {
		return "", ""
	}

	row, _ := nets.table.GetSelection()
	netID := nets.table.GetCell(row, 0).Text
	netName := nets.table.GetCell(row, 1).Text

	return netID, netName
}

func (nets *Networks) getInnerDialogs() []utils.UIDialog {
	dialogs := []utils.UIDialog{
		nets.errorDialog,
		nets.progressDialog,
		nets.confirmDialog,
		nets.cmdDialog,
		nets.messageDialog,
		nets.connectDialog,
		nets.createDialog,
		nets.disconnectDialog,
		nets.sortDialog,
	}

	return dialogs
}
