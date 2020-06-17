package simpleventsrc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

type EventSource struct {
	topic         string
	eventRegistry map[string]reflect.Type
	eventProvider IEventProvider
	newPosition   int64
	eventQueue    chan eventQueueItem
}

type eventQueueItem struct {
	event    IEvent
	callback func(error)
}

type EventEntry struct {
	ID        string      `json:"id"`
	Position  int64       `json:"position"`
	Body      interface{} `json:"body"`
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
}

func NewEventSource(topic string, eventProvider IEventProvider) *EventSource {
	es := &EventSource{
		topic:         topic,
		eventRegistry: make(map[string]reflect.Type),
		eventProvider: eventProvider,
		newPosition:   -1,
		eventQueue:    make(chan eventQueueItem, 50),
	}

	es.handleQueue()
	return es
}

func (es *EventSource) handleQueue() {
	go func() {
		for eq := range es.eventQueue {
			eventType := typeOf(eq.event).Name()
			if _, ok := es.eventRegistry[eventType]; !ok {
				eq.callback(fmt.Errorf("Event type %v is unregistered", eventType))
				return
			}

			if es.newPosition == -1 {
				eq.callback(fmt.Errorf("EventSource needs to be replayed before new events can be added"))
				return
			}

			eventEntry := EventEntry{
				Position:  es.newPosition,
				Body:      eq.event,
				Type:      eventType,
				Timestamp: time.Now().UTC(),
			}
			if err := es.eventProvider.SaveEvent(eventEntry); err != nil {
				eq.callback(err)
			}

			es.newPosition++
			eq.event.Consume()
			eq.callback(nil)
		}
	}()
}

func (es *EventSource) ReplayEvents() {
	es.newPosition = -1
	newPosition := int64(0)

	for eventEntry := range es.eventProvider.ReadEvents() {
		if eventEntry.Position < newPosition {
			panic(fmt.Errorf("Error replaying events, invalid order of event positions"))
		}
		newPosition = eventEntry.Position

		data, err := json.Marshal(eventEntry.Body)
		if err != nil {
			panic(err)
		}

		if _, ok := es.eventRegistry[eventEntry.Type]; !ok {
			panic("Error replaying events, event type has not been registered yet")
		}

		event := reflect.New(es.eventRegistry[eventEntry.Type]).Interface()
		if err := json.Unmarshal(data, event); err != nil {
			panic(err)
		}

		(event.(IEvent)).Consume()
	}

	es.newPosition = newPosition
}

func (es *EventSource) RegisterEventType(et IEvent) {
	eventType := typeOf(et)
	es.eventRegistry[eventType.Name()] = eventType
}

func (es *EventSource) AddEvent(e IEvent, cb func(error)) {
	es.eventQueue <- eventQueueItem{event: e, callback: cb}
}
