package handlers

import (
	"bufio"
	"encoding/json"
	"io"
	"mini-paas/internals/dto"
	"mini-paas/internals/services"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	service services.LogService
}

func NewLogHandler(service *services.LogService) *LogHandler {
	return &LogHandler{
		service: *service,
	}
}

func (h *LogHandler) GetAllLogs(c *gin.Context) {
	var req dto.LogRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})

		return
	}

	result, err := h.service.GetContainerLogs(c.Request.Context(), req.ContainerName, req.LineNo)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
	}

	c.JSON(200, gin.H{
		"Result": result,
	})
}

func (h *LogHandler) StreamLogs(c *gin.Context) {
	containerName := c.Param("name")
	tail := c.DefaultQuery("tail", "50")

	streamReader, err := h.service.StreamContainerLogs(c.Request.Context(), containerName, tail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	defer streamReader.Close()

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	scanner := bufio.NewScanner(streamReader)

	c.Stream(func(w io.Writer) bool {
		if scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimRight(line, "\r")

			if len(strings.TrimSpace(line)) == 0 {
				return true
			}

			logEntry, _ := json.Marshal(dto.LogMessage{
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Log: line,
				Source: "stdout",
			})
			

			c.Writer.WriteString("data: " + string(logEntry) + "\n\n")

			if flusher, ok := c.Writer.(http.Flusher); ok {
				flusher.Flush()
			}
			return true
		}
		return false
	})
}