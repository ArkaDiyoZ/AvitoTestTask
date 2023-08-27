package main

import (
	"DynamicUserSegmentationService/database"
	"DynamicUserSegmentationService/handlers"
	"DynamicUserSegmentationService/internal/repository"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {

	logger := log.New(os.Stdout, "log:: ", log.LstdFlags)

	db := database.Connection(logger)
	if db != nil {
		logger.Printf("\nConnected to database")
	}

	userRepository := repository.NewUserRepository(db)
	segmentRepository := repository.NewSegmentRepository(db)
	historyRepository := repository.NewHistoryRepository(db)

	router := gin.Default()
	router.POST("/segments", func(context *gin.Context) {
		handlers.AddSegmentHandler(context, segmentRepository)
	})
	router.DELETE("/segments/:slug", func(context *gin.Context) {
		handlers.DeleteSegmentHandler(context, segmentRepository)
	})
	router.GET("/users/:id", func(context *gin.Context) {
		handlers.GetUserHandler(context, userRepository)
	})
	router.POST("/users", func(context *gin.Context) {
		handlers.AddUserHandler(context, userRepository)
	})
	router.POST("/users/:id/segments", func(context *gin.Context) {
		handlers.AddUserToSegmentHandler(context, userRepository)
	})
	router.GET("/users/:id/segments", func(context *gin.Context) {
		handlers.GetUserSegmentsHandler(context, userRepository)
	})
	router.DELETE("/users/:id/segments/slug", func(context *gin.Context) {
		handlers.DeleteUserSegmentsHandler(context, userRepository)
	})
	router.GET("/history", func(context *gin.Context) {
		handlers.GenerateReportHandler(context, historyRepository)
	})

	router.Run(":8080")

}
