package simpleventsrc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisterEventType(t *testing.T) {
	es := NewEventSource(&EventProviderMock{})
	require.Len(t, es.EventRegistry, 0)

	es.RegisterEventType(&EventMock{})
	require.Len(t, es.EventRegistry, 1)
	require.NotNil(t, es.EventRegistry["EventMock"])
}

func TestReplayEvents(t *testing.T) {
	// TODO
}

// func TestAddEvent(t *testing.T) {
// 	ep := EventProviderMock{
// 		ReadEventsMock: func() <-chan EventEntry {
// 			c := make(chan EventEntry)

// 			go func() {
// 				c <- EventEntry{Type: "EventMock"}
// 				c <- EventEntry{Type: "EventMock"}
// 				close(c)
// 			}()

// 			return c
// 		},
// 	}
// 	es := NewEventSource(&ep)
// 	es.RegisterEventType(&EventMock{})
// 	// es.ReplayEvents()
// }
