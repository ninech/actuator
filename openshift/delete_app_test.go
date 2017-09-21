package openshift_test

import (
	"testing"

	"github.com/ninech/actuator/openshift"
	"github.com/ninech/actuator/test"
	"github.com/stretchr/testify/assert"
)

func TestDeleteApp(t *testing.T) {
	shell := &test.MockShell{OutputToReturn: "command executed"}
	openshiftClient := openshift.CommandLineClient{CommandExecutor: shell}

	t.Run("runs the command to delete the objects", func(t *testing.T) {
		output, err := openshiftClient.DeleteApp(&openshift.ObjectLabels{"label1": "value1"})

		assert.Nil(t, err)
		assert.Equal(t, "command executed", output.Raw)

		expectedCommandArguments := []string{"delete", "all", "-l", "label1=value1"}
		assert.Equal(t, expectedCommandArguments, shell.ReceivedArguments)
	})
}
