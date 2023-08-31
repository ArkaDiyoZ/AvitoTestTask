package user

import (
	"DynamicUserSegmentationService/internal/models"
	"DynamicUserSegmentationService/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	userService *service.UserService
}

func NewHandler(userService *service.UserService) Handler {
	return Handler{userService: userService}
}

type SegmentsRequest struct {
	Segment        []string  `json:"segments"`
	ExpirationTime time.Time `json:"expiration_time"`
}

func (h Handler) GetUserHandler(c *gin.Context) {
	userIdParam := c.Param("id")
	userID, err := strconv.Atoi(userIdParam)

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.String(http.StatusNotFound, "Failed to get user")
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h Handler) AddUserHandler(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.String(http.StatusBadRequest, "Invalid data")
		return
	}

	userService := h.userService
	if err := userService.AddUser(user); err != nil {
		c.String(http.StatusInternalServerError, "Failed to add user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "user": user.Name})
}

func (h Handler) AddUserToSegmentHandler(c *gin.Context) {
	userIdParam := c.Param("id")
	userID, err := strconv.Atoi(userIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if exist, err := h.userService.UserExist(userID); err != nil || !exist {
		c.JSON(http.StatusNotFound, "User not found")
		return
	}

	var segments = SegmentsRequest{} // убрать в хендлер из модели
	if err := c.ShouldBindJSON(&segments); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid data")
		return
	}

	if err := h.userService.AddUserToSegment(userID, segments.Segment, segments.ExpirationTime); err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to add user to segment")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "user": userID, "segment": segments.Segment})
}

func (h Handler) GetUserSegmentsHandler(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if exist, err := h.userService.UserExist(userID); err != nil || !exist {
		c.JSON(http.StatusNotFound, "User not found")
		return
	}
	userSegments, err := h.userService.GetUserSegments(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to get user segments")
		return
	}

	c.JSON(http.StatusOK, userSegments)
}

func (h Handler) DeleteUserSegmentsHandler(c *gin.Context) {
	userIdParam := c.Param("id")
	userID, err := strconv.Atoi(userIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if exist, err := h.userService.UserExist(userID); err != nil || !exist {
		c.JSON(http.StatusNotFound, "User not found")
		return
	}

	var segments = SegmentsRequest{}
	if err := c.ShouldBindJSON(&segments); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid data")
		return
	}

	if err := h.userService.DeleteUserSegments(userID, segments.Segment); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "deleted": segments.Segment, "user": userID})
}
