package test

import "github.com/ninech/actuator/actuator"

func NewMockEventHandler(message string) *MockEventHandler {
	response := actuator.EventResponse{Message: message}
	return &MockEventHandler{EventResponse: response}
}

type MockEventHandler struct {
	EventResponse   actuator.EventResponse
	EventWasHandled bool
}

func (h *MockEventHandler) GetEventResponse() *actuator.EventResponse {
	return &h.EventResponse
}

func (h *MockEventHandler) HandleEvent() {
	h.EventWasHandled = true
}
