package config

// Define structs matching the YAML structure
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type TelegramConfig struct {
	Token string `yaml:"token"`
}

type ServerConfig struct {
	Port     int    `yaml:"port"`
	LogLevel string `yaml:"logLevel"`
}

type FeaturesConfig struct {
	EnableCache bool `yaml:"enableCache"`
	MaxItems    int  `yaml:"maxItems"`
}

// Top-level config struct
type Config struct {
	Database       DatabaseConfig `yaml:"database"`
	TelegramConfig TelegramConfig `yaml:"telegram"`
	Server         ServerConfig   `yaml:"server"`
	Features       FeaturesConfig `yaml:"features"`
}
