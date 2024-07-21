package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Configs struct {
	vp          *viper.Viper
	FilePath    *FilePath    `yaml:"FILEPATH"`
	Application *Application `yaml:"APPLICATION"`
}

var Settings *Configs

// init 函数在程序运行时只执行一次，
func init() {
	Settings = NewViperConfig()
}

func NewViperConfig() *Configs {
	// 读取配置
	vp := viper.New()
	viper.SetConfigName("app")  // name of config file (without extension)
	viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	viper.AddConfigPath("../")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// Config file found and successfully parsed
	configs := new(Configs)
	err = viper.Unmarshal(configs)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	configs.vp = vp
	return configs
}
