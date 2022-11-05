package containers

import (
	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/pdcs/utils"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/rs/zerolog/log"
)

// CntCheckPointOptions is container checkpoint options.
type CntCheckPointOptions struct {
	Export         string
	CreateImage    string
	IgnoreRootFs   bool
	Keep           bool
	LeaveRunning   bool
	TCPEstablished bool
	PrintStats     bool
	PreCheckpoint  bool
	WithPrevious   bool
	FileLocks      bool
}

// Checkpoint a running container.
func Checkpoint(id string, opts CntCheckPointOptions) (string, error) {
	log.Debug().Msgf("pdcs: podman container checkpoint %s", id)

	conn, err := registry.GetConnection()
	if err != nil {
		return "", err
	}

	checkpointOptions := new(containers.CheckpointOptions)

	if opts.Export != "" {
		checkpointOptions.WithExport(opts.Export)
	}

	if opts.CreateImage != "" {
		checkpointOptions.WithCreateImage(opts.CreateImage)
	}

	checkpointOptions.WithIgnoreRootfs(opts.IgnoreRootFs)
	checkpointOptions.WithKeep(opts.Keep)
	checkpointOptions.WithLeaveRunning(opts.LeaveRunning)
	checkpointOptions.WithTCPEstablished(opts.TCPEstablished)
	checkpointOptions.WithPrintStats(opts.PrintStats)
	checkpointOptions.WithPreCheckpoint(opts.PreCheckpoint)
	checkpointOptions.WithWithPrevious(opts.WithPrevious)
	checkpointOptions.WithFileLocks(opts.FileLocks)

	response, err := containers.Checkpoint(conn, id, checkpointOptions)
	if err != nil {
		return "", err
	}

	report, err := utils.GetJSONOutput(response)
	if err != nil {
		return "", err
	}

	log.Debug().Msgf("pdcs: checkpoint %s", report)

	return report, nil
}
