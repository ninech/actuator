package testutils

import "github.com/ninech/actuator/openshift"

type OpenshiftMock struct {
	AppliedTemplate   string
	AppliedLabels     openshift.ObjectLabels
	AppliedParameters openshift.TemplateParameters
}

// ApplyOpenshiftTemplate mocks the apply template method and records the passed values
func (om *OpenshiftMock) ApplyOpenshiftTemplate(templateName string, templateParameters openshift.TemplateParameters, labels openshift.ObjectLabels) (string, error) {
	om.AppliedTemplate = templateName
	om.AppliedLabels = labels
	om.AppliedParameters = templateParameters
	return "", nil
}
