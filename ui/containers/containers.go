package containers

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/containers/podman-tui/ui/containers/cntdialogs"
	"github.com/containers/podman-tui/ui/containers/cntdialogs/vterm"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/rivo/tview"
)

const (
	viewContainersIDColIndex = 0 + iota
	viewContainersImageColIndex
	viewContainersPodColIndex
	viewContainersCreatedAtColIndex
	viewContainersStatusColIndex
	viewContainersNamesColIndex
	viewContainersPortsColIndex
)

var (
	errNoContainerAttach       = errors.New("there is no container to attach")
	errNoContainerHealthCheck  = errors.New("there is no container to perform healthcheck")
	errNoContainerCommit       = errors.New("there is no container to commit")
	errNoContainerStat         = errors.New("there is no container to display stats")
	errNoContainerCheckpoint   = errors.New("there is no container to perform checkpoint")
	errNoContainerExec         = errors.New("there is no container to perform exec")
	errNoContainerDiff         = errors.New("there is no container to display diff")
	errNoContainerInspect      = errors.New("there is no container to inspect")
	errNoContainerKill         = errors.New("there is no container to kill")
	errNoContainerLogs         = errors.New("there is no container to display logs")
	errNoContainerPause        = errors.New("there is no container to pause")
	errNoContainerUnpause      = errors.New("there is no container to unpause")
	errNoContainerPorts        = errors.New("there is no container to display ports")
	errNoContainerRename       = errors.New("there is no container to rename")
	errNoContainerRemove       = errors.New("there is no container to remove")
	errNoContainerStart        = errors.New("there is no container to start")
	errNoContainerStop         = errors.New("there is no container to stop")
	errNoContainerTop          = errors.New("there is no container to display top")
	errEmptyContainerImageName = errors.New("empty container image name")
)

// Containers implements the containers page primitive.
type Containers struct {
	*tview.Box

	title            string
	headers          []string
	table            *tview.Table
	errorDialog      *dialogs.ErrorDialog
	cmdDialog        *dialogs.CommandDialog
	cmdInputDialog   *dialogs.SimpleInputDialog
	confirmDialog    *dialogs.ConfirmDialog
	messageDialog    *dialogs.MessageDialog
	progressDialog   *dialogs.ProgressDialog
	topDialog        *dialogs.TopDialog
	createDialog     *cntdialogs.ContainerCreateDialog
	runDialog        *cntdialogs.ContainerCreateDialog
	execDialog       *cntdialogs.ContainerExecDialog
	terminalDialog   *vterm.VtermDialog
	statsDialog      *cntdialogs.ContainerStatsDialog
	commitDialog     *cntdialogs.ContainerCommitDialog
	checkpointDialog *cntdialogs.ContainerCheckpointDialog
	restoreDialog    *cntdialogs.ContainerRestoreDialog
	sortDialog       *dialogs.SortDialog
	containersList   containerListReport
	selectedID       string
	selectedName     string
	confirmData      string
	fastRefreshChan  chan bool
	appFocusHandler  func()
}

type containerListReport struct {
	mu        sync.Mutex
	report    []entities.ListContainer
	sortBy    string
	ascending bool
}

// NewContainers returns containers page view.
func NewContainers() *Containers {
	containers := &Containers{
		Box:              tview.NewBox(),
		title:            "containers",
		headers:          []string{"container id", "image", "pod", "created", "status", "names", "ports"},
		errorDialog:      dialogs.NewErrorDialog(),
		cmdInputDialog:   dialogs.NewSimpleInputDialog(""),
		messageDialog:    dialogs.NewMessageDialog(""),
		progressDialog:   dialogs.NewProgressDialog(),
		confirmDialog:    dialogs.NewConfirmDialog(),
		topDialog:        dialogs.NewTopDialog(),
		createDialog:     cntdialogs.NewContainerCreateDialog(cntdialogs.ContainerCreateOnlyDialogMode),
		runDialog:        cntdialogs.NewContainerCreateDialog(cntdialogs.ContainerCreateAndRunDialogMode),
		execDialog:       cntdialogs.NewContainerExecDialog(),
		terminalDialog:   vterm.NewVtermDialog(),
		statsDialog:      cntdialogs.NewContainerStatsDialog(),
		commitDialog:     cntdialogs.NewContainerCommitDialog(),
		checkpointDialog: cntdialogs.NewContainerCheckpointDialog(),
		restoreDialog:    cntdialogs.NewContainerRestoreDialog(),
		sortDialog:       dialogs.NewSortDialog([]string{"name", "pod", "image", "created"}, 3), //nolint:mnd
		containersList:   containerListReport{sortBy: "created", ascending: true},
	}

	containers.topDialog.SetTitle("podman container top")

	containers.cmdDialog = dialogs.NewCommandDialog([][]string{
		{"attach", "attach to a running container"},
		{"checkpoint", "checkpoints a running container"},
		{"commit", "create an image from a container's changes"},
		{"create", "create a new container but do not start"},
		{"diff", "inspect changes to the selected container's file systems"},
		{"exec", "execute the specified command inside a running container"},
		{"healthcheck", "run the health check of a container"},
		{"inspect", "display the configuration of a container"},
		{"kill", "kill the selected running container with a SIGKILL signal"},
		{"logs", "fetch the logs of the selected container"},
		{"pause", "pause all the processes in the selected container"},
		{"port", "list port mappings for the selected container"},
		{"prune", "remove all non running containers"},
		{"rename", "rename the selected container"},
		{"restore", "restores a container from a checkpoint"},
		{"rm", "remove the selected container"},
		{"run", "runs a command in a new container from the given image"},
		{"start", "start the selected containers"},
		{"stats", "display container resource usage statistics"},
		{"stop", "stop the selected containers"},
		{"top", "display the running processes of the selected container"},
		{"unpause", "unpause the selected container that was paused before"},
	})

	containers.table = tview.NewTable()
	containers.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(containers.title)))
	containers.table.SetBorderColor(style.BorderColor)
	containers.table.SetTitleColor(style.FgColor)
	containers.table.SetBackgroundColor(style.BgColor)
	containers.table.SetBorder(true)

	for i := range containers.headers {
		containers.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[black::b]%s", strings.ToUpper(containers.headers[i]))). //nolint:perfsprint
															SetExpansion(1).
															SetBackgroundColor(style.PageHeaderBgColor).
															SetTextColor(style.PageHeaderFgColor).
															SetAlign(tview.AlignLeft).
															SetSelectable(false))
	}

	containers.table.SetFixed(1, 1)
	containers.table.SetSelectable(true, false)

	// set command dialog functions
	containers.cmdDialog.SetSelectedFunc(func() {
		containers.cmdDialog.Hide()
		containers.runCommand(containers.cmdDialog.GetSelectedItem())
	})
	containers.cmdDialog.SetCancelFunc(containers.cmdDialog.Hide)

	// set input cmd dialog functions
	containers.cmdInputDialog.SetCancelFunc(containers.cmdInputDialog.Hide)
	containers.cmdInputDialog.SetSelectedFunc(containers.cmdInputDialog.Hide)

	// set message dialog functions
	containers.messageDialog.SetCancelFunc(containers.messageDialog.Hide)

	// set container top dialog functions
	containers.topDialog.SetCancelFunc(containers.topDialog.Hide)

	// set confirm dialogs functions
	containers.confirmDialog.SetSelectedFunc(func() {
		containers.confirmDialog.Hide()

		switch containers.confirmData {
		case utils.PruneCommandLabel:
			containers.prune()
		case "rm":
			containers.remove()
		}
	})

	containers.confirmDialog.SetCancelFunc(containers.confirmDialog.Hide)

	// set create dialog functions
	containers.createDialog.SetCancelFunc(func() {
		containers.createDialog.Hide()
	})

	containers.createDialog.SetHandlerFunc(func() {
		containers.createDialog.Hide()
		containers.create()
	})

	// set run dialog functions
	containers.runDialog.SetCancelFunc(func() {
		containers.runDialog.Hide()
	})

	containers.runDialog.SetHandlerFunc(func() {
		containers.runDialog.Hide()
		containers.run()
	})

	// set exec dialog functions
	containers.execDialog.SetCancelFunc(containers.execDialog.Hide)
	containers.execDialog.SetExecFunc(containers.exec)

	// terminal dialog
	containers.terminalDialog.SetCancelFunc(containers.terminalDialog.Hide)
	containers.terminalDialog.SetFastRefreshHandler(func() {
		containers.fastRefreshChan <- true
	})

	// set stats dialogs functions
	containers.statsDialog.SetDoneFunc(containers.statsDialog.Hide)

	// set commit dialog functions
	containers.commitDialog.SetCommitFunc(containers.commit)
	containers.commitDialog.SetCancelFunc(containers.commitDialog.Hide)

	// set checkpoint dialog functions
	containers.checkpointDialog.SetCheckpointFunc(containers.checkpoint)
	containers.checkpointDialog.SetCancelFunc(containers.checkpointDialog.Hide)

	// set restore dialog functions
	containers.restoreDialog.SetRestoreFunc(containers.restore)
	containers.restoreDialog.SetCancelFunc(containers.restoreDialog.Hide)

	// set sort dialog functions
	containers.sortDialog.SetSelectFunc(containers.SortView)
	containers.sortDialog.SetCancelFunc(containers.sortDialog.Hide)

	return containers
}

// SetAppFocusHandler sets application focus handler.
func (cnt *Containers) SetAppFocusHandler(handler func()) {
	cnt.appFocusHandler = handler
}

// GetTitle returns primitive title.
func (cnt *Containers) GetTitle() string {
	return cnt.title
}

// HasFocus returns whether or not this primitive has focus.
func (cnt *Containers) HasFocus() bool { //nolint:cyclop
	if cnt.table.HasFocus() || cnt.errorDialog.HasFocus() {
		return true
	}

	if cnt.cmdDialog.HasFocus() || cnt.progressDialog.HasFocus() {
		return true
	}

	if cnt.topDialog.HasFocus() || cnt.messageDialog.HasFocus() {
		return true
	}

	if cnt.confirmDialog.HasFocus() || cnt.cmdInputDialog.HasFocus() {
		return true
	}

	if cnt.createDialog.HasFocus() || cnt.execDialog.HasFocus() {
		return true
	}

	if cnt.statsDialog.HasFocus() || cnt.commitDialog.HasFocus() {
		return true
	}

	if cnt.checkpointDialog.HasFocus() || cnt.restoreDialog.HasFocus() {
		return true
	}

	if cnt.runDialog.HasFocus() || cnt.terminalDialog.HasFocus() {
		return true
	}

	if cnt.sortDialog.HasFocus() || cnt.Box.HasFocus() {
		return true
	}

	return false
}

// SubDialogHasFocus returns whether or not sub dialog primitive has focus.
func (cnt *Containers) SubDialogHasFocus() bool { //nolint:cyclop
	if cnt.statsDialog.HasFocus() || cnt.errorDialog.HasFocus() {
		return true
	}

	if cnt.cmdDialog.HasFocus() || cnt.progressDialog.HasFocus() {
		return true
	}

	if cnt.topDialog.HasFocus() || cnt.messageDialog.HasFocus() {
		return true
	}

	if cnt.confirmDialog.HasFocus() || cnt.cmdInputDialog.HasFocus() {
		return true
	}

	if cnt.createDialog.HasFocus() || cnt.execDialog.HasFocus() {
		return true
	}

	if cnt.commitDialog.HasFocus() || cnt.checkpointDialog.HasFocus() {
		return true
	}

	if cnt.restoreDialog.HasFocus() || cnt.terminalDialog.HasFocus() {
		return true
	}

	if cnt.sortDialog.HasFocus() || cnt.runDialog.HasFocus() {
		return true
	}

	return false
}

// Focus is called when this primitive receives focus.
func (cnt *Containers) Focus(delegate func(p tview.Primitive)) { //nolint:cyclop
	// error dialog
	if cnt.errorDialog.IsDisplay() {
		delegate(cnt.errorDialog)

		return
	}

	// command dialog
	if cnt.cmdDialog.IsDisplay() {
		delegate(cnt.cmdDialog)

		return
	}

	// command input dialog
	if cnt.cmdInputDialog.IsDisplay() {
		delegate(cnt.cmdInputDialog)

		return
	}

	// message dialog
	if cnt.messageDialog.IsDisplay() {
		delegate(cnt.messageDialog)

		return
	}

	// container top dialog
	if cnt.topDialog.IsDisplay() {
		delegate(cnt.topDialog)

		return
	}

	// confirm dialog
	if cnt.confirmDialog.IsDisplay() {
		delegate(cnt.confirmDialog)

		return
	}

	// create dialog
	if cnt.createDialog.IsDisplay() {
		delegate(cnt.createDialog)

		return
	}

	// run dialog
	if cnt.runDialog.IsDisplay() {
		delegate(cnt.runDialog)

		return
	}

	// exec dialog
	if cnt.execDialog.IsDisplay() {
		delegate(cnt.execDialog)

		return
	}

	// stats dialog
	if cnt.statsDialog.IsDisplay() {
		delegate(cnt.statsDialog)

		return
	}

	// commit dialog
	if cnt.commitDialog.IsDisplay() {
		delegate(cnt.commitDialog)

		return
	}

	// checkpoint dialog
	if cnt.checkpointDialog.IsDisplay() {
		delegate(cnt.checkpointDialog)

		return
	}

	// restore dialog
	if cnt.restoreDialog.IsDisplay() {
		delegate(cnt.restoreDialog)

		return
	}

	// terminal dialog
	if cnt.terminalDialog.IsDisplay() {
		delegate(cnt.terminalDialog)

		return
	}

	// sort dialog
	if cnt.sortDialog.IsDisplay() {
		delegate(cnt.sortDialog)

		return
	}

	delegate(cnt.table)
}

// SetFastRefreshChannel sets channel for fastRefresh func.
func (cnt *Containers) SetFastRefreshChannel(refresh chan bool) {
	cnt.fastRefreshChan = refresh
}

// HideAllDialogs hides all sub dialogs.
func (cnt *Containers) HideAllDialogs() { //nolint:cyclop
	if cnt.errorDialog.IsDisplay() {
		cnt.errorDialog.Hide()
	}

	if cnt.progressDialog.IsDisplay() {
		cnt.progressDialog.Hide()
	}

	if cnt.confirmDialog.IsDisplay() {
		cnt.confirmDialog.Hide()
	}

	if cnt.cmdDialog.IsDisplay() {
		cnt.cmdDialog.Hide()
	}

	if cnt.cmdInputDialog.IsDisplay() {
		cnt.cmdInputDialog.Hide()
	}

	if cnt.messageDialog.IsDisplay() {
		cnt.messageDialog.Hide()
	}

	if cnt.topDialog.IsDisplay() {
		cnt.topDialog.Hide()
	}

	if cnt.createDialog.IsDisplay() {
		cnt.createDialog.Hide()
	}

	if cnt.runDialog.IsDisplay() {
		cnt.runDialog.Hide()
	}

	if cnt.execDialog.IsDisplay() {
		cnt.execDialog.Hide()
	}

	if cnt.statsDialog.IsDisplay() {
		cnt.statsDialog.Hide()
	}

	if cnt.commitDialog.IsDisplay() {
		cnt.commitDialog.Hide()
	}

	if cnt.checkpointDialog.IsDisplay() {
		cnt.checkpointDialog.Hide()
	}

	if cnt.restoreDialog.IsDisplay() {
		cnt.restoreDialog.Hide()
	}

	if cnt.terminalDialog.IsDisplay() {
		cnt.terminalDialog.Hide()
	}

	if cnt.sortDialog.IsDisplay() {
		cnt.sortDialog.Hide()
	}
}

func (cnt *Containers) getSelectedItem() (string, string) {
	var (
		cntID   string
		cntName string
	)

	if cnt.table.GetRowCount() <= 1 {
		return cntID, cntName
	}

	row, _ := cnt.table.GetSelection()
	cntID = cnt.table.GetCell(row, viewContainersIDColIndex).Text
	cntName = cnt.table.GetCell(row, viewContainersNamesColIndex).Text

	return cntID, cntName
}
