package simpleventsrc

type IEventProvider interface {
	ReadEvents() <-chan EventEntry
	SaveEvent(EventEntry) error
}

type EventProviderMock struct {
	ReadEventsMock        func() <-chan EventEntry
	ReadEventsInvocations int
	SaveEventMock         func(EventEntry) error
	SaveEventInvocations  int
	TotalInvocations      int
}

var _ IEventProvider = &EventProviderMock{}

func (ep *EventProviderMock) ReadEvents() <-chan EventEntry {
	ep.ReadEventsInvocations++
	ep.TotalInvocations++

	if ep.ReadEventsMock == nil {
		panic("EventProviderMock.ReadEventsMock() not implemented")
	}
	return ep.ReadEventsMock()
}

func (ep *EventProviderMock) SaveEvent(entry EventEntry) error {
	ep.SaveEventInvocations++
	ep.TotalInvocations++

	if ep.SaveEventMock == nil {
		panic("EventProviderMock.SaveEventMock() not implemented")
	}
	return ep.SaveEventMock(entry)
}
