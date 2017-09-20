package openshift

import (
	"errors"
	"os/exec"
)

type Shell interface {
	RunWithArgs(args ...string) (string, error)
}

type OpenshiftShell struct{}

func (s OpenshiftShell) RunWithArgs(args ...string) (string, error) {
	cmd := exec.Command("oc", args...)

	output, err := cmd.Output()
	if exitError, ok := err.(*exec.ExitError); ok {
		err = errors.New(string(exitError.Stderr))
	}

	return string(output), err
}
