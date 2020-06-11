package config

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func init() {
	newFile("")
}

func newFile(base string) {
	if base == "" {
		base = "../configfile"
	}
	viper.SetConfigType("yaml")
	viper.SetConfigName("config_test")
	viper.AddConfigPath(base)
	loadFile()
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		switch in.Op {
		case fsnotify.Create, fsnotify.Write:
			loadFile()
		}
	})
}

func loadFile() {
	err := viper.ReadInConfig()
	if err != nil {
		panic("ReadConfig File " + err.Error())
	}
}

func SetFile(fileName string) {
	viper.SetConfigFile(fileName)
	loadFile()
}

func Decode(key string, data interface{}) error {
	if !viper.IsSet(key) {
		return errors.New("cur key not exist in any Config File")
	}
	return viper.UnmarshalKey(key, data)
}
