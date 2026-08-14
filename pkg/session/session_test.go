package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserSession(t *testing.T) {
	session := NewUserSession()

	assert.Equal(t, QuestionName, session.CurrentQuestion, "New session should start with QuestionName")
	assert.Empty(t, session.Answers.Name, "Answers.Name should be zero value")
	assert.Equal(t, float32(0), session.Answers.Amount, "Answers.Amount should be zero value")
	assert.Empty(t, session.Answers.Currency, "Answers.Currency should be zero value")
	assert.Empty(t, session.Answers.Date, "Answers.Date should be zero value")
	assert.False(t, session.Answers.IsClaimable, "Answers.IsClaimable should be false")
	assert.False(t, session.Answers.PaidForFamily, "Answers.PaidForFamily should be false")
	assert.Empty(t, session.Answers.Category, "Answers.Category should be zero value")
	assert.Equal(t, 0, session.LastQuestionMessageID, "LastQuestionMessageID should be zero")
}

func TestIsSessionComplete(t *testing.T) {
	tests := []struct {
		name             string
		currentQuestion  int
		expectedComplete bool
	}{
		{"given question less than count when checked then returns false", QuestionName, false},
		{"given question at count when checked then returns true", QuestionCount, true},
		{"given question beyond count when checked then returns true", QuestionCount + 1, true},
		{"given last question answered when checked then returns true", QuestionCategory + 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &UserSession{
				CurrentQuestion: tt.currentQuestion,
			}
			result := session.IsSessionComplete()
			assert.Equal(t, tt.expectedComplete, result)
		})
	}
}
