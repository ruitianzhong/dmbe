package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	App *App `yaml:"app"`
	Db  *Db  `yaml:"db"`
}
type App struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
type Db struct {
	Password string `yaml:"password"`
	DbName   string `yaml:"dbname"`
	Username string `yaml:"username"`
	Port     string `yaml:"port"`
	Address  string `yaml:"address"`
}

func InitConfig() *Config {
	var config *Config
	yamlFile, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal("failed to read config.yaml")
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	return config
}
