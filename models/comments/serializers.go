package comments

import "time"

type CommentResponse struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	ParentID  string    `json:"parent_id"`
	FullName  string    `json:"full_name"`
	Avatar    string    `json:"avatar"`
	Content   string    `json:"content"`
	LikeNum   int64     `json:"like_number"`
	IsLike    bool      `json:"is_like"`
	IsMe      bool      `json:"is_me"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CommentGroupResponse groups a root comment and its replies.
type CommentGroupResponse struct {
	Parent  CommentResponse   `json:"parent"`
	Replies []CommentResponse `json:"replies"`
}

type CreateCommentRequest struct {
	Content  string  `json:"content" binding:"required"`
	ParentID *string `json:"parent_id"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}