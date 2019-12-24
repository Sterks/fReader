package main

import (
	"github.com/BurntSushi/toml"
	"github.com/Sterks/FReader/config"
	"github.com/Sterks/FReader/services"
)

func main() {
	configPath := "config/config.toml"
	config := config.NewConf()
	toml.DecodeFile(configPath, &config)
	ftpreader := services.New(config)
	ftpreader.Start()
}
