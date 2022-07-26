package utils

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var conf *Config

type Config struct {
	Application struct {
		Host      string `yaml:"host"`
		Port      string `yaml:"port"`
		Cors      string `yaml:"cors-allow-origin"`
		Migration bool   `yaml:"migration"`
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
	MySQL struct {
		Username  string `yaml:"username"`
		Password  string `yaml:"password"`
		Host      string `yaml:"host"`
		Port      string `yaml:"port"`
		Database  string `yaml:"database"`
		Charset   string `yaml:"charset"`
		ParseTime bool   `yaml:"parse-time"`
	} `yaml:"mysql"`
}

func init() {
	file, err := os.Open("config/config.yaml")
	if err != nil {
		logrus.Fatalf("config.yml file not found in config directory, error=%v", err)
		panic(err)
	}

	conf = &Config{}
	err = yaml.NewDecoder(file).Decode(conf)
	if err != nil {
		logrus.Fatalf("config.yml file incorrect format cannot be resolved, error=%v", err)
		panic(err)
	}
}

func GetConfig() *Config {
	return conf
}
