package blogs

import (
	"time"
)

type CreateNewsRequest struct {
	Title         string  `json:"title" binding:"required"`
	Description   *string `json:"description"`
	Content       string  `json:"content" binding:"required"`
	CoverImageURL *string `json:"cover_image_url"`
	BlogTagID     *string `json:"blog_tag_id"`
	Status        string  `json:"status"` // draft or published
}

type NewsResponse struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Description   *string    `json:"description"`
	BlogTagID     *string    `json:"blog_tag_id"`
	BlogTagName   *string    `json:"blog_tag_name,omitempty"`
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
	Description   *string    `json:"description"`
	Content       string     `json:"content"`
	BlogTagID     *string    `json:"blog_tag_id"`
	BlogTagName   *string    `json:"blog_tag_name,omitempty"`
	CoverImageURL *string    `json:"cover_image_url"`
	Status        string     `json:"status"`
	PublishedAt   *time.Time `json:"published_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UserID        string     `json:"user_id"`
	FullName      string     `json:"full_name"`
	Avatar        string     `json:"avatar"`
}

// BlogSummaryResponse is a lightweight blog representation
// used for embedding blogs under sub guide stages.
type BlogSummaryResponse struct {
	Title         string  `json:"title"`
	Description   *string `json:"description,omitempty"`
	Content       string  `json:"content"`
	BlogTagID     *string `json:"blog_tag_id,omitempty"`
	BlogTagName   *string `json:"blog_tag_name,omitempty"`
	CoverImageURL *string `json:"cover_image_url"`
}

// ToBlogSummaryResponse converts Blog to BlogSummaryResponse.
func (b *Blog) ToBlogSummaryResponse() BlogSummaryResponse {
	return BlogSummaryResponse{
		Title:         b.Title,
		Description:   b.Description,
		Content:       b.Content,
		BlogTagID:     b.BlogTagID,
		CoverImageURL: b.CoverImageURL,
	}
}
