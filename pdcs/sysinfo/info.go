package sysinfo

import (
	"encoding/json"

	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/containers/podman/v3/pkg/bindings/system"
	"github.com/rs/zerolog/log"
)

// Info returns podman system information
func Info() (string, error) {
	log.Debug().Msgf("pdcs: podman system info")
	var report string
	conn, err := connection.GetConnection()
	if err != nil {
		return report, err
	}

	sysInfo, err := system.Info(conn, new(system.InfoOptions))
	if err != nil {
		return report, err
	}

	b, err := json.MarshalIndent(sysInfo, "", "  ")
	if err != nil {
		return report, nil
	}
	report = string(b)

	return report, nil
}
