package handlers

import (
	"DynamicUserSegmentationService/internal/repository"
	"DynamicUserSegmentationService/models"
	"DynamicUserSegmentationService/service"
	"fmt"
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
	fmt.Println(user.Name, user.ID)

	c.JSON(http.StatusOK, gin.H{"status": "success", "user": user.Name})
	//todo вернуть id пользователя
}
