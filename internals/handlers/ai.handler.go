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
	
	err := c.ShouldBindJSON(&req); if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})

		return
	}

	logs, logErr := h.LogService.GetContainerLogs(c.Request.Context(), req.ContainerName, req.LineNo)

	if logErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": logErr.Error(),
		})
	}

	result, analyzeErr := h.AiService.AnalyzeContainerLogs(c.Request.Context(), logs)

	if analyzeErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": analyzeErr.Error(),
		})
	}

	c.JSON(200, gin.H{
		"message": "Logs Verified",
		"Logs": result,
	})

} 