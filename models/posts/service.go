package posts

import (
	"plantheon-backend/common"

	"gorm.io/gorm"
)

type PostsService struct {
	db *gorm.DB
}

func NewPostsService() *PostsService {
	return &PostsService{
		db: common.GetDB(),
	}
}
func CreatePost(post *Post) error {
	service := NewPostsService()
	if err := service.db.Create(post).Error; err != nil {
		return err
	}
	return nil
}