package handlers

import (
	"DynamicUserSegmentationService/internal/repository"
	"DynamicUserSegmentationService/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GenerateReportHandler(c *gin.Context, historyRepo *repository.HistoryRepository) {
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

	historyService := service.NewHistoryService(historyRepo)
	reportFilePath, err := historyService.GenerateReport(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to generate report")
		return
	}

	c.JSON(http.StatusOK, gin.H{"report_link": reportFilePath})
}
