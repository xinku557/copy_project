package models

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title       string `json:"title" gorm:"not null"`
	Body        string `json:"body"`
	AuthorID    *uint  `json:"authorId" gorm:"not null"`
	Author      *User  `json:"author"`
	IsPublished bool   `json:"isPublished" gorm:"default:false"`
}

type PostDTO struct {
	Title       string `json:"title"`
	Body        string `json:"body"`
	IsPublished bool   `json:"isPublished"`
}

func DTOToPost(dto PostDTO) Post {
	return Post{
		Title:       dto.Title,
		Body:        dto.Body,
		IsPublished: dto.IsPublished,
	}
}

func (p Post) Validate(action string) error {
	switch strings.ToLower(action) {
	case "create":
		if len(p.Title) < 3 {
			return errors.New("title must be at least 3 characters long")
		}
		if len(p.Body) < 3 {
			return errors.New("content must be at least 3 characters long")
		}
	case "update":
		if len(p.Title) < 3 && p.Title != "" {
			return errors.New("title must be at least 3 characters long")
		}
		if len(p.Body) < 3 && p.Body != "" {
			return errors.New("content must be at least 3 characters long")
		}
	}

	return nil
}
