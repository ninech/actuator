package actuator

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ninech/actuator/github"
	"github.com/ninech/actuator/openshift"
)

// SupportedPullRequestActions defines all pull request event actions which are supported by this app.
const (
	ActionOpened   = "opened"
	ActionClosed   = "closed"
	ActionReopened = "reopened"
)

// SupportedPullRequestActions defines all actions which are currently supported to be handled
var SupportedPullRequestActions = [2]string{ActionOpened, ActionClosed}

// PullRequestEventHandler handles pull request events
type PullRequestEventHandler struct {
	Event        *github.Event
	Config       Configuration
	GithubClient github.Client
	Openshift    openshift.OpenshiftClient
}

// NewPullRequestEventHandler creates a new event handler for pull requests
func NewPullRequestEventHandler(event *github.Event, config Configuration) *PullRequestEventHandler {
	return &PullRequestEventHandler{
		Event:        event,
		Config:       config,
		GithubClient: github.NewAuthenticatedGithubClient(config.GithubAccessToken),
		Openshift:    openshift.NewCommandLineClient()}
}

// HandleEvent handles a pull request event from github
func (h *PullRequestEventHandler) HandleEvent() (string, error) {
	if !h.actionIsSupported() {
		return "Event is not relevant and will be ignored.", nil
	}

	repositoryName := h.Event.RepositoryFullname
	repositoryConfig := h.Config.GetRepositoryConfig(repositoryName)
	if repositoryConfig == nil {
		return fmt.Sprintf("Repository %s is not configured. Doing nothing.", repositoryName), nil
	}

	if !repositoryConfig.Enabled {
		return fmt.Sprintf("Repository %s is disabled. Doing nothing.", repositoryName), nil
	}

	switch h.Event.Action {
	case ActionOpened:
		output, err := h.CreateEnvironmentOnOpenshift(repositoryConfig.Template)
		if err != nil {
			return err.Error(), err
		}

		routeName := output.RouteName()
		comment := h.BuildCommentForRoute(routeName)

		if err := h.PostCommentOnGithub(comment); err != nil {
			return err.Error(), err
		}

		return fmt.Sprintf("Event for pull request #%d received. Thank you.", h.Event.IssueNumber), nil

	case ActionClosed:
		_, err := h.DeleteEnvironmentOnOpenshift()
		if err != nil {
			return err.Error(), err
		}
		return fmt.Sprintf("Event for pull request #%d received. Thank you.", h.Event.IssueNumber), nil

	default:
		return "No handler for this action defined.", errors.New("no action handled")
	}
}

func (h *PullRequestEventHandler) CreateEnvironmentOnOpenshift(template string) (*openshift.NewAppOutput, error) {
	labels := h.buildLabelsFromEvent(h.Event)
	params := h.buildTemplateParamsFromEvent(h.Event)
	output, err := h.Openshift.NewAppFromTemplate(template, params, labels)
	if err != nil {
		return output, err
	}

	Logger.Println(output.Raw)
	return output, nil
}

// DeleteEnvironmentOnOpenshift deletes an environment on openshift based on the pull request number
func (h *PullRequestEventHandler) DeleteEnvironmentOnOpenshift() (*openshift.DeleteAppOutput, error) {
	pullRequestNumber := h.Event.IssueNumber
	labels := openshift.ObjectLabels{"actuator.nine.ch/pull-request": strconv.Itoa(pullRequestNumber)}
	output, err := h.Openshift.DeleteApp(&labels)

	Logger.Println(output.Raw)
	return output, err
}

// PostCommentOnGithub posts a comment on Github, based on data from the event
func (h *PullRequestEventHandler) PostCommentOnGithub(body string) error {
	owner := h.Event.RepositoryOwner
	repo := h.Event.RepositoryName
	issueNumber := h.Event.IssueNumber

	comment, err := h.GithubClient.CreateComment(owner, repo, issueNumber, body)
	if err != nil {
		return err
	}

	Logger.Printf("Created comment on Github: %v.\n", comment.GetHTMLURL())
	return nil
}

// BuildCommentForRoute tries to get the url for the route and compiles a comment
func (h *PullRequestEventHandler) BuildCommentForRoute(routeName string) string {
	var commentBuffer bytes.Buffer
	commentBuffer.WriteString("Your environment is being set-up on Openshift.")

	if routeName != "" {
		url, _ := h.Openshift.GetURLForRoute(routeName)
		commentBuffer.WriteString(" ")
		commentBuffer.WriteString(url)
	} else {
		commentBuffer.WriteString(" There is no route I can point you to.")
	}

	return strings.TrimSpace(commentBuffer.String())
}

// actionIsSupported returns true when the provided action is currently supported by the app
func (h *PullRequestEventHandler) actionIsSupported() bool {
	for _, a := range SupportedPullRequestActions {
		if a == h.Event.Action {
			return true
		}
	}
	return false
}

func (h *PullRequestEventHandler) buildLabelsFromEvent(event *github.Event) openshift.ObjectLabels {
	return openshift.ObjectLabels{
		"actuator.nine.ch/create-reason": "GithubWebhook",
		"actuator.nine.ch/branch":        event.HeadRef,
		"actuator.nine.ch/pull-request":  strconv.Itoa(event.IssueNumber)}
}

func (h *PullRequestEventHandler) buildTemplateParamsFromEvent(event *github.Event) openshift.TemplateParameters {
	return openshift.TemplateParameters{
		"BRANCH_NAME": event.HeadRef}
}
