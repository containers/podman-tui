package images

import (
	"fmt"
	"strings"
	"sync"

	"github.com/containers/podman-tui/pdcs/images"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/images/imgdialogs"
	"github.com/containers/podman-tui/ui/utils"

	"github.com/rivo/tview"
)

// Images implements the images primitive
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
	imagesList      imageListReport
	selectedID      string
	selectedName    string
	confirmData     string
	fastRefreshChan chan bool
}

type imageListReport struct {
	mu     sync.Mutex
	report []images.ImageListReporter
}

// NewImages returns images page view
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
		progressDialog: dialogs.NewProgressDialog(),
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

	fgColor := utils.Styles.PageTable.FgColor
	bgColor := utils.Styles.PageTable.BgColor
	imgTable := tview.NewTable()
	imgTable.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(images.title)))
	imgTable.SetBorderColor(bgColor)
	imgTable.SetTitleColor(fgColor)
	imgTable.SetBorder(true)
	fgColor = utils.Styles.PageTable.HeaderRow.FgColor
	bgColor = utils.Styles.PageTable.HeaderRow.BgColor

	for i := 0; i < len(images.headers); i++ {
		imgTable.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[black::b]%s", strings.ToUpper(images.headers[i]))).
				SetExpansion(1).
				SetBackgroundColor(bgColor).
				SetTextColor(fgColor).
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
		case "prune":
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

	return images
}

// GetTitle returns primitive title
func (img *Images) GetTitle() string {
	return img.title
}

// HasFocus returns whether or not this primitive has focus
func (img *Images) HasFocus() bool {
	if img.table.HasFocus() || img.messageDialog.HasFocus() {
		return true
	}
	if img.cmdDialog.HasFocus() || img.cmdInputDialog.HasFocus() {
		return true
	}
	if img.confirmDialog.HasFocus() || img.errorDialog.HasFocus() {
		return true
	}
	if img.searchDialog.HasFocus() || img.progressDialog.HasFocus() {
		return true
	}
	if img.historyDialog.HasFocus() || img.buildDialog.HasFocus() {
		return true
	}
	if img.buildPrgDialog.HasFocus() || img.saveDialog.HasFocus() {
		return true
	}
	if img.importDialog.HasFocus() || img.pushDialog.HasFocus() {
		return true
	}
	return img.Box.HasFocus()
}

// SubDialogHasFocus returns whether or not sub dialog primitive has focus
func (img *Images) SubDialogHasFocus() bool {
	if img.historyDialog.HasFocus() || img.messageDialog.HasFocus() {
		return true
	}
	if img.cmdDialog.HasFocus() || img.cmdInputDialog.HasFocus() {
		return true
	}
	if img.confirmDialog.HasFocus() || img.errorDialog.HasFocus() {
		return true
	}
	if img.searchDialog.HasFocus() || img.progressDialog.HasFocus() {
		return true
	}
	if img.buildDialog.HasFocus() || img.buildPrgDialog.HasFocus() {
		return true
	}
	if img.saveDialog.HasFocus() || img.importDialog.HasFocus() {
		return true
	}
	return img.pushDialog.HasFocus()
}

// Focus is called when this primitive receives focus
func (img *Images) Focus(delegate func(p tview.Primitive)) {

	// error dialog
	if img.errorDialog.IsDisplay() {
		delegate(img.errorDialog)
		return
	}
	// command dialog
	if img.cmdDialog.IsDisplay() {
		delegate(img.cmdDialog)
		return
	}
	// command input dialog
	if img.cmdInputDialog.IsDisplay() {
		delegate(img.cmdInputDialog)
		return
	}
	// message dialog
	if img.messageDialog.IsDisplay() {
		delegate(img.messageDialog)
		return
	}
	// confirm dialog
	if img.confirmDialog.IsDisplay() {
		delegate(img.confirmDialog)
		return
	}
	// search dialog
	if img.searchDialog.IsDisplay() {
		delegate(img.searchDialog)
		return
	}
	// history dialog
	if img.historyDialog.IsDisplay() {
		delegate(img.historyDialog)
		return
	}
	// build dialog
	if img.buildDialog.IsDisplay() {
		delegate(img.buildDialog)
		return
	}
	// build progress dialog
	if img.buildPrgDialog.IsDisplay() {
		delegate(img.buildPrgDialog)
		return
	}
	// save dialog
	if img.saveDialog.IsDisplay() {
		delegate(img.saveDialog)
		return
	}
	// import dialog
	if img.importDialog.IsDisplay() {
		delegate(img.importDialog)
		return
	}
	// push dialog
	if img.pushDialog.IsDisplay() {
		delegate(img.pushDialog)
		return
	}
	delegate(img.table)
}

func (img *Images) getSelectedItem() (string, string) {
	if img.table.GetRowCount() <= 1 {
		return "", ""
	}
	row, _ := img.table.GetSelection()
	imageRepo := img.table.GetCell(row, 0).Text
	imageTag := img.table.GetCell(row, 1).Text
	imageName := imageRepo + ":" + imageTag
	imageID := img.table.GetCell(row, 2).Text
	return imageID, imageName
}

// HideAllDialogs hides all sub dialogs
func (img *Images) HideAllDialogs() {
	if img.errorDialog.IsDisplay() {
		img.errorDialog.Hide()
	}
	if img.progressDialog.IsDisplay() {
		img.progressDialog.Hide()
	}
	if img.cmdDialog.IsDisplay() {
		img.cmdDialog.Hide()
	}
	if img.cmdInputDialog.IsDisplay() {
		img.cmdInputDialog.Hide()
	}
	if img.messageDialog.IsDisplay() {
		img.messageDialog.Hide()
	}
	if img.searchDialog.IsDisplay() {
		img.searchDialog.Hide()
	}
	if img.confirmDialog.IsDisplay() {
		img.confirmDialog.Hide()
	}
	if img.historyDialog.IsDisplay() {
		img.historyDialog.Hide()
	}
	if img.buildDialog.IsDisplay() {
		img.buildDialog.Hide()
	}
	if img.buildPrgDialog.IsDisplay() {
		img.buildPrgDialog.Hide()
	}
	if img.saveDialog.IsDisplay() {
		img.saveDialog.Hide()
	}
	if img.importDialog.IsDisplay() {
		img.importDialog.Hide()
	}
	if img.pushDialog.IsDisplay() {
		img.pushDialog.Hide()
	}
}

// SetFastRefreshChannel sets channel for fastRefresh func
func (img *Images) SetFastRefreshChannel(refresh chan bool) {
	img.fastRefreshChan = refresh
}
