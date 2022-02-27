package networks

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/networks"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// UpdateData retrieves networks list data
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
	fgColor := utils.Styles.PageTable.HeaderRow.FgColor
	bgColor := utils.Styles.PageTable.HeaderRow.BgColor

	for i := 0; i < len(nets.headers); i++ {
		nets.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[black::b]%s", strings.ToUpper(nets.headers[i]))).
				SetExpansion(expand).
				SetBackgroundColor(bgColor).
				SetTextColor(fgColor).
				SetAlign(tview.AlignLeft).
				SetSelectable(false))
	}
	rowIndex := 1

	nets.table.SetTitle(fmt.Sprintf("[::b]%s[%d]", strings.ToUpper(nets.title), len(netList)))
	for i := 0; i < len(netList); i++ {
		netName := netList[i][0]
		netVersion := netList[i][1]
		netPlugins := netList[i][2]

		// name name column
		nets.table.SetCell(rowIndex, 0,
			tview.NewTableCell(netName).
				SetExpansion(expand).
				SetAlign(alignment))

		// version column
		nets.table.SetCell(rowIndex, 1,
			tview.NewTableCell(netVersion).
				SetExpansion(expand).
				SetAlign(alignment))

		// plugins at column
		nets.table.SetCell(rowIndex, 2,
			tview.NewTableCell(netPlugins).
				SetExpansion(expand).
				SetAlign(alignment))

		rowIndex++
	}

}

// ClearData clears table data
func (nets *Networks) ClearData() {
	nets.table.Clear()
	expand := 1
	fgColor := utils.Styles.PageTable.HeaderRow.FgColor
	bgColor := utils.Styles.PageTable.HeaderRow.BgColor

	for i := 0; i < len(nets.headers); i++ {
		nets.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[black::b]%s", strings.ToUpper(nets.headers[i]))).
				SetExpansion(expand).
				SetBackgroundColor(bgColor).
				SetTextColor(fgColor).
				SetAlign(tview.AlignLeft).
				SetSelectable(false))
	}
	nets.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(nets.title)))
}
