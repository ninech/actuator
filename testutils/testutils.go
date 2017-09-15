package testutils

import (
	"github.com/google/go-github/github"

	"github.com/ninech/actuator/actuator"
)

func NewTestEvent(number int, action string, repoName string) *github.PullRequestEvent {
	repo := github.Repository{FullName: &repoName}
	return &github.PullRequestEvent{Action: &action, Number: &number, Repo: &repo}
}

func NewDefaultTestEvent() *github.PullRequestEvent {
	return NewTestEvent(1, actuator.ActionOpened, "ninech/actuator")
}

func NewDefaultConfig() actuator.Configuration {
	repoConfig := actuator.RepositoryConfig{Fullname: "ninech/actuator", Enabled: true, Exclude: "master"}
	return actuator.Configuration{Repositories: []actuator.RepositoryConfig{repoConfig}}
}
