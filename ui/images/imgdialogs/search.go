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
	searchFieldMaxSize   = 60
	searchButtonWidth    = 10
	searchInpuLabelWidth = 13

	// focus elements.
	sInputElement        = 1
	sSearchButtonElement = 2
	sSearchResultElement = 3
	sFormElement         = 4
)

const (
	searchResultIndexIndex = 0 + iota
	searchResultNameIndex
	searchResultDescIndex
	searchResultStarsIndex
	searchResultOfficialIndex
	searchResultAutomatedIndex
)

const (
	searchResultIndexColIndex = 0 + iota
	searchResultNameColIndex
	searchResultStarsColIndex
	searchResultOfficialColIndex
	searchResultAutomatedColIndex
	searchResultDescColIndex
)

// ImageSearchDialog represents image search dialogs.
type ImageSearchDialog struct {
	*tview.Box
	layout              *tview.Flex
	searchLayout        *tview.Flex
	input               *tview.InputField
	searchButton        *tview.Button
	searchResult        *tview.Table
	form                *tview.Form
	result              [][]string
	display             bool
	focusElement        int
	cancelHandler       func()
	searchSelectHandler func()
	pullSelectHandler   func()
}

// NewImageSearchDialog returns new image search dialog primitive.
func NewImageSearchDialog() *ImageSearchDialog {
	dialog := &ImageSearchDialog{
		Box:          tview.NewBox(),
		input:        tview.NewInputField(),
		searchButton: tview.NewButton("Search"),
		searchResult: tview.NewTable(),
		display:      false,
		focusElement: sInputElement,
	}

	bgColor := style.DialogBgColor
	fgColor := style.DialogFgColor
	inputFieldBgColor := style.InputFieldBgColor
	buttonBgColor := style.ButtonBgColor

	dialog.searchButton.SetBackgroundColor(buttonBgColor)
	dialog.searchButton.SetLabelColorActivated(buttonBgColor)

	dialog.input.SetLabel("search term: ")
	dialog.input.SetLabelColor(fgColor)
	dialog.input.SetFieldWidth(searchFieldMaxSize)
	dialog.input.SetBackgroundColor(bgColor)
	dialog.input.SetFieldBackgroundColor(inputFieldBgColor)

	dialog.searchLayout = tview.NewFlex().SetDirection(tview.FlexColumn)

	dialog.searchLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	dialog.searchLayout.AddItem(dialog.input, searchFieldMaxSize+searchInpuLabelWidth, 10, true) //nolint:gomnd
	dialog.searchLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	dialog.searchLayout.AddItem(dialog.searchButton, searchButtonWidth, 0, true)
	dialog.searchLayout.SetBackgroundColor(bgColor)

	dialog.searchResult.SetBackgroundColor(style.BgColor)
	dialog.searchResult.SetTitleColor(style.TableHeaderFgColor)
	dialog.searchResult.SetBorder(true)
	dialog.searchResult.SetBorderColor(style.DialogSubBoxBorderColor)

	searchResultLayout := tview.NewFlex().SetDirection(tview.FlexColumn)

	searchResultLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)
	searchResultLayout.AddItem(dialog.searchResult, 0, 1, true)
	searchResultLayout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, false)

	dialog.initTable()

	dialog.form = tview.NewForm().
		AddButton("Cancel", nil).
		AddButton("Pull", nil).
		SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(bgColor)
	dialog.form.SetButtonBackgroundColor(buttonBgColor)

	dialog.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	dialog.layout.SetBorder(true)
	dialog.layout.SetBorderColor(style.DialogBorderColor)
	dialog.layout.SetBackgroundColor(bgColor)
	dialog.layout.SetTitle("PODMAN IMAGE SEARCH/PULL")
	dialog.layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	dialog.layout.AddItem(dialog.searchLayout, 1, 0, true)
	dialog.layout.AddItem(utils.EmptyBoxSpace(bgColor), 1, 0, true)
	dialog.layout.AddItem(searchResultLayout, 0, 1, true)
	dialog.layout.AddItem(dialog.form, dialogs.DialogFormHeight, 0, true)

	return dialog
}

func (d *ImageSearchDialog) initTable() {
	bgColor := style.TableHeaderBgColor
	fgColor := style.TableHeaderFgColor

	d.searchResult.Clear()
	d.searchResult.SetCell(0, searchResultIndexColIndex,
		tview.NewTableCell(fmt.Sprintf("[%s::b]INDEX", style.GetColorHex(fgColor))).
			SetExpansion(1).
			SetBackgroundColor(bgColor).
			SetTextColor(fgColor).
			SetAlign(tview.AlignLeft).
			SetSelectable(false))

	d.searchResult.SetCell(0, searchResultNameColIndex,
		tview.NewTableCell(fmt.Sprintf("[%s::b]NAME", style.GetColorHex(fgColor))).
			SetExpansion(1).
			SetBackgroundColor(bgColor).
			SetTextColor(fgColor).
			SetAlign(tview.AlignLeft).
			SetSelectable(false))

	d.searchResult.SetCell(0, searchResultStarsColIndex,
		tview.NewTableCell(fmt.Sprintf("[%s::b]STARS", style.GetColorHex(fgColor))).
			SetExpansion(1).
			SetBackgroundColor(bgColor).
			SetTextColor(fgColor).
			SetAlign(tview.AlignCenter).
			SetSelectable(false))

	d.searchResult.SetCell(0, searchResultOfficialColIndex,
		tview.NewTableCell(fmt.Sprintf("[%s::b]OFFICIAL", style.GetColorHex(fgColor))).
			SetExpansion(1).
			SetBackgroundColor(bgColor).
			SetTextColor(fgColor).
			SetAlign(tview.AlignCenter).
			SetSelectable(false))
	d.searchResult.SetCell(0, searchResultAutomatedColIndex,
		tview.NewTableCell(fmt.Sprintf("[%s::b]AUTOMATED", style.GetColorHex(fgColor))).
			SetExpansion(1).
			SetBackgroundColor(bgColor).
			SetTextColor(fgColor).
			SetAlign(tview.AlignCenter).
			SetSelectable(false))

	d.searchResult.SetCell(0, searchResultDescColIndex,
		tview.NewTableCell(fmt.Sprintf("[%s::b]DESCRIPTION", style.GetColorHex(fgColor))).
			SetExpansion(1).
			SetBackgroundColor(bgColor).
			SetTextColor(fgColor).
			SetAlign(tview.AlignCenter).
			SetSelectable(false))

	d.searchResult.SetFixed(1, 1)
	d.searchResult.SetSelectable(true, false)
}

// Display displays this primitive.
func (d *ImageSearchDialog) Display() {
	d.display = true
}

// IsDisplay returns true if primitive is shown.
func (d *ImageSearchDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *ImageSearchDialog) Hide() {
	d.focusElement = sInputElement
	d.display = false

	d.input.SetText("")
	d.ClearResults()
}

// Focus is called when this primitive receives focus.
func (d *ImageSearchDialog) Focus(delegate func(p tview.Primitive)) { //nolint:cyclop
	switch d.focusElement {
	case sInputElement:
		d.input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = sSearchButtonElement
				d.Focus(delegate)

				return nil
			}

			if event.Key() == tcell.KeyDown {
				d.focusElement = sSearchResultElement
				d.Focus(delegate)

				return nil
			}

			if event.Key() == tcell.KeyEnter {
				d.searchSelectHandler()

				return nil
			}

			return event
		})

		delegate(d.input)

		return
	case sSearchButtonElement:
		d.searchButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = sSearchResultElement
				d.Focus(delegate)

				return nil
			}

			if event.Key() == tcell.KeyEnter {
				d.searchSelectHandler()

				return nil
			}

			return event
		})

		delegate(d.searchButton)

		return
	case sSearchResultElement:
		d.searchResult.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = sFormElement
				d.Focus(delegate)

				return nil
			}

			if event.Key() == tcell.KeyEnter {
				d.pullSelectHandler()

				return nil
			}

			return event
		})

		delegate(d.searchResult)

		return
	case sFormElement:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = sInputElement
				d.Focus(delegate)
				d.form.SetFocus(0)

				return nil
			}
			if event.Key() == tcell.KeyEnter {
				d.pullSelectHandler()

				return nil
			}

			return event
		})

		delegate(d.form)
	}
}

// HasFocus returns whether or not this primitive has focus.
func (d *ImageSearchDialog) HasFocus() bool {
	return d.form.HasFocus() || d.input.HasFocus() || d.searchResult.HasFocus() || d.searchButton.HasFocus()
}

// SetRect set rects for this primitive.
func (d *ImageSearchDialog) SetRect(x, y, width, height int) {
	paddingX := 1
	paddingY := 1
	dX := x + paddingX
	dY := y + paddingY
	dWidth := width - (2 * paddingX)   //nolint:gomnd
	dHeight := height - (2 * paddingY) //nolint:gomnd

	// set search input field size
	iwidth := dWidth - searchInpuLabelWidth - searchButtonWidth - 5 //nolint:gomnd
	if iwidth > searchFieldMaxSize {
		iwidth = searchFieldMaxSize
	}

	d.input.SetFieldWidth(iwidth)
	d.searchLayout.ResizeItem(d.input, iwidth+searchInpuLabelWidth, 0)

	// set table height size
	d.layout.ResizeItem(d.searchResult, dHeight-dialogs.DialogFormHeight-5, 0) //nolint:gomnd
	d.Box.SetRect(dX, dY, dWidth, dHeight)
}

// Draw draws this primitive onto the screen.
func (d *ImageSearchDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	bgColor := style.DialogBgColor
	d.Box.SetBackgroundColor(bgColor)
	d.Box.DrawForSubclass(screen, d)
	x, y, width, height := d.Box.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.SetBorder(true)
	d.layout.SetBackgroundColor(bgColor)

	d.layout.Draw(screen)
}

// InputHandler returns input handler function for this primitive.
func (d *ImageSearchDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("confirm dialog: event %v received", event)

		if event.Key() == tcell.KeyEsc {
			d.cancelHandler()

			return
		}

		if d.searchResult.HasFocus() {
			if searchResultHandler := d.searchResult.InputHandler(); searchResultHandler != nil {
				searchResultHandler(event, setFocus)

				return
			}
		}

		if d.input.HasFocus() {
			if inputFieldHandler := d.input.InputHandler(); inputFieldHandler != nil {
				inputFieldHandler(event, setFocus)

				return
			}
		}

		if d.searchButton.HasFocus() {
			if searchButtonHandler := d.searchButton.InputHandler(); searchButtonHandler != nil {
				searchButtonHandler(event, setFocus)

				return
			}
		}

		if d.form.HasFocus() {
			if formHandler := d.form.InputHandler(); formHandler != nil {
				formHandler(event, setFocus)

				return
			}
		}
	})
}

// SetCancelFunc sets form cancel button selected function.
func (d *ImageSearchDialog) SetCancelFunc(handler func()) *ImageSearchDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2) //nolint:gomnd
	cancelButton.SetSelectedFunc(handler)

	return d
}

// SetSearchFunc sets form cancel button selected function.
func (d *ImageSearchDialog) SetSearchFunc(handler func()) *ImageSearchDialog {
	d.searchSelectHandler = handler

	return d
}

// SetPullFunc sets form pull button selected function.
func (d *ImageSearchDialog) SetPullFunc(handler func()) *ImageSearchDialog {
	d.pullSelectHandler = handler

	return d
}

// GetSearchText returns search input field text.
func (d *ImageSearchDialog) GetSearchText() string {
	return d.input.GetText()
}

// GetSelectedItem returns selected image name from search result table.
func (d *ImageSearchDialog) GetSelectedItem() string {
	row, _ := d.searchResult.GetSelection()
	if row >= 0 {
		return d.result[row-1][1]
	}

	return ""
}

// ClearResults clear image search result table.
func (d *ImageSearchDialog) ClearResults() {
	d.UpdateResults([][]string{})
}

// UpdateResults updates result table.
func (d *ImageSearchDialog) UpdateResults(data [][]string) {
	d.result = data

	d.initTable()

	alignment := tview.AlignLeft
	rowIndex := 1
	expand := 1

	for i := 0; i < len(data); i++ {
		index := data[i][searchResultIndexIndex]
		name := data[i][searchResultNameIndex]
		desc := data[i][searchResultDescIndex]
		stars := data[i][searchResultStarsIndex]
		official := data[i][searchResultOfficialIndex]
		automated := data[i][searchResultAutomatedIndex]

		if official == "[OK]" {
			official = style.HeavyGreenCheckMark
		}

		if automated == "[OK]" {
			automated = style.HeavyGreenCheckMark
		}

		if strings.Index(name, index+"/") == 0 {
			name = strings.Replace(name, index+"/", "", 1)
		}

		// index column
		d.searchResult.SetCell(rowIndex, searchResultIndexColIndex,
			tview.NewTableCell(index).
				SetExpansion(expand).
				SetAlign(alignment))

		// name column
		d.searchResult.SetCell(rowIndex, searchResultNameColIndex,
			tview.NewTableCell(name).
				SetExpansion(expand).
				SetAlign(alignment))

		// stars column
		d.searchResult.SetCell(rowIndex, searchResultStarsColIndex,
			tview.NewTableCell(stars).
				SetExpansion(expand).
				SetAlign(tview.AlignCenter))

		// official column
		d.searchResult.SetCell(rowIndex, searchResultOfficialColIndex,
			tview.NewTableCell(official).
				SetExpansion(expand).
				SetAlign(tview.AlignCenter))

		// autoamted column
		d.searchResult.SetCell(rowIndex, searchResultAutomatedColIndex,
			tview.NewTableCell(automated).
				SetExpansion(expand).
				SetAlign(tview.AlignCenter))

		// description column
		d.searchResult.SetCell(rowIndex, searchResultDescColIndex,
			tview.NewTableCell(desc).
				SetExpansion(expand).
				SetAlign(alignment))

		rowIndex++
	}

	if len(data) > 0 {
		d.searchResult.Select(1, 1)
		d.searchResult.ScrollToBeginning()
	}
}
