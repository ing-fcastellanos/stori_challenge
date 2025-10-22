package models

import "time"

type Transaction struct {
	ID        int       `json:"id"`
	LegacyId  int       `json:"legacy_id"`
	UserID    int       `json:"user_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"datetime"`
}
