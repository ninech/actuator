package actuator

// GenericEventHandler is a simple event handler that does nothing
type GenericEventHandler struct{}

func NewGenericEventHandler() *GenericEventHandler {
	return &GenericEventHandler{}
}

// GetEventResponse just returns a generic message. The hook was received but doing nothing in this case.
func (h *GenericEventHandler) GetEventResponse() *EventResponse {
	return &EventResponse{Message: "Request received. Doing nothing."}
}

// HandleEvent does nothing
func (h *GenericEventHandler) HandleEvent() {
	return
}
