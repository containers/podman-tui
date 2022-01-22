package images

import (
	"github.com/containers/podman/v3/pkg/bindings/images"
	"github.com/containers/podman/v3/pkg/errorhandling"
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/rs/zerolog/log"
)

// Prune removes all un used specified images
func Prune() error {
	log.Debug().Msgf("pdcs: podman image prune")
	var errReport []error

	conn, err := connection.GetConnection()
	if err != nil {
		return err
	}
	response, err := images.Prune(conn, new(images.PruneOptions).WithAll(true))
	if err != nil {
		return nil
	}
	for _, r := range response {
		if r.Err != nil {
			errReport = append(errReport, r.Err)
		}
	}

	return errorhandling.JoinErrors(errReport)
}
