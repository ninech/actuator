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

func TestGetMessageWithUnsetMessageAndSetEvent(t *testing.T) {
	event := getEvent(1, "opened")
	handler := actuator.PullRequestEventHandler{Event: event}

	assert.Equal(t, "Event for pull request #1 received. Thank you.", handler.GetMessage())
}

func TestGetMessageWithSetMessage(t *testing.T) {
	event := getEvent(1, "opened")
	handler := actuator.PullRequestEventHandler{Event: event, Message: "hello"}

	assert.Equal(t, "hello", handler.GetMessage())
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
