package openshift_test

import (
	"testing"

	"github.com/ninech/actuator/openshift"
	"github.com/ninech/actuator/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNewAppFromTemplate(t *testing.T) {
	shell := &testutils.MockShell{OutputToReturn: "command executed"}
	openshiftClient := openshift.CommandLineClient{CommandExecutor: shell}

	t.Run("empty template name", func(t *testing.T) {
		_, err := openshiftClient.NewAppFromTemplate("", openshift.TemplateParameters{}, openshift.ObjectLabels{})

		assert.Equal(t, "a template name has to be set", err.Error())
	})

	t.Run("runs the command", func(t *testing.T) {
		output, err := openshiftClient.NewAppFromTemplate(
			"actuator",
			openshift.TemplateParameters{"PARAM1": "yolo"},
			openshift.ObjectLabels{"label1": "value1"})

		assert.Nil(t, err)
		assert.Equal(t, "command executed", output.Raw)

		expectedCommandArguments := []string{"new-app", "--template", "actuator", "--param", "PARAM1=yolo", "--labels", "label1=value1"}
		assert.Equal(t, expectedCommandArguments, shell.ReceivedArguments)
	})
}

func TestExtractRouteNameFromNewAppOutput(t *testing.T) {
	t.Run("when there is a route to get", func(t *testing.T) {
		cmdOutput := `--> Creating resources with label actuator.nine.ch/branch=changes ...
		route "actuator-changes" created
		configmap "actuator-test-5ybwccanc4" created`
		output := openshift.NewAppOutput{Raw: cmdOutput}

		assert.Equal(t, "actuator-changes", output.RouteName())
	})

	t.Run("when there is no route", func(t *testing.T) {
		cmdOutput := `no route at all`
		output := openshift.NewAppOutput{Raw: cmdOutput}

		assert.Empty(t, output.RouteName())
	})
}
