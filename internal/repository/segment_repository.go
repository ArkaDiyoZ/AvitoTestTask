package repository

import (
	"DynamicUserSegmentationService/models"
	"errors"
	"gorm.io/gorm"
)

type SegmentRepository struct {
	db *gorm.DB
}

func NewSegmentRepository(db *gorm.DB) *SegmentRepository {
	return &SegmentRepository{db}
}

func (r *SegmentRepository) GetSegmentByID(id int) (*models.Segment, error) {
	var segment models.Segment
	if err := r.db.First(&segment, id).Error; err != nil {
		return nil, err
	}
	return &segment, nil
}

func (r *SegmentRepository) FindSegmentBySlug(slug string) (bool, error) {
	var segment models.Segment
	if err := r.db.Where("slug = ?", slug).Take(&segment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *SegmentRepository) AddNewSegment(slug string) error {
	segment := models.Segment{Slug: slug}
	if err := r.db.Create(&segment).Error; err != nil {
		return err
	}
	return nil
}

func (r *SegmentRepository) DeleteSegmentBySlug(slug string) error {
	var segment models.Segment

	if err := r.db.Where("slug = ?", slug).Delete(&segment).Error; err != nil {
		return err
	}

	return nil
}
