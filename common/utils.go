package common

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var JWTSecret = []byte("your-secret-key") // Thay đổi secret key này trong production

// Claims struct for JWT
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// HashPassword hashes password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT generates JWT token
func GenerateJWT(userID, email, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// ValidateJWT validates JWT token
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

// Time constants for OTP
const (
	Minute = time.Minute
)

// Custom errors for OTP operations
var (
	ErrTooManyRequests  = &CustomError{Code: 429, Message: "Vui lòng đợi 1 phút trước khi request OTP mới"}
	ErrInvalidOTP       = &CustomError{Code: 400, Message: "OTP không hợp lệ"}
	ErrOTPExpired       = &CustomError{Code: 400, Message: "OTP đã hết hạn"}
	ErrTooManyAttempts  = &CustomError{Code: 400, Message: "Đã nhập sai quá 3 lần, vui lòng request OTP mới"}
)

// CustomError represents a custom error with code
type CustomError struct {
	Code    int
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

// GetCurrentTime returns current time (useful for testing)
func GetCurrentTime() time.Time {
	return time.Now()
}

// GenerateRandomCode generates a cryptographically secure random numeric code
func GenerateRandomCode(length int) string {
	const digits = "0123456789"
	code := make([]byte, length)
	
	for i := range code {
		// Generate random number from 0-9
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			// Fallback to time-based if crypto/rand fails (very unlikely)
			code[i] = digits[time.Now().UnixNano()%10]
		} else {
			code[i] = digits[num.Int64()]
		}
	}
	
	return string(code)
}
