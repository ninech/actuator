package openshift_test

import (
	"testing"

	"github.com/ninech/actuator/openshift"
	"github.com/ninech/actuator/testutils"
	"github.com/stretchr/testify/assert"
)

func TestRunOcCommand(t *testing.T) {
	t.Run("command succeeds", func(t *testing.T) {
		shell := &testutils.MockShell{OutputToReturn: "command executed"}
		openshift.CommandExecutor = shell

		output, err := openshift.RunOcCommand("additional", "args")

		assert.Nil(t, err)
		assert.Equal(t, "command executed", output)
		assert.Equal(t, []string{"additional", "args"}, shell.ReceivedArguments)
	})
}

func TestNewAppFromTemplate(t *testing.T) {
	shell := &testutils.MockShell{OutputToReturn: "command executed"}
	openshift.CommandExecutor = shell

	t.Run("empty template name", func(t *testing.T) {
		_, err := openshift.NewAppFromTemplate("", openshift.TemplateParameters{}, openshift.ObjectLabels{})

		assert.Equal(t, "a template name has to be set", err.Error())
	})

	t.Run("runs the command", func(t *testing.T) {
		output, err := openshift.NewAppFromTemplate(
			"actuator",
			openshift.TemplateParameters{"PARAM1": "yolo"},
			openshift.ObjectLabels{"label1": "value1"})

		assert.Nil(t, err)
		assert.Equal(t, "command executed", string(output))

		expectedCommandArguments := []string{"new-app", "--template", "actuator", "--param", "PARAM1=yolo", "--labels", "label1=value1"}
		assert.Equal(t, expectedCommandArguments, shell.ReceivedArguments)
	})
}
