package utils

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
)

const (
	// IDLength max ID length to display.
	IDLength = 12
	// RefreshInterval application refresh interval.
	RefreshInterval = 1000 * time.Millisecond

	ContainerIDLabel  = "CONTAINER ID:"
	PruneCommandLabel = "prune"

	idLimit = 12
)

var (
	ErrURLMissingScheme = errors.New("url missing scheme")
	ErrInvalidFilename  = errors.New("invalid filename (should not contain ':')")
)

// GetIDWithLimit return ID string with limited string characters.
func GetIDWithLimit(id string) string {
	if len(id) > 0 {
		if len(id) >= idLimit {
			return id[:idLimit]
		}
	}

	return id
}

// LabelWidthLeftPadding adds left space padding.
func LabelWidthLeftPadding(input string, padding int) string {
	label := input
	for range padding {
		label = " " + label
	}

	return label
}

// StringToInputLabel create string with max width required for input fields.
func StringToInputLabel(input string, maxWidth int) string {
	label := ""
	labelIndex := 0

	for index, char := range input {
		if index >= maxWidth {
			break
		}

		label += string(char)
		labelIndex++
	}

	for index := labelIndex; index < maxWidth; index++ {
		label += " "
	}

	return label
}

// AlignStringListWidth returns max string len in the list.
func AlignStringListWidth(list []string) ([]string, int) {
	var (
		maxAlignment = 0
		alignedList  = make([]string, 0)
	)

	for _, item := range list {
		if len(item) > maxAlignment {
			maxAlignment = len(item)
		}
	}

	for _, item := range list {
		if len(item) < maxAlignment {
			need := maxAlignment - len(item)
			for range need {
				item += " "
			}
		}

		alignedList = append(alignedList, item)
	}

	return alignedList, maxAlignment
}

// EmptyBoxSpace returns simple Box without border with bgColor as background.
func EmptyBoxSpace(bgColor tcell.Color) *tview.Box {
	box := tview.NewBox()
	box.SetBackgroundColor(bgColor)
	box.SetBorder(false)

	return box
}

// ResolveHomeDir converts a path referencing the home directory via "~"
// to an absolute path.
func ResolveHomeDir(path string) (string, error) {
	// check if the path references the home dir to avoid work
	// don't use strings.HasPrefix(path, "~") as this doesn't match "~" alone
	// use strings.HasPrefix(...) to not match "something/~/something"
	if path != "~" && !strings.HasPrefix(path, "~/") {
		// path does not reference home dir -> Nothing to do
		return path, nil
	}

	// only get HomeDir when necessary
	home, err := UserHomeDir()
	if err != nil {
		return "", err
	}

	// replace the first "~" (start of path) with the HomeDir to resolve "~"
	return strings.Replace(path, "~", home, 1), nil
}

// Following codes are from https://github.com/containers/podman/blob/main/cmd/podman/parse/net.go

// ValidateFileName returns an error if filename contains ":"
// as it is currently not supported.
func ValidateFileName(filename string) error {
	if strings.Contains(filename, ":") {
		return fmt.Errorf("%w %q", ErrInvalidFilename, filename)
	}

	return nil
}

// ValidURL checks a string urlStr is a url or not.
func ValidURL(urlStr string) error {
	url, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return errors.Wrapf(err, "invalid url %q", urlStr)
	}

	if url.Scheme == "" {
		return fmt.Errorf("%w %q", ErrURLMissingScheme, urlStr)
	}

	return nil
}
