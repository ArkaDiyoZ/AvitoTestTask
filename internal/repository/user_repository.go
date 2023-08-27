package repository

import (
	"DynamicUserSegmentationService/models"
	"errors"
	"gorm.io/gorm"
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

	for _, segment := range segments {
		// Связываем сегмент с пользователем
		userSegment := models.UserSegment{
			UserID:    id,
			SegmentID: segment.ID,
		}
		if err := r.db.Create(&userSegment).Error; err != nil {
			return err
		}
	}

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
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return errors.New("test")
	}

	//todo check if slug is really exist
	var segmentIDs []int
	if err := r.db.Model(&models.Segment{}).
		Select("id").
		Where("slug IN ?", segmentsSlugs).
		Find(&segmentIDs).Error; err != nil {
		return err
	}

	if len(segmentIDs) == 0 {
		return errors.New("no valid segment IDs found")
	}

	if len(segmentIDs) != len(segmentsSlugs) {
		return errors.New("not all segments found or belong to the user")
	}

	if err := r.db.
		Where("user_id = ?", id).
		Where("segment_id IN ?", segmentIDs).
		Delete(&models.UserSegment{}); err != nil {
		return err.Error
	}

	return nil
}
