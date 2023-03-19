package volumes

import (
	"fmt"
	"strings"
	"time"

	"github.com/containers/podman-tui/ui/style"
	"github.com/docker/go-units"
	"github.com/rivo/tview"
)

const (
	volsTableDriverColIndex = 0 + iota
	volsTableNameColIndex
	volsTableCreatedAtColIndex
	volsTableMountPointColIndex
)

func (vols *Volumes) refresh() {
	vols.table.Clear()

	expand := 1
	alignment := tview.AlignLeft

	for i := 0; i < len(vols.headers); i++ {
		vols.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[::b]%s", strings.ToUpper(vols.headers[i]))).
				SetExpansion(expand).
				SetBackgroundColor(style.PageHeaderBgColor).
				SetTextColor(style.PageHeaderFgColor).
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
		vols.table.SetCell(rowIndex, volsTableDriverColIndex,
			tview.NewTableCell(volDriver).
				SetExpansion(expand).
				SetAlign(alignment))

		// name column
		vols.table.SetCell(rowIndex, volsTableNameColIndex,
			tview.NewTableCell(volName).
				SetExpansion(expand).
				SetAlign(alignment))

		// created at column
		vols.table.SetCell(rowIndex, volsTableCreatedAtColIndex,
			tview.NewTableCell(volCreatedAt).
				SetExpansion(expand).
				SetAlign(alignment))

		// mount point at column
		vols.table.SetCell(rowIndex, volsTableCreatedAtColIndex,
			tview.NewTableCell(volMountPoint).
				SetExpansion(expand).
				SetAlign(alignment))

		rowIndex++
	}
}
