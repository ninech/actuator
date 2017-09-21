package actuator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ninech/actuator/actuator"
)

func TestHandleEvent(t *testing.T) {
	handler := actuator.GenericEventHandler{}

	message, err := handler.HandleEvent()

	assert.Equal(t, "Not processing this type of event", message)
	assert.Nil(t, err)
}
