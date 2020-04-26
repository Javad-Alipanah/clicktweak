package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Init reads config file into config struct and calls callback on file change
//
// Returns nil on success
func Init(path string, callback func()) error {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		callback()
	})

	return nil
}

// Unmarshal reads in config to specified destination
//
// Returns a err on failure
func Unmarshal(dst interface{}) error {
	return viper.Unmarshal(dst)
}
