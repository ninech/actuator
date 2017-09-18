package actuator_test

import (
	"os"
	"testing"

	"github.com/ninech/actuator/actuator"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// LoadConfiguration

func TestLoadConfiguration(t *testing.T) {
	assert.Nil(t, actuator.LoadConfiguration())
}

// LoadConfigFile

func TestLoadConfigFileNotFound(t *testing.T) {
	config := actuator.Configuration{}
	assert.Nil(t, config.LoadConfigFile(afero.NewMemMapFs(), &config))
}

func TestLoadConfigFileValidFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	file, _ := fs.Create(actuator.ConfigFileName)
	file.WriteString("---\ngithub_webhook_secret: abcd")

	config := actuator.Configuration{}
	err := config.LoadConfigFile(fs, &config)

	assert.Nil(t, err)
	assert.Equal(t, "abcd", config.GithubWebhookSecret)
}

func TestLoadConfigFileInvalidYaml(t *testing.T) {
	fs := afero.NewMemMapFs()
	file, _ := fs.Create(actuator.ConfigFileName)
	file.WriteString("1234")

	config := actuator.Configuration{}
	err := config.LoadConfigFile(fs, &config)

	assert.NotNil(t, err)
}

func TestLoadConfigFileRepositories(t *testing.T) {
	fs := afero.NewMemMapFs()
	file, _ := fs.Create(actuator.ConfigFileName)
	file.WriteString(`---
repositories:
- enabled: true
  fullname: ninech/actuator
  exclude: master
`)

	config := actuator.Configuration{}
	config.LoadConfigFile(fs, &config)

	assert.Len(t, config.Repositories, 1)

	repository := config.Repositories[0]
	assert.True(t, repository.Enabled)
	assert.Equal(t, "ninech/actuator", repository.Fullname)
	assert.Equal(t, "master", repository.Exclude)
}

// LoadGithubWebhookSecret

func TestLoadGithubWebhookSecretFromEnvironment(t *testing.T) {
	os.Setenv(actuator.GithubWebhookSecretEnvVariable, "superyolo")

	config := actuator.Configuration{}
	config.LoadGithubWebhookSecret()
	assert.Equal(t, "superyolo", config.GithubWebhookSecret)

	os.Setenv(actuator.GithubWebhookSecretEnvVariable, "")
}

// GetRepositoryConfig

func TestGetRepositoryConfig(t *testing.T) {
	repoConfig := actuator.RepositoryConfig{Fullname: "ninech/actuator", Enabled: true}
	repositories := []actuator.RepositoryConfig{repoConfig}
	config := actuator.Configuration{Repositories: repositories}

	assert.Equal(t, &repoConfig, config.GetRepositoryConfig("ninech/actuator"))
}

func TestGetRepositoryConfigNotFound(t *testing.T) {
	repositories := []actuator.RepositoryConfig{}
	config := actuator.Configuration{Repositories: repositories}

	assert.Nil(t, config.GetRepositoryConfig("ninech/actuator"))
}
