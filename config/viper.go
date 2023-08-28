package viper_config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadViperConfig() {
	viper.AddConfigPath("$HOME/.account/")
	viper.SetConfigName("viper.config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
}
