package pods

import (
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/containers/podman/v4/pkg/bindings/pods"
	"github.com/containers/podman/v4/pkg/errorhandling"
	"github.com/rs/zerolog/log"
)

// Restart restarts a pod's containers
func Restart(id string) error {
	log.Debug().Msgf("pdcs: podman pod restart %s", id)
	conn, err := connection.GetConnection()
	if err != nil {
		return err
	}
	response, err := pods.Restart(conn, id, new(pods.RestartOptions))
	if err != nil {
		return err
	}
	var errs error
	if len(response.Errs) > 0 {
		errs = errorhandling.JoinErrors(response.Errs)
	}

	return errs
}
