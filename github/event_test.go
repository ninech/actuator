package github_test

import (
	"testing"

	gh "github.com/google/go-github/github"
	"github.com/ninech/actuator/github"
	"github.com/stretchr/testify/assert"
)

func TestEvent(t *testing.T) {
	t.Run("ConvertGithubEvent", func(t *testing.T) {
		t.Run("for the ping event", func(t *testing.T) {
			originalEvent := &gh.PingEvent{}

			internalEvent, ok := github.ConvertGithubEvent(originalEvent)

			assert.True(t, ok)
			assert.IsType(t, &github.Event{}, internalEvent)
			assert.Equal(t, github.PingEvent, internalEvent.Type)
		})

		t.Run("for a pull request event", func(t *testing.T) {
			number := 1
			action := "opened"
			ownerName := "ninech"
			owner := gh.User{Login: &ownerName}
			fullname := "ninech/actuator"
			name := "actuator"
			branchName := "pr-1"
			repo := gh.Repository{FullName: &fullname, Name: &name, Owner: &owner}
			pr := gh.PullRequest{Number: &number, Head: &gh.PullRequestBranch{Ref: &branchName}}

			originalEvent := &gh.PullRequestEvent{Action: &action, Number: &number, Repo: &repo, PullRequest: &pr}

			internalEvent, ok := github.ConvertGithubEvent(originalEvent)

			assert.True(t, ok)

			assert.Equal(t, action, internalEvent.Action)
			assert.Equal(t, number, internalEvent.IssueNumber)
			assert.Equal(t, name, internalEvent.RepositoryName)
			assert.Equal(t, fullname, internalEvent.RepositoryFullname)
			assert.Equal(t, ownerName, internalEvent.RepositoryOwner)
			assert.Equal(t, branchName, internalEvent.HeadRef)
			assert.Equal(t, github.PullRequestEvent, internalEvent.Type)

			assert.Equal(t, originalEvent, internalEvent.OriginalEvent)
		})

		t.Run("for an unknown event", func(t *testing.T) {
			internalEvent, ok := github.ConvertGithubEvent("yolo")

			assert.False(t, ok)
			assert.Nil(t, internalEvent)
		})
	})
}
