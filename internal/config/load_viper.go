package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func LoadViper(env string) *viper.Viper {

	v := viper.New()
	v.SetConfigName(fmt.Sprintf("application-%s", env))
	v.AddConfigPath(".")
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, err{%v}", err)
	}

	return v
}
