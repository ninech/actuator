package test

import (
	"github.com/ninech/actuator/actuator"
	"github.com/ninech/actuator/github"
)

func NewMockEventHandler(message string) *MockEventHandler {
	response := actuator.EventResponse{Message: message}
	return &MockEventHandler{EventResponse: response}
}

type MockEventHandler struct {
	EventResponse   actuator.EventResponse
	EventWasHandled bool
	LastEvent       *github.Event
}

func (h *MockEventHandler) GetEventResponse(event *github.Event) *actuator.EventResponse {
	h.LastEvent = event
	return &h.EventResponse
}

func (h *MockEventHandler) HandleEvent(event *github.Event) {
	h.LastEvent = event
	h.EventWasHandled = true
}
