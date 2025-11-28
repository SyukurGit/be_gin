package main

import (
	"log"
	"backend-gin/database"
	"backend-gin/handlers"
	"backend-gin/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/gin-contrib/cors" // <--- 1. JANGAN LUPA IMPORT INI
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDatabase()

	r := gin.Default()

	// 2. PASANG CORS DI SINI
	// Ini yang bikin HTML kamu boleh ambil data
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{
"Origin", "Content-Length", "Content-Type", "Authorization", "ngrok-skip-browser-warning", }
	r.Use(cors.New(config))

	// Route
	r.POST("/login", handlers.Login)
	r.POST("/register-admin", handlers.RegisterAdmin)
	r.POST("/telegram/webhook", handlers.TelegramWebhook)

	api := r.Group("/api")
	api.Use(middleware.JwtAuthMiddleware())
	{
		// 3. DAFTARKAN ROUTE BARU INI
		api.GET("/transactions", handlers.GetTransactions)
		api.GET("/summary", handlers.GetSummary)
		api.GET("/chart/daily", handlers.GetDailyChart)      // Baru
		api.GET("/categories", handlers.GetCategorySummary)  // Baru
	}

	r.Run(":8080")
}
