package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleAnswer_GivenNameQuestion_whenValidInput_thenStoresAndAdvances(t *testing.T) {
	session := NewUserSession()
	session.CurrentQuestion = QuestionName

	err := session.HandleAnswer("Test Lunch")

	assert.NoError(t, err)
	assert.Equal(t, "Test Lunch", session.Answers.Name)
	assert.Equal(t, QuestionAmount, session.CurrentQuestion, "Should advance to next question")
}

func TestHandleAnswer_GivenAmountQuestion_whenValidInput_thenStoresAndAdvances(t *testing.T) {
	session := NewUserSession()
	session.CurrentQuestion = QuestionAmount

	err := session.HandleAnswer("50.00")

	assert.NoError(t, err)
	assert.Equal(t, float32(50.00), session.Answers.Amount)
	assert.Equal(t, QuestionCurrency, session.CurrentQuestion, "Should advance to next question")
}

func TestHandleAnswer_GivenAmountQuestion_whenInvalidInput_thenReturnsErrorAndDoesNotAdvance(t *testing.T) {
	session := NewUserSession()
	session.CurrentQuestion = QuestionAmount

	err := session.HandleAnswer("not-a-number")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid amount")
	assert.Equal(t, QuestionAmount, session.CurrentQuestion, "Should not advance on error")
}

func TestHandleAnswer_GivenCurrencyQuestion_whenValidInput_thenStoresAndAdvances(t *testing.T) {
	session := NewUserSession()
	session.CurrentQuestion = QuestionCurrency

	err := session.HandleAnswer("SGD")

	assert.NoError(t, err)
	assert.Equal(t, "SGD", session.Answers.Currency)
	assert.Equal(t, QuestionDate, session.CurrentQuestion, "Should advance to next question")
}

func TestHandleAnswer_GivenDateQuestion_whenTInput_thenStoresTodayAndAdvances(t *testing.T) {
	session := NewUserSession()
	session.CurrentQuestion = QuestionDate

	err := session.HandleAnswer("t")

	assert.NoError(t, err)
	assert.NotEmpty(t, session.Answers.Date)
	assert.Equal(t, QuestionIsClaimable, session.CurrentQuestion, "Should advance to next question")
}

func TestHandleAnswer_GivenDateQuestion_whenValidFormatInput_thenStoresAndAdvances(t *testing.T) {
	session := NewUserSession()
	session.CurrentQuestion = QuestionDate

	err := session.HandleAnswer("15.04.25")

	assert.NoError(t, err)
	assert.Equal(t, "2025-04-15", session.Answers.Date)
	assert.Equal(t, QuestionIsClaimable, session.CurrentQuestion, "Should advance to next question")
}

func TestHandleAnswer_GivenIsClaimableQuestion_whenYesInput_thenReturnsError(t *testing.T) {
	session := NewUserSession()
	session.CurrentQuestion = QuestionIsClaimable

	err := session.HandleAnswer("yes")

	assert.Error(t, err, "yes/no not accepted by strconv.ParseBool")
	assert.Contains(t, err.Error(), "invalid input")
}

func TestHandleAnswer_GivenIsClaimableQuestion_whenTrueInput_thenStoresAndAdvances(t *testing.T) {
	session := NewUserSession()
	session.CurrentQuestion = QuestionIsClaimable

	err := session.HandleAnswer("true")

	assert.NoError(t, err)
	assert.True(t, session.Answers.IsClaimable)
	assert.Equal(t, QuestionPaidForFamily, session.CurrentQuestion, "Should advance to next question")
}

func TestHandleAnswer_GivenIsClaimableQuestion_whenInvalidInput_thenReturnsErrorAndDoesNotAdvance(t *testing.T) {
	session := NewUserSession()
	session.CurrentQuestion = QuestionIsClaimable

	err := session.HandleAnswer("maybe")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid input")
	assert.Equal(t, QuestionIsClaimable, session.CurrentQuestion, "Should not advance on error")
}

func TestHandleAnswer_GivenPaidForFamilyQuestion_whenFalseInput_thenStoresAndAdvances(t *testing.T) {
	session := NewUserSession()
	session.CurrentQuestion = QuestionPaidForFamily

	err := session.HandleAnswer("false")

	assert.NoError(t, err)
	assert.False(t, session.Answers.PaidForFamily)
	assert.Equal(t, QuestionCategory, session.CurrentQuestion, "Should advance to next question")
}

func TestHandleAnswer_GivenPaidForFamilyQuestion_whenInvalidInput_thenReturnsErrorAndDoesNotAdvance(t *testing.T) {
	session := NewUserSession()
	session.CurrentQuestion = QuestionPaidForFamily

	err := session.HandleAnswer("maybe")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid input")
	assert.Equal(t, QuestionPaidForFamily, session.CurrentQuestion, "Should not advance on error")
}

func TestHandleAnswer_GivenCategoryQuestion_whenValidInput_thenStoresAndAdvances(t *testing.T) {
	session := NewUserSession()
	session.CurrentQuestion = QuestionCategory

	err := session.HandleAnswer("Food")

	assert.NoError(t, err)
	assert.Equal(t, "Food", session.Answers.Category)
	assert.Equal(t, QuestionCount, session.CurrentQuestion, "Should advance to QuestionCount (complete)")
}

func TestHandleAnswer_GivenInvalidQuestionIndex_whenHandled_thenReturnsError(t *testing.T) {
	session := NewUserSession()
	session.CurrentQuestion = 99 // Invalid question index

	err := session.HandleAnswer("any answer")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid question number")
}
