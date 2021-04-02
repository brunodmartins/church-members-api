package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

var scope = ""

func BootStrapConfiguration() {
	loadScope()
	loadViper()
	if IsProd() {
		loadEnvIntoViper()
	}
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

func GetScope() string {
	return scope
}