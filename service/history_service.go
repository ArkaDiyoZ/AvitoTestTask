package service

import (
	"DynamicUserSegmentationService/internal/repository"
	"encoding/csv"
	"os"
	"strconv"
	"time"
)

type HistoryService struct {
	repo *repository.HistoryRepository
}

func NewHistoryService(repo *repository.HistoryRepository) *HistoryService {
	return &HistoryService{repo}
}

func (s *HistoryService) GenerateReport(start time.Time, end time.Time) (string, error) {
	historyRecords, err := s.repo.GetRecordsForPeriod(start, end)
	if err != nil {
		return "", err
	}

	reportFilePath := "report.csv"
	file, err := os.Create(reportFilePath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range historyRecords {
		row := []string{
			strconv.Itoa(record.UserID),
			strconv.Itoa(record.SegmentID),
			record.Operation,
			record.Timestamp.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}

	return reportFilePath, nil
}
