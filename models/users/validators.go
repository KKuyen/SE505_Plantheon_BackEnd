package users

import (
	"errors"
	"regexp"
	"strings"
)

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email is required")
	}
	
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(strings.ToLower(email)) {
		return errors.New("invalid email format")
	}
	
	return nil
}

// ValidateUsername validates username
func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return errors.New("username is required")
	}
	
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	
	if len(username) > 50 {
		return errors.New("username must be less than 50 characters")
	}
	
	// Only allow alphanumeric and underscore
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !usernameRegex.MatchString(username) {
		return errors.New("username can only contain letters, numbers, and underscores")
	}
	
	return nil
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password is required")
	}
	
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	
	if len(password) > 100 {
		return errors.New("password must be less than 100 characters")
	}
	
	return nil
}

// ValidateRole validates user role
func ValidateRole(role string) error {
	if role != "user" && role != "admin" {
		return errors.New("role must be either 'user' or 'admin'")
	}
	return nil
}

// ValidateRegisterRequest validates registration request
func ValidateRegisterRequest(req *RegisterRequest) error {
	if err := ValidateEmail(req.Email); err != nil {
		return err
	}
	
	if err := ValidateUsername(req.Username); err != nil {
		return err
	}
	
	if err := ValidatePassword(req.Password); err != nil {
		return err
	}
	
	req.FullName = strings.TrimSpace(req.FullName)
	if req.FullName == "" {
		return errors.New("full name is required")
	}
	
	if len(req.FullName) > 100 {
		return errors.New("full name must be less than 100 characters")
	}
	
	
	return nil
}
