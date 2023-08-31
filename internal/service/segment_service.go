package service

import (
	"DynamicUserSegmentationService/internal/models"
	"DynamicUserSegmentationService/internal/repository"
)

type SegmentService struct {
	repo *repository.SegmentRepository
}

func NewSegmentService(repo *repository.SegmentRepository) *SegmentService {
	return &SegmentService{repo: repo}
}

func (s *SegmentService) GetSegmentByID(id int) (*models.Segment, error) {
	return s.repo.GetSegmentByID(id)
}

func (s *SegmentService) FindSegmentBySlug(slug string) bool {
	return s.repo.FindSegmentBySlug(slug)
}

func (s *SegmentService) AddNewSegment(slug string) error {
	return s.repo.AddNewSegment(slug)
}

func (s *SegmentService) DeleteSegmentBySlug(slug string) error {
	return s.repo.DeleteSegmentBySlug(slug)
}
