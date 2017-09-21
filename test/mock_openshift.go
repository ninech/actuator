package test

import (
	"fmt"

	"github.com/ninech/actuator/openshift"
)

// OpenshiftMock implements the OpenshiftClient interface to mock calls to the cli
type OpenshiftMock struct {
	AppliedTemplate   string
	AppliedLabels     openshift.ObjectLabels
	AppliedParameters openshift.TemplateParameters

	DeletedLabels openshift.ObjectLabels

	NewAppOutputToReturn    openshift.NewAppOutput
	DeleteAppOutputToReturn openshift.DeleteAppOutput
}

var ensureMockImplementsInterface openshift.OpenshiftClient = &OpenshiftMock{}

// NewAppFromTemplate mocks the apply template method and records the passed values
func (om *OpenshiftMock) NewAppFromTemplate(templateName string, templateParameters openshift.TemplateParameters, labels openshift.ObjectLabels) (*openshift.NewAppOutput, error) {
	om.AppliedTemplate = templateName
	om.AppliedLabels = labels
	om.AppliedParameters = templateParameters
	return &om.NewAppOutputToReturn, nil
}

// GetURLForRoute turns the route into an url
func (om *OpenshiftMock) GetURLForRoute(routeName string) (string, error) {
	if routeName != "" {
		return fmt.Sprintf("http://%v.domain.com", routeName), nil
	}
	return "", nil
}

// DeleteApp mocks the delte operation
func (om *OpenshiftMock) DeleteApp(labels *openshift.ObjectLabels) (*openshift.DeleteAppOutput, error) {
	om.DeletedLabels = *labels
	return &om.DeleteAppOutputToReturn, nil
}
