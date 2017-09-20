package openshift

import (
	"errors"
	"fmt"
	"regexp"
)

// NewAppOutput represents the output of a call to `oc new-app`
type NewAppOutput struct {
	Raw string
}

// RouteName extracts the name of the first route which was created
// The result is empty when no routes were created
func (o *NewAppOutput) RouteName() string {
	r, _ := regexp.Compile(`route "([a-z-]+)" created`)
	matches := r.FindStringSubmatch(o.Raw)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// NewAppFromTemplate applies a template using the command `oc new-app`
// It returns the output of the command and an error
// labels defines some lables which are applied to all created objects
func (c *CommandLineClient) NewAppFromTemplate(templateName string, templateParameters TemplateParameters, labels ObjectLabels) (*NewAppOutput, error) {
	if templateName == "" {
		return &NewAppOutput{}, errors.New("a template name has to be set")
	}

	arguments := []string{"new-app", "--template", templateName}
	arguments = appendKeyValueArgument(arguments, "--param", templateParameters)
	arguments = append(arguments, "--labels", labels.Combined())

	output, err := c.RunOcCommand(arguments...)
	return &NewAppOutput{Raw: output}, err
}

func appendKeyValueArgument(appendTarget []string, argumentName string, keyValuePairs map[string]string) []string {
	for key, value := range keyValuePairs {
		combinedKeyAndValue := fmt.Sprintf("%s=%s", key, value)
		appendTarget = append(appendTarget, argumentName, combinedKeyAndValue)
	}
	return appendTarget
}
