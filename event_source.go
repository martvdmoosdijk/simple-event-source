package simple_event_source

import (
	"encoding/json"
	"reflect"
	"sync"
	"time"
)

var _ IEventSource = &EventSource{}

type IEventSource interface {
	RegisterEventType(eventType IEvent)
	ReplayEvents()
	AddEvent(event IEvent) error
}

type EventSource struct {
	// TODO - Add ID?
	EventRegistry map[string]reflect.Type
	EventProvider IEventProvider
	Position      int64 // TODO - Time string instead of counter?
	AddEventLock  sync.Mutex
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
		EventRegistry: make(map[string]reflect.Type),
		EventProvider: eventProvider,
		Position:      0,
	}
}

func (self *EventSource) RegisterEventType(event IEvent) {
	eventName := reflect.TypeOf(event).Name()
	self.EventRegistry[eventName] = reflect.TypeOf(event)
}

func (self *EventSource) ReplayEvents() {
	for event := range self.EventProvider.ReadEvents() {
		if event.Position > self.Position {
			self.Position = event.Position
		}

		typedEvent := reflect.New(self.EventRegistry[event.Type])
		bodyData, err := json.Marshal(event.Body)
		if err != nil {
			panic(err)
		}

		// Write body to typedEvent
		if err := json.Unmarshal(bodyData, typedEvent.Interface()); err != nil {
			panic(err)
		}

		// Parse and consume event
		(typedEvent.Interface().(IEvent)).Consume()
	}
}

func (self *EventSource) AddEvent(event IEvent) error {
	self.AddEventLock.Lock()
	defer self.AddEventLock.Unlock()

	self.Position++
	eventEntry := EventEntry{
		Position:  self.Position,
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
