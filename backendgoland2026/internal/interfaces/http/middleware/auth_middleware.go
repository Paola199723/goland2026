package middleware

import (
	"net/http"
	"strings"

	"github.com/Paola199723/backendgoland2026/internal/infrastructure/services"
	"github.com/Paola199723/backendgoland2026/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	authService := services.NewAuthService()

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusBadRequest, dto.AuthErrorResponse{
				Error:   "missing_token",
				Message: "Authorization token is required",
			})
			c.Abort()
			return
		}

		// Extraer token del header "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusBadRequest, dto.AuthErrorResponse{
				Error:   "invalid_token_format",
				Message: "Token format should be: Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Verificar token
		claims, err := authService.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.AuthErrorResponse{
				Error:   "token_verification_failed",
				Message: "Token has expired or is invalid. Please login again.",
			})
			c.Abort()
			return
		}

		// Agregar información del usuario al contexto
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}
