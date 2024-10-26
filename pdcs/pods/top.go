package pods

import (
	"regexp"
	"strings"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/pdcs/utils"
	"github.com/containers/podman/v5/pkg/bindings/pods"
	"github.com/rs/zerolog/log"
)

// Top returns running processes on the pod.
func Top(id string) ([][]string, error) {
	log.Debug().Msgf("pdcs: podman pod top %s", id)

	report := [][]string{}

	conn, err := registry.GetConnection()
	if err != nil {
		return report, err
	}

	response, err := pods.Top(conn, id, new(pods.TopOptions))
	if err != nil {
		return report, err
	}

	for i := range response {
		if response[i] == "" {
			continue
		}

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

	if len(report) == 0 {
		return report, utils.ErrTopPodNotRunning
	}

	log.Debug().Msgf("pdcs: %v", report)

	return report, nil
}
