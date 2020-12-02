package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

//Config ...
type Config struct {
	MainSettings MainSettings
	Directory    Directory
	Tasks        Tasks
	Rabbit       Rabbit
	FTPServer44  FTPServer44
	FTPServer223 FTPServer223
	Postgres     Postgres
	TimeDownloader TimeDownloader
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

// TimeDownloader ...
type TimeDownloader struct {
	From string `toml:"from"`
}

// Tasks ...
type Tasks struct {
	Notifications44  int64 `toml:"notifications44"`
	Notifications223 int64 `toml:"notifications223"`
	Protocols44      int64 `toml:"protocols44"`
	Protocols223     int64 `toml:"protocols223"`
}

// Rabbit соединение для RabbitMQ
type Rabbit struct {
	ConnectRabbit string `toml:"connection_rabbit"`
}

// FTPServer44 Настройки для FTP 44
type FTPServer44 struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	Url44    string `toml:"url44"`
	RootPath string `toml:"root_path"`
}

type FTPServer223 struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	Url223   string `toml:"url44"`
	RootPath string `toml:"root_path"`
}

// Postgres настройки для подключения к базе данных
type Postgres struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"dbname"`
}

// NewConf инициализация конфигурации
func NewConf() *Config {
	return &Config{
		MainSettings: MainSettings{},
		Directory:    Directory{},
		Tasks:        Tasks{},
		Postgres:     Postgres{},
		FTPServer223: FTPServer223{},
		FTPServer44:  FTPServer44{},
		Rabbit:       Rabbit{},
		TimeDownloader: TimeDownloader{},
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
