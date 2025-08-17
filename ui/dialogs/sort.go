package dialogs

import (
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	sortDialogMaxWidth  = 30
	sortDialogMaxHeight = 8
)

const (
	sortDialogOptionsFocus = 0 + iota
	sortDialogOrderFocus
	sortDialogFormFocus
)

type SortDialog struct {
	*tview.Box

	layout        *tview.Flex
	sortBy        *tview.DropDown
	sortOrder     *tview.DropDown
	form          *tview.Form
	display       bool
	selectHandler func()
	cancelHandler func()
	focusElement  int
}

func NewSortDialog(options []string, defaultOption int) *SortDialog {
	sd := SortDialog{
		Box:          tview.NewBox(),
		form:         tview.NewForm(),
		layout:       tview.NewFlex(),
		focusElement: sortDialogOptionsFocus,
	}

	sortOrder := "Sort order:"
	sd.sortOrder = tview.NewDropDown()
	sd.sortOrder.SetLabel(sortOrder)
	sd.sortOrder.SetLabelWidth(len(sortOrder) + 1)
	sd.sortOrder.SetTitleAlign(tview.AlignRight)
	sd.sortOrder.SetLabelColor(style.DialogFgColor)
	sd.sortOrder.SetBackgroundColor(style.DialogBgColor)
	sd.sortOrder.SetOptions([]string{"ascending", "descending"}, nil)
	sd.sortOrder.SetListStyles(style.DropDownUnselected, style.DropDownSelected)
	sd.sortOrder.SetFieldBackgroundColor(style.InputFieldBgColor)
	sd.sortOrder.SetFieldWidth(0)
	sd.sortOrder.SetCurrentOption(0)

	sd.sortBy = tview.NewDropDown()
	sd.sortBy.SetLabel("Sort by:")
	sd.sortBy.SetLabelWidth(len(sortOrder) + 1)
	sd.sortBy.SetTitleAlign(tview.AlignRight)
	sd.sortBy.SetLabelColor(style.DialogFgColor)
	sd.sortBy.SetBackgroundColor(style.DialogBgColor)
	sd.sortBy.SetOptions(options, nil)
	sd.sortBy.SetListStyles(style.DropDownUnselected, style.DropDownSelected)
	sd.sortBy.SetFieldBackgroundColor(style.InputFieldBgColor)
	sd.sortBy.SetFieldWidth(0)

	if len(options) > 0 {
		sd.sortBy.SetCurrentOption(defaultOption)
	}

	// form
	sd.form.AddButton("Cancel", nil)
	sd.form.AddButton(" Sort ", nil)
	sd.form.SetButtonsAlign(tview.AlignRight)
	sd.form.SetBackgroundColor(style.DialogBgColor)
	sd.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// main layout
	sd.layout.SetDirection(tview.FlexRow)
	sd.layout.SetBackgroundColor(style.DialogBgColor)
	sd.layout.SetBorder(true)
	sd.layout.SetBorderColor(style.DialogBorderColor)
	sd.layout.AddItem(sd.sortBy, 0, 1, true)
	sd.layout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 0, 1, false)
	sd.layout.AddItem(sd.sortOrder, 0, 1, true)
	sd.layout.AddItem(sd.form, DialogFormHeight, 0, true)

	return &sd
}

// Display displays this primitive.
func (d *SortDialog) Display() {
	d.display = true
}

// IsDisplay returns true if primitive is shown.
func (d *SortDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *SortDialog) Hide() {
	d.display = false
	d.focusElement = sortDialogOptionsFocus
	d.form.SetFocus(0)
}

// HasFocus returns whether or not this primitive has focus.
func (d *SortDialog) HasFocus() bool {
	if d.sortBy.HasFocus() || d.sortOrder.HasFocus() {
		return true
	}

	if d.form.HasFocus() || d.layout.HasFocus() {
		return true
	}

	return d.Box.HasFocus()
}

// Focus is called when this primitive receives focus.
func (d *SortDialog) Focus(delegate func(p tview.Primitive)) {
	switch d.focusElement {
	case sortDialogOptionsFocus:
		delegate(d.sortBy)
	case sortDialogOrderFocus:
		delegate(d.sortOrder)
	case sortDialogFormFocus:
		sortButton := d.form.GetButton(1)
		sortButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = sortDialogOptionsFocus

				d.Focus(delegate)
				d.form.SetFocus(0)

				return nil
			}

			if event.Key() == tcell.KeyEnter {
				d.selectHandler()
			}

			return event
		})

		delegate(d.form)
	}
}

// InputHandler returns input handler function for this primitive.
func (d *SortDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) { //nolint:cyclop
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("sort dialog dialog: event %v received", event)

		if event.Key() == utils.SwitchFocusKey.Key {
			switch d.focusElement {
			case sortDialogOptionsFocus:
				d.focusElement = sortDialogOrderFocus
			case sortDialogOrderFocus:
				d.focusElement = sortDialogFormFocus
			}
		}

		// dropdown widgets shall handle events before "Esc" key handler
		if d.sortBy.HasFocus() {
			event = utils.ParseKeyEventKey(event)
			if sortByHandler := d.sortBy.InputHandler(); sortByHandler != nil {
				sortByHandler(event, setFocus)

				return
			}
		}

		if d.sortOrder.HasFocus() {
			event = utils.ParseKeyEventKey(event)
			if sortOrderHandler := d.sortOrder.InputHandler(); sortOrderHandler != nil {
				sortOrderHandler(event, setFocus)

				return
			}
		}

		if d.form.HasFocus() {
			if event.Key() == tcell.KeyEsc {
				d.cancelHandler()

				return
			}

			if formHandler := d.form.InputHandler(); formHandler != nil {
				formHandler(event, setFocus)

				return
			}
		}
	})
}

// SetRect set rects for this primitive.
func (d *SortDialog) SetRect(x, y, width, height int) {
	if width > sortDialogMaxWidth {
		emptySpace := (width - sortDialogMaxWidth) / 2 //nolint:mnd
		x += emptySpace
		width = sortDialogMaxWidth
	}

	if height > sortDialogMaxHeight {
		emptySpace := (height - sortDialogMaxHeight) / 2 //nolint:mnd
		y += emptySpace
		height = sortDialogMaxHeight
	}

	d.Box.SetRect(x, y, width, height)
}

// Draw draws this primitive onto the screen.
func (d *SortDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.DrawForSubclass(screen, d)

	x, y, width, height := d.GetInnerRect()

	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetSelectFunc sets form sort button selected function.
func (d *SortDialog) SetSelectFunc(handler func(string, bool)) *SortDialog {
	selectHandler := func() {
		if d.sortBy.GetOptionCount() > 0 {
			_, sortOpt := d.sortBy.GetCurrentOption()
			_, order := d.sortOrder.GetCurrentOption()
			ascending := true

			if order == "descending" {
				ascending = false
			}

			handler(sortOpt, ascending)
		}

		d.Hide()
	}

	d.selectHandler = selectHandler

	return d
}

// SetCancelFunc sets form cancel button selected function.
func (d *SortDialog) SetCancelFunc(handler func()) *SortDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:mnd

	cancelButton.SetSelectedFunc(handler)

	return d
}
