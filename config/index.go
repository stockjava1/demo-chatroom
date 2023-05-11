package config

import (
	"github.com/JabinGP/demo-chatroom/infra/logger"
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
			logger.Info("Read config successfully: %s", Viper.ConfigFileUsed())
		} else {
			logger.Info("Read failed: %s \n", err)
			panic(err)
		}
	})
}
