package main

import (
	"backend-gin/database"
	"backend-gin/handlers"
	"backend-gin/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Konek Database
	database.ConnectDatabase()

	r := gin.Default()

	// 2. Route Publik (Bisa diakses siapa saja)
	r.POST("/login", handlers.Login)
	r.POST("/register-admin", handlers.RegisterAdmin) // Panggil ini sekali aja nanti via Postman/Curl
	r.POST("/telegram/webhook", handlers.TelegramWebhook)

	// 3. Route Privat (Harus Login / Pakai Token)
	api := r.Group("/api")
	api.Use(middleware.JwtAuthMiddleware()) // Pasang Satpam
	{
		api.GET("/transactions", handlers.GetTransactions)
		api.GET("/summary", handlers.GetSummary)
	}

	r.Run(":8080")
}