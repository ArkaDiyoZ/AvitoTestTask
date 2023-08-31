package scheduler

import (
	"DynamicUserSegmentationService/internal/repository"
	"fmt"
	"log"
	"time"
)

func StartSegmentCleanupRoutine(userRepo *repository.UserRepository, logger *log.Logger) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	logger.Printf("searching user segments to delete ...")
	for {
		select {
		case <-ticker.C:
			currentTime := time.Now()
			if err := userRepo.DeleteExpiredUserSegments(currentTime); err != nil {
				fmt.Errorf("shceduler delete error", err)
			}
		}
	}
}
