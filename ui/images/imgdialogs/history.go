package imgdialogs

import (
	"fmt"
	"strings"

	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	commentCellMaxWidth = 20
)

// ImageHistoryDialog represents image history dialog primitive
type ImageHistoryDialog struct {
	*tview.Box
	layout        *tview.Flex
	imageInfo     *tview.InputField
	table         *tview.Table
	form          *tview.Form
	tableHeaders  []string
	results       [][]string
	display       bool
	cancelHandler func()
}

// NewImageHistoryDialog returns new image history dialog
func NewImageHistoryDialog() *ImageHistoryDialog {
	dialog := &ImageHistoryDialog{
		Box:       tview.NewBox(),
		imageInfo: tview.NewInputField(),
		tableHeaders: []string{
			"id", "created", "create by", "size", "comment",
		},
		display: false,
	}

	bgColor := style.DialogBgColor
	historyTableBgColor := style.DialogBgColor

	// image info field.
	imageInfoLabel := "IMAGE ID:"
	dialog.imageInfo.SetBackgroundColor(style.DialogBgColor)
	dialog.imageInfo.SetLabel("[::b]" + imageInfoLabel)
	dialog.imageInfo.SetLabelWidth(len(imageInfoLabel) + 1)
	dialog.imageInfo.SetFieldBackgroundColor(style.DialogBgColor)
	dialog.imageInfo.SetLabelStyle(tcell.StyleDefault.
		Background(style.DialogBorderColor).
		Foreground(style.DialogFgColor))

	dialog.table = tview.NewTable()
	dialog.table.SetBackgroundColor(historyTableBgColor)
	dialog.initTable()

	dialog.form = tview.NewForm().
		AddButton("Cancel", nil).
		SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(bgColor)
	dialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// table layout
	tableLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	tableLayout.SetBackgroundColor(bgColor)
	tableLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)

	tableLayout.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(dialog.imageInfo, 1, 0, false).
		AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false).
		AddItem(dialog.table, 0, 1, true),
		0, 1, true)

	tableLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)

	dialog.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	dialog.layout.SetTitle("PODMAN IMAGE HISTORY")
	dialog.layout.SetBorder(true)
	dialog.layout.SetBorderColor(style.DialogBorderColor)
	dialog.layout.SetBackgroundColor(bgColor)

	dialog.layout.AddItem(tableLayout, 0, 1, true)
	dialog.layout.AddItem(dialog.form, dialogs.DialogFormHeight, 0, true)

	return dialog
}

// Display displays this primitive
func (d *ImageHistoryDialog) Display() {
	d.display = true
}

// IsDisplay returns true if primitive is shown
func (d *ImageHistoryDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive
func (d *ImageHistoryDialog) Hide() {
	d.display = false
}

// HasFocus returns whether or not this primitive has focus
func (d *ImageHistoryDialog) HasFocus() bool {
	return d.Box.HasFocus() || d.form.HasFocus()
}

// Focus is called when this primitive receives focus
func (d *ImageHistoryDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(d.form)
}

// InputHandler returns input handler function for this primitive
func (d *ImageHistoryDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("image history dialog: event %v received", event)
		if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyEnter {
			d.cancelHandler()
			return
		}
		event = utils.ParseKeyEventKey(event)
		if event.Key() == tcell.KeyDown || event.Key() == tcell.KeyUp || event.Key() == tcell.KeyPgDn || event.Key() == tcell.KeyPgUp {
			if tableHandler := d.table.InputHandler(); tableHandler != nil {
				tableHandler(event, setFocus)
				return
			}
		}
	})
}

// SetRect set rects for this primitive.
func (d *ImageHistoryDialog) SetRect(x, y, width, height int) {
	dX := x + dialogs.DialogPadding
	dWidth := width - (2 * dialogs.DialogPadding)
	dHeight := len(d.results) + dialogs.DialogFormHeight + 5

	if dHeight > height {
		dHeight = height
	}
	tableHeight := dHeight - dialogs.DialogFormHeight - 2

	hs := ((height - dHeight) / 2)
	dY := y + hs

	d.Box.SetRect(dX, dY, dWidth, dHeight)
	//set table height size
	d.layout.ResizeItem(d.table, tableHeight, 0)
	cWidth := d.getCreatedByWidth()
	for i := 0; i < d.table.GetRowCount(); i++ {
		cell := d.table.GetCell(i, 2)
		cell.SetMaxWidth(cWidth / 2)
		d.table.SetCell(i, 2, cell)
	}

}

// Draw draws this primitive onto the screen.
func (d *ImageHistoryDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}
	d.Box.DrawForSubclass(screen, d)
	x, y, width, height := d.Box.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetCancelFunc sets form cancel button selected function
func (d *ImageHistoryDialog) SetCancelFunc(handler func()) *ImageHistoryDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 1)
	cancelButton.SetSelectedFunc(handler)
	return d
}

func (d *ImageHistoryDialog) initTable() {
	bgColor := style.TableHeaderBgColor
	fgColor := style.TableHeaderFgColor

	d.table.Clear()
	d.table.SetFixed(1, 1)
	d.table.SetSelectable(true, false)
	for i := 0; i < len(d.tableHeaders); i++ {
		d.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[%s::b]%s", style.GetColorHex(fgColor), strings.ToUpper(d.tableHeaders[i]))).
				SetExpansion(1).
				SetBackgroundColor(bgColor).
				SetTextColor(fgColor).
				SetAlign(tview.AlignLeft).
				SetSelectable(false))
	}
}

// UpdateResults updates result table
func (d *ImageHistoryDialog) UpdateResults(data [][]string) {
	d.results = data
	d.initTable()
	alignment := tview.AlignLeft
	rowIndex := 1
	expand := 0
	for i := 0; i < len(data); i++ {
		id := data[i][0]
		if len(id) > utils.IDLength {
			id = id[0:utils.IDLength]
		}
		created := data[i][1]
		createdBy := data[i][2]
		size := data[i][3]
		comment := data[i][4]
		if len(comment) > commentCellMaxWidth {
			comment = comment[0:commentCellMaxWidth]
		}

		// id column
		d.table.SetCell(rowIndex, 0,
			tview.NewTableCell(id).
				SetExpansion(expand).
				SetAlign(alignment))

		// created column
		d.table.SetCell(rowIndex, 1,
			tview.NewTableCell(created).
				SetExpansion(expand).
				SetAlign(alignment))

		// createdBy column
		d.table.SetCell(rowIndex, 2,
			tview.NewTableCell(createdBy).
				SetExpansion(1).
				SetAlign(alignment))

		// size column
		d.table.SetCell(rowIndex, 3,
			tview.NewTableCell(size).
				SetExpansion(expand).
				SetAlign(alignment))

		// comment column
		d.table.SetCell(rowIndex, 4,
			tview.NewTableCell(comment).
				SetExpansion(expand).
				SetAlign(alignment))

		rowIndex++
	}
	if len(data) > 0 {
		d.table.Select(1, 1)
		d.table.ScrollToBeginning()
	}
}

func (d *ImageHistoryDialog) SetImageInfo(id string, name string) {
	imageInfo := fmt.Sprintf("%12s (%s)", id, name)
	d.imageInfo.SetText(imageInfo)
}

func (d *ImageHistoryDialog) getCreatedByWidth() int {
	var idWidth int
	var createdWidth int
	var createdByWidth int
	var sizeWidth int
	var commentWidth int
	// get table inner rect
	_, _, width, _ := d.table.GetInnerRect()

	// get width used by other columns
	for _, row := range d.results {
		if len(row[0]) > idWidth && len(row[0]) <= utils.IDLength {
			idWidth = len(row[0])
		}
		if len(row[1]) > createdWidth {
			createdWidth = len(row[1])
		}
		if len(row[3]) > sizeWidth {
			sizeWidth = len(row[3])
		}
		if len(row[4]) > commentWidth && len(row[4]) < 40 {
			commentWidth = len(row[4])
		}
	}

	usedWidth := idWidth + createdWidth + sizeWidth + commentWidth
	createdByWidth = width - usedWidth*2 + 8
	if createdByWidth <= 0 {
		createdByWidth = 0
	}
	return createdByWidth
}
