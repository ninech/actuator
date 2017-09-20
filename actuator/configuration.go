package actuator

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/afero"
	yaml "gopkg.in/yaml.v2"
)

// Config contains the actual configuration of the app
var Config = Configuration{}

const (
	// GithubWebhookSecretEnvVariable is the name of the environment variable to use for the secret
	GithubWebhookSecretEnvVariable = "ACTUATOR_WEBHOOK_SECRET"
	// GithubAccessTokenEnvVariable is the name of the environment variable to use for the access token
	GithubAccessTokenEnvVariable = "ACTUATOR_GITHUB_ACCESS_TOKEN"
)

// ConfigFileName defines the name of the configuration file
const ConfigFileName = "actuator.yml"

// Configuration contains all configuration options for the app
type Configuration struct {
	GithubWebhookSecret string             `yaml:"github_webhook_secret"` // used to identify incomming webhooks
	GithubAccessToken   string             `yaml:"github_access_token"`   // used to identify against the github api
	Repositories        []RepositoryConfig `yaml:"repositories"`          // repository configurations
}

// RepositoryConfig contains configuration for each repository for which events
// are being received.
type RepositoryConfig struct {
	Enabled  bool   `yaml:"enabled"`  // Run actuator for this repo?
	Fullname string `yaml:"fullname"` // Full name of the repository. ex. ninech/actuator
	Exclude  string `yaml:"exclude"`  // Pattern to exclude branches. ex. ^master$
	Template string `yaml:"template"` // Defines the openshift template to apply
}

// LoadConfiguration loads the configuration into an internal struct
func LoadConfiguration() error {
	err := Config.LoadConfigFile(afero.NewOsFs(), &Config)
	if err != nil {
		return err
	}

	Config.LoadGithubWebhookSecret()
	Config.LoadGithubAccessToken()

	return nil
}

// LoadGithubWebhookSecret loads the github token from an environment variable
func (c *Configuration) LoadGithubWebhookSecret() {
	secretFromEnv := os.Getenv(GithubWebhookSecretEnvVariable)
	if secretFromEnv != "" {
		c.GithubWebhookSecret = secretFromEnv
	}
	if c.GithubWebhookSecret == "" {
		color.Yellow("Warning: %s is not set.\n", GithubWebhookSecretEnvVariable)
	}
}

// LoadGithubAccessToken loads the github token from an environment variable
func (c *Configuration) LoadGithubAccessToken() {
	tokenFromEnv := os.Getenv(GithubAccessTokenEnvVariable)
	if tokenFromEnv != "" {
		c.GithubAccessToken = tokenFromEnv
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

// GetRepositoryConfig returns the configuration for the repository with the given name
func (c *Configuration) GetRepositoryConfig(fullname string) *RepositoryConfig {
	for _, repo := range c.Repositories {
		if repo.Fullname == fullname {
			return &repo
		}
	}
	return nil
}
