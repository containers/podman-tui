package pods

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/pods/poddialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/rivo/tview"
)

const (
	viewPodIDColIndex = 0 + iota
	viewPodNameColIndex
	viewPodStatusColIndex
	viewPodCreatedColIndex
	viewPodInfraIDColIndex
	viewPodContainersColIndex
)

var (
	errNoPodUnpause = errors.New("there is no pod to unpause")
	errNoPodPause   = errors.New("there is no pod to pause")
	errNoPodTop     = errors.New("there is no pod to display top")
	errNoPodStop    = errors.New("there is no pod to stop")
	errNoPodStart   = errors.New("there is no pod to start")
	errNoPodRemove  = errors.New("there is no pod to remove")
	errNoPodRestart = errors.New("there is no pod to restart")
	errNoPodKill    = errors.New("there is no pod to kill")
	errNoPodInspect = errors.New("there is no pod to display inspect")
	errNoPodStat    = errors.New("there is no pod to display stats")
	errPodRemove    = errors.New("remove error")
	errPodPrune     = errors.New("prune error")
)

// Pods implemnents the pods page primitive.
type Pods struct {
	*tview.Box

	title           string
	headers         []string
	table           *tview.Table
	errorDialog     *dialogs.ErrorDialog
	progressDialog  *dialogs.ProgressDialog
	confirmDialog   *dialogs.ConfirmDialog
	cmdDialog       *dialogs.CommandDialog
	messageDialog   *dialogs.MessageDialog
	topDialog       *dialogs.TopDialog
	sortDialog      *dialogs.SortDialog
	createDialog    *poddialogs.PodCreateDialog
	statsDialog     *poddialogs.PodStatsDialog
	podsList        podsListReport
	selectedID      string
	confirmData     string
	appFocusHandler func()
}

type podsListReport struct {
	mu        sync.Mutex
	report    []*entities.ListPodsReport
	sortBy    string
	ascending bool
}

// NewPods returns pods page view.
func NewPods() *Pods {
	pods := &Pods{
		Box:            tview.NewBox(),
		title:          "pods",
		headers:        []string{"pod id", "name", "status", "created", "infra id", "# of containers"},
		errorDialog:    dialogs.NewErrorDialog(),
		confirmDialog:  dialogs.NewConfirmDialog(),
		progressDialog: dialogs.NewProgressDialog(),
		messageDialog:  dialogs.NewMessageDialog(""),
		topDialog:      dialogs.NewTopDialog(),
		sortDialog:     dialogs.NewSortDialog([]string{"name", "created", "status", "# of containers"}, 1),
		createDialog:   poddialogs.NewPodCreateDialog(),
		statsDialog:    poddialogs.NewPodStatsDialog(),
		podsList:       podsListReport{sortBy: "created", ascending: true},
	}

	pods.topDialog.SetTitle("podman pod top")

	pods.cmdDialog = dialogs.NewCommandDialog([][]string{
		{"create", "create a new pod"},
		{"inspect", "display information describing the selected pod"},
		{"kill", "send SIGTERM signal to containers in the pod"},
		{"pause", "pause  the selected pod"},
		{"prune", "remove all stopped pods and their containers"},
		{"restart", "restart  the selected pod"},
		{"rm", "remove the selected pod"},
		{"start", "start  the selected pod"},
		{"stats", "display live stream of resource usage"},
		{"stop", "stop the selected pod"},
		{"top", "display the running processes of the pod's containers"},
		{"unpause", "unpause  the selected pod"},
	})

	pods.table = tview.NewTable()
	pods.table.SetBackgroundColor(style.BgColor)
	pods.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(pods.title)))
	pods.table.SetBorderColor(style.BorderColor)
	pods.table.SetTitleColor(style.FgColor)
	pods.table.SetBorder(true)

	for i := range pods.headers {
		pods.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[black::b]%s", strings.ToUpper(pods.headers[i]))). //nolint:perfsprint
														SetExpansion(1).
														SetBackgroundColor(style.PageHeaderBgColor).
														SetTextColor(style.PageHeaderFgColor).
														SetAlign(tview.AlignLeft).
														SetSelectable(false))
	}

	pods.table.SetFixed(1, 1)
	pods.table.SetSelectable(true, false)

	// set command dialog functions
	pods.cmdDialog.SetSelectedFunc(func() {
		pods.cmdDialog.Hide()
		pods.runCommand(pods.cmdDialog.GetSelectedItem())
	})
	pods.cmdDialog.SetCancelFunc(func() {
		pods.cmdDialog.Hide()
	})

	// set message dialog functions
	pods.messageDialog.SetCancelFunc(func() {
		pods.messageDialog.Hide()
	})

	// set top dialog functions
	pods.topDialog.SetCancelFunc(func() {
		pods.topDialog.Hide()
	})

	// set confirm dialogs functions
	pods.confirmDialog.SetSelectedFunc(func() {
		pods.confirmDialog.Hide()

		switch pods.confirmData {
		case utils.PruneCommandLabel:
			pods.prune()
		case "rm":
			pods.remove()
		}
	})

	pods.confirmDialog.SetCancelFunc(func() {
		pods.confirmDialog.Hide()
	})

	// set create dialog functions
	pods.createDialog.SetCancelFunc(func() {
		pods.createDialog.Hide()
	})

	pods.createDialog.SetCreateFunc(func() {
		pods.createDialog.Hide()
		pods.create()
	})

	// set stats dialog functions
	pods.statsDialog.SetDoneFunc(pods.statsDialog.Hide)

	// set sort dialog functions
	pods.sortDialog.SetCancelFunc(pods.sortDialog.Hide)
	pods.sortDialog.SetSelectFunc(pods.SortView)

	return pods
}

// SetAppFocusHandler sets application focus handler.
func (pods *Pods) SetAppFocusHandler(handler func()) {
	pods.appFocusHandler = handler
}

// GetTitle returns primitive title.
func (pods *Pods) GetTitle() string {
	return pods.title
}

// HasFocus returns whether or not this primitive has focus.
func (pods *Pods) HasFocus() bool { //nolint:cyclop
	if pods.table.HasFocus() || pods.errorDialog.HasFocus() {
		return true
	}

	if pods.cmdDialog.HasFocus() || pods.messageDialog.IsDisplay() {
		return true
	}

	if pods.progressDialog.HasFocus() || pods.topDialog.HasFocus() {
		return true
	}

	if pods.confirmDialog.HasFocus() || pods.createDialog.HasFocus() {
		return true
	}

	if pods.statsDialog.HasFocus() || pods.sortDialog.HasFocus() {
		return true
	}

	return pods.Box.HasFocus()
}

// SubDialogHasFocus returns whether or not sub dialog primitive has focus.
func (pods *Pods) SubDialogHasFocus() bool {
	if pods.statsDialog.HasFocus() || pods.errorDialog.HasFocus() {
		return true
	}

	if pods.cmdDialog.HasFocus() || pods.messageDialog.IsDisplay() {
		return true
	}

	if pods.progressDialog.HasFocus() || pods.topDialog.HasFocus() {
		return true
	}

	if pods.confirmDialog.HasFocus() || pods.createDialog.HasFocus() {
		return true
	}

	return pods.sortDialog.HasFocus()
}

// Focus is called when this primitive receives focus.
func (pods *Pods) Focus(delegate func(p tview.Primitive)) {
	// error dialog
	if pods.errorDialog.IsDisplay() {
		delegate(pods.errorDialog)

		return
	}

	// command dialog
	if pods.cmdDialog.IsDisplay() {
		delegate(pods.cmdDialog)

		return
	}

	// message dialog
	if pods.messageDialog.IsDisplay() {
		delegate(pods.messageDialog)

		return
	}

	// top dialog
	if pods.topDialog.IsDisplay() {
		delegate(pods.topDialog)

		return
	}

	// confirm dialog
	if pods.confirmDialog.IsDisplay() {
		delegate(pods.confirmDialog)

		return
	}

	// create dialog
	if pods.createDialog.IsDisplay() {
		delegate(pods.createDialog)

		return
	}

	// stats dialog
	if pods.statsDialog.IsDisplay() {
		delegate(pods.statsDialog)

		return
	}

	// sort dialog
	if pods.sortDialog.IsDisplay() {
		delegate(pods.sortDialog)

		return
	}

	delegate(pods.table)
}

// HideAllDialogs hides all sub dialogs.
func (pods *Pods) HideAllDialogs() {
	if pods.errorDialog.IsDisplay() {
		pods.errorDialog.Hide()
	}

	if pods.progressDialog.IsDisplay() {
		pods.progressDialog.Hide()
	}

	if pods.confirmDialog.IsDisplay() {
		pods.confirmDialog.Hide()
	}

	if pods.cmdDialog.IsDisplay() {
		pods.cmdDialog.Hide()
	}

	if pods.messageDialog.IsDisplay() {
		pods.messageDialog.Hide()
	}

	if pods.topDialog.IsDisplay() {
		pods.topDialog.Hide()
	}

	if pods.createDialog.IsDisplay() {
		pods.createDialog.Hide()
	}

	if pods.statsDialog.IsDisplay() {
		pods.statsDialog.Hide()
	}

	if pods.sortDialog.IsDisplay() {
		pods.sortDialog.Hide()
	}
}

func (pods *Pods) getSelectedItem() (string, string) {
	var (
		id   string
		name string
	)

	if pods.table.GetRowCount() <= 1 {
		return id, name
	}

	row, _ := pods.table.GetSelection()
	id = pods.table.GetCell(row, 0).Text
	name = pods.table.GetCell(row, 1).Text

	return id, name
}

func (pods *Pods) getAllItemsForStats() []poddialogs.PodStatsDropDownOptions {
	var items []poddialogs.PodStatsDropDownOptions

	rows := pods.table.GetRowCount()

	for i := 1; i < rows; i++ {
		podID := pods.table.GetCell(i, 0).Text
		podName := pods.table.GetCell(i, 1).Text

		items = append(items, poddialogs.PodStatsDropDownOptions{
			ID:   podID,
			Name: podName,
		})
	}

	return items
}
