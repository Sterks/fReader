package config

//Config ...
type Config struct {
	ServerConnect string `toml:"server_connect"`
}

// NewConf инициализация конфигурации
func NewConf() *Config {
	return &Config{
		ServerConnect: "",
	}
}
