package posts

import (
	"plantheon-backend/common"
	"plantheon-backend/models/users"

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

func GetAllPosts() (PostListResponse, error) {
	service := NewPostsService()
	var posts []Post
	if err := service.db.Find(&posts).Error; err != nil {
		return PostListResponse{}, err
	}
	
	var postResponses []PostResponse
	for _, post := range posts {
		// Lấy thông tin user từ UserID
		var user users.User
		if err := service.db.Where("id = ?", post.UserID).First(&user).Error; err != nil {
			// Nếu không tìm thấy user, vẫn trả về post nhưng với thông tin user rỗng
			postResponses = append(postResponses, PostResponse{
				ID:         post.ID,
				UserID:     post.UserID,
				FullName:   "Unknown User",
				Avatar:     "",
				Content:    post.Content,
				Tags:       post.Tags,
				LikeNum:    post.LikeNum,
				CommentNum: post.CommentNum,
				ShareNum:   post.ShareNum,
				CreatedAt:  post.CreatedAt,
				UpdatedAt:  post.UpdatedAt,
			})
			continue
		}
		
		postResponses = append(postResponses, PostResponse{
			ID:         post.ID,
			UserID:     post.UserID,
			FullName:   user.FullName,
			Avatar:     user.Avatar,
			Content:    post.Content,
			Tags:       post.Tags,
			LikeNum:    post.LikeNum,
			CommentNum: post.CommentNum,
			ShareNum:   post.ShareNum,
			CreatedAt:  post.CreatedAt,
			UpdatedAt:  post.UpdatedAt,
		})
	}
	
	return PostListResponse{
		Posts: postResponses,
		Total: len(postResponses),
	}, nil
}

func GetPostByID(id string) (*PostResponse, error) {
	service := NewPostsService()
	var post Post
	if err := service.db.Where("id = ?", id).First(&post).Error; err != nil {
		return nil, err
	}

	// Lấy thông tin user từ UserID
	var user users.User
	if err := service.db.Where("id = ?", post.UserID).First(&user).Error; err != nil {
		// Nếu không tìm thấy user, trả về post nhưng với thông tin user rỗng
		return &PostResponse{
			ID:         post.ID,
			UserID:     post.UserID,
			FullName:   "Unknown User",
			Avatar:     "",
			Content:    post.Content,
			Tags:       post.Tags,
			LikeNum:    post.LikeNum,
			CommentNum: post.CommentNum,
			ShareNum:   post.ShareNum,
			CreatedAt:  post.CreatedAt,
			UpdatedAt:  post.UpdatedAt,
		}, nil
	}

	return &PostResponse{
		ID:         post.ID,
		UserID:     post.UserID,
		FullName:   user.FullName,
		Avatar:     user.Avatar,
		Content:    post.Content,
		Tags:       post.Tags,
		LikeNum:    post.LikeNum,
		CommentNum: post.CommentNum,
		ShareNum:   post.ShareNum,
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
	}, nil
}

