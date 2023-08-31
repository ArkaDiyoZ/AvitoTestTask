package service

import (
	"DynamicUserSegmentationService/internal/models"
	"DynamicUserSegmentationService/internal/repository"
	"time"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByID(id int) (models.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) AddUser(user models.User) error {
	return s.repo.AddUser(user.Name)
}

func (s *UserService) UserExist(id int) (bool, error) {
	return s.repo.UserExist(id)
}

func (s *UserService) AddUserToSegment(userId int, segments []string, expirationTime time.Time) error {
	return s.repo.AddUserToSegment(userId, segments, expirationTime)
}

func (s *UserService) GetUserSegments(id int) ([]models.Segment, error) {
	return s.repo.GetUserSegments(id)
}

func (s *UserService) DeleteUserSegments(id int, segments []string) error {
	return s.repo.DeleteUserSegments(id, segments)
}
