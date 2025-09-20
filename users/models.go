package users

import (
	"time"

	"plantheon-backend/common"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRole enum for user roles
type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type User struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"` // "-" để không serialize password
	FullName  string    `json:"full_name"`
	Avatar    string    `json:"avatar"`
	Role      UserRole  `json:"role" gorm:"type:varchar(20);default:'user';not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// CreateUser creates a new user
func CreateUser(user *User) error {
	hashedPassword, err := common.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	
	// Set default role if not specified
	if user.Role == "" {
		user.Role = RoleUser
	}
	
	return common.GetDB().Create(user).Error
}

// GetUserByEmail finds user by email
func GetUserByEmail(email string) (*User, error) {
	var user User
	err := common.GetDB().Where("email = ?", email).First(&user).Error
	return &user, err
}

// GetUserByUsername finds user by username
func GetUserByUsername(username string) (*User, error) {
	var user User
	err := common.GetDB().Where("username = ?", username).First(&user).Error
	return &user, err
}

// GetUserByID finds user by ID
func GetUserByID(id string) (*User, error) {
	var user User
	err := common.GetDB().Where("id = ?", id).First(&user).Error
	return &user, err
}

// UpdateUser updates user information
func UpdateUser(user *User) error {
	return common.GetDB().Save(user).Error
}

// DeleteUser deletes user by ID
func DeleteUser(id string) error {
	return common.GetDB().Delete(&User{}, "id = ?", id).Error
}

// IsAdmin checks if user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsUser checks if user has user role
func (u *User) IsUser() bool {
	return u.Role == RoleUser
}
