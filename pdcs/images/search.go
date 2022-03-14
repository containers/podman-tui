package images

import (
	"fmt"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v4/pkg/bindings/images"
	"github.com/rs/zerolog/log"
)

// Search search repostiroy for images matche the search term
func Search(term string) ([][]string, error) {
	log.Debug().Msgf("pdcs: podman image search %s", term)
	var report [][]string
	conn, err := registry.GetConnection()
	if err != nil {
		return report, err
	}
	response, err := images.Search(conn, term, new(images.SearchOptions))
	if err != nil {
		return report, err
	}

	for _, sReport := range response {
		report = append(report, []string{
			sReport.Index,
			sReport.Name,
			sReport.Description,
			fmt.Sprintf("%d", sReport.Stars),
			sReport.Official,
			sReport.Automated,
		})

	}
	log.Debug().Msgf("pdcs: %s", report)
	return report, nil
}
