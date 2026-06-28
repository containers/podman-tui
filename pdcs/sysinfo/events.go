package sysinfo

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"go.podman.io/podman/v6/pkg/bindings/system"
	"go.podman.io/podman/v6/pkg/domain/entities/types"
)

// Events returns libpod events.
func Events(eventChan chan types.Event, cancelChan chan bool) error {
	conn, err := registry.GetConnection()
	if err != nil {
		return err
	}

	return system.Events(conn, eventChan, cancelChan, new(system.EventsOptions).WithStream(true))
}
