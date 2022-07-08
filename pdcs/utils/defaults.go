package utils

import "fmt"

const (
	// DefaultContainerDetachKeys container's default attach keys string.
	DefaultContainerDetachKeys = "ctrl-p,ctrl-q"
)

// ErrEmptyVolDest empty volume destination error.
var ErrEmptyVolDest = fmt.Errorf("volume destination cannot be empty")
