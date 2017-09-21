package openshift

// DeleteAppOutput contains the output of the DeleteApp command
type DeleteAppOutput struct {
	Raw string
}

// DeleteApp removes all objects created for the given pull request
func (c *CommandLineClient) DeleteApp(labels *ObjectLabels) (*DeleteAppOutput, error) {
	arguments := []string{"delete", "all", "--labels", labels.Combined()}

	output, err := c.RunOcCommand(arguments...)
	return &DeleteAppOutput{Raw: output}, err
}
