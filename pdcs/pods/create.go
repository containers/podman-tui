package pods

import (
	"fmt"
	"net"

	"github.com/containers/podman/v3/pkg/bindings/pods"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/containers/podman/v3/pkg/errorhandling"
	"github.com/containers/podman/v3/pkg/specgen"
	"github.com/containers/podman/v3/pkg/util"
	"github.com/containers/podman-tui/pdcs/connection"
	"github.com/rs/zerolog/log"
)

// CreateOptions implements pods create spec options.
type CreateOptions struct {
	Name            string
	NoHost          bool
	Labels          map[string]string
	DNSServer       []string
	DNSOptions      []string
	DNSSearchDomain []string
	Infra           bool
	InfraCommand    []string
	InfraImage      string
	Hostname        string
	IPAddress       string
	MacAddress      string
	HostToIP        []string
	Network         string
	//NetworkAliases  []string
}

// Create creates a new pod.
func Create(opts CreateOptions) error {
	log.Debug().Msgf("pdcs: podman pod create %v", opts)
	var errList []error
	conn, err := connection.GetConnection()
	if err != nil {
		return err
	}
	podSpecGenerator := specgen.NewPodSpecGenerator()
	podSpecGenerator.Name = opts.Name
	podSpecGenerator.Labels = opts.Labels
	podSpecGenerator.NoManageHosts = opts.NoHost
	podSpecGenerator.DNSOption = opts.DNSOptions
	podSpecGenerator.DNSSearch = opts.DNSSearchDomain
	podSpecGenerator.NoInfra = !opts.Infra
	podSpecGenerator.InfraCommand = opts.InfraCommand
	podSpecGenerator.InfraImage = opts.InfraImage
	podSpecGenerator.Hostname = opts.Hostname
	podSpecGenerator.HostAdd = opts.HostToIP
	if opts.Network != "" {
		podSpecGenerator.CNINetworks = []string{opts.Network}
	}

	if opts.MacAddress != "" {
		mac, err := net.ParseMAC(opts.MacAddress)
		if err != nil {
			errList = append(errList, err)
		} else {
			podSpecGenerator.StaticMAC = &mac
		}

	}
	if opts.IPAddress != "" {
		addr := net.ParseIP(opts.IPAddress)
		if addr != nil {
			podSpecGenerator.StaticIP = &addr
		} else {
			errList = append(errList, fmt.Errorf("invalid ip address: %s", opts.IPAddress))
		}
	}

	var dnsServers []net.IP
	for _, d := range opts.DNSServer {
		addr := net.ParseIP(d)
		if addr != nil {
			dnsServers = append(dnsServers, addr)
			continue
		}
		errList = append(errList, fmt.Errorf("invalid dns server: %s", d))
	}
	if len(dnsServers) > 0 {
		podSpecGenerator.DNSServer = dnsServers
	}

	if err := podSpecGenerator.Validate(); err != nil {
		errList = append(errList, err)
	}

	if len(errList) > 0 {
		return errorhandling.JoinErrors(errList)
	}

	podSpec := entities.PodSpec{
		PodSpecGen: *podSpecGenerator,
	}

	_, err = pods.CreatePodFromSpec(conn, &podSpec)
	if err != nil {
		return err
	}
	return nil
}

// DefaultPodInfraImage returns default infra container image.
func DefaultPodInfraImage() string {
	containerConfig := util.DefaultContainerConfig()
	return containerConfig.Engine.InfraImage
}
