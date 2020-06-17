package simpleventsrc

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisterEventType(t *testing.T) {
	es := NewEventSource("test", &EventProviderMock{})
	es.RegisterEventType(&EventMock{})

	require.Len(t, es.eventRegistry, 1)
	require.Contains(t, es.eventRegistry, "EventMock")
	require.NotNil(t, es.eventRegistry["EventMock"])
}

func TestAddEvent_FailWithoutRegister(t *testing.T) {
	es := NewEventSource("test", &EventProviderMock{})

	var err error
	var wg sync.WaitGroup

	wg.Add(1)
	es.AddEvent(&EventMock{}, func(addErr error) {
		err = addErr
		wg.Done()
	})
	wg.Wait()

	require.Error(t, err)
}

func TestAddEvent(t *testing.T) {
	require := require.New(t)
	ep := &EventProviderMock{
		ReadEventsMock: func() <-chan EventEntry {
			c := make(chan EventEntry, 1)
			c <- EventEntry{Type: "EventMock"}
			close(c)
			return c
		},
		SaveEventMock: func(EventEntry) error {
			return nil
		},
	}
	e := &EventMock{
		ConsumeMock: func() {},
	}

	es := NewEventSource("some-topic", ep)
	require.Equal(es.newPosition, int64(-1))
	es.RegisterEventType(&EventMock{})
	es.ReplayEvents()
	require.Equal(es.newPosition, int64(0))

	var err error
	var wg sync.WaitGroup

	wg.Add(1)
	err = fmt.Errorf("AddEvent callback not invoked")
	es.AddEvent(e, func(addErr error) {
		err = addErr
		wg.Done()
	})
	wg.Wait()

	require.NoError(err)
	require.Equal(ep.ReadEventsInvocations, 1)
	require.Equal(ep.SaveEventInvocations, 1)
	require.Equal(ep.TotalInvocations, 2)
	require.Equal(e.ConsumeInvocations, 1)
	require.Equal(e.TotalInvocations, 1)

	require.Equal(es.newPosition, int64(1))
}
