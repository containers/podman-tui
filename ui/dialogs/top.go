package dialogs

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// TopDialog is a simple dialog with pod/container top result table
type TopDialog struct {
	*tview.Box
	layout       *tview.Flex
	table        *tview.Table
	form         *tview.Form
	display      bool
	tableHeaders []string
	results      [][]string
	height       int
	doneHandler  func()
}

// NewTopDialog returns new TopDialog primitive
func NewTopDialog() *TopDialog {
	dialog := &TopDialog{
		Box:          tview.NewBox(),
		tableHeaders: []string{"user", "pid", "ppid", "%cpu", "elapsed", "tty", "time", "command"},
		display:      false,
	}
	bgColor := utils.Styles.CommandDialog.BgColor
	dialog.table = tview.NewTable()
	dialog.table.SetBackgroundColor(bgColor)
	dialog.table.SetBorder(true)
	dialog.table.SetBorderColor(bgColor)
	dialog.initTable()

	dialog.height = len(dialog.results) + TableHeightOffset + DialogFormHeight + TableHeightOffset

	dialog.form = tview.NewForm().
		AddButton("Enter", nil).
		SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(bgColor)

	dialog.layout = tview.NewFlex().SetDirection(tview.FlexRow)

	dialog.layout.SetBorder(true)
	dialog.layout.SetBackgroundColor(bgColor)

	dialog.layout.AddItem(dialog.table, 1, 0, true)
	dialog.layout.AddItem(dialog.form, DialogFormHeight, 0, true)

	return dialog
}

// SetTitle sets title for the dialog
func (d *TopDialog) SetTitle(title string) {
	d.layout.SetTitle(strings.ToUpper(title))
}

// Display displays this primitive
func (d *TopDialog) Display() {
	d.display = true
}

// IsDisplay returns true if primitive is shown
func (d *TopDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive
func (d *TopDialog) Hide() {
	d.display = false
}

// Focus is called when this primitive receives focus
func (d *TopDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(d.form)
}

// HasFocus returns true if this primitive has focus
func (d *TopDialog) HasFocus() bool {
	return d.form.HasFocus() || d.table.HasFocus()
}

// InputHandler returns input handler function for this primitive
func (d *TopDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("top dialog: event %v received", event.Key())
		if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyEnter {
			d.doneHandler()
			return
		}
		if event.Key() == tcell.KeyDown || event.Key() == tcell.KeyUp || event.Key() == tcell.KeyPgDn || event.Key() == tcell.KeyPgUp {
			if tableHandler := d.table.InputHandler(); tableHandler != nil {
				tableHandler(event, setFocus)
				return
			}
		}
	})
}

// SetRect set rects for this primitive.
func (d *TopDialog) SetRect(x, y, width, height int) {
	dX := x + DialogPadding
	dWidth := width - (2 * DialogPadding)
	dHeight := d.height

	if dHeight > height {
		dHeight = height
	}

	hs := ((height - dHeight) / 2)
	dY := y + hs

	d.Box.SetRect(dX, dY, dWidth, dHeight)
	//set table height size
	d.layout.ResizeItem(d.table, dHeight-DialogFormHeight-2, 0)

	cWidth := d.getCommandWidth()
	for i := 0; i < d.table.GetRowCount(); i++ {
		cell := d.table.GetCell(i, 7)
		cell.SetMaxWidth(cWidth / 2)
		d.table.SetCell(i, 7, cell)
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

// SetDoneFunc sets form button selected function
func (d *TopDialog) SetDoneFunc(handler func()) *TopDialog {
	d.doneHandler = handler
	return d
}

func (d *TopDialog) initTable() {
	bgColor := utils.Styles.CommandDialog.HeaderRow.BgColor
	fgColor := utils.Styles.CommandDialog.HeaderRow.FgColor

	d.table.Clear()
	d.table.SetFixed(1, 1)
	d.table.SetSelectable(true, false)
	for i := 0; i < len(d.tableHeaders); i++ {
		d.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[%s::]%s", utils.GetColorName(fgColor), strings.ToUpper(d.tableHeaders[i]))).
				SetExpansion(0).
				SetBackgroundColor(bgColor).
				SetTextColor(fgColor).
				SetAlign(tview.AlignLeft).
				SetSelectable(false))
	}
}

// UpdateResults updates result table
func (d *TopDialog) UpdateResults(data [][]string) {
	d.results = data
	d.initTable()
	alignment := tview.AlignLeft
	rowIndex := 1
	expand := 1
	d.height = len(data) + TableHeightOffset + DialogFormHeight + TableHeightOffset - 2

	if len(data) < 2 {
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
		d.table.SetCell(rowIndex, 0,
			tview.NewTableCell(user).
				SetExpansion(expand).
				SetAlign(alignment))

		// pid column
		d.table.SetCell(rowIndex, 1,
			tview.NewTableCell(pid).
				SetExpansion(expand).
				SetAlign(alignment))

		// ppid column
		d.table.SetCell(rowIndex, 2,
			tview.NewTableCell(ppid).
				SetExpansion(expand).
				SetAlign(alignment))

		// cpu column
		d.table.SetCell(rowIndex, 3,
			tview.NewTableCell(cpu).
				SetExpansion(expand).
				SetAlign(alignment))

		// elapsed column
		d.table.SetCell(rowIndex, 4,
			tview.NewTableCell(elapsed).
				SetExpansion(expand).
				SetAlign(alignment))

		// tty column
		d.table.SetCell(rowIndex, 5,
			tview.NewTableCell(tty).
				SetExpansion(expand).
				SetAlign(alignment))

		// time column
		d.table.SetCell(rowIndex, 6,
			tview.NewTableCell(time).
				SetExpansion(expand).
				SetAlign(alignment))

		// command column
		d.table.SetCell(rowIndex, 7,
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
	var commandWidth int
	var usedWidth int
	// get table inner rect
	_, _, width, _ := d.table.GetInnerRect()

	// get width used by other columns
	for _, row := range d.results {
		user := len(row[0])
		pid := len(row[1])
		ppid := len(row[2])
		cpu := len(row[3])
		elapsed := len(row[4])
		tty := len(row[5])
		time := len(row[6])
		tmpUsed := user + pid + ppid + cpu + elapsed + tty + time
		if tmpUsed > usedWidth {
			usedWidth = tmpUsed
		}
	}

	commandWidth = width - usedWidth*2 + 8
	if commandWidth <= 0 {
		commandWidth = 0
	}
	return commandWidth
}
