package handlers

import (
	"net/http"
	"time"

	"github.com/Paola199723/backendgoland2026/internal/infrastructure/services"
	"github.com/gin-gonic/gin"
)

type RecommendationHandler struct{}

func NewRecommendationHandler() *RecommendationHandler {
	return &RecommendationHandler{}
}

func (h *RecommendationHandler) GetTodayRecommendation(c *gin.Context) {
	today := time.Now()
	recommendation, err := services.RecommendBestActionOfDayV2(today)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recommendation)
}
