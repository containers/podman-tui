package containers

import (
	"fmt"
	"net"

	"github.com/containers/common/libnetwork/types"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/specgen"
	"github.com/containers/podman/v4/pkg/specgenutil"
	"github.com/rs/zerolog/log"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/pdcs/volumes"
)

//CreateOptions container create options.
type CreateOptions struct {
	Name            string
	Labels          map[string]string
	Image           string
	Remove          bool
	Pod             string
	Hostname        string
	IPAddress       string
	Network         string
	MacAddress      string
	Publish         []string
	Expose          []string
	PublishAll      bool
	DNSServer       []string
	DNSOptions      []string
	DNSSearchDomain []string
	Volume          string
	ImageVolume     string
}

//Create creates a new container.
func Create(opts CreateOptions) ([]string, error) {
	var (
		warningResponse []string
		macAddress      net.HardwareAddr
		ipAddr          net.IP
		dnsServers      []net.IP
		networks        = make(map[string]types.PerNetworkOptions)
	)
	log.Debug().Msgf("pdcs: podman container create %v", opts)
	conn, err := registry.GetConnection()
	if err != nil {
		return warningResponse, err
	}
	containerSpecGen := specgen.NewSpecGenerator(opts.Name, false)
	containerSpecGen.Name = opts.Name
	if opts.Pod != "" {
		containerSpecGen.Pod = opts.Pod
	}
	containerSpecGen.Image = opts.Image
	containerSpecGen.Labels = opts.Labels
	containerSpecGen.Remove = opts.Remove
	containerSpecGen.Hostname = opts.Hostname

	var perNetworkOpt types.PerNetworkOptions
	if opts.MacAddress != "" {
		macAddress, err = net.ParseMAC(opts.MacAddress)
		if err != nil {
			return warningResponse, err
		}
		perNetworkOpt.StaticMAC = types.HardwareAddr(macAddress)
	}

	if opts.IPAddress != "" {
		ipAddr = net.ParseIP(opts.IPAddress)
		if ipAddr == nil {
			return warningResponse, fmt.Errorf("invalid IP address: %s", ipAddr)
		}
		perNetworkOpt.StaticIPs = []net.IP{ipAddr}
	}
	if opts.Network != "" {
		networks[opts.Network] = perNetworkOpt
	}
	containerSpecGen.Networks = networks

	for _, d := range opts.DNSServer {
		addr := net.ParseIP(d)
		if addr == nil {
			return warningResponse, fmt.Errorf("invalid DNS address: %s", ipAddr)
		}
		dnsServers = append(dnsServers, addr)
	}
	if len(dnsServers) > 0 {
		containerSpecGen.DNSServers = dnsServers
	}
	if len(opts.DNSOptions) > 0 {
		containerSpecGen.DNSOptions = opts.DNSOptions
	}
	if len(opts.DNSSearchDomain) > 0 {
		containerSpecGen.DNSSearch = opts.DNSSearchDomain
	}

	// ports
	if len(opts.Publish) > 0 {
		containerSpecGen.PortMappings, err = specgenutil.CreatePortBindings(opts.Publish)
		if err != nil {
			return warningResponse, err
		}
	}
	if len(opts.Expose) > 0 {
		containerSpecGen.Expose, err = createExpose(opts.Expose)
		if err != nil {
			return warningResponse, err
		}
	}
	containerSpecGen.PublishExposedPorts = opts.PublishAll

	// volume
	if opts.ImageVolume != "" {
		containerSpecGen.ImageVolumeMode = opts.ImageVolume
	}
	if opts.Volume != "" {
		// get volumes dest
		volumeDest, err := volumes.VolumeDest(opts.Volume)
		if err != nil {
			return warningResponse, err
		}
		containerSpecGen.Volumes = append(containerSpecGen.Volumes, &specgen.NamedVolume{
			Name: opts.Volume,
			Dest: volumeDest,
		})
	}

	// validate spec
	if err := containerSpecGen.Validate(); err != nil {
		return warningResponse, err
	}
	response, err := containers.CreateWithSpec(conn, containerSpecGen, &containers.CreateOptions{})
	if err != nil {
		return warningResponse, err
	}
	warningResponse = response.Warnings
	return warningResponse, nil
}
