package comments

import "time"

type CommentResponse struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	FullName  string    `json:"full_name"`
	Avatar    string    `json:"avatar"`
	Content   string    `json:"content"`
	LikeNum   int       `json:"like_number"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}