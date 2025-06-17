package custom_config

import (
	"errors"
	"fmt"

	"github.com/New-Tatthep/microservice/util/stringutil"
	"github.com/spf13/viper"
)

const (
	DefaultCustomConfigFilePath = "./conf/service_conf.yaml"
)

type CustomConfig struct {
	FilePath string
	Config   *CustomConfiguration
}

var instance *CustomConfig

func New(filePath string) (*CustomConfiguration, error) {
	if stringutil.IsEmptyString(filePath) {
		return nil, errors.New("filePath is required")
	}

	config, err := read(filePath)
	if err != nil {
		return nil, err
	}

	instance = &CustomConfig{
		FilePath: filePath,
		Config:   config,
	}

	return instance.Config, nil
}

func Config() (*CustomConfiguration, error) {
	if instance == nil {
		config, err := read(DefaultCustomConfigFilePath)
		if err != nil {
			return nil, err
		}

		instance = &CustomConfig{
			FilePath: DefaultCustomConfigFilePath,
			Config:   config,
		}
	} else if stringutil.IsNotEmptyString(instance.FilePath) && instance.Config == nil {
		config, err := read(instance.FilePath)
		if err != nil {
			return nil, err
		}

		instance.Config = config
	}

	return instance.Config, nil
}

func Reload() (*CustomConfiguration, error) {
	config, err := read(instance.FilePath)
	if err != nil {
		return nil, err
	}

	instance.Config = config

	return instance.Config, nil
}

func read(filePath string) (*CustomConfiguration, error) {
	// Create a new viper instance
	viperObj := viper.New()

	// Set the config file name and path
	viperObj.SetConfigFile(filePath)

	// Read the config file
	if err := viperObj.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file fail, %s", err.Error())
	}

	// Unmarshal the config into a struct
	var config CustomConfiguration
	if err := viperObj.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unable to decode into struct, %s", err.Error())
	}

	// Validate config
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}
