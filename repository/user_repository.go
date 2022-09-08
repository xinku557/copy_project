package repository

import (
	"gorm.io/gorm"
	"sheinko.tk/copy_project/models"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r userRepository) Save(p *models.User) error {
	if err := p.Validate("register"); err != nil {
		return err
	}

	if err := r.db.Create(&p).Error; err != nil {
		return err
	}
	return nil
}

func (r userRepository) FindById(id uint) (models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r userRepository) FindByEmailOrUsername(email string, username string) (models.User, error) {
	var user models.User
	if err := r.db.First(&user, "email = ? OR username = ?", email, username).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r userRepository) UpdateById(value *models.User, newValue *models.User) error {
	if err := newValue.Validate("update"); err != nil {
		return err
	}

	if err := r.db.Model(&value).Updates(&newValue).Error; err != nil {
		return err
	}

	return nil
}
