package pods

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/containers/podman/v3/pkg/bindings/pods"
	"github.com/rs/zerolog/log"
)

// Top returns running processes on the pod
func Top(id string) ([][]string, error) {
	log.Debug().Msgf("pdcs: podman pod top %s", id)
	var report [][]string
	conn, err := connection.GetConnection()
	if err != nil {
		return report, err
	}
	response, err := pods.Top(conn, id, new(pods.TopOptions))
	if err != nil {
		return report, err
	}
	for i := 0; i < len(response); i++ {
		if response[i] == "" {
			continue
		}
		space := regexp.MustCompile(`\s+`)
		line := space.ReplaceAllString(response[i], " ")
		splited := strings.Split(line, " ")
		user := splited[0]
		pid := splited[1]
		ppid := splited[2]
		cpu := splited[3]
		elapsed := splited[4]
		tty := splited[5]
		time := splited[6]
		command := splited[7:]
		cmd := strings.Join(command, " ")

		report = append(report, []string{
			user, pid, ppid, cpu, elapsed, tty, time, cmd,
		})
	}
	if len(report) == 0 {
		return report, fmt.Errorf("pods top can only be used on running pods")
	}
	log.Debug().Msgf("pdcs: %v", report)
	return report, nil
}
