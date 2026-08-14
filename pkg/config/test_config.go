package config

// TestConfig returns a configuration suitable for tests
func TestConfig() Config {
	return Config{
		Server: ServerConfig{Port: 0}, // Use random port
		FeaturesConfig: FeaturesConfig{
			SaveToDB: false, // Use file storage for unit tests
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "test",
			Password: "test",
			DBName:   "expenditure_test",
			SSLMode:  "disable",
		},
		TelegramConfig:      TelegramConfig{Token: "test-token"},
		ExpenseCategories:   []string{"Food", "Transport", "Utilities"},
		FrequentExpenses:    []FrequentExpense{},
		SupportedCurrencies: []string{"USD", "EUR", "SGD"},
	}
}
