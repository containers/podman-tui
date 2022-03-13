package pods

import (
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/containers/podman/v4/pkg/bindings/pods"
	"github.com/containers/podman/v4/pkg/errorhandling"
	"github.com/rs/zerolog/log"
)

// Pause pauses a pod's containers
func Pause(id string) error {
	log.Debug().Msgf("pdcs: podman pod pause %s", id)
	conn, err := connection.GetConnection()
	if err != nil {
		return err
	}
	response, err := pods.Pause(conn, id, new(pods.PauseOptions))
	if err != nil {
		return err
	}
	var errs error
	if len(response.Errs) > 0 {
		errs = errorhandling.JoinErrors(response.Errs)
	}

	return errs
}

// Unpause unpauses a pod's containers
func Unpause(id string) error {
	log.Debug().Msgf("pdcs: podman pod unpause %s", id)
	conn, err := connection.GetConnection()
	if err != nil {
		return err
	}
	response, err := pods.Unpause(conn, id, new(pods.UnpauseOptions))
	if err != nil {
		return err
	}
	var errs error
	if len(response.Errs) > 0 {
		errs = errorhandling.JoinErrors(response.Errs)
	}

	return errs
}
