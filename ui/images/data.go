package images

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/containers/podman-tui/pdcs/images"
	"github.com/containers/podman-tui/ui/style"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// SortView sorts data view called from sort dialog.
func (img *Images) SortView(option string, ascending bool) {
	log.Debug().Msgf("view: images sort by %s", option)

	img.imagesList.mu.Lock()
	defer img.imagesList.mu.Unlock()

	img.imagesList.sortBy = option
	img.imagesList.ascending = ascending

	sort.Sort(imgListSorted{img.imagesList.report, option, ascending})
}

// UpdateData retrieves images list data.
func (img *Images) UpdateData() {
	images, err := images.List()
	if err != nil {
		log.Error().Msgf("view: images update %v", err)
		img.errorDialog.SetText(fmt.Sprintf("%v", err))
		img.errorDialog.Display()

		return
	}

	img.imagesList.mu.Lock()
	defer img.imagesList.mu.Unlock()

	sort.Sort(imgListSorted{images, img.imagesList.sortBy, img.imagesList.ascending})

	img.imagesList.report = images
}

func (img *Images) getData() []images.ImageListReporter {
	img.imagesList.mu.Lock()
	defer img.imagesList.mu.Unlock()

	data := img.imagesList.report

	return data
}

// ClearData clears table data.
func (img *Images) ClearData() {
	img.imagesList.mu.Lock()
	defer img.imagesList.mu.Unlock()

	img.imagesList.report = nil

	img.table.Clear()

	expand := 1
	fgColor := style.PageHeaderFgColor
	bgColor := style.PageHeaderBgColor

	for i := range img.headers {
		img.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[::b]%s", strings.ToUpper(img.headers[i]))). //nolint:perfsprint
													SetExpansion(expand).
													SetBackgroundColor(bgColor).
													SetTextColor(fgColor).
													SetAlign(tview.AlignLeft).
													SetSelectable(false))
	}

	img.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(img.title)))
}

type lprSort []images.ImageListReporter

func (a lprSort) Len() int      { return len(a) }
func (a lprSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type imgListSorted struct {
	lprSort

	option    string
	ascending bool
}

func (a imgListSorted) Less(i, j int) bool {
	switch a.option {
	case "size":
		if a.ascending {
			return a.lprSort[i].Size < a.lprSort[j].Size
		}

		return a.lprSort[i].Size > a.lprSort[j].Size
	case "created":
		icreated := time.Unix(a.lprSort[i].Created, 0).UTC()
		jcreated := time.Unix(a.lprSort[j].Created, 0).UTC()

		if a.ascending {
			return icreated.After(jcreated)
		}

		return icreated.Before(jcreated)
	}

	if a.ascending {
		return a.lprSort[i].Repository < a.lprSort[j].Repository
	}

	return a.lprSort[i].Repository > a.lprSort[j].Repository
}
