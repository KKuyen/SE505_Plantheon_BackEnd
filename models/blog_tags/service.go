package blog_tags

import (
	"plantheon-backend/common"

	"gorm.io/gorm"
)

type BlogTagsService struct {
	db *gorm.DB
}

func NewBlogTagsService() *BlogTagsService {
	return &BlogTagsService{db: common.GetDB()}
}

func CreateBlogTag(tag *BlogTag) error {
	return NewBlogTagsService().db.Create(tag).Error
}

func UpdateBlogTag(tag *BlogTag) error {
	res := NewBlogTagsService().db.Model(&BlogTag{}).Where("id = ?", tag.ID).Updates(map[string]interface{}{
		"name": tag.Name,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func DeleteBlogTag(id string) error {
	return NewBlogTagsService().db.Delete(&BlogTag{}, "id = ?", id).Error
}

func GetBlogTagByID(id string) (*BlogTag, error) {
	var tag BlogTag
	err := NewBlogTagsService().db.Where("id = ?", id).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetAllBlogTags returns all blog tags
func GetAllBlogTags() ([]BlogTag, error) {
	var tags []BlogTag
	err := NewBlogTagsService().db.Order("created_at ASC").Find(&tags).Error
	return tags, err
}
