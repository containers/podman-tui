package images

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/images"
	"github.com/containers/podman-tui/ui/style"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// UpdateData retrieves images list data
func (img *Images) UpdateData() {
	images, err := images.List()
	if err != nil {
		log.Error().Msgf("view: images update %v", err)
		img.errorDialog.SetText(fmt.Sprintf("%v", err))
		img.errorDialog.Display()
		return
	}
	img.imagesList.mu.Lock()
	img.imagesList.report = images
	img.imagesList.mu.Unlock()

}

func (img *Images) getData() []images.ImageListReporter {
	img.imagesList.mu.Lock()
	data := img.imagesList.report
	img.imagesList.mu.Unlock()
	return data
}

// ClearData clears table data
func (img *Images) ClearData() {
	img.imagesList.mu.Lock()
	img.imagesList.report = nil
	img.imagesList.mu.Unlock()
	img.table.Clear()
	expand := 1
	fgColor := style.PageHeaderFgColor
	bgColor := style.PageHeaderBgColor
	for i := 0; i < len(img.headers); i++ {
		img.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[::b]%s", strings.ToUpper(img.headers[i]))).
				SetExpansion(expand).
				SetBackgroundColor(bgColor).
				SetTextColor(fgColor).
				SetAlign(tview.AlignLeft).
				SetSelectable(false))
	}
	img.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(img.title)))
}
