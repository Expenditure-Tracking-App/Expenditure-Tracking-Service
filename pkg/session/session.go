package session

import (
	"main/pkg/transaction"
)

// UserSession represents a user's Q&A session.
type UserSession struct {
	CurrentQuestion int
	Answers         transaction.Transaction // Assuming this struct has Name, Amount, Category etc.
}

// NewUserSession creates a new user session.
func NewUserSession() *UserSession {
	return &UserSession{
		CurrentQuestion: QuestionName,
		// Answers field is implicitly initialized to its zero value
	}
}

// IsSessionComplete checks if the session is complete.
func (s *UserSession) IsSessionComplete() bool {
	return s.CurrentQuestion >= QuestionCount
}
