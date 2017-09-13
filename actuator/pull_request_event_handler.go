package actuator

import (
	"fmt"

	"github.com/google/go-github/github"
)

// SupportedPullRequestActions defines all pull request event actions which are supported by this app.
var SupportedPullRequestActions = [...]string{"opened", "closed", "reopened"}

// PullRequestEventHandler handles pull request events
type PullRequestEventHandler struct {
	Event   *github.PullRequestEvent
	Message string
}

// GetMessage returns the end message of this handler to be sent to the client
func (h *PullRequestEventHandler) GetMessage() string {
	if h.Event != nil && h.Message == "" {
		return fmt.Sprintf("Event for pull request #%d received. Thank you.", h.Event.GetNumber())
	}
	return h.Message
}

// HandleEvent handles a pull request event from github
func (h *PullRequestEventHandler) HandleEvent() error {
	if !h.actionIsSupported() {
		h.Message = "Event is not relevant and will be ignored."
		return nil
	}

	// TODO: Implement the handling of this specific event

	return nil
}

// actionIsSupported returns true when the provided action is currently supported by the app
func (h *PullRequestEventHandler) actionIsSupported() bool {
	for _, a := range SupportedPullRequestActions {
		if a == h.Event.GetAction() {
			return true
		}
	}
	return false
}
