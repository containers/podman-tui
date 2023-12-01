package images

import (
	"fmt"
	"strings"

	putils "github.com/containers/podman-tui/pdcs/utils"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/rivo/tview"
)

func (img *Images) refresh() {
	img.table.Clear()

	expand := 1
	alignment := tview.AlignLeft

	for i := 0; i < len(img.headers); i++ {
		img.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[::b]%s", strings.ToUpper(img.headers[i]))).
				SetExpansion(expand).
				SetBackgroundColor(style.PageHeaderBgColor).
				SetTextColor(style.PageHeaderFgColor).
				SetAlign(tview.AlignLeft).
				SetSelectable(false))
	}

	rowIndex := 1
	images := img.getData()

	img.table.SetTitle(fmt.Sprintf("[::b]%s[%d]", strings.ToUpper(img.title), len(images)))

	for i := 0; i < len(images); i++ {
		repo := images[i].Repository
		tag := images[i].Tag
		imgID := images[i].ID
		imgIDString := imgID

		if len(imgID) > utils.IDLength {
			imgIDString = imgIDString[:utils.IDLength]
		}

		size := putils.SizeToStr(images[i].Size)
		created := putils.CreatedToStr(images[i].Created)

		// repository name column
		img.table.SetCell(rowIndex, viewImageRepoNameColIndex,
			tview.NewTableCell(repo).
				SetExpansion(expand).
				SetAlign(alignment))

		// tag column
		img.table.SetCell(rowIndex, viewImageTagColIndex,
			tview.NewTableCell(tag).
				SetExpansion(expand).
				SetAlign(alignment))

		// id column
		img.table.SetCell(rowIndex, viewImageIDColIndex,
			tview.NewTableCell(imgIDString).
				SetExpansion(expand).
				SetAlign(alignment))

		// created at column
		img.table.SetCell(rowIndex, viewImageCreatedAtColIndex,
			tview.NewTableCell(created).
				SetExpansion(expand).
				SetAlign(alignment))

		// size column
		img.table.SetCell(rowIndex, viewImageSizeColIndex,
			tview.NewTableCell(size).
				SetExpansion(expand).
				SetAlign(alignment))

		rowIndex++
	}
}
