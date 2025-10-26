package blogs

import (
	"time"
)

type CreateNewsRequest struct {
	Title         string  `json:"title" binding:"required"`
	Content       string  `json:"content" binding:"required"`
	CoverImageURL *string `json:"cover_image_url"`
	Status        string  `json:"status"` // draft or published
}

type NewsResponse struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	CoverImageURL *string    `json:"cover_image_url"`
	Status        string     `json:"status"`
	PublishedAt   *time.Time `json:"published_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UserID        string     `json:"user_id"`
	FullName      string     `json:"full_name"`
	Avatar        string     `json:"avatar"`
}

type NewsListResponse struct {
	News  []NewsResponse `json:"news"`
	Total int            `json:"total"`
}

type NewsDetailResponse struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	CoverImageURL *string    `json:"cover_image_url"`
	Status        string     `json:"status"`
	PublishedAt   *time.Time `json:"published_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UserID        string     `json:"user_id"`
	FullName      string     `json:"full_name"`
	Avatar        string     `json:"avatar"`
}
