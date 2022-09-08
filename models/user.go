package models

import (
	"errors"
	"strings"

	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email       string `json:"-" gorm:"not null;unique"`
	Username    string `json:"username" gorm:"not null;unique"`
	Password    string `json:"-" gorm:"not null"`
	DisplayName string `json:"displayName"`
	IsActive    bool   `json:"isActive" gorm:"default:true"`
	IsLocked    bool   `json:"isLocked" gorm:"default:false"`
}

type UserPayload struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func PayloadToUser(p UserPayload) User {
	return User{
		Email:    p.Email,
		Username: p.Username,
		Password: p.Password,
	}
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	tx.Statement.SetColumn("password", hashedPassword)

	return nil
}

type UserDTO struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
}

func DTOToUser(dto UserDTO) User {
	return User{
		Username:    dto.Username,
		Password:    dto.Password,
		DisplayName: dto.DisplayName,
	}
}

func (u User) Validate(action string) error {
	validate := validator.New()

	switch strings.ToLower(action) {
	case "register":
		if err := validate.Var(u.Email, "required,email"); err != nil {
			return errors.New("you have to provide a valid email")
		}

		if err := validate.Var(u.Username, "required"); err != nil {
			return errors.New("you have to provide a username")
		}

		if len(u.Password) < 8 {
			return errors.New("password must be at least 8 characters")
		}
	case "update":
		if err := validate.Var(u.Email, "required,email"); err != nil && u.Email != "" {
			return errors.New("you have to provide a valid email")
		}

		if len(u.Password) < 8 && u.Password != "" {
			return errors.New("password must be at least 8 characters")
		}
	}

	return nil
}
