package models

type BalanceResult struct {
	TotalCredits float64 `json:"total_credits"`
	TotalDebits  float64 `json:"total_debits"`
}
