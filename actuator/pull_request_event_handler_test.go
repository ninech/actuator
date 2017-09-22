package actuator_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ninech/actuator/actuator"
	"github.com/ninech/actuator/openshift"
	"github.com/ninech/actuator/test"
)

func TestHandleEventFails(t *testing.T) {
	t.Run("the event's repository is not configured", func(t *testing.T) {
		event := test.NewTestEvent(1, actuator.ActionOpened, "ninech/yoloproject")
		handler := actuator.PullRequestEventHandler{Event: event}

		message, err := handler.HandleEvent()
		assert.Nil(t, err)
		assert.Equal(t, "Repository ninech/yoloproject is not configured. Doing nothing.", message)
	})

	t.Run("the event repository is disabled in the config file", func(t *testing.T) {
		config := test.NewDefaultConfig()
		config.Repositories[0].Enabled = false

		event := test.NewTestEvent(1, actuator.ActionOpened, config.Repositories[0].Fullname)
		handler := actuator.PullRequestEventHandler{Event: event, Config: config}

		message, _ := handler.HandleEvent()
		assert.Contains(t, message, "is disabled.")
	})

	t.Run("the action is not supported", func(t *testing.T) {
		event := test.NewTestEvent(1, "yolo", "ninech/actuator")
		handler := actuator.PullRequestEventHandler{Event: event}

		_, err := handler.HandleEvent()
		assert.Nil(t, err)
	})
}

func TestHandleEventActionOpened(t *testing.T) {
	test.DisableLogging()

	event := test.NewTestEvent(1, actuator.ActionOpened, "ninech/actuator")
	config := test.NewDefaultConfig()
	githubClient := test.NewMockGithubClient()
	openshiftClient := &test.OpenshiftMock{}

	handler := actuator.PullRequestEventHandler{
		Event:        event,
		Config:       config,
		GithubClient: githubClient,
		Openshift:    openshiftClient}

	t.Run("applies the template in openshift", func(t *testing.T) {
		message, err := handler.HandleEvent()
		assert.Nil(t, err)
		assert.Equal(t, "Event for pull request #1 received. Thank you.", message)
		assert.Equal(t, config.GetRepositoryConfig(event.RepositoryFullname).Template, openshiftClient.AppliedTemplate, "it instantiates the template from the config")

		assert.Equal(t, openshiftClient.AppliedLabels["actuator.nine.ch/create-reason"], "GithubWebhook")
		assert.Equal(t, openshiftClient.AppliedLabels["actuator.nine.ch/branch"], event.HeadRef)
		assert.Equal(t, openshiftClient.AppliedLabels["actuator.nine.ch/pull-request"], strconv.Itoa(event.IssueNumber))

		assert.Equal(t, openshiftClient.AppliedParameters["BRANCH_NAME"], "pr-1")
	})

	t.Run("writes a comment on github", func(t *testing.T) {
		handler.HandleEvent()
		githubComment := githubClient.LastComment
		assert.NotNil(t, githubComment, "creates a comment on github")
		assert.Equal(t, "Your environment is being set-up on Openshift. There is no route I can point you to.", githubComment.GetBody())
		assert.Equal(t, "https://github.com/ninech/actuator/issues/1#issuecomment-330230087", githubComment.GetHTMLURL())
	})

	t.Run("posts the url as comment", func(t *testing.T) {
		openshiftClient.NewAppOutputToReturn = openshift.NewAppOutput{Raw: `route "actuator" created`}
		handler.HandleEvent()

		githubComment := githubClient.LastComment
		assert.Equal(t, "Your environment is being set-up on Openshift. http://actuator.domain.com", githubComment.GetBody())
	})
}

func TestHandleEventActionClosed(t *testing.T) {
	test.DisableLogging()

	event := test.NewTestEvent(1, actuator.ActionClosed, "ninech/actuator")
	config := test.NewDefaultConfig()
	openshiftClient := &test.OpenshiftMock{}

	handler := actuator.PullRequestEventHandler{
		Event:     event,
		Config:    config,
		Openshift: openshiftClient}

	t.Run("deletes the objects in openshift", func(t *testing.T) {
		message, err := handler.HandleEvent()

		assert.Nil(t, err)
		assert.Equal(t, "Event for pull request #1 received. Thank you.", message)
		assert.Equal(t, openshiftClient.DeletedLabels["actuator.nine.ch/pull-request"], strconv.Itoa(event.IssueNumber))
	})
}
