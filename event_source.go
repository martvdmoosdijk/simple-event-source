package simpleventsrc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"
)

// TODO Snapshots

var _ IEventSource = &EventSource{}

type IEventSource interface {
	RegisterEventType(eventType IEvent)
	ReplayEvents()
	AddEvent(event IEvent) error
}

type EventSource struct {
	// TODO - Add ID/Aggregate?
	EventRegistry  map[string]reflect.Type
	EventProvider  IEventProvider
	LatestPosition int64
	AddEventLock   sync.Mutex
}

type EventEntry struct {
	Position  int64       `json:"position"`
	Body      interface{} `json:"body"`
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
}

func NewEventSource(eventProvider IEventProvider) EventSource {
	return EventSource{
		EventRegistry:  make(map[string]reflect.Type),
		EventProvider:  eventProvider,
		LatestPosition: 0,
	}
}

func (es *EventSource) RegisterEventType(event IEvent) {
	es.EventRegistry[nameOf(event)] = reflect.TypeOf(event)
}

func (es *EventSource) ReplayEvents() {
	for eventEntry := range es.EventProvider.ReadEvents() {
		if eventEntry.Position < es.LatestPosition {
			panic(fmt.Errorf("ReplayEvents panic: Event position %v is lower than latest known position %v", eventEntry.Position, es.LatestPosition))
		}
		es.LatestPosition = eventEntry.Position

		event := reflect.New(es.EventRegistry[eventEntry.Type])
		data, err := json.Marshal(eventEntry.Body)
		if err != nil {
			panic(err)
		}

		if err := json.Unmarshal(data, event.Interface()); err != nil {
			panic(err)
		}

		(event.Interface().(IEvent)).Consume()
	}
}

func (es *EventSource) AddEvent(event IEvent) error {
	es.AddEventLock.Lock()
	defer es.AddEventLock.Unlock()

	// TODO Check if event type is in registry

	es.LatestPosition++
	eventEntry := EventEntry{
		Position:  es.LatestPosition,
		Body:      event,
		Type:      nameOf(event),
		Timestamp: time.Now().UTC(),
	}

	if err := es.EventProvider.SaveEvent(eventEntry); err != nil {
		return err
	}

	event.Consume()
	return nil
}
