package config

// Top-level config struct
type Config struct {
	Database          DatabaseConfig    `yaml:"database"`
	TelegramConfig    TelegramConfig    `yaml:"telegram"`
	ExpenseCategories []string          `yaml:"expense_categories"`
	FrequentExpenses  []FrequentExpense `yaml:"frequent_expenses"`
}

/*func GetConfig() Config {
	return Config{}
}*/
