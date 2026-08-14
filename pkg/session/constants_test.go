package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuestionCount(t *testing.T) {
	assert.Equal(t, QuestionCount, len(Questions), "QuestionCount should equal length of Questions array")
}

func TestQuestionOrder(t *testing.T) {
	// Verify iota-based constants are sequential starting from 0
	assert.Equal(t, 0, QuestionName)
	assert.Equal(t, 1, QuestionAmount)
	assert.Equal(t, 2, QuestionCurrency)
	assert.Equal(t, 3, QuestionDate)
	assert.Equal(t, 4, QuestionIsClaimable)
	assert.Equal(t, 5, QuestionPaidForFamily)
	assert.Equal(t, 6, QuestionCategory)
	assert.Equal(t, 7, QuestionCount)
}
