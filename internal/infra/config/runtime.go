package config

import (
	"github.com/spf13/viper"
)

func IsAWS() bool {
	return viper.GetString("cloud") == "AWS"
}
