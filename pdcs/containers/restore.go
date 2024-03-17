package containers

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/pdcs/utils"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/rs/zerolog/log"
)

// CntRestoreOptions is container restore options.
type CntRestoreOptions struct {
	ContainerID     string
	PodID           string
	Name            string
	Publish         []string
	Import          string
	Keep            bool
	IgnoreStaticIP  bool
	IgnoreStaticMAC bool
	FileLocks       bool
	PrintStats      bool
	TCPEstablished  bool
	IgnoreVolumes   bool
	IgnoreRootfs    bool
}

// Restore performs container checkpoint restore.
func Restore(opts CntRestoreOptions) (string, error) {
	log.Debug().Msgf("pdcs: podman container restore %v", opts)

	conn, err := registry.GetConnection()
	if err != nil {
		return "", err
	}

	restoreOptions := new(containers.RestoreOptions)

	if opts.PodID != "" {
		restoreOptions.WithPod(opts.PodID)
	}

	if opts.Name != "" {
		restoreOptions.WithName(opts.Name)
	}

	if len(opts.Publish) != 0 {
		restoreOptions.WithPublishPorts(opts.Publish)
	}

	if opts.Import != "" {
		restoreOptions.WithImportArchive(opts.Import)
	}

	restoreOptions.WithKeep(opts.Keep)
	restoreOptions.WithIgnoreStaticIP(opts.IgnoreStaticIP)
	restoreOptions.WithIgnoreStaticMAC(opts.IgnoreStaticMAC)
	restoreOptions.WithFileLocks(opts.FileLocks)
	restoreOptions.WithPrintStats(opts.PrintStats)
	restoreOptions.WithTCPEstablished(opts.TCPEstablished)
	restoreOptions.WithIgnoreVolumes(opts.IgnoreVolumes)
	restoreOptions.WithIgnoreRootfs(opts.IgnoreRootfs)

	response, err := containers.Restore(conn, opts.ContainerID, restoreOptions)
	if err != nil {
		return "", err
	}

	report, err := utils.GetJSONOutput(response)
	if err != nil {
		return "", err
	}

	log.Debug().Msgf("pdcs: restore %s", report)

	return report, nil
}
