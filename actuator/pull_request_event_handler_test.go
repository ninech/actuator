package actuator_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ninech/actuator/actuator"
	"github.com/ninech/actuator/openshift"
	"github.com/ninech/actuator/testutils"
)

func TestGetMessage(t *testing.T) {
	t.Run("with message not set but event is set", func(t *testing.T) {
		event := testutils.NewDefaultTestEvent()
		handler := actuator.PullRequestEventHandler{Event: event}

		assert.Equal(t, "Event for pull request #1 received. Thank you.", handler.GetMessage(),
			"it sets the default message")
	})

	t.Run("with set message", func(t *testing.T) {
		event := testutils.NewDefaultTestEvent()
		handler := actuator.PullRequestEventHandler{Event: event, Message: "hello"}

		assert.Equal(t, "hello", handler.GetMessage(),
			"it returns the given message")
	})
}

func TestHandleEvent(t *testing.T) {
	t.Run("the event's repository is not configured", func(t *testing.T) {
		event := testutils.NewTestEvent(1, actuator.ActionOpened, "ninech/yoloproject")
		handler := actuator.PullRequestEventHandler{Event: event}

		err := handler.HandleEvent()
		assert.Nil(t, err)
		assert.Equal(t, "Repository ninech/yoloproject is not configured. Doing nothing.", handler.GetMessage())
	})

	t.Run("the event repository is disabled in the config file", func(t *testing.T) {
		config := testutils.NewDefaultConfig()
		config.Repositories[0].Enabled = false

		event := testutils.NewTestEvent(1, actuator.ActionOpened, config.Repositories[0].Fullname)
		handler := actuator.PullRequestEventHandler{Event: event, Config: config}

		handler.HandleEvent()
		assert.Contains(t, handler.GetMessage(), "is disabled.")
	})

	t.Run("event action opened", func(t *testing.T) {
		event := testutils.NewTestEvent(1, actuator.ActionOpened, "ninech/actuator")
		config := testutils.NewDefaultConfig()
		handler := actuator.PullRequestEventHandler{Event: event, Config: config}

		var templateInstantiated string
		var labelsApplied openshift.ObjectLabels
		actuator.NewAppFromTemplate = func(templateName string, templateParameters openshift.TemplateParameters, labels openshift.ObjectLabels) (string, error) {
			templateInstantiated = templateName
			labelsApplied = labels
			return "", nil
		}

		err := handler.HandleEvent()
		assert.Nil(t, err)
		assert.Equal(t, config.GetRepositoryConfig(*event.Repo.FullName).Template, templateInstantiated, "it instantiates the template from the config")

		assert.Equal(t, labelsApplied["actuator.nine.ch/create-reason"], "GithubWebhook")
		assert.Equal(t, labelsApplied["actuator.nine.ch/branch"], event.PullRequest.Head.GetRef())
		assert.Equal(t, labelsApplied["actuator.nine.ch/pull-request"], strconv.Itoa(event.PullRequest.GetNumber()))
	})
}

func TestActionIsSupported(t *testing.T) {
	t.Run("with supported action", func(t *testing.T) {
		event := testutils.NewDefaultTestEvent()
		handler := actuator.PullRequestEventHandler{Event: event}

		err := handler.HandleEvent()
		assert.Nil(t, err)
	})

	t.Run("with unsupported action", func(t *testing.T) {
		event := testutils.NewTestEvent(1, "yolo", "ninech/actuator")
		handler := actuator.PullRequestEventHandler{Event: event}

		err := handler.HandleEvent()
		assert.Nil(t, err)
	})
}
