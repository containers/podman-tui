package dialogs

import (
	"fmt"

	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	cmdWidthOffset = 6
)

// CommandDialog is a commands list dialog
type CommandDialog struct {
	*tview.Box
	layout        *tview.Flex
	table         *tview.Table
	form          *tview.Form
	display       bool
	options       [][]string
	width         int
	height        int
	cancelHandler func()
	selectHandler func()
}

// NewCommandDialog returns a command list primitive.
func NewCommandDialog(options [][]string) *CommandDialog {

	form := tview.NewForm().
		AddButton("Cancel", nil).
		AddButton("Enter", nil).
		SetButtonsAlign(tview.AlignRight)

	form.SetBackgroundColor(utils.Styles.CommandDialog.BgColor)

	bgColor := utils.Styles.CommandDialog.HeaderRow.BgColor
	fgColor := utils.Styles.CommandDialog.HeaderRow.FgColor
	cmdsTable := tview.NewTable()

	cmdWidth := 0
	// command table header
	cmdsTable.SetCell(0, 0,
		tview.NewTableCell(fmt.Sprintf("[%s::b]COMMAND", utils.GetColorName(fgColor))).
			SetExpansion(1).
			SetBackgroundColor(bgColor).
			SetTextColor(fgColor).
			SetAlign(tview.AlignLeft).
			SetSelectable(false))
	cmdsTable.SetCell(0, 1,
		tview.NewTableCell(fmt.Sprintf("[%s::b]DESCRIPTION", utils.GetColorName(fgColor))).
			SetExpansion(1).
			SetBackgroundColor(bgColor).
			SetTextColor(fgColor).
			SetAlign(tview.AlignCenter).
			SetSelectable(false))

	// command table items
	bgColor = utils.Styles.CommandDialog.BgColor
	fgColor = utils.Styles.CommandDialog.FgColor
	col1Width := 0
	col2Width := 0
	for i := 0; i < len(options); i++ {
		cmdsTable.SetCell(i+1, 0,
			tview.NewTableCell(options[i][0]).
				SetAlign(tview.AlignLeft).
				SetSelectable(true).SetTextColor(fgColor))
		cmdsTable.SetCell(i+1, 1,
			tview.NewTableCell(options[i][1]).
				SetAlign(tview.AlignLeft).
				SetSelectable(true).SetTextColor(fgColor))

		if len(options[i][0]) > col1Width {
			col1Width = len(options[i][0])
		}
		if len(options[i][1]) > col2Width {
			col2Width = len(options[i][1])
		}
	}
	cmdWidth = col1Width + col2Width + 2
	cmdsTable.SetFixed(1, 1)
	cmdsTable.SetSelectable(true, false)
	cmdsTable.SetBackgroundColor(bgColor)

	cmdLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	cmdLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	cmdLayout.AddItem(cmdsTable, 0, 1, true)
	cmdLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)

	// layout
	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	layout.AddItem(cmdLayout, 0, 1, true)
	layout.AddItem(form, DialogFormHeight, 0, true)
	layout.SetBorder(true)
	layout.SetBackgroundColor(bgColor)

	// returns the command primitive
	return &CommandDialog{
		Box:     tview.NewBox().SetBorder(false),
		layout:  layout,
		table:   cmdsTable,
		form:    form,
		display: false,
		options: options,
		width:   cmdWidth + cmdWidthOffset,
		height:  len(options) + TableHeightOffset + DialogFormHeight,
	}
}

// GetSelectedItem returns selected row item
func (cmd *CommandDialog) GetSelectedItem() string {
	row, _ := cmd.table.GetSelection()
	if row >= 0 {
		return cmd.options[row-1][0]
	}
	return ""
}

// GetCommandCount returns number of commands
func (cmd *CommandDialog) GetCommandCount() int {
	return cmd.table.GetRowCount()
}

// Display displays this primitive
func (cmd *CommandDialog) Display() {
	cmd.table.Select(1, 0)
	cmd.form.SetFocus(1)
	cmd.display = true
}

// IsDisplay returns true if primitive is shown
func (cmd *CommandDialog) IsDisplay() bool {
	return cmd.display
}

// Hide stops displaying this primitive
func (cmd *CommandDialog) Hide() {
	cmd.display = false
}

// HasFocus returns whether or not this primitive has focus
func (cmd *CommandDialog) HasFocus() bool {
	return cmd.form.HasFocus()
}

// Focus is called when this primitive receives focus
func (cmd *CommandDialog) Focus(delegate func(p tview.Primitive)) {
	delegate(cmd.form)
}

// InputHandler returns input handler function for this primitive
func (cmd *CommandDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return cmd.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("commands dialog: event %v received", event)
		if event.Key() == tcell.KeyEsc {
			cmd.cancelHandler()
			return
		}
		// select command
		if event.Key() == tcell.KeyEnter || event.Key() == tcell.KeyTab {
			if formHandler := cmd.form.InputHandler(); formHandler != nil {
				formHandler(event, setFocus)
				return
			}
		}
		// scroll between command items
		if tableHandler := cmd.table.InputHandler(); tableHandler != nil {
			tableHandler(event, setFocus)
			return
		}
	})
}

// SetSelectedFunc sets form enter button selected function
func (cmd *CommandDialog) SetSelectedFunc(handler func()) *CommandDialog {
	cmd.selectHandler = handler
	enterButton := cmd.form.GetButton(cmd.form.GetButtonCount() - 1)
	enterButton.SetSelectedFunc(handler)
	return cmd
}

// SetCancelFunc sets form cancel button selected function
func (cmd *CommandDialog) SetCancelFunc(handler func()) *CommandDialog {
	cmd.cancelHandler = handler
	cancelButton := cmd.form.GetButton(cmd.form.GetButtonCount() - 2)
	cancelButton.SetSelectedFunc(handler)
	return cmd
}

// SetRect set rects for this primitive.
func (cmd *CommandDialog) SetRect(x, y, width, height int) {

	ws := (width - cmd.width) / 2
	hs := ((height - cmd.height) / 2)
	dy := y + hs
	bWidth := cmd.width
	if cmd.width > width {
		ws = 0
		bWidth = width - 1
	}
	bHeight := cmd.height
	if cmd.height >= height {
		dy = y + 1
		bHeight = height - 1
	}
	cmd.Box.SetRect(x+ws, dy, bWidth, bHeight)
	x, y, width, height = cmd.Box.GetInnerRect()
	cmd.layout.SetRect(x, y, width, height)

}

// Draw draws this primitive onto the screen.
func (cmd *CommandDialog) Draw(screen tcell.Screen) {

	if !cmd.display {
		return
	}
	cmd.Box.DrawForSubclass(screen, cmd)
	cmd.layout.Draw(screen)
}
