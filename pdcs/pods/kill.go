package pods

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/pods"
	"go.podman.io/podman/v6/pkg/errorhandling"
)

// Kill sends SIGKILL signal to a pod's containers processeses.
func Kill(id string) error {
	log.Debug().Msgf("pdcs: podman pod kill %s", id)

	conn, err := registry.GetConnection()
	if err != nil {
		return err
	}

	response, err := pods.Kill(conn, id, new(pods.KillOptions))
	if err != nil {
		return err
	}

	var errs error
	if len(response.Errs) > 0 {
		errs = errorhandling.JoinErrors(response.Errs)
	}

	return errs
}
