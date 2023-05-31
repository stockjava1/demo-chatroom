package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"sync"
)

var once sync.Once

// Viper viper global instance
var Viper *viper.Viper

func init() {
	once.Do(func() {
		Viper = viper.New()
		// scan the file named config in the root directory
		Viper.AddConfigPath("./")
		Viper.SetConfigName("config")

		// read config, if failed, configure by default
		if err := Viper.ReadInConfig(); err == nil {
			log.Info().Msgf("Read config successfully: %s", Viper.ConfigFileUsed())
		} else {
			log.Info().Msgf("Read failed: %s \n", err)
			panic(err)
		}
	})
}
