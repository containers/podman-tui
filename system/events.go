package system

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/containers/podman-tui/pdcs/registry"
	"github.com/containers/podman-tui/pdcs/sysinfo"
	"github.com/containers/podman/v5/pkg/domain/entities/types"
	"github.com/rs/zerolog/log"
)

var eventChannelSize = 20

type podmanEvents struct {
	mu                sync.Mutex
	status            bool
	eventChan         chan types.Event
	eventCancelChan   chan bool
	cancelChan        chan bool
	eventBuffer       []string
	messageBuffer     []string
	messageBufferSize int
	hasNewEvent       bool
}

func (engine *Engine) startEventStreamer() {
	log.Debug().Msgf("health check: start event streamer")

	engine.sysEvents.mu.Lock()
	engine.sysEvents.status = true
	engine.sysEvents.eventBuffer = []string{}
	engine.sysEvents.messageBuffer = []string{}
	engine.sysEvents.cancelChan = make(chan bool)
	engine.sysEvents.eventCancelChan = make(chan bool)
	engine.sysEvents.eventChan = make(chan types.Event, eventChannelSize)
	engine.sysEvents.mu.Unlock()

	go engine.eventReader()
	go engine.streamEvents()
}

func (engine *Engine) streamEvents() {
	log.Debug().Msg("health check: pdcs event steamer started")

	for {
		if err := sysinfo.Events(engine.sysEvents.eventChan, engine.sysEvents.eventCancelChan); err != nil {
			log.Error().Msgf("health check: pdcs event streamer %v", err)
			engine.sysEvents.cancelChan <- true
			engine.sysEvents.mu.Lock()
			engine.sysEvents.status = false
			engine.sysEvents.mu.Unlock()
			log.Debug().Msgf("health check: pdcs event steamer cancel sent")

			break
		}
	}

	log.Debug().Msg("health check: pdcs event streamer stopped")

	close(engine.sysEvents.eventCancelChan)
}

func (engine *Engine) eventReader() {
	log.Debug().Msg("health check: event reader started")

	for {
		select {
		case <-engine.sysEvents.cancelChan:
			log.Debug().Msg("health check: event reader stopped")

			close(engine.sysEvents.cancelChan)
			registry.CancelContext()

			return

		case event := <-engine.sysEvents.eventChan:
			{
				msg := engine.convertEventToHumanReadable(event)
				if strings.TrimSpace(msg) != "" {
					log.Debug().Msgf("health check: event reader received %s", msg)
				}

				engine.addEvent(event)
				engine.addEventMessage(msg)
			}
		}
	}
}

// GetEventMessages returns events buffer messages.
func (engine *Engine) GetEventMessages() []string {
	var events []string

	engine.sysEvents.mu.Lock()
	events = engine.sysEvents.messageBuffer
	engine.sysEvents.hasNewEvent = false
	engine.sysEvents.mu.Unlock()

	return events
}

// HasNewEvent returns true if there is new event added to event buffer.
func (engine *Engine) HasNewEvent() bool {
	hasEvent := false

	engine.sysEvents.mu.Lock()
	hasEvent = engine.sysEvents.hasNewEvent
	engine.sysEvents.mu.Unlock()

	return hasEvent
}

func (engine *Engine) addEventMessage(msg string) {
	engine.sysEvents.mu.Lock()
	if len(engine.sysEvents.messageBuffer) == engine.sysEvents.messageBufferSize {
		// empty first 10 entries
		engine.sysEvents.messageBuffer = engine.sysEvents.messageBuffer[20:]
	}

	engine.sysEvents.messageBuffer = append(engine.sysEvents.messageBuffer, msg)
	engine.sysEvents.hasNewEvent = true
	engine.sysEvents.mu.Unlock()
}

// GetEvents returns event buffer types.
func (engine *Engine) GetEvents() []string {
	var events []string

	engine.sysEvents.mu.Lock()
	events = engine.sysEvents.eventBuffer
	// empty buffer.
	engine.sysEvents.eventBuffer = []string{}
	engine.sysEvents.mu.Unlock()

	events = unique(events)

	return events
}

// EventStatus returns event stats.
func (engine *Engine) EventStatus() bool {
	engine.sysEvents.mu.Lock()
	defer engine.sysEvents.mu.Unlock()

	return engine.sysEvents.status
}

func (engine *Engine) addEvent(event types.Event) {
	engine.sysEvents.mu.Lock()
	engine.sysEvents.hasNewEvent = true
	engine.sysEvents.eventBuffer = append(engine.sysEvents.eventBuffer, string(event.Type))
	engine.sysEvents.mu.Unlock()
}

// convertEventToHumanReadable returns human readable event as a formatted string.
func (engine *Engine) convertEventToHumanReadable(event types.Event) string {
	var humanFormat string

	id := event.Actor.ID
	evtime := time.Unix(event.Time, event.TimeNano).String()

	switch string(event.Type) {
	case "container", "pod":
		humanFormat = fmt.Sprintf("%s %s %s %s (image=%s, name=%s",
			evtime,
			event.Type,
			event.Action,
			id,
			event.Actor.Attributes["image"],
			event.Actor.Attributes["name"],
		)

		// check if the container has labels and add it to the output
		if len(event.Actor.Attributes) > 0 {
			for k, v := range event.Actor.Attributes {
				humanFormat += fmt.Sprintf(", %s=%s", k, v)
			}
		}

		humanFormat += ")"

	case "network":
		humanFormat = fmt.Sprintf("%s %s %s %s (container=%s, name=%s)",
			evtime,
			event.Type,
			event.Action,
			id,
			event.Actor.Attributes["id"],
			event.Actor.Attributes["network"],
		)
	case "image":
		humanFormat = fmt.Sprintf("%s %s %s %s %s",
			evtime,
			event.Type,
			event.Action,
			id,
			event.Actor.Attributes["name"],
		)
	case "system":
		humanFormat = fmt.Sprintf("%s %s %s",
			evtime,
			event.Type,
			event.Action,
		)
	case "volume":
		humanFormat = fmt.Sprintf("%s %s %s %s",
			evtime,
			event.Type,
			event.Action,
			event.Actor.Attributes["name"],
		)
	}

	return humanFormat
}
