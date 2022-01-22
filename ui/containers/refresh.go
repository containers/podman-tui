package containers

import (
	"fmt"
	"strings"
	"time"

	"github.com/docker/go-units"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/rivo/tview"
)

func (cnt *Containers) refresh() {
	cnt.table.Clear()
	expand := 1
	alignment := tview.AlignLeft
	fgColor := utils.Styles.PageTable.HeaderRow.FgColor
	bgColor := utils.Styles.PageTable.HeaderRow.BgColor
	for i := 0; i < len(cnt.headers); i++ {
		cnt.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[black::b]%s", strings.ToUpper(cnt.headers[i]))).
				SetExpansion(expand).
				SetBackgroundColor(bgColor).
				SetTextColor(fgColor).
				SetAlign(tview.AlignLeft).
				SetSelectable(false))
	}
	rowIndex := 1

	cntList := cnt.getData()
	cnt.table.SetTitle(fmt.Sprintf("[::b]%s[%d]", strings.ToUpper(cnt.title), len(cntList)))
	for i := 0; i < len(cntList); i++ {
		cntID := cntList[i].ID
		if len(cntID) > utils.IDLength {
			cntID = cntID[:utils.IDLength]
		}
		cntImage := cntList[i].Image
		cntPodName := cntList[i].PodName
		cntCreated := units.HumanDuration(time.Since(cntList[i].Created)) + " ago"
		cntStatus := conReporter{cntList[i]}.status()
		cntPorts := conReporter{cntList[i]}.ports()
		cntNames := conReporter{cntList[i]}.names()

		//cntStatus = strings.Split(cntStatus, " ")[0]

		if strings.Contains(strings.ToLower(cntStatus), "up") {
			cntStatus = fmt.Sprintf("[green::]%s[-::] %s", "\u25B2", cntStatus)
		} else {
			cntStatus = fmt.Sprintf("[red::]%s[-::] %s", "\u25BC", cntStatus)
		}

		// id name column
		cnt.table.SetCell(rowIndex, 0,
			tview.NewTableCell(cntID).
				SetExpansion(expand).
				SetAlign(alignment))

		// image name column
		cnt.table.SetCell(rowIndex, 1,
			tview.NewTableCell(cntImage).
				SetExpansion(expand).
				SetAlign(alignment))

		// pod column
		cnt.table.SetCell(rowIndex, 2,
			tview.NewTableCell(cntPodName).
				SetExpansion(expand).
				SetAlign(alignment))

		// created at column
		cnt.table.SetCell(rowIndex, 3,
			tview.NewTableCell(cntCreated).
				SetExpansion(expand).
				SetAlign(alignment))

		// status column
		cnt.table.SetCell(rowIndex, 4,
			tview.NewTableCell(cntStatus).
				SetExpansion(expand).
				SetAlign(alignment))

		// names column
		cnt.table.SetCell(rowIndex, 5,
			tview.NewTableCell(cntNames).
				SetExpansion(expand).
				SetAlign(alignment))

		// ports column
		cnt.table.SetCell(rowIndex, 6,
			tview.NewTableCell(cntPorts).
				SetExpansion(expand).
				SetAlign(alignment))

		rowIndex++
	}

}
