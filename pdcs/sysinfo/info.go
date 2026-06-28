package sysinfo

import (
	"encoding/json"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/rs/zerolog/log"
	"go.podman.io/podman/v6/pkg/bindings/system"
)

// Info returns podman system information.
func Info() (string, error) {
	log.Debug().Msgf("pdcs: podman system info")

	var report string

	conn, err := registry.GetConnection()
	if err != nil {
		return report, err
	}

	sysInfo, err := system.Info(conn, new(system.InfoOptions))
	if err != nil {
		return report, err
	}

	b, err := json.MarshalIndent(sysInfo, "", "  ")
	if err != nil {
		return report, err
	}

	report = string(b)

	return report, nil
}
