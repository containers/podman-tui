package pods

import (
	"github.com/containers/podman/v3/pkg/bindings/pods"
	"github.com/containers/podman/v3/pkg/errorhandling"
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/rs/zerolog/log"
)

// Start starts a pod's containers
func Start(id string) error {
	log.Debug().Msgf("pdcs: podman pod start %s", id)
	conn, err := connection.GetConnection()
	if err != nil {
		return err
	}
	response, err := pods.Start(conn, id, new(pods.StartOptions))
	if err != nil {
		return err
	}
	var errs error
	if len(response.Errs) > 0 {
		errs = errorhandling.JoinErrors(response.Errs)
	}

	return errs
}
