package handlers

import (
	"net/http"
	"backend-gin/database"
	"backend-gin/models"
	"backend-gin/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// Cari user di database
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	// Cek Password (Hash vs Input)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	// Bikin Token
	token, _ := utils.GenerateToken(user.ID)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Fungsi sementara buat daftarin admin pertama kali
func RegisterAdmin(c *gin.Context) {
	// Hash password "admin123" (bisa diganti)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	
	user := models.User{Username: "admin", Password: string(hashedPassword)}
	database.DB.Create(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Admin berhasil dibuat!"})
}