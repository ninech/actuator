package actuator

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ninech/actuator/github"
)

type WebhookParser interface {
	ValidateAndParseWebhook(*http.Request) (interface{}, error)
}

// EventHandler defines an interface for all event handlers
type EventHandler interface {
	HandleEvent() (string, error)
}

// EventEndpoint is an api endpoint to handle Github event webhooks.
// It needs a parser and an event handler
type EventEndpoint struct {
	WebhookParser WebhookParser
	EventHandler  EventHandler
	Request       *http.Request
}

// NewEventEndpoint produces a new endpoint to handle github events
func NewEventEndpoint(request *http.Request) *EventEndpoint {
	parser := github.NewWebhookParser(Config.GithubWebhookSecret)
	return &EventEndpoint{
		WebhookParser: parser,
		Request:       request,
		EventHandler:  &GenericEventHandler{}}
}

// Handle parses the request into a github event and handles it
func (e *EventEndpoint) Handle() (int, interface{}) {
	event, err := e.WebhookParser.ValidateAndParseWebhook(e.Request)
	if err != nil {
		return 400, gin.H{"message": err.Error()}
	}

	if e.EventHandler == nil {
		e.EventHandler = e.getHandlerForEvent(event)
	}

	message, handleError := e.EventHandler.HandleEvent()
	if handleError != nil {
		Logger.Println(handleError.Error())
		return 500, gin.H{"message": message}
	}

	return 200, gin.H{"message": message}
}

func (e *EventEndpoint) getHandlerForEvent(githubEvent interface{}) EventHandler {
	if event, ok := github.ConvertGithubEvent(githubEvent); ok {
		switch event.Type {
		case github.PullRequestEvent:
			return NewPullRequestEventHandler(event, Config)
		}
	}
	return &GenericEventHandler{}
}
