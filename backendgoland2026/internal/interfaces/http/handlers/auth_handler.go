package handlers

import (
	"net/http"

	"github.com/Paola199723/backendgoland2026/internal/application/usecases"
	"github.com/Paola199723/backendgoland2026/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authLoginUC *usecases.AuthLoginUseCase
}

func NewAuthHandler(authLoginUC *usecases.AuthLoginUseCase) *AuthHandler {
	return &AuthHandler{
		authLoginUC: authLoginUC,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password, returns JWT token
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.AuthTokenResponse
// @Failure 400 {object} dto.AuthErrorResponse
// @Failure 401 {object} dto.AuthErrorResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.AuthErrorResponse{
			Error:   "invalid_request",
			Message: "Email and password are required",
		})
		return
	}

	response, err := h.authLoginUC.Execute(req.Email, req.Password)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMsg := err.Error()

		if errorMsg == "invalid email format" {
			statusCode = http.StatusBadRequest
		} else if errorMsg == "invalid credentials" {
			statusCode = http.StatusUnauthorized
		}

		c.JSON(statusCode, dto.AuthErrorResponse{
			Error:   "authentication_failed",
			Message: errorMsg,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Verify godoc
// @Summary Verify token
// @Description Verify if a token is valid
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} dto.AuthErrorResponse
// @Router /api/auth/verify [get]
func (h *AuthHandler) VerifyToken(c *gin.Context) {
	// Middleware will handle token verification
	c.JSON(http.StatusOK, gin.H{
		"message": "token_valid",
		"status":  "authenticated",
	})
}
