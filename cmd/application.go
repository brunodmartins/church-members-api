package cmd

import (
	"github.com/BrunoDM2943/church-members-api/platform/config"
	"github.com/spf13/viper"
)

// Application interface defines a single run method to be executed
type Application interface {
	//Run defines a way to start a application
	Run()
}

// ProvideRunner defines which Application runner should be initialized
func ProvideRunner() Application {
	if config.IsAWS() {
		if viper.Get("application") == "JOB" {
			return JobApplication{}
		}
		return LambdaApplication{}
	} else {
		return FiberApplication{}
	}
}
