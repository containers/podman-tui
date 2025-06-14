package utils

import "github.com/containers/podman-tui/pdcs/registry"

type ConnSort []registry.Connection

func (a ConnSort) Len() int      { return len(a) }
func (a ConnSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type ConnectionListSortedName struct{ ConnSort }

func (a ConnectionListSortedName) Less(i, j int) bool {
	return a.ConnSort[i].Name < a.ConnSort[j].Name
}
