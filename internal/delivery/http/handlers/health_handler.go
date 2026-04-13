package handlers

import (
	"net/http"

	"mini-payment-system/pkg/response"

	"github.com/gin-gonic/gin"
)

// Health godoc
// @Summary Health check
// @Description Check whether the service is running.
// @Tags health
// @Produce json
// @Success 200 {object} response.HealthResponse
// @Router /health [get]
func Health(c *gin.Context) {
	response.JSON(c, http.StatusOK, response.HealthData{Status: "ok"})
}
