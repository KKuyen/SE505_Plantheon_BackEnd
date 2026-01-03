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

// DeleteUser deletes user by ID and all related data (cascade delete)
// This ensures GDPR compliance and Google Play data deletion requirements
func DeleteUser(id string) error {
	service := NewUserService()
	
	// Start a transaction to ensure all deletes succeed or none do
	return service.db.Transaction(func(tx *gorm.DB) error {
		// 1. Delete all post likes by this user
		tx.Exec("DELETE FROM post_likes WHERE user_id = ?", id)
		
		// 2. Delete all comment likes by this user
		tx.Exec("DELETE FROM comment_likes WHERE user_id = ?", id)
		
		// 3. Delete all comments by this user
		tx.Exec("DELETE FROM comments WHERE user_id = ?", id)
		
		// 4. Delete all posts by this user
		tx.Exec("DELETE FROM posts WHERE user_id = ?", id)
		
		// 5. Delete all scan history by this user
		tx.Exec("DELETE FROM scan_histories WHERE user_id = ?", id)
		
		// 6. Delete all complaints by this user
		tx.Exec("DELETE FROM complaints WHERE user_id = ?", id)
		
		// 7. Delete all notifications for this user
		tx.Exec("DELETE FROM noti WHERE user_id = ?", id)
		
		// 8. Finally, delete the user account
		result := tx.Delete(&User{}, "id = ?", id)
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return result.Error
	})
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

// GenerateOTP creates a random 6-digit OTP and sends it via email
func GenerateOTP(email string) error {
	service := NewUserService()
	
	// Find user by email
	user, err := GetUserByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return gorm.ErrRecordNotFound
		}
		return err
	}
	
	// Check rate limiting (1 minute between requests)
	if user.LastOTPRequest != nil {
		timeSinceLastRequest := common.GetCurrentTime().Sub(*user.LastOTPRequest)
		if timeSinceLastRequest.Minutes() < 1 {
			return common.ErrTooManyRequests
		}
	}
	
	// Generate random 6-digit OTP
	otp := generateRandomOTP()
	
	// Set expiry time (5 minutes from now)
	expiryTime := common.GetCurrentTime().Add(5 * common.Minute)
	
	// Update user with OTP data
	user.OTPCode = &otp
	user.OTPExpiry = &expiryTime
	user.OTPAttempts = 0
	now := common.GetCurrentTime()
	user.LastOTPRequest = &now
	
	if err := service.db.Save(user).Error; err != nil {
		return err
	}
	
	// Send OTP via email
	emailService := common.NewEmailService()
	if err := emailService.SendOTP(email, otp); err != nil {
		return err
	}
	
	return nil
}

// VerifyOTP checks if the provided OTP is valid
func VerifyOTP(email, otp string) (valid bool, attemptsRemaining int, err error) {
	service := NewUserService()
	
	// Find user by email
	user, err := GetUserByEmail(email)
	if err != nil {
		return false, 0, err
	}
	
	// Check if OTP exists
	if user.OTPCode == nil || user.OTPExpiry == nil {
		return false, 0, common.ErrInvalidOTP
	}
	
	// Check if OTP has expired
	if common.GetCurrentTime().After(*user.OTPExpiry) {
		return false, 0, common.ErrOTPExpired
	}
	
	// Check if max attempts exceeded
	if user.OTPAttempts >= 3 {
		return false, 0, common.ErrTooManyAttempts
	}
	
	// Verify OTP
	if *user.OTPCode != otp {
		// Increment attempts
		user.OTPAttempts++
		service.db.Save(user)
		remaining := 3 - user.OTPAttempts
		return false, remaining, nil
	}
	
	// OTP is valid
	return true, 0, nil
}

// ResetPasswordWithOTP resets user password after verifying OTP
func ResetPasswordWithOTP(email, otp, newPassword string) error {
	service := NewUserService()
	
	// Verify OTP one last time
	valid, _, err := VerifyOTP(email, otp)
	if err != nil {
		return err
	}
	if !valid {
		return common.ErrInvalidOTP
	}
	
	// Find user
	user, err := GetUserByEmail(email)
	if err != nil {
		return err
	}
	
	// Hash new password
	hashedPassword, err := common.HashPassword(newPassword)
	if err != nil {
		return err
	}
	
	// Update password and clear OTP data
	user.Password = hashedPassword
	user.OTPCode = nil
	user.OTPExpiry = nil
	user.OTPAttempts = 0
	user.LastOTPRequest = nil
	
	return service.db.Save(user).Error
}

// generateRandomOTP creates a random 6-digit OTP
func generateRandomOTP() string {
	return common.GenerateRandomCode(6)
}
