package repository

import (
	"DynamicUserSegmentationService/models"
	"gorm.io/gorm"
	"time"
)

type HistoryRepository struct {
	db *gorm.DB
}

func NewHistoryRepository(db *gorm.DB) *HistoryRepository {
	return &HistoryRepository{db}
}

func (r *HistoryRepository) GetRecordsForPeriod(start time.Time, end time.Time) ([]models.History, error) {
	var historyRecords []models.History
	if err := r.db.Table("history").Where("timestamp >= ? AND timestamp <= ?", start, end).
		Find(&historyRecords).Error; err != nil {
		return nil, err
	}
	return historyRecords, nil
}
