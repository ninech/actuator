package openshift_test

import (
	"os/exec"
	"testing"

	"github.com/ninech/actuator/openshift"
	"github.com/ninech/actuator/testutils"
	"github.com/stretchr/testify/assert"
)

func TestRunOcCommand(t *testing.T) {
	t.Run("command succeeds", func(t *testing.T) {
		executer := testutils.NewDefaultTestCommandExecuter()

		openshift.ExecCommand = executer.ExecCommand
		defer func() { openshift.ExecCommand = exec.Command }()

		output, err := openshift.RunOcCommand("additional", "args")

		assert.Nil(t, err)
		assert.Equal(t, "command executed\n", output, "should return the command output")
		assert.Equal(t, "oc", executer.CommandReceived, "should run `oc` as the command")
		assert.Equal(t, []string{"additional", "args"}, executer.ArgsReceived, "should receive the real oc arguments")
	})

	t.Run("command fails", func(t *testing.T) {
		failingExecuter := testutils.TestCommandExecuter{
			TestCommandToRun:     "cat",
			TestCommandArguments: []string{"/tmp/nonexistingfile.txt"}}

		openshift.ExecCommand = failingExecuter.ExecCommand
		defer func() { openshift.ExecCommand = exec.Command }()

		_, err := openshift.RunOcCommand("")

		assert.NotNil(t, err)
		assert.Equal(t, "cat: /tmp/nonexistingfile.txt: No such file or directory\n", err.Error())
	})
}

func TestNewAppFromTemplate(t *testing.T) {
	executer := testutils.NewDefaultTestCommandExecuter()
	openshift.ExecCommand = executer.ExecCommand
	defer func() { openshift.ExecCommand = exec.Command }()

	t.Run("empty template name", func(t *testing.T) {
		_, err := openshift.NewAppFromTemplate("", openshift.TemplateParameters{}, openshift.ObjectLabels{})

		assert.Equal(t, "a template name has to be set", err.Error())
	})

	t.Run("runs the command", func(t *testing.T) {
		output, err := openshift.NewAppFromTemplate(
			"actuator",
			openshift.TemplateParameters{"PARAM1": "yolo"},
			openshift.ObjectLabels{"label1": "value1", "label2": "value2"})

		assert.Nil(t, err)
		assert.Equal(t, "command executed\n", string(output))
		assert.Equal(t, "oc", executer.CommandReceived)

		expectedCommandArguments := []string{"new-app", "--template", "actuator", "--param", "PARAM1=yolo", "--labels", "label1=value1,label2=value2"}
		assert.Equal(t, expectedCommandArguments, executer.ArgsReceived)
	})
}
