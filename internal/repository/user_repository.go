package repository

import (
	"DynamicUserSegmentationService/internal/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) GetUserByID(id int) (models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (r *UserRepository) AddUser(name string) error {
	exists, err := r.CheckUserExists(name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("user already exists", err)
	}

	user := models.User{
		Name: name,
	}
	if err := r.db.Create(&user).Error; err != nil {
		return fmt.Errorf("user creation error", err)
	}
	return nil
}

func (r *UserRepository) CheckUserExists(name string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) UserExist(id int) (bool, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, err
		}
		return false, err
	}
	return true, nil
}

func (r *UserRepository) AddUserToSegment(userId int, segmentsSlugs []string, expirationTime time.Time) error {
	var user models.User
	var defaultExpirationTime = time.Date(1099, time.December, 31, 23, 59, 59, 0, time.UTC)

	if err := r.db.First(&user, userId).Error; err != nil {
		return err
	}

	var segments []models.Segment
	if err := r.db.Where("slug IN ?", segmentsSlugs).Find(&segments).Error; err != nil {
		return err
	}

	var existingSegments []models.Segment
	if err := r.db.Where("slug IN ?", segmentsSlugs).Find(&existingSegments).Error; err != nil {
		return err
	} //подумать зачем это

	if expirationTime.IsZero() {
		expirationTime = defaultExpirationTime
	}

	transaction := r.db.Begin() // Начать транзакцию

	for _, segment := range existingSegments { // батч вместо цикла
		userSegment := models.UserSegment{
			UserID:         userId,
			SegmentID:      segment.ID,
			ExpirationTime: expirationTime,
		}

		if err := transaction.Create(&userSegment).Error; err != nil {
			transaction.Rollback()
			return err
		}

		// Создать запись в таблице истории
		historyRecord := models.History{
			UserID:      userId,
			SegmentID:   segment.ID,
			Operation:   "user adding segments",
			CreatedTime: time.Now(),
		}
		if err := transaction.Table("history").Create(&historyRecord).Error; err != nil {
			transaction.Rollback()
			return err
		}
	}

	transaction.Commit()

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
			UserID:      id,
			SegmentID:   segmentID,
			Operation:   "user removing segment",
			CreatedTime: time.Now(),
		}
		if err := transaction.Table("history").Create(&historyRecord).Error; err != nil {
			transaction.Rollback()
			return err
		}
	}

	return transaction.Commit().Error
}

func (r *UserRepository) DeleteExpiredUserSegments(currentTime time.Time) error {
	return r.db.
		Where("expiration_time <= ?", currentTime).
		Delete(&models.UserSegment{}).
		Error
}
