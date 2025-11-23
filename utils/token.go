package utils

import (
	"os" // Kita butuh package OS untuk baca ENV
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID uint) (string, error) {
	// AMBIL DARI ENV
	apiSecret := os.Getenv("JWT_SECRET") 
	
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Ubah secret jadi []byte saat dipakai
	return token.SignedString([]byte(apiSecret))
}

func ApiSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}