package posts

import (
	"fmt"

	"plantheon-backend/common"
	"plantheon-backend/models/comments"
	"plantheon-backend/models/noti"
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
	
	// Exclude viewer's own posts if viewerID is provided
	query := service.db
	if viewerID != "" {
		query = query.Where("user_id != ?", viewerID)
	}
	
	if err := query.Find(&posts).Error; err != nil {
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
			user.FullName = "Unknown User"
			user.Avatar = ""
		}

		// Lấy thông tin disease nếu có disease_link
		var diseaseName *string
		var diseaseDescription *string
		var diseaseSolution *string
		var diseaseImageLink []string
		if post.DiseaseLink != nil && *post.DiseaseLink != "" {
			var disease struct {
				Name        *string
				Description *string
				Solution    *string
				ImageLink   interface{}
			}
			if err := service.db.Table("diseases").Select("name, description, solution, image_link").Where("id = ?", *post.DiseaseLink).First(&disease).Error; err == nil {
				diseaseName = disease.Name
				diseaseDescription = disease.Description
				diseaseSolution = disease.Solution
				if imgLinks, ok := disease.ImageLink.([]interface{}); ok {
					for _, link := range imgLinks {
						if strLink, ok := link.(string); ok {
							diseaseImageLink = append(diseaseImageLink, strLink)
						}
					}
				}
			}
		}

		postResponses = append(postResponses, PostResponse{
			ID:                 post.ID,
			UserID:             post.UserID,
			FullName:           user.FullName,
			Avatar:             user.Avatar,
			Content:            post.Content,
			ImageLink:          post.ImageLink,
			DiseaseLink:        post.DiseaseLink,
			DiseaseName:        diseaseName,
			DiseaseDescription: diseaseDescription,
			DiseaseSolution:    diseaseSolution,
			DiseaseImageLink:   diseaseImageLink,
			ScanHistoryID:      post.ScanHistoryID,
			Tags:               post.Tags,
			LikeNum:            post.LikeNum,
			Liked:              likedMap[post.ID],
			IsMyPost:           viewerID != "" && post.UserID == viewerID,
			CommentNum:         post.CommentNum,
			ShareNum:           post.ShareNum,
			CreatedAt:          post.CreatedAt,
			UpdatedAt:          post.UpdatedAt,
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
		user.FullName = "Unknown User"
		user.Avatar = ""
	}

	// Lấy thông tin disease nếu có disease_link
	var diseaseName *string
	var diseaseDescription *string
	var diseaseSolution *string
	var diseaseImageLink []string
	if post.DiseaseLink != nil && *post.DiseaseLink != "" {
		var disease struct {
			Name        *string
			Description *string
			Solution    *string
			ImageLink   interface{}
		}
		if err := service.db.Table("diseases").Select("name, description, solution, image_link").Where("id = ?", *post.DiseaseLink).First(&disease).Error; err == nil {
			diseaseName = disease.Name
			diseaseDescription = disease.Description
			diseaseSolution = disease.Solution
			if imgLinks, ok := disease.ImageLink.([]interface{}); ok {
				for _, link := range imgLinks {
					if strLink, ok := link.(string); ok {
						diseaseImageLink = append(diseaseImageLink, strLink)
					}
				}
			}
		}
	}
	
	// Lấy danh sách bình luận phẳng (bao gồm parent_id, is_like)
	commentList, err := comments.GetCommentsByPostID(post.ID, viewerID)
	if err != nil {
		return nil, err
	}
	
	return &PostDetailResponse{
		ID:                 post.ID,
		UserID:             post.UserID,
		FullName:           user.FullName,
		Avatar:             user.Avatar,
		Content:            post.Content,
		ImageLink:          post.ImageLink,
		DiseaseLink:        post.DiseaseLink,
		DiseaseName:        diseaseName,
		DiseaseDescription: diseaseDescription,
		DiseaseSolution:    diseaseSolution,
		DiseaseImageLink:   diseaseImageLink,
		ScanHistoryID:      post.ScanHistoryID,
		Tags:               post.Tags,
		LikeNum:            post.LikeNum,
		Liked:              liked,
		IsMyPost:           viewerID != "" && post.UserID == viewerID,
		CommentNum:         post.CommentNum,
		ShareNum:           post.ShareNum,
		CreatedAt:          post.CreatedAt,
		UpdatedAt:          post.UpdatedAt,
		CommentList:        commentList,
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
	if err := service.db.Transaction(func(tx *gorm.DB) error {
		// Insert into post_likes table
		if err := tx.Exec("INSERT INTO post_likes (user_id, post_id, created_at) VALUES (?, ?, NOW())", userID, postID).Error; err != nil {
			return err
		}

		// Increment like_num in posts table
		if err := tx.Model(&Post{}).Where("id = ?", postID).UpdateColumn("like_num", gorm.Expr("like_num + ?", 1)).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	// After successful like, send notification to post owner (if not self)
	var post Post
	if err := service.db.Select("id, user_id").Where("id = ?", postID).First(&post).Error; err == nil {
		if post.UserID != userID {
			var liker users.User
			if err := service.db.Select("full_name").Where("id = ?", userID).First(&liker).Error; err != nil {
				liker.FullName = "Ai đó"
			}
			title := "Bài viết của bạn được thích"
			content := fmt.Sprintf("%s đã thích bài viết của bạn.", liker.FullName)
			_ = noti.CreatePostNotification(post.UserID, &post.ID, title, content)
		}
	}

	return nil
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
			user.FullName = "Unknown User"
			user.Avatar = ""
		}

		// Lấy thông tin disease nếu có disease_link
		var diseaseName *string
		var diseaseDescription *string
		var diseaseSolution *string
		var diseaseImageLink []string
		if post.DiseaseLink != nil && *post.DiseaseLink != "" {
			var disease struct {
				Name        *string
				Description *string
				Solution    *string
				ImageLink   interface{}
			}
			if err := service.db.Table("diseases").Select("name, description, solution, image_link").Where("id = ?", *post.DiseaseLink).First(&disease).Error; err == nil {
				diseaseName = disease.Name
				diseaseDescription = disease.Description
				diseaseSolution = disease.Solution
				if imgLinks, ok := disease.ImageLink.([]interface{}); ok {
					for _, link := range imgLinks {
						if strLink, ok := link.(string); ok {
							diseaseImageLink = append(diseaseImageLink, strLink)
						}
					}
				}
			}
		}

		postResponses = append(postResponses, PostResponse{
			ID:                 post.ID,
			UserID:             post.UserID,
			FullName:           user.FullName,
			Avatar:             user.Avatar,
			Content:            post.Content,
			ImageLink:          post.ImageLink,
			DiseaseLink:        post.DiseaseLink,
			DiseaseName:        diseaseName,
			DiseaseDescription: diseaseDescription,
			DiseaseSolution:    diseaseSolution,
			DiseaseImageLink:   diseaseImageLink,
			ScanHistoryID:      post.ScanHistoryID,
			Tags:               post.Tags,
			LikeNum:            post.LikeNum,
			Liked:              likedMap[post.ID],
			IsMyPost:           viewerID != "" && post.UserID == viewerID,
			CommentNum:         post.CommentNum,
			ShareNum:           post.ShareNum,
			CreatedAt:          post.CreatedAt,
			UpdatedAt:          post.UpdatedAt,
		})
	}

	return PostListResponse{
		Posts: postResponses,
		Total: len(postResponses),
	}, nil
}

// GetPublicPostsByUserID returns only posts that are publicly visible for a user.
// A post is treated as public when its status is AVAILABLE or not set.
func GetPublicPostsByUserID(userID string, viewerID string) (PostListResponse, error) {
	service := NewPostsService()
	var posts []Post

	if err := service.db.Where("user_id = ?", userID).
		Where("status IS NULL OR status = '' OR status = ?", "AVAILABLE").
		Find(&posts).Error; err != nil {
		return PostListResponse{}, err
	}

	// Batch-fetch likes for all posts by viewerID to avoid N+1 lookups
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
		var user users.User
		if err := service.db.Where("id = ?", post.UserID).First(&user).Error; err != nil {
			user.FullName = "Unknown User"
			user.Avatar = ""
		}

		var diseaseName *string
		var diseaseDescription *string
		var diseaseSolution *string
		var diseaseImageLink []string
		if post.DiseaseLink != nil && *post.DiseaseLink != "" {
			var disease struct {
				Name        *string
				Description *string
				Solution    *string
				ImageLink   interface{}
			}
			if err := service.db.Table("diseases").Select("name, description, solution, image_link").Where("id = ?", *post.DiseaseLink).First(&disease).Error; err == nil {
				diseaseName = disease.Name
				diseaseDescription = disease.Description
				diseaseSolution = disease.Solution
				if imgLinks, ok := disease.ImageLink.([]interface{}); ok {
					for _, link := range imgLinks {
						if strLink, ok := link.(string); ok {
							diseaseImageLink = append(diseaseImageLink, strLink)
						}
					}
				}
			}
		}

		postResponses = append(postResponses, PostResponse{
			ID:                 post.ID,
			UserID:             post.UserID,
			FullName:           user.FullName,
			Avatar:             user.Avatar,
			Content:            post.Content,
			ImageLink:          post.ImageLink,
			DiseaseLink:        post.DiseaseLink,
			DiseaseName:        diseaseName,
			DiseaseDescription: diseaseDescription,
			DiseaseSolution:    diseaseSolution,
			DiseaseImageLink:   diseaseImageLink,
			ScanHistoryID:      post.ScanHistoryID,
			Tags:               post.Tags,
			LikeNum:            post.LikeNum,
			Liked:              likedMap[post.ID],
			IsMyPost:           viewerID != "" && post.UserID == viewerID,
			CommentNum:         post.CommentNum,
			ShareNum:           post.ShareNum,
			CreatedAt:          post.CreatedAt,
			UpdatedAt:          post.UpdatedAt,
		})
	}

	return PostListResponse{
		Posts: postResponses,
		Total: len(postResponses),
	}, nil
}
