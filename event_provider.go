package simpleventsrc

type IEventProvider interface {
	ReadEvents() <-chan EventEntry
	SaveEvent(EventEntry) error
}

type EventProviderMock struct {
	TotalInvocations      int
	ReadEventsInvocations int
	ReadEventsMock        func() <-chan EventEntry
	SaveEventInvocations  int
	SaveEventMock         func(EventEntry) error
}

var _ IEventProvider = &EventProviderMock{}

func (ep *EventProviderMock) ReadEvents() <-chan EventEntry {
	ep.ReadEventsInvocations++
	ep.TotalInvocations++

	if ep.ReadEventsMock == nil {
		return ep.ReadEventsMock()
	}

	c := make(chan EventEntry)
	close(c)
	return c
}

func (ep *EventProviderMock) SaveEvent(entry EventEntry) error {
	ep.SaveEventInvocations++
	ep.TotalInvocations++

	if ep.SaveEventMock == nil {
		return ep.SaveEventMock(entry)
	}

	return nil
}
