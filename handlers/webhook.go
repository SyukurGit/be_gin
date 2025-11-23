package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"backend-gin/database"
	"backend-gin/models"
	"github.com/gin-gonic/gin"
)

func TelegramWebhook(c *gin.Context) {
	var payload struct {
		Message struct {
			Text string `json:"text"`
		} `json:"message"`
	}

	// Telegram ngirim JSON, kita tangkap
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "ignored"}) // Tetap 200 biar gak error di Telegram
		return
	}

	text := payload.Message.Text
	// Format: "+50000 Gaji Bulanan" atau "-20000 Makan Siang"
	parts := strings.Fields(text) // Pisahkan spasi

	if len(parts) < 2 {
		c.JSON(http.StatusOK, gin.H{"status": "format salah"})
		return
	}

	// Analisis kata pertama (Nominal & Tipe)
	nominalStr := parts[0]
	tipe := ""
	if strings.HasPrefix(nominalStr, "+") {
		tipe = "income"
	} else if strings.HasPrefix(nominalStr, "-") {
		tipe = "expense"
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "bukan transaksi"})
		return
	}

	// Bersihkan simbol + dan - lalu ubah ke angka
	cleanNominal := strings.TrimPrefix(strings.TrimPrefix(nominalStr, "+"), "-")
	amount, _ := strconv.Atoi(cleanNominal)

	// Simpan Transaksi
	trx := models.Transaction{
		Amount:   amount,
		Type:     tipe,
		Category: parts[1],              // Kata kedua (contoh: Makan)
		Note:     strings.Join(parts[2:], " "), // Sisanya adalah catatan
	}
	database.DB.Create(&trx)

	c.JSON(http.StatusOK, gin.H{"status": "saved"})
}