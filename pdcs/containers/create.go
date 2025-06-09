package containers

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/containers/common/libnetwork/types"
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/pdcs/utils"
	"github.com/containers/podman/v5/libpod/define"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/containers/podman/v5/pkg/specgenutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var ErrInvalidCreateTimeout = errors.New("invalid container create timeout value")

// CreateOptions container create options.
type CreateOptions struct {
	Name                  string
	Command               string
	Labels                []string
	Image                 string
	Remove                bool
	Privileged            bool
	Timeout               string
	Interactive           bool
	TTY                   bool
	Detach                bool
	Secret                []string
	WorkDir               string
	EnvVars               []string
	EnvFile               []string
	EnvMerge              []string
	UnsetEnv              []string
	EnvHost               bool
	UnsetEnvAll           bool
	Umask                 string
	User                  string
	HostUsers             []string
	GroupEntry            string
	PasswdEntry           string
	Pod                   string
	Hostname              string
	IPAddress             string
	Network               string
	MacAddress            string
	Publish               []string
	Expose                []string
	PublishAll            bool
	DNSServer             []string
	DNSOptions            []string
	DNSSearchDomain       []string
	Volume                string
	ImageVolume           string
	Mount                 string
	SelinuxOpts           []string
	ApparmorProfile       string
	Seccomp               string
	SecNoNewPriv          bool
	SecMask               string
	SecUnmask             string
	HealthCmd             string
	HealthInterval        string
	HealthRetries         string
	HealthStartPeroid     string
	HealthTimeout         string
	HealthOnFailure       string
	HealthStartupCmd      string
	HealthStartupInterval string
	HealthStartupRetries  string
	HealthStartupSuccess  string
	HealthStartupTimeout  string
	HealthLogDestination  string
	HealthMaxLogSize      string
	HealthMaxLogCount     string
	Memory                string
	MemoryReservation     string
	MemorySwap            string
	MemorySwappiness      string
	CPUs                  string
	CPUShares             string
	CPUPeriod             string
	CPUQuota              string
	CPURtPeriod           string
	CPURtRuntime          string
	CPUSetCPUs            string
	CPUSetMems            string
	SHMSize               string
	SHMSizeSystemd        string
	NamespaceCgroup       string
	NamespacePid          string
	NamespaceIpc          string
	NamespaceUser         string
	NamespaceUts          string
	NamespaceUidmap       string
	NamespaceSubuidName   string
	NamespaceGidmap       string
	NamespaceSubgidName   string
}

// Create creates a new container.
func Create(opts CreateOptions, run bool) ([]string, string, error) { //nolint:cyclop,gocognit,gocyclo,maintidx
	var (
		warningResponse []string
		containerID     string
		createOptions   entities.ContainerCreateOptions
	)

	log.Debug().Msgf("pdcs: podman container create %v", opts)
	utils.DefineCreateDefaults(&createOptions)

	conn, err := registry.GetConnection()
	if err != nil {
		return warningResponse, containerID, err
	}

	if len(opts.Labels) > 0 {
		createOptions.Label = opts.Labels
	}

	createOptions.Name = opts.Name
	createOptions.Rm = opts.Remove
	createOptions.Privileged = opts.Privileged

	if opts.Timeout != "" {
		timeout, err := strconv.Atoi(opts.Timeout)
		if err != nil {
			return warningResponse, containerID, fmt.Errorf("%w: %s", ErrInvalidCreateTimeout, opts.Timeout)
		}

		createOptions.Timeout = uint(timeout) //nolint:gosec
	}

	createOptions.Hostname = opts.Hostname

	if len(opts.Expose) > 0 {
		createOptions.Expose = opts.Expose
	}

	createOptions.PublishAll = opts.PublishAll

	createOptions.Net, err = containerNetworkOptions(opts)
	if err != nil {
		return warningResponse, containerID, err
	}

	if opts.Pod != "" {
		createOptions.Pod = opts.Pod
		createOptions.Net.Network.NSMode = specgen.FromPod
	} else {
		createOptions.Net.Network.NSMode = specgen.Default
	}

	if opts.Volume != "" {
		createOptions.Volume = strings.Split(opts.Volume, ",")
	}

	if opts.Mount != "" {
		for _, mopts := range strings.Split(opts.Mount, " ") {
			if mopts != "" {
				createOptions.Mount = append(createOptions.Mount, mopts)
			}
		}
	}

	createOptions.ImageVolume = opts.ImageVolume

	// security options
	if opts.ApparmorProfile != "" {
		createOptions.SecurityOpt = append(createOptions.SecurityOpt, "apparmor="+opts.ApparmorProfile)
	}

	if opts.SecMask != "" {
		createOptions.SecurityOpt = append(createOptions.SecurityOpt, "mask="+opts.SecMask)
	}

	if opts.SecUnmask != "" {
		createOptions.SecurityOpt = append(createOptions.SecurityOpt, "unmask="+opts.SecUnmask)
	}

	if opts.Seccomp != "" {
		createOptions.SeccompPolicy = opts.Seccomp
	}

	if opts.SecNoNewPriv {
		createOptions.SecurityOpt = append(createOptions.SecurityOpt, "no-new-privileges")
	}

	if len(opts.SelinuxOpts) > 0 {
		for _, selinuxLabel := range opts.SelinuxOpts {
			createOptions.SecurityOpt = append(createOptions.SecurityOpt, "label="+selinuxLabel)
		}
	}

	// environment options
	if opts.WorkDir != "" {
		createOptions.Workdir = opts.WorkDir
	}

	if len(opts.EnvVars) > 0 {
		createOptions.Env = opts.EnvVars
	}

	if len(opts.EnvFile) > 0 {
		createOptions.EnvFile = opts.EnvFile
	}

	if len(opts.EnvMerge) > 0 {
		createOptions.EnvMerge = opts.EnvMerge
	}

	if len(opts.UnsetEnv) > 0 {
		createOptions.UnsetEnv = opts.UnsetEnv
	}

	createOptions.EnvHost = opts.EnvHost
	createOptions.UnsetEnvAll = opts.UnsetEnvAll

	if opts.Umask != "" {
		createOptions.Umask = opts.Umask
	}

	// user and groups
	if opts.User != "" {
		createOptions.User = opts.User
	}

	if len(opts.HostUsers) > 0 {
		createOptions.HostUsers = opts.HostUsers
	}

	if opts.PasswdEntry != "" {
		createOptions.PasswdEntry = opts.PasswdEntry
	}

	if opts.GroupEntry != "" {
		createOptions.GroupEntry = opts.GroupEntry
	}

	// add secrets
	if len(opts.Secret) > 0 {
		createOptions.Secrets = opts.Secret
	}

	if run {
		createOptions.Interactive = opts.Interactive
		createOptions.TTY = opts.TTY
	}

	// add healthcheck options
	if err := containerHealthOptions(&createOptions, opts); err != nil {
		return warningResponse, containerID, err
	}

	// add resources
	if opts.Memory != "" {
		createOptions.Memory = opts.Memory
	}

	if opts.MemoryReservation != "" {
		createOptions.MemoryReservation = opts.MemoryReservation
	}

	if opts.MemorySwap != "" {
		createOptions.MemorySwap = opts.MemorySwap
	}

	if opts.MemorySwappiness != "" {
		val, err := strconv.Atoi(opts.MemorySwappiness)
		if err != nil {
			return warningResponse, containerID, err
		}

		createOptions.MemorySwappiness = int64(val)
	}

	if opts.CPUs != "" {
		val, err := strconv.ParseFloat(opts.CPUs, 64)
		if err != nil {
			return warningResponse, containerID, err
		}

		createOptions.CPUS = val
	}

	if opts.CPUShares != "" {
		val, err := strconv.ParseUint(opts.CPUShares, 10, 64)
		if err != nil {
			return warningResponse, containerID, err
		}

		createOptions.CPUShares = val
	}

	if opts.CPUPeriod != "" {
		val, err := strconv.ParseUint(opts.CPUPeriod, 10, 64)
		if err != nil {
			return warningResponse, containerID, err
		}

		createOptions.CPUPeriod = val
	}

	if opts.CPURtPeriod != "" {
		val, err := strconv.ParseUint(opts.CPURtPeriod, 10, 64)
		if err != nil {
			return warningResponse, containerID, err
		}

		createOptions.CPURTPeriod = val
	}

	if opts.CPUQuota != "" {
		val, err := strconv.Atoi(opts.CPUQuota)
		if err != nil {
			return warningResponse, containerID, err
		}

		createOptions.CPUQuota = int64(val)
	}

	if opts.CPURtRuntime != "" {
		val, err := strconv.Atoi(opts.CPURtRuntime)
		if err != nil {
			return warningResponse, containerID, err
		}

		createOptions.CPURTRuntime = int64(val)
	}

	if opts.CPUSetCPUs != "" {
		createOptions.CPUSetCPUs = opts.CPUSetCPUs
	}

	if opts.CPUSetMems != "" {
		createOptions.CPUSetMems = opts.CPUSetMems
	}

	if opts.SHMSize != "" {
		createOptions.ShmSize = opts.SHMSize
	}

	if opts.SHMSizeSystemd != "" {
		createOptions.ShmSizeSystemd = opts.SHMSizeSystemd
	}

	// namespace options
	if opts.NamespaceCgroup != "" {
		createOptions.CgroupNS = opts.NamespaceCgroup
	}

	if opts.NamespaceIpc != "" {
		createOptions.IPC = opts.NamespaceIpc
	}

	if opts.NamespacePid != "" {
		createOptions.PID = opts.NamespacePid
	}

	if opts.NamespaceUser != "" {
		createOptions.UserNS = opts.NamespaceUser
	}

	if opts.NamespaceUts != "" {
		createOptions.UTS = opts.NamespaceUts
	}

	if opts.NamespaceUidmap != "" {
		createOptions.UIDMap = []string{opts.NamespaceUidmap}
	}

	if opts.NamespaceSubuidName != "" {
		createOptions.SubUIDName = opts.NamespaceSubuidName
	}

	if opts.NamespaceGidmap != "" {
		createOptions.GIDMap = []string{opts.NamespaceGidmap}
	}

	if opts.NamespaceSubgidName != "" {
		createOptions.SubGIDName = opts.NamespaceSubgidName
	}

	// generate spec
	s := specgen.NewSpecGenerator(opts.Name, false)
	if err := specgenutil.FillOutSpecGen(s, &createOptions, nil); err != nil {
		return warningResponse, containerID, err
	}

	// container image
	s.Image = opts.Image

	// command
	cmd := strings.TrimSpace(opts.Command)
	if cmd != "" {
		s.Command = strings.Split(cmd, " ")
	}

	// validate spec
	if err := s.Validate(); err != nil {
		return warningResponse, containerID, err
	}

	response, err := containers.CreateWithSpec(conn, s, &containers.CreateOptions{})
	if err != nil {
		return warningResponse, containerID, err
	}

	warningResponse = response.Warnings
	containerID = response.ID

	return warningResponse, containerID, nil
}

func containerHealthOptions(createOptions *entities.ContainerCreateOptions, opts CreateOptions) error { //nolint:cyclop
	createOptions.HealthInterval = define.DefaultHealthCheckInterval
	createOptions.StartupHCInterval = define.DefaultHealthCheckInterval
	createOptions.HealthRetries = define.DefaultHealthCheckRetries
	createOptions.HealthStartPeriod = define.DefaultHealthCheckStartPeriod
	createOptions.HealthTimeout = define.DefaultHealthCheckTimeout
	createOptions.StartupHCTimeout = define.DefaultHealthCheckTimeout
	createOptions.HealthOnFailure = opts.HealthOnFailure
	createOptions.HealthLogDestination = opts.HealthLogDestination

	if opts.HealthCmd == "" {
		createOptions.HealthCmd = "none"

		return nil
	}

	createOptions.HealthCmd = opts.HealthCmd

	if opts.HealthMaxLogCount != "" {
		logCount, err := strconv.ParseUint(opts.HealthMaxLogCount, 10, 32)
		if err != nil {
			return err
		}

		logCountWd := uint(logCount)
		createOptions.HealthMaxLogCount = logCountWd
	}

	if opts.HealthMaxLogSize != "" {
		logSize, err := strconv.ParseUint(opts.HealthMaxLogSize, 10, 32)
		if err != nil {
			return err
		}

		logSizeWd := uint(logSize)
		createOptions.HealthMaxLogSize = logSizeWd
	}

	if opts.HealthInterval != "" {
		createOptions.HealthInterval = opts.HealthInterval
	}

	if opts.HealthStartPeroid != "" {
		createOptions.HealthStartPeriod = opts.HealthStartPeroid
	}

	if opts.HealthTimeout != "" {
		createOptions.HealthTimeout = opts.HealthTimeout
	}

	if opts.HealthStartupCmd != "" {
		createOptions.StartupHCCmd = opts.HealthStartupCmd
	}

	if opts.HealthStartupInterval != "" {
		createOptions.StartupHCInterval = opts.HealthStartupInterval
	}

	if opts.HealthStartupTimeout != "" {
		createOptions.StartupHCTimeout = opts.HealthStartupTimeout
	}

	if opts.HealthRetries != "" {
		retries, err := strconv.ParseUint(opts.HealthRetries, 10, 32)
		if err != nil {
			return err
		}

		retriesWd := uint(retries)
		createOptions.HealthRetries = retriesWd
	}

	if opts.HealthStartupRetries != "" {
		startupRetries, err := strconv.ParseUint(opts.HealthStartupRetries, 10, 32)
		if err != nil {
			return err
		}

		startupRetriesWd := uint(startupRetries)
		createOptions.StartupHCRetries = startupRetriesWd
	}

	if opts.HealthStartupSuccess != "" {
		startupSuccess, err := strconv.ParseUint(opts.HealthStartupSuccess, 10, 32)
		if err != nil {
			return err
		}

		startupSuccessWd := uint(startupSuccess)
		createOptions.StartupHCRetries = startupSuccessWd
	}

	return nil
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
