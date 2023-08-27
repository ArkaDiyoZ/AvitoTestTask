package handlers

import (
	"DynamicUserSegmentationService/internal/repository"
	"DynamicUserSegmentationService/models"
	"DynamicUserSegmentationService/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetUserHandler(c *gin.Context, userRepo *repository.UserRepository) {
	userIdParam := c.Param("id")
	userID, err := strconv.Atoi(userIdParam)

	userService := service.NewUserService(userRepo)
	user, err := userService.GetUserByID(userID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to get user")
		return
	}

	c.JSON(http.StatusOK, user)
}

func AddUserHandler(c *gin.Context, userRepo *repository.UserRepository) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.String(http.StatusBadRequest, "Invalid data")
		return
	}

	userService := service.NewUserService(userRepo)
	if err := userService.AddUser(user.Name); err != nil {
		c.String(http.StatusInternalServerError, "Failed to add user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "user": user.Name})
	//todo вернуть id пользователя
}

func AddUserToSegmentHandler(c *gin.Context, userRepo *repository.UserRepository) {
	userIdParam := c.Param("id")
	userID, err := strconv.Atoi(userIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userService := service.NewUserService(userRepo)

	if userService.FindUserById(userID) != true {
		c.JSON(http.StatusNotFound, "User not found")
		return
	}

	var segments = models.SegmentsRequest{}
	if err := c.ShouldBindJSON(&segments); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid data")
		return
	}
	//todo add check to empty segments array and mb to incorrect name in body
	//todo check if this user already in this segment ( optional )
	//todo log the errors !!!!!!

	if err := userService.AddUserToSegment(userID, segments.Segment); err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to add user to segment")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "user": userID, "segment": segments.Segment})
}

func GetUserSegmentsHandler(c *gin.Context, userRepo *repository.UserRepository) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userService := service.NewUserService(userRepo)
	if userService.FindUserById(userID) != true {
		c.JSON(http.StatusNotFound, "User not found")
		return
	}

	userSegments, err := userService.GetUserSegments(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to get user segments")
		return
	}

	c.JSON(http.StatusOK, userSegments)
}

func DeleteUserSegmentsHandler(c *gin.Context, userRepo *repository.UserRepository) {
	userIdParam := c.Param("id")
	userID, err := strconv.Atoi(userIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userService := service.NewUserService(userRepo)

	if userService.FindUserById(userID) != true {
		c.JSON(http.StatusNotFound, "User not found")
		return
	}

	var segments = models.SegmentsRequest{}
	if err := c.ShouldBindJSON(&segments); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid data")
		return
	}

	if err := userService.DeleteUserSegments(userID, segments.Segment); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "deleted": segments.Segment, "user": userID})
}
