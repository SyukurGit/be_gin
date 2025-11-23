package utils

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

// Ganti string ini dengan kunci rahasiamu sendiri nanti
var API_SECRET = []byte("rahasia_dapur_syukur_123")

func GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token berlaku 24 jam

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(API_SECRET)
}