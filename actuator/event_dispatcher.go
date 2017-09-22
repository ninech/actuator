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

	// Data contains some additional data which was gathered during the answering process
	Data map[string]interface{}
}

// EventHandler is an interface for every event handler in the system
type EventHandler interface {
	// GetEventResponse gets a response for the sender of the hook. This method doesn't have
	// to do the heavy work as in most cases the caller of the endpoint doesn't care if the work
	// involved succeeded or not. It's just interested if the hook was received correctly and
	// can be handled.
	GetEventResponse() *EventResponse

	// HandleEvent does the actual work of handling an event. The result of this operation will
	// be logged but not returned to the caller. This makes it possible to do the heavy work
	// in the background and get a fast response for the caller of the hook.
	HandleEvent()
}

// EventDispatcher is responsible to find the right handler for an incoming
// Github event. It then forwards the request to the respective handler.
type EventDispatcher struct {
	Event             github.Event
	LastEventHandler  EventHandler
	LastEventResponse *EventResponse
}

// GetEventResponse looks for an appropriate handler to care for this event. It then
// calls its GetEventResponse method and returns the data.
func (d *EventDispatcher) GetEventResponse() *EventResponse {
	handler := d.FindEventHandler()
	response := handler.GetEventResponse()

	d.LastEventHandler = handler
	d.LastEventResponse = response

	return response
}

// HandleEvent dispatches the call to the last used event handler
func (d *EventDispatcher) HandleEvent() {
	if d.LastEventHandler != nil {
		d.LastEventHandler.HandleEvent()
	}
}

// FindEventHandler finds a handler for the provided event
func (d *EventDispatcher) FindEventHandler() EventHandler {
	switch d.Event.Type {
	case github.PullRequestEvent:
		return NewPullRequestEventHandler(d.Event)
	default:
		return NewGenericEventHandler()
	}
}
