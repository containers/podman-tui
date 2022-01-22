package networks

import (
	"github.com/containers/podman/v3/pkg/bindings/network"
	"github.com/containers/podman/v3/pkg/errorhandling"
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/rs/zerolog/log"
)

// Remove removes the network.
func Remove(name string) error {
	var errorReport []error
	log.Debug().Msgf("pdcs: podman network remove %s", name)
	conn, err := connection.GetConnection()
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
