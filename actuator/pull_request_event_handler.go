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
var SupportedPullRequestActions = [2]string{github.EventActionOpened, github.EventActionClosed}

// PullRequestEventHandler handles pull request events
type PullRequestEventHandler struct {
	Event            github.Event
	RepositoryConfig RepositoryConfig
	GithubClient     github.Client
	Openshift        openshift.OpenshiftClient
}

// NewPullRequestEventHandler creates a new event handler for pull requests
func NewPullRequestEventHandler(event github.Event) *PullRequestEventHandler {
	repositoryConfig, _ := Config.GetRepositoryConfig(event.RepositoryFullname)

	return &PullRequestEventHandler{
		Event:            event,
		RepositoryConfig: repositoryConfig,
		GithubClient:     github.NewAuthenticatedGithubClient(Config.GithubAccessToken),
		Openshift:        openshift.NewCommandLineClient()}
}

// GetEventResponse validates the event and checks if it can be handled.
func (h *PullRequestEventHandler) GetEventResponse(event *github.Event) *EventResponse {
	h.Event = *event
	response := &EventResponse{}

	if h.Event.Type != github.PullRequestEvent {
		response.Message = "Invalid event for this handler."
		return response
	}

	if !h.actionIsSupported() {
		response.Message = "Event is not relevant and will be ignored."
		return response
	}

	if !h.RepositoryConfig.Enabled {
		response.Message = fmt.Sprintf("Repository %s is not configured or disabled. Doing nothing.", h.Event.RepositoryFullname)
		return response
	}

	response.Message = fmt.Sprintf("Event for pull request #%d received. Thank you.", h.Event.IssueNumber)
	response.HandleEvent = true
	return response
}

// HandleEvent handles a pull request event from github
func (h *PullRequestEventHandler) HandleEvent(event *github.Event) {
	h.Event = *event
	Logger.Printf("Starting to handle action %v.", h.Event.Action)

	var err error
	switch h.Event.Action {
	case github.EventActionOpened:
		err = h.HandleActionOpened()
		break
	case github.EventActionClosed:
		err = h.DeleteEnvironmentOnOpenshift()
		break
	}

	if err != nil {
		Logger.Printf("There were some errors while handling the event.\n%v", err)
	} else {
		Logger.Printf("%v action handled without errors.", h.Event.Action)
	}
}

func (h *PullRequestEventHandler) HandleActionOpened() error {
	output, err := h.CreateEnvironmentOnOpenshift()
	if err != nil {
		return err
	}

	routeName := output.RouteName()
	comment := h.BuildCommentForRoute(routeName)

	return h.PostCommentOnGithub(comment)
}

func (h *PullRequestEventHandler) CreateEnvironmentOnOpenshift() (*openshift.NewAppOutput, error) {
	template := h.RepositoryConfig.Template
	labels := buildLabelsFromEvent(&h.Event)
	params := buildTemplateParamsFromEvent(&h.Event)

	output, err := h.Openshift.NewAppFromTemplate(template, params, labels)
	if err != nil {
		return output, err
	}

	Logger.Println(output.Raw)
	return output, nil
}

// DeleteEnvironmentOnOpenshift deletes an environment on openshift based on the pull request number
func (h *PullRequestEventHandler) DeleteEnvironmentOnOpenshift() error {
	pullRequestNumber := h.Event.IssueNumber
	labels := openshift.ObjectLabels{"actuator.nine.ch/pull-request": strconv.Itoa(pullRequestNumber)}
	output, err := h.Openshift.DeleteApp(&labels)

	Logger.Println(output.Raw)
	return err
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
