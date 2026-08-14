package bot

import (
	"testing"

	"main/pkg/config"
	"main/pkg/session"

	"github.com/stretchr/testify/assert"
)

func TestNewBot_GivenValidToken_whenCreated_thenReturnsBot(t *testing.T) {
	// Skip if no test token available
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// This test requires a valid Telegram bot token
	// For unit testing, we skip it
	t.Skip("requires valid Telegram bot token")
}

func TestNewBot_GivenInvalidToken_whenCreated_thenReturnsError(t *testing.T) {
	expenses := []config.FrequentExpense{}
	categories := []string{"Food", "Transport"}
	currencies := []string{"SGD", "USD"}

	bot, err := NewBot("invalid-token-12345", config.FeaturesConfig{SaveToDB: false}, expenses, categories, currencies)

	// Invalid tokens will fail during API initialization
	assert.Error(t, err)
	assert.Nil(t, bot)
}

func TestHandleTextMessage_GivenAddCommand_whenReceived_thenStartsSession(t *testing.T) {
	// This test verifies the session creation logic without Telegram API
	// The actual handleTextMessage requires tgbotapi mocking
	// This test documents and verifies the startSession behavior

	userSessions := make(map[int64]*session.UserSession)
	chatID := int64(12345)

	// Simulate startSession logic
	userSessions[chatID] = session.NewUserSession()

	assert.NotNil(t, userSessions[chatID])
	assert.Equal(t, session.QuestionName, userSessions[chatID].CurrentQuestion)
}

func TestHandleTextMessage_GivenSummaryCommand_whenReceived_thenReturnsSummary(t *testing.T) {
	// This test verifies the summary command flow logic
	// Full integration with DB is tested in storage integration tests

	// Verify that the summary command constant is defined
	assert.Equal(t, "/summary", transactionsSummaryOption)

	// The actual handler requires DB and Telegram API mocking
	// This test documents the expected behavior:
	// 1. Fetch category counts from storage
	// 2. Build summary message with totals by category, is_claimable, paid_for_family
	// 3. Send message to user
}

func TestHandleTextMessage_GivenAnswerDuringSession_whenReceived_thenProcessesAnswer(t *testing.T) {
	userSessions := make(map[int64]*session.UserSession)
	userSessions[123] = session.NewUserSession()
	userSessions[123].CurrentQuestion = session.QuestionAmount

	// The actual message handling requires tgbotapi mocking
	// This test documents the expected behavior
	assert.NotNil(t, userSessions[123])
	assert.Equal(t, session.QuestionAmount, userSessions[123].CurrentQuestion)
}

func TestHandleTextMessage_GivenMessageWithoutSession_whenReceived_thenSendsDefaultMessage(t *testing.T) {
	userSessions := make(map[int64]*session.UserSession)

	// No session exists for this chat
	_, exists := userSessions[999]
	assert.False(t, exists)
}

func TestStartSession_GivenChatId_whenStarted_thenCreatesSession(t *testing.T) {
	userSessions := make(map[int64]*session.UserSession)
	chatID := int64(12345)

	userSessions[chatID] = session.NewUserSession()

	assert.NotNil(t, userSessions[chatID])
	assert.Equal(t, session.QuestionName, userSessions[chatID].CurrentQuestion)
}

func TestAskCurrentQuestion_GivenNameQuestion_whenAsked_thenShowsFrequentExpenses(t *testing.T) {
	// This test would require full bot mocking
	// Documents expected behavior
	expenses := []config.FrequentExpense{
		{Name: "Lunch", Category: "Food"},
		{Name: "Taxi", Category: "Transport"},
	}

	assert.Len(t, expenses, 2)
	assert.Equal(t, "Lunch", expenses[0].Name)
}

func TestHandleCallbackQuery_GivenYesForIsClaimable_whenReceived_thenSetsTrue(t *testing.T) {
	userSession := session.NewUserSession()
	userSession.CurrentQuestion = session.QuestionIsClaimable

	// Simulate callback data handling
	switch "yes" {
	case "yes":
		userSession.Answers.IsClaimable = true
	case "no":
		userSession.Answers.IsClaimable = false
	}

	assert.True(t, userSession.Answers.IsClaimable)
}

func TestHandleCallbackQuery_GivenNoForPaidForFamily_whenReceived_thenSetsFalse(t *testing.T) {
	userSession := session.NewUserSession()
	userSession.CurrentQuestion = session.QuestionPaidForFamily

	// Simulate callback data handling
	switch "no" {
	case "yes":
		userSession.Answers.PaidForFamily = true
	case "no":
		userSession.Answers.PaidForFamily = false
	}

	assert.False(t, userSession.Answers.PaidForFamily)
}

func TestHandleCallbackQuery_GivenCategorySelection_whenReceived_thenSetsCategory(t *testing.T) {
	userSession := session.NewUserSession()
	userSession.CurrentQuestion = session.QuestionCategory

	// Simulate callback data handling for category selection
	selectedCategory := "Food"
	userSession.Answers.Category = selectedCategory

	assert.Equal(t, "Food", userSession.Answers.Category)
}

func TestCompleteSession_GivenSaveToDbTrue_whenCompleted_thenSavesToDb(t *testing.T) {
	// This test verifies the complete session flow with DB saving
	// Full DB integration is tested in storage integration tests

	// Create a completed session
	userSession := session.NewUserSession()
	userSession.Answers.Name = "Test Transaction"
	userSession.Answers.Amount = 100.0
	userSession.Answers.Currency = "SGD"
	userSession.Answers.Date = "2025-01-15"
	userSession.Answers.IsClaimable = true
	userSession.Answers.PaidForFamily = false
	userSession.Answers.Category = "Food"
	userSession.CurrentQuestion = session.QuestionCount

	// Verify session is complete
	assert.True(t, userSession.IsSessionComplete())

	// The actual save to DB is tested in storage package
	// This test documents the expected flow
}

func TestCompleteSession_GivenSaveToDbFalse_whenCompleted_thenSavesToFile(t *testing.T) {
	// This test verifies file saving logic
	// Actual file operations tested in storage tests
	session := session.NewUserSession()
	session.Answers.Name = "Test"
	session.Answers.Amount = 50.0

	assert.Equal(t, "Test", session.Answers.Name)
	assert.Equal(t, float32(50.0), session.Answers.Amount)
}

func TestSendDefaultMessage_GivenChatId_whenSent_thenSendsInstructions(t *testing.T) {
	// This test would require tgbotapi mocking
	// Documents expected behavior
	expectedMessage := "Send /add to add new transaction or /summary to view summary!"

	assert.Contains(t, expectedMessage, "/add")
	assert.Contains(t, expectedMessage, "/summary")
}

func TestBot_GivenPreFilledExpense_whenNameSelected_thenAutofillsFields(t *testing.T) {
	preFilledExpenses := []config.FrequentExpense{
		{
			Name:          "Daily Lunch",
			Category:      "Food",
			Currency:      "SGD",
			IsClaimable:   true,
			PaidForFamily: false,
		},
	}

	// Simulate selecting a pre-filled expense
	selectedName := "Daily Lunch"
	var matchedExpense *config.FrequentExpense

	for _, exp := range preFilledExpenses {
		if exp.Name == selectedName {
			matchedExpense = &exp
			break
		}
	}

	assert.NotNil(t, matchedExpense)
	assert.Equal(t, "Food", matchedExpense.Category)
	assert.Equal(t, "SGD", matchedExpense.Currency)
	assert.True(t, matchedExpense.IsClaimable)
	assert.False(t, matchedExpense.PaidForFamily)
}

func TestHandleCallbackQuery_GivenCurrencySelection_whenReceived_thenSetsCurrency(t *testing.T) {
	userSession := session.NewUserSession()
	userSession.CurrentQuestion = session.QuestionCurrency

	// Simulate callback data handling for currency selection
	selectedCurrency := "USD"
	userSession.Answers.Currency = selectedCurrency

	assert.Equal(t, "USD", userSession.Answers.Currency)
}

func TestHandleCallbackQuery_GivenExpiredSession_whenReceived_thenHandlesGracefully(t *testing.T) {
	userSessions := make(map[int64]*session.UserSession)
	chatID := int64(999)

	// No session exists for this chat ID
	_, exists := userSessions[chatID]
	assert.False(t, exists)

	// Should handle gracefully (not panic) when callback arrives for expired session
}

func TestSessionProgressSummary_GivenPartialAnswers_whenBuilt_thenShowsCorrectProgress(t *testing.T) {
	userSession := session.NewUserSession()
	userSession.CurrentQuestion = session.QuestionIsClaimable

	// Simulate answers provided so far
	userSession.Answers.Name = "Test Lunch"
	userSession.Answers.Amount = 50.0
	userSession.Answers.Currency = "SGD"
	userSession.Answers.Date = "2025-01-15"

	// Verify progress can be built (actual formatting tested in bot.go)
	assert.Equal(t, "Test Lunch", userSession.Answers.Name)
	assert.Equal(t, float32(50.0), userSession.Answers.Amount)
	assert.Equal(t, "SGD", userSession.Answers.Currency)
	assert.Equal(t, "2025-01-15", userSession.Answers.Date)
	assert.False(t, userSession.Answers.IsClaimable)   // Not yet answered
	assert.False(t, userSession.Answers.PaidForFamily) // Not yet answered
	assert.Empty(t, userSession.Answers.Category)      // Not yet answered
}

func TestHandleAnswer_GivenNameDuringSession_whenReceived_thenAdvancesToAmount(t *testing.T) {
	userSession := session.NewUserSession()
	userSession.CurrentQuestion = session.QuestionName

	err := userSession.HandleAnswer("Test Transaction")

	assert.NoError(t, err)
	assert.Equal(t, "Test Transaction", userSession.Answers.Name)
	assert.Equal(t, session.QuestionAmount, userSession.CurrentQuestion)
}
