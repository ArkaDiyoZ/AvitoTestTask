package main

import (
	_ "DynamicUserSegmentationService/docs"
	"DynamicUserSegmentationService/internal/config"
	"DynamicUserSegmentationService/internal/handlers/history"
	"DynamicUserSegmentationService/internal/handlers/segment"
	"DynamicUserSegmentationService/internal/handlers/user"
	"DynamicUserSegmentationService/internal/repository"
	"DynamicUserSegmentationService/internal/scheduler"
	"DynamicUserSegmentationService/internal/service"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.

// @host      localhost:8080
// @BasePath  /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// close db (defer)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.New(os.Stdout, "log:: ", log.LstdFlags)

	db := config.Connection()
	if db != nil {
		logger.Printf("\nConnected to database")
	}

	userRepository := repository.NewUserRepository(db)
	segmentRepository := repository.NewSegmentRepository(db)
	historyRepository := repository.NewHistoryRepository(db)

	userService := service.NewUserService(userRepository)
	segmentService := service.NewSegmentService(segmentRepository)
	historyService := service.NewHistoryService(historyRepository)

	go scheduler.StartSegmentCleanupRoutine(ctx, userRepository, logger)

	router := gin.Default()

	userHandler := user.NewHandler(userService) //repo by pointer
	segmentHandler := segment.NewHandler(segmentService)
	historyHandler := history.NewHandler(historyService)

	router.POST("/segments", segmentHandler.AddSegmentHandler)
	router.DELETE("/segments/:slug", segmentHandler.DeleteSegmentHandler)

	router.POST("/users", userHandler.AddUserHandler)
	router.GET("/users/:id", userHandler.GetUserHandler)
	router.POST("/users/:id/segments", userHandler.AddUserToSegmentHandler)
	router.GET("/users/:id/segments", userHandler.GetUserSegmentsHandler)
	router.DELETE("/users/:id/segments/slug", userHandler.DeleteUserSegmentsHandler)

	router.GET("/history", historyHandler.GenerateReportHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	go router.Run(":8080")
	<-signalCh
	fmt.Println("Interrupt signal received")
}
