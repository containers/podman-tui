package containers

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func createExpose(expose []string) (map[uint16]string, error) {
	toReturn := make(map[uint16]string)

	for _, e := range expose {
		// Check for protocol
		proto := "tcp"
		splitProto := strings.Split(e, "/")

		if len(splitProto) > 2 { // nolint:gomnd
			return nil, errors.Errorf("invalid expose format - protocol can only be specified once")
		} else if len(splitProto) == 2 { // nolint:gomnd
			proto = splitProto[1]
		}

		// Check for a range
		start, rangeLen, err := parseAndValidateRange(splitProto[0])
		if err != nil {
			return nil, err
		}

		var index uint16
		for index = 0; index < rangeLen; index++ {
			portNum := start + index
			protocols, ok := toReturn[portNum]

			if !ok {
				toReturn[portNum] = proto
			} else {
				newProto := strings.Join(append(strings.Split(protocols, ","), strings.Split(proto, ",")...), ",")
				toReturn[portNum] = newProto
			}
		}
	}

	return toReturn, nil
}

func parseAndValidateRange(portRange string) (uint16, uint16, error) {
	splitRange := strings.Split(portRange, "-")
	if len(splitRange) > 2 { // nolint:gomnd
		return 0, 0, errors.Errorf("invalid port format - port ranges are formatted as startPort-stopPort")
	}

	if splitRange[0] == "" {
		return 0, 0, errors.Errorf("port numbers cannot be negative")
	}

	startPort, err := parseAndValidatePort(splitRange[0])
	if err != nil {
		return 0, 0, err
	}

	var rangeLen uint16 = 1

	if len(splitRange) == 2 { // nolint:gomnd
		if splitRange[1] == "" {
			return 0, 0, errors.Errorf("must provide ending number for port range")
		}

		endPort, err := parseAndValidatePort(splitRange[1])
		if err != nil {
			return 0, 0, err
		}

		if endPort <= startPort {
			return 0, 0, errors.Errorf("the end port of a range must be higher than the start port - %d is not higher than %d", endPort, startPort) // nolint:lll
		}
		// Our range is the total number of ports
		// involved, so we need to add 1 (8080:8081 is
		// 2 ports, for example, not 1)
		rangeLen = endPort - startPort + 1
	}

	return startPort, rangeLen, nil
}

// Turn a single string into a valid U16 port.
func parseAndValidatePort(port string) (uint16, error) {
	num, err := strconv.Atoi(port)
	if err != nil {
		return 0, errors.Wrapf(err, "invalid port number")
	}

	if num < 1 || num > 65535 {
		return 0, errors.Errorf("port numbers must be between 1 and 65535 (inclusive), got %d", num)
	}

	return uint16(num), nil
}
