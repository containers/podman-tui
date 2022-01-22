package images

import (
	"time"

	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/containers/podman-tui/pdcs/utils"
	"github.com/containers/podman/v3/pkg/bindings/images"
	"github.com/docker/go-units"
	"github.com/rs/zerolog/log"
)

// History returns history of an image
func History(id string) ([][]string, error) {
	log.Debug().Msgf("pdcs: podman image history %s", id)
	var report [][]string
	conn, err := connection.GetConnection()
	if err != nil {
		return report, err
	}
	response, err := images.History(conn, id, new(images.HistoryOptions))
	if err != nil {
		return report, err
	}
	for i := 0; i < len(response); i++ {
		report = append(report, []string{
			response[i].ID,
			units.HumanDuration(time.Since(time.Unix(response[i].Created, 0))) + " ago",
			response[i].CreatedBy,
			utils.SizeToStr(response[i].Size),
			response[i].Comment,
		})
	}
	log.Debug().Msgf("pdcs: %v", report)
	return report, nil
}
