package users

import (
	"plantheon-backend/common"

	"gorm.io/gorm"
)

// UserService handles all database operations for users
type UserService struct {
	db *gorm.DB
}

// NewUserService creates a new user service instance
func NewUserService() *UserService {
	return &UserService{
		db: common.GetDB(),
	}
}

// CreateUser creates a new user
func CreateUser(user *User) error {
	service := NewUserService()
	hashedPassword, err := common.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// Set default role if not specified
	if user.Role == "" {
		user.Role = RoleUser
	}
	// Activate by default
	if !user.IsActive {
		user.IsActive = true
	}

	return service.db.Create(user).Error
}

// ListUsers returns paginated users and total count
func ListUsers(page, limit int) ([]User, int64, error) {
	service := NewUserService()
	var users []User
	var total int64

	query := service.db.Model(&User{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&users).Error
	return users, total, err
}

// GetUserByEmail finds user by email
func GetUserByEmail(email string) (*User, error) {
	service := NewUserService()
	var user User
	err := service.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

// GetUserByUsername finds user by username
func GetUserByUsername(username string) (*User, error) {
	service := NewUserService()
	var user User
	err := service.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

// GetUserByID finds user by ID
func GetUserByID(id string) (*User, error) {
	service := NewUserService()
	var user User
	err := service.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

// UpdateUser updates user information
func UpdateUser(user *User) error {
	service := NewUserService()
	return service.db.Save(user).Error
}

// DeleteUser deletes user by ID
func DeleteUser(id string) error {
	service := NewUserService()
	result := service.db.Delete(&User{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// DisableUser sets is_active to false for the given user ID
func DisableUser(id string) error {
	service := NewUserService()
	result := service.db.Model(&User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_active": false,
	})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// EnableUser sets is_active to true for the given user ID
func EnableUser(id string) error {
	service := NewUserService()
	result := service.db.Model(&User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_active": true,
	})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
