package simpleventsrc

type IEvent interface {
	Validate() error
	Consume()
}

type EventMock struct {
	ValidateMock        func() error
	ValidateInvocations int
	ConsumeMock         func()
	ConsumeInvocations  int
	TotalInvocations    int
}

var _ IEvent = &EventMock{}

func (e *EventMock) Validate() error {
	e.ValidateInvocations++
	e.TotalInvocations++

	if e.ValidateMock == nil {
		panic("EventMock.ValidateMock() not implemented")
	}
	return e.ValidateMock()
}

func (e *EventMock) Consume() {
	e.ConsumeInvocations++
	e.TotalInvocations++

	// TODO
	if e.ConsumeMock == nil {
		panic("EventMock.ConsumeMock() not implemented")
	}
	e.ConsumeMock()
}
