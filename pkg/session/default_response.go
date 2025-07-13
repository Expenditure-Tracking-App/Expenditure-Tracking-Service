package session

import "main/pkg/config"

func CheckPreFilledExpense(transactionName string, preFilledExpenses []config.FrequentExpense) *config.FrequentExpense {
	for i := 0; i < len(preFilledExpenses); i++ {
		if preFilledExpenses[i].Name == transactionName {
			return &preFilledExpenses[i]
		}
	}

	return nil
}

func DefaultCategory(transactionName string, preFilledExpense *config.FrequentExpense) string {
	if preFilledExpense == nil {
		return ""
	}

	if preFilledExpense.Name == transactionName {
		return preFilledExpense.Category
	}

	return ""
}

func DefaultPaidForFamily(transactionName string, preFilledExpense *config.FrequentExpense) (bool, bool) {
	if preFilledExpense == nil {
		return false, false
	}

	if preFilledExpense.Name == transactionName {
		return preFilledExpense.PaidForFamily, true
	}

	return false, false
}

func DefaultCurrency(transactionName string, preFilledExpense *config.FrequentExpense) string {
	if preFilledExpense == nil {
		return ""
	}

	if preFilledExpense.Name == transactionName {
		return preFilledExpense.Currency
	}

	return ""
}
