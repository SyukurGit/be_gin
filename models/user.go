package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique" json:"username"`
	Password  string    `json:"-"` // Password tidak akan dikirim balik di JSON (aman)
	CreatedAt time.Time `json:"created_at"`
}