package networks

import (
	"net"

	"github.com/containers/common/libnetwork/types"
	"github.com/containers/common/libnetwork/util"
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/containers/podman/v4/pkg/bindings/network"
	"github.com/containers/podman/v4/pkg/errorhandling"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// CreateOptions implements network create options.
type CreateOptions struct {
	Name           string
	Labels         map[string]string
	Internal       bool
	Drivers        string
	DriversOptions map[string]string
	IPv6           bool
	Gateways       []string
	IPRanges       []string
	Subnets        []string
	DisableDNS     bool
}

// Create creates a new pod.
func Create(opts CreateOptions) (string, error) {
	log.Debug().Msgf("pdcs: podman network create %v", opts)

	var (
		errList  []error
		filename string
	)

	createOptions := &types.Network{
		Name:        opts.Name,
		Labels:      opts.Labels,
		Driver:      opts.Drivers,
		DNSEnabled:  !opts.DisableDNS,
		Internal:    opts.Internal,
		IPv6Enabled: opts.IPv6,
	}

	if len(opts.Subnets) > 0 {
		for i := range opts.Subnets {
			subnet, err := types.ParseCIDR(opts.Subnets[i])
			if err != nil {
				return "", err
			}
			s := types.Subnet{
				Subnet: subnet,
			}
			if len(opts.IPRanges) > i {
				leaseRange, err := parseRange(opts.IPRanges[i])
				if err != nil {
					return "", err
				}
				s.LeaseRange = leaseRange
			}
			if len(opts.Gateways) > i {
				s.Gateway = net.ParseIP(opts.Gateways[i])
			}
			createOptions.Subnets = append(createOptions.Subnets, s)
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
	return report.Name, nil

}

// DefaultNetworkDriver returns default network driver name.
func DefaultNetworkDriver() string {
	return types.DefaultNetworkDriver
}

func parseRange(iprange string) (*types.LeaseRange, error) {
	_, subnet, err := net.ParseCIDR(iprange)
	if err != nil {
		return nil, err
	}

	startIP, err := util.FirstIPInSubnet(subnet)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get first ip in range")
	}
	lastIP, err := util.LastIPInSubnet(subnet)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get last ip in range")
	}
	return &types.LeaseRange{
		StartIP: startIP,
		EndIP:   lastIP,
	}, nil
}
