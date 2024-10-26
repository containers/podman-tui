package utils

import "github.com/gdamore/tcell/v2"

// StringToEventKey returns list of key events equvalant to the input string.
func StringToEventKey(input string) []*tcell.EventKey {
	events := []*tcell.EventKey{}

	for i := range input {
		ch := rune(input[i])
		events = append(events, tcell.NewEventKey(256, ch, tcell.ModNone)) //nolint:mnd
	}

	return events
}
