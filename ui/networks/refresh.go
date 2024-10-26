package networks

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/networks"
	"github.com/containers/podman-tui/ui/style"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// UpdateData retrieves networks list data.
func (nets *Networks) UpdateData() {
	netList, err := networks.List()
	if err != nil {
		log.Error().Msgf("view: networks update %v", err)
		nets.errorDialog.SetText(fmt.Sprintf("%v", err))
		nets.errorDialog.Display()
	}

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

	nets.table.SetTitle(fmt.Sprintf("[::b]%s[%d]", strings.ToUpper(nets.title), len(netList)))

	for i := range netList {
		netID := netList[i][0]
		netName := netList[i][1]
		netDriver := netList[i][2]

		// name column
		nets.table.SetCell(rowIndex, viewNetworkNameColIndex,
			tview.NewTableCell(netID[:12]).
				SetExpansion(expand).
				SetAlign(alignment))

		// version column
		nets.table.SetCell(rowIndex, viewNetworkVersionColIndex,
			tview.NewTableCell(netName).
				SetExpansion(expand).
				SetAlign(alignment))

		// plugins at column
		nets.table.SetCell(rowIndex, viewNetworkPluginColIndex,
			tview.NewTableCell(netDriver).
				SetExpansion(expand).
				SetAlign(alignment))

		rowIndex++
	}
}

// ClearData clears table data.
func (nets *Networks) ClearData() {
	nets.table.Clear()

	expand := 1
	fgColor := style.PageHeaderFgColor
	bgColor := style.PageHeaderBgColor

	for i := range nets.headers {
		nets.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[::b]%s", strings.ToUpper(nets.headers[i]))). //nolint:perfsprint
													SetExpansion(expand).
													SetBackgroundColor(bgColor).
													SetTextColor(fgColor).
													SetAlign(tview.AlignLeft).
													SetSelectable(false))
	}

	nets.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(nets.title)))
}
