package actuator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ninech/actuator/actuator"
	"github.com/ninech/actuator/testutils"
)

// GetMessage

func TestGetMessageWithUnsetMessageAndSetEvent(t *testing.T) {
	event := testutils.NewDefaultTestEvent()
	handler := actuator.PullRequestEventHandler{Event: event}

	assert.Equal(t, "Event for pull request #1 received. Thank you.", handler.GetMessage())
}

func TestGetMessageWithSetMessage(t *testing.T) {
	event := testutils.NewDefaultTestEvent()
	handler := actuator.PullRequestEventHandler{Event: event, Message: "hello"}

	assert.Equal(t, "hello", handler.GetMessage())
}

// HandleEvent

func TestHandleEventRepositoryIsNotConfigured(t *testing.T) {
	event := testutils.NewTestEvent(1, actuator.ActionOpened, "ninech/yoloproject")
	handler := actuator.PullRequestEventHandler{Event: event}

	err := handler.HandleEvent()
	assert.Nil(t, err)
	assert.Equal(t, "Repository ninech/yoloproject is not configured. Doing nothing.", handler.GetMessage())
}

func TestHandleEventActionOpened(t *testing.T) {
	event := testutils.NewTestEvent(1, actuator.ActionOpened, "ninech/actuator")
	handler := actuator.PullRequestEventHandler{Event: event, Config: testutils.NewDefaultConfig()}

	err := handler.HandleEvent()
	assert.Nil(t, err)
	// check apply template
}

// actionIsSupported

func TestActionIsSupportedWithSupportedAction(t *testing.T) {
	event := testutils.NewDefaultTestEvent()
	handler := actuator.PullRequestEventHandler{Event: event}

	err := handler.HandleEvent()
	assert.Nil(t, err)
}

func TestActionIsSupportedWithUnsupportedAction(t *testing.T) {
	event := testutils.NewTestEvent(1, "yolo", "ninech/actuator")
	handler := actuator.PullRequestEventHandler{Event: event}

	err := handler.HandleEvent()
	assert.Nil(t, err)
}
