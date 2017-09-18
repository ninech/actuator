package openshift

import (
	"fmt"
	"strings"
)

// ObjectLabels is a hash of labels for openshift objects
type ObjectLabels map[string]string

// Combined combines the labels into one string
// example: label1=value1,label2=value2
// This is used when passing the labels into an `oc` command
func (o ObjectLabels) Combined() string {
	labelList := []string{}
	for labelName, labelValue := range o {
		labelList = append(labelList, fmt.Sprintf("%s=%s", labelName, labelValue))
	}
	return strings.Join(labelList, ",")
}
