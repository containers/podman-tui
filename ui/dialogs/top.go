package dialogs

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	viewTopUserColIndex = 0 + iota
	viewTopPIDColIndex
	viewTopPPIDColIndex
	viewTopCPUColIndex
	viewTopElapsedColIndex
	viewTopTTYColIndex
	viewTopTimeColIndex
	viewTopCommandColIndex
)

// TopDialog is a simple dialog with pod/container top result table.
type TopDialog struct {
	*tview.Box
	layout        *tview.Flex
	table         *tview.Table
	info          *tview.InputField
	form          *tview.Form
	display       bool
	tableHeaders  []string
	results       [][]string
	cancelHandler func()
}

type topInfo int

const (
	// top dialog header label.
	TopPodInfo topInfo = 0 + iota
	TopContainerInfo
)

// NewTopDialog returns new TopDialog primitive.
func NewTopDialog() *TopDialog {
	dialog := &TopDialog{
		Box:          tview.NewBox(),
		info:         tview.NewInputField(),
		tableHeaders: []string{"user", "pid", "ppid", "%cpu", "elapsed", "tty", "time", "command"},
		display:      false,
	}
	dialog.table = tview.NewTable()
	dialog.table.SetBackgroundColor(style.DialogBgColor)
	dialog.table.SetBorder(true)
	dialog.table.SetBorderColor(style.DialogSubBoxBorderColor)
	dialog.initTable()

	dialog.info.SetBackgroundColor(style.DialogBgColor)
	dialog.info.SetFieldBackgroundColor(style.DialogBgColor)
	dialog.info.SetLabelStyle(tcell.StyleDefault.
		Background(style.DialogBorderColor).
		Foreground(style.DialogFgColor))

	dialog.form = tview.NewForm().
		AddButton("Cancel", nil).
		SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(style.DialogBgColor)
	dialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// table layout
	tableLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	tableLayout.SetBackgroundColor(style.DialogBgColor)
	tableLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)
	tableLayout.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(dialog.info, 1, 0, false).
		AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false).
		AddItem(dialog.table, 0, 1, true), 0, 1, true)
	tableLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)

	dialog.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	dialog.layout.SetBorder(true)
	dialog.layout.SetBorderColor(style.DialogBorderColor)
	dialog.layout.SetBackgroundColor(style.DialogBgColor)
	dialog.layout.AddItem(tableLayout, 0, 1, true)
	dialog.layout.AddItem(dialog.form, DialogFormHeight, 0, true)

	return dialog
}

// SetTitle sets title for the dialog.
func (d *TopDialog) SetTitle(title string) {
	d.layout.SetTitle(strings.ToUpper(title))
}

// Display displays this primitive.
func (d *TopDialog) Display() {
	d.display = true
}

// IsDisplay returns true if primitive is shown.
func (d *TopDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *TopDialog) Hide() {
	d.display = false
	d.info.SetText("")
}

// Focus is called when this primitive receives focus.
func (d *TopDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(d.form)
}

// HasFocus returns true if this primitive has focus.
func (d *TopDialog) HasFocus() bool {
	return d.form.HasFocus() || d.table.HasFocus()
}

// InputHandler returns input handler function for this primitive.
func (d *TopDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("top dialog: event %v received", event)

		if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyEnter {
			d.cancelHandler()

			return
		}

		// scroll between top items
		if tableHandler := d.table.InputHandler(); tableHandler != nil {
			tableHandler(event, setFocus)

			return
		}
	})
}

// SetRect set rects for this primitive.
func (d *TopDialog) SetRect(x, y, width, height int) {
	dX := x + DialogPadding
	dWidth := width - (2 * DialogPadding)            //nolint:mnd
	dHeight := len(d.results) + DialogFormHeight + 6 //nolint:mnd

	if dHeight > height {
		dHeight = height
	}

	tableHeight := dHeight - DialogFormHeight - 2 //nolint:mnd

	hs := ((height - dHeight) / 2) //nolint:mnd
	dY := y + hs

	d.Box.SetRect(dX, dY, dWidth, dHeight)
	// set table height size
	d.layout.ResizeItem(d.table, tableHeight, 0)

	cWidth := d.getCommandWidth()

	for i := range d.table.GetRowCount() {
		cell := d.table.GetCell(i, viewTopCommandColIndex)
		cell.SetMaxWidth(cWidth / 2) //nolint:mnd
		d.table.SetCell(i, viewTopCommandColIndex, cell)
	}
}

// Draw draws this primitive onto the screen.
func (d *TopDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.Box.DrawForSubclass(screen, d)
	x, y, width, height := d.Box.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetCancelFunc sets form button selected function.
func (d *TopDialog) SetCancelFunc(handler func()) *TopDialog {
	d.cancelHandler = handler

	return d
}

func (d *TopDialog) initTable() {
	bgColor := style.TableHeaderBgColor
	fgColor := style.TableHeaderFgColor

	d.table.Clear()
	d.table.SetFixed(1, 1)
	d.table.SetSelectable(true, false)

	for i := range d.tableHeaders {
		d.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[%s::b]%s", style.GetColorHex(fgColor), strings.ToUpper(d.tableHeaders[i]))).
				SetExpansion(0).
				SetBackgroundColor(bgColor).
				SetTextColor(fgColor).
				SetAlign(tview.AlignLeft).
				SetSelectable(false))
	}
}

// UpdateResults updates result table.
func (d *TopDialog) UpdateResults(infoType topInfo, id string, name string, data [][]string) {
	headerInfo := "CONTAINER ID:"
	if infoType == TopPodInfo {
		headerInfo = "POD ID:"
	}

	d.info.SetLabel("[b::b]" + headerInfo)
	d.info.SetLabelWidth(len(headerInfo) + 1)

	infoMessage := fmt.Sprintf("%12s (%s)", id, name)
	d.info.SetText(infoMessage)

	d.results = data
	d.initTable()

	alignment := tview.AlignLeft
	rowIndex := 1
	expand := 1

	if len(data) < 2 { //nolint:mnd
		return
	}

	for i := 1; i < len(data); i++ {
		user := data[i][0]
		pid := data[i][1]
		ppid := data[i][2]
		cpu := data[i][3]
		elapsed := data[i][4]
		tty := data[i][5]
		time := data[i][6]
		command := data[i][7]

		// user column
		d.table.SetCell(rowIndex, viewTopUserColIndex,
			tview.NewTableCell(user).
				SetExpansion(expand).
				SetAlign(alignment))

		// pid column
		d.table.SetCell(rowIndex, viewTopPIDColIndex,
			tview.NewTableCell(pid).
				SetExpansion(expand).
				SetAlign(alignment))

		// ppid column
		d.table.SetCell(rowIndex, viewTopPPIDColIndex,
			tview.NewTableCell(ppid).
				SetExpansion(expand).
				SetAlign(alignment))

		// cpu column
		d.table.SetCell(rowIndex, viewTopCPUColIndex,
			tview.NewTableCell(cpu).
				SetExpansion(expand).
				SetAlign(alignment))

		// elapsed column
		d.table.SetCell(rowIndex, viewTopElapsedColIndex,
			tview.NewTableCell(elapsed).
				SetExpansion(expand).
				SetAlign(alignment))

		// tty column
		d.table.SetCell(rowIndex, viewTopTTYColIndex,
			tview.NewTableCell(tty).
				SetExpansion(expand).
				SetAlign(alignment))

		// time column
		d.table.SetCell(rowIndex, viewTopTimeColIndex,
			tview.NewTableCell(time).
				SetExpansion(expand).
				SetAlign(alignment))

		// command column
		d.table.SetCell(rowIndex, viewTopCommandColIndex,
			tview.NewTableCell(command).
				SetExpansion(1).
				SetAlign(alignment))

		rowIndex++
	}

	if len(data) > 0 {
		d.table.Select(1, 1)
		d.table.ScrollToBeginning()
	}
}

func (d *TopDialog) getCommandWidth() int {
	var (
		commandWidth int
		usedWidth    int
	)

	// get table inner rect
	_, _, width, _ := d.table.GetInnerRect() //nolint:dogsled

	// get width used by other columns
	for _, row := range d.results {
		user := len(row[viewTopUserColIndex])
		pid := len(row[viewTopPIDColIndex])
		ppid := len(row[viewTopPPIDColIndex])
		cpu := len(row[viewTopCPUColIndex])
		elapsed := len(row[viewTopElapsedColIndex])
		tty := len(row[viewTopTTYColIndex])
		time := len(row[viewTopTimeColIndex])
		tmpUsed := user + pid + ppid + cpu + elapsed + tty + time

		if tmpUsed > usedWidth {
			usedWidth = tmpUsed
		}
	}

	commandWidth = width - usedWidth*2 + 8 //nolint:mnd
	if commandWidth <= 0 {
		commandWidth = 0
	}

	return commandWidth
}
