package config

import "github.com/spf13/viper"

func loadEnvIntoViper() {
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
