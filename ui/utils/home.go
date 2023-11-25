//go:build !windows
// +build !windows

package utils

import "github.com/containers/storage/pkg/unshare"

// UserHomeDir returns user's home directory.
func UserHomeDir() (string, error) {
	// only get HomeDir when necessary.
	home, err := unshare.HomeDir()
	if err != nil {
		return "", err
	}

	return home, nil
}
