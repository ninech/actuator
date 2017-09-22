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
		EventHandler:  &handler}
	parser.SetEventData(1, "opened")

	code, message := endpoint.Handle()
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, gin.H{"message": handler.Message}, message)
}

func TestUnsupportedEventType(t *testing.T) {
	handler := test.NewMockEventHandler("unsupported!")
	parser := MockWebhookParser{ValidRequest: true, Event: &github.IssueEvent{}}
	endpoint := actuator.EventEndpoint{
		WebhookParser: &parser,
		EventHandler:  &handler}

	code, message := endpoint.Handle()
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, gin.H{"message": handler.Message}, message)
}

func TestFailingEventHandler(t *testing.T) {
	test.DisableLogging()

	handler := MockGithubEventHandler{Error: errors.New("something went wrong")}
	parser := MockWebhookParser{ValidRequest: true}
	endpoint := actuator.EventEndpoint{
		WebhookParser: &parser,
		EventHandler:  &handler}
	parser.SetEventData(1, "opened")

	code, message := endpoint.Handle()
	assert.Equal(t, http.StatusInternalServerError, code)
	assert.Equal(t, gin.H{"message": "something went wrong"}, message)
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
