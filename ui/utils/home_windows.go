//go:build windows
// +build windows

package utils

import (
	"os"
)

// UserHomeDir returns user's home directory
func UserHomeDir() (string, error) {
	return os.UserHomeDir()
}
