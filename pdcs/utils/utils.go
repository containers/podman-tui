package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/docker/go-units"
	"go.podman.io/common/libnetwork/types"
)

// SizeToStr converts size to human readable format.
func SizeToStr(size int64) string {
	return units.HumanSizeWithPrecision(float64(size), 3) //nolint:mnd
}

// CreatedToStr converts duration to human readable format.
func CreatedToStr(duration int64) string {
	created := time.Unix(duration, 0).UTC()

	return units.HumanDuration(time.Since(created)) + " ago"
}

// PrintJSON convert data interface to json string.
func PrintJSON(data []interface{}) (string, error) {
	buf, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

// Following code are from https://github.com/containers/podman/blob/main/cmd/podman/containers/ps.go

// PortsToString converts the ports used to a string of the from "port1, port2"
// and also groups a continuous list of ports into a readable format.
// The format is IP:HostPort(-Range)->ContainerPort(-Range)/Proto.
func PortsToString(ports []types.PortMapping) string {
	if len(ports) == 0 {
		return ""
	}

	sb := &strings.Builder{}

	for _, port := range ports {
		hostIP := port.HostIP
		if hostIP == "" {
			hostIP = "0.0.0.0"
		}

		protocols := strings.Split(port.Protocol, ",")

		for _, protocol := range protocols {
			if port.Range > 1 {
				fmt.Fprintf(sb, "%s:%d-%d->%d-%d/%s, ",
					hostIP, port.HostPort, port.HostPort+port.Range-1,
					port.ContainerPort, port.ContainerPort+port.Range-1, protocol)
			} else {
				fmt.Fprintf(sb, "%s:%d->%d/%s, ",
					hostIP, port.HostPort,
					port.ContainerPort, protocol)
			}
		}
	}

	display := sb.String()

	// make sure to trim the last ", " of the string
	return display[:len(display)-2]
}
