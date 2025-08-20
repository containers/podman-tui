package pods

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	ppods "github.com/containers/podman-tui/pdcs/pods"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// SortView sorts data view called from sort dialog.
func (pods *Pods) SortView(option string, ascending bool) {
	log.Debug().Msgf("view: pods sort by %s", option)

	pods.podsList.mu.Lock()
	defer pods.podsList.mu.Unlock()

	pods.podsList.sortBy = option
	pods.podsList.ascending = ascending
	sort.Sort(containerListSorted{pods.podsList.report, option, ascending})
}

// UpdateData retrieves pods list data.
func (pods *Pods) UpdateData() {
	podList, err := ppods.List()
	if err != nil {
		log.Error().Msgf("view: pods update %v", err)
		pods.errorDialog.SetText(fmt.Sprintf("%v", err))
		pods.errorDialog.Display()

		return
	}

	pods.podsList.mu.Lock()
	defer pods.podsList.mu.Unlock()

	sort.Sort(containerListSorted{podList, pods.podsList.sortBy, pods.podsList.ascending})
	pods.podsList.report = podList
}

func (pods *Pods) getData() []*entities.ListPodsReport {
	pods.podsList.mu.Lock()
	defer pods.podsList.mu.Unlock()

	data := pods.podsList.report

	return data
}

// ClearData clears table data.
func (pods *Pods) ClearData() {
	pods.podsList.mu.Lock()
	defer pods.podsList.mu.Unlock()

	pods.podsList.report = nil

	pods.table.Clear()

	expand := 1
	fgColor := style.PageHeaderFgColor
	bgColor := style.PageHeaderBgColor

	for i := range pods.headers {
		pods.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[::b]%s", strings.ToUpper(pods.headers[i]))). //nolint:perfsprint
													SetExpansion(expand).
													SetBackgroundColor(bgColor).
													SetTextColor(fgColor).
													SetAlign(tview.AlignLeft).
													SetSelectable(false))
	}

	pods.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(pods.title)))
}

type lprSort []*entities.ListPodsReport

func (a lprSort) Len() int      { return len(a) }
func (a lprSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type containerListSorted struct {
	lprSort

	option    string
	ascending bool
}

func (a containerListSorted) Less(i, j int) bool {
	switch a.option {
	case "# of containers":
		iNumOfCnt := strconv.Itoa(len(a.lprSort[i].Containers))
		jNumOfCnt := strconv.Itoa(len(a.lprSort[j].Containers))

		if a.ascending {
			return iNumOfCnt < jNumOfCnt
		}

		return iNumOfCnt > jNumOfCnt
	case "status":
		if a.ascending {
			return a.lprSort[i].Status < a.lprSort[j].Status
		}

		return a.lprSort[i].Status > a.lprSort[j].Status
	case "created":
		if a.ascending {
			return a.lprSort[i].Created.After(a.lprSort[j].Created)
		}

		return a.lprSort[i].Created.Before(a.lprSort[j].Created)
	}

	if a.ascending {
		return a.lprSort[i].Name < a.lprSort[j].Name
	}

	return a.lprSort[i].Name > a.lprSort[j].Name
}
