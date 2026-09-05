package handlers

import (
	"mini-paas/internals/dto"
	"mini-paas/internals/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	AiService services.AIService
	LogService services.LogService
}

func NewAIHandler(aiService *services.AIService, logService *services.LogService) *AIHandler {
	return &AIHandler{
		AiService: *aiService,
		LogService: *logService,
	}
}

func (h *AIHandler) AnalyzeLogs(c *gin.Context) {
	var req dto.AnalyzeLogResponse

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	logs, logErr := h.LogService.GetContainerLogs(c.Request.Context(), req.ContainerName, req.LineNo)
	if logErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": logErr.Error(),
		})
		return 
	}

	result, analyzeErr := h.AiService.AnalyzeContainerLogs(c.Request.Context(), logs, "")
	if analyzeErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": analyzeErr.Error(),
		})
		return 
	}

	c.JSON(http.StatusOK, gin.H{
		"severity":  result.Severity,
		"summary":   result.Summary,
		"diagnosis": result.Diagonsis,
		"solution":  result.Solution,
	})
}