package actuator

import (
	"os"

	"github.com/fatih/color"
)

// Config contains the actual configuration of the app
var Config = Configuration{}

// GithubWebhookSecretEnvVariable is the name of the environment variable to use for the token
const GithubWebhookSecretEnvVariable = "ACTUATOR_WEBHOOK_SECRET"

// Configuration contains all configuration options for the app
type Configuration struct {
	GithubToken string
}

// LoadConfiguration loads the configuration into an internal struct
func LoadConfiguration() error {
	Config.LoadGithubWebhookSecret()
	return nil
}

// LoadGithubWebhookSecret loads the github token from an environment variable
func (c *Configuration) LoadGithubWebhookSecret() {
	c.GithubWebhookSecret = os.Getenv(GithubWebhookSecretEnvVariable)
	if c.GithubWebhookSecret == "" {
		color.Yellow("Warning: %s is not set.\n", GithubWebhookSecretEnvVariable)
	}
}
