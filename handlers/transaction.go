package handlers

import (
	"net/http"
	"time"
	"backend-gin/database"
	"backend-gin/models"
	"github.com/gin-gonic/gin"
)

// 1. GET /api/transactions (Dengan Filter)
func GetTransactions(c *gin.Context) {
	var trx []models.Transaction
	query := database.DB.Order("created_at desc")

	// Filter Tipe (income/expense)
	if tipe := c.Query("type"); tipe != "" {
		query = query.Where("type = ?", tipe)
	}

	query.Find(&trx)
	c.JSON(http.StatusOK, gin.H{"data": trx})
}

// 2. GET /api/summary
func GetSummary(c *gin.Context) {
	var income, expense int64

	// Pakai SQL Sum biar cepat & akurat
	database.DB.Model(&models.Transaction{}).Where("type = ?", "income").Select("ifnull(sum(amount),0)").Scan(&income)
	database.DB.Model(&models.Transaction{}).Where("type = ?", "expense").Select("ifnull(sum(amount),0)").Scan(&expense)

	c.JSON(http.StatusOK, gin.H{
		"total_income":  income,
		"total_expense": expense,
		"balance":       income - expense,
	})
}

// 3. GET /api/chart/daily (DATA GRAFIK)
func GetDailyChart(c *gin.Context) {
	var trx []models.Transaction
	// Ambil 30 hari terakhir
	last30Days := time.Now().AddDate(0, 0, -30)
	database.DB.Where("created_at >= ?", last30Days).Order("created_at asc").Find(&trx)

	// Format data agar mudah dibaca Chart.js
	type DailyStats struct {
		Date    string `json:"date"`
		Income  int    `json:"income"`
		Expense int    `json:"expense"`
	}
	
	// Kita gabungkan data per tanggal
	statsMap := make(map[string]*DailyStats)
	for _, t := range trx {
		dateStr := t.CreatedAt.Format("2006-01-02") // Ambil tanggalnya saja
		if _, exists := statsMap[dateStr]; !exists {
			statsMap[dateStr] = &DailyStats{Date: dateStr}
		}
		if t.Type == "income" {
			statsMap[dateStr].Income += t.Amount
		} else {
			statsMap[dateStr].Expense += t.Amount
		}
	}

	// Ubah Map jadi List
	var result []DailyStats
	for _, v := range statsMap {
		result = append(result, *v)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// 4. GET /api/categories (DATA KATEGORI)
func GetCategorySummary(c *gin.Context) {
	type CatStats struct {
		Category string `json:"category"`
		Total    int    `json:"total"`
		Type     string `json:"type"`
	}
	var results []CatStats

	// SQL Group By
	database.DB.Model(&models.Transaction{}).
		Select("category, type, sum(amount) as total").
		Group("category, type").
		Scan(&results)

	c.JSON(http.StatusOK, gin.H{"data": results})
}