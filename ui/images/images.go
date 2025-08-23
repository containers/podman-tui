package images

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/containers/podman-tui/pdcs/images"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/images/imgdialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/rivo/tview"
)

const (
	viewImageRepoNameColIndex = 0 + iota
	viewImageTagColIndex
	viewImageIDColIndex
	viewImageCreatedAtColIndex
	viewImageSizeColIndex
)

var (
	errNoImageToTree       = errors.New("here is no image to display tree")
	errNoImageToUntag      = errors.New("here is no image to untag")
	errNoImageToTag        = errors.New("here is no image to tag")
	errNoImageToSave       = errors.New("here is no image to save")
	errNoImageToPush       = errors.New("here is no image to push")
	errNoImageToHistory    = errors.New("here is no image to display history")
	errNoImageToDiff       = errors.New("here is no image to display diff")
	errNoImageToRemove     = errors.New("there is no image to remove")
	errNoImageToInspect    = errors.New("there is no image to display inspect")
	errNoBuildDirOrCntFile = errors.New("both context directory path and container files fields are empty")
)

// Images implements the images primitive.
type Images struct {
	*tview.Box

	title           string
	headers         []string
	table           *tview.Table
	errorDialog     *dialogs.ErrorDialog
	cmdDialog       *dialogs.CommandDialog
	cmdInputDialog  *dialogs.SimpleInputDialog
	messageDialog   *dialogs.MessageDialog
	confirmDialog   *dialogs.ConfirmDialog
	searchDialog    *imgdialogs.ImageSearchDialog
	historyDialog   *imgdialogs.ImageHistoryDialog
	importDialog    *imgdialogs.ImageImportDialog
	buildDialog     *imgdialogs.ImageBuildDialog
	buildPrgDialog  *imgdialogs.ImageBuildProgressDialog
	progressDialog  *dialogs.ProgressDialog
	saveDialog      *imgdialogs.ImageSaveDialog
	pushDialog      *imgdialogs.ImagePushDialog
	sortDialog      *dialogs.SortDialog
	imagesList      imageListReport
	selectedID      string
	selectedName    string
	confirmData     string
	fastRefreshChan chan bool
	appFocusHandler func()
}

type imageListReport struct {
	mu        sync.Mutex
	report    []images.ImageListReporter
	sortBy    string
	ascending bool
}

// NewImages returns images page view.
func NewImages() *Images {
	images := &Images{
		Box:            tview.NewBox(),
		title:          "images",
		headers:        []string{"repository", "tag", "image id", "created at", "size"},
		errorDialog:    dialogs.NewErrorDialog(),
		cmdInputDialog: dialogs.NewSimpleInputDialog(""),
		messageDialog:  dialogs.NewMessageDialog(""),
		confirmDialog:  dialogs.NewConfirmDialog(),
		searchDialog:   imgdialogs.NewImageSearchDialog(),
		historyDialog:  imgdialogs.NewImageHistoryDialog(),
		importDialog:   imgdialogs.NewImageImportDialog(),
		buildDialog:    imgdialogs.NewImageBuildDialog(),
		buildPrgDialog: imgdialogs.NewImageBuildProgressDialog(),
		saveDialog:     imgdialogs.NewImageSaveDialog(),
		pushDialog:     imgdialogs.NewImagePushDialog(),
		sortDialog:     dialogs.NewSortDialog([]string{"repository", "created", "size"}, 1),
		progressDialog: dialogs.NewProgressDialog(),
		imagesList:     imageListReport{sortBy: "created", ascending: true},
	}

	images.cmdDialog = dialogs.NewCommandDialog([][]string{
		{"build", "build an image from Containerfile"},
		{"diff", "inspect changes to the image's file systems"},
		{"history", "show history of the selected image"},
		{"import", "create a container image from a tarball"},
		{"inspect", "display the configuration of the selected image"},
		{"prune", "remove all unused images"},
		{"push", "push a source image to a specified destination"},
		{"rm", "removes the selected  image from local storage"},
		{"save", "save an image to docker-archive or oci-archive"},
		{"search/pull", "search and pull image from registry"},
		{"tag", "add an additional name to the selected  image"},
		{"tree", "display layer hierarchy of an image"},
		{"untag", "remove a name from the selected image"},
	})

	imgTable := tview.NewTable()
	imgTable.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(images.title)))
	imgTable.SetBorderColor(style.BorderColor)
	imgTable.SetBackgroundColor(style.BgColor)
	imgTable.SetTitleColor(style.FgColor)
	imgTable.SetBorder(true)

	for i := range images.headers {
		imgTable.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[black::b]%s", strings.ToUpper(images.headers[i]))). //nolint:perfsprint
														SetExpansion(1).
														SetBackgroundColor(style.PageHeaderBgColor).
														SetTextColor(style.PageHeaderFgColor).
														SetAlign(tview.AlignLeft).
														SetSelectable(false))
	}

	imgTable.SetFixed(1, 1)
	imgTable.SetSelectable(true, false)
	images.table = imgTable

	// set message dialog functions
	images.messageDialog.SetCancelFunc(func() {
		images.messageDialog.Hide()
	})

	// set input cmd dialog functions
	images.cmdInputDialog.SetCancelFunc(func() {
		images.cmdInputDialog.Hide()
	})

	images.cmdInputDialog.SetSelectedFunc(func() {
		images.cmdInputDialog.Hide()
	})

	// set command dialogs functions
	images.cmdDialog.SetSelectedFunc(func() {
		images.cmdDialog.Hide()
		images.runCommand(images.cmdDialog.GetSelectedItem())
	})
	images.cmdDialog.SetCancelFunc(func() {
		images.cmdDialog.Hide()
	})

	// set confirm dialogs functions
	images.confirmDialog.SetSelectedFunc(func() {
		images.confirmDialog.Hide()

		switch images.confirmData {
		case utils.PruneCommandLabel:
			images.prune()
		case "rm":
			images.remove()
		}
	})

	images.confirmDialog.SetCancelFunc(func() {
		images.confirmDialog.Hide()
	})

	// set history dialogs functions
	images.historyDialog.SetCancelFunc(func() {
		images.historyDialog.Hide()
	})

	// set search dialogs functions
	images.searchDialog.SetCancelFunc(func() {
		images.searchDialog.Hide()
	})

	images.searchDialog.SetSearchFunc(func() {
		term := images.searchDialog.GetSearchText()
		if term == "" {
			return
		}

		images.search(term)
	})

	images.searchDialog.SetPullFunc(func() {
		name := images.searchDialog.GetSelectedItem()
		images.pull(name)
	})

	// set build dialogs functions
	images.buildDialog.SetCancelFunc(images.buildDialog.Hide)
	images.buildDialog.SetBuildFunc(images.build)
	images.buildPrgDialog.SetFastRefreshHandler(func() {
		images.fastRefreshChan <- true
	})

	// set save dialog functions
	images.saveDialog.SetCancelFunc(images.saveDialog.Hide)
	images.saveDialog.SetSaveFunc(images.save)

	// set import dialog functions
	images.importDialog.SetCancelFunc(images.importDialog.Hide)
	images.importDialog.SetImportFunc(images.imageImport)

	// set push dialog functions
	images.pushDialog.SetPushFunc(images.push)
	images.pushDialog.SetCancelFunc(images.pushDialog.Hide)

	// set sort dialog functions
	images.sortDialog.SetSelectFunc(images.SortView)
	images.sortDialog.SetCancelFunc(images.sortDialog.Hide)

	return images
}

// SetAppFocusHandler sets application focus handler.
func (img *Images) SetAppFocusHandler(handler func()) {
	img.appFocusHandler = handler
}

// GetTitle returns primitive title.
func (img *Images) GetTitle() string {
	return img.title
}

// HasFocus returns whether or not this primitive has focus.
func (img *Images) HasFocus() bool {
	if img.SubDialogHasFocus() || img.table.HasFocus() {
		return true
	}

	return img.Box.HasFocus()
}

// SubDialogHasFocus returns whether or not sub dialog primitive has focus.
func (img *Images) SubDialogHasFocus() bool {
	for _, dialog := range img.getInnerDialogs() {
		if dialog.HasFocus() {
			return true
		}
	}

	for _, dialog := range img.getInnerTopDialogs() {
		if dialog.HasFocus() {
			return true
		}
	}

	return false
}

// Focus is called when this primitive receives focus.
func (img *Images) Focus(delegate func(p tview.Primitive)) {
	// since error and confirm dialog can get focus on top of other dialogs
	if img.errorDialog.IsDisplay() {
		delegate(img.errorDialog)

		return
	}

	if img.confirmDialog.IsDisplay() {
		delegate(img.confirmDialog)

		return
	}

	for _, dialog := range img.getInnerDialogs() {
		if dialog.IsDisplay() {
			delegate(dialog)

			return
		}
	}

	for _, dialog := range img.getInnerTopDialogs() {
		if dialog.IsDisplay() {
			delegate(dialog)

			return
		}
	}

	delegate(img.table)
}

// HideAllDialogs hides all sub dialogs.
func (img *Images) HideAllDialogs() {
	for _, dialog := range img.getInnerDialogs() {
		if dialog.IsDisplay() {
			dialog.Hide()
		}
	}

	for _, dialog := range img.getInnerTopDialogs() {
		if dialog.IsDisplay() {
			dialog.Hide()
		}
	}
}

// SetFastRefreshChannel sets channel for fastRefresh func.
func (img *Images) SetFastRefreshChannel(refresh chan bool) {
	img.fastRefreshChan = refresh
}

func (img *Images) getSelectedItem() (string, string) {
	if img.table.GetRowCount() <= 1 {
		return "", ""
	}

	row, _ := img.table.GetSelection()
	imageRepo := img.table.GetCell(row, 0).Text
	imageTag := img.table.GetCell(row, 1).Text
	imageName := imageRepo + ":" + imageTag
	imageID := img.table.GetCell(row, 2).Text //nolint:mnd

	return imageID, imageName
}

func (img *Images) getInnerDialogs() []utils.UIDialog {
	dialogs := []utils.UIDialog{
		img.cmdDialog,
		img.cmdInputDialog,
		img.messageDialog,
		img.searchDialog,
		img.historyDialog,
		img.importDialog,
		img.buildDialog,
		img.buildPrgDialog,
		img.saveDialog,
		img.pushDialog,
		img.sortDialog,
	}

	return dialogs
}

func (img *Images) getInnerTopDialogs() []utils.UIDialog {
	dialogs := []utils.UIDialog{
		img.errorDialog,
		img.progressDialog,
		img.confirmDialog,
	}

	return dialogs
}
