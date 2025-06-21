package session

import (
	"fmt"
	"main/pkg/transaction"
)

// UserSession represents a user's Q&A session.
type UserSession struct {
	CurrentQuestion int
	Answers         transaction.Transaction // Assuming this struct has Name, Amount, Category etc.
}

// Question constants
const (
	QuestionName = iota
	QuestionAmount
	QuestionCurrency
	QuestionDate
	QuestionIsClaimable
	QuestionPaidForFamily
	QuestionCategory
	QuestionCount // Should be last; represents the total number of questions

	DinnerForTheFamily         = "Dinner for the family"
	DailyTransportExpenses     = "Daily transport expenses"
	GroceriesFromPandamart     = "Groceries from Pandamart"
	MonthlyGymMembership       = "Monthly gym membership"
	GOMOMobilePlan             = "GOMO mobile plan"
	SpotifyMonthlySubscription = "Spotify monthly subscription"
	AppleICloudSubscription    = "Apple iCloud subscription"
	GoogleOneSubscription      = "Google One subscription"

	TransportCategory        = "Transport"
	FoodCategory             = "Food"
	EntertainmentCategory    = "Entertainment"
	TravelCategory           = "Travel"
	HealthAndFitnessCategory = "Health and Fitness"
	EducationCategory        = "Education"
	OtherCategory            = "Other"

	SGDCurrency = "SGD"
	USDCurrency = "USD"
	JPYCurrency = "JPY"
	CNYCurrency = "CNY"
	MYRCurrency = "MYR"
)

// Questions array for the process
var Questions = []string{
	"What is the name of the transaction?",
	"How much is the transaction?",
	"What currency is the transaction in?",
	"What is the date of transaction? (i.e., DD.MM.YY)", // Added format hint
	"Is it claimable? (yes/no)",                         // Added format hint
	"Is it paid for the family? (yes/no)",               // Added format hint
	"What is the category of transaction?",
}

// Currencies array - available for suggestions or validation
var Currencies = []string{USDCurrency, CNYCurrency, JPYCurrency, SGDCurrency, MYRCurrency}

// QuickInput array - for quick suggestions for the transaction name
var QuickInput = []string{DailyTransportExpenses, DinnerForTheFamily, GroceriesFromPandamart, MonthlyGymMembership, GOMOMobilePlan, AppleICloudSubscription, SpotifyMonthlySubscription, GoogleOneSubscription}

// TransactionCategory array - available for suggestions or validation
var TransactionCategory = []string{TransportCategory, FoodCategory, EntertainmentCategory, TravelCategory, HealthAndFitnessCategory, EducationCategory, OtherCategory}

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

// HandleAnswer processes the user's answer, updates the session,
// applies auto-fill logic, and skips already answered questions.
func (s *UserSession) HandleAnswer(answer string) error {
	var err error

	switch s.CurrentQuestion {
	case QuestionName:
		s.Answers.Name = answer
	case QuestionAmount:
		s.Answers.Amount, err = transaction.ValidateAmount(answer)
		if err != nil {
			// Return a user-friendly error message
			return fmt.Errorf("invalid amount: %w. Please enter a valid number", err)
		}

		if len(DefaultCurrency(s.Answers.Name)) > 0 {
			s.Answers.Currency = DefaultCurrency(s.Answers.Name)
			s.CurrentQuestion++
		}
	case QuestionCurrency:
		// TODO: Consider adding validation for currency (e.g., check if 'answer' is in 'Currencies' list)
		s.Answers.Currency = answer
	case QuestionDate:
		s.Answers.Date = transaction.ProcessDate(answer)
	case QuestionIsClaimable:
		s.Answers.IsClaimable, err = transaction.ValidateBool(answer)
		if err != nil {
			return fmt.Errorf("invalid input for 'claimable': %w. Please answer 'yes' or 'no'", err)
		}
	case QuestionPaidForFamily:
		s.Answers.PaidForFamily, err = transaction.ValidateBool(answer)
		if err != nil {
			return fmt.Errorf("invalid input for 'paid for family': %w. Please answer 'yes' or 'no'", err)
		}
	case QuestionCategory:
		// TODO: Consider adding validation for category (e.g., check if 'answer' is in 'TransactionCategory' list)
		s.Answers.Category = answer
	default:
		// This state should ideally not be reached if IsSessionComplete is checked before calling HandleAnswer.
		return fmt.Errorf("invalid question number: %d", s.CurrentQuestion)
	}

	// If an error occurred during the specific answer processing (e.g., validation failed),
	// return the error. s.CurrentQuestion is NOT advanced, so the same question will be asked again.
	if err != nil {
		return err
	}

	// Advance to what would normally be the next question index.
	s.CurrentQuestion++

	return nil
}

func DefaultCategory(transactionName string) string {
	switch transactionName {
	case DailyTransportExpenses:
		return TransportCategory
	case DinnerForTheFamily, GroceriesFromPandamart:
		return FoodCategory
	case MonthlyGymMembership:
		return HealthAndFitnessCategory
	case GOMOMobilePlan, AppleICloudSubscription, SpotifyMonthlySubscription, GoogleOneSubscription:
		return EntertainmentCategory
	default:
		return ""
	}
}

func DefaultPaidForFamily(transactionName string) bool {
	switch transactionName {
	case DinnerForTheFamily, GroceriesFromPandamart:
		return true
	case MonthlyGymMembership, GOMOMobilePlan, AppleICloudSubscription, SpotifyMonthlySubscription, GoogleOneSubscription:
		return false
	default:
		return false
	}
}

func DefaultCurrency(transactionName string) string {
	switch transactionName {
	case GroceriesFromPandamart, MonthlyGymMembership, GOMOMobilePlan, AppleICloudSubscription, GoogleOneSubscription:
		return SGDCurrency
	default:
		return ""
	}
}
