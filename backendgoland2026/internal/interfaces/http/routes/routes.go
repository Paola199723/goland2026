package routes

import (
	"github.com/Paola199723/backendgoland2026/internal/application/usecases"
	"github.com/Paola199723/backendgoland2026/internal/infrastructure/persistence"
	"github.com/Paola199723/backendgoland2026/internal/interfaces/http/handlers"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	// Challenge routes
	challengeRepo := persistence.NewChallengeRepository()
	loginUC := usecases.NewLoginUseCase(challengeRepo)
	getChallengesUC := usecases.NewGetChallengesUseCase(challengeRepo)
	challengeHandler := handlers.NewChallengeHandler(loginUC, getChallengesUC)

	api := r.Group("/api")
	{
		api.POST("/login", challengeHandler.Login)
		api.GET("/challenges", challengeHandler.GetChallenges)
	}

	// User routes (mantener existentes)
	userRepo := persistence.NewUserRepository()
	createUC := usecases.NewCreateUserUseCase(userRepo)
	userHandler := handlers.NewUserHandler(createUC)

	r.POST("/users", userHandler.Create)
}
