package repository

import (
	"DynamicUserSegmentationService/models"
	"errors"
	"gorm.io/gorm"
	"time"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) AddUser(name string) error {
	user := models.User{Name: name}
	if err := r.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindUserById(id int) bool {
	var user models.User
	if err := r.db.Where("id = ?", id).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
		return false
	}
	return true
}

func (r *UserRepository) AddUserToSegment(id int, segmentsSlugs []string) error {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return err
	}

	var segments []models.Segment
	if err := r.db.Where("slug IN ?", segmentsSlugs).Find(&segments).Error; err != nil {
		return err
	}

	var existingSegments []models.Segment
	if err := r.db.Where("slug IN ?", segmentsSlugs).Find(&existingSegments).Error; err != nil {
		return err
	}

	if len(existingSegments) != len(segmentsSlugs) {
		return errors.New("not all segments exist")
	}

	transaction := r.db.Begin() // Начать транзакцию

	for _, segment := range existingSegments {
		userSegment := models.UserSegment{
			UserID:    id,
			SegmentID: segment.ID,
		}

		if err := transaction.Create(&userSegment).Error; err != nil {
			transaction.Rollback()
			return err
		}

		// Создать запись в таблице истории
		historyRecord := models.History{
			UserID:    id,
			SegmentID: segment.ID,
			Operation: "user adding segments",
			Timestamp: time.Now(),
		}
		if err := transaction.Table("history").Create(&historyRecord).Error; err != nil {
			transaction.Rollback() // Откатить транзакцию при ошибке
			return err
		}
	}

	transaction.Commit() // Подтвердить транзакцию

	return nil
}

func (r *UserRepository) GetUserSegments(id int) ([]models.Segment, error) {
	var segments []models.Segment
	err := r.db.
		Joins("JOIN user_segments ON segments.id = user_segments.segment_id").
		Where("user_segments.user_id = ?", id).
		Find(&segments).Error
	if err != nil {
		return nil, err
	}
	return segments, nil
}

func (r *UserRepository) DeleteUserSegments(id int, segmentsSlugs []string) error {
	transaction := r.db.Begin()

	var user models.User
	if err := transaction.First(&user, id).Error; err != nil {
		transaction.Rollback() // Откат транзакции в случае ошибки
		return errors.New("user not found")
	}

	var segmentIDs []int
	if err := transaction.Model(&models.Segment{}).
		Select("id").
		Where("slug IN ?", segmentsSlugs).
		Find(&segmentIDs).Error; err != nil {
		transaction.Rollback()
		return err
	}

	if len(segmentIDs) == 0 {
		transaction.Rollback()
		return errors.New("no valid segment IDs found")
	}

	if len(segmentIDs) != len(segmentsSlugs) {
		transaction.Rollback()
		return errors.New("not all segments found or belong to the user")
	}

	var userSegments []models.UserSegment
	if err := transaction.Model(&models.UserSegment{}).
		Where("user_id = ?", id).
		Find(&userSegments).Error; err != nil {
		transaction.Rollback()
		return err
	}

	userSegmentMap := make(map[int]bool)
	for _, userSegment := range userSegments {
		userSegmentMap[userSegment.SegmentID] = true
	}

	for _, segmentSlug := range segmentsSlugs {
		var segment models.Segment
		if err := transaction.Model(&models.Segment{}).
			Where("slug = ?", segmentSlug).
			First(&segment).Error; err != nil {
			transaction.Rollback()
			return err
		}

		if !userSegmentMap[segment.ID] {
			transaction.Rollback()
			return errors.New("segment not found or doesn't belong to the user")
		}
	}

	if err := transaction.
		Where("user_id = ?", id).
		Where("segment_id IN ?", segmentIDs).
		Delete(&models.UserSegment{}).Error; err != nil {
		transaction.Rollback()
		return err
	}

	for _, segmentID := range segmentIDs {
		historyRecord := models.History{
			UserID:    id,
			SegmentID: segmentID,
			Operation: "user removing segment",
			Timestamp: time.Now(),
		}
		if err := transaction.Table("history").Create(&historyRecord).Error; err != nil {
			transaction.Rollback()
			return err
		}
	}

	return transaction.Commit().Error
}
