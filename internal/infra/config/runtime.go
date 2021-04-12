package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

var scope = ""

func init() {
	loadScope()
	loadViper()
}

func loadScope() {
	scope = os.Getenv("SCOPE")
	if scope == "" {
		scope = "local"
	}
	log.Info(fmt.Sprintf("Running on scope: %s", scope))
}

func IsLocal() bool {
	return scope == "local"
}

func IsProd() bool {
	return scope == "prod"
}

func IsTest() bool {
	return scope == "test"
}

func IsAWS() bool {
	return viper.GetString("cloud") == "AWS"
}

func GetScope() string {
	return scope
}