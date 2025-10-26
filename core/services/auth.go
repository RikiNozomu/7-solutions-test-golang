package service

import (
	domain "7-solutions-test-golang/core/domains"
	util "7-solutions-test-golang/utils"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userService *UserService // Dependency on UserService for user lookup
}

// CreateAuthService initializes a new AuthService instance.
func CreateAuthService(userService *UserService) *AuthService {
	return &AuthService{
		userService: userService,
	}
}

// Login authenticates a user and returns a signed JWT token with expiration.
func (s *AuthService) Login(login domain.LoginPayload) (*string, int64, error) {
	// Load secret key and expiration time from environment
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	expSecond, err := strconv.Atoi(os.Getenv("JWT_TIME_EXPIRED_SECOND"))
	if err != nil {
		return nil, 0, err
	}

	// Find user by email
	user, err := s.userService.GetByEmail(login.Email)
	if err != nil {
		return nil, 0, err
	}

	// Verify password using bcrypt
	if !util.CheckPasswordHash(login.Password, user.Password) {
		return nil, 0, util.ErrorAuthenticated
	}

	// Calculate token expiration time (Unix timestamp)
	expire_in_unix := time.Now().Add(time.Duration(expSecond) * time.Second).Unix()

	// Create JWT token with user ID and expiration
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID.Hex(),
		"exp": expire_in_unix,
	})

	// Sign token with secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, 0, util.ErrorAuthenticated
	}

	// Return token and expiration
	return &tokenString, expire_in_unix, nil
}
