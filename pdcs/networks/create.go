package networks

import (
	"fmt"
	"net"

	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/containers/podman/v3/libpod/network/types"
	"github.com/containers/podman/v3/pkg/bindings/network"
	"github.com/containers/podman/v3/pkg/errorhandling"
	"github.com/rs/zerolog/log"
)

// CreateOptions implements network create options.
type CreateOptions struct {
	Name           string
	Labels         map[string]string
	Internal       bool
	Macvlan        string
	Drivers        string
	DriversOptions map[string]string
	IPv6           bool
	Gateway        string
	IPRange        string
	Subnet         string
	DisableDNS     bool
}

// Create creates a new pod.
func Create(opts CreateOptions) (string, error) {
	log.Debug().Msgf("pdcs: podman network create %v", opts)

	var (
		errList  []error
		filename string
	)

	createOptions := &network.CreateOptions{
		Name:       &opts.Name,
		MacVLAN:    &opts.Macvlan,
		Labels:     opts.Labels,
		Driver:     &opts.Drivers,
		DisableDNS: &opts.DisableDNS,
		Internal:   &opts.Internal,
		IPv6:       &opts.IPv6,
	}

	if opts.Gateway != "" {
		addr := net.ParseIP(opts.Gateway)
		if addr != nil {
			createOptions.Gateway = &addr
		} else {
			errList = append(errList, fmt.Errorf("invalid gateway address: %s", opts.Gateway))
		}
	}

	if opts.Subnet != "" {
		_, ipnet, err := net.ParseCIDR(opts.Subnet)
		if err != nil {
			errList = append(errList, fmt.Errorf("invalid subnet: %s", opts.Subnet))
		} else {
			createOptions.Subnet = ipnet
		}
	}

	if opts.IPRange != "" {
		_, ipnet, err := net.ParseCIDR(opts.IPRange)
		if err != nil {
			errList = append(errList, fmt.Errorf("invalid ip range: %s", opts.IPRange))
		} else {
			createOptions.IPRange = ipnet
		}
	}

	if len(errList) > 0 {
		return filename, errorhandling.JoinErrors(errList)
	}

	conn, err := connection.GetConnection()
	if err != nil {
		return filename, err
	}

	report, err := network.Create(conn, createOptions)
	if err != nil {
		return filename, err
	}
	filename = report.Filename
	return filename, nil

}

// DefaultNetworkDriver returns default network driver name.
func DefaultNetworkDriver() string {
	return types.DefaultNetworkDriver
}
