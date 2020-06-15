package config

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"tomm/utils"
)

const (
	CONFIGDIR = "configfile"
)

func init() {
	newFile("")
}

func newFile(base string) {
	if base == "" {
		base = utils.GetProDirAbs() + CONFIGDIR
	}

	viper.SetConfigName("config_test")
	viper.SetConfigType("yaml")
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
