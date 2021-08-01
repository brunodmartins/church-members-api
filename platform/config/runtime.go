package config

import (
	"github.com/spf13/viper"
	"os"
)

//InitConfiguration loads global configuration variables on Viper
func InitConfiguration() {
	viper.Set("cloud", os.Getenv("SERVER"))
	viper.Set("church.name", os.Getenv("CHURCH_NAME"))
	viper.Set("lang", os.Getenv("APP_LANG"))
	viper.Set("bundles.location", "bundles")
	viper.Set("tables.member", "member")
	viper.Set("tables.member_history", "member_history")
	viper.Set("pdf.font.path", "./fonts/Arial.ttf")
}

func IsAWS() bool {
	return viper.GetString("cloud") == "AWS"
}
