package simple_event_source

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

type IEventProvider interface {
	ReadEvents() <-chan EventEntry
	SaveEvent(EventEntry) error
}

type IEvent interface {
	Validate() error
	Consume()
}

type EventEntry struct {
	Position  int64       `json:"position"`
	Body      interface{} `json:"body"`
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
}

func New(eventProvider IEventProvider) EventSource {
	return EventSource{
		EventRegistry:  make(map[string]reflect.Type),
		EventProvider:  eventProvider,
		LatestPosition: 0,
	}
}

func (self *EventSource) RegisterEventType(event IEvent) {
	eventName := reflect.TypeOf(event).Name()
	self.EventRegistry[eventName] = reflect.TypeOf(event)
}

func (self *EventSource) ReplayEvents() {
	for eventEntry := range self.EventProvider.ReadEvents() {
		if eventEntry.Position < self.LatestPosition {
			panic(fmt.Errorf("ReplayEvents panic: Event position %v is lower than latest known position %v", eventEntry.Position, self.LatestPosition))
		}
		self.LatestPosition = eventEntry.Position

		event := reflect.New(self.EventRegistry[eventEntry.Type])
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

func (self *EventSource) AddEvent(event IEvent) error {
	self.AddEventLock.Lock()
	defer self.AddEventLock.Unlock()

	self.LatestPosition++
	eventEntry := EventEntry{
		Position:  self.LatestPosition,
		Body:      event,
		Type:      reflect.TypeOf(event).Name(),
		Timestamp: time.Now().UTC(),
	}

	if err := self.EventProvider.SaveEvent(eventEntry); err != nil {
		return err
	}

	event.Consume()
	return nil
}
