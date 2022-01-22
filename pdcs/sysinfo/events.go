package sysinfo

import (
	"github.com/containers/podman/v3/pkg/bindings/system"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/containers/podman-tui/pdcs/connection"
)

// Events returns libpod events
func Events(eventChan chan entities.Event, cancelChan chan bool) error {
	conn, err := connection.GetConnection()
	if err != nil {
		return err
	}

	return system.Events(conn, eventChan, cancelChan, new(system.EventsOptions))
}
