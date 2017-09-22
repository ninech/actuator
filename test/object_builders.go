package test

import (
	"strings"

	gh "github.com/google/go-github/github"

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
	return NewTestEvent(1, github.EventActionOpened, "ninech/actuator")
}

func NewDefaultConfig() actuator.Configuration {
	repoConfig := actuator.RepositoryConfig{
		Fullname: "ninech/actuator",
		Enabled:  true,
		Exclude:  "master",
		Template: "actuator-template"}
	return actuator.Configuration{Repositories: []actuator.RepositoryConfig{repoConfig}}
}

func NewDefaultOriginalPullRequestEvent(number int, action string) *gh.PullRequestEvent {
	login := "ninech"
	name := "actuator"
	fullname := "ninech/actuator"
	owner := gh.User{Login: &login}
	repo := gh.Repository{Owner: &owner, Name: &name, FullName: &fullname}
	branch := "pr-1"
	head := gh.PullRequestBranch{Ref: &branch}
	pr := gh.PullRequest{Head: &head}
	return &gh.PullRequestEvent{
		Number:      &number,
		Action:      &action,
		Repo:        &repo,
		PullRequest: &pr}
}
