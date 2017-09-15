package openshift

// Cli exposes the oc utility
var Cli = CommandLineClient{}

// Client is interface for every openshift cli client
type Client interface {
	NewApp(templateName string, params TemplateParameters) error
}

// CommandLineClient wrapps the openshift cli (`oc`)
type CommandLineClient struct{}

// NewApp applies a template using the command `oc new-app`
func (c *CommandLineClient) NewApp(templateName string, params TemplateParameters) error {
	return nil
}
