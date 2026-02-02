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
		"http://localhost:8081",
		"http://127.0.0.1:5173",
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}

	r.Use(cors.New(config))

	routes.Register(r)

	log.Println("🚀 Server running on http://localhost:8080")
	log.Println("📊 Frontend available at http://localhost:5173")

	r.Run(":8081")
}
