package handlers

import (
	"backend-gin/database"
	"backend-gin/models"
	"net/http"
	"strings" // <--- JANGAN LUPA TAMBAH INI
	"time"

	"github.com/gin-gonic/gin"
)

// 1. GET /api/transactions (Tetap sama)
func GetTransactions(c *gin.Context) {
	var trx []models.Transaction
	query := database.DB.Order("created_at desc")

	if tipe := c.Query("type"); tipe != "" {
		query = query.Where("type = ?", tipe)
	}

	query.Find(&trx)
	c.JSON(http.StatusOK, gin.H{"data": trx})
}

// 2. GET /api/summary (DIPERBAIKI: Pakai Hitungan Manual Biar Akurat)
func GetSummary(c *gin.Context) {
	var trx []models.Transaction
	
	// Ambil semua data transaksi
	database.DB.Find(&trx)

	var income, expense int

	// Loop satu per satu (Cara Manual = Paling Aman)
	for _, t := range trx {
		if t.Type == "income" {
			income += t.Amount
		} else if t.Type == "expense" {
			expense += t.Amount
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total_income":  income,
		"total_expense": expense,
		"balance":       income - expense,
	})
}

// 3. GET /api/chart/daily (Tetap sama)
func GetDailyChart(c *gin.Context) {
	var trx []models.Transaction
	last30Days := time.Now().AddDate(0, 0, -30)
	database.DB.Where("created_at >= ?", last30Days).Order("created_at asc").Find(&trx)

	type DailyStats struct {
		Date    string `json:"date"`
		Income  int    `json:"income"`
		Expense int    `json:"expense"`
	}
	
	statsMap := make(map[string]*DailyStats)
	for _, t := range trx {
		dateStr := t.CreatedAt.Format("2006-01-02")
		if _, exists := statsMap[dateStr]; !exists {
			statsMap[dateStr] = &DailyStats{Date: dateStr}
		}
		if t.Type == "income" {
			statsMap[dateStr].Income += t.Amount
		} else {
			statsMap[dateStr].Expense += t.Amount
		}
	}

	var result []DailyStats
	for _, v := range statsMap {
		result = append(result, *v)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// 4. GET /api/categories (DIPERBAIKI: Pakai Looping juga biar konsisten)
func GetCategorySummary(c *gin.Context) {
	var trx []models.Transaction
	database.DB.Find(&trx)

	type CatStats struct {
		Category string `json:"category"`
		Total    int    `json:"total"`
		Type     string `json:"type"`
	}

	// Pakai Map untuk mengelompokkan
	// Key string formatnya: "Tipe-Kategori" (contoh: "expense-Makan")
	tempMap := make(map[string]int)

	for _, t := range trx {
		key := t.Type + "-" + t.Category
		tempMap[key] += t.Amount
	}

	var results []CatStats
	for key, total := range tempMap {
		// Pecah key "expense-Makan" jadi Tipe & Kategori
		parts := parseKey(key) 
		results = append(results, CatStats{
			Type:     parts[0],
			Category: parts[1],
			Total:    total,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}

// Helper kecil untuk memecah key map
func parseKey(key string) []string {
	// Simple split manual tanpa import strings biar ringkas di file ini
	// Tapi sebaiknya pakai import "strings" di atas jika mau strings.Split
	// Di sini saya asumsikan kamu sudah import "strings" di atas
	return strings.Split(key, "-")
}