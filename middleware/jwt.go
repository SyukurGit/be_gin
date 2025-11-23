package middleware

import (
	"net/http"
	"strings"
	"backend-gin/utils" // Pastikan sesuai nama module di go.mod
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari Header: "Authorization: Bearer <token>"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Butuh token akses!"})
			c.Abort()
			return
		}

		// Buang kata "Bearer " supaya sisa tokennya saja
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Cek validitas token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return utils.API_SECRET, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			c.Abort()
			return
		}

		c.Next() // Lanjut masuk ke dalam
	}
}