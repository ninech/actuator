package testutils

import (
	"os/exec"

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
	repoConfig := actuator.RepositoryConfig{
		Fullname: "ninech/actuator",
		Enabled:  true,
		Exclude:  "master",
		Template: "actuator-template"}
	return actuator.Configuration{Repositories: []actuator.RepositoryConfig{repoConfig}}
}

///////// TestCommandExecuter ////////////

// TestCommandExecuter abstracts system calls and makes it possible to test them
type TestCommandExecuter struct {
	TestCommandToRun     string
	TestCommandArguments []string

	// The following values will be set by ExecCommand and can be checked with assert
	CommandReceived string
	ArgsReceived    []string
}

func NewDefaultTestCommandExecuter() TestCommandExecuter {
	return TestCommandExecuter{
		TestCommandToRun:     "echo",
		TestCommandArguments: []string{"command executed"}}
}

// ExecCommand runs some predefined command but also remembers provided parameters
func (c *TestCommandExecuter) ExecCommand(cmd string, args ...string) *exec.Cmd {
	c.CommandReceived = cmd
	c.ArgsReceived = args

	return exec.Command(c.TestCommandToRun, c.TestCommandArguments...)
}
