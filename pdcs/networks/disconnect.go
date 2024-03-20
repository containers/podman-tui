package networks

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v5/pkg/bindings/network"
	"github.com/rs/zerolog/log"
)

// Disconnect disconnects a container from a network.
func Disconnect(networkName string, containerID string) error {
	log.Debug().Msgf("pdcs: podman network (%s) disconnect (%s)", networkName, containerID)

	conn, err := registry.GetConnection()
	if err != nil {
		return err
	}

	return network.Disconnect(conn, networkName, containerID, new(network.DisconnectOptions).WithForce(true))
}
