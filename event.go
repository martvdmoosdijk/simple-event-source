package simpleventsrc

type IEvent interface {
	Validate() error
	Consume()
}

type EventMock struct {
	TotalInvocations    int
	ValidateInvocations int
	ValidateMock        func() error
	ConsumeInvocations  int
	ConsumeMock         func()
}

var _ IEvent = &EventMock{}

func (e *EventMock) Validate() error {
	e.ValidateInvocations++
	e.TotalInvocations++

	if e.ValidateMock == nil {
		return e.ValidateMock()
	}

	return nil
}

func (e *EventMock) Consume() {
	e.ConsumeInvocations++
	e.TotalInvocations++

	if e.ConsumeMock == nil {
		e.ConsumeMock()
	}
}
