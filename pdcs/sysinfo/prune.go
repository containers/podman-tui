package sysinfo

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/pdcs/utils"
	"github.com/containers/podman/v4/pkg/bindings/system"
	"github.com/rs/zerolog/log"
)

// Prune removes all unused pod, container, image and volume data.
func Prune() (string, error) {
	log.Debug().Msgf("pdcs: podman system prune")

	var report string

	conn, err := registry.GetConnection()
	if err != nil {
		return report, err
	}

	all := true
	allVolumes := true

	response, err := system.Prune(conn, &system.PruneOptions{
		All:     &all,
		Volumes: &allVolumes,
	})
	if err != nil {
		return report, err
	}

	report, err = utils.GetJSONOutput(response)
	if err != nil {
		return report, err
	}

	return report, nil
}
