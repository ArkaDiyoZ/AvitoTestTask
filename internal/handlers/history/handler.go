package history

import (
	"DynamicUserSegmentationService/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Handler struct {
	historyService *service.HistoryService
}

func NewHandler(historyService *service.HistoryService) Handler {
	return Handler{historyService: historyService}
}

func (h Handler) GenerateReportHandler(c *gin.Context) {
	startParam := c.Query("start")
	endParam := c.Query("end")

	start, err := time.Parse("2006-01-02", startParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid start date")
		return
	}

	end, err := time.Parse("2006-01-02", endParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid end date")
		return
	}

	reportFilePath, err := h.historyService.GenerateReport(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to generate report")
		return
	}

	c.JSON(http.StatusOK, gin.H{"report_link": reportFilePath})
}
