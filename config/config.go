package config

//MainSettings ...
type MainSettings struct {
	Config Config
}

//Config ...
type Config struct {
	ServerConnect string `toml:"server_connect"`
	LogLevel      string `toml:"log_level"`
	Directory     Directory
	Tasks         Tasks
}

// Directory ...
type Directory struct {
	RootPath string `toml:"root_path"`
}

// Tasks ...
type Tasks struct {
	Notifications int64 `toml:"notifications"`
	Protocols     int64 `toml:"protocols"`
}

// NewConf инициализация конфигурации
func NewConf() *Config {
	return &Config{
		ServerConnect: "",
		LogLevel:      "",
	}
}
