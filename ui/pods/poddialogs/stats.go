package poddialogs

import (
	"fmt"
	"sync"
	"time"

	ppods "github.com/containers/podman-tui/pdcs/pods"
	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	podStatDialogFormFocus = 0 + iota
	podStatDialogPodDropDownFocus
	podStatDialogPodDropDownSoryByFocus
	podStatDialogResultTableFocus
)

const (
	podStatTablePodIDIndex = 0 + iota
	podStatTableCIDIndex
	podStatTableNameIndex
	podStatTableCPUPercIndex
	podStatTableMemUsageIndex
	podStatTableMemPercIndex
	podStatTableNetIOIndex
	podStatTableBlockIOIndex
	podStatTablePidsIndex
)

// PodStatsDialog implements the pods stats dialog primitive.
type PodStatsDialog struct {
	*tview.Box

	layout               *tview.Flex
	controlLayout        *tview.Flex
	form                 *tview.Form
	table                *tview.Table
	podDropDown          *tview.DropDown
	podSortByDropDown    *tview.DropDown
	mu                   sync.Mutex
	doneHandler          func()
	podDropDownOptions   []PodStatsDropDownOptions
	focusElement         int
	statQueryOpts        *ppods.StatsOptions
	statsResult          []ppods.StatReporter
	queryRefreshInterval time.Duration
	doneChan             chan bool
	display              bool
}

// PodStatsDropDownOptions implements pods dropdown options.
type PodStatsDropDownOptions struct {
	ID   string
	Name string
}

// NewPodStatsDialog returns new pod stats dialog.
func NewPodStatsDialog() *PodStatsDialog {
	statsDialog := PodStatsDialog{
		Box:                  tview.NewBox(),
		podDropDown:          tview.NewDropDown(),
		podSortByDropDown:    tview.NewDropDown(),
		statQueryOpts:        &ppods.StatsOptions{},
		queryRefreshInterval: 3000 * time.Millisecond, //nolint:mnd
	}

	ddUnselectedStyle := style.DropDownUnselected
	ddselectedStyle := style.DropDownSelected

	// pod dropdown
	pddLabel := "POD ID:"
	labelBgColor := fmt.Sprintf("#%x", style.DialogBorderColor.Hex())

	statsDialog.podDropDown.SetLabel(fmt.Sprintf("[:%s:b]%s[::-]", labelBgColor, pddLabel))
	statsDialog.podDropDown.SetLabelWidth(len(pddLabel) + 1)
	statsDialog.podDropDown.SetBackgroundColor(style.DialogBgColor)
	statsDialog.podDropDown.SetLabelColor(style.DialogFgColor)
	statsDialog.podDropDown.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	statsDialog.podDropDown.SetFocusedStyle(style.DropDownFocused)
	statsDialog.podDropDown.SetFieldBackgroundColor(style.InputFieldBgColor)

	// pod sortby dropdown
	pddSortByLabel := "SORT BY:"

	statsDialog.podSortByDropDown.SetLabel(fmt.Sprintf("[:%s:b]%s[::-]", labelBgColor, pddSortByLabel))
	statsDialog.podSortByDropDown.SetLabelWidth(len(pddSortByLabel) + 1)
	statsDialog.podSortByDropDown.SetBackgroundColor(style.DialogBgColor)
	statsDialog.podSortByDropDown.SetLabelColor(style.DialogFgColor)
	statsDialog.podSortByDropDown.SetListStyles(ddUnselectedStyle, ddselectedStyle)
	statsDialog.podSortByDropDown.SetFocusedStyle(style.DropDownFocused)
	statsDialog.podSortByDropDown.SetOptions([]string{
		"pod ID",
		"container name",
		"cpu %",
		"mem %",
	}, statsDialog.setStatsQuerySortBy)
	statsDialog.podSortByDropDown.SetFieldBackgroundColor(style.InputFieldBgColor)

	// table
	statsDialog.table = tview.NewTable()
	statsDialog.table.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	statsDialog.table.SetBorder(true)
	statsDialog.table.SetBorderColor(style.DialogSubBoxBorderColor)
	statsDialog.initTableUI()

	// form
	statsDialog.form = tview.NewForm().
		AddButton("Cancel", nil).
		SetButtonsAlign(tview.AlignRight)
	statsDialog.form.SetBackgroundColor(style.DialogBgColor)
	statsDialog.form.SetButtonBackgroundColor(style.ButtonBgColor)

	// pod dropdown and sort by dropdown
	statsDialog.controlLayout = tview.NewFlex().SetDirection(tview.FlexColumn)
	statsDialog.controlLayout.SetBackgroundColor(style.DialogBgColor)
	statsDialog.controlLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)
	statsDialog.controlLayout.AddItem(statsDialog.podDropDown, 0, 1, false)
	statsDialog.controlLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)
	statsDialog.controlLayout.AddItem(statsDialog.podSortByDropDown, 0, 1, false)
	statsDialog.controlLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)

	// table layout
	statLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	statLayout.SetBackgroundColor(style.DialogBgColor)
	statLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)
	statLayout.AddItem(statsDialog.table, 0, 1, false)
	statLayout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, false)

	// main dialog layout
	statsDialog.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	statsDialog.layout.SetBorder(true)
	statsDialog.layout.SetBorderColor(style.BorderColor)
	statsDialog.layout.SetBackgroundColor(style.DialogBgColor)
	statsDialog.layout.SetTitle("PODMAN POD STATS")

	statsDialog.layout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	statsDialog.layout.AddItem(statsDialog.controlLayout, 1, 0, true)
	statsDialog.layout.AddItem(utils.EmptyBoxSpace(style.DialogBgColor), 1, 0, true)
	statsDialog.layout.AddItem(statLayout, 0, 1, true)
	statsDialog.layout.AddItem(statsDialog.form, dialogs.DialogFormHeight, 0, true)

	return &statsDialog
}

// Display displays this primitive.
func (d *PodStatsDialog) Display() {
	d.display = true

	d.podSortByDropDown.SetCurrentOption(0)

	d.focusElement = podStatDialogResultTableFocus
	d.doneChan = make(chan bool)

	d.startStatsQueryLoop()
}

// IsDisplay returns true if primitive is shown.
func (d *PodStatsDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive.
func (d *PodStatsDialog) Hide() {
	d.display = false
	d.doneChan <- true

	d.SetPodsOptions([]PodStatsDropDownOptions{})

	d.mu.Lock()

	defer d.mu.Unlock()

	close(d.doneChan)
}

// HasFocus returns whether or not this primitive has focus.
func (d *PodStatsDialog) HasFocus() bool {
	if d.podDropDown.HasFocus() || d.podSortByDropDown.HasFocus() {
		return true
	}

	if d.table.HasFocus() || d.form.HasFocus() {
		return true
	}

	return d.Box.HasFocus()
}

// Focus is called when this primitive receives focus.
func (d *PodStatsDialog) Focus(delegate func(p tview.Primitive)) {
	switch d.focusElement {
	case podStatDialogFormFocus:
		delegate(d.form)
	case podStatDialogPodDropDownFocus:
		delegate(d.podDropDown)
	case podStatDialogPodDropDownSoryByFocus:
		delegate(d.podSortByDropDown)
	case podStatDialogResultTableFocus:
		delegate(d.table)
	}
}

// InputHandler  returns input handler function for this primitive.
func (d *PodStatsDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) { //nolint:gocognit,lll,cyclop
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("pod stats dialog: event %v received", event)
		// pod ID dropdown
		if d.podDropDown.HasFocus() {
			if event.Key() == tcell.KeyTab {
				d.focusElement = podStatDialogPodDropDownSoryByFocus
				setFocus(d)

				return
			}

			if podDropDownHandler := d.podDropDown.InputHandler(); podDropDownHandler != nil {
				event = utils.ParseKeyEventKey(event)
				podDropDownHandler(event, setFocus)

				return
			}
		}

		// sortby dropdown
		if d.podSortByDropDown.HasFocus() {
			if event.Key() == tcell.KeyTab {
				d.focusElement = podStatDialogResultTableFocus
				setFocus(d)

				return
			}

			if podSortByDropDownHandler := d.podSortByDropDown.InputHandler(); podSortByDropDownHandler != nil {
				event = utils.ParseKeyEventKey(event)
				podSortByDropDownHandler(event, setFocus)

				return
			}
		}

		// Esc key shall be after drop down so it won't overwrite default
		// dropdown handler
		if event.Key() == tcell.KeyEsc {
			d.doneHandler()

			return
		}

		// form
		if d.form.HasFocus() {
			if event.Key() == tcell.KeyTab {
				d.focusElement = podStatDialogPodDropDownFocus
				setFocus(d)

				return
			}

			if formHandler := d.form.InputHandler(); formHandler != nil {
				formHandler(event, setFocus)

				return
			}
		}

		// stats table
		if d.table.HasFocus() {
			if event.Key() == tcell.KeyTab {
				d.focusElement = podStatDialogFormFocus
				setFocus(d)

				return
			}

			if tableHanlder := d.table.InputHandler(); tableHanlder != nil {
				tableHanlder(event, setFocus)

				return
			}
		}
	})
}

// Draw draws this primitive onto the screen.
func (d *PodStatsDialog) Draw(screen tcell.Screen) {
	if !d.display {
		return
	}

	d.DrawForSubclass(screen, d)

	x, y, width, height := d.GetInnerRect()

	d.layout.SetRect(x, y, width, height)
	d.layout.Draw(screen)
}

// SetRect set rects for this primitive.
func (d *PodStatsDialog) SetRect(x, y, width, height int) {
	dX := x + 1
	dY := y + 1
	dWidth := width - 2   //nolint:mnd
	dHeight := height - 2 //nolint:mnd

	d.Box.SetRect(dX, dY, dWidth, dHeight)
}

// SetDoneFunc sets form cancel button selected function.
func (d *PodStatsDialog) SetDoneFunc(handler func()) *PodStatsDialog {
	d.doneHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 1)

	cancelButton.SetSelectedFunc(handler)

	return d
}

// SetPodsOptions sets pod drop down options.
func (d *PodStatsDialog) SetPodsOptions(options []PodStatsDropDownOptions) {
	maxWidth := 0
	d.podDropDownOptions = options

	if len(options) == 0 {
		return
	}

	podOptions := []string{"all"}

	for i := range options {
		item := options[i].ID
		if options[i].Name != "" {
			item = fmt.Sprintf("%s (%s)", item, options[i].Name)
		}

		if len(item) > maxWidth {
			maxWidth = len(item)
		}

		podOptions = append(podOptions, item)
	}

	maxWidth += 11
	d.controlLayout.ResizeItem(d.podDropDown, maxWidth, 0)
	d.podDropDown.SetOptions(podOptions, d.setStatsQueryPodIDs)
	d.podDropDown.SetCurrentOption(0)
}

func (d *PodStatsDialog) query() {
	opts := d.getStatsQueryOptions()

	podStats, err := ppods.Stats(opts)
	if err != nil {
		log.Error().Msgf("pod stats dialog: query error: %v", err)

		return
	}

	d.updateData(podStats)
}

func (d *PodStatsDialog) startStatsQueryLoop() {
	log.Debug().Msgf("pod stats dialog: starting pod stats query loop")

	go func() {
		tick := time.NewTicker(d.queryRefreshInterval)
		// initial query
		d.query()

		for {
			select {
			case <-tick.C:
				d.query()
			case <-d.doneChan:
				log.Debug().Msgf("pod stats dialog: stats query loop stopped")

				return
			}
		}
	}()
}

func (d *PodStatsDialog) getStatsQueryOptions() *ppods.StatsOptions {
	d.mu.Lock()
	opts := d.statQueryOpts
	d.mu.Unlock()

	return opts
}

func (d *PodStatsDialog) setStatsQueryPodIDs(name string, index int) {
	if index == -1 {
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if index == 0 {
		d.statQueryOpts.IDs = d.getAllPodIDs()
	} else {
		d.statQueryOpts.IDs = []string{
			d.podDropDownOptions[index-1].ID,
		}
	}

	go d.query()
}

func (d *PodStatsDialog) setStatsQuerySortBy(name string, index int) {
	if index == -1 {
		return
	}

	d.mu.Lock()

	defer d.mu.Unlock()

	d.statQueryOpts.SortBy = index

	go d.query()
}

func (d *PodStatsDialog) getAllPodIDs() []string {
	ids := make([]string, 0)

	for _, item := range d.podDropDownOptions {
		ids = append(ids, item.ID)
	}

	return ids
}

func (d *PodStatsDialog) initTableUI() {
	tableHeaders := []string{"POD ID", "CID", "NAME", "CPU %", "MEM USAGE / LIMIT", "MEM %", "NET IO", "BLOCK IO", "PIDS"}

	headerBgColor := style.TableHeaderBgColor
	headerFgColor := style.TableHeaderFgColor

	d.table.Clear()

	for index, header := range tableHeaders {
		headerItem := fmt.Sprintf("[::b]%s[::-]", header)
		d.table.SetCell(0, index,
			tview.NewTableCell(headerItem).
				SetExpansion(1).
				SetAlign(tview.AlignLeft).
				SetBackgroundColor(headerBgColor).
				SetTextColor(headerFgColor).
				SetSelectable(false))
	}

	d.table.SetFixed(1, 1)
	d.table.SetSelectable(true, false)
}

func (d *PodStatsDialog) updateData(statReport []ppods.StatReporter) {
	d.mu.Lock()
	d.statsResult = statReport
	d.mu.Unlock()

	fgColor := style.DialogFgColor
	row := 1

	d.initTableUI()

	for i := range d.statsResult {
		podID := d.statsResult[i].Pod
		cntID := d.statsResult[i].CID
		cntName := d.statsResult[i].Name
		cpuPerc := d.statsResult[i].CPU
		memUsage := d.statsResult[i].MemUsage
		memPerc := d.statsResult[i].Mem
		netIO := d.statsResult[i].NetIO
		blockIO := d.statsResult[i].BlockIO
		pids := d.statsResult[i].PIDS

		// POD ID
		d.table.SetCell(row, podStatTablePodIDIndex,
			tview.NewTableCell(podID).
				SetExpansion(1).
				SetAlign(tview.AlignLeft).
				SetTextColor(fgColor))

		// CID
		d.table.SetCell(row, podStatTableCIDIndex,
			tview.NewTableCell(cntID).
				SetExpansion(1).
				SetAlign(tview.AlignLeft).
				SetTextColor(fgColor))

		// Name
		d.table.SetCell(row, podStatTableNameIndex,
			tview.NewTableCell(cntName).
				SetExpansion(1).
				SetAlign(tview.AlignLeft).
				SetTextColor(fgColor))

		// CPU %
		d.table.SetCell(row, podStatTableCPUPercIndex,
			tview.NewTableCell(cpuPerc).
				SetExpansion(1).
				SetAlign(tview.AlignLeft).
				SetTextColor(fgColor))

		// MEM USAGE
		d.table.SetCell(row, podStatTableMemUsageIndex,
			tview.NewTableCell(memUsage).
				SetExpansion(1).
				SetAlign(tview.AlignLeft).
				SetTextColor(fgColor))

		// MEM %
		d.table.SetCell(row, podStatTableMemPercIndex,
			tview.NewTableCell(memPerc).
				SetExpansion(1).
				SetAlign(tview.AlignLeft).
				SetTextColor(fgColor))

		// NET IO
		d.table.SetCell(row, podStatTableNetIOIndex,
			tview.NewTableCell(netIO).
				SetExpansion(1).
				SetAlign(tview.AlignLeft).
				SetTextColor(fgColor))

		// BLOCK IO
		d.table.SetCell(row, podStatTableBlockIOIndex,
			tview.NewTableCell(blockIO).
				SetExpansion(1).
				SetAlign(tview.AlignLeft).
				SetTextColor(fgColor))

		// PIDs IO
		d.table.SetCell(row, podStatTablePidsIndex,
			tview.NewTableCell(pids).
				SetExpansion(1).
				SetAlign(tview.AlignLeft).
				SetTextColor(fgColor))

		row++
	}
}
