package sysinfo

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v5/pkg/bindings/system"
	"github.com/containers/podman/v5/pkg/domain/entities"
)

// Events returns libpod events.
func Events(eventChan chan entities.Event, cancelChan chan bool) error {
	conn, err := registry.GetConnection()
	if err != nil {
		return err
	}

	return system.Events(conn, eventChan, cancelChan, new(system.EventsOptions))
}
