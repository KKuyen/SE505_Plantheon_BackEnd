package posts

import (
	"plantheon-backend/models/comments"
	"plantheon-backend/models/users"
	"time"
)



type CreatePostRequest struct {
	Content       string   `json:"content"`
	ImageLink     []string `json:"image_link"`
	DiseaseLink   *string  `json:"disease_link"`
	ScanHistoryID *string  `json:"scan_history_id"`
	Tags          []string `json:"tags" binding:"required"`
}

type PostResponse struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	FullName           string    `json:"full_name"`
	Avatar             string    `json:"avatar"`
	Content            string    `json:"content"`
	ImageLink          []string  `json:"image_link"`
	DiseaseLink        *string   `json:"disease_link"`
	DiseaseName        *string   `json:"disease_name"`
	DiseaseDescription *string   `json:"disease_description"`
	DiseaseSolution    *string   `json:"disease_solution"`
	DiseaseImageLink   []string  `json:"disease_image_link"`
	ScanHistoryID      *string   `json:"scan_history_id"`
	Tags               []string  `json:"tags"`
	LikeNum            int       `json:"like_number"`
	Liked              bool      `json:"liked"`
	IsMyPost           bool      `json:"is_my_post"`
	CommentNum         int       `json:"comment_number"`
	ShareNum           int       `json:"share_number"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type PostListResponse struct {
	Posts []PostResponse `json:"posts"`
	Total int            `json:"total"`
}
type PostDetailResponse struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	FullName           string    `json:"full_name"`
	Avatar             string    `json:"avatar"`
	Content            string    `json:"content"`
	ImageLink          []string  `json:"image_link"`
	DiseaseLink        *string   `json:"disease_link"`
	DiseaseName        *string   `json:"disease_name"`
	DiseaseDescription *string   `json:"disease_description"`
	DiseaseSolution    *string   `json:"disease_solution"`
	DiseaseImageLink   []string  `json:"disease_image_link"`
	ScanHistoryID      *string   `json:"scan_history_id"`
	Tags               []string  `json:"tags"`
	LikeNum            int       `json:"like_number"`
	Liked              bool      `json:"liked"`
	IsMyPost           bool      `json:"is_my_post"`
	CommentNum         int       `json:"comment_number"`
	CommentList        []comments.CommentResponse `json:"comment_list"`
	ShareNum           int       `json:"share_number"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// PublicUserPostsResponse bundles basic user info with that user's public posts.
// Useful for public endpoints where no authentication context is available.
type PublicUserPostsResponse struct {
	User  users.UserResponse `json:"user"`
	Posts PostListResponse   `json:"posts"`
}
