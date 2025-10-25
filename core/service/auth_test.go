package service

import (
	domain "7-solutions-test-golang/core/domains"
	mock "7-solutions-test-golang/mocks"
	util "7-solutions-test-golang/utils"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestLogin(t *testing.T) {
	os.Setenv("JWT_SECRET", "jwt-secret")
	os.Setenv("JWT_TIME_EXPIRED_SECOND", "3600")

	userRepo := mock.NewMockUserRepository()
	userService := CreateUserService(userRepo)
	AuthService := CreateAuthService(userService)
	email := "email@com.com"
	password := "test1234"

	_, err := userService.Create(domain.UserCreate{
		Name:     "test",
		Email:    email,
		Password: password,
	})
	if err != nil {
		t.Errorf("Create() = %s", err.Error())
	}

	count, err := userService.GetCount()
	if err != nil {
		t.Errorf("GetCount() = %s; want 1", err.Error())
	}
	if count != 1 {
		t.Errorf("GetCount() = %d; want 1", count)
	}

	// wrong password
	_, _, err = AuthService.Login(domain.LoginPayload{Email: email, Password: "55555555"})
	if err != util.ErrorAuthenticated {
		t.Errorf("Login() = %s; want = %s", err.Error(), util.ErrorAuthenticated.Error())
	}

	// wrong email
	_, _, err = AuthService.Login(domain.LoginPayload{Email: "ppp@po.co", Password: password})
	if err != mongo.ErrNoDocuments {
		t.Errorf("Login() = %s; want = %s", err.Error(), mongo.ErrNoDocuments.Error())
	}

	// correct user
	_, _, err = AuthService.Login(domain.LoginPayload{Email: email, Password: password})
	if err != nil {
		t.Errorf("Login() = %s;", err.Error())
	}
}
