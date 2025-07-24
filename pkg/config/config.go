package config

// Top-level config struct
type Config struct {
	FeaturesConfig      FeaturesConfig    `yaml:"features"`
	Database            DatabaseConfig    `yaml:"database"`
	TelegramConfig      TelegramConfig    `yaml:"telegram"`
	ExpenseCategories   []string          `yaml:"expense_categories"`
	FrequentExpenses    []FrequentExpense `yaml:"frequent_expenses"`
	SupportedCurrencies []string          `yaml:"supported_currencies"`
}

/*func GetConfig() Config {
	return Config{}
}*/
