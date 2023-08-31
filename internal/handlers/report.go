package handlers

//
//import (
//	"github.com/gin-gonic/gin"
//	"net/http"
//	"time"
//)
//
//type service interface {
//	GenerateReport(start time.Time, end time.Time) (string, error)
//}
//
//type Report struct {
//	Service service
//}
//
//func (r *Report) HandlerHistory(c *gin.Context) {
//	startParam := c.Query("start")
//	endParam := c.Query("end")
//
//	start, err := time.Parse("2006-01-02", startParam)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, "Invalid start date")
//		return
//	}
//
//	end, err := time.Parse("2006-01-02", endParam)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, "Invalid end date")
//		return
//	}
//
//	reportFilePath, err := r.Service.GenerateReport(start, end)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, "Failed to generate report")
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"report_link": reportFilePath})
//}
