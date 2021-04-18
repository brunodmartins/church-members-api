package cmd

import (
	"github.com/BrunoDM2943/church-members-api/internal/infra/config"
)

//Application interface defines a single run method to be executed
type Application interface {
	//Run defines a way to start a application
	Run()
}

//ProvideRunner defines which Application runner should be initialized
func ProvideRunner() Application {
	if config.IsAWS() && !config.IsLocal() {
		return LambdaApplication{}
	} else {
		return GinApplication{}
	}
}
