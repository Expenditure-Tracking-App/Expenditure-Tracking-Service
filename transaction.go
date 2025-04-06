package main

type Transaction struct {
	Name          string  `json:"name"`
	Amount        float32 `json:"amount"`
	Currency      string  `json:"currency"`
	Category      int     `json:"category"`
	Date          string  `json:"date"`
	IsClaimable   bool    `json:"is_claimable"`
	PaidForFamily bool    `json:"paid_for_family"`
}
