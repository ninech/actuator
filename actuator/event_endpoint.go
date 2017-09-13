package actuator

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
)

// EventHandler defines an interface for all event handlers
type EventHandler interface {
	HandleEvent() error
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

	switch event := event.(type) {
	case *github.PullRequestEvent:
		if handleError := e.EventHandler.HandleEvent(); handleError == nil {
			message := fmt.Sprintf("Event for pull request #%d received. Thank you.", event.GetNumber())
			return 200, gin.H{"message": message}
		} else {
			return 200, gin.H{"message": handleError.Error()}
		}

	default:
		return 200, gin.H{"message": "Not processing this type of event."}
	}
}
