package openshift

import (
	"errors"
	"fmt"
	"os/exec"
)

var ExecCommand = exec.Command

// RunOcCommand runs an oc command with the given arguments
func RunOcCommand(args ...string) (string, error) {
	cmd := ExecCommand("oc", args...)

	output, err := cmd.Output()
	if exitError, ok := err.(*exec.ExitError); ok {
		err = errors.New(string(exitError.Stderr))
	}

	return string(output), err
}

// NewAppFromTemplate applies a template using the command `oc new-app`
// It returns the output of the command and an error
// labels defines some lables which are applied to all created objects
func NewAppFromTemplate(templateName string, templateParameters TemplateParameters, labels ObjectLabels) (string, error) {
	if templateName == "" {
		return "", errors.New("a template name has to be set")
	}

	arguments := []string{"new-app", "--template", templateName}
	arguments = appendKeyValueArgument(arguments, "--param", templateParameters)
	arguments = append(arguments, "--labels", labels.Combined())

	return RunOcCommand(arguments...)
}

func appendKeyValueArgument(appendTarget []string, argumentName string, keyValuePairs map[string]string) []string {
	for key, value := range keyValuePairs {
		combinedKeyAndValue := fmt.Sprintf("%s=%s", key, value)
		appendTarget = append(appendTarget, argumentName, combinedKeyAndValue)
	}
	return appendTarget
}
