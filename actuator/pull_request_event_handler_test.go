package actuator_test

import (
	"testing"

	"github.com/google/go-github/github"
	"github.com/ninech/actuator/actuator"
	"github.com/stretchr/testify/assert"
)

func getEvent(number int, action string) *github.PullRequestEvent {
	return &github.PullRequestEvent{Action: &action, Number: &number}
}

func TestActionIsSupportedWithSupportedAction(t *testing.T) {
	event := getEvent(1, "opened")
	handler := actuator.PullRequestEventHandler{Event: event}

	err := handler.HandleEvent()
	assert.Nil(t, err)
}

func TestActionIsSupportedWithUnsupportedAction(t *testing.T) {
	event := getEvent(1, "yolo")
	handler := actuator.PullRequestEventHandler{Event: event}

	err := handler.HandleEvent()
	assert.Nil(t, err)
}
