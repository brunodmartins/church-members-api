package config

import (
	"github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/viper"
)

//InitConfiguration loads global configuration variables on Viper
func InitConfiguration() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	viper.Set("cloud", os.Getenv("SERVER"))
	viper.Set("application", os.Getenv("APPLICATION"))
	viper.Set("lang", os.Getenv("APP_LANG"))
	viper.Set("bundles.location", "bundles")
	viper.Set("tables.member", os.Getenv("TABLE_MEMBER"))
	viper.Set("tables.member_history", os.Getenv("TABLE_MEMBER_HISTORY"))
	viper.Set("tables.user", os.Getenv("TABLE_USER"))
	viper.Set("tables.church", os.Getenv("TABLE_CHURCH"))
	viper.Set("pdf.font.path", "./fonts/Arial.ttf")
	viper.Set("reports.topic", os.Getenv("REPORTS_TOPIC"))
	viper.Set("security.token.secret", os.Getenv("TOKEN_SECRET"))
	viper.Set("security.token.expiration", os.Getenv("TOKEN_EXPIRATION"))
	viper.Set("email.sender", os.Getenv("EMAIL_SENDER"))
	viper.Set("storage.name", os.Getenv("STORAGE"))
}

func IsAWS() bool {
	return viper.GetString("cloud") == "AWS"
}
