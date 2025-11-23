package handlers

import (
	"net/http"
	"backend-gin/database"
	"backend-gin/models"
	"github.com/gin-gonic/gin"
)

func GetTransactions(c *gin.Context) {
	var trx []models.Transaction
	database.DB.Order("created_at desc").Find(&trx)
	c.JSON(http.StatusOK, gin.H{"data": trx})
}

func GetSummary(c *gin.Context) {
	// Hitung total manual (bisa pakai SQL SUM nanti biar lebih cepat)
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

	c.JSON(http.StatusOK, gin.H{
		"total_income":  income,
		"total_expense": expense,
		"balance":       income - expense,
	})
}