package testutils

import (
	"fmt"

	"github.com/google/go-github/github"
)

type MockGithubClient struct {
	LastComment *github.IssueComment
}

func NewMockGithubClient() *MockGithubClient {
	return &MockGithubClient{}
}

// CreateComment returns a fake comment and saves it as LastComment
func (tgc *MockGithubClient) CreateComment(owner string, repo string, issueNumber int, body string) (*github.IssueComment, error) {
	url := fmt.Sprintf("https://github.com/%v/%v/issues/%d#issuecomment-330230087", owner, repo, issueNumber)
	tgc.LastComment = &github.IssueComment{HTMLURL: &url, Body: &body}
	return tgc.LastComment, nil
}
