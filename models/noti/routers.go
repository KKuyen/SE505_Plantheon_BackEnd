package noti

import (
	"net/http"

	"plantheon-backend/models/users"

	"github.com/gin-gonic/gin"
)

func GetNotificationsHandler(c *gin.Context) {
	// Lấy thông tin user từ JWT token context
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

	notifications, err := GetNotificationsByUserID(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Lấy danh sách thông báo thất bại",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": notifications,
	})
}

func MarkNotificationAsSeenHandler(c *gin.Context) {
	// Lấy thông tin user từ JWT token context
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

	id := c.Param("id")
	if err := ValidateIdParam(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := MarkNotificationAsSeen(id, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Đánh dấu thông báo đã xem thất bại",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Đánh dấu thông báo đã xem thành công",
	})
}

func DeleteNotificationHandler(c *gin.Context) {
	// Lấy thông tin user từ JWT token context
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

	id := c.Param("id")
	if err := ValidateIdParam(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := DeleteNotificationByID(id, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Xóa thông báo thất bại",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Xóa thông báo thành công",
	})
}

func DeleteAllNotificationsHandler(c *gin.Context) {
	// Lấy thông tin user từ JWT token context
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

	if err := DeleteAllNotificationsByUserID(user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Xóa tất cả thông báo thất bại",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Xóa tất cả thông báo thành công",
	})
}

func CreateNotificationHandler(c *gin.Context) {
	// Lấy thông tin user từ JWT token context
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

	var req CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Yêu cầu không hợp lệ",
		})
		return
	}

	if err := ValidateCreateNotificationRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	notification := &Notification{
		UserID:  user.ID,
		Title:   req.Title,
		Content: req.Content,
		PostID:  req.PostID,
		IsRead:  false,
	}

	if err := CreateNotification(notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Tạo thông báo thất bại",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tạo thông báo thành công",
		"data": NotificationResponse{
			ID:        notification.ID,
			UserID:    notification.UserID,
			Title:     notification.Title,
			Content:   notification.Content,
			IsRead:    notification.IsRead,
			PostID:    notification.PostID,
			CreatedAt: notification.CreatedAt,
		},
	})
}

