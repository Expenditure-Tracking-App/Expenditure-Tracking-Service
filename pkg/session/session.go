package session

import (
	"fmt"
	"main/pkg/transaction"
)

// UserSession represents a user's Q&A session.
type UserSession struct {
	CurrentQuestion int
	Answers         transaction.Transaction
}

// Question constants
const (
	QuestionName = iota
	QuestionAmount
	QuestionCurrency
	QuestionDate
	QuestionIsClaimable
	QuestionPaidForFamily
	QuestionCount // Should be last
)

// Questions array for the process
var Questions = []string{
	"What is the name of the transaction?",
	"How much is the transaction?",
	"What currency is the transaction in?",
	"What is the date of transaction?",
	"Is it claimable?",
	"Is it payable for the family?",
}

// Currencies array
var Currencies = []string{"USD", "EUR", "JPY", "SGD", "MYR"}

// NewUserSession creates a new user session.
func NewUserSession() *UserSession {
	return &UserSession{
		CurrentQuestion: QuestionName,
	}
}

// IsSessionComplete checks if the session is complete.
func (s *UserSession) IsSessionComplete() bool {
	return s.CurrentQuestion >= QuestionCount
}

// HandleAnswer processes the user's answer and updates the session.
func (s *UserSession) HandleAnswer(answer string) error {
	var err error
	switch s.CurrentQuestion {
	case QuestionName:
		s.Answers.Name = answer
	case QuestionAmount:
		s.Answers.Amount, err = transaction.ValidateAmount(answer)
		if err != nil {
			return fmt.Errorf("invalid amount: %w", err)
		}
	case QuestionCurrency:
		s.Answers.Currency = answer
	case QuestionDate:
		s.Answers.Date = answer
	case QuestionIsClaimable:
		s.Answers.IsClaimable, err = transaction.ValidateBool(answer)
		if err != nil {
			return fmt.Errorf("invalid claimable value: %w", err)
		}
	case QuestionPaidForFamily:
		s.Answers.PaidForFamily, err = transaction.ValidateBool(answer)
		if err != nil {
			return fmt.Errorf("invalid paid for family value: %w", err)
		}
	default:
		return fmt.Errorf("invalid question number: %d", s.CurrentQuestion)
	}

	s.CurrentQuestion++
	return nil
}
