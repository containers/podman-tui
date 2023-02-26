package utils

import (
	"fmt"

	"github.com/containers/podman/v4/libpod/define"
	"github.com/containers/podman/v4/pkg/domain/entities"
)

const (
	// DefaultContainerDetachKeys container's default attach keys string.
	DefaultContainerDetachKeys = "ctrl-p,ctrl-q"
)

var (
	// ErrEmptyVolDest empty volume destination error.
	ErrEmptyVolDest = fmt.Errorf("volume destination cannot be empty")
	// ErrTopPodNotRunning top error while pod not running.
	ErrTopPodNotRunning = fmt.Errorf("pods top can only be used on running pods")
	// ErrInvalidIPAddress invalid IP address error.
	ErrInvalidIPAddress = fmt.Errorf("invalid IP address")
	// ErrInvalidDNSAddress invalid DNS server address error.
	ErrInvalidDNSAddress = fmt.Errorf("invalid DNS address")
)

// DefineCreateDefaults sets default container create options.
func DefineCreateDefaults(opts *entities.ContainerCreateOptions) {
	opts.LogDriver = ""
	opts.CgroupParent = ""
	opts.MemorySwappiness = -1
	opts.Pull = ""
	// opts.ReadOnlyTmpFS = true
	opts.SdNotifyMode = define.SdNotifyModeContainer
	opts.Systemd = "true"
	opts.Ulimit = nil
	opts.SeccompPolicy = "default"
	opts.Volume = nil
}
