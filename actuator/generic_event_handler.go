package actuator

// GenericEventHandler is a simple event handler that does nothing
type GenericEventHandler struct{}

// HandleEvent just return a message that it does nothing
func (h *GenericEventHandler) HandleEvent() (string, error) {
	return "Not processing this type of event", nil
}
