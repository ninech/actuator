package testutils

import "os/exec"

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
