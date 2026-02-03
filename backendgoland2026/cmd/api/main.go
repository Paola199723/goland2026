package main

import (
	"log"
	"os"

	"github.com/Paola199723/backendgoland2026/internal/infrastructure/db"
	"github.com/Paola199723/backendgoland2026/internal/interfaces/http/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using system environment variables")
	}

	hostfront := os.Getenv("FRONTEND_HOST")
	if hostfront == "" {
		hostfront = "http://localhost:5173"
	}

	port := os.Getenv("PORT")
	if hostfront == "" {
		hostfront = "http://localhost:5173"
	}
	// Verificar que TOKEN esté disponible
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Println("⚠️ WARNING: TOKEN environment variable not set. External API calls may fail.")
	} else {
		log.Println("✅ TOKEN loaded from environment")
	}

	db.Connect()

	r := gin.Default()

	// Configurar CORS para permitir requests desde el frontend
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		hostfront,
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.AllowCredentials = true

	r.Use(cors.New(config))

	routes.Register(r)

	log.Println("🚀 Server running on http://localhost:8081")
	log.Println("📊 Frontend available at http://localhost:5173")

	r.Run(":" + port)
}
