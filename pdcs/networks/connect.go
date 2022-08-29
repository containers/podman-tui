package networks

import (
	"net"

	"github.com/containers/common/libnetwork/types"
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman/v4/pkg/bindings/network"
	"github.com/rs/zerolog/log"
)

// NetworkConnect networks connect options.
type NetworkConnect struct {
	Container  string
	Network    string
	IPv4       string
	IPv6       string
	MacAddress string
	Aliases    []string
}

// Connect connects a container to a network.
func Connect(opts NetworkConnect) error {
	log.Debug().Msgf("pdcs: podman network connect %v", opts)

	var networkConnectOptions types.PerNetworkOptions

	if opts.MacAddress != "" {
		mac, err := net.ParseMAC(opts.MacAddress)
		if err != nil {
			return err
		}

		networkConnectOptions.StaticMAC = types.HardwareAddr(mac)
	}

	for _, ipaddr := range []string{opts.IPv4, opts.IPv6} {
		ip := net.ParseIP(ipaddr)
		if ip != nil {
			networkConnectOptions.StaticIPs = append(networkConnectOptions.StaticIPs, ip)
		}
	}

	for _, alias := range opts.Aliases {
		if alias != "" {
			networkConnectOptions.Aliases = append(networkConnectOptions.Aliases, alias)
		}
	}

	conn, err := registry.GetConnection()
	if err != nil {
		return err
	}

	return network.Connect(conn, opts.Network, opts.Container, &networkConnectOptions)
}
