package pods

import (
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/containers/podman/v4/pkg/bindings/pods"
	"github.com/containers/podman/v4/pkg/errorhandling"
	"github.com/rs/zerolog/log"
)

// Stop stops a pod's containers
func Stop(id string) error {
	log.Debug().Msgf("pdcs: podman pod stop %s", id)
	conn, err := connection.GetConnection()
	if err != nil {
		return err
	}
	response, err := pods.Stop(conn, id, new(pods.StopOptions))
	if err != nil {
		return err
	}
	var errs error
	if len(response.Errs) > 0 {
		errs = errorhandling.JoinErrors(response.Errs)
	}

	return errs
}
