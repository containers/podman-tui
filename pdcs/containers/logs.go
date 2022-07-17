package containers

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/rs/zerolog/log"
)

// Logs returns container's log.
func Logs(id string) ([]string, error) {
	var logs []string

	logsBuffer := 20

	log.Debug().Msgf("pdcs: podman container logs %s", id)

	conn, err := registry.GetConnection()
	if err != nil {
		return logs, err
	}

	done := make(chan bool)
	logout := make(chan string, logsBuffer)
	logerr := make(chan string, logsBuffer)

	logAppender := func() {
		for {
			select {
			case msg := <-logout:
				logs = append(logs, msg)
			case msg := <-logerr:
				logs = append(logs, msg)
			case <-done:
				return
			}
		}
	}

	go logAppender()

	options := new(containers.LogOptions).WithFollow(false)

	err = containers.Logs(conn, id, options, logout, logerr)
	if err != nil {
		return logs, err
	}
	done <- true
	// logs = append(logs, <-logout)
	// logs = append(logs, <-logerr)

	return logs, nil
}
