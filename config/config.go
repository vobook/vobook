package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v3"
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
		Host      string
		Port      int
		User      string
		Password  string
		Driver    string
		From      string
		Stub      string
		Templates string
	}

	Server struct {
		Port string
		Host string
	}

	FileStorage struct {
		Dir string
	}

	TelegramBotAPI string `yaml:"telegram_bot_api"`

	LogsFilePath              string
	EmailVerificationLifetime time.Duration `yaml:"email_verification_lifetime"`
	AuthTokenLifetime         time.Duration `yaml:"auth_token_lifetime"`
	PasswordResetLifetime     time.Duration `yaml:"password_reset_lifetime"`

	ApiBasePath   string `yaml:"api_base_path"`
	WebClientAddr string `yaml:"web_client_addr"`

	DateFormat     string `yaml:"date_format"`
	TimeFormat     string `yaml:"time_format"`
	DateTimeFormat string `yaml:"date_time_format"`
}

var conf *Config

// Get returns config from .config.yaml
func Get() *Config {
	if conf != nil {
		return conf
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
	err = yaml.Unmarshal(yamlConf, &conf)
	if err != nil {
		panic(err)
	}

	if conf.WebClientAddr == "" {
		conf.WebClientAddr = conf.Server.Host
		if conf.Server.Port != "" {
			conf.WebClientAddr += ":" + conf.Server.Port
		}
	}

	if conf.DateFormat == "" {
		conf.DateFormat = "2006-01-02"
	}
	if conf.TimeFormat == "" {
		conf.TimeFormat = "15:04:05"
	}
	if conf.DateTimeFormat == "" {
		conf.DateTimeFormat = conf.DateFormat + " " + conf.DateTimeFormat
	}

	return conf
}

func Reload() *Config {
	conf = nil
	return Get()
}

func IsDevEnv() bool {
	return conf.App.Env == DevEnv
}

func IsTestEnv() bool {
	return conf.App.Env == TestEnv
}

func IsReleaseEnv() bool {
	return conf.App.Env == ReleaseEnv
}
