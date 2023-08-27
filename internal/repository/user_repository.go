package repository

import (
	"DynamicUserSegmentationService/models"
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
