package repository

import (
	"gorm.io/gorm"
	"sheinko.tk/copy_project/models"
)

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *postRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Save(p *models.Post) error {
	if err := p.Validate("create"); err != nil {
		return err
	}

	if err := r.db.Create(&p).Error; err != nil {
		return err
	}
	return nil
}

func (r *postRepository) FindById(id uint) (models.Post, error) {
	var post models.Post
	if err := r.db.Preload("Author").First(&post, id).Error; err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (r *postRepository) UpdateById(post *models.Post, newPost models.Post) error {
	if err := newPost.Validate("update"); err != nil {
		return err
	}

	if err := r.db.Model(&post).Updates(newPost).Error; err != nil {
		return err
	}

	return nil
}

func (r *postRepository) DeleteById(id uint) error {
	if err := r.db.Delete(&models.Post{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (r *postRepository) FindMany(limit int) ([]models.Post, error) {
	if limit == 0 {
		limit = 10
	}

	var posts []models.Post
	if err := r.db.
		Limit(limit).
		Order("created_at desc").
		Preload("Author").
		Where("is_published = ?", true).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func (r postRepository) FindPostsByUserId(uid uint) ([]models.Post, error) {
	var posts []models.Post

	if err := r.db.
		Order("created_at desc").
		Preload("Author").
		Where("author_id = ?", uid).
		Where("is_published = ?", true).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func (r postRepository) FindMyPosts(uid uint) ([]models.Post, error) {
	var posts []models.Post

	if err := r.db.
		Order("created_at desc").
		Preload("Author").
		Where("author_id = ?", uid).
		Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}
