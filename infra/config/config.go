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
	viper.BindEnv("CHURCH_MEMBERS_ACCESS_TOKEN")

	viper.Set("mongo.url", viper.GetString("CHURCH_MEMBERS_DATABASE_URL"))
	viper.Set("auth.token", viper.GetString("CHURCH_MEMBERS_ACCESS_TOKEN"))
}

func loadViper() {
	configFile := fmt.Sprintf("config_%s", scope)
	viper.SetConfigName(configFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

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

func GetScope() string {
	return scope
}
