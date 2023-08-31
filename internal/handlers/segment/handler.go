package segment

import (
	"DynamicUserSegmentationService/internal/models"
	"DynamicUserSegmentationService/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	segmentService *service.SegmentService
}

func NewHandler(segmentService *service.SegmentService) Handler {
	return Handler{segmentService: segmentService}
}

func (h Handler) AddSegmentHandler(c *gin.Context) {
	var segment models.Segment

	if err := c.ShouldBindJSON(&segment); err != nil {
		c.String(http.StatusBadRequest, "Invalid data")
		return
	}

	if h.segmentService.FindSegmentBySlug(segment.Slug) != false {
		c.JSON(http.StatusConflict, "Segment already exist")
		return
	}

	if err := h.segmentService.AddNewSegment(segment.Slug); err != nil {
		c.String(http.StatusInternalServerError, "Failed to add segment")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "slug": segment.Slug})
}

func (h Handler) DeleteSegmentHandler(c *gin.Context) {
	segmentSlugParam := c.Param("slug")

	if h.segmentService.FindSegmentBySlug(segmentSlugParam) != true {
		c.JSON(http.StatusNotFound, "Segment not found")
		return
	}

	if err := h.segmentService.DeleteSegmentBySlug(segmentSlugParam); err != nil {
		c.String(http.StatusInternalServerError, "Failed to delete segment")
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}
