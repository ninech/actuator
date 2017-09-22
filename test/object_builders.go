package test

import (
	"strings"

	"github.com/ninech/actuator/actuator"
	"github.com/ninech/actuator/github"
)

func NewTestEvent(number int, action string, repoName string) *github.Event {
	return &github.Event{
		Action:             action,
		IssueNumber:        number,
		RepositoryName:     strings.Split(repoName, "/")[1],
		RepositoryFullname: repoName,
		RepositoryOwner:    "ninech",
		HeadRef:            "pr-1",
		Type:               github.PullRequestEvent}
}

func NewDefaultTestEvent() *github.Event {
	return NewTestEvent(1, actuator.ActionOpened, "ninech/actuator")
}

func NewDefaultConfig() actuator.Configuration {
	repoConfig := actuator.RepositoryConfig{
		Fullname: "ninech/actuator",
		Enabled:  true,
		Exclude:  "master",
		Template: "actuator-template"}
	return actuator.Configuration{Repositories: []actuator.RepositoryConfig{repoConfig}}
}
