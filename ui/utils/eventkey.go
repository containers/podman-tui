package utils

import "github.com/gdamore/tcell/v2"

// StringToEventKey returns list of key events equvalant to the input string.
func StringToEventKey(input string) []*tcell.EventKey {
	var events []*tcell.EventKey

	for i := 0; i < len(input); i++ {
		ch := rune(input[i])
		events = append(events, tcell.NewEventKey(256, ch, tcell.ModNone)) //nolint:gomnd
	}

	return events
}
