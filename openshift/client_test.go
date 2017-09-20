package openshift_test

import (
	"testing"

	"github.com/ninech/actuator/openshift"
	"github.com/ninech/actuator/test"
	"github.com/stretchr/testify/assert"
)

func TestRunOcCommand(t *testing.T) {
	shell := &test.MockShell{OutputToReturn: "command executed"}
	openshiftClient := openshift.CommandLineClient{CommandExecutor: shell}

	output, err := openshiftClient.RunOcCommand("additional", "args")

	assert.Nil(t, err)
	assert.Equal(t, "command executed", output)
	assert.Equal(t, []string{"additional", "args"}, shell.ReceivedArguments)
}
