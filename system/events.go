package system

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/containers/podman-tui/pdcs/sysinfo"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/docker/docker/api/types/events"
	"github.com/rs/zerolog/log"
)

type podmanEvents struct {
	mu                sync.Mutex
	status            bool
	eventChan         chan entities.Event
	eventCancelChan   chan bool
	cancelChan        chan bool
	eventBuffer       []string
	messageBuffer     []string
	messageBufferSize int
	hasNewEvent       bool
}

func (engine *Engine) startEventStreamer() {
	engine.sysEvents.status = true
	engine.sysEvents.eventCancelChan = make(chan bool)
	engine.sysEvents.cancelChan = make(chan bool)
	engine.sysEvents.eventChan = make(chan entities.Event, 20)
	engine.sysEvents.eventBuffer = []string{}
	engine.sysEvents.messageBuffer = []string{}
	go engine.eventReader()
	go engine.streamEvents()
}

func (engine *Engine) stopEventStreamer() {
	engine.sysEvents.mu.Lock()
	engine.sysEvents.status = false
	engine.sysEvents.eventCancelChan <- true
	engine.sysEvents.mu.Unlock()
}

func (engine *Engine) streamEvents() {
	log.Debug().Msg("health check: event steamer started")
	if err := sysinfo.Events(engine.sysEvents.eventChan, engine.sysEvents.eventCancelChan); err != nil {
		// TODO error check for events
		log.Error().Msgf("health check: event streamer %v", err)
	}
	log.Debug().Msg("health check: event steamer stopped")
}

func (engine *Engine) eventReader() {
	log.Debug().Msg("health check:: event reader started")
	for {
		select {
		case <-engine.sysEvents.cancelChan:
			log.Debug().Msg("health check: event reader stopped")
			engine.stopEventStreamer()
			return
		case event := <-engine.sysEvents.eventChan:
			{
				msg := engine.convertEventToHumanReadable(event.Message)
				engine.addEvent(event.Message)
				engine.addEventMessage(msg)
				if strings.TrimSpace(msg) != "" {
					log.Debug().Msgf("health check: event reader received %s", msg)
				}

			}
		}
	}
}

// GetEventMessages returns events buffer messages
func (engine *Engine) GetEventMessages() []string {
	var events []string
	engine.sysEvents.mu.Lock()
	events = engine.sysEvents.messageBuffer
	engine.sysEvents.hasNewEvent = false
	engine.sysEvents.mu.Unlock()
	return events
}

// HasNewEvent returns true if there is new event added to event buffer
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

// GetEvents returns event buffer types
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

func (engine *Engine) addEvent(event events.Message) {
	engine.sysEvents.mu.Lock()
	engine.sysEvents.hasNewEvent = true
	engine.sysEvents.eventBuffer = append(engine.sysEvents.eventBuffer, event.Type)
	engine.sysEvents.mu.Unlock()
}

// convertEventToHumanReadable returns human readable event as a formatted string
func (engine *Engine) convertEventToHumanReadable(event events.Message) string {
	var humanFormat string
	id := event.Actor.ID
	//id = stringid.TruncateID(id)
	evtime := time.Unix(event.Time, event.TimeNano).String()
	switch event.Type {
	case "container", "pod":
		humanFormat = fmt.Sprintf("%s %s %s %s (image=%s, name=%s", evtime, event.Type, event.Action, id, event.Actor.Attributes["image"], event.Actor.Attributes["name"])
		// check if the container has labels and add it to the output
		if len(event.Actor.Attributes) > 0 {
			for k, v := range event.Actor.Attributes {
				humanFormat += fmt.Sprintf(", %s=%s", k, v)
			}
		}
		humanFormat += ")"

	case "network":
		humanFormat = fmt.Sprintf("%s %s %s %s (container=%s, name=%s)", evtime, event.Type, event.Action, id, event.Actor.Attributes["id"], event.Actor.Attributes["network"])
	case "image":
		humanFormat = fmt.Sprintf("%s %s %s %s %s", evtime, event.Type, event.Action, id, event.Actor.Attributes["name"])
	case "system":
		humanFormat = fmt.Sprintf("%s %s %s", evtime, event.Type, event.Action)
	case "volume":
		humanFormat = fmt.Sprintf("%s %s %s %s", evtime, event.Type, event.Action, event.Actor.Attributes["name"])

	}

	return humanFormat
}
