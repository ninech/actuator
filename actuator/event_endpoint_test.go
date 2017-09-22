package actuator_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	gh "github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"

	"github.com/ninech/actuator/actuator"
	"github.com/ninech/actuator/test"
)

func TestEventEndpoint(t *testing.T) {
	t.Run("Handle", func(t *testing.T) {
		t.Run("validate request fails", func(t *testing.T) {
			parser := MockWebhookParser{ValidRequest: false}
			endpoint := actuator.EventEndpoint{WebhookParser: &parser}

			code, message := endpoint.Handle()
			assert.Equal(t, http.StatusBadRequest, code)
			assert.Equal(t, gin.H{"message": "Request validation failed."}, message)
		})

		t.Run("with an unsupported event type", func(t *testing.T) {
			parser := MockWebhookParser{ValidRequest: true, Event: &gh.IssueEvent{}}
			endpoint := actuator.EventEndpoint{
				WebhookParser: &parser}

			code, message := endpoint.Handle()
			assert.Equal(t, gin.H{"message": "Invalid or unsupported event payload."}, message)
			assert.Equal(t, http.StatusBadRequest, code)
		})

		t.Run("with a valid event", func(t *testing.T) {
			parser := MockWebhookParser{ValidRequest: true}
			parser.Event = test.NewDefaultOriginalPullRequestEvent(1, "opened")

			handler := test.NewMockEventHandler("All is fine!")
			handler.EventResponse.HandleEvent = true

			endpoint := actuator.EventEndpoint{
				WebhookParser: &parser,
				EventHandler:  handler}

			code, message := endpoint.Handle()
			assert.Equal(t, gin.H{"message": "All is fine!"}, message)
			assert.Equal(t, http.StatusOK, code)
			assert.True(t, handler.EventWasHandled)
		})
	})
}

/// HELPERS ////

type MockWebhookParser struct {
	ValidRequest bool
	Event        interface{}
}

func (p *MockWebhookParser) ValidateAndParseWebhook(request *http.Request) (interface{}, error) {
	if p.ValidRequest {
		return p.Event, nil
	}
	return nil, errors.New("Request validation failed.")
}
