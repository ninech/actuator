package actuator

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
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
	Event        *github.PullRequestEvent
	Config       Configuration
	GithubClient GithubClient
}

// ApplyOpenshiftTemplate creates new objects in openshift using a template
var ApplyOpenshiftTemplate = openshift.NewAppFromTemplate
var GetURLForRoute = openshift.GetURLForRoute

// NewPullRequestEventHandler creates a new event handler for pull requests
func NewPullRequestEventHandler(event *github.PullRequestEvent, config Configuration) *PullRequestEventHandler {
	return &PullRequestEventHandler{
		Event:        event,
		Config:       config,
		GithubClient: NewAuthenticatedGithubClient()}
}

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
		output, err := h.CreateEnvironmentOnOpenshift(repositoryConfig.Template)
		if err != nil {
			return err.Error(), err
		}

		var commentBuffer bytes.Buffer
		commentBuffer.WriteString("Your environment is being set-up on Openshift.")
		routeName := h.ExtractRouteNameFromNewAppOutput(output)
		if routeName != "" {
			url, _ := GetURLForRoute(routeName)
			commentBuffer.WriteString(url)
		} else {
			commentBuffer.WriteString(" There is no route I can point you to.")
		}
		commentText := commentBuffer.String()

		if err := h.PostCommentOnGithub(commentText); err != nil {
			return err.Error(), err
		}

		return fmt.Sprintf("Event for pull request #%d received. Thank you.", h.Event.GetNumber()), nil
	default:
		return "No handler for this action defined.", errors.New("no action handled")
	}
}

func (h *PullRequestEventHandler) CreateEnvironmentOnOpenshift(template string) (string, error) {
	labels := h.buildLabelsFromEvent(h.Event)
	params := h.buildTemplateParamsFromEvent(h.Event)
	output, err := ApplyOpenshiftTemplate(template, params, labels)
	if err != nil {
		return output, err
	}

	Logger.Println(output)
	return output, nil
}

func (h *PullRequestEventHandler) PostCommentOnGithub(body string) error {
	owner := h.Event.Repo.Owner.GetLogin()
	repo := h.Event.Repo.GetName()
	issueNumber := h.Event.PullRequest.GetNumber()

	if h.GithubClient == nil {
		panic("GithubClient is not set")
	}

	comment, err := h.GithubClient.CreateComment(owner, repo, issueNumber, body)
	if err != nil {
		return err
	}

	Logger.Printf("Created comment on Github: %v.\n", comment.GetHTMLURL())
	return nil
}

// ExtractRouteNameFromNewAppOutput extracts the name of a route from the output of the newapp call
func (h *PullRequestEventHandler) ExtractRouteNameFromNewAppOutput(output string) string {
	r, _ := regexp.Compile(`route "([a-z-]+)" created`)
	matches := r.FindStringSubmatch(output)
	if len(matches) > 1 {
		return r.FindStringSubmatch(output)[1]
	}
	return ""
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

func (h *PullRequestEventHandler) buildTemplateParamsFromEvent(event *github.PullRequestEvent) openshift.TemplateParameters {
	return openshift.TemplateParameters{
		"BRANCH_NAME": event.PullRequest.Head.GetRef()}
}
