package dialogs

import (
	"fmt"

	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	cmdWidthOffset = 6
)

const (
	cmdTableFocus = 0 + iota
	cmdFormFocus
)

// CommandDialog is a commands list dialog.
type CommandDialog struct {
	*tview.Box

	layout        *tview.Flex
	table         *tview.Table
	form          *tview.Form
	display       bool
	options       [][]string
	width         int
	height        int
	focusElement  int
	selectedStyle tcell.Style
	cancelHandler func()
	selectHandler func()
}

// NewCommandDialog returns a command list primitive.
func NewCommandDialog(options [][]string) *CommandDialog {
	var cmdWidth int

	// command table items
	col1Width := 0
	col2Width := 0

	form := tview.NewForm().
		AddButton("Cancel", nil).
		SetButtonsAlign(tview.AlignRight)

	form.SetBackgroundColor(style.DialogBgColor)
	form.SetButtonBackgroundColor(style.ButtonBgColor)

	cmdsTable := tview.NewTable()
	cmdsTable.SetBackgroundColor(style.DialogBgColor)

	// command table header
	cmdsTable.SetCell(0, 0,
		tview.NewTableCell(fmt.Sprintf("[%s::b]COMMAND", style.GetColorHex(style.TableHeaderFgColor))).
			SetExpansion(1).
			SetBackgroundColor(style.TableHeaderBgColor).
			SetTextColor(style.TableHeaderFgColor).
			SetAlign(tview.AlignLeft).
			SetSelectable(false))

	cmdsTable.SetCell(0, 1,
		tview.NewTableCell(fmt.Sprintf("[%s::b]DESCRIPTION", style.GetColorHex(style.TableHeaderFgColor))).
			SetExpansion(1).
			SetBackgroundColor(style.TableHeaderBgColor).
			SetTextColor(style.TableHeaderFgColor).
			SetAlign(tview.AlignCenter).
			SetSelectable(false))

	for i := range options {
		cmdsTable.SetCell(i+1, 0,
			tview.NewTableCell(options[i][0]).
				SetAlign(tview.AlignLeft).
				SetSelectable(true).SetTextColor(style.DialogFgColor))
		cmdsTable.SetCell(i+1, 1,
			tview.NewTableCell(options[i][1]).
				SetAlign(tview.AlignLeft).
				SetSelectable(true).SetTextColor(style.DialogFgColor))

		if len(options[i][0]) > col1Width {
			col1Width = len(options[i][0])
		}

		if len(options[i][1]) > col2Width {
			col2Width = len(options[i][1])
		}
	}

	cmdWidth = col1Width + col2Width + 2 //nolint:mnd

	cmdsTable.SetFixed(1, 1)
	cmdsTable.SetSelectable(true, false)
	cmdsTable.SetBackgroundColor(style.DialogBgColor)

	cmdLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	cmdLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)
	cmdLayout.AddItem(cmdsTable, 0, 1, true)
	cmdLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)

	// layout
	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	layout.AddItem(cmdLayout, 0, 1, true)
	layout.AddItem(form, DialogFormHeight, 0, true)
	layout.SetBorder(true)
	layout.SetBorderColor(style.DialogBorderColor)
	layout.SetBackgroundColor(style.DialogBgColor)

	// returns the command primitive
	return &CommandDialog{
		Box:          tview.NewBox().SetBorder(false),
		layout:       layout,
		table:        cmdsTable,
		form:         form,
		display:      false,
		focusElement: cmdTableFocus,
		selectedStyle: tcell.StyleDefault.
			Background(style.DialogFgColor).
			Foreground(style.DialogBgColor),
		options: options,
		width:   cmdWidth + cmdWidthOffset,
		height:  len(options) + TableHeightOffset + DialogFormHeight,
	}
}

// GetSelectedItem returns selected row item.
func (cmd *CommandDialog) GetSelectedItem() string {
	row, _ := cmd.table.GetSelection()
	if row >= 0 {
		return cmd.options[row-1][0]
	}

	return ""
}

// GetCommandCount returns number of commands.
func (cmd *CommandDialog) GetCommandCount() int {
	return cmd.table.GetRowCount()
}

// Display displays this primitive.
func (cmd *CommandDialog) Display() {
	cmd.table.Select(1, 0)
	cmd.form.SetFocus(1)

	cmd.display = true
}

// IsDisplay returns true if primitive is shown.
func (cmd *CommandDialog) IsDisplay() bool {
	return cmd.display
}

// Hide stops displaying this primitive.
func (cmd *CommandDialog) Hide() {
	cmd.display = false
	cmd.focusElement = cmdTableFocus

	cmd.table.SetSelectedStyle(cmd.selectedStyle)
}

// HasFocus returns whether or not this primitive has focus.
func (cmd *CommandDialog) HasFocus() bool {
	if cmd.table.HasFocus() || cmd.form.HasFocus() {
		return true
	}

	return false
}

// Focus is called when this primitive receives focus.
func (cmd *CommandDialog) Focus(delegate func(p tview.Primitive)) {
	if cmd.focusElement == cmdTableFocus {
		delegate(cmd.table)

		return
	}

	button := cmd.form.GetButton(cmd.form.GetButtonCount() - 1)

	button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == utils.SwitchFocusKey.Key {
			cmd.focusElement = cmdTableFocus

			cmd.Focus(delegate)
			cmd.form.SetFocus(0)

			return nil
		}

		return event
	})

	delegate(cmd.form)
}

// InputHandler returns input handler function for this primitive.
func (cmd *CommandDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return cmd.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("commands dialog: event %v received", event)

		if event.Key() == utils.CloseDialogKey.Key {
			cmd.cancelHandler()

			return
		}

		if event.Key() == utils.SwitchFocusKey.Key {
			cmd.setFocusElement()
		}

		if cmd.form.HasFocus() {
			if formHandler := cmd.form.InputHandler(); formHandler != nil {
				formHandler(event, setFocus)

				return
			}
		}

		// command table handler
		if cmd.table.HasFocus() {
			if event.Key() == tcell.KeyEnter {
				cmd.selectHandler()

				return
			}

			if tableHandler := cmd.table.InputHandler(); tableHandler != nil {
				tableHandler(event, setFocus)

				return
			}
		}
	})
}

// SetSelectedFunc sets form enter button selected function.
func (cmd *CommandDialog) SetSelectedFunc(handler func()) *CommandDialog {
	cmd.selectHandler = handler

	return cmd
}

// SetCancelFunc sets form cancel button selected function.
func (cmd *CommandDialog) SetCancelFunc(handler func()) *CommandDialog {
	cmd.cancelHandler = handler
	cancelButton := cmd.form.GetButton(cmd.form.GetButtonCount() - 1)

	cancelButton.SetSelectedFunc(handler)

	return cmd
}

// SetRect set rects for this primitive.
func (cmd *CommandDialog) SetRect(x, y, width, height int) {
	ws := (width - cmd.width) / 2     //nolint:mnd
	hs := ((height - cmd.height) / 2) //nolint:mnd
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

	x, y, width, height = cmd.GetInnerRect()

	cmd.layout.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen.
func (cmd *CommandDialog) Draw(screen tcell.Screen) {
	if !cmd.display {
		return
	}

	cmd.DrawForSubclass(screen, cmd)
	cmd.layout.Draw(screen)
}

func (cmd *CommandDialog) setFocusElement() {
	if cmd.focusElement == cmdTableFocus {
		cmd.focusElement = cmdFormFocus
		cmd.table.SetSelectedStyle(tcell.StyleDefault.
			Background(style.DialogBgColor).
			Foreground(style.DialogFgColor))
	} else {
		cmd.focusElement = cmdTableFocus
		cmd.table.SetSelectedStyle(cmd.selectedStyle)
	}
}
