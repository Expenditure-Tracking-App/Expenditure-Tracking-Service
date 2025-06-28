package session

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

func DefaultPaidForFamily(transactionName string) (bool, bool) {
	switch transactionName {
	case DinnerForTheFamily, GroceriesFromPandamart:
		return true, true
	case MonthlyGymMembership, GOMOMobilePlan, AppleICloudSubscription, SpotifyMonthlySubscription, GoogleOneSubscription:
		return false, true
	default:
		return false, false
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
