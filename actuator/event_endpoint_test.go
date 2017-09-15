package actuator_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/ninech/actuator/actuator"
	"github.com/stretchr/testify/assert"
)

func TestNoParserDefined(t *testing.T) {
	endpoint := actuator.EventEndpoint{}
	assert.Panics(t, func() { endpoint.Handle() })
}

func TestValidRequestationFails(t *testing.T) {
	parser := MockGithubWebhookParser{ValidRequest: false}
	endpoint := actuator.EventEndpoint{WebhookParser: &parser}

	code, message := endpoint.Handle()
	assert.Equal(t, http.StatusBadRequest, code)
	assert.Equal(t, gin.H{"message": "Request validation failed."}, message)
}

func TestValidRequestPullRequestEvent(t *testing.T) {
	handler := MockGithubEventHandler{Message: "success!"}
	parser := MockGithubWebhookParser{ValidRequest: true}
	endpoint := actuator.EventEndpoint{
		WebhookParser: &parser,
		EventHandler:  &handler}
	parser.SetEventData(1, "opened")

	code, message := endpoint.Handle()
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, gin.H{"message": handler.Message}, message)
}

func TestUnsupportedEventType(t *testing.T) {
	handler := MockGithubEventHandler{Message: "unsupported!"}
	parser := MockGithubWebhookParser{ValidRequest: true, Event: &github.IssueEvent{}}
	endpoint := actuator.EventEndpoint{
		WebhookParser: &parser,
		EventHandler:  &handler}

	code, message := endpoint.Handle()
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, gin.H{"message": handler.Message}, message)
}

func TestFailingEventHandler(t *testing.T) {
	handler := MockGithubEventHandler{}
	handler.Error = errors.New("something went wrong")
	parser := MockGithubWebhookParser{ValidRequest: true}
	endpoint := actuator.EventEndpoint{
		WebhookParser: &parser,
		EventHandler:  &handler}
	parser.SetEventData(1, "opened")

	code, message := endpoint.Handle()
	assert.Equal(t, http.StatusInternalServerError, code)
	assert.Equal(t, gin.H{"message": "something went wrong"}, message)
}

/// HELPERS ////

type MockGithubWebhookParser struct {
	ValidRequest bool
	Event        interface{}
}

func (p *MockGithubWebhookParser) ValidateAndParseWebhook() (interface{}, error) {
	if p.ValidRequest {
		return p.Event, nil
	}
	return nil, errors.New("Request validation failed.")
}

func (p *MockGithubWebhookParser) SetEventData(number int, action string) {
	p.Event = &github.PullRequestEvent{Number: &number, Action: &action}
}

type MockGithubEventHandler struct {
	Error   error
	Message string
}

func (h *MockGithubEventHandler) HandleEvent() error {
	return h.Error
}

func (h *MockGithubEventHandler) GetMessage() string {
	return h.Message
}
