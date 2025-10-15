package posts

type CreatePostRequest struct {
	Content   string   `json:"content"`
	ImageLink []string `json:"image_link"`
	Tags      []string `json:"tags" binding:"required"`
}