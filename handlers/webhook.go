package handlers

import (
	"backend-gin/database"
	"backend-gin/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type TelegramResponse struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
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

	// --- 🔒 KEAMANAN ---
	const MyTelegramID int64 = 5321617875 // Ganti dengan ID Kamu
	if chatID != MyTelegramID {
		c.JSON(http.StatusOK, gin.H{"status": "ignored_unauthorized"})
		return
	}

	// ===========================
	//   LOGIKA PERINTAH (COMMANDS)
	// ===========================

	// 1. FITUR DELETE - LANGKAH 1 (Preview)
	// Format: /del 26
	if strings.HasPrefix(text, "/del ") {
		idStr := strings.TrimPrefix(text, "/del ")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			sendReply(chatID, "⚠️ ID harus angka. Contoh: `/del 26`")
			c.JSON(http.StatusOK, gin.H{"status": "error"})
			return
		}

		var trx models.Transaction
		if err := database.DB.First(&trx, id).Error; err != nil {
			sendReply(chatID, fmt.Sprintf("❌ Data ID %d tidak ditemukan.", id))
			c.JSON(http.StatusOK, gin.H{"status": "not_found"})
			return
		}

		// Preview Data
		preview := fmt.Sprintf("⚠️ *KONFIRMASI HAPUS*\n\nID: %d\nKet: %s\nJml: %d\n\nKlik ini untuk menghapus permanen:", trx.ID, trx.Category, trx.Amount)
		
		// FIX: Format disamakan jadi /yakinhapus106 (tanpa underscore)
		confirmCmd := fmt.Sprintf("/yakinhapus%d", trx.ID) 
		
		sendReply(chatID, preview+"\n"+confirmCmd)
		c.JSON(http.StatusOK, gin.H{"status": "confirming"})
		return
	}

	// 2. FITUR DELETE - LANGKAH 2 (Eksekusi Hapus)
	// Format: /yakinhapus106 (Sesuai log chat kamu)
	if strings.HasPrefix(text, "/yakinhapus") {
		// Ambil angka dibelakang "/yakinhapus"
		idStr := strings.TrimPrefix(text, "/yakinhapus")
		id, _ := strconv.Atoi(idStr)

		// DELETE DATABASE
		result := database.DB.Delete(&models.Transaction{}, id)

		if result.Error != nil {
			sendReply(chatID, "❌ Gagal menghapus database.")
		} else if result.RowsAffected == 0 {
			sendReply(chatID, "⚠️ Data tersebut sudah tidak ada.")
		} else {
			// FIX: Memberikan pesan sukses dengan ID yang dihapus
			sendReply(chatID, fmt.Sprintf("✅ Sukses! Data ID %d berhasil dihapus permanen.", id))
		}
		
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
		return
	}

	// 3. FITUR SALDO
	if text == "/saldo" || text == "/summary" || text == "cek" {
		handleCekSaldo(chatID)
		c.JSON(http.StatusOK, gin.H{"status": "replied"})
		return
	}

	// 4. FITUR BANTUAN
	if text == "/start" || text == "/help" {
		sendHelpMessage(chatID)
		c.JSON(http.StatusOK, gin.H{"status": "replied"})
		return
	}

	// ===========================
	//   LOGIKA TRANSAKSI (+/-)
	// ===========================

	// Cek apakah pesan dimulai dengan '+' atau '-'
	isTransaction := strings.HasPrefix(text, "+") || strings.HasPrefix(text, "-")

	// JIKA BUKAN TRANSAKSI DAN BUKAN COMMAND DI ATAS
	if !isTransaction {
		pesanError := "⚠️ Perintah tidak dikenali.\nLihat panduan di bawah ini:"
		sendReply(chatID, pesanError)
		sendHelpMessage(chatID)
		c.JSON(http.StatusOK, gin.H{"status": "unknown_input"})
		return
	}

	// Jika sampai sini, berarti depannya '+' atau '-'
	parts := strings.Fields(text)
	
	if len(parts) < 2 {
		sendReply(chatID, "⚠️ Format kurang lengkap!\n\nJangan lupa kategorinya.\nContoh: `+50000 Gaji`")
		c.JSON(http.StatusOK, gin.H{"status": "format_incomplete"})
		return
	}

	nominalStr := parts[0]
	tipe := ""
	if strings.HasPrefix(nominalStr, "+") {
		tipe = "income"
	} else if strings.HasPrefix(nominalStr, "-") {
		tipe = "expense"
	}

	cleanNominal := strings.TrimPrefix(strings.TrimPrefix(nominalStr, "+"), "-")
	amount, err := strconv.Atoi(cleanNominal)
	if err != nil {
		sendReply(chatID, "⚠️ Nominal harus angka!")
		return
	}

	trx := models.Transaction{
		Amount:   amount,
		Type:     tipe,
		Category: parts[1],
		Note:     strings.Join(parts[2:], " "),
	}
	
	if err := database.DB.Create(&trx).Error; err != nil {
		sendReply(chatID, "❌ Gagal menyimpan ke database.")
		return
	}

	icon := "Dn"
	if tipe == "income" {
		icon = "UP"
	}
	
	pesan := fmt.Sprintf("✅ *Tersimpan! (ID: %d)*\n\n%s Rp %d\n📂 %s", trx.ID, icon, amount, parts[1])
	sendReply(chatID, pesan)

	c.JSON(http.StatusOK, gin.H{"status": "saved"})
}

// Helper: Menghitung Saldo
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
		"💰 *Saldo: Rp %d*", income, expense, balance)
	
	sendReply(chatID, pesan)
}

// Helper: Kirim Pesan Bantuan
func sendHelpMessage(chatID int64) {
	pesan := "🤖 *Panduan Bot Keuangan*\n\n" +
		"• `+50000 Gaji` (Pemasukan)\n" +
		"• `-20000 Makan` (Pengeluaran)\n" +
		"• `/del <ID>` (Hapus Data)\n" +
		"• /saldo (Cek Laporan)\n" +         // <— newline DI SINI WAJIB
		"• /yakinhapusid (hapus final tanpa konfirmasi)"

	sendReply(chatID, pesan)
}

// Helper: Kirim Balasan ke Telegram
func sendReply(chatID int64, text string) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	msg := TelegramResponse{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "Markdown",
	}
	body, _ := json.Marshal(msg)
	http.Post(url, "application/json", bytes.NewBuffer(body))
}