package openshift

import (
	"errors"
	"os/exec"
	"strings"
)

var ExecCommand = exec.Command

// RunOcCommand runs an oc command with the given arguments
func RunOcCommand(args ...string) (string, error) {
	cmd := ExecCommand("oc", args...)
	output, err := cmd.Output()
	return string(output), err
}

// NewAppFromTemplate applies a template using the command `oc new-app`
// It returns the output of the command and an error
func NewAppFromTemplate(templateName string, templateParameters TemplateParameters) (string, error) {
	if templateName == "" {
		return "", errors.New("a template name has to be set")
	}

	arguments := []string{"new-app", "--template", templateName}
	for paramName, paramValue := range templateParameters {
		combinedKeyAndValue := strings.Join([]string{paramName, paramValue}, "=")
		arguments = append(arguments, "--param", combinedKeyAndValue)
	}

	return RunOcCommand(arguments...)
}
