package main

type Transaction struct {
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	Category      int    `json:"category"`
	IsClaimable   string `json:"is_claimable"`
	date          string
	paidForFamily bool
}
