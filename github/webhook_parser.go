package github

import (
	"net/http"

	"github.com/google/go-github/github"
)

// WebhookParser is responsible for parsing data from github into internal structs
type WebhookParser struct {
	Secret string
}

// NewWebhookParser returns an initialized WebhookParser
func NewWebhookParser(secret string) *WebhookParser {
	return &WebhookParser{Secret: secret}
}

// ValidateAndParseWebhook validates a request from github against the token and
// also parses the data into an internal event struct
func (p *WebhookParser) ValidateAndParseWebhook(request *http.Request) (interface{}, error) {
	payload, err := github.ValidatePayload(request, []byte(p.Secret))
	if err != nil {
		return nil, err
	}

	event, err := github.ParseWebHook(github.WebHookType(request), payload)
	if err != nil {
		return nil, err
	}

	return event, nil
}
