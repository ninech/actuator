package openshift

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

type Route struct {
	Spec RouteSpec
}

type RouteSpec struct {
	Host string
}

// GetURLForRoute fetches data for a route and then extracts the domain name and returns it
func (c *CommandLineClient) GetURLForRoute(routeName string) (string, error) {
	arguments := []string{"export", "-o", "yaml", "route", routeName}
	output, err := c.RunOcCommand(arguments...)
	if err != nil {
		return "", err
	}

	routeDefinition := Route{}
	err = yaml.Unmarshal([]byte(output), &routeDefinition)
	if err != nil {
		return "", err
	}

	host := routeDefinition.Spec.Host

	return fmt.Sprintf("http://%v", host), nil
}
