package testutils

import (
	"strings"

	"github.com/google/go-github/github"

	"github.com/ninech/actuator/actuator"
)

func NewTestEvent(number int, action string, repoName string) *github.PullRequestEvent {
	ownerName := "ninech"
	owner := github.User{Login: &ownerName}
	name := strings.Split(repoName, "/")[1]
	repo := github.Repository{FullName: &repoName, Name: &name, Owner: &owner}

	branchName := "pr-1"
	prNumber := 1408
	pr := github.PullRequest{Number: &prNumber, Head: &github.PullRequestBranch{Ref: &branchName}}

	return &github.PullRequestEvent{Action: &action, Number: &number, Repo: &repo, PullRequest: &pr}
}

func NewDefaultTestEvent() *github.PullRequestEvent {
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
