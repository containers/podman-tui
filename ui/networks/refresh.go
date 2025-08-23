package networks

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/ui/style"
	"github.com/rivo/tview"
)

func (nets *Networks) refresh(_ int) {
	nets.table.Clear()

	expand := 1
	alignment := tview.AlignLeft

	for i := range nets.headers {
		nets.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[::b]%s", strings.ToUpper(nets.headers[i]))). //nolint:perfsprint
													SetExpansion(expand).
													SetBackgroundColor(style.PageHeaderBgColor).
													SetTextColor(style.PageHeaderFgColor).
													SetAlign(tview.AlignLeft).
													SetSelectable(false))
	}

	rowIndex := 1
	netList := nets.getData()

	nets.table.SetTitle(fmt.Sprintf("[::b]%s[%d]", strings.ToUpper(nets.title), len(netList)))

	for _, net := range netList {
		// name column
		nets.table.SetCell(rowIndex, viewNetworkNameColIndex,
			tview.NewTableCell(net.ID[:12]).
				SetExpansion(expand).
				SetAlign(alignment))

		// version column
		nets.table.SetCell(rowIndex, viewNetworkVersionColIndex,
			tview.NewTableCell(net.Name).
				SetExpansion(expand).
				SetAlign(alignment))

		// plugins at column
		nets.table.SetCell(rowIndex, viewNetworkPluginColIndex,
			tview.NewTableCell(net.Driver).
				SetExpansion(expand).
				SetAlign(alignment))

		rowIndex++
	}
}
