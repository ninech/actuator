package github

import (
	"context"

	gh "github.com/google/go-github/github"

	"golang.org/x/oauth2"
)

// Client defines the interface for a github client
type Client interface {
	CreateComment(owner string, repo string, issueNumber int, body string) (*gh.IssueComment, error)
}

// AuthenticatedGithubClient abstracts the github client from google's client
type AuthenticatedGithubClient struct {
	Client *gh.Client
	ctx    context.Context
}

// NewAuthenticatedGithubClient produces a new, authenticated github client
func NewAuthenticatedGithubClient(accessToken string) *AuthenticatedGithubClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(ctx, ts)

	client := gh.NewClient(tc)

	return &AuthenticatedGithubClient{Client: client, ctx: ctx}
}

// CreateComment creates a new comment on github
func (agc *AuthenticatedGithubClient) CreateComment(owner string, repo string, issueNumber int, body string) (*gh.IssueComment, error) {
	comment := gh.IssueComment{Body: &body}
	issueComment, _, error := agc.Client.Issues.CreateComment(agc.ctx, owner, repo, issueNumber, &comment)
	return issueComment, error
}
