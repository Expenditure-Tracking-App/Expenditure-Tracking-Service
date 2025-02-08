package main

type Transaction struct {
	Name          string  `json:"name"`
	Amount        float32 `json:"amount"`
	Currency      string  `json:"currency"`
	Category      int     `json:"category"`
	IsClaimable   string  `json:"is_claimable"`
	date          string
	paidForFamily bool
}
