package main

import (
	"log"

	"github.com/Paola199723/backendgoland2026/internal/infrastructure/db"
	"github.com/Paola199723/backendgoland2026/internal/interfaces/http/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()

	r := gin.Default()

	// Configurar CORS para permitir requests desde el frontend
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:5173",
		"http://localhost:5174",
		"http://localhost:5175",
		"http://localhost:5176",
		"http://localhost:5177",
		"http://localhost:8081",
		"http://127.0.0.1:5173",
		"http://127.0.0.1:5174",
		"http://127.0.0.1:5175",
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	routes.Register(r)

	log.Println("🚀 Server running on http://localhost:8081")
	log.Println("📊 Frontend available at http://localhost:5173")

	r.Run(":8081")
}
