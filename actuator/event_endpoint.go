package actuator

import (
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
)

// EventHandler defines an interface for all event handlers
type EventHandler interface {
	HandleEvent() (string, error)
}

// EventEndpoint is an api endpoint to handle Github event webhooks.
// It needs a parser and an event handler
type EventEndpoint struct {
	WebhookParser WebhookParser
	EventHandler  EventHandler
}

// Handle parses the request into a github event and handles it
func (e *EventEndpoint) Handle() (int, interface{}) {
	if e.WebhookParser == nil {
		panic("A webhook parser must be provided.")
	}

	event, err := e.WebhookParser.ValidateAndParseWebhook()
	if err != nil {
		return 400, gin.H{"message": err.Error()}
	}

	if e.EventHandler == nil {
		eventHandler := e.getHandlerForEvent(event)
		if eventHandler == nil {
			return 200, gin.H{"message": "Not processing this type of event."}
		}
		e.EventHandler = eventHandler
	}

	message, handleError := e.EventHandler.HandleEvent()
	if handleError != nil {
		Logger.Println(handleError.Error())
		return 500, gin.H{"message": message}
	}

	return 200, gin.H{"message": message}
}

func (e *EventEndpoint) getHandlerForEvent(event interface{}) EventHandler {
	switch event := event.(type) {
	case *github.PullRequestEvent:
		return NewPullRequestEventHandler(event, Config)
	default:
		return nil
	}
}
