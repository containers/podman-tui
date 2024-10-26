package cntdialogs

import (
	"fmt"
	"sync"

	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/docker/go-units"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// ContainerStatsDialog implements the containers stats dialog primitive.
type ContainerStatsDialog struct {
	*tview.Box
	layout        *tview.Flex
	form          *tview.Form
	table         *tview.Table
	containerInfo *tview.InputField
	resultChan    *chan entities.ContainerStatsReport
	statsStream   *bool
	mu            sync.Mutex
	doneHandler   func()
	doneChan      chan bool
	display       bool
	maxHeight     int
	maxWidth      int
}

// NewContainerStatsDialog returns new container stats dialog.
func NewContainerStatsDialog() *ContainerStatsDialog {
	statsDialog := ContainerStatsDialog{
		Box:           tview.NewBox(),
		containerInfo: tview.NewInputField(),
		maxHeight:     14, //nolint:mnd
		maxWidth:      92, //nolint:mnd
	}

	// table
	statsDialog.table = tview.NewTable()
	statsDialog.table.SetBackgroundColor(style.BgColor)
	statsDialog.table.SetBorder(true)
	statsDialog.table.SetBorderColor(style.DialogSubBoxBorderColor)
	statsDialog.initTableUI()

	// container info text view
	cntInfoLabel := "CONTAINER ID:"

	statsDialog.containerInfo.SetBackgroundColor(style.DialogBgColor)
	statsDialog.containerInfo.SetLabel("[::b]" + cntInfoLabel)
	statsDialog.containerInfo.SetLabelWidth(len(cntInfoLabel) + 1)
	statsDialog.containerInfo.SetFieldBackgroundColor(style.DialogBgColor)
	statsDialog.containerInfo.SetLabelStyle(tcell.StyleDefault.
		Background(style.DialogBorderColor).
		Foreground(style.DialogFgColor))

	// form
	statsDialog.form = tview.NewForm().
		AddButton("Cancel", nil).
		SetButtonsAlign(tview.AlignRight)
	statsDialog.form.SetBackgroundColor(style.DialogBgColor)
	statsDialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// table layout
	statTableLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	statTableLayout.AddItem(statsDialog.table, 0, 1, false)

	statResultLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	statResultLayout.AddItem(statsDialog.containerInfo, 1, 0, true)
	statResultLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)
	statResultLayout.AddItem(statTableLayout, 0, 1, true)
	statResultLayout.SetBackgroundColor(style.BgColor)
	statResultLayout.SetBorder(false)

	// main dialog layout
	statsDialog.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	statsDialog.layout.SetBorder(true)
	statsDialog.layout.SetBorderColor(style.DialogBorderColor)
	statsDialog.layout.SetBackgroundColor(style.DialogBgColor)
	statsDialog.layout.SetTitle("PODMAN CONTAINER STATS")

	statDialogResultLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	statDialogResultLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)
	statDialogResultLayout.AddItem(statResultLayout, 0, 1, false)
	statDialogResultLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)

	statsDialog.layout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	statsDialog.layout.AddItem(statDialogResultLayout, 0, 1, true)
	statsDialog.layout.AddItem(statsDialog.form, dialogs.DialogFormHeight, 0, true)

	return &statsDialog
}

// Display displays this primitive.
func (d *ContainerStatsDialog) Display() {
	d.display = true
	d.doneChan = make(chan bool)
	d.startReportReader()
}

// IsDisplay returns true if primitive is shown.
func (d *ContainerStatsDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *ContainerStatsDialog) Hide() {
	d.display = false
	d.doneChan <- true

	d.SetContainerInfo("", "")
	d.setContainerPID(0)
	d.setContainerCPUPerc(0.0)
	d.setContainerMemPerc(0.0)
	d.setContainerMemUsage(0, 0)
	d.setContainerBlockInput(0)
	d.setContainerBlockOutput(0)
	d.setContainerNetInput(0)
	d.setContainerNetOutput(0)
	d.mu.Lock()
	defer d.mu.Unlock()

	*d.statsStream = false

	close(d.doneChan)
}

// HasFocus returns whether or not this primitive has focus.
func (d *ContainerStatsDialog) HasFocus() bool {
	if d.form.HasFocus() {
		return true
	}

	return d.Box.HasFocus()
}

// Focus is called when this primitive receives focus.
func (d *ContainerStatsDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(d.form)
}

// InputHandler  returns input handler function for this primitive.
func (d *ContainerStatsDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("container stats dialog: event %v received", event)

		if event.Key() == tcell.KeyEsc {
			d.doneHandler()

			return
		}

		if formHandler := d.form.InputHandler(); formHandler != nil {
			formHandler(event, setFocus)

			return
		}
	})
}

// Draw draws this primitive onto the screen.
func (d *ContainerStatsDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.Box.DrawForSubclass(screen, d)
	x, y, width, height := d.Box.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetRect set rects for this primitive.
func (d *ContainerStatsDialog) SetRect(x, y, width, height int) {
	dX := x + dialogs.DialogPadding
	dY := y + dialogs.DialogPadding - 1
	dWidth := width - (2 * dialogs.DialogPadding)         //nolint:mnd
	dHeight := height - (2 * (dialogs.DialogPadding - 1)) //nolint:mnd

	if dHeight > d.maxHeight {
		emptySpace := dHeight - d.maxHeight
		dY += (emptySpace / 2) //nolint:mnd
		dHeight = d.maxHeight
	}

	if dWidth > d.maxWidth {
		emptySpace := dWidth - d.maxWidth
		dX += (emptySpace / 2) //nolint:mnd
		dWidth = d.maxWidth
	}

	d.Box.SetRect(dX, dY, dWidth, dHeight)
}

// SetDoneFunc sets form cancel button selected function.
func (d *ContainerStatsDialog) SetDoneFunc(handler func()) *ContainerStatsDialog {
	d.doneHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 1)

	cancelButton.SetSelectedFunc(handler)

	return d
}

// SetContainerInfo sets container ID and name.
func (d *ContainerStatsDialog) SetContainerInfo(id string, name string) {
	info := fmt.Sprintf("%s (%s)", id, name)

	d.containerInfo.SetText(info)
}

// SetStatsChannel sets stats result read channel.
func (d *ContainerStatsDialog) SetStatsChannel(reportChan *chan entities.ContainerStatsReport) {
	d.resultChan = reportChan
}

// SetStatsStream sets stats stream state. if true it will stream the stats
// and false will stop the process.
func (d *ContainerStatsDialog) SetStatsStream(stream *bool) {
	d.statsStream = stream
}

func (d *ContainerStatsDialog) startReportReader() {
	log.Debug().Msgf("container stats dialog: starting stats reader")

	go func() {
		for {
			select {
			case result := <-*d.resultChan:
				log.Debug().Msgf("%v", result)

				if result.Error != nil {
					log.Error().Msgf("container stats error: %v", result.Error)

					continue
				}

				if len(result.Stats) > 0 {
					metric := result.Stats[0]

					var (
						netInput  uint64
						netOutput uint64
					)

					for _, net := range metric.Network {
						netInput += net.RxBytes
						netOutput += net.TxBytes
					}

					d.setContainerPID(metric.PIDs)
					d.setContainerMemPerc(metric.MemPerc)
					d.setContainerMemUsage(metric.MemUsage, metric.MemLimit)
					d.setContainerCPUPerc(metric.CPU)
					d.setContainerBlockInput(metric.BlockInput)
					d.setContainerBlockOutput(metric.BlockOutput)
					d.setContainerNetInput(netInput)
					d.setContainerNetOutput(netOutput)
				}

			case <-d.doneChan:
				log.Debug().Msgf("container stats dialog: stats reader stopped")

				return
			}
		}
	}()
}

var (
	containerMemUsageCell = tableCell{
		row: 0,
		col: 1,
	}
	containerMemPercCell = tableCell{
		row: 1,
		col: 1,
	}
	containerBlockInputCell = tableCell{
		row: 2, //nolint:mnd
		col: 1,
	}
	containerBlockOutputCell = tableCell{
		row: 3, //nolint:mnd
		col: 1,
	}
	containerPidsCell = tableCell{
		row: 0,
		col: 3, //nolint:mnd
	}
	containerCPUPercCell = tableCell{
		row: 1,
		col: 3, //nolint:mnd
	}
	containerNetInputCell = tableCell{
		row: 2, //nolint:mnd
		col: 3, //nolint:mnd
	}
	containerNetOutputCell = tableCell{
		row: 3, //nolint:mnd
		col: 3, //nolint:mnd
	}
)

type tableCell struct {
	row int
	col int
}

func (d *ContainerStatsDialog) initTableUI() {
	headerFgColor := style.TableHeaderFgColor

	// first column
	d.table.SetCell(containerMemUsageCell.row, containerMemUsageCell.col-1,
		tview.NewTableCell("mem usage/limit:").SetTextColor(headerFgColor))
	d.table.SetCell(containerMemUsageCell.row, containerMemUsageCell.col, tview.NewTableCell(""))

	d.table.SetCell(containerMemPercCell.row, containerMemPercCell.col-1,
		tview.NewTableCell("memory %:").SetTextColor(headerFgColor))
	d.table.SetCell(containerMemPercCell.row, containerMemPercCell.col, tview.NewTableCell(""))
	d.setContainerMemPerc(0.00) //nolint:mnd

	d.table.SetCell(containerBlockInputCell.row, containerBlockInputCell.col-1,
		tview.NewTableCell("block input:").SetTextColor(headerFgColor))
	d.table.SetCell(containerBlockInputCell.row, containerBlockInputCell.col, tview.NewTableCell(""))

	d.table.SetCell(containerBlockOutputCell.row, containerBlockOutputCell.col-1,
		tview.NewTableCell("block output:").SetTextColor(headerFgColor))
	d.table.SetCell(containerBlockOutputCell.row, containerBlockOutputCell.col, tview.NewTableCell(""))

	// second column
	d.table.SetCell(containerPidsCell.row, containerPidsCell.col-1,
		tview.NewTableCell("pids:").SetTextColor(headerFgColor))
	d.table.SetCell(containerPidsCell.row, containerPidsCell.col, tview.NewTableCell(""))

	d.table.SetCell(containerCPUPercCell.row, containerCPUPercCell.col-1,
		tview.NewTableCell("cpu %:").SetTextColor(headerFgColor))
	d.table.SetCell(containerCPUPercCell.row, containerCPUPercCell.col, tview.NewTableCell(""))
	d.setContainerCPUPerc(0.00) //nolint:mnd

	d.table.SetCell(containerNetInputCell.row, containerNetInputCell.col-1,
		tview.NewTableCell("net input:").SetTextColor(headerFgColor))
	d.table.SetCell(containerNetInputCell.row, containerNetInputCell.col, tview.NewTableCell(""))

	d.table.SetCell(containerNetOutputCell.row, containerNetOutputCell.col-1,
		tview.NewTableCell("net output:").SetTextColor(headerFgColor))
	d.table.SetCell(containerNetOutputCell.row, containerNetOutputCell.col, tview.NewTableCell(""))
}

func (d *ContainerStatsDialog) setContainerPID(pids uint64) {
	cntPIDS := "--"

	fgColor := style.DialogFgColor

	d.mu.Lock()
	defer d.mu.Unlock()

	if pids != 0 {
		cntPIDS = fmt.Sprintf("%d", pids) //nolint:perfsprint
	}

	d.table.GetCell(containerPidsCell.row, containerPidsCell.col).SetText(cntPIDS).SetTextColor(fgColor)
}

func (d *ContainerStatsDialog) setContainerCPUPerc(usage float64) {
	usageBar := utils.ProgressUsageString(usage)

	d.mu.Lock()
	defer d.mu.Unlock()

	d.table.GetCell(containerCPUPercCell.row, containerCPUPercCell.col).SetText(usageBar)
}

func (d *ContainerStatsDialog) setContainerMemPerc(usage float64) {
	usageBar := utils.ProgressUsageString(usage)

	d.mu.Lock()
	defer d.mu.Unlock()

	d.table.GetCell(containerMemPercCell.row, containerMemPercCell.col).SetText(usageBar)
}

func (d *ContainerStatsDialog) setContainerMemUsage(memUsage uint64, memLimit uint64) {
	usage := "-- / --"

	fgColor := style.DialogFgColor

	d.mu.Lock()
	defer d.mu.Unlock()

	if memUsage != 0 && memLimit != 0 {
		usage = fmt.Sprintf("%s / %s", units.HumanSize(float64(memUsage)), units.HumanSize(float64(memLimit)))
	}

	d.table.GetCell(containerMemUsageCell.row, containerMemUsageCell.col).SetText(usage).SetTextColor(fgColor)
}

func (d *ContainerStatsDialog) setContainerBlockInput(binput uint64) {
	blockInput := "--"

	fgColor := style.DialogFgColor

	d.mu.Lock()
	defer d.mu.Unlock()

	if binput != 0 {
		blockInput = units.HumanSize(float64(binput))
	}

	d.table.GetCell(containerBlockInputCell.row, containerBlockInputCell.col).SetText(blockInput).SetTextColor(fgColor)
}

func (d *ContainerStatsDialog) setContainerBlockOutput(boutput uint64) {
	blockOutput := "--"

	fgColor := style.DialogFgColor

	d.mu.Lock()
	defer d.mu.Unlock()

	if boutput != 0 {
		blockOutput = units.HumanSize(float64(boutput))
	}

	d.table.GetCell(containerBlockOutputCell.row, containerBlockOutputCell.col).SetText(blockOutput).SetTextColor(fgColor)
}

func (d *ContainerStatsDialog) setContainerNetInput(ninput uint64) {
	netInput := "--"

	fgColor := style.DialogFgColor

	d.mu.Lock()
	defer d.mu.Unlock()

	if ninput != 0 {
		netInput = units.HumanSize(float64(ninput))
	}

	d.table.GetCell(containerNetInputCell.row, containerNetInputCell.col).SetText(netInput).SetTextColor(fgColor)
}

func (d *ContainerStatsDialog) setContainerNetOutput(noutput uint64) {
	netOutput := "--"

	fgColor := style.DialogFgColor

	d.mu.Lock()
	defer d.mu.Unlock()

	if noutput != 0 {
		netOutput = units.HumanSize(float64(noutput))
	}

	d.table.GetCell(containerNetOutputCell.row, containerNetOutputCell.col).SetText(netOutput).SetTextColor(fgColor)
}
