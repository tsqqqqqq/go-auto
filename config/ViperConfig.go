package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Setting struct {
	vp          *viper.Viper
	Application *Application `yaml:"APPLICATION"`
}

func NewViperConfig() (*Setting, error) {
	// 读取配置
	vp := viper.New()
	viper.SetConfigName("app")  // name of config file (without extension)
	viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("../")  // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		//panic(fmt.Errorf("fatal error config file: %w", err))
		return nil, err
	}

	// Config file found and successfully parsed
	setting := new(Setting)
	err = viper.Unmarshal(setting)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))

	}
	setting.vp = vp
	return setting, nil
}
