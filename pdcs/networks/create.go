package networks

import (
	"net"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v5/pkg/bindings/network"
	"github.com/containers/podman/v5/pkg/errorhandling"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.podman.io/common/libnetwork/types"
	"go.podman.io/common/libnetwork/util"
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
func Create(opts CreateOptions) (types.Network, error) {
	log.Debug().Msgf("pdcs: podman network create %v", opts)

	var (
		errList []error
		report  types.Network
	)

	createOptions := &types.Network{
		Name:        opts.Name,
		Labels:      opts.Labels,
		Driver:      opts.Drivers,
		DNSEnabled:  !opts.DisableDNS,
		Internal:    opts.Internal,
		IPv6Enabled: opts.IPv6,
	}

	if len(opts.Subnets) > 0 { //nolint:nestif
		for i := range opts.Subnets {
			subnet, err := types.ParseCIDR(opts.Subnets[i])
			if err != nil {
				return report, err
			}

			s := types.Subnet{
				Subnet: subnet,
			}

			if len(opts.IPRanges) > i {
				leaseRange, err := parseRange(opts.IPRanges[i])
				if err != nil {
					return report, err
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
		return report, errorhandling.JoinErrors(errList)
	}

	conn, err := registry.GetConnection()
	if err != nil {
		return report, err
	}

	report, err = network.Create(conn, createOptions)
	if err != nil {
		return report, err
	}

	return report, nil
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
