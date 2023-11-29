package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var config *Config

type Config struct {
	App *App `yaml:"app"`
	Db  *Db  `yaml:"db"`
}
type App struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
type Db struct {
	Password string `yaml:"password"`
	DbName   string `yaml:"dbname"`
	Username string `yaml:"username"`
}

func InitConfig() {
	yamlFile, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal("failed to read config.yaml")
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal(err.Error())
	}

}
