package service

import (
	domain "7-solutions-test-golang/core/domains"
	mock "7-solutions-test-golang/mocks"
	util "7-solutions-test-golang/utils"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// TestLogin validates the AuthService.Login behavior under different scenarios.
func TestLogin(t *testing.T) {
	// Set environment variables for JWT configuration
	os.Setenv("JWT_SECRET", "jwt-secret")
	os.Setenv("JWT_TIME_EXPIRED_SECOND", "3600")

	// Initialize mock repository and services
	userRepo := mock.NewMockUserRepository()
	userService := CreateUserService(userRepo)
	AuthService := CreateAuthService(userService)

	email := "email@com.com"
	password := "test1234"

	// Create a user for login testing
	_, err := userService.Create(domain.UserCreate{
		Name:     "test",
		Email:    email,
		Password: password,
	})
	if err != nil {
		t.Errorf("Create() = %s", err.Error())
	}

	// Verify user count is 1
	count, err := userService.GetCount()
	if err != nil {
		t.Errorf("GetCount() = %s; want 1", err.Error())
	}
	if count != 1 {
		t.Errorf("GetCount() = %d; want 1", count)
	}

	// Test login with incorrect password
	_, _, err = AuthService.Login(domain.LoginPayload{Email: email, Password: "55555555"})
	if err != util.ErrorAuthenticated {
		t.Errorf("Login() = %s; want = %s", err.Error(), util.ErrorAuthenticated.Error())
	}

	// Test login with non-existent email
	_, _, err = AuthService.Login(domain.LoginPayload{Email: "ppp@po.co", Password: password})
	if err != mongo.ErrNoDocuments {
		t.Errorf("Login() = %s; want = %s", err.Error(), mongo.ErrNoDocuments.Error())
	}

	// Test login with correct credentials
	_, _, err = AuthService.Login(domain.LoginPayload{Email: email, Password: password})
	if err != nil {
		t.Errorf("Login() = %s;", err.Error())
	}
}
