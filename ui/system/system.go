package system

import (
	"strings"

	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/system/sysdialogs"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/rivo/tview"
)

// System implemnents the system information page primitive
type System struct {
	*tview.Box
	title          string
	textview       *tview.TextView
	cmdDialog      *dialogs.CommandDialog
	confirmDialog  *dialogs.ConfirmDialog
	messageDialog  *dialogs.MessageDialog
	progressDialog *dialogs.ProgressDialog
	errorDialog    *dialogs.ErrorDialog
	dfDialog       *sysdialogs.DfDialog
	confirmData    string
}

// NewSystem returns new system page view
func NewSystem() *System {
	sys := &System{
		Box:            tview.NewBox(),
		title:          "system",
		confirmDialog:  dialogs.NewConfirmDialog(),
		progressDialog: dialogs.NewProgressDialog(),
		errorDialog:    dialogs.NewErrorDialog(),
		messageDialog:  dialogs.NewMessageDialog(""),
		dfDialog:       sysdialogs.NewDfDialog(),
	}
	fgColor := utils.Styles.PageTable.FgColor
	bgColor := utils.Styles.PageTable.BgColor

	sys.textview = tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetTextAlign(tview.AlignLeft)

	sys.textview.SetTitleAlign(tview.AlignCenter)
	sys.textview.SetTitle("podman system")
	sys.textview.SetBorder(true)
	sys.textview.SetBorderColor(bgColor)
	sys.textview.SetTitleColor(fgColor)

	sys.cmdDialog = dialogs.NewCommandDialog([][]string{
		{"disk usage", "display Podman related system information"},
		{"info", "display system information"},
		{"prune", "remove all unused pod, container, image and volume data"},
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
		}
	})
	sys.confirmDialog.SetCancelFunc(func() {
		sys.confirmDialog.Hide()
	})

	// set message dialog functions
	sys.messageDialog.SetSelectedFunc(func() {
		sys.messageDialog.Hide()
	})
	sys.messageDialog.SetCancelFunc(func() {
		sys.messageDialog.Hide()
	})

	// set disk usage function
	sys.dfDialog.SetDoneFunc(func() {
		sys.dfDialog.Hide()
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
	if sys.textview.HasFocus() || sys.dfDialog.HasFocus() {
		return true
	}
	if sys.messageDialog.HasFocus() {
		return true
	}
	return sys.Box.HasFocus()
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
	delegate(sys.textview)
}

// SetEventMessage appends podman events to textview
func (sys *System) SetEventMessage(messages []string) {
	msg := strings.Join(messages, "\n")
	sys.textview.SetText(msg)
	sys.textview.ScrollToEnd()
}
