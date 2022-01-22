package volumes

import (
	"fmt"
	"strings"
	"time"

	"github.com/containers/podman-tui/ui/utils"
	"github.com/docker/go-units"
	"github.com/rivo/tview"
)

func (vols *Volumes) refresh() {

	vols.table.Clear()
	expand := 1
	alignment := tview.AlignLeft
	fgColor := utils.Styles.PageTable.HeaderRow.FgColor
	bgColor := utils.Styles.PageTable.HeaderRow.BgColor

	for i := 0; i < len(vols.headers); i++ {
		vols.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[black::b]%s", strings.ToUpper(vols.headers[i]))).
				SetExpansion(expand).
				SetBackgroundColor(bgColor).
				SetTextColor(fgColor).
				SetAlign(tview.AlignLeft).
				SetSelectable(false))
	}
	rowIndex := 1

	volList := vols.getData()

	vols.table.SetTitle(fmt.Sprintf("[::b]%s[%d]", strings.ToUpper(vols.title), len(volList)))
	for i := 0; i < len(volList); i++ {
		volDriver := volList[i].Driver
		volName := volList[i].Name
		volCreatedAt := units.HumanDuration(time.Since(volList[i].CreatedAt)) + " ago"
		volMountPoint := volList[i].Mountpoint

		// driver name column
		vols.table.SetCell(rowIndex, 0,
			tview.NewTableCell(volDriver).
				SetExpansion(expand).
				SetAlign(alignment))

		// name name column
		vols.table.SetCell(rowIndex, 1,
			tview.NewTableCell(volName).
				SetExpansion(expand).
				SetAlign(alignment))

		// created at column
		vols.table.SetCell(rowIndex, 2,
			tview.NewTableCell(volCreatedAt).
				SetExpansion(expand).
				SetAlign(alignment))

		// mount point at column
		vols.table.SetCell(rowIndex, 3,
			tview.NewTableCell(volMountPoint).
				SetExpansion(expand).
				SetAlign(alignment))

		rowIndex++
	}

}
