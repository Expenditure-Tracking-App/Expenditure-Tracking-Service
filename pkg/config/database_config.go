package config

// Define structs matching the YAML structure
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

/*type ServerConfig struct {
	Port     int    `yaml:"port"`
	LogLevel string `yaml:"logLevel"`
}

type FeaturesConfig struct {
	EnableCache bool `yaml:"enableCache"`
	MaxItems    int  `yaml:"maxItems"`
}
*/
