package system

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/ui/style"
	"github.com/rivo/tview"
)

func (sys *System) refresh() {
	connections := sys.getConnectionsData()
	sys.connTable.Clear()
	sys.updateConnTableTitle(len(connections))

	expand := 1
	alignment := tview.AlignLeft
	defaultAlignment := tview.AlignCenter

	for i := 0; i < len(sys.connTableHeaders); i++ {
		headerAlignment := alignment
		if sys.connTableHeaders[i] == "default" {
			headerAlignment = defaultAlignment
		}

		header := fmt.Sprintf("[::b]%s", strings.ToUpper(sys.connTableHeaders[i]))
		sys.connTable.SetCell(0, i,
			tview.NewTableCell(header).
				SetExpansion(1).
				SetBackgroundColor(style.PageHeaderBgColor).
				SetTextColor(style.PageHeaderFgColor).
				SetAlign(headerAlignment).
				SetSelectable(false))
	}

	rowIndex := 1

	for i := 0; i < len(connections); i++ {
		isDefault := ""
		conn := connections[i]
		status := connectionItemStatus{conn.Status}.StatusString()

		if conn.Default {
			isDefault = style.HeavyGreenCheckMark
		}

		// name column
		sys.connTable.SetCell(rowIndex, 0,
			tview.NewTableCell(conn.Name).
				SetExpansion(expand).
				SetAlign(alignment))

		// default column
		sys.connTable.SetCell(rowIndex, 1,
			tview.NewTableCell(isDefault).
				SetExpansion(expand).
				SetAlign(defaultAlignment))

		// status column
		sys.connTable.SetCell(rowIndex, 2, //nolint:gomnd
			tview.NewTableCell(status).
				SetExpansion(expand).
				SetAlign(alignment))

		// uri column
		sys.connTable.SetCell(rowIndex, 3, //nolint:gomnd
			tview.NewTableCell(conn.URI).
				SetExpansion(expand).
				SetAlign(alignment))

		// identity column
		sys.connTable.SetCell(rowIndex, 4, //nolint:gomnd
			tview.NewTableCell(conn.Identity).
				SetExpansion(expand).
				SetAlign(alignment))

		rowIndex++
	}
}

func (sys *System) updateConnTableTitle(total int) {
	title := fmt.Sprintf("[::b]SYSTEM CONNECTIONS[%d]", total)
	sys.connTable.SetTitle(title)
}
