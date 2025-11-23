package models

import "time"

type Transaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Amount    int       `json:"amount"`    // Contoh: 50000
	Type      string    `json:"type"`      // "income" atau "expense"
	Category  string    `json:"category"`  // Contoh: "makan", "gaji"
	Note      string    `json:"note"`      // Catatan tambahan
	CreatedAt time.Time `json:"created_at"`
}