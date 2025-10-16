package posts

import (
	"plantheon-backend/common"
	"plantheon-backend/models/comments"
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

func UpdatePost(post *Post) error {	
	service := NewPostsService()
	if err := service.db.Save(post).Error; err != nil {
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

func GetPostByID(id string) (*PostDetailResponse, error) {
	service := NewPostsService()
	var post Post
	if err := service.db.Where("id = ?", id).First(&post).Error; err != nil {
		return nil, err
	}

	// Lấy thông tin user từ UserID
	var user users.User
	if err := service.db.Where("id = ?", post.UserID).First(&user).Error; err != nil {
		// Nếu không tìm thấy user, trả về post nhưng với thông tin user rỗng
		return &PostDetailResponse{
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
			CommentList: []comments.CommentResponse{},
		}, nil
	}
	// Lấy danh sách bình luận liên quan đến bài viết
	var commentList []comments.CommentResponse
	if err := service.db.Table("comments").Where("post_id = ?", post.ID).Find(&commentList).Error; err != nil {
		return nil, err
	}
	return &PostDetailResponse{
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
		CommentList: commentList,
	}, nil
}

func DeletePostByID(id string) error {
	service := NewPostsService()
	if err := service.db.Delete(&Post{}, "id = ?", id).Error; err != nil {
		return err	
	}
	return nil
}

func LikePost(id string) error {
	service := NewPostsService()
	if err := service.db.Model(&Post{}).Where("id = ?", id).UpdateColumn("like_num", gorm.Expr("like_num + ?", 1)).Error; err != nil {
		return err	
	}
	return nil
}

func UnlikePost(id string) error {
	service := NewPostsService()
	if err := service.db.Model(&Post{}).Where("id = ? AND like_num > 0", id).UpdateColumn("like_num", gorm.Expr("like_num - ?", 1)).Error; err != nil {
		return err
	}
	return nil
}

func SharePost(id string) error {
	service := NewPostsService()
	if err := service.db.Model(&Post{}).Where("id = ?", id).UpdateColumn("share_num", gorm.Expr("share_num + ?", 1)).Error; err != nil {
		return err
	}
	return nil
}
