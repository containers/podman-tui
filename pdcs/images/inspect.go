package images

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/pdcs/utils"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/rs/zerolog/log"
)

// Inspect inspects the specified image ID.
func Inspect(id string) (string, error) {
	log.Debug().Msgf("pdcs: podman image inspect %s", id)

	var report string

	conn, err := registry.GetConnection()
	if err != nil {
		return report, err
	}

	response, err := images.GetImage(conn, id, new(images.GetOptions))
	if err != nil {
		return report, err
	}

	report, err = utils.GetJSONOutput(response.ImageData)
	if err != nil {
		return report, err
	}

	log.Debug().Msgf("pdcs: %s", report)

	return report, nil
}
