package config

import (
	"auto-record/utils"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

type Configs struct {
	vp          *viper.Viper
	Application *Application `yaml:"APPLICATION"`
	FilePath    *FilePath    `yaml:"FILEPATH"`
}

const (
	DefaultApplicationName = "go-auto"
	DefaultRecordPath      = "record-files"
	DefaultFileName        = "app.yaml"
)

var Settings *Configs

// init 函数在程序运行时只执行一次，
func init() {
	InitConfig()
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

func appFileIsExist() (string, error) {
	//app_file := path.Join()
	rootPath, err := utils.Rootname()
	if err != nil {
		return "", err
	}
	appFile := filepath.Join(rootPath, DefaultFileName)

	if _, err := os.Stat(appFile); err != nil {
		if os.IsNotExist(err) {
			return appFile, err
		}
	}
	return appFile, nil
}

// InitConfig 为了保证开发流程不受音响，档没有app.yaml的时候我来给他创建一个
func InitConfig() {
	appFile, err := appFileIsExist()
	if err == nil {
		return
	}

	application := new(Application)
	application.Name = DefaultApplicationName

	filepath2 := new(FilePath)
	filepath2.Record = DefaultRecordPath

	config := new(Configs)
	config.Application = application
	config.FilePath = filepath2

	data, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}
	fi, err := os.OpenFile(appFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer fi.Close()
	//err = os.WriteFile(appFile, data, os.ModePerm)
	if err != nil {
		panic(err)
	}
	_, err = fi.Write(data)
	if err != nil {
		panic(err)
	}

}
