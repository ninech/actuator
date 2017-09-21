package openshift

type OpenshiftClient interface {
	NewAppFromTemplate(string, TemplateParameters, ObjectLabels) (*NewAppOutput, error)
	DeleteApp(labels *ObjectLabels) (*DeleteAppOutput, error)
	GetURLForRoute(routeName string) (string, error)
}

// CommandLineClient represents an interface to the openshift command line
type CommandLineClient struct {
	CommandExecutor Shell
}

// NewCommandLineClient initializes a command line client with an openshift shell
func NewCommandLineClient() *CommandLineClient {
	return &CommandLineClient{CommandExecutor: &OpenshiftShell{}}
}

// RunOcCommand runs an oc command with the given arguments
func (c *CommandLineClient) RunOcCommand(args ...string) (string, error) {
	return c.CommandExecutor.RunWithArgs(args...)
}
