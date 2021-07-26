package config

import (
	"os"

	"github.com/spf13/viper"
)

//InitConfiguration loads global configuration variables on Viper
func InitConfiguration() {
	viper.Set("cloud", os.Getenv("SERVER"))
	viper.Set("church.name", os.Getenv("CHURCH_NAME"))
	viper.Set("lang", os.Getenv("APP_LANG"))
	viper.Set("bundles.location", "bundles")
}
