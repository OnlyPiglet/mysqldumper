package config

import (
	"github.com/ory/viper"
)

func init() {
	//viper.SetConfigFile("./config.yaml")
	viper.SetConfigFile("./config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
