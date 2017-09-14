package actuator_test

import (
	"os"
	"testing"

	"github.com/ninech/actuator/actuator"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfiguration(t *testing.T) {
	assert.Nil(t, actuator.LoadConfiguration())
}

func TestLoadGithubWebhookSecretFromEnvironment(t *testing.T) {
	os.Setenv(actuator.GithubWebhookSecretEnvVariable, "superyolo")

	config := actuator.Configuration{}
	config.LoadGithubWebhookSecret()
	assert.Equal(t, "superyolo", config.GithubWebhookSecret)

	os.Setenv(actuator.GithubWebhookSecretEnvVariable, "")
}
