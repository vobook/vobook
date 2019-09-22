package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	AppEnvDev     = "dev"
	AppEnvTest    = "test"
	AppEnvRelease = "release"
)

type Config struct {
	App struct {
		Name    string
		Version string
		Build   string
		Env     string
	}

	DB struct {
		User       string
		Password   string
		Name       string
		Port       string
		Schema     string
		LogQueries bool `yaml:"log_queries"`
	}

	Mail struct {
		Server   string
		Port     string
		User     string
		Password string
		Driver   string
	}

	Server struct {
		Port string
		Host string
	}

	TelegramBotAPI string `yaml:"telegram_bot_api"`

	LogsFilePath string
}

var Instance *Config

// Get returns config from .config.yaml
func Get() *Config {
	if Instance != nil {
		return Instance
	}

	name := ".config.yaml"
	yamlConf, err := ioutil.ReadFile(name)
	if err != nil {
		msg := ".config.yaml missing from working directory"
		wd, err := os.Getwd()
		if err == nil {
			msg += fmt.Sprintf(" (%s)", wd)
		}
		panic(msg)
	}
	err = yaml.Unmarshal(yamlConf, &Instance)
	if err != nil {
		panic(err)
	}

	return Instance
}

func Reload() *Config {
	Instance = nil
	return Get()
}

func IsDevEnv() bool {
	return Instance.App.Env == AppEnvDev
}

func IsTestEnv() bool {
	return Instance.App.Env == AppEnvTest
}

func IsReleaseEnv() bool {
	return Instance.App.Env == AppEnvRelease
}