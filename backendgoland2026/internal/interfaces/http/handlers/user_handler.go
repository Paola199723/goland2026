package handlers

import (
	"net/http"

	"github.com/Paola199723/backendgoland2026/internal/application/usecases"
	"github.com/Paola199723/backendgoland2026/internal/domain/entities"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	createUC *usecases.CreateUserUseCase
}

func NewUserHandler(uc *usecases.CreateUserUseCase) *UserHandler {
	return &UserHandler{createUC: uc}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := entities.User{
		ID:    req.ID,
		Name:  req.Name,
		Email: req.Email,
	}

	if err := h.createUC.Execute(user); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "User created"})
}
