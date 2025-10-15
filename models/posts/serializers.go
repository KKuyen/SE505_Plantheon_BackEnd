package posts

import (
	"plantheon-backend/models/comments"
	"time"
)



type CreatePostRequest struct {
	Content   string   `json:"content"`
	ImageLink []string `json:"image_link"`
	Tags      []string `json:"tags" binding:"required"`
}

type PostResponse struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	FullName   string    `json:"full_name"`
	Avatar     string    `json:"avatar"`
	Content    string    `json:"content"`
	Tags       []string  `json:"tags"`
	LikeNum    int       `json:"like_number"`
	CommentNum int       `json:"comment_number"`
	ShareNum   int       `json:"share_number"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type PostListResponse struct {
	Posts []PostResponse `json:"posts"`
	Total int            `json:"total"`
}
type PostDetailResponse struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	FullName   string    `json:"full_name"`
	Avatar     string    `json:"avatar"`
	Content    string    `json:"content"`
	Tags       []string  `json:"tags"`
	LikeNum    int       `json:"like_number"`
	CommentNum int       `json:"comment_number"`
	CommentList []comments.CommentResponse `json:"comment_list"`
	ShareNum   int       `json:"share_number"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
