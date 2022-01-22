package utils

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cri-o/ocicni/pkg/ocicni"
	"github.com/docker/go-units"
)

// SizeToStr converts size to human readable format
func SizeToStr(size int64) string {
	return units.HumanSizeWithPrecision(float64(size), 3)
}

// CreatedToStr converts duration to human readable format
func CreatedToStr(duration int64) string {
	created := time.Unix(duration, 0).UTC()
	return units.HumanDuration(time.Since(created)) + " ago"
}

// PrintJSON convert data interface to json string
func PrintJSON(data []interface{}) (string, error) {
	buf, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// Following code are from https://github.com/containers/podman/blob/v3.4.2/cmd/podman/containers/ps.go

// PortsToString converts the ports used to a string of the from "port1, port2"
// and also groups a continuous list of ports into a readable format.
func PortsToString(ports []ocicni.PortMapping) string {
	if len(ports) == 0 {
		return ""
	}
	// Sort the ports, so grouping continuous ports become easy.
	sort.Slice(ports, func(i, j int) bool {
		return comparePorts(ports[i], ports[j])
	})

	portGroups := [][]ocicni.PortMapping{}
	currentGroup := []ocicni.PortMapping{}
	for i, v := range ports {
		var prevPort, nextPort *int32
		if i > 0 {
			prevPort = &ports[i-1].ContainerPort
		}
		if i+1 < len(ports) {
			nextPort = &ports[i+1].ContainerPort
		}

		port := v.ContainerPort

		// Helper functions
		addToCurrentGroup := func(x ocicni.PortMapping) {
			currentGroup = append(currentGroup, x)
		}

		addToPortGroup := func(x ocicni.PortMapping) {
			portGroups = append(portGroups, []ocicni.PortMapping{x})
		}

		finishCurrentGroup := func() {
			portGroups = append(portGroups, currentGroup)
			currentGroup = []ocicni.PortMapping{}
		}

		// Single entry slice
		if prevPort == nil && nextPort == nil {
			addToPortGroup(v)
		}

		// Start of the slice with len > 0
		if prevPort == nil && nextPort != nil {
			isGroup := *nextPort-1 == port

			if isGroup {
				// Start with a group
				addToCurrentGroup(v)
			} else {
				// Start with single item
				addToPortGroup(v)
			}

			continue
		}

		// Middle of the slice with len > 0
		if prevPort != nil && nextPort != nil {
			currentIsGroup := *prevPort+1 == port
			nextIsGroup := *nextPort-1 == port

			if currentIsGroup {
				// Maybe in the middle of a group
				addToCurrentGroup(v)

				if !nextIsGroup {
					// End of a group
					finishCurrentGroup()
				}
			} else if nextIsGroup {
				// Start of a new group
				addToCurrentGroup(v)
			} else {
				// No group at all
				addToPortGroup(v)
			}

			continue
		}

		// End of the slice with len > 0
		if prevPort != nil && nextPort == nil {
			isGroup := *prevPort+1 == port

			if isGroup {
				// End group
				addToCurrentGroup(v)
				finishCurrentGroup()
			} else {
				// End single item
				addToPortGroup(v)
			}
		}
	}

	portDisplay := []string{}
	for _, group := range portGroups {
		if len(group) == 0 {
			// Usually should not happen, but better do not crash.
			continue
		}

		first := group[0]

		hostIP := first.HostIP
		if hostIP == "" {
			hostIP = "0.0.0.0"
		}

		// Single mappings
		if len(group) == 1 {
			portDisplay = append(portDisplay,
				fmt.Sprintf(
					"%s:%d->%d/%s",
					hostIP, first.HostPort, first.ContainerPort, first.Protocol,
				),
			)
			continue
		}

		// Group mappings
		last := group[len(group)-1]
		portDisplay = append(portDisplay, formatGroup(
			fmt.Sprintf("%s/%s", hostIP, first.Protocol),
			first.HostPort, last.HostPort,
			first.ContainerPort, last.ContainerPort,
		))
	}
	return strings.Join(portDisplay, ", ")
}

func comparePorts(i, j ocicni.PortMapping) bool {
	if i.ContainerPort != j.ContainerPort {
		return i.ContainerPort < j.ContainerPort
	}

	if i.HostIP != j.HostIP {
		return i.HostIP < j.HostIP
	}

	if i.HostPort != j.HostPort {
		return i.HostPort < j.HostPort
	}

	return i.Protocol < j.Protocol
}

// formatGroup returns the group in the format:
// <IP:firstHost:lastHost->firstCtr:lastCtr/Proto>
// e.g 0.0.0.0:1000-1006->2000-2006/tcp.
func formatGroup(key string, firstHost, lastHost, firstCtr, lastCtr int32) string {
	parts := strings.Split(key, "/")
	groupType := parts[0]
	var ip string
	if len(parts) > 1 {
		ip = parts[0]
		groupType = parts[1]
	}

	group := func(first, last int32) string {
		group := strconv.Itoa(int(first))
		if first != last {
			group = fmt.Sprintf("%s-%d", group, last)
		}
		return group
	}
	hostGroup := group(firstHost, lastHost)
	ctrGroup := group(firstCtr, lastCtr)

	return fmt.Sprintf("%s:%s->%s/%s", ip, hostGroup, ctrGroup, groupType)
}
