package volumes

import (
	"fmt"
	"strings"
	"sync"

	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/containers/podman-tui/ui/volumes/voldialogs"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/rivo/tview"
)

// Volumes implemnents the volumes page primitive
type Volumes struct {
	*tview.Box
	title          string
	headers        []string
	table          *tview.Table
	errorDialog    *dialogs.ErrorDialog
	progressDialog *dialogs.ProgressDialog
	confirmDialog  *dialogs.ConfirmDialog
	cmdDialog      *dialogs.CommandDialog
	messageDialog  *dialogs.MessageDialog
	createDialog   *voldialogs.VolumeCreateDialog
	volumeList     volListReport
	selectedID     string
	confirmData    string
}

type volListReport struct {
	mu     sync.Mutex
	report []*entities.VolumeListReport
}

// NewVolumes returns new vols page view
func NewVolumes() *Volumes {
	vols := &Volumes{
		Box:            tview.NewBox(),
		title:          "volumes",
		headers:        []string{"driver", "volume name", "created at", "mount point"},
		errorDialog:    dialogs.NewErrorDialog(),
		progressDialog: dialogs.NewProgressDialog(),
		confirmDialog:  dialogs.NewConfirmDialog(),
		messageDialog:  dialogs.NewMessageDialog(""),
		createDialog:   voldialogs.NewVolumeCreateDialog(),
	}

	vols.cmdDialog = dialogs.NewCommandDialog([][]string{
		{"create", "create a new volume"},
		{"inspect", "display detailed volume's information"},
		{"prune", "remove all unused volumes"},
		{"rm", "remove the selected volume"},
	})
	fgColor := utils.Styles.PageTable.FgColor
	bgColor := utils.Styles.PageTable.BgColor
	vols.table = tview.NewTable()
	vols.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(vols.title)))
	vols.table.SetBorderColor(bgColor)
	vols.table.SetTitleColor(fgColor)
	vols.table.SetBorder(true)

	fgColor = utils.Styles.PageTable.HeaderRow.FgColor
	bgColor = utils.Styles.PageTable.HeaderRow.BgColor

	for i := 0; i < len(vols.headers); i++ {
		vols.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[black::b]%s", strings.ToUpper(vols.headers[i]))).
				SetExpansion(1).
				SetBackgroundColor(bgColor).
				SetTextColor(fgColor).
				SetAlign(tview.AlignLeft).
				SetSelectable(false))
	}

	vols.table.SetFixed(1, 1)
	vols.table.SetSelectable(true, false)

	// set command dialog functions
	vols.cmdDialog.SetSelectedFunc(func() {
		vols.cmdDialog.Hide()
		vols.runCommand(vols.cmdDialog.GetSelectedItem())
	})
	vols.cmdDialog.SetCancelFunc(func() {
		vols.cmdDialog.Hide()
	})

	// set message dialog functions
	vols.messageDialog.SetSelectedFunc(func() {
		vols.messageDialog.Hide()
	})
	vols.messageDialog.SetCancelFunc(func() {
		vols.messageDialog.Hide()
	})
	// set confirm dialogs functions
	vols.confirmDialog.SetSelectedFunc(func() {
		vols.confirmDialog.Hide()
		switch vols.confirmData {
		case "prune":
			vols.prune()
		case "rm":
			vols.remove()
		}
	})
	vols.confirmDialog.SetCancelFunc(func() {
		vols.confirmDialog.Hide()
	})
	// set create dialog functions
	vols.createDialog.SetCancelFunc(func() {
		vols.createDialog.Hide()
	})
	vols.createDialog.SetCreateFunc(func() {
		vols.createDialog.Hide()
		vols.create()
	})

	return vols
}

// GetTitle returns primitive title
func (vols *Volumes) GetTitle() string {
	return vols.title
}

// HasFocus returns whether or not this primitive has focus
func (vols *Volumes) HasFocus() bool {
	if vols.table.HasFocus() || vols.errorDialog.HasFocus() {
		return true
	}
	if vols.cmdDialog.HasFocus() || vols.messageDialog.IsDisplay() {
		return true
	}
	if vols.progressDialog.HasFocus() || vols.confirmDialog.HasFocus() {
		return true
	}
	if vols.createDialog.HasFocus() {
		return true
	}
	return vols.Box.HasFocus()
}

// SubDialogHasFocus returns whether or not sub dialog primitive has focus
func (vols *Volumes) SubDialogHasFocus() bool {
	if vols.errorDialog.HasFocus() || vols.createDialog.HasFocus() {
		return true
	}
	if vols.cmdDialog.HasFocus() || vols.messageDialog.IsDisplay() {
		return true
	}
	if vols.progressDialog.HasFocus() || vols.confirmDialog.HasFocus() {
		return true
	}
	return false
}

// Focus is called when this primitive receives focus
func (vols *Volumes) Focus(delegate func(p tview.Primitive)) {
	// error dialog
	if vols.errorDialog.IsDisplay() {
		delegate(vols.errorDialog)
		return
	}
	// command dialog
	if vols.cmdDialog.IsDisplay() {
		delegate(vols.cmdDialog)
		return
	}
	// message dialog
	if vols.messageDialog.IsDisplay() {
		delegate(vols.messageDialog)
		return
	}
	// confirm dialog
	if vols.confirmDialog.IsDisplay() {
		delegate(vols.confirmDialog)
		return
	}
	// create dialog
	if vols.createDialog.IsDisplay() {
		delegate(vols.createDialog)
		return
	}
	delegate(vols.table)
}

func (vols *Volumes) getSelectedItem() string {
	if vols.table.GetRowCount() <= 1 {
		return ""
	}
	row, _ := vols.table.GetSelection()
	podID := vols.table.GetCell(row, 1).Text
	return podID
}

// HideAllDialogs hides all sub dialogs
func (vols *Volumes) HideAllDialogs() {
	if vols.errorDialog.IsDisplay() {
		vols.errorDialog.Hide()
	}
	if vols.progressDialog.IsDisplay() {
		vols.progressDialog.Hide()
	}
	if vols.confirmDialog.IsDisplay() {
		vols.confirmDialog.Hide()
	}
	if vols.cmdDialog.IsDisplay() {
		vols.cmdDialog.Hide()
	}
	if vols.messageDialog.IsDisplay() {
		vols.messageDialog.Hide()
	}
	if vols.createDialog.IsDisplay() {
		vols.createDialog.Hide()
	}
}
