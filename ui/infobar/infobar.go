package infobar

import (
	"fmt"

	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"

	"github.com/rivo/tview"
)

// InfoBarViewHeight info bar height
const (
	InfoBarViewHeight = 5
	connectionCellRow = 0
	hostnameCellRow   = 1
	osCellRow         = 2
	memCellRow        = 3
	swapCellRow       = 4
	connOK            = "\u2705"
	connERR           = "\u274C"
)

// InfoBar implements the info bar primitive
type InfoBar struct {
	*tview.Box
	table  *tview.Table
	title  string
	connOK bool
}

// NewInfoBar returns info bar view
func NewInfoBar() *InfoBar {
	table := tview.NewTable()
	headerColor := utils.GetColorName(utils.Styles.InfoBar.ItemFgColor)
	emptyCell := func() *tview.TableCell {
		return tview.NewTableCell("")
	}

	// empty column
	for i := 0; i < 5; i++ {
		table.SetCell(i, 0, emptyCell())
	}

	// valueColor := Styles.InfoBar.ValueFgColor
	table.SetCell(connectionCellRow, 1, tview.NewTableCell(fmt.Sprintf("[%s::]%s", headerColor, "Connection:")))
	table.SetCell(connectionCellRow, 2, emptyCell())

	table.SetCell(hostnameCellRow, 1, tview.NewTableCell(fmt.Sprintf("[%s::]%s", headerColor, "Hostname:")))
	table.SetCell(hostnameCellRow, 2, emptyCell())

	table.SetCell(osCellRow, 1, tview.NewTableCell(fmt.Sprintf("[%s::]%s", headerColor, "OS type:")))
	table.SetCell(osCellRow, 2, emptyCell())

	table.SetCell(memCellRow, 1, tview.NewTableCell(fmt.Sprintf("[%s::]%s", headerColor, "Memory usage:")))
	table.SetCell(memCellRow, 2, tview.NewTableCell(utils.ProgressUsageString(0.00)))

	table.SetCell(swapCellRow, 1, tview.NewTableCell(fmt.Sprintf("[%s::]%s", headerColor, "Swap usage:")))
	table.SetCell(swapCellRow, 2, tview.NewTableCell(utils.ProgressUsageString(0.00)))

	// empty column
	for i := 0; i < 5; i++ {
		table.SetCell(i, 3, emptyCell())
	}

	table.SetCell(connectionCellRow, 4, tview.NewTableCell(fmt.Sprintf("[%s::]%s", headerColor, "Kernel version:")))
	table.SetCell(connectionCellRow, 5, emptyCell())

	table.SetCell(hostnameCellRow, 4, tview.NewTableCell(fmt.Sprintf("[%s::]%s", headerColor, "API version:")))
	table.SetCell(hostnameCellRow, 5, emptyCell())

	table.SetCell(osCellRow, 4, tview.NewTableCell(fmt.Sprintf("[%s::]%s", headerColor, "OCI runtime:")))
	table.SetCell(osCellRow, 5, emptyCell())

	table.SetCell(memCellRow, 4, tview.NewTableCell(fmt.Sprintf("[%s::]%s", headerColor, "Conmon version:")))
	table.SetCell(memCellRow, 5, emptyCell())

	table.SetCell(swapCellRow, 4, tview.NewTableCell(fmt.Sprintf("[%s::]%s", headerColor, "Buildah version:")))
	table.SetCell(swapCellRow, 5, emptyCell())

	// infobar
	infoBar := &InfoBar{
		Box:    tview.NewBox(),
		title:  "infobar",
		table:  table,
		connOK: false,
	}
	return infoBar
}

// UpdatePodmanInfo updates api, conmon and buildah version values
func (info *InfoBar) UpdatePodmanInfo(apiVersion string, ociRuntime string, conmonVersion string, buildahVersion string) {
	info.table.GetCell(hostnameCellRow, 5).SetText(apiVersion)
	info.table.GetCell(osCellRow, 5).SetText(ociRuntime)
	info.table.GetCell(memCellRow, 5).SetText(conmonVersion)
	info.table.GetCell(swapCellRow, 5).SetText(buildahVersion)
}

// UpdateBasicInfo updates hostname, kernel and os type values
func (info *InfoBar) UpdateBasicInfo(hostname string, kernel string, ostype string) {
	info.table.GetCell(hostnameCellRow, 2).SetText(hostname)
	info.table.GetCell(osCellRow, 2).SetText(ostype)
	info.table.GetCell(connectionCellRow, 5).SetText(kernel)
}

// UpdateSystemUsageInfo updates memory and swap values
func (info *InfoBar) UpdateSystemUsageInfo(memUsage float64, swapUsage float64) {
	memUsageText := utils.ProgressUsageString(memUsage)
	swapUsageText := utils.ProgressUsageString(swapUsage)
	info.table.GetCell(memCellRow, 2).SetText(memUsageText)
	info.table.GetCell(swapCellRow, 2).SetText(swapUsageText)
}

// UpdateConnStatus updates connection status value
func (info *InfoBar) UpdateConnStatus(status bool) {
	info.connOK = status
	connStatus := ""
	if info.connOK {
		connStatus = fmt.Sprintf("%s STATUS_OK", connOK)

	} else {
		connStatus = fmt.Sprintf("%s STATUS_ERR", connERR)
	}
	info.table.GetCell(connectionCellRow, 2).SetText(connStatus)
}

// Draw draws this primitive onto the screen.
func (info *InfoBar) Draw(screen tcell.Screen) {
	info.Box.DrawForSubclass(screen, info)
	info.Box.SetBorder(false)
	x, y, width, height := info.GetInnerRect()
	info.table.SetRect(x, y, width, height)
	info.table.SetBorder(false)
	info.table.Draw(screen)
}
