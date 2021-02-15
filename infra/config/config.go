package config

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var scope = ""

func init() {
	loadScope()
	loadViper()
	if IsProd() {
		wrapEnvVariables()
	}

}

func wrapEnvVariables() {
	viper.SetEnvPrefix("VPR")
	viper.BindEnv("CHURCH_MEMBERS_DATABASE_URL")
	viper.BindEnv("CHURCH_NAME")

	viper.BindEnv("AUTH_ENABLE")
	viper.BindEnv("AUTH_TOKEN")
	viper.BindEnv("AUTH_JWK")
	viper.BindEnv("AUTH_ISS")
	viper.BindEnv("AUTH_AUD")

	viper.Set("mongo.url", viper.GetString("CHURCH_MEMBERS_DATABASE_URL"))

	viper.Set("church.name", viper.GetString("CHURCH_NAME"))

	//viper Configs
	viper.Set("auth.enable", viper.GetBool("AUTH_ENABLE"))
	viper.Set("auth.token", viper.GetString("AUTH_TOKEN"))
	viper.Set("auth.jwk", viper.GetString("AUTH_JWK"))
	viper.Set("auth.iss", viper.GetString("AUTH_ISS"))
	viper.Set("auth.aud", viper.GetString("AUTH_AUD"))

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
