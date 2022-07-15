package utils

import "fmt"

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
