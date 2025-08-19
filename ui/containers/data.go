package containers

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/containers/podman-tui/pdcs/containers"
	putils "github.com/containers/podman-tui/pdcs/utils"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/docker/go-units"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// SortView sorts data view called from sort dialog.
func (cnt *Containers) SortView(option string, ascending bool) {
	log.Debug().Msgf("view: containers sort by %s", option)

	cnt.containersList.mu.Lock()
	defer cnt.containersList.mu.Unlock()

	cnt.containersList.sortBy = option
	cnt.containersList.ascending = ascending
	sort.Sort(containerListSorted{cnt.containersList.report, option, ascending})
}

// UpdateData retrieves containers list data.
func (cnt *Containers) UpdateData() {
	cntList, err := containers.List()
	if err != nil {
		log.Error().Msgf("view: containers update %v", err)
		cnt.errorDialog.SetText(fmt.Sprintf("%v", err))
		cnt.errorDialog.Display()

		return
	}

	cnt.containersList.mu.Lock()
	defer cnt.containersList.mu.Unlock()

	sort.Sort(containerListSorted{cntList, cnt.containersList.sortBy, cnt.containersList.ascending})

	cnt.containersList.report = cntList
}

func (cnt *Containers) getData() []entities.ListContainer {
	cnt.containersList.mu.Lock()
	defer cnt.containersList.mu.Unlock()

	data := cnt.containersList.report

	return data
}

// ClearData clears table data.
func (cnt *Containers) ClearData() {
	cnt.containersList.mu.Lock()
	defer cnt.containersList.mu.Unlock()

	cnt.containersList.report = nil
	cnt.table.Clear()

	expand := 1
	fgColor := style.PageHeaderFgColor
	bgColor := style.PageHeaderBgColor

	for i := range cnt.headers {
		cnt.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[::b]%s", strings.ToUpper(cnt.headers[i]))). //nolint:perfsprint
													SetExpansion(expand).
													SetBackgroundColor(bgColor).
													SetTextColor(fgColor).
													SetAlign(tview.AlignLeft).
													SetSelectable(false))
	}

	cnt.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(cnt.title)))
}

type conReporter struct {
	entities.ListContainer
}

func (con conReporter) names() string {
	return strings.Join(con.Names, ",")
}

func (con conReporter) state() string {
	var state string

	switch con.State {
	case "running":
		t := units.HumanDuration(time.Since(time.Unix(con.StartedAt, 0)))
		state = "Up " + t + " ago"
	case "configured":
		state = "Created"
	case "exited", "stopped":
		t := units.HumanDuration(time.Since(time.Unix(con.ExitedAt, 0)))
		state = fmt.Sprintf("Exited (%d) %s ago", con.ExitCode, t)
	default:
		state = con.State
	}

	return state
}

func (con conReporter) status() string {
	hc := con.Status
	if hc != "" {
		return con.state() + " (" + hc + ")"
	}

	return con.state()
}

func (con conReporter) ports() string {
	if len(con.Ports) < 1 {
		return ""
	}

	return putils.PortsToString(con.Ports)
}

type lprSort []entities.ListContainer

func (a lprSort) Len() int      { return len(a) }
func (a lprSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type containerListSorted struct {
	lprSort

	option    string
	ascending bool
}

func (a containerListSorted) Less(i, j int) bool {
	switch a.option {
	case "pod":
		if a.ascending {
			return a.lprSort[i].Pod < a.lprSort[j].Pod
		}

		return a.lprSort[i].Pod > a.lprSort[j].Pod
	case "image":
		if a.ascending {
			return a.lprSort[i].Image < a.lprSort[j].Image
		}

		return a.lprSort[i].Image > a.lprSort[j].Image
	case "created":
		if a.ascending {
			return a.lprSort[i].Created.After(a.lprSort[j].Created)
		}

		return a.lprSort[i].Created.Before(a.lprSort[j].Created)
	}

	if a.ascending {
		return a.lprSort[i].Names[0] < a.lprSort[j].Names[0]
	}

	return a.lprSort[i].Names[0] > a.lprSort[j].Names[0]
}
