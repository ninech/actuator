package actuator

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/ninech/actuator/openshift"
)

// SupportedPullRequestActions defines all pull request event actions which are supported by this app.
const (
	ActionOpened   = "opened"
	ActionClosed   = "closed"
	ActionReopened = "reopened"
)

// SupportedPullRequestActions defines all actions which are currently supported to be handled
var SupportedPullRequestActions = [1]string{ActionOpened}

// PullRequestEventHandler handles pull request events
type PullRequestEventHandler struct {
	Event   *github.PullRequestEvent
	Message string
	Config  Configuration
}

// ApplyOpenshiftTemplate creates new objects in openshift using a template
var ApplyOpenshiftTemplate = openshift.NewAppFromTemplate

// HandleEvent handles a pull request event from github
func (h *PullRequestEventHandler) HandleEvent() (string, error) {
	if !h.actionIsSupported() {
		return "Event is not relevant and will be ignored.", nil
	}

	repositoryName := h.Event.Repo.GetFullName()
	repositoryConfig := h.Config.GetRepositoryConfig(repositoryName)
	if repositoryConfig == nil {
		return fmt.Sprintf("Repository %s is not configured. Doing nothing.", repositoryName), nil
	}

	if !repositoryConfig.Enabled {
		return fmt.Sprintf("Repository %s is disabled. Doing nothing.", repositoryName), nil
	}

	switch h.Event.GetAction() {
	case ActionOpened:
		labels := h.buildLabelsFromEvent(h.Event)
		// TODO: pass template params from config
		output, err := ApplyOpenshiftTemplate(repositoryConfig.Template, openshift.TemplateParameters{}, labels)
		if err != nil {
			return err.Error(), err
		}

		Logger.Println(output)
		return fmt.Sprintf("Event for pull request #%d received. Thank you.", h.Event.GetNumber()), nil
	default:
		return "No handler for this action defined.", errors.New("no action handled")
	}
}

// actionIsSupported returns true when the provided action is currently supported by the app
func (h *PullRequestEventHandler) actionIsSupported() bool {
	for _, a := range SupportedPullRequestActions {
		if a == h.Event.GetAction() {
			return true
		}
	}
	return false
}

func (h *PullRequestEventHandler) buildLabelsFromEvent(event *github.PullRequestEvent) openshift.ObjectLabels {
	return openshift.ObjectLabels{
		"actuator.nine.ch/create-reason": "GithubWebhook",
		"actuator.nine.ch/branch":        event.PullRequest.Head.GetRef(),
		"actuator.nine.ch/pull-request":  strconv.Itoa(event.PullRequest.GetNumber())}
}
