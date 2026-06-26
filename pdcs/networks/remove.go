package networks

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/network"
	"go.podman.io/podman/v6/pkg/errorhandling"
)

// Remove removes the network.
func Remove(name string) error {
	var errorReport []error

	log.Debug().Msgf("pdcs: podman network remove %s", name)

	conn, err := registry.GetConnection()
	if err != nil {
		return err
	}

	response, err := network.Remove(conn, name, new(network.RemoveOptions))
	if err != nil {
		return err
	}

	for _, r := range response {
		if r.Err != nil {
			errorReport = append(errorReport, r.Err)
		}
	}

	return errorhandling.JoinErrors(errorReport)
}
