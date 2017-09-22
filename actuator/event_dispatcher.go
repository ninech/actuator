package actuator

import "github.com/ninech/actuator/github"

// EventResponse is returned from a handler's GetEventResponse method
// and contains information to be passed to the sender of the hook
type EventResponse struct {
	// Message can be returned to the sender of the hook
	Message string

	// HandleEvent tells the receiver if the event should actually be handled
	// if not, the the HandleEvent method should not be called
	HandleEvent bool
}

// EventHandler is an interface for every event handler in the system
type EventHandler interface {
	// GetEventResponse gets a response for the sender of the hook. This method doesn't have
	// to do the heavy work as in most cases the caller of the endpoint doesn't care if the work
	// involved succeeded or not. It's just interested if the hook was received correctly and
	// can be handled.
	GetEventResponse(event *github.Event) *EventResponse

	// HandleEvent does the actual work of handling an event. The result of this operation will
	// be logged but not returned to the caller. This makes it possible to do the heavy work
	// in the background and get a fast response for the caller of the hook.
	HandleEvent(event *github.Event)
}

// EventDispatcher is responsible to find the right handler for an incoming
// Github event. It then forwards the request to the respective handler.
type EventDispatcher struct {
	LastEventHandler  EventHandler
	LastEventResponse *EventResponse
}

// GetEventResponse looks for an appropriate handler to care for this event. It then
// calls its GetEventResponse method and returns the data.
func (d *EventDispatcher) GetEventResponse(event *github.Event) *EventResponse {
	handler := d.FindEventHandler(event)
	response := handler.GetEventResponse(event)

	d.LastEventHandler = handler
	d.LastEventResponse = response

	return response
}

// HandleEvent dispatches the call to the last used event handler
func (d *EventDispatcher) HandleEvent(event *github.Event) {
	if d.LastEventHandler != nil {
		d.LastEventHandler.HandleEvent(event)
	}
}

// FindEventHandler finds a handler for the provided event
func (d *EventDispatcher) FindEventHandler(event *github.Event) EventHandler {
	switch event.Type {
	case github.PullRequestEvent:
		return NewPullRequestEventHandler(*event)
	default:
		return NewGenericEventHandler()
	}
}
