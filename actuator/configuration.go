package actuator

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/afero"
	yaml "gopkg.in/yaml.v2"
)

// Config contains the actual configuration of the app
var Config = Configuration{}

// GithubWebhookSecretEnvVariable is the name of the environment variable to use for the token
const GithubWebhookSecretEnvVariable = "ACTUATOR_WEBHOOK_SECRET"

// ConfigFileName defines the name of the configuration file
const ConfigFileName = "actuator.yml"

// Configuration contains all configuration options for the app
type Configuration struct {
	GithubWebhookSecret string             `yaml:"github_webhook_secret"`
	Repositories        []RepositoryConfig `yaml:"repositories"`
}

// RepositoryConfig contains configuration for each repository for which events
// are being received.
type RepositoryConfig struct {
	Enabled  bool   `yaml:"enabled,omitempty"` // Run actuator for this repo?
	Fullname string `yaml:"fullname"`          // Full name of the repository. ex. ninech/actuator
	Exclude  string `yaml:"exclude,omitempty"` // Pattern to exclude branches. ex. ^master$
}

// LoadConfiguration loads the configuration into an internal struct
func LoadConfiguration() error {
	err := Config.LoadConfigFile(afero.NewOsFs(), &Config)
	if err != nil {
		return err
	}

	if Config.GithubWebhookSecret == "" {
		Config.LoadGithubWebhookSecret()
	}

	return nil
}

// LoadGithubWebhookSecret loads the github token from an environment variable
func (c *Configuration) LoadGithubWebhookSecret() {
	c.GithubWebhookSecret = os.Getenv(GithubWebhookSecretEnvVariable)
	if c.GithubWebhookSecret == "" {
		color.Yellow("Warning: %s is not set.\n", GithubWebhookSecretEnvVariable)
	}
}

// LoadConfigFile loads the configuration file into the Config variable
// It takes an afero.Fs struct to abstract the file system
func (c *Configuration) LoadConfigFile(fs afero.Fs, config *Configuration) error {
	if _, err := fs.Stat(ConfigFileName); os.IsNotExist(err) {
		color.Yellow("The file %s does not exist. Using default config.", ConfigFileName)
		return nil
	}

	yamlFile, err := afero.ReadFile(fs, ConfigFileName)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return err
	}

	return nil
}