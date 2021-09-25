package config

import (
	"os"

	"github.com/spf13/viper"
)

//InitConfiguration loads global configuration variables on Viper
func InitConfiguration() {
	viper.Set("cloud", os.Getenv("SERVER"))
	viper.Set("application", os.Getenv("APPLICATION"))
	viper.Set("church.name", os.Getenv("CHURCH_NAME"))
	viper.Set("church.shortname", os.Getenv("CHURCH_NAME_SHORT"))
	viper.Set("lang", os.Getenv("APP_LANG"))
	viper.Set("bundles.location", "bundles")
	viper.Set("tables.member", os.Getenv("TABLE_MEMBER"))
	viper.Set("tables.member_history", os.Getenv("TABLE_MEMBER_HISTORY"))
	viper.Set("tables.user", os.Getenv("TABLE_USER"))
	viper.Set("pdf.font.path", "./fonts/Arial.ttf")
	viper.Set("jobs.daily.phones", os.Getenv("JOBS_DAILY_PHONE"))
	viper.Set("jobs.weekly.emails", os.Getenv("JOBS_WEEKLY_EMAIL"))
	viper.Set("reports.topic", os.Getenv("REPORTS_TOPIC"))
	viper.Set("security.token.secret", os.Getenv("TOKEN_SECRET"))
	viper.Set("security.token.expiration", os.Getenv("TOKEN_EXPIRATION"))
	viper.Set("email.sender", os.Getenv("EMAIL_SENDER"))
}

func IsAWS() bool {
	return viper.GetString("cloud") == "AWS"
}
