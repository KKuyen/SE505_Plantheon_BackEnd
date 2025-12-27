package users

import (
	"net/http"
	"strconv"
	"strings"

	"plantheon-backend/common"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Register handles user registration
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Validate request
	if err := ValidateRegisterRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if email already exists
	if _, err := GetUserByEmail(req.Email); err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Email already exists",
		})
		return
	}

	// Check if username already exists
	if _, err := GetUserByUsername(req.Username); err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Username already exists",
		})
		return
	}

	// Create user
	user := &User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
		FullName: req.FullName,
		Role:     UserRole(req.Role), // Set role from request
	}

	if err := CreateUser(user); err != nil {
		// Check for specific database errors
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			if strings.Contains(err.Error(), "users_email_key") {
				c.JSON(http.StatusConflict, gin.H{
					"error": "Email already exists",
				})
			} else if strings.Contains(err.Error(), "users_username_key") {
				c.JSON(http.StatusConflict, gin.H{
					"error": "Username already exists",
				})
			} else {
				c.JSON(http.StatusConflict, gin.H{
					"error": "User already exists",
				})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create user",
			})
		}
		return
	}

	// Generate JWT token
	token, err := common.GenerateJWT(user.ID, user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	// Return response
	response := LoginResponse{
		User:  user.ToUserResponse(),
		Token: token,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data":    response,
	})
}

// Login handles user login
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Find user by email
	user, err := GetUserByEmail(req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid email or password",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to find user",
		})
		return
	}

	// Check password
	if !common.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Check if account is active
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Tài khoản đã bị vô hiệu hóa",
		})
		return
	}

	// Generate JWT token
	token, err := common.GenerateJWT(user.ID, user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	// Return response
	response := LoginResponse{
		User:  user.ToUserResponse(),
		Token: token,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    response,
	})
}

// GetProfile gets current user profile
func GetProfile(c *gin.Context) {
	user, exists := GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found in context",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user.ToUserResponse(),
	})
}

// UpdateProfile updates current user profile
func UpdateProfile(c *gin.Context) {
	user, exists := GetCurrentUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found in context",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Update user fields if provided
	if req.Username != "" {
		if err := ValidateUsername(req.Username); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		user.Username = req.Username
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}

	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	// Save updated user
	if err := UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"data":    user.ToUserResponse(),
	})
}

// GetAllUsers handles admin retrieval of all users with pagination
func GetAllUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid page parameter",
		})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid limit parameter",
		})
		return
	}

	users, total, err := ListUsers(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch users",
		})
		return
	}

	responses := make([]UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, user.ToUserResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responses,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetUserByIDHandler handles admin retrieval of a single user by ID
func GetUserByIDHandler(c *gin.Context) {
	userID := c.Param("id")
	user, err := GetUserByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user.ToUserResponse(),
	})
}

// DisableUserHandler disables a user by ID (admin only)
func DisableUserHandler(c *gin.Context) {
	userID := c.Param("id")

	if _, err := GetUserByID(userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user",
		})
		return
	}

	if err := DisableUser(userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to disable user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User has been disabled successfully",
	})
}

// EnableUserHandler enables a user by ID (admin only)
func EnableUserHandler(c *gin.Context) {
	userID := c.Param("id")

	if _, err := GetUserByID(userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user",
		})
		return
	}

	if err := EnableUser(userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to enable user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User has been enabled successfully",
	})
}

// UpdateUserByAdminHandler updates user info (admin only)
func UpdateUserByAdminHandler(c *gin.Context) {
	userID := c.Param("id")
	var req UpdateUserAdminRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	user, err := GetUserByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user",
		})
		return
	}

	// Email update
	if strings.TrimSpace(req.Email) != "" {
		if err := ValidateEmail(req.Email); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Check duplicate email (excluding current user)
		if existing, err := GetUserByEmail(req.Email); err == nil && existing.ID != userID {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Email already exists",
			})
			return
		}
		user.Email = strings.TrimSpace(req.Email)
	}

	// Username update
	if strings.TrimSpace(req.Username) != "" {
		if err := ValidateUsername(req.Username); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if existing, err := GetUserByUsername(req.Username); err == nil && existing.ID != userID {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Username already exists",
			})
			return
		}
		user.Username = strings.TrimSpace(req.Username)
	}

	// Full name update
	if strings.TrimSpace(req.FullName) != "" {
		fullName := strings.TrimSpace(req.FullName)
		if len(fullName) > 100 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Full name must be less than 100 characters",
			})
			return
		}
		user.FullName = fullName
	}

	// Avatar update
	if strings.TrimSpace(req.Avatar) != "" {
		user.Avatar = strings.TrimSpace(req.Avatar)
	}

	if err := UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"data":    user.ToUserResponse(),
	})
}

// ForgotPassword handles forgot password request - generates and sends OTP
func ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Generate and send OTP
	err := GenerateOTP(req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Email không tồn tại trong hệ thống",
			})
			return
		}
		if customErr, ok := err.(*common.CustomError); ok {
			c.JSON(customErr.Code, gin.H{
				"error": customErr.Message,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to send OTP",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "OTP đã được gửi đến email của bạn",
		"expires_in": "5m",
	})
}

// VerifyOTPHandler handles OTP verification
func VerifyOTPHandler(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Verify OTP
	valid, attemptsRemaining, err := VerifyOTP(req.Email, req.OTP)
	if err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			c.JSON(customErr.Code, gin.H{
				"error": customErr.Message,
			})
			return
		}
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Email không tồn tại trong hệ thống",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to verify OTP",
		})
		return
	}

	if !valid {
		response := VerifyOTPResponse{
			Valid:             false,
			Message:           "OTP không đúng",
			AttemptsRemaining: &attemptsRemaining,
		}
		c.JSON(http.StatusOK, response)
		return
	}

	response := VerifyOTPResponse{
		Valid:   true,
		Message: "OTP hợp lệ",
	}
	c.JSON(http.StatusOK, response)
}

// ResetPasswordWithOTPHandler handles password reset with OTP
func ResetPasswordWithOTPHandler(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Reset password
	err := ResetPasswordWithOTP(req.Email, req.OTP, req.NewPassword)
	if err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			c.JSON(customErr.Code, gin.H{
				"error": customErr.Message,
			})
			return
		}
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Email không tồn tại trong hệ thống",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to reset password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Đặt lại mật khẩu thành công",
	})
}
