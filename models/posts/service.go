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

func GetAllPosts(viewerID string) (PostListResponse, error) {
	service := NewPostsService()
	var posts []Post
	if err := service.db.Find(&posts).Error; err != nil {
		return PostListResponse{}, err
	}
	
	// Batch-fetch likes for all posts by viewerID to avoid N+1
	postIDs := make([]string, len(posts))
	for i, p := range posts {
		postIDs[i] = p.ID
	}
	var likedPostIDs []string
	if viewerID != "" && len(postIDs) > 0 {
		type LikeResult struct {
			PostID string
		}
		var likeResults []LikeResult
		if err := service.db.Table("post_likes").Select("post_id").Where("user_id = ? AND post_id IN ?", viewerID, postIDs).Find(&likeResults).Error; err == nil {
			for _, lr := range likeResults {
				likedPostIDs = append(likedPostIDs, lr.PostID)
			}
		}
	}
	likedMap := make(map[string]bool)
	for _, pid := range likedPostIDs {
		likedMap[pid] = true
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
				ImageLink:  post.ImageLink,
				Tags:       post.Tags,
				LikeNum:    post.LikeNum,
				Liked:      likedMap[post.ID],
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
			ImageLink:  post.ImageLink,
			Tags:       post.Tags,
			LikeNum:    post.LikeNum,
			Liked:      likedMap[post.ID],
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

func GetPostByID(id string, viewerID string) (*PostDetailResponse, error) {
	service := NewPostsService()
	var post Post
	if err := service.db.Where("id = ?", id).First(&post).Error; err != nil {
		return nil, err
	}

	// Check if viewer liked this post
	liked := false
	if viewerID != "" {
		var count int64
		service.db.Table("post_likes").Where("user_id = ? AND post_id = ?", viewerID, post.ID).Count(&count)
		liked = count > 0
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
			ImageLink:  post.ImageLink,
			Tags:       post.Tags,
			LikeNum:    post.LikeNum,
			Liked:      liked,
			CommentNum: post.CommentNum,
			ShareNum:   post.ShareNum,
			CreatedAt:  post.CreatedAt,
			UpdatedAt:  post.UpdatedAt,
			CommentList: []comments.CommentResponse{},
		}, nil
	}
	// Lấy danh sách bình luận liên quan đến bài viết với thông tin user
	var rawComments []comments.Comments
	if err := service.db.Where("post_id = ?", post.ID).Order("created_at DESC").Find(&rawComments).Error; err != nil {
		return nil, err
	}
	
	var commentList []comments.CommentResponse
	for _, comment := range rawComments {
		// Query user info for each comment
		var userInfo struct {
			FullName string
			Avatar   string
		}
		if err := service.db.Table("users").Select("full_name, avatar").Where("id = ?", comment.UserID).First(&userInfo).Error; err != nil {
			// If user not found, use default values
			userInfo.FullName = "Unknown User"
			userInfo.Avatar = ""
		}
		
		commentList = append(commentList, comments.CommentResponse{
			ID:        comment.ID,
			PostID:    comment.PostID,
			UserID:    comment.UserID,
			FullName:  userInfo.FullName,
			Avatar:    userInfo.Avatar,
			Content:   comment.Content,
			LikeNum:   comment.LikeNum,
			IsMe:      viewerID != "" && comment.UserID == viewerID,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		})
	}
	
	return &PostDetailResponse{
		ID:         post.ID,
		UserID:     post.UserID,
		FullName:   user.FullName,
		Avatar:     user.Avatar,
		Content:    post.Content,
		ImageLink:  post.ImageLink,
		Tags:       post.Tags,
		LikeNum:    post.LikeNum,
		Liked:      liked,
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

func LikePost(postID string, userID string) error {
	service := NewPostsService()
	
	// Check if like already exists to prevent duplicate
	var count int64
	if err := service.db.Table("post_likes").Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil // Already liked, do nothing
	}
	
	// Use transaction to ensure both operations succeed or fail together
	return service.db.Transaction(func(tx *gorm.DB) error {
		// Insert into post_likes table
		if err := tx.Exec("INSERT INTO post_likes (user_id, post_id, created_at) VALUES (?, ?, NOW())", userID, postID).Error; err != nil {
			return err
		}
		
		// Increment like_num in posts table
		if err := tx.Model(&Post{}).Where("id = ?", postID).UpdateColumn("like_num", gorm.Expr("like_num + ?", 1)).Error; err != nil {
			return err
		}
		
		return nil
	})
}

func UnlikePost(postID string, userID string) error {
	service := NewPostsService()
	
	// Check if like exists
	var count int64
	if err := service.db.Table("post_likes").Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return nil // Not liked, do nothing
	}
	
	// Use transaction to ensure both operations succeed or fail together
	return service.db.Transaction(func(tx *gorm.DB) error {
		// Delete from post_likes table
		if err := tx.Exec("DELETE FROM post_likes WHERE user_id = ? AND post_id = ?", userID, postID).Error; err != nil {
			return err
		}
		
		// Decrement like_num in posts table
		if err := tx.Model(&Post{}).Where("id = ? AND like_num > 0", postID).UpdateColumn("like_num", gorm.Expr("like_num - ?", 1)).Error; err != nil {
			return err
		}
		
		return nil
	})
}

func SharePost(id string) error {
	service := NewPostsService()
	if err := service.db.Model(&Post{}).Where("id = ?", id).UpdateColumn("share_num", gorm.Expr("share_num + ?", 1)).Error; err != nil {
		return err
	}
	return nil
}

func GetPostsByUserID(userID string, viewerID string) (PostListResponse, error) {
	service := NewPostsService()
	var posts []Post
	if err := service.db.Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		return PostListResponse{}, err
	}
	
	// Batch-fetch likes for all posts by viewerID to avoid N+1
	postIDs := make([]string, len(posts))
	for i, p := range posts {
		postIDs[i] = p.ID
	}
	var likedPostIDs []string
	if viewerID != "" {
		type LikeResult struct {
			PostID string
		}
		var likeResults []LikeResult
		if err := service.db.Table("post_likes").Select("post_id").Where("user_id = ? AND post_id IN ?", viewerID, postIDs).Find(&likeResults).Error; err == nil {
			for _, lr := range likeResults {
				likedPostIDs = append(likedPostIDs, lr.PostID)
			}
		}
	}
	likedMap := make(map[string]bool)
	for _, pid := range likedPostIDs {
		likedMap[pid] = true
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
				ImageLink:  post.ImageLink,
				Tags:       post.Tags,
				LikeNum:    post.LikeNum,
				Liked:      likedMap[post.ID],
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
			ImageLink:  post.ImageLink,
			Tags:       post.Tags,
			LikeNum:    post.LikeNum,
			Liked:      likedMap[post.ID],
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