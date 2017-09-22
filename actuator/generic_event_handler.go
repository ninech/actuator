package actuator

import "github.com/ninech/actuator/github"

// GenericEventHandler is a simple event handler that does nothing
type GenericEventHandler struct{}

func NewGenericEventHandler() *GenericEventHandler {
	return &GenericEventHandler{}
}

// GetEventResponse just returns a generic message. The hook was received but doing nothing in this case.
func (h *GenericEventHandler) GetEventResponse(event *github.Event) *EventResponse {
	return &EventResponse{Message: "Request received. Doing nothing."}
}

// HandleEvent does nothing
func (h *GenericEventHandler) HandleEvent(event *github.Event) {
	return
}
