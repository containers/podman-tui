package secrets

import (
	"fmt"
	"sort"
	"strings"

	"github.com/containers/podman-tui/pdcs/secrets"
	"github.com/containers/podman-tui/ui/style"
	"github.com/containers/podman/v5/pkg/domain/entities/types"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// SortView sorts data view called from sort dialog.
func (s *Secrets) SortView(option string, ascending bool) {
	log.Debug().Msgf("view: secrets sort by %s", option)

	s.secretList.mu.Lock()
	defer s.secretList.mu.Unlock()

	s.secretList.sortBy = option
	s.secretList.ascending = ascending

	sort.Sort(secretsListSorted{s.secretList.report, option, ascending})
}

// UpdateData retrieves secrets list data.
func (s *Secrets) UpdateData() {
	secResponse, err := secrets.List()
	if err != nil {
		log.Error().Msgf("view: secrets update %v", err)

		s.errorDialog.SetText(fmt.Sprintf("%v", err))
		s.errorDialog.Display()
	}

	s.secretList.mu.Lock()
	defer s.secretList.mu.Unlock()

	sort.Sort(secretsListSorted{secResponse, s.secretList.sortBy, s.secretList.ascending})

	s.secretList.report = secResponse
}

func (s *Secrets) getData() []*types.SecretInfoReport {
	s.secretList.mu.Lock()
	defer s.secretList.mu.Unlock()

	data := s.secretList.report

	return data
}

// ClearData clears table data.
func (s *Secrets) ClearData() {
	s.secretList.mu.Lock()
	defer s.secretList.mu.Unlock()

	s.secretList.report = nil

	s.table.Clear()

	expand := 1

	for i := range s.headers {
		s.table.SetCell(0, i,
			tview.NewTableCell(fmt.Sprintf("[::b]%s", strings.ToUpper(s.headers[i]))). //nolint:perfsprint
													SetExpansion(expand).
													SetBackgroundColor(style.PageHeaderBgColor).
													SetTextColor(style.PageHeaderFgColor).
													SetAlign(tview.AlignLeft).
													SetSelectable(false))
	}

	s.table.SetTitle(fmt.Sprintf("[::b]%s[0]", strings.ToUpper(s.title)))
}

type lprSort []*types.SecretInfoReport

func (a lprSort) Len() int      { return len(a) }
func (a lprSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type secretsListSorted struct {
	lprSort

	option    string
	ascending bool
}

func (a secretsListSorted) Less(i, j int) bool {
	switch a.option {
	case "driver":
		if a.ascending {
			return a.lprSort[i].Spec.Driver.Name < a.lprSort[j].Spec.Driver.Name
		}

		return a.lprSort[i].Spec.Driver.Name > a.lprSort[j].Spec.Driver.Name
	case "created":
		if a.ascending {
			return a.lprSort[i].CreatedAt.After(a.lprSort[j].CreatedAt)
		}

		return a.lprSort[i].CreatedAt.Before(a.lprSort[j].CreatedAt)
	case "updated":
		if a.ascending {
			return a.lprSort[i].UpdatedAt.After(a.lprSort[j].UpdatedAt)
		}

		return a.lprSort[i].UpdatedAt.Before(a.lprSort[j].UpdatedAt)
	}

	if a.ascending {
		return a.lprSort[i].Spec.Name < a.lprSort[j].Spec.Name
	}

	return a.lprSort[i].Spec.Name > a.lprSort[j].Spec.Name
}
