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
	userService *UserService
}

func CreateAuthService(userService *UserService) *AuthService {
	return &AuthService{
		userService: userService,
	}
}

func (s *AuthService) Login(login domain.LoginPayload) (*string, int64, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	expSecond, err := strconv.Atoi(os.Getenv("JWT_TIME_EXPIRED_SECOND"))
	if err != nil {
		return nil, 0, err
	}

	user, err := s.userService.GetByEmail(login.Email)
	if err != nil {
		return nil, 0, err
	}
	if !util.CheckPasswordHash(login.Password, user.Password) {
		return nil, 0, util.ErrorAuthenticated
	}

	expire_in_unix := time.Now().Add(time.Duration(expSecond) * time.Second).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  user.ID.Hex(),
			"exp": expire_in_unix,
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, 0, util.ErrorAuthenticated
	}
	return &tokenString, expire_in_unix, nil
}
