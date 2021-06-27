package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

type Config struct {
	Addr     string   `toml:"addr"`
	Secret   string   `toml:"secret"`
	Issure   string   `toml:"issure"`
}

var config *Config

func InitConfig(path string) {
	config = &Config{}
	if _, err := toml.DecodeFile(path, config); err != nil {
		log.Println(err)
		log.Fatalln("cannot load config")
	}
}

func Get() *Config {
	return config
}

