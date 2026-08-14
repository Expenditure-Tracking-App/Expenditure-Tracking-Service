package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestConfigStruct_GivenValidYaml_whenParsed_thenLoadsCorrectly(t *testing.T) {
	yamlContent := `
server:
  port: 8081
  logLevel: info
features:
  save_to_database: true
database:
  host: localhost
  port: 5432
  user: testuser
  password: testpass
  dbname: expenditure
  ssl_mode: disable
telegram:
  token: test-bot-token
expense_categories:
  - Food
  - Transport
  - Utilities
frequent_expenses:
  - name: Daily Lunch
    category: Food
    currency: SGD
    is_claimable: true
    paid_for_family: false
supported_currencies:
  - SGD
  - USD
  - EUR
`

	var cfg Config
	err := yaml.Unmarshal([]byte(yamlContent), &cfg)

	require.NoError(t, err)
	assert.Equal(t, 8081, cfg.Server.Port)
	assert.Equal(t, "info", cfg.Server.LogLevel)
	assert.True(t, cfg.FeaturesConfig.SaveToDB)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.Equal(t, "testuser", cfg.Database.User)
	assert.Equal(t, "test-bot-token", cfg.TelegramConfig.Token)
	assert.Len(t, cfg.ExpenseCategories, 3)
	assert.Contains(t, cfg.ExpenseCategories, "Food")
}

func TestConfigStruct_GivenMinimalYaml_whenParsed_thenUsesDefaults(t *testing.T) {
	yamlContent := `
server:
  port: 8081
features:
  save_to_database: false
telegram:
  token: minimal-token
expense_categories: []
supported_currencies: []
frequent_expenses: []
`

	var cfg Config
	err := yaml.Unmarshal([]byte(yamlContent), &cfg)

	require.NoError(t, err)
	assert.Equal(t, 8081, cfg.Server.Port)
	assert.False(t, cfg.FeaturesConfig.SaveToDB)
	assert.Empty(t, cfg.ExpenseCategories)
	assert.Empty(t, cfg.SupportedCurrencies)
	assert.Empty(t, cfg.FrequentExpenses)
}

func TestFrequentExpense_GivenYamlWithAllFields_whenParsed_thenLoadsCorrectly(t *testing.T) {
	yamlContent := `
name: Weekly Taxi
category: Transport
currency: SGD
is_claimable: true
paid_for_family: false
`

	var expense FrequentExpense
	err := yaml.Unmarshal([]byte(yamlContent), &expense)

	require.NoError(t, err)
	assert.Equal(t, "Weekly Taxi", expense.Name)
	assert.Equal(t, "Transport", expense.Category)
	assert.Equal(t, "SGD", expense.Currency)
	assert.True(t, expense.IsClaimable)
	assert.False(t, expense.PaidForFamily)
}

func TestDatabaseConfig_GivenYaml_whenParsed_thenLoadsCorrectly(t *testing.T) {
	yamlContent := `
host: db.example.com
port: 5433
user: appuser
password: secretpass
dbname: prod_expenses
ssl_mode: require
`

	var dbCfg DatabaseConfig
	err := yaml.Unmarshal([]byte(yamlContent), &dbCfg)

	require.NoError(t, err)
	assert.Equal(t, "db.example.com", dbCfg.Host)
	assert.Equal(t, 5433, dbCfg.Port)
	assert.Equal(t, "require", dbCfg.SSLMode)
}

func TestFeaturesConfig_GivenSaveToDatabaseTrue_whenParsed_thenEnablesDbSaving(t *testing.T) {
	yamlContent := `save_to_database: true`

	var features FeaturesConfig
	err := yaml.Unmarshal([]byte(yamlContent), &features)

	require.NoError(t, err)
	assert.True(t, features.SaveToDB)
}

func TestServerConfig_GivenYaml_whenParsed_thenLoadsCorrectly(t *testing.T) {
	yamlContent := `
port: 9000
logLevel: debug
`

	var server ServerConfig
	err := yaml.Unmarshal([]byte(yamlContent), &server)

	require.NoError(t, err)
	assert.Equal(t, 9000, server.Port)
	assert.Equal(t, "debug", server.LogLevel)
}

func TestTestConfig_GivenNoArgs_whenCalled_thenReturnsTestConfiguration(t *testing.T) {
	cfg := TestConfig()

	assert.Equal(t, 0, cfg.Server.Port) // Random port
	assert.False(t, cfg.FeaturesConfig.SaveToDB)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.Equal(t, "test", cfg.Database.User)
	assert.Equal(t, "test", cfg.Database.Password)
	assert.Equal(t, "expenditure_test", cfg.Database.DBName)
	assert.Equal(t, "test-token", cfg.TelegramConfig.Token)
	assert.Contains(t, cfg.ExpenseCategories, "Food")
	assert.Contains(t, cfg.SupportedCurrencies, "SGD")
}

func TestConfigStruct_GivenInvalidYaml_whenParsed_thenReturnsError(t *testing.T) {
	// Test that invalid YAML values are caught during unmarshalling
	yamlContent := `
server:
  port: "not_a_number"
`
	var cfg Config
	err := yaml.Unmarshal([]byte(yamlContent), &cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot unmarshal")
}
