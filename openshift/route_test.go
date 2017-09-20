package openshift_test

import (
	"errors"
	"testing"

	"github.com/ninech/actuator/openshift"
	"github.com/ninech/actuator/testutils"
	"github.com/stretchr/testify/assert"
)

func TestGetURLForRoute(t *testing.T) {
	sampleRouteExport := `apiVersion: v1
kind: Route
metadata:
  creationTimestamp: null
  name: actuator
spec:
  host: actuator.openshift.nine.ch
  port:
    targetPort: 8080-tcp
  to:
    kind: Service
    name: actuator
    weight: 100
  wildcardPolicy: None
status: {}`

	t.Run("when the command works", func(t *testing.T) {
		openshift.CommandExecutor = &testutils.MockShell{OutputToReturn: sampleRouteExport}

		url, _ := openshift.GetURLForRoute("actuator")
		assert.Equal(t, "http://actuator.openshift.nine.ch", url)
	})

	t.Run("when there is no such route", func(t *testing.T) {
		openshift.CommandExecutor = &testutils.MockShell{ErrorToReturn: errors.New(`Error from server (NotFound): routes "actuator" not found`)}

		url, err := openshift.GetURLForRoute("actuator")
		assert.Empty(t, url)
		assert.NotNil(t, err)
	})

	t.Run("when the yaml is not valid", func(t *testing.T) {
		openshift.CommandExecutor = &testutils.MockShell{OutputToReturn: "12345"}

		url, err := openshift.GetURLForRoute("actuator")
		assert.Empty(t, url)
		assert.NotNil(t, err)
	})
}
