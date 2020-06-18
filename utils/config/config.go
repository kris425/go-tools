package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

//
func NewViperConfig(path string, name string) *ViperConfig {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigType("yaml")
	v.SetConfigName(name)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
		return nil
	}
	conf := &ViperConfig{
		Viper:      v,
		configPath: path,
		configName: name,
		funcs:      make([]NotifyFunc, 0),
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		for _, f := range conf.funcs {
			f(e.Name)
		}
	})
	return conf
}

type NotifyFunc func(name string)

type ViperConfig struct {
	*viper.Viper
	configPath string
	configName string
	funcs      []NotifyFunc
}

func (v *ViperConfig) AddNotifyFunc(notifyFunc NotifyFunc) {
	v.funcs = append(v.funcs, notifyFunc)
}
