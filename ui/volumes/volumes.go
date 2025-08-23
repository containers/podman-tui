package volumes

import (
	"fmt"
	"strings"
	"sync"

	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/containers/podman-tui/ui/volumes/voldialogs"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/rivo/tview"
)

// Volumes implemnents the volumes page primitive.
type Volumes struct {
	*tview.Box

	title           string
	headers         []string
	table           *tview.Table
	errorDialog     *dialogs.ErrorDialog
	progressDialog  *dialogs.ProgressDialog
	confirmDialog   *dialogs.ConfirmDialog
	cmdDialog       *dialogs.CommandDialog
	messageDialog   *dialogs.MessageDialog
	sortDialog      *dialogs.SortDialog
	createDialog    *voldialogs.VolumeCreateDialog
	volumeList      volListReport
	confirmData     string
	appFocusHandler func()
}

type volListReport struct {
	mu        sync.Mutex
	report    []*entities.VolumeListReport
	sortBy    string
	ascending bool
}

// NewVolumes returns new vols page view.
func NewVolumes() *Volumes {
	vols := &Volumes{
		Box:            tview.NewBox(),
		title:          "volumes",
		headers:        []string{"driver", "volume name", "created at", "mount point"},
		errorDialog:    dialogs.NewErrorDialog(),
		progressDialog: dialogs.NewProgressDialog(),
		confirmDialog:  dialogs.NewConfirmDialog(),
		messageDialog:  dialogs.NewMessageDialog(""),
		sortDialog:     dialogs.NewSortDialog([]string{"driver", "name", "created", "mount point"}, 2), //nolint:mnd
		createDialog:   voldialogs.NewVolumeCreateDialog(),
		volumeList:     volListReport{sortBy: "created", ascending: true},
	}

	vols.initUI()

	return vols
}

// SetAppFocusHandler sets application focus handler.
func (vols *Volumes) SetAppFocusHandler(handler func()) {
	vols.appFocusHandler = handler
}

// GetTitle returns primitive title.
func (vols *Volumes) GetTitle() string {
	return vols.title
}

// HasFocus returns whether or not this primitive has focus.
func (vols *Volumes) HasFocus() bool {
	if vols.SubDialogHasFocus() {
		return true
	}

	if vols.table.HasFocus() || vols.Box.HasFocus() {
		return true
	}

	return false
}

// SubDialogHasFocus returns whether or not sub dialog primitive has focus.
func (vols *Volumes) SubDialogHasFocus() bool {
	for _, dialog := range vols.getInnerDialogs() {
		if dialog.HasFocus() {
			return true
		}
	}

	return false
}

// Focus is called when this primitive receives focus.
func (vols *Volumes) Focus(delegate func(p tview.Primitive)) {
	for _, dialog := range vols.getInnerDialogs() {
		if dialog.IsDisplay() {
			delegate(dialog)

			return
		}
	}

	delegate(vols.table)
}

// HideAllDialogs hides all sub dialogs.
func (vols *Volumes) HideAllDialogs() {
	for _, dialog := range vols.getInnerDialogs() {
		if dialog.IsDisplay() {
			dialog.Hide()
		}
	}
}

func (vols *Volumes) getInnerDialogs() []utils.UIDialog {
	dialogs := []utils.UIDialog{
		vols.errorDialog,
		vols.progressDialog,
		vols.confirmDialog,
		vols.cmdDialog,
		vols.messageDialog,
		vols.createDialog,
		vols.sortDialog,
	}

	return dialogs
}

func (vols *Volumes) getSelectedItem() string {
	if vols.table.GetRowCount() <= 1 {
		return ""
	}

	row, _ := vols.table.GetSelection()
	volID := vols.table.GetCell(row, 1).Text

	return volID
}

func (vols *Volumes) initUI() {
	vols.cmdDialog = dialogs.NewCommandDialog([][]string{
		{"create", "create a new volume"},
		{"inspect", "display detailed volume's information"},
		{"prune", "remove all unused volumes"},
		{"rm", "remove the selected volume"},
	})

	vols.table = tview.NewTable()

	vols.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(vols.title)))
	vols.table.SetBorderColor(style.BorderColor)
	vols.table.SetBackgroundColor(style.BgColor)
	vols.table.SetTitleColor(style.FgColor)
	vols.table.SetBorder(true)

	for i := range vols.headers {
		vols.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[black::b]%s", strings.ToUpper(vols.headers[i]))). //nolint:perfsprint
														SetExpansion(1).
														SetBackgroundColor(style.PageHeaderBgColor).
														SetTextColor(style.PageHeaderFgColor).
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
	vols.messageDialog.SetCancelFunc(func() {
		vols.messageDialog.Hide()
	})

	// set confirm dialogs functions
	vols.confirmDialog.SetSelectedFunc(func() {
		vols.confirmDialog.Hide()

		switch vols.confirmData {
		case utils.PruneCommandLabel:
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

	// set sort dialog functions
	vols.sortDialog.SetSelectFunc(vols.SortView)
	vols.sortDialog.SetCancelFunc(vols.sortDialog.Hide)
}
