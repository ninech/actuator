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

func TestLoadGithubTokenFromEnvironment(t *testing.T) {
	os.Setenv(actuator.GithubTokenEnvVariable, "superyolo")

	config := actuator.Configuration{}
	config.LoadGithubToken()
	assert.Equal(t, "superyolo", config.GithubToken)

	os.Setenv(actuator.GithubTokenEnvVariable, "")
}
