package pods

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/pods"
	"go.podman.io/podman/v6/pkg/errorhandling"
)

// Stop stops a pod's containers.
func Stop(id string) error {
	log.Debug().Msgf("pdcs: podman pod stop %s", id)

	conn, err := registry.GetConnection()
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
