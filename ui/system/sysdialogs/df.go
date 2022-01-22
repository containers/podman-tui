package sysdialogs

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/sysinfo"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	dfDialogMaxWidth = 60
)

// DfDialog is a simple dialog with disk usage result table
type DfDialog struct {
	*tview.Box
	layout       *tview.Flex
	table        *tview.Table
	form         *tview.Form
	display      bool
	tableHeaders []string
	doneHandler  func()
}

// NewDfDialog returns new DfDialog primitive
func NewDfDialog() *DfDialog {
	dialog := &DfDialog{
		Box:          tview.NewBox(),
		tableHeaders: []string{"type", "total", "active", "size", "reclaimable"},
		display:      false,
	}
	bgColor := utils.Styles.CommandDialog.BgColor
	dialog.table = tview.NewTable()
	dialog.table.SetBackgroundColor(bgColor)
	dialog.table.SetBorder(false)
	dialog.initTable()

	dialog.form = tview.NewForm().
		AddButton("Enter", nil).
		SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(bgColor)

	dialog.layout = tview.NewFlex().SetDirection(tview.FlexRow)

	dialog.layout.SetBorder(true)
	dialog.layout.SetBackgroundColor(bgColor)

	dialog.layout.AddItem(dialog.table, 4, 0, true)
	dialog.layout.AddItem(dialog.form, dialogs.DialogFormHeight, 0, true)

	return dialog
}

// SetTitle sets title for the dialog
func (d *DfDialog) SetTitle(title string) {
	d.layout.SetTitle(strings.ToUpper(title))
}

// Display displays this primitive
func (d *DfDialog) Display() {
	d.display = true
}

// IsDisplay returns true if primitive is shown
func (d *DfDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive
func (d *DfDialog) Hide() {
	d.display = false
}

// Focus is called when this primitive receives focus
func (d *DfDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(d.form)
}

// HasFocus returns true if this primitive has focus
func (d *DfDialog) HasFocus() bool {
	return d.form.HasFocus() || d.table.HasFocus()
}

// InputHandler returns input handler function for this primitive
func (d *DfDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("disk usage dialog: event %v received", event.Key())
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
func (d *DfDialog) SetRect(x, y, width, height int) {
	dX := x + dialogs.DialogPadding
	dY := y
	dWidth := width - (2 * dialogs.DialogPadding)
	if dWidth > dfDialogMaxWidth {
		dWidth = dfDialogMaxWidth
		emptySpace := (width - dWidth) / 2
		dX = x + emptySpace
	}

	dHeight := dialogs.DialogFormHeight + 4 + 2
	if height > dHeight {
		dY = y + ((height - dHeight) / 2)
		height = dHeight
	}

	d.Box.SetRect(dX, dY, dWidth, height)

}

// Draw draws this primitive onto the screen.
func (d *DfDialog) Draw(screen tcell.Screen) {

	if !d.display {
		return
	}

	d.Box.DrawForSubclass(screen, d)
	x, y, width, height := d.Box.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetDoneFunc sets form button selected function
func (d *DfDialog) SetDoneFunc(handler func()) *DfDialog {
	d.doneHandler = handler
	return d
}

func (d *DfDialog) initTable() {
	bgColor := utils.Styles.CommandDialog.HeaderRow.BgColor
	fgColor := utils.Styles.CommandDialog.HeaderRow.FgColor

	d.table.Clear()
	d.table.SetFixed(1, 1)
	d.table.SetSelectable(true, false)
	// add headers
	for i := 0; i < len(d.tableHeaders); i++ {
		d.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[%s::]%s", utils.GetColorName(fgColor), strings.ToUpper(d.tableHeaders[i]))).
				SetExpansion(1).
				SetBackgroundColor(bgColor).
				SetTextColor(fgColor).
				SetAlign(tview.AlignLeft).
				SetSelectable(false))
	}
	d.table.SetFixed(1, 1)
	d.table.SetSelectable(true, false)
}

// UpdateDiskSummary udpates disk summary table result
func (d *DfDialog) UpdateDiskSummary(sum []*sysinfo.DfSummary) {
	// add summaries
	rowIndex := 1
	for _, dfReport := range sum {
		d.table.SetCell(rowIndex, 0,
			tview.NewTableCell(dfReport.Type()).
				SetExpansion(1).
				SetAlign(tview.AlignLeft))
		d.table.SetCell(rowIndex, 1,
			tview.NewTableCell(dfReport.Total()).
				SetExpansion(1).
				SetAlign(tview.AlignLeft))
		d.table.SetCell(rowIndex, 2,
			tview.NewTableCell(dfReport.Active()).
				SetExpansion(1).
				SetAlign(tview.AlignLeft))

		d.table.SetCell(rowIndex, 3,
			tview.NewTableCell(dfReport.Size()).
				SetExpansion(1).
				SetAlign(tview.AlignLeft))

		d.table.SetCell(rowIndex, 4,
			tview.NewTableCell(dfReport.Reclaimable()).
				SetExpansion(1).
				SetAlign(tview.AlignLeft))

		rowIndex = rowIndex + 1
	}
}
