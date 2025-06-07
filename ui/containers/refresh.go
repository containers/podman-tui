package containers

import (
	"fmt"
	"strings"
	"time"

	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/docker/go-units"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (cnt *Containers) refresh(maxWidth int) {
	imageColMaxWidth := maxWidth / 5 //nolint:mnd

	cnt.table.Clear()

	expand := 1
	alignment := tview.AlignLeft

	for i := range cnt.headers {
		cnt.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[::b]%s", strings.ToUpper(cnt.headers[i]))). //nolint:perfsprint
													SetExpansion(expand).
													SetBackgroundColor(style.PageHeaderBgColor).
													SetTextColor(style.PageHeaderFgColor).
													SetAlign(tview.AlignLeft).
													SetSelectable(false))
	}

	rowIndex := 1
	cntList := cnt.getData()

	cnt.table.SetTitle(fmt.Sprintf("[::b]%s[%d]", strings.ToUpper(cnt.title), len(cntList)))

	for i := range cntList {
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

		var cellTextColor tcell.Color

		cntShortStatus := strings.Split(strings.ToLower(cntStatus), " ")[0]

		switch cntShortStatus {
		case "up":
			cntStatus = fmt.Sprintf("[green::]%s[-::] %s", "\u25B2", cntStatus)
			cellTextColor = style.RunningStatusFgColor
		case "paused":
			cntStatus = fmt.Sprintf("[red::]%s[-::] %s", "\u25BC", cntStatus)
			cellTextColor = style.PausedStatusFgColor
		default:
			cntStatus = fmt.Sprintf("[red::]%s[-::] %s", "\u25BC", cntStatus)
			cellTextColor = style.FgColor
		}

		// id name column
		cnt.table.SetCell(rowIndex, viewContainersIDColIndex,
			tview.NewTableCell(cntID).
				SetTextColor(cellTextColor).
				SetExpansion(expand).
				SetAlign(alignment))

		// image name column
		cnt.table.SetCell(rowIndex, viewContainersImageColIndex,
			tview.NewTableCell(cntImage).
				SetMaxWidth(imageColMaxWidth).
				SetTextColor(cellTextColor).
				SetExpansion(expand).
				SetAlign(alignment))

		// pod column
		cnt.table.SetCell(rowIndex, viewContainersPodColIndex,
			tview.NewTableCell(cntPodName).
				SetTextColor(cellTextColor).
				SetExpansion(expand).
				SetAlign(alignment))

		// created at column
		cnt.table.SetCell(rowIndex, viewContainersCreatedAtColIndex,
			tview.NewTableCell(cntCreated).
				SetTextColor(cellTextColor).
				SetExpansion(expand).
				SetAlign(alignment))

		// status column
		cnt.table.SetCell(rowIndex, viewContainersStatusColIndex,
			tview.NewTableCell(cntStatus).
				SetTextColor(cellTextColor).
				SetExpansion(expand).
				SetAlign(alignment))

		// names column
		cnt.table.SetCell(rowIndex, viewContainersNamesColIndex,
			tview.NewTableCell(cntNames).
				SetTextColor(cellTextColor).
				SetExpansion(expand).
				SetAlign(alignment))

		// ports column
		cnt.table.SetCell(rowIndex, viewContainersPortsColIndex,
			tview.NewTableCell(cntPorts).
				SetTextColor(cellTextColor).
				SetExpansion(expand).
				SetAlign(alignment))

		rowIndex++
	}
}
