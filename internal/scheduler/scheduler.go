package scheduler

import (
	"DynamicUserSegmentationService/internal/repository"
	"context"
	"fmt"
	"log"
	"time"
)

func StartSegmentCleanupRoutine(ctx context.Context, userRepo *repository.UserRepository, logger *log.Logger) {
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
		case <-ctx.Done():
			return
		}
	}
}
