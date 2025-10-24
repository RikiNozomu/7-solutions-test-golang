package service

import (
	domain "7-solutions-test-golang/core/domains"
	util "7-solutions-test-golang/utils"
)

type AuthService struct {
	userService *UserService
}

func CreateAuthService(userService *UserService) *AuthService {
	return &AuthService{
		userService: userService,
	}
}

func (s *AuthService) Login(login domain.LoginPayload) (*domain.User, error) {
	user, err := s.userService.GetByEmail(login.Email)
	if err != nil {
		return nil, err
	}
	if !util.CheckPasswordHash(login.Password, user.Password) {
		return nil, util.ErrorAuthenticated
	}
	return user, nil
}
