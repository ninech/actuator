package actuator

import (
	"net/http"

	"github.com/google/go-github/github"
)

// GithubToken contains the shared token with which github calculates the HMAC signature
var GithubToken string

// GithubTokenEnvVariable is the name of the environment variable to use for the token
const GithubTokenEnvVariable = "ACTUATOR_GITHUB_TOKEN"

type WebhookParser interface {
	ValidateAndParseWebhook() (interface{}, error)
}

// GithubWebhookParser is responsible for parsing data from github into internal structs
type GithubWebhookParser struct {
	request *http.Request
}

// ValidateAndParseWebhook validates a request from github against the token and
// also parses the data into an internal event struct
func (p *GithubWebhookParser) ValidateAndParseWebhook() (interface{}, error) {
	payload, err := github.ValidatePayload(p.request, []byte(GithubToken))
	if err != nil {
		return nil, err
	}

	event, err := github.ParseWebHook(github.WebHookType(p.request), payload)
	if err != nil {
		return nil, err
	}

	return event, nil
}
