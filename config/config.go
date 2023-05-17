package config

import (
	"github.com/JabinGP/demo-chatroom/infra/logger"
	"github.com/spf13/viper"
	"sync"
)

var once sync.Once

// Viper viper global instance
var Viper *viper.Viper

var log *logger.CustZeroLogger

func init() {
	log = logger.NewLogger()
	log.SetModule("config")
	once.Do(func() {
		Viper = viper.New()
		// scan the file named config in the root directory
		Viper.AddConfigPath("./")
		Viper.SetConfigName("config")

		// read config, if failed, configure by default
		if err := Viper.ReadInConfig(); err == nil {
			log.SetLogLevel(Viper.GetString("loglevel.config"))
			log.Info("Read config successfully: %s", Viper.ConfigFileUsed())
		} else {
			log.Info("Read failed: %s \n", err)
			panic(err)
		}
	})
}
