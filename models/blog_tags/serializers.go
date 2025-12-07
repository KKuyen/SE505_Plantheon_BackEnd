package blog_tags

import "time"

// BlogTagRequest represents create/update payload
type BlogTagRequest struct {
	Name string `json:"name" binding:"required"`
}

// BlogTagResponse represents tag data returned to clients
type BlogTagResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToBlogTagResponse converts model to response
func (bt *BlogTag) ToBlogTagResponse() BlogTagResponse {
	return BlogTagResponse{
		ID:        bt.ID,
		Name:      bt.Name,
		CreatedAt: bt.CreatedAt,
		UpdatedAt: bt.UpdatedAt,
	}
}
