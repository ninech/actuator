package test

import (
	"io/ioutil"

	"github.com/ninech/actuator/actuator"
)

// DisableLogger discards all logging
func DisableLogging() {
	actuator.Logger.SetOutput(ioutil.Discard)
}
