package actuator_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ninech/actuator/actuator"
	"github.com/ninech/actuator/openshift"
	"github.com/ninech/actuator/testutils"
)

func TestHandleEvent(t *testing.T) {
	t.Run("the event's repository is not configured", func(t *testing.T) {
		event := testutils.NewTestEvent(1, actuator.ActionOpened, "ninech/yoloproject")
		handler := actuator.PullRequestEventHandler{Event: event}

		message, err := handler.HandleEvent()
		assert.Nil(t, err)
		assert.Equal(t, "Repository ninech/yoloproject is not configured. Doing nothing.", message)
	})

	t.Run("the event repository is disabled in the config file", func(t *testing.T) {
		config := testutils.NewDefaultConfig()
		config.Repositories[0].Enabled = false

		event := testutils.NewTestEvent(1, actuator.ActionOpened, config.Repositories[0].Fullname)
		handler := actuator.PullRequestEventHandler{Event: event, Config: config}

		message, _ := handler.HandleEvent()
		assert.Contains(t, message, "is disabled.")
	})

	t.Run("event action opened", func(t *testing.T) {
		event := testutils.NewTestEvent(1, actuator.ActionOpened, "ninech/actuator")
		config := testutils.NewDefaultConfig()
		handler := actuator.PullRequestEventHandler{Event: event, Config: config}

		var (
			appliedTemplate   string
			appliedLabels     openshift.ObjectLabels
			appliedParameters openshift.TemplateParameters
		)
		actuator.ApplyOpenshiftTemplate = func(templateName string, templateParameters openshift.TemplateParameters, labels openshift.ObjectLabels) (string, error) {
			appliedTemplate = templateName
			appliedLabels = labels
			appliedParameters = templateParameters
			return "", nil
		}

		message, err := handler.HandleEvent()

		assert.Nil(t, err)
		assert.Equal(t, "Event for pull request #1 received. Thank you.", message)
		assert.Equal(t, config.GetRepositoryConfig(*event.Repo.FullName).Template, appliedTemplate, "it instantiates the template from the config")

		assert.Equal(t, appliedLabels["actuator.nine.ch/create-reason"], "GithubWebhook")
		assert.Equal(t, appliedLabels["actuator.nine.ch/branch"], event.PullRequest.Head.GetRef())
		assert.Equal(t, appliedLabels["actuator.nine.ch/pull-request"], strconv.Itoa(event.PullRequest.GetNumber()))

		assert.Equal(t, appliedParameters["BRANCH_NAME"], "pr-1")
	})
}

func TestActionIsSupported(t *testing.T) {
	t.Run("with supported action", func(t *testing.T) {
		event := testutils.NewDefaultTestEvent()
		handler := actuator.PullRequestEventHandler{Event: event}

		_, err := handler.HandleEvent()
		assert.Nil(t, err)
	})

	t.Run("with unsupported action", func(t *testing.T) {
		event := testutils.NewTestEvent(1, "yolo", "ninech/actuator")
		handler := actuator.PullRequestEventHandler{Event: event}

		_, err := handler.HandleEvent()
		assert.Nil(t, err)
	})
}
