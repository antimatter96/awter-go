// Pacakge contants stores all the constants used all over the service

package constants

import (
	"fmt"

	"github.com/spf13/viper"
)

var configs map[string]interface{}

// Init reads the config file specified and makes it available to use by others
func Init(configFileName string, paths ...string) error {
	viper.SetConfigName(configFileName)
	for _, configPath := range paths {
		viper.AddConfigPath(configPath)
	}
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("cant read config file: %v", err)
	}

	configs = make(map[string]interface{})

	keys := viper.GetStringSlice("keys")
	for _, key := range keys {
		configs[key] = viper.Get(key)
	}

	return nil
}

//Value is used to get config values without anyone able to modify it
func Value(key string) interface{} {
	return configs[key]
}
