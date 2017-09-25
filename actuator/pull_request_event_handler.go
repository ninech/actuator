package actuator

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/ninech/actuator/github"
	"github.com/ninech/actuator/openshift"
)

// SupportedPullRequestActions defines all actions which are currently supported to be handled
var SupportedPullRequestActions = [3]string{
	github.EventActionOpened,
	github.EventActionClosed,
	github.EventActionReopened}

// PullRequestEventHandler handles pull request events
type PullRequestEventHandler struct {
	RepositoryConfig RepositoryConfig
	GithubClient     github.Client
	Openshift        openshift.OpenshiftClient
}

// NewPullRequestEventHandler creates a new event handler for pull requests
func NewPullRequestEventHandler(event github.Event) *PullRequestEventHandler {
	repositoryConfig, _ := Config.GetRepositoryConfig(event.RepositoryFullname)

	return &PullRequestEventHandler{
		RepositoryConfig: repositoryConfig,
		GithubClient:     github.NewAuthenticatedGithubClient(Config.GithubAccessToken),
		Openshift:        openshift.NewCommandLineClient()}
}

// GetEventResponse validates the event and checks if it can be handled.
func (h *PullRequestEventHandler) GetEventResponse(event *github.Event) *EventResponse {
	response := &EventResponse{}

	if event.Type != github.PullRequestEvent {
		response.Message = "Invalid event for this handler."
		return response
	}

	if !isActionSupported(event.Action) {
		response.Message = "Event is not relevant and will be ignored."
		return response
	}

	if !h.RepositoryConfig.Enabled {
		response.Message = fmt.Sprintf("Repository %s is not configured or disabled. Doing nothing.", event.RepositoryFullname)
		return response
	}

	response.Message = fmt.Sprintf("Event for pull request #%d received. Thank you.", event.IssueNumber)
	response.HandleEvent = true
	return response
}

// HandleEvent handles a pull request event from github
func (h *PullRequestEventHandler) HandleEvent(event *github.Event) {
	Logger.Printf("Starting to handle action %v.", event.Action)

	var err error
	switch event.Action {
	case github.EventActionOpened, github.EventActionReopened:
		err = h.HandleActionOpened(event)
	case github.EventActionClosed:
		labels := openshift.ObjectLabels{"actuator.nine.ch/pull-request": strconv.Itoa(event.IssueNumber)}
		err = h.DeleteEnvironmentOnOpenshift(&labels)
	}

	if err != nil {
		Logger.Printf("There were some errors while handling the event.\n%v", err)
	} else {
		Logger.Printf("%v action handled without errors.", event.Action)
	}
}

// HandleActionOpened is called when we receive an opened or reopened event from Github
func (h *PullRequestEventHandler) HandleActionOpened(event *github.Event) error {
	output, err := h.CreateEnvironmentOnOpenshift(event)
	if err != nil {
		return err
	}

	routeName := output.RouteName()
	comment := h.BuildCommentForRoute(routeName)

	return h.PostCommentOnGithub(event, comment)
}

// CreateEnvironmentOnOpenshift does everything necessary to create the environment on openshift
func (h *PullRequestEventHandler) CreateEnvironmentOnOpenshift(event *github.Event) (*openshift.NewAppOutput, error) {
	template := h.RepositoryConfig.Template
	labels := buildLabelsFromEvent(event)
	params := buildTemplateParamsFromEvent(event)

	output, err := h.Openshift.NewAppFromTemplate(template, params, labels)
	if err != nil {
		return output, err
	}

	Logger.Println(output.Raw)
	return output, nil
}

// DeleteEnvironmentOnOpenshift deletes an environment on openshift based on the pull request number
func (h *PullRequestEventHandler) DeleteEnvironmentOnOpenshift(labels *openshift.ObjectLabels) error {
	output, err := h.Openshift.DeleteApp(labels)
	Logger.Println(output.Raw)
	return err
}

// PostCommentOnGithub posts a comment on Github, based on data from the event
func (h *PullRequestEventHandler) PostCommentOnGithub(event *github.Event, body string) error {
	owner := event.RepositoryOwner
	repo := event.RepositoryName
	issueNumber := event.IssueNumber

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
func isActionSupported(action string) bool {
	for _, a := range SupportedPullRequestActions {
		if a == action {
			return true
		}
	}
	return false
}

func buildLabelsFromEvent(event *github.Event) openshift.ObjectLabels {
	return openshift.ObjectLabels{
		"actuator.nine.ch/create-reason": "GithubWebhook",
		"actuator.nine.ch/branch":        event.HeadRef,
		"actuator.nine.ch/pull-request":  strconv.Itoa(event.IssueNumber)}
}

func buildTemplateParamsFromEvent(event *github.Event) openshift.TemplateParameters {
	return openshift.TemplateParameters{
		"BRANCH_NAME": event.HeadRef}
}
