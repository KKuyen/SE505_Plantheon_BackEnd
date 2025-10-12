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
	
	return service.db.Create(user).Error
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
	return service.db.Delete(&User{}, "id = ?", id).Error
}