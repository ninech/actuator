package actuator_test

import (
	"testing"

	"github.com/ninech/actuator/actuator"
	"github.com/ninech/actuator/github"

	"github.com/stretchr/testify/assert"
)

func TestPingEventHandler(t *testing.T) {
	t.Run("GetEventResponse", func(t *testing.T) {
		handler := actuator.NewPingEventHandler()
		response := handler.GetEventResponse(&github.Event{})

		assert.Equal(t, "Request received. Thank you.", response.Message)
	})

	t.Run("HandleEvent", func(t *testing.T) {
		handler := actuator.NewGenericEventHandler()
		handler.HandleEvent(&github.Event{})
	})
}
