package pods

import (
	"fmt"
	"strings"
	"time"

	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/docker/go-units"
	"github.com/rivo/tview"
)

func (pods *Pods) refresh() {
	pods.table.Clear()
	expand := 1
	alignment := tview.AlignLeft

	for i := 0; i < len(pods.headers); i++ {
		pods.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[::b]%s", strings.ToUpper(pods.headers[i]))).
				SetExpansion(expand).
				SetBackgroundColor(style.PageHeaderBgColor).
				SetTextColor(style.PageHeaderFgColor).
				SetAlign(tview.AlignLeft).
				SetSelectable(false))
	}
	rowIndex := 1

	podList := pods.getData()

	pods.table.SetTitle(fmt.Sprintf("[::b]%s[%d]", strings.ToUpper(pods.title), len(podList)))
	for i := 0; i < len(podList); i++ {
		podID := podList[i].Id
		podID = podID[0:utils.IDLength]
		podName := podList[i].Name
		podStatus := podList[i].Status
		podCreated := units.HumanDuration(time.Since(podList[i].Created)) + " ago"
		podInfraID := podList[i].InfraId

		if len(podInfraID) > utils.IDLength {
			podInfraID = podInfraID[0:utils.IDLength]
		}

		podNumCtn := fmt.Sprintf("%d", len(podList[i].Containers))

		cellTextColor := style.FgColor

		switch strings.ToLower(podStatus) {
		case "running":
			podStatus = fmt.Sprintf("[green::]%s[-::] %s", "\u25B2", podStatus)
			cellTextColor = style.RunningStatusFgColor
		case "paused":
			podStatus = fmt.Sprintf("[red::]%s[-::] %s", "\u25BC", podStatus)
			cellTextColor = style.PausedStatusFgColor
		default:
			podStatus = fmt.Sprintf("[red::]%s[-::] %s", "\u25BC", podStatus)
			cellTextColor = style.FgColor
		}

		// id column
		pods.table.SetCell(rowIndex, 0,
			tview.NewTableCell(podID).
				SetTextColor(cellTextColor).
				SetExpansion(expand).
				SetAlign(alignment))

		// name column
		pods.table.SetCell(rowIndex, 1,
			tview.NewTableCell(podName).
				SetTextColor(cellTextColor).
				SetExpansion(expand).
				SetAlign(alignment))

		// status column
		pods.table.SetCell(rowIndex, 2,
			tview.NewTableCell(podStatus).
				SetTextColor(cellTextColor).
				SetExpansion(expand).
				SetAlign(alignment))

		// created column
		pods.table.SetCell(rowIndex, 3,
			tview.NewTableCell(podCreated).
				SetTextColor(cellTextColor).
				SetExpansion(expand).
				SetAlign(alignment))

		// infra id at column
		pods.table.SetCell(rowIndex, 4,
			tview.NewTableCell(podInfraID).
				SetTextColor(cellTextColor).
				SetExpansion(expand).
				SetAlign(alignment))
		// # of container column
		pods.table.SetCell(rowIndex, 5,
			tview.NewTableCell(podNumCtn).
				SetTextColor(cellTextColor).
				SetExpansion(expand).
				SetAlign(alignment))

		rowIndex++
	}
}
