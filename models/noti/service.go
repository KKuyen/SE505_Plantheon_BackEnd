package noti

import (
	"plantheon-backend/common"

	"gorm.io/gorm"
)

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		db: common.GetDB(),
	}
}

func GetNotificationsByUserID(userID string) (NotificationListResponse, error) {
	service := NewNotificationService()
	var notifications []Notification
	if err := service.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&notifications).Error; err != nil {
		return NotificationListResponse{}, err
	}

	var notificationResponses []NotificationResponse
	for _, notification := range notifications {
		notificationResponses = append(notificationResponses, NotificationResponse{
			ID:        notification.ID,
			UserID:    notification.UserID,
			Title:     notification.Title,
			Content:   notification.Content,
			IsRead:    notification.IsRead,
			PostID:    notification.PostID,
			CreatedAt: notification.CreatedAt,
		})
	}

	return NotificationListResponse{
		Notifications: notificationResponses,
		Total:         len(notificationResponses),
	}, nil
}

func MarkNotificationAsSeen(id string, userID string) error {
	service := NewNotificationService()
	if err := service.db.Model(&Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true).Error; err != nil {
		return err
	}
	return nil
}

func DeleteNotificationByID(id string, userID string) error {
	service := NewNotificationService()
	if err := service.db.Where("id = ? AND user_id = ?", id, userID).Delete(&Notification{}).Error; err != nil {
		return err
	}
	return nil
}

func DeleteAllNotificationsByUserID(userID string) error {
	service := NewNotificationService()
	if err := service.db.Where("user_id = ?", userID).Delete(&Notification{}).Error; err != nil {
		return err
	}
	return nil
}

func CreateNotification(notification *Notification) error {
	service := NewNotificationService()
	if err := service.db.Create(notification).Error; err != nil {
		return err
	}
	return nil
}

// CreatePostNotification is a helper to create notification linked to a post.
func CreatePostNotification(userID string, postID *string, title, content string) error {
	notification := &Notification{
		UserID:  userID,
		Title:   title,
		Content: content,
		PostID:  postID,
		IsRead:  false,
	}
	return CreateNotification(notification)
}
