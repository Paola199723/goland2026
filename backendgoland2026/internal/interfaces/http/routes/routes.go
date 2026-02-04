package routes

import (
	"github.com/Paola199723/backendgoland2026/internal/application/usecases"
	"github.com/Paola199723/backendgoland2026/internal/infrastructure/persistence"
	"github.com/Paola199723/backendgoland2026/internal/interfaces/http/handlers"
	"github.com/Paola199723/backendgoland2026/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	api := r.Group("/api")

	// Auth routes (sin autenticación)
	userRepo := persistence.NewUserRepository()
	authLoginUC := usecases.NewAuthLoginUseCase(userRepo)
	authHandler := handlers.NewAuthHandler(authLoginUC)
	recommendationHandler := handlers.NewRecommendationHandler()

	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
	}

	// Challenge routes (con autenticación)
	challengeRepo := persistence.NewChallengeRepository()
	getChallengesUC := usecases.NewGetChallengesUseCase(challengeRepo)
	challengeHandler := handlers.NewChallengeHandler(nil, getChallengesUC)

	challenges := api.Group("/challenges")
	challenges.Use(middleware.AuthMiddleware())
	{
		challenges.GET("", challengeHandler.GetChallenges)
		// Endpoint temporal para forzar refresh desde la API externa
		challenges.POST("/refresh", challengeHandler.RefreshChallenges)
	}

	// User routes (mantener existentes)
	createUC := usecases.NewCreateUserUseCase(userRepo)
	userHandler := handlers.NewUserHandler(createUC)

	r.POST("/users", userHandler.Create)

	// Verify token route (con autenticación)
	authVerify := api.Group("/auth")
	authVerify.Use(middleware.AuthMiddleware())
	{
		authVerify.GET("/verify", authHandler.VerifyToken)
	}
	recommendation := api.Group("/recommendation")
	recommendation.Use(middleware.AuthMiddleware())
	{
		recommendation.GET("/today", recommendationHandler.GetTodayRecommendation)
	}
}
