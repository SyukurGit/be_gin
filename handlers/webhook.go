package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"backend-gin/database"
	"backend-gin/models"
	"github.com/gin-gonic/gin"
	// "github.com/dustin/go-humanize" // Opsional: Biar angka ada titiknya (Rp 50.000)
)

// ⚠️ GANTI DENGAN TOKEN BOT KAMU YANG ASLI
const BOT_TOKEN = "8575911688:AAGMWG9TKJd2Lz6Lvf2aI7TMLQlkFIpOgmI"

type TelegramResponse struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
	ParseMode string `json:"parse_mode"` // Supaya bisa pakai huruf tebal (Bold)
}

func TelegramWebhook(c *gin.Context) {
	var payload struct {
		Message struct {
			Text string `json:"text"`
			Chat struct {
				ID int64 `json:"id"`
			} `json:"chat"`
		} `json:"message"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	text := payload.Message.Text
	chatID := payload.Message.Chat.ID

	// --- LOGIKA BARU: CEK PERINTAH ---

	// 1. Fitur Cek Saldo
	if text == "/saldo" || text == "/summary" || text == "cek" {
		handleCekSaldo(chatID)
		c.JSON(http.StatusOK, gin.H{"status": "replied"})
		return
	}

	// 2. Fitur Bantuan
	if text == "/start" || text == "/help" {
		pesan := "Halo! Saya Bot Keuangan Pribadi 🤖\n\n" +
			"Cara pakai:\n" +
			"• `+50000 Gaji` : Catat Pemasukan\n" +
			"• `-20000 Makan` : Catat Pengeluaran\n" +
			"• `/saldo` : Cek Saldo Saat Ini"
		sendReply(chatID, pesan)
		c.JSON(http.StatusOK, gin.H{"status": "replied"})
		return
	}

	// 3. Logika Transaksi (+/-) yang Lama
	parts := strings.Fields(text)
	if len(parts) < 2 {
		sendReply(chatID, "⚠️ Format salah! Ketik `/help` untuk bantuan.")
		c.JSON(http.StatusOK, gin.H{"status": "format salah"})
		return
	}

	nominalStr := parts[0]
	tipe := ""
	if strings.HasPrefix(nominalStr, "+") {
		tipe = "income"
	} else if strings.HasPrefix(nominalStr, "-") {
		tipe = "expense"
	} else {
		sendReply(chatID, "⚠️ Gunakan tanda '+' untuk pemasukan atau '-' untuk pengeluaran.")
		c.JSON(http.StatusOK, gin.H{"status": "bukan transaksi"})
		return
	}

	cleanNominal := strings.TrimPrefix(strings.TrimPrefix(nominalStr, "+"), "-")
	amount, _ := strconv.Atoi(cleanNominal)

	trx := models.Transaction{
		Amount:   amount,
		Type:     tipe,
		Category: parts[1],
		Note:     strings.Join(parts[2:], " "),
	}
	database.DB.Create(&trx)

	// Berhasil Simpan
	icon := "Dn" // Expense
	if tipe == "income" {
		icon = "UP" // Income
	}
	pesan := fmt.Sprintf("✅ *Transaksi Tersimpan!*\n\n%s Nominal: Rp %d\n📂 Ket: %s", icon, amount, parts[1])
	sendReply(chatID, pesan)

	c.JSON(http.StatusOK, gin.H{"status": "saved"})
}

// Fungsi Khusus Menghitung & Melaporkan Saldo
func handleCekSaldo(chatID int64) {
	var trx []models.Transaction
	database.DB.Find(&trx)

	var income, expense int
	for _, t := range trx {
		if t.Type == "income" {
			income += t.Amount
		} else {
			expense += t.Amount
		}
	}
	balance := income - expense

	pesan := fmt.Sprintf("📊 *Laporan Keuangan*\n\n"+
		"Total Pemasukan: Rp %d\n"+
		"Total Pengeluaran: Rp %d\n"+
		"--------------------------\n"+
		"💰 *Sisa Saldo: Rp %d*", income, expense, balance)
	
	sendReply(chatID, pesan)
}

func sendReply(chatID int64, text string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", BOT_TOKEN)
	
	msg := TelegramResponse{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "Markdown", // Biar tulisan jadi Bold/Indah
	}

	body, _ := json.Marshal(msg)
	http.Post(url, "application/json", bytes.NewBuffer(body))
}