package users

import (
	"time"

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

// IsAdmin checks if user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsUser checks if user has user role
func (u *User) IsUser() bool {
	return u.Role == RoleUser
}
