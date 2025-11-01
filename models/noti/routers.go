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
			"error": "User not found in context",
		})
		return
	}

	user, ok := userInterface.(*users.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user format",
		})
		return
	}

	notifications, err := GetNotificationsByUserID(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get notifications",
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
			"error": "User not found in context",
		})
		return
	}

	user, ok := userInterface.(*users.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user format",
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
			"error": "Failed to mark notification as seen",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification marked as seen successfully",
	})
}

func DeleteNotificationHandler(c *gin.Context) {
	// Lấy thông tin user từ JWT token context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found in context",
		})
		return
	}

	user, ok := userInterface.(*users.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user format",
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
			"error": "Failed to delete notification",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification deleted successfully",
	})
}

func DeleteAllNotificationsHandler(c *gin.Context) {
	// Lấy thông tin user từ JWT token context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found in context",
		})
		return
	}

	user, ok := userInterface.(*users.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user format",
		})
		return
	}

	if err := DeleteAllNotificationsByUserID(user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete all notifications",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All notifications deleted successfully",
	})
}

func CreateNotificationHandler(c *gin.Context) {
	// Lấy thông tin user từ JWT token context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found in context",
		})
		return
	}

	user, ok := userInterface.(*users.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user format",
		})
		return
	}

	var req CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
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
			"error": "Failed to create notification",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Notification created successfully",
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

