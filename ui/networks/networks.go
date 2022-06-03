package networks

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/networks/netdialogs"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/rivo/tview"
)

// Networks implemnents the Networks page primitive
type Networks struct {
	*tview.Box
	title          string
	headers        []string
	table          *tview.Table
	errorDialog    *dialogs.ErrorDialog
	progressDialog *dialogs.ProgressDialog
	confirmDialog  *dialogs.ConfirmDialog
	cmdDialog      *dialogs.CommandDialog
	messageDialog  *dialogs.MessageDialog
	createDialog   *netdialogs.NetworkCreateDialog
	selectedID     string
	confirmData    string
}

// NewNetworks returns nets page view
func NewNetworks() *Networks {
	nets := &Networks{
		Box:            tview.NewBox(),
		title:          "networks",
		headers:        []string{"id", "name", "driver"},
		errorDialog:    dialogs.NewErrorDialog(),
		progressDialog: dialogs.NewProgressDialog(),
		confirmDialog:  dialogs.NewConfirmDialog(),
		messageDialog:  dialogs.NewMessageDialog(""),
		createDialog:   netdialogs.NewNetworkCreateDialog(),
	}

	nets.cmdDialog = dialogs.NewCommandDialog([][]string{
		//{"connect", "connect a container to a network"},
		{"create", "create a Podman CNI network"},
		//{"disconnect", "disconnect a container from a network"},
		{"inspect", "displays the raw CNI network configuration"},
		{"prune", "remove all unused networks"},
		//{"reload", "reload the network for containers"},
		{"rm", "remove a CNI networks"},
	})
	fgColor := utils.Styles.PageTable.FgColor
	bgColor := utils.Styles.PageTable.BgColor
	nets.table = tview.NewTable()
	nets.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(nets.title)))
	nets.table.SetBorderColor(bgColor)
	nets.table.SetTitleColor(fgColor)
	nets.table.SetBorder(true)

	fgColor = utils.Styles.PageTable.HeaderRow.FgColor
	bgColor = utils.Styles.PageTable.HeaderRow.BgColor

	for i := 0; i < len(nets.headers); i++ {
		nets.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[::b]%s", strings.ToUpper(nets.headers[i]))).
				SetExpansion(1).
				SetBackgroundColor(bgColor).
				SetTextColor(fgColor).
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
		case "prune":
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
	return nets
}

// GetTitle returns primitive title
func (nets *Networks) GetTitle() string {
	return nets.title
}

// HasFocus returns whether or not this primitive has focus
func (nets *Networks) HasFocus() bool {
	if nets.table.HasFocus() || nets.errorDialog.HasFocus() {
		return true
	}
	if nets.cmdDialog.HasFocus() || nets.messageDialog.IsDisplay() {
		return true
	}
	if nets.progressDialog.HasFocus() || nets.confirmDialog.IsDisplay() {
		return true
	}
	if nets.createDialog.HasFocus() {
		return true
	}
	return nets.Box.HasFocus()
}

// SubDialogHasFocus returns whether or not sub dialog primitive has focus
func (nets *Networks) SubDialogHasFocus() bool {
	if nets.createDialog.HasFocus() || nets.errorDialog.HasFocus() {
		return true
	}
	if nets.cmdDialog.HasFocus() || nets.messageDialog.IsDisplay() {
		return true
	}
	if nets.progressDialog.HasFocus() || nets.confirmDialog.IsDisplay() {
		return true
	}
	return false
}

// Focus is called when this primitive receives focus
func (nets *Networks) Focus(delegate func(p tview.Primitive)) {
	// error dialog
	if nets.errorDialog.IsDisplay() {
		delegate(nets.errorDialog)
		return
	}
	// command dialog
	if nets.cmdDialog.IsDisplay() {
		delegate(nets.cmdDialog)
		return
	}
	// message dialog
	if nets.messageDialog.IsDisplay() {
		delegate(nets.messageDialog)
		return
	}
	// confirm dialog
	if nets.confirmDialog.IsDisplay() {
		delegate(nets.confirmDialog)
		return
	}
	// create dialog
	if nets.createDialog.IsDisplay() {
		delegate(nets.createDialog)
		return
	}
	delegate(nets.table)
}

func (nets *Networks) getSelectedItem() string {
	if nets.table.GetRowCount() <= 1 {
		return ""
	}
	row, _ := nets.table.GetSelection()
	netID := nets.table.GetCell(row, 0).Text
	netName := nets.table.GetCell(row, 1).Text

	netIDorName := netID
	if netIDorName == "" {
		netIDorName = netName
	}
	return netIDorName
}

// HideAllDialogs hides all sub dialogs
func (nets *Networks) HideAllDialogs() {
	if nets.errorDialog.IsDisplay() {
		nets.errorDialog.Hide()
	}
	if nets.progressDialog.IsDisplay() {
		nets.progressDialog.Hide()
	}
	if nets.confirmDialog.IsDisplay() {
		nets.confirmDialog.Hide()
	}
	if nets.cmdDialog.IsDisplay() {
		nets.cmdDialog.Hide()
	}
	if nets.messageDialog.IsDisplay() {
		nets.messageDialog.Hide()
	}
	if nets.createDialog.IsDisplay() {
		nets.createDialog.Hide()
	}
}
