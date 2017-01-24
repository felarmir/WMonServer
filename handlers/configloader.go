package handlers

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Login    string
	Password string
	Ip       string
	Base     string
	Port     string
}

func GetConfigData() Config {
	file, err := ioutil.ReadFile("server_cfg.yaml")
	if err != nil {
		log.Fatal("Can't load config filr. %v", err)
	}
	config := Config{}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatal("Config error. %v", err)
	}
	return config
}
