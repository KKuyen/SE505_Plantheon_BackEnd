package noti

import (
	"time"
)

type NotificationResponse struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	IsRead    bool       `json:"is_read"`
	PostID    *string    `json:"post_id"`
	CreatedAt time.Time  `json:"created_at"`
}

type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	Total         int                    `json:"total"`
}

type MarkNotificationAsSeenRequest struct {
	IsRead bool `json:"is_read"`
}

type CreateNotificationRequest struct {
	Title   string  `json:"title" binding:"required"`
	Content string  `json:"content" binding:"required"`
	PostID  *string `json:"post_id"`
}

