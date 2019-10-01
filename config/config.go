package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	DevEnv     = "dev"
	TestEnv    = "test"
	ReleaseEnv = "release"
)

type Config struct {
	App struct {
		Name    string
		Version string
		Build   string
		Env     string
	}

	DB struct {
		Addr       string
		User       string
		Password   string
		Name       string
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

	FileStorage struct {
		Dir string
	}

	TelegramBotAPI string `yaml:"telegram_bot_api"`

	LogsFilePath string
}

var instance *Config

// Get returns config from .config.yaml
func Get() *Config {
	if instance != nil {
		return instance
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
	err = yaml.Unmarshal(yamlConf, &instance)
	if err != nil {
		panic(err)
	}

	return instance
}

func Reload() *Config {
	instance = nil
	return Get()
}

func IsDevEnv() bool {
	return instance.App.Env == DevEnv
}

func IsTestEnv() bool {
	return instance.App.Env == TestEnv
}

func IsReleaseEnv() bool {
	return instance.App.Env == ReleaseEnv
}
