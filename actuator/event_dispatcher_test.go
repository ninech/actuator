package actuator_test

import (
	"testing"

	"github.com/ninech/actuator/actuator"
	"github.com/ninech/actuator/github"
	"github.com/ninech/actuator/test"
	"github.com/stretchr/testify/assert"
)

func TestEventDispatcher(t *testing.T) {
	event := test.NewDefaultTestEvent()

	t.Run("GetEventResponse", func(t *testing.T) {
		event.Type = 999

		dispatcher := actuator.EventDispatcher{Event: *event}

		response := dispatcher.GetEventResponse()

		assert.Equal(t, "Request received. Doing nothing.", response.Message)
		assert.IsType(t, &actuator.GenericEventHandler{}, dispatcher.LastEventHandler)
	})

	t.Run("HandleEvent", func(t *testing.T) {
		handler := test.NewMockEventHandler("yay!")
		dispatcher := actuator.EventDispatcher{
			Event:            *event,
			LastEventHandler: handler}

		dispatcher.HandleEvent()

		assert.True(t, handler.EventWasHandled)
	})

	t.Run("FindEventHandler", func(t *testing.T) {
		t.Run("PullRequestEvent", func(t *testing.T) {
			event.Type = github.PullRequestEvent
			dispatcher := actuator.EventDispatcher{Event: *event}

			handler := dispatcher.FindEventHandler()

			assert.IsType(t, &actuator.PullRequestEventHandler{}, handler)
		})

		t.Run("any other event type", func(t *testing.T) {
			event.Type = 999
			dispatcher := actuator.EventDispatcher{Event: *event}

			handler := dispatcher.FindEventHandler()

			assert.IsType(t, &actuator.GenericEventHandler{}, handler)
		})
	})
}
