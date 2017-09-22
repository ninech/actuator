package actuator

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ninech/actuator/github"
)

type WebhookParser interface {
	ValidateAndParseWebhook(*http.Request) (interface{}, error)
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
		Request:       request}
}

// Handle parses the request into a github event and handles it
func (e *EventEndpoint) Handle() (int, interface{}) {
	githubEvent, err := e.WebhookParser.ValidateAndParseWebhook(e.Request)
	if err != nil {
		return 400, gin.H{"message": err.Error()}
	}

	event, ok := github.ConvertGithubEvent(githubEvent)
	if !ok {
		return 400, gin.H{"message": "Invalid or unsupported event payload."}
	}

	dispatcher := EventDispatcher{Event: *event}
	response := dispatcher.GetEventResponse()

	if response.HandleEvent {
		dispatcher.HandleEvent()
	}

	return 200, gin.H{"message": response.Message}
}

func (e *EventEndpoint) getHandlerForEvent(githubEvent interface{}) EventHandler {
	if event, ok := github.ConvertGithubEvent(githubEvent); ok {
		switch event.Type {
		case github.PullRequestEvent:
			return NewPullRequestEventHandler(*event)
		default:
			return NewGenericEventHandler()
		}
	}
	return &GenericEventHandler{}
}
