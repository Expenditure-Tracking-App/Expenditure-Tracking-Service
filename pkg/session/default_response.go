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
