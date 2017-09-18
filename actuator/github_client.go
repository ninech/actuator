package actuator

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GithubClient defines the interface for a github client
type GithubClient interface {
	CreateComment(owner string, repo string, issueNumber int, body string) (*github.IssueComment, error)
}

// AuthenticatedGithubClient abstracts the github client from google's client
type AuthenticatedGithubClient struct {
	Client *github.Client
	ctx    context.Context
}

// NewAuthenticatedGithubClient produces a new, authenticated github client
func NewAuthenticatedGithubClient() *AuthenticatedGithubClient {
	fmt.Printf("Creating github client with token %s\n", Config.GithubAccessToken)
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: Config.GithubAccessToken})
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return &AuthenticatedGithubClient{Client: client, ctx: ctx}
}

// CreateComment creates a new comment on github
func (agc *AuthenticatedGithubClient) CreateComment(owner string, repo string, issueNumber int, body string) (*github.IssueComment, error) {
	comment := github.IssueComment{Body: &body}
	issueComment, _, error := agc.Client.Issues.CreateComment(agc.ctx, owner, repo, issueNumber, &comment)
	return issueComment, error
}
