package service

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// LoadConfig :: apply config file
func LoadConfig() error {
	// setup config
	viper.AddConfigPath("./conf")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Fatal error config file: %s", err)
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.MergeConfig(strings.NewReader(viper.GetString("configs"))); err != nil {
		log.Panic(err.Error())
	} else {
		log.Println("loaded config " + viper.GetString("app.name"))
	}
	log.Println(viper.AllSettings())
	return nil
}
