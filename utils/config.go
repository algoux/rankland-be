package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

var conf *Config

type Config struct {
	Application struct {
		Host      string `yaml:"host"`
		Port      string `yaml:"port"`
		Cors      string `yaml:"cors-allow-origin"`
		Migration bool   `yaml:"migration"`
		Env       string `yaml:"env"`
	} `yaml:"application"`
	Log struct {
		Dir string `yaml:"dir"`
	} `yaml:"log"`
	Jwt struct {
		Secret  string `yaml:"secret"`
		Timeout int    `yaml:"timeout"`
	}
	Redis struct {
		Host       string `yaml:"host"`
		Port       string `yaml:"port"`
		Password   string `yaml:"password"`
		Database   int    `yaml:"database"`
		MaxRetries int    `yaml:"max-retry"`
	} `yaml:"redis"`
	PostgreSQL struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		DBname   string `yaml:"dbname"`
		TimeZone string `yaml:"time-zone"`
	} `yaml:"postgresql"`
}

func InitConfig() error {
	file, err := os.Open("config/config.dev.yaml")
	if err != nil {
		return err
	}

	conf = &Config{}
	err = yaml.NewDecoder(file).Decode(conf)
	if err != nil {
		return err
	}
	return nil
}

func GetConfig() *Config {
	return conf
}
