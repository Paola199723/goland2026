package handlers

import (
	"net/http"
	"strconv"

	"github.com/Paola199723/backendgoland2026/internal/application/usecases"
	"github.com/Paola199723/backendgoland2026/internal/interfaces/dto"
	"github.com/gin-gonic/gin"
)

type ChallengeHandler struct {
	loginUC         *usecases.LoginUseCase
	getChallengesUC *usecases.GetChallengesUseCase
}

func NewChallengeHandler(loginUC *usecases.LoginUseCase, getChallengesUC *usecases.GetChallengesUseCase) *ChallengeHandler {
	return &ChallengeHandler{
		loginUC:         loginUC,
		getChallengesUC: getChallengesUC,
	}
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and get challenges
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.LoginResponse
// @Router /api/login [post]
func (h *ChallengeHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.loginUC.Execute(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetChallenges godoc
// @Summary Get challenges by page
// @Description Get challenges with pagination using next_page cursor or page number
// @Produce json
// @Param next_page query string false "Next page cursor"
// @Param page query int false "Page number" default(1)
// @Success 200 {object} dto.LoginResponse
// @Router /api/challenges [get]
func (h *ChallengeHandler) GetChallenges(c *gin.Context) {
	// Obtener email del contexto (puesto por el middleware de auth)
	email, exists := c.Get("email")
	if !exists {
		email = ""
	}
	
	// Soportar tanto next_page cursor como page number
	nextPageCursor := c.Query("next_page")
	pageStr := c.DefaultQuery("page", "1")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Si hay next_page, usarlo para obtener datos
	// Si no, usar el page number
	var response *dto.LoginResponse
	
	if nextPageCursor != "" {
		// Usar cursor para paginación
		response, err = h.getChallengesUC.ExecuteByPageCursor(nextPageCursor)
	} else {
		// Usar page number para paginación
		response, err = h.getChallengesUC.ExecuteByPage(page)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Agregar email a la respuesta
	response.Email = email.(string)

	c.JSON(http.StatusOK, response)
}

// RefreshChallenges - fuerza la obtención desde la API externa y guarda en BD
func (h *ChallengeHandler) RefreshChallenges(c *gin.Context) {
	// Forzamos la ejecución para la página 1
	response, err := h.getChallengesUC.ExecuteByPage(1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

