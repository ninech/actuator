package actuator_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/ninech/actuator/actuator"
	"github.com/ninech/actuator/test"
	"github.com/stretchr/testify/assert"
)

func TestValidateRequestFails(t *testing.T) {
	parser := MockWebhookParser{ValidRequest: false}
	endpoint := actuator.EventEndpoint{WebhookParser: &parser}

	code, message := endpoint.Handle()
	assert.Equal(t, http.StatusBadRequest, code)
	assert.Equal(t, gin.H{"message": "Request validation failed."}, message)
}

func TestValidRequestPullRequestEvent(t *testing.T) {
	handler := test.NewMockEventHandler("success!")
	parser := MockWebhookParser{ValidRequest: true}
	endpoint := actuator.EventEndpoint{
		WebhookParser: &parser,
		EventHandler:  handler}
	parser.SetEventData(1, "opened")

	t.Skip("validate test")

	code, message := endpoint.Handle()
	assert.True(t, handler.EventWasHandled)
	assert.Equal(t, "tbd", message)
	assert.Equal(t, http.StatusOK, code)
}

func TestUnsupportedEventType(t *testing.T) {
	handler := test.NewMockEventHandler("unsupported!")
	parser := MockWebhookParser{ValidRequest: true, Event: &github.IssueEvent{}}
	endpoint := actuator.EventEndpoint{
		WebhookParser: &parser,
		EventHandler:  handler}

	code, message := endpoint.Handle()
	assert.False(t, handler.EventWasHandled)
	assert.Equal(t, gin.H{"message": "Invalid or unsupported event payload."}, message)
	assert.Equal(t, http.StatusBadRequest, code)
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

func (p *MockWebhookParser) SetEventData(number int, action string) {
	p.Event = &github.PullRequestEvent{Number: &number, Action: &action}
}
