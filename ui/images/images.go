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
	title          string
	headers        []string
	table          *tview.Table
	errorDialog    *dialogs.ErrorDialog
	cmdDialog      *dialogs.CommandDialog
	cmdInputDialog *dialogs.SimpleInputDialog
	messageDialog  *dialogs.MessageDialog
	confirmDialog  *dialogs.ConfirmDialog
	searchDialog   *imgdialogs.ImageSearchDialog
	historyDialog  *imgdialogs.ImageHistoryDialog
	progressDialog *dialogs.ProgressDialog
	imagesList     imageListReport
	selectedID     string
	selectedName   string
	confirmData    string
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
		progressDialog: dialogs.NewProgressDialog(),
	}

	images.cmdDialog = dialogs.NewCommandDialog([][]string{
		{"diff", "inspect changes to the image's file systems"},
		{"history", "show history of the selected image"},
		{"inspect", "display the configuration of the selected image"},
		{"prune", "remove all unused images"},
		{"rm", "removes the selected  image from local storage"},
		{"search/pull", "search and pull image from registry"},
		{"tag", "add an additional name to the selected  image"},
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
	images.messageDialog.SetSelectedFunc(func() {
		images.messageDialog.Hide()
	})

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
	if img.historyDialog.HasFocus() {
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
	return false
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
}
