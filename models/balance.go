package models

// Estructura para almacenar el resultado de la consulta de balance
type BalanceResult struct {
	TotalCredits float64 `json:"total_credits"`
	TotalDebits  float64 `json:"total_debits"`
}
