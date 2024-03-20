package pods

import (
	"encoding/json"
	"net"

	"github.com/containers/common/libnetwork/types"
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/pdcs/utils"
	"github.com/containers/podman/v5/pkg/bindings/pods"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/containers/podman/v5/pkg/errorhandling"
	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/containers/podman/v5/pkg/specgenutil"
	"github.com/containers/podman/v5/pkg/util"
	"github.com/pkg/errors"
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
	InfraCommand    string
	InfraImage      string
	Hostname        string
	IPAddress       string
	MacAddress      string
	AddHost         []string
	Network         string
	Publish         []string
	SecurityOpts    []string
}

// Create creates a new pod.
func Create(opts CreateOptions) error { //nolint:cyclop
	log.Debug().Msgf("pdcs: podman pod create %v", opts)

	var createOptions entities.PodCreateOptions

	var (
		infraOptions = entities.NewInfraContainerCreateOptions()
		errList      = make([]error, 0)
	)

	utils.DefineCreateDefaults(&infraOptions)

	conn, err := registry.GetConnection()
	if err != nil {
		return err
	}

	createOptions.Name = opts.Name
	createOptions.Labels = opts.Labels

	// network options
	podNetworkOptions, err := podNetworkOptions(opts)
	if err != nil {
		return err
	}

	createOptions.Infra = opts.Infra

	if createOptions.Infra {
		if opts.InfraImage != "" {
			createOptions.InfraImage = opts.InfraImage
		} else {
			createOptions.InfraImage = defaultPodInfraImage()
		}

		infraOptions.Net = podNetworkOptions
		createOptions.InfraCommand = &(opts.InfraCommand)

		err = containerToPodOptions(&infraOptions, &createOptions)
		if err != nil {
			return err
		}
	} else {
		createOptions.Share = nil
		createOptions.Net = podNetworkOptions
	}

	createOptions.Hostname = opts.Hostname
	createOptions.SecurityOpt = opts.SecurityOpts

	podSpec := specgen.NewPodSpecGenerator()

	podSpec, err = entities.ToPodSpecGen(*podSpec, &createOptions)
	if err != nil {
		return err
	}

	if createOptions.Infra {
		imageName := opts.InfraImage
		podSpec.InfraContainerSpec = specgen.NewSpecGenerator(imageName, false)
		podSpec.InfraContainerSpec.RawImageName = imageName

		err = specgenutil.FillOutSpecGen(podSpec.InfraContainerSpec, &infraOptions, []string{})
		if err != nil {
			return err
		}

		podSpec.Volumes = podSpec.InfraContainerSpec.Volumes
		podSpec.ImageVolumes = podSpec.InfraContainerSpec.ImageVolumes
		podSpec.OverlayVolumes = podSpec.InfraContainerSpec.OverlayVolumes
		podSpec.Mounts = podSpec.InfraContainerSpec.Mounts

		wrapped, err := json.Marshal(podSpec.InfraContainerSpec)
		if err != nil {
			return err
		}

		err = json.Unmarshal(wrapped, podSpec)
		if err != nil {
			return err
		}
	}

	// validate spec
	if err := podSpec.Validate(); err != nil {
		errList = append(errList, err)
	}

	if len(errList) > 0 {
		return errorhandling.JoinErrors(errList)
	}

	newPodSpec := entities.PodSpec{PodSpecGen: *podSpec}

	_, err = pods.CreatePodFromSpec(conn, &newPodSpec)
	if err != nil {
		return err
	}

	return nil
}

func defaultPodInfraImage() string {
	containerConfig := util.DefaultContainerConfig()

	return containerConfig.Engine.InfraImage
}

func podNetworkOptions(opts CreateOptions) (*entities.NetOptions, error) { //nolint:cyclop
	var (
		err           error
		perNetworkOpt types.PerNetworkOptions
	)

	netOptions := &entities.NetOptions{}
	netOptions.Networks = make(map[string]types.PerNetworkOptions)

	if len(opts.AddHost) > 0 {
		netOptions.AddHosts = opts.AddHost
	}

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

	netOptions.NoHosts = opts.NoHost

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

func containerToPodOptions(containerCreate *entities.ContainerCreateOptions, podCreate *entities.PodCreateOptions) error { //nolint:lll
	contMarshal, err := json.Marshal(containerCreate)
	if err != nil {
		return err
	}

	return json.Unmarshal(contMarshal, podCreate)
}
