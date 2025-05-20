package config

// Top-level config struct
type Config struct {
	Database       DatabaseConfig `yaml:"database"`
	TelegramConfig TelegramConfig `yaml:"telegram"`
}

/*func GetConfig() Config {
	return Config{}
}*/
