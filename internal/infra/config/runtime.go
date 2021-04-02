package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

func loadViper() {
	configFile := fmt.Sprintf("config_%s", scope)
	viper.SetConfigName(configFile)
	viper.SetConfigType("yaml")
	if IsTest() {
		viper.AddConfigPath("/home/runner/work/church-members-api/church-members-api/config")
	} else {
		viper.AddConfigPath("./config")
	}
	err := viper.ReadInConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"config_file": configFile,
		}).Error("Error reading config file", err)
	}
}

func loadScope() {
	scope = os.Getenv("SCOPE")
	if scope == "" {
		scope = "local"
	}
	log.Info(fmt.Sprintf("Running on scope: %s", scope))
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