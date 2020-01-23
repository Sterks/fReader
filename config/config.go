package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

//Config ...
type Config struct {
	MainSettings MainSettings
	Directory    Directory
	Tasks        Tasks
}

//MainSettings ...
type MainSettings struct {
	ServerConnect string `toml:"server_connect"`
	LogLevel      string `toml:"log_level"`
}

// Directory ...
type Directory struct {
	RootPath   string `toml:"root_path"`
	MainFolder string `toml:"main_folder"`
}

// Tasks ...
type Tasks struct {
	Notifications int64 `toml:"notifications"`
	Protocols     int64 `toml:"protocols"`
}

// NewConf инициализация конфигурации
func NewConf() *Config {
	return &Config{
		MainSettings: MainSettings{},
		Directory:    Directory{},
		Tasks:        Tasks{},
	}
}

// ConfigConfigure ...
func (conf *Config) ConfigConfigure() {
	configPath := "config/config.prod.toml"
	_, err := toml.DecodeFile(configPath, conf)
	if err != nil {
		log.Printf("Ошибка - %v", err)
	}
}
