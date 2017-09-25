package actuator

import "github.com/ninech/actuator/github"

// PingEventHandler responds to ping events from github
// https://developer.github.com/webhooks/#ping-event
type PingEventHandler struct{}

// NewPingEventHandler creates an empty PingEventHandler
func NewPingEventHandler() *PingEventHandler {
	return &PingEventHandler{}
}

// GetEventResponse just returns a generic message. The hook was received but doing nothing in this case.
func (h *PingEventHandler) GetEventResponse(event *github.Event) *EventResponse {
	return &EventResponse{Message: "Request received. Thank you."}
}

// HandleEvent does nothing
func (h *PingEventHandler) HandleEvent(event *github.Event) {
	return
}
