package containers

import (
	"net"

	"github.com/containers/common/libnetwork/types"
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/pdcs/utils"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/containers/podman/v4/pkg/specgen"
	"github.com/containers/podman/v4/pkg/specgenutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// CreateOptions container create options.
type CreateOptions struct {
	Name            string
	Labels          []string
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
	SelinuxOpts     []string
	ApparmorProfile string
	Seccomp         string
	NoNewPriv       bool
	Mask            string
	Unmask          string
}

// Create creates a new container.
func Create(opts CreateOptions) ([]string, error) { //nolint:cyclop
	var (
		warningResponse []string
		createOptions   entities.ContainerCreateOptions
	)

	log.Debug().Msgf("pdcs: podman container create %v", opts)
	utils.DefineCreateDefaults(&createOptions)

	conn, err := registry.GetConnection()
	if err != nil {
		return warningResponse, err
	}

	if len(opts.Labels) > 0 {
		createOptions.Label = opts.Labels
	}

	createOptions.Name = opts.Name
	createOptions.Rm = opts.Remove

	createOptions.Hostname = opts.Hostname

	if len(opts.Expose) > 0 {
		createOptions.Expose = opts.Expose
	}

	createOptions.PublishAll = opts.PublishAll

	createOptions.Net, err = containerNetworkOptions(opts)
	if err != nil {
		return warningResponse, err
	}

	if opts.Pod != "" {
		createOptions.Pod = opts.Pod
		createOptions.Net.Network.NSMode = specgen.FromPod
	} else {
		createOptions.Net.Network.NSMode = specgen.Default
	}

	if opts.Volume != "" {
		createOptions.Volume = []string{opts.Volume}
	}

	createOptions.ImageVolume = opts.ImageVolume

	// security options
	if opts.ApparmorProfile != "" {
		createOptions.SecurityOpt = append(createOptions.SecurityOpt, "apparmor="+opts.ApparmorProfile)
	}

	if opts.Mask != "" {
		createOptions.SecurityOpt = append(createOptions.SecurityOpt, "mask="+opts.Mask)
	}

	if opts.Unmask != "" {
		createOptions.SecurityOpt = append(createOptions.SecurityOpt, "unmask="+opts.Unmask)
	}

	if opts.Seccomp != "" {
		createOptions.SeccompPolicy = opts.Seccomp
	}

	if opts.NoNewPriv {
		createOptions.SecurityOpt = append(createOptions.SecurityOpt, "no-new-privileges")
	}

	if len(opts.SelinuxOpts) > 0 {
		for _, selinuxLabel := range opts.SelinuxOpts {
			createOptions.SecurityOpt = append(createOptions.SecurityOpt, "label="+selinuxLabel)
		}
	}

	s := specgen.NewSpecGenerator(opts.Name, false)
	if err := specgenutil.FillOutSpecGen(s, &createOptions, nil); err != nil {
		return warningResponse, err
	}

	s.Image = opts.Image

	// validate spec
	if err := s.Validate(); err != nil {
		return warningResponse, err
	}

	response, err := containers.CreateWithSpec(conn, s, &containers.CreateOptions{})
	if err != nil {
		return warningResponse, err
	}

	warningResponse = response.Warnings

	return warningResponse, nil
}

func containerNetworkOptions(opts CreateOptions) (*entities.NetOptions, error) { //nolint:cyclop
	var (
		err           error
		perNetworkOpt types.PerNetworkOptions
	)

	netOptions := &entities.NetOptions{}
	netOptions.Networks = make(map[string]types.PerNetworkOptions)

	var dnsServers []net.IP

	for _, d := range opts.DNSServer {
		addr := net.ParseIP(d)
		if addr != nil {
			dnsServers = append(dnsServers, addr)

			continue
		}

		return nil, errors.Wrap(utils.ErrInvalidDNSAddress, d)
	}

	if len(dnsServers) > 0 {
		netOptions.DNSServers = dnsServers
		netOptions.DNSOptions = opts.DNSOptions
		netOptions.DNSSearch = opts.DNSSearchDomain
	}

	if len(opts.Publish) > 0 {
		netOptions.PublishPorts, err = specgenutil.CreatePortBindings(opts.Publish)
		if err != nil {
			return nil, err
		}
	}

	if opts.Network != "" { //nolint:nestif
		if opts.MacAddress != "" {
			mac, err := net.ParseMAC(opts.MacAddress)
			if err != nil {
				return nil, err
			}

			perNetworkOpt.StaticMAC = types.HardwareAddr(mac)
		}

		if opts.IPAddress != "" {
			addr := net.ParseIP(opts.IPAddress)

			if addr == nil {
				return nil, errors.Wrap(utils.ErrInvalidIPAddress, opts.IPAddress)
			}

			perNetworkOpt.StaticIPs = []net.IP{addr}
		}

		netOptions.Networks[opts.Network] = perNetworkOpt
	}

	return netOptions, nil
}
