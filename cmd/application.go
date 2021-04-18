package cmd

import (
	"github.com/BrunoDM2943/church-members-api/internal/infra/config"
)

type Application interface {
	Run()
}

func ProvideRunner() Application {
	if config.IsAWS() && !config.IsLocal() {
		return LambdaApplication{}
	} else {
		return GinApplication{}
	}
}