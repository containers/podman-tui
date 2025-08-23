package networks

import (
	"github.com/containers/common/libnetwork/types"
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v5/pkg/bindings/network"
	"github.com/rs/zerolog/log"
)

// List returns list of podman networks.
func List() ([]types.Network, error) {
	log.Debug().Msg("pdcs: podman network ls")

	report := make([]types.Network, 0)

	conn, err := registry.GetConnection()
	if err != nil {
		return report, err
	}

	response, err := network.List(conn, new(network.ListOptions))
	if err != nil {
		return report, err
	}

	log.Debug().Msgf("pdcs: %v", response)

	return response, nil
}
