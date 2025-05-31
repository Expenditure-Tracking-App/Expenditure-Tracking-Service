package transaction

import (
	"fmt"
	"strconv"
	"time"
)

const inputDateFormat = "02.01.06"

// Transaction represents a user's transaction data.
type Transaction struct {
	Name          string  `json:"name"`
	Amount        float32 `json:"amount"`
	Currency      string  `json:"currency"`
	Date          string  `json:"date"`
	IsClaimable   bool    `json:"is_claimable"`
	PaidForFamily bool    `json:"paid_for_family"`
	Category      string  `json:"category"`
}

// TODO: Deprecate all versions and combine to a single one
type TransactionV2 struct {
	Name          string  `db:"name"`
	Amount        float32 `db:"amount"`
	Currency      string  `db:"currency"`
	Date          string  `db:"date"`
	IsClaimable   bool    `db:"is_claimable"`
	PaidForFamily bool    `db:"paid_for_family"`
	Category      string  `db:"category"`
}

// TODO: Deprecate all versions and combine to a single one
type TransactionV3 struct {
	ID            int64     `db:"id"`
	Name          string    `db:"name"`
	Amount        float32   `db:"amount"`
	Currency      string    `db:"currency"`
	Date          string    `db:"date"`
	IsClaimable   bool      `db:"is_claimable"`
	PaidForFamily bool      `db:"paid_for_family"`
	Category      string    `db:"category"`
	CreatedAt     time.Time `db:"created_at"`
}

// ValidateAmount checks if the amount is a valid number.
func ValidateAmount(amountStr string) (float32, error) {
	f, err := strconv.ParseFloat(amountStr, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid amount: %w", err)
	}
	return float32(f), nil
}

// ValidateBool checks if the string is a valid boolean.
func ValidateBool(s string) (bool, error) {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false, fmt.Errorf("invalid boolean value: %w", err)
	}
	return b, nil
}

func ProcessDate(answer string) string {
	if answer == "t" {
		return time.Now().Format("2006-01-02")
	}

	t, err := time.Parse(inputDateFormat, answer)
	if err != nil {
		return time.Now().Format("2006-01-02")
	}

	return t.Format("2006-01-02")
}
