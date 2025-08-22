package volumes

import (
	"fmt"
	"sort"
	"strings"

	"github.com/containers/podman-tui/pdcs/volumes"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// SortView sorts data view called from sort dialog.
func (vols *Volumes) SortView(option string, ascending bool) {
	log.Debug().Msgf("view: vols sort by %s", option)

	vols.volumeList.mu.Lock()
	defer vols.volumeList.mu.Unlock()

	vols.volumeList.sortBy = option
	vols.volumeList.ascending = ascending
	sort.Sort(volListSorted{vols.volumeList.report, option, ascending})
}

// UpdateData retrieves pods list data.
func (vols *Volumes) UpdateData() {
	volList, err := volumes.List()
	if err != nil {
		log.Error().Msgf("view: volumes update %v", err)
		vols.errorDialog.SetText(fmt.Sprintf("%v", err))
		vols.errorDialog.Display()

		return
	}

	vols.volumeList.mu.Lock()
	defer vols.volumeList.mu.Unlock()

	sort.Sort(volListSorted{volList, vols.volumeList.sortBy, vols.volumeList.ascending})

	vols.volumeList.report = volList
}

// ClearData clears table data.
func (vols *Volumes) ClearData() {
	vols.volumeList.mu.Lock()
	defer vols.volumeList.mu.Unlock()

	vols.volumeList.report = nil

	vols.table.Clear()

	expand := 1
	fgColor := style.PageHeaderFgColor
	bgColor := style.PageHeaderBgColor

	for i := range vols.headers {
		vols.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[::b]%s", strings.ToUpper(vols.headers[i]))). //nolint:perfsprint
													SetExpansion(expand).
													SetBackgroundColor(bgColor).
													SetTextColor(fgColor).
													SetAlign(tview.AlignLeft).
													SetSelectable(false))
	}

	vols.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(vols.title)))
}

func (vols *Volumes) getData() []*entities.VolumeListReport {
	vols.volumeList.mu.Lock()
	defer vols.volumeList.mu.Unlock()

	data := vols.volumeList.report

	return data
}

type lprSort []*entities.VolumeListReport

func (a lprSort) Len() int      { return len(a) }
func (a lprSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type volListSorted struct {
	lprSort

	option    string
	ascending bool
}

func (a volListSorted) Less(i, j int) bool {
	switch a.option {
	case "driver":
		if a.ascending {
			return a.lprSort[i].Driver < a.lprSort[j].Driver
		}

		return a.lprSort[i].Driver > a.lprSort[j].Driver
	case "mount point":
		if a.ascending {
			return a.lprSort[i].Mountpoint < a.lprSort[j].Mountpoint
		}

		return a.lprSort[i].Mountpoint > a.lprSort[j].Mountpoint
	case "created":
		if a.ascending {
			return a.lprSort[i].CreatedAt.After(a.lprSort[j].CreatedAt)
		}

		return a.lprSort[i].CreatedAt.Before(a.lprSort[j].CreatedAt)
	}

	if a.ascending {
		return a.lprSort[i].Name < a.lprSort[j].Name
	}

	return a.lprSort[i].Name > a.lprSort[j].Name
}
