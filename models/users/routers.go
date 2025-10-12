package users

import (
	"net/http"
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
