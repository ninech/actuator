package github

import (
	gh "github.com/google/go-github/github"
)

type EventType int

const (
	PullRequestEvent EventType = iota
	PingEvent
)

// Event is the internal structure used for an event
type Event struct {
	Action             string
	IssueNumber        int
	RepositoryName     string
	RepositoryFullname string
	RepositoryOwner    string
	HeadRef            string
	Type               EventType

	OriginalEvent interface{}
}

// SupportedPullRequestActions defines all pull request event actions which are supported by this app.
const (
	EventActionOpened   = "opened"
	EventActionClosed   = "closed"
	EventActionReopened = "reopened"
)

// ConvertGithubEvent turns an original Github Event into the internal structure
func ConvertGithubEvent(original interface{}) (*Event, bool) {
	switch event := original.(type) {
	case *gh.PullRequestEvent:
		return convertPullRequestEvent(event), true
	case *gh.PingEvent:
		return &Event{Type: PingEvent, OriginalEvent: original}, true
	default:
		return nil, false
	}
}

// ConvertPullRequestEvent turns an original PullRequestEvent into the internal structure
func convertPullRequestEvent(original *gh.PullRequestEvent) *Event {
	return &Event{
		IssueNumber:        original.GetNumber(),
		Action:             original.GetAction(),
		RepositoryName:     original.Repo.GetName(),
		RepositoryFullname: original.Repo.GetFullName(),
		RepositoryOwner:    original.Repo.Owner.GetLogin(),
		HeadRef:            original.PullRequest.Head.GetRef(),
		Type:               PullRequestEvent,

		OriginalEvent: original}
}
