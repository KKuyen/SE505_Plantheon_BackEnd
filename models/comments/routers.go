package comments

import (
	"net/http"
	"plantheon-backend/common"
	"plantheon-backend/models/noti"
	"plantheon-backend/models/users"
	"strings"

	"github.com/gin-gonic/gin"
)

func AddCommentHandler(c *gin.Context) {
	// Implementation of adding a comment to a post
	postID := c.Param("id")
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ValidateCreateCommentRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Không tìm thấy người dùng",
		})
		return
	}

	user, ok := userInterface.(*users.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Định dạng người dùng không hợp lệ",
		})
		return
	}

	// Validate parent if provided
	if req.ParentID != nil && *req.ParentID != "" {
		if err := ValidateParentComment(common.GetDB(), *req.ParentID, postID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	comment := &Comments{
		PostID:  postID,
		Content: req.Content,
		UserID:  user.ID, // Sử dụng user.ID từ JWT token
	}

	if err := AddComment(comment, req.ParentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Tăng số comment cho post
	if err := IncreasePostCommentCount(postID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tăng số bình luận thất bại"})
		return
	}

	// Gửi thông báo cho chủ bài viết (nếu không phải chính họ)
	var postOwner struct{ UserID string }
	db := common.GetDB()
	if err := db.Table("posts").Select("user_id").Where("id = ?", postID).First(&postOwner).Error; err == nil {
		if postOwner.UserID != user.ID {
			snippet := strings.TrimSpace(req.Content)
			if len(snippet) > 100 {
				snippet = snippet[:100] + "..."
			}
			title := "Bài viết của bạn có bình luận mới"
			content := user.FullName + " đã bình luận: " + snippet
			_ = noti.CreatePostNotification(postOwner.UserID, &comment.PostID, title, content)
		}
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Thêm bình luận thành công",
		"data": gin.H{
			"id":        comment.ID,
			"post_id":   comment.PostID,
			"user_id":   comment.UserID,
			"full_name": user.FullName, // Thêm thông tin user
			"avatar":    user.Avatar,   // Thêm thông tin user
			"content":   comment.Content,
			"parent_id": comment.ParentID,
			"created_at": comment.CreatedAt,
		},
	})
}

func UpdateCommentHandler(c *gin.Context) {
	commentID := c.Param("id")
	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ValidateUpdateCommentRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Không tìm thấy người dùng",
		})
		return
	}
	user, ok := userInterface.(*users.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Định dạng người dùng không hợp lệ",
		})
		return
	}
	comment, err := GetCommentByID(commentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy bình luận"})
		return
	}
	if comment.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bạn không có quyền cập nhật bình luận này"})
		return
	}
	comment.Content = req.Content
	if err := UpdateComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cập nhật bình luận thành công"})
}

func DeleteCommentHandler(c *gin.Context) {
	commentID := c.Param("commentId")
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Không tìm thấy người dùng",
		})
		return
	}
	user, ok := userInterface.(*users.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Định dạng người dùng không hợp lệ",
		})
		return
	}
	comment, err := GetCommentByID(commentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy bình luận"})
		return
	}
	if comment.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bạn không có quyền xóa bình luận này"})
		return
	}
	if err := DeleteCommentByID(commentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Xóa bình luận thành công"})
}

func GetCommentsByPostIDHandler(c *gin.Context) {
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cần có ID bài viết"})
		return
	}

	// Extract viewer's user ID from JWT token
	viewerID := ""
	userInterface, exists := c.Get("user")
	if exists {
		user, ok := userInterface.(*users.User)
		if ok {
			viewerID = user.ID
		}
	}

	comments, err := GetCommentsByPostID(postID, viewerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lấy danh sách bình luận thất bại"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  comments,
		"total": len(comments),
	})
}

// LikeCommentHandler handles liking a comment.
func LikeCommentHandler(c *gin.Context) {
	commentID := c.Param("commentId")
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Không tìm thấy người dùng"})
		return
	}
	user, ok := userInterface.(*users.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Định dạng người dùng không hợp lệ"})
		return
	}

	if err := LikeComment(commentID, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Đã thích bình luận"})
}

// UnlikeCommentHandler handles unliking a comment.
func UnlikeCommentHandler(c *gin.Context) {
	commentID := c.Param("commentId")
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Không tìm thấy người dùng"})
		return
	}
	user, ok := userInterface.(*users.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Định dạng người dùng không hợp lệ"})
		return
	}

	if err := UnlikeComment(commentID, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Đã bỏ thích bình luận"})
}

// ============ ADMIN HANDLERS ============

// AdminUpdateCommentIsDeletedHandler updates the is_deleted status of a comment (admin only)
func AdminUpdateCommentIsDeletedHandler(c *gin.Context) {
	commentID := c.Param("commentId")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cần có ID bình luận"})
		return
	}

	var req struct {
		IsDeleted bool `json:"is_deleted"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Định dạng yêu cầu không hợp lệ"})
		return
	}

	// Check if comment exists
	comment, err := GetCommentByIDAdmin(commentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy bình luận"})
		return
	}

	if err := UpdateCommentIsDeleted(commentID, req.IsDeleted); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cập nhật bình luận thất bại"})
		return
	}

	action := "restored"
	if req.IsDeleted {
		action = "deleted"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment " + action + " successfully",
		"data": gin.H{
			"id":         comment.ID,
			"post_id":    comment.PostID,
			"is_deleted": req.IsDeleted,
		},
	})
}

// AdminGetCommentByIDHandler gets a comment by ID without IsDeleted filter (admin only)
func AdminGetCommentByIDHandler(c *gin.Context) {
	commentID := c.Param("commentId")
	if commentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cần có ID bình luận"})
		return
	}

	comment, err := GetCommentByIDAdmin(commentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy bình luận"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":         comment.ID,
			"post_id":    comment.PostID,
			"user_id":    comment.UserID,
			"parent_id":  comment.ParentID,
			"content":    comment.Content,
			"like_num":   comment.LikeNum,
			"status":     comment.Status,
			"is_deleted": comment.IsDeleted,
			"created_at": comment.CreatedAt,
			"updated_at": comment.UpdatedAt,
		},
	})
}
