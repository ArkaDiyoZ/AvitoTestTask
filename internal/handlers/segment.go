package handlers

import (
	"DynamicUserSegmentationService/internal/models"
	"DynamicUserSegmentationService/internal/repository"
	"DynamicUserSegmentationService/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddSegmentHandler(c *gin.Context, segmentRepo *repository.SegmentRepository) {
	var segment models.Segment

	if err := c.ShouldBindJSON(&segment); err != nil {
		c.String(http.StatusBadRequest, "Invalid data")
		return
	}

	segmentService := service.NewSegmentService(segmentRepo)

	if segmentService.FindSegmentBySlug(segment.Slug) != false {
		c.JSON(http.StatusConflict, "Segment already exist")
		return
	}

	if err := segmentService.AddNewSegment(segment.Slug); err != nil {
		c.String(http.StatusInternalServerError, "Failed to add segment")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "slug": segment.Slug})
	//todo вернуть id slug
}

func DeleteSegmentHandler(c *gin.Context, segmentRepo *repository.SegmentRepository) {
	segmentSlugParam := c.Param("slug")

	segmentService := service.NewSegmentService(segmentRepo)

	if segmentService.FindSegmentBySlug(segmentSlugParam) != true {
		c.JSON(http.StatusNotFound, "Segment not found")
		return
	}

	if err := segmentService.DeleteSegmentBySlug(segmentSlugParam); err != nil {
		c.String(http.StatusInternalServerError, "Failed to delete segment")
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}
