package blogs

import (
	"plantheon-backend/common"
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
		"content":         blog.Content,
		"cover_image_url": blog.CoverImageURL,
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

func GetAllNews() (NewsListResponse, error) {
	service := NewBlogsService()
	var blogs []Blog
	if err := service.db.Where("status = ?", "published").Find(&blogs).Error; err != nil {
		return NewsListResponse{}, err
	}

	var newsResponses []NewsResponse
	for _, blog := range blogs {
		// Lấy thông tin user từ UserID
		var user users.User
		if err := service.db.Where("id = ?", blog.UserID).First(&user).Error; err != nil {
			// Nếu không tìm thấy user, vẫn trả về blog nhưng với thông tin user rỗng
			newsResponses = append(newsResponses, NewsResponse{
				ID:           blog.ID,
				Title:        blog.Title,
				Content:     blog.Content,
				CoverImageURL: blog.CoverImageURL,
				Status:       blog.Status,
				PublishedAt:  blog.PublishedAt,
				CreatedAt:    blog.CreatedAt,
				UpdatedAt:    blog.UpdatedAt,
				UserID:       blog.UserID,
				FullName:     "Unknown User",
				Avatar:       "",
			})
			continue
		}

		newsResponses = append(newsResponses, NewsResponse{
			ID:           blog.ID,
			Title:        blog.Title,
			Content:      blog.Content,
			CoverImageURL: blog.CoverImageURL,
			Status:       blog.Status,
			PublishedAt:  blog.PublishedAt,
			CreatedAt:    blog.CreatedAt,
			UpdatedAt:    blog.UpdatedAt,
			UserID:       blog.UserID,
			FullName:     user.FullName,
			Avatar:       user.Avatar,
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
			ID:           blog.ID,
			Title:        blog.Title,
			Content:      blog.Content,
			CoverImageURL: blog.CoverImageURL,
			Status:       blog.Status,
			PublishedAt:  blog.PublishedAt,
			CreatedAt:    blog.CreatedAt,
			UpdatedAt:    blog.UpdatedAt,
			UserID:       blog.UserID,
			FullName:     "Unknown User",
			Avatar:       "",
		}, nil
	}

	return &NewsDetailResponse{
		ID:           blog.ID,
		Title:        blog.Title,
		Content:      blog.Content,
		CoverImageURL: blog.CoverImageURL,
		Status:       blog.Status,
		PublishedAt:  blog.PublishedAt,
		CreatedAt:    blog.CreatedAt,
		UpdatedAt:    blog.UpdatedAt,
		UserID:       blog.UserID,
		FullName:     user.FullName,
		Avatar:       user.Avatar,
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
				ID:           blog.ID,
				Title:        blog.Title,
				Content:      blog.Content,
				CoverImageURL: blog.CoverImageURL,
				Status:       blog.Status,
				PublishedAt:  blog.PublishedAt,
				CreatedAt:    blog.CreatedAt,
				UpdatedAt:    blog.UpdatedAt,
				UserID:       blog.UserID,
				FullName:     "Unknown User",
				Avatar:       "",
			})
			continue
		}

		newsResponses = append(newsResponses, NewsResponse{
			ID:           blog.ID,
			Title:        blog.Title,
			Content:      blog.Content,
			CoverImageURL: blog.CoverImageURL,
			Status:       blog.Status,
			PublishedAt:  blog.PublishedAt,
			CreatedAt:    blog.CreatedAt,
			UpdatedAt:    blog.UpdatedAt,
			UserID:       blog.UserID,
			FullName:     user.FullName,
			Avatar:       user.Avatar,
		})
	}

	return NewsListResponse{
		News:  newsResponses,
		Total: len(newsResponses),
	}, nil
}
