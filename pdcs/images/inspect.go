package images

import (
	"strings"

	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/containers/podman-tui/pdcs/utils"
	"github.com/containers/podman/v4/pkg/bindings/images"
	"github.com/rs/zerolog/log"
)

// Inspect inspects the specified image ID
func Inspect(id string) (string, error) {
	log.Debug().Msgf("pdcs: podman image inspect %s", id)
	var report string
	conn, err := connection.GetConnection()
	if err != nil {
		return report, err
	}
	response, err := images.List(conn, new(images.ListOptions).WithAll(true))
	if err != nil {
		return report, err
	}

	for index, imgSumm := range response {
		if strings.Index(imgSumm.ID, id) == 0 {
			report, err = utils.GetJSONOutput(response[index])
			if err != nil {
				return report, err
			}
		}

	}
	log.Debug().Msgf("pdcs: %s", report)
	return report, nil
}
