package actuator_test

import (
	"testing"

	"github.com/ninech/actuator/actuator"

	"github.com/stretchr/testify/assert"
)

func TestGenericEventHandler(t *testing.T) {
	t.Run("GetEventResponse", func(t *testing.T) {
		handler := actuator.NewGenericEventHandler()
		response := handler.GetEventResponse()

		assert.Equal(t, "Request received. Doing nothing.", response.Message)
	})

	t.Run("HandleEvent", func(t *testing.T) {
		handler := actuator.NewGenericEventHandler()
		handler.HandleEvent()
	})
}
