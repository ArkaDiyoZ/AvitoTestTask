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
