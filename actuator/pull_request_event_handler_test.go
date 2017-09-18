package actuator_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ninech/actuator/actuator"
	"github.com/ninech/actuator/testutils"
)

func TestHandleEventFails(t *testing.T) {
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

	t.Run("the action is not supported", func(t *testing.T) {
		event := testutils.NewTestEvent(1, "yolo", "ninech/actuator")
		handler := actuator.PullRequestEventHandler{Event: event}

		_, err := handler.HandleEvent()
		assert.Nil(t, err)
	})
}

func TestHandleEventActionOpened(t *testing.T) {
	event := testutils.NewTestEvent(1, actuator.ActionOpened, "ninech/actuator")
	config := testutils.NewDefaultConfig()
	githubClient := testutils.NewMockGithubClient()
	openshiftMock := testutils.OpenshiftMock{}
	actuator.ApplyOpenshiftTemplate = openshiftMock.ApplyOpenshiftTemplate

	handler := actuator.PullRequestEventHandler{
		Event:        event,
		Config:       config,
		GithubClient: githubClient}

	t.Run("applies the template in openshift", func(t *testing.T) {
		message, err := handler.HandleEvent()

		assert.Nil(t, err)
		assert.Equal(t, "Event for pull request #1 received. Thank you.", message)
		assert.Equal(t, config.GetRepositoryConfig(*event.Repo.FullName).Template, openshiftMock.AppliedTemplate, "it instantiates the template from the config")

		assert.Equal(t, openshiftMock.AppliedLabels["actuator.nine.ch/create-reason"], "GithubWebhook")
		assert.Equal(t, openshiftMock.AppliedLabels["actuator.nine.ch/branch"], event.PullRequest.Head.GetRef())
		assert.Equal(t, openshiftMock.AppliedLabels["actuator.nine.ch/pull-request"], strconv.Itoa(event.PullRequest.GetNumber()))

		assert.Equal(t, openshiftMock.AppliedParameters["BRANCH_NAME"], "pr-1")
	})

	t.Run("writes a comment on openshift", func(t *testing.T) {
		githubComment := githubClient.LastComment
		assert.NotNil(t, githubComment, "creates a comment on github")
		assert.Equal(t, "Your environment is being set-up on Openshift.", githubComment.GetBody())
		assert.Equal(t, "https://github.com/ninech/actuator/issues/1408#issuecomment-330230087", githubComment.GetHTMLURL())
	})
}
