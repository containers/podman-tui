package containers

import (
	"regexp"
	"strings"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/rs/zerolog/log"
)

// Top returns running processes on the container.
func Top(id string) ([][]string, error) {
	log.Debug().Msgf("pdcs: podman container top %s", id)

	var report [][]string

	conn, err := registry.GetConnection()
	if err != nil {
		return report, err
	}

	response, err := containers.Top(conn, id, new(containers.TopOptions))
	if err != nil {
		return report, err
	}

	for i := 0; i < len(response); i++ {
		space := regexp.MustCompile(`\s+`)
		line := space.ReplaceAllString(response[i], " ")
		split := strings.Split(line, " ")
		user := split[0]
		pid := split[1]
		ppid := split[2]
		cpu := split[3]
		elapsed := split[4]
		tty := split[5]
		time := split[6]
		command := split[7:]
		cmd := strings.Join(command, " ")

		report = append(report, []string{
			user, pid, ppid, cpu, elapsed, tty, time, cmd,
		})
	}

	log.Debug().Msgf("pdcs: %v", report)

	return report, nil
}
