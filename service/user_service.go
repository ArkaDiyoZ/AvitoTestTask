package service

import (
	"DynamicUserSegmentationService/internal/repository"
	"DynamicUserSegmentationService/models"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) AddUser(name string) error {
	return s.repo.AddUser(name)
}

func (s *UserService) FindUserById(id int) bool {
	return s.repo.FindUserById(id)
}

func (s *UserService) AddUserToSegment(id int, segments []string) error {
	return s.repo.AddUserToSegment(id, segments)
}

func (s *UserService) GetUserSegments(id int) ([]models.Segment, error) {
	return s.repo.GetUserSegments(id)
}

func (s *UserService) DeleteUserSegments(id int, segments []string) error {
	return s.repo.DeleteUserSegments(id, segments)
}
