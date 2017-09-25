package test

import (
	"github.com/ninech/actuator/actuator"
	"github.com/ninech/actuator/github"
)

func NewMockEventHandler(message string) *MockEventHandler {
	response := actuator.EventResponse{Message: message}
	done := make(chan bool, 1)
	return &MockEventHandler{EventResponse: response, Done: done}
}

type MockEventHandler struct {
	EventResponse   actuator.EventResponse
	EventWasHandled bool
	LastEvent       *github.Event
	Done            chan bool
}

func (h *MockEventHandler) GetEventResponse(event *github.Event) *actuator.EventResponse {
	h.LastEvent = event
	return &h.EventResponse
}

func (h *MockEventHandler) HandleEvent(event *github.Event) {
	h.LastEvent = event
	h.EventWasHandled = true
	h.Done <- true
}
