package comments

import "time"

type Comments struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	PostID    string    `json:"post_id" gorm:"not null;type:uuid"`
	UserID    string    `json:"user_id" gorm:"not null;type:uuid"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	LikeNum   int64     `json:"like_number" gorm:"default:0"`
	Status    string    `json:"status" gorm:"type:varchar"` // REPORTED AVAILABLE UNAVAILABLE
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}