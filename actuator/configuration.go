package actuator

import (
	"os"

	"github.com/fatih/color"
)

// Config contains the actual configuration of the app
var Config = Configuration{}

// GithubTokenEnvVariable is the name of the environment variable to use for the token
const GithubTokenEnvVariable = "ACTUATOR_GITHUB_TOKEN"

// Configuration contains all configuration options for the app
type Configuration struct {
	GithubToken string
}

// LoadConfiguration loads the configuration into an internal struct
func LoadConfiguration() error {
	Config.LoadGithubToken()
	return nil
}

// LoadGithubToken loads the github token from an environment variable
func (c *Configuration) LoadGithubToken() {
	c.GithubToken = os.Getenv(GithubTokenEnvVariable)
	if c.GithubToken == "" {
		color.Yellow("Warning: %s is not set.\n", GithubTokenEnvVariable)
	}
}
