package blogs

import (
	"plantheon-backend/common"
	"plantheon-backend/models/blog_tags"
	"plantheon-backend/models/users"

	"gorm.io/gorm"
)

type BlogsService struct {
	db *gorm.DB
}

func NewBlogsService() *BlogsService {
	return &BlogsService{
		db: common.GetDB(),
	}
}

func CreateNews(blog *Blog) error {
	service := NewBlogsService()
	if err := service.db.Create(blog).Error; err != nil {
		return err
	}
	return nil
}

func UpdateNews(blog *Blog) error {
	service := NewBlogsService()
	// Use Updates to only update the fields that are provided
	updates := map[string]interface{}{
		"title":           blog.Title,
		"description":     blog.Description,
		"content":         blog.Content,
		"cover_image_url": blog.CoverImageURL,
		"blog_tag_id":     blog.BlogTagID,
		"status":          blog.Status,
	}

	// Only update published_at if it's being set (not nil)
	if blog.PublishedAt != nil {
		updates["published_at"] = blog.PublishedAt
	}

	if err := service.db.Model(&Blog{}).Where("id = ?", blog.ID).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

func GetAllNews(size *int) (NewsListResponse, error) {
	service := NewBlogsService()
	var blogs []Blog
	// Only fetch published blogs that are NOT linked to any sub guide stage
	query := service.db.Where("status = ? AND sub_guide_stages_id IS NULL", "published").Order("created_at DESC")
	if size != nil {
		query = query.Limit(*size)
	}
	if err := query.Find(&blogs).Error; err != nil {
		return NewsListResponse{}, err
	}

	var newsResponses []NewsResponse
	for _, blog := range blogs {
		// Lấy thông tin user từ UserID
		var user users.User
		if err := service.db.Where("id = ?", blog.UserID).First(&user).Error; err != nil {
			// Nếu không tìm thấy user, vẫn trả về blog nhưng với thông tin user rỗng
			newsResponses = append(newsResponses, NewsResponse{
				ID:            blog.ID,
				Title:         blog.Title,
				Description:   blog.Description,
				BlogTagID:     blog.BlogTagID,
				BlogTagName:   resolveBlogTagName(service.db, blog.BlogTagID),
				CoverImageURL: blog.CoverImageURL,
				Status:        blog.Status,
				PublishedAt:   blog.PublishedAt,
				CreatedAt:     blog.CreatedAt,
				UpdatedAt:     blog.UpdatedAt,
				UserID:        blog.UserID,
				FullName:      "Unknown User",
				Avatar:        "",
			})
			continue
		}

		newsResponses = append(newsResponses, NewsResponse{
			ID:            blog.ID,
			Title:         blog.Title,
			Description:   blog.Description,
			BlogTagID:     blog.BlogTagID,
			BlogTagName:   resolveBlogTagName(service.db, blog.BlogTagID),
			CoverImageURL: blog.CoverImageURL,
			Status:        blog.Status,
			PublishedAt:   blog.PublishedAt,
			CreatedAt:     blog.CreatedAt,
			UpdatedAt:     blog.UpdatedAt,
			UserID:        blog.UserID,
			FullName:      user.FullName,
			Avatar:        user.Avatar,
		})
	}

	return NewsListResponse{
		News:  newsResponses,
		Total: len(newsResponses),
	}, nil
}

func GetNewsByID(id string) (*NewsDetailResponse, error) {
	service := NewBlogsService()
	var blog Blog
	if err := service.db.Where("id = ?", id).First(&blog).Error; err != nil {
		return nil, err
	}

	// Lấy thông tin user từ UserID
	var user users.User
	if err := service.db.Where("id = ?", blog.UserID).First(&user).Error; err != nil {
		// Nếu không tìm thấy user, trả về blog nhưng với thông tin user rỗng
		return &NewsDetailResponse{
			ID:            blog.ID,
			Title:         blog.Title,
			Description:   blog.Description,
			Content:       blog.Content,
			BlogTagID:     blog.BlogTagID,
			BlogTagName:   resolveBlogTagName(service.db, blog.BlogTagID),
			CoverImageURL: blog.CoverImageURL,
			Status:        blog.Status,
			PublishedAt:   blog.PublishedAt,
			CreatedAt:     blog.CreatedAt,
			UpdatedAt:     blog.UpdatedAt,
			UserID:        blog.UserID,
			FullName:      "Unknown User",
			Avatar:        "",
		}, nil
	}

	return &NewsDetailResponse{
		ID:            blog.ID,
		Title:         blog.Title,
		Description:   blog.Description,
		Content:       blog.Content,
		BlogTagID:     blog.BlogTagID,
		BlogTagName:   resolveBlogTagName(service.db, blog.BlogTagID),
		CoverImageURL: blog.CoverImageURL,
		Status:        blog.Status,
		PublishedAt:   blog.PublishedAt,
		CreatedAt:     blog.CreatedAt,
		UpdatedAt:     blog.UpdatedAt,
		UserID:        blog.UserID,
		FullName:      user.FullName,
		Avatar:        user.Avatar,
	}, nil
}

func DeleteNewsByID(id string) error {
	service := NewBlogsService()
	if err := service.db.Delete(&Blog{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func GetNewsByUserID(userID string) (NewsListResponse, error) {
	service := NewBlogsService()
	var blogs []Blog
	if err := service.db.Where("user_id = ? AND status = ?", userID, "published").Find(&blogs).Error; err != nil {
		return NewsListResponse{}, err
	}

	var newsResponses []NewsResponse
	for _, blog := range blogs {
		// Lấy thông tin user từ UserID
		var user users.User
		if err := service.db.Where("id = ?", blog.UserID).First(&user).Error; err != nil {
			newsResponses = append(newsResponses, NewsResponse{
				ID:            blog.ID,
				Title:         blog.Title,
				Description:   blog.Description,
				BlogTagID:     blog.BlogTagID,
				BlogTagName:   resolveBlogTagName(service.db, blog.BlogTagID),
				CoverImageURL: blog.CoverImageURL,
				Status:        blog.Status,
				PublishedAt:   blog.PublishedAt,
				CreatedAt:     blog.CreatedAt,
				UpdatedAt:     blog.UpdatedAt,
				UserID:        blog.UserID,
				FullName:      "Unknown User",
				Avatar:        "",
			})
			continue
		}

		newsResponses = append(newsResponses, NewsResponse{
			ID:            blog.ID,
			Title:         blog.Title,
			Description:   blog.Description,
			BlogTagID:     blog.BlogTagID,
			BlogTagName:   resolveBlogTagName(service.db, blog.BlogTagID),
			CoverImageURL: blog.CoverImageURL,
			Status:        blog.Status,
			PublishedAt:   blog.PublishedAt,
			CreatedAt:     blog.CreatedAt,
			UpdatedAt:     blog.UpdatedAt,
			UserID:        blog.UserID,
			FullName:      user.FullName,
			Avatar:        user.Avatar,
		})
	}

	return NewsListResponse{
		News:  newsResponses,
		Total: len(newsResponses),
	}, nil
}

// resolveBlogTagName safely fetches tag name for given tag id using provided db.
func resolveBlogTagName(db *gorm.DB, tagID *string) *string {
	if tagID == nil || *tagID == "" {
		return nil
	}
	var tag blog_tags.BlogTag
	if err := db.Select("id", "name").Where("id = ?", *tagID).First(&tag).Error; err != nil {
		return nil
	}
	return &tag.Name
}

// GetBlogsBySubGuideStageID returns published blogs linked to a sub guide stage.
func GetBlogsBySubGuideStageID(subGuideStageID string) ([]Blog, error) {
	service := NewBlogsService()
	var blogs []Blog
	if err := service.db.
		Where("sub_guide_stages_id = ? AND status = ?", subGuideStageID, "published").
		Order("created_at ASC").
		Find(&blogs).Error; err != nil {
		return nil, err
	}
	return blogs, nil
}
