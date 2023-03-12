package load

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	EnvDev  = "dev"
	EnvProd = "prod"
)

var Conf *Config

type Config struct {
	Application struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Cors string `yaml:"cors-allow-origin"`
		Env  string `yaml:"env"`
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

func Init() {
	file, err := os.Open("config/config.yaml")
	if err != nil {
		panic(err)
	}

	Conf = &Config{}
	err = yaml.NewDecoder(file).Decode(Conf)
	if err != nil {
		panic(err)
	}
}
