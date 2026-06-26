package networks

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/network"
	"go.podman.io/podman/v6/pkg/errorhandling"
)

// Prune removes all unused network.
func Prune() error {
	var errorReport []error

	log.Debug().Msgf("pdcs: podman network remove")

	conn, err := registry.GetConnection()
	if err != nil {
		return err
	}

	response, err := network.Prune(conn, new(network.PruneOptions))
	if err != nil {
		return err
	}

	for _, r := range response {
		if r.Error != nil {
			errorReport = append(errorReport, r.Error)
		}
	}

	return errorhandling.JoinErrors(errorReport)
}
