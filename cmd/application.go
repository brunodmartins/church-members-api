package cmd

import (
	config2 "github.com/BrunoDM2943/church-members-api/platform/config"
)

//Application interface defines a single run method to be executed
type Application interface {
	//Run defines a way to start a application
	Run()
}

//ProvideRunner defines which Application runner should be initialized
func ProvideRunner() Application {
	if config2.IsAWS() {
		return LambdaApplication{}
	} else {
		return GinApplication{}
	}
}
