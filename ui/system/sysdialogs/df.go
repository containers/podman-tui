package sysdialogs

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/pdcs/sysinfo"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	dfDialogMaxWidth = 60
)

// DfDialog is a simple dialog with disk usage result table.
type DfDialog struct {
	*tview.Box
	layout        *tview.Flex
	serviceName   *tview.InputField
	table         *tview.Table
	form          *tview.Form
	display       bool
	tableHeaders  []string
	cancelHandler func()
}

// NewDfDialog returns new DfDialog primitive.
func NewDfDialog() *DfDialog {
	dialog := &DfDialog{
		Box:          tview.NewBox(),
		serviceName:  tview.NewInputField(),
		tableHeaders: []string{"type", "total", "active", "size", "reclaimable"},
		display:      false,
	}

	// service name input field
	serviceNameLabel := "SERVICE NAME:"

	dialog.serviceName.SetBackgroundColor(style.DialogBgColor)
	dialog.serviceName.SetLabel("[::b]" + serviceNameLabel)
	dialog.serviceName.SetLabelWidth(len(serviceNameLabel) + 1)
	dialog.serviceName.SetFieldBackgroundColor(style.DialogBgColor)
	dialog.serviceName.SetLabelStyle(tcell.StyleDefault.
		Background(style.DialogBorderColor).
		Foreground(style.DialogFgColor))

	// disk usage table
	dialog.table = tview.NewTable()

	dialog.table.SetBackgroundColor(style.DialogBgColor)
	dialog.table.SetBorder(true)
	dialog.table.SetBorderColor(style.DialogSubBoxBorderColor)
	dialog.initTable()

	dialog.form = tview.NewForm().
		AddButton("Cancel", nil).
		SetButtonsAlign(tview.AlignRight)

	dialog.form.SetBackgroundColor(style.DialogBgColor)
	dialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	dialog.layout = tview.NewFlex().SetDirection(tview.FlexRow)

	dialog.layout.SetBorder(true)
	dialog.layout.SetBorderColor(style.DialogBorderColor)
	dialog.layout.SetBackgroundColor(style.DialogBgColor)

	tableLayout := tview.NewFlex().SetDirection(tview.FlexColumn)

	tableLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	tableLayout.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false).
		AddItem(dialog.serviceName, 1, 0, false).
		AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false).
		AddItem(dialog.table, 0, 1, true), 0, 1, true)
	tableLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)

	dialog.layout.AddItem(tableLayout, 9, 0, true) //nolint:mnd
	dialog.layout.AddItem(dialog.form, dialogs.DialogFormHeight, 0, true)

	return dialog
}

// SetServiceName sets event dialog service (connection) name.
func (d *DfDialog) SetServiceName(name string) {
	d.layout.SetTitle("SYSTEM DISK USAGE")
	d.serviceName.SetText(name)
}

// Display displays this primitive.
func (d *DfDialog) Display() {
	d.display = true
}

// IsDisplay returns true if primitive is shown.
func (d *DfDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *DfDialog) Hide() {
	d.display = false
}

// Focus is called when this primitive receives focus.
func (d *DfDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(d.form)
}

// HasFocus returns true if this primitive has focus.
func (d *DfDialog) HasFocus() bool {
	return d.form.HasFocus() || d.table.HasFocus()
}

// InputHandler returns input handler function for this primitive.
func (d *DfDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("disk usage dialog: event %v received", event)

		if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyEnter {
			d.cancelHandler()

			return
		}

		// scroll between df items
		if tableHandler := d.table.InputHandler(); tableHandler != nil {
			tableHandler(event, setFocus)

			return
		}
	})
}

// SetRect set rects for this primitive.
func (d *DfDialog) SetRect(x, y, width, height int) {
	dX := x + dialogs.DialogPadding
	dY := y
	dWidth := width - (2 * dialogs.DialogPadding) //nolint:mnd

	if dWidth > dfDialogMaxWidth {
		dWidth = dfDialogMaxWidth
		emptySpace := (width - dWidth) / 2 //nolint:mnd
		dX = x + emptySpace
	}

	dHeight := dialogs.DialogFormHeight + 11 //nolint:mnd
	if height > dHeight {
		dY = y + ((height - dHeight) / 2) //nolint:mnd
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

// SetCancelFunc sets form cancel button selected function.
func (d *DfDialog) SetCancelFunc(handler func()) *DfDialog {
	d.cancelHandler = handler

	return d
}

func (d *DfDialog) initTable() {
	bgColor := style.TableHeaderBgColor
	fgColor := style.TableHeaderFgColor

	d.table.Clear()
	d.table.SetFixed(1, 1)
	d.table.SetSelectable(true, false)

	// add headers
	for i := range d.tableHeaders {
		d.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[%s::b]%s", style.GetColorHex(fgColor), strings.ToUpper(d.tableHeaders[i]))).
				SetExpansion(1).
				SetBackgroundColor(bgColor).
				SetTextColor(fgColor).
				SetAlign(tview.AlignLeft).
				SetSelectable(false))
	}

	d.table.SetFixed(1, 1)
	d.table.SetSelectable(true, false)
}

// UpdateDiskSummary updates disk summary table result.
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
		d.table.SetCell(rowIndex, 2, //nolint:mnd
			tview.NewTableCell(dfReport.Active()).
				SetExpansion(1).
				SetAlign(tview.AlignLeft))

		d.table.SetCell(rowIndex, 3, //nolint:mnd
			tview.NewTableCell(dfReport.Size()).
				SetExpansion(1).
				SetAlign(tview.AlignLeft))

		d.table.SetCell(rowIndex, 4, //nolint:mnd
			tview.NewTableCell(dfReport.Reclaimable()).
				SetExpansion(1).
				SetAlign(tview.AlignLeft))

		rowIndex++
	}
}
