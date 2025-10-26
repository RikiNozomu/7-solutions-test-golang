package service

import (
	domain "7-solutions-test-golang/core/domains"
	mock "7-solutions-test-golang/mocks"
	util "7-solutions-test-golang/utils"
	"testing"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// TestCreate validates user creation logic, including duplicate email handling.
func TestCreate(t *testing.T) {
	userRepo := mock.NewMockUserRepository()
	userService := CreateUserService(userRepo)
	password := "test1234"

	// Ensure initial user count is 0
	_, err := userService.GetCount()
	if err != nil {
		t.Errorf("GetCount() = %s; want 0", err.Error())
	}

	// Create a new user
	_, err = userService.Create(domain.UserCreate{
		Name:     "test",
		Email:    "email@com.com",
		Password: password,
	})
	if err != nil {
		t.Errorf("Create() = %s", err.Error())
	}

	// Verify user count is now 1
	count, err := userService.GetCount()
	if err != nil {
		t.Errorf("GetCount() = %s; want 1", err.Error())
	}

	// Attempt to create a user with duplicate email
	_, err = userService.Create(domain.UserCreate{
		Name:     "test",
		Email:    "email@com.com",
		Password: password,
	})
	if err != util.ErrorDuplicateKey {
		t.Errorf("Create() = %s; want = %s", err.Error(), util.ErrorDuplicateKey.Error())
	}

	// Final count should still be 1
	if count != 1 {
		t.Errorf("GetCount() = %d; want 1", count)
	}
}

// TestUpdate validates user update logic and duplicate email conflict.
func TestUpdate(t *testing.T) {
	userRepo := mock.NewMockUserRepository()
	userService := CreateUserService(userRepo)
	password := "test1234"

	// Ensure initial user count is 0
	_, err := userService.GetCount()
	if err != nil {
		t.Errorf("GetCount() = %s; want 0", err.Error())
	}

	// Create a user
	data, err := userService.Create(domain.UserCreate{
		Name:     "test",
		Email:    "email@com.com",
		Password: password,
	})
	if err != nil {
		t.Errorf("Create() = %s", err.Error())
	}

	// Update the user's name
	data, err = userService.Update(data.ID.Hex(), domain.UserUpdate{
		Name: "12345",
	})
	if err != nil {
		t.Errorf("Update() = %s;", err.Error())
	}
	if data.Name != "12345" {
		t.Errorf("Update() = %s; want 12345", data.Name)
	}
	if data.Email != "email@com.com" {
		t.Errorf("Update() = %s; want email@com.com", data.Email)
	}

	// Create another user with a different email
	data, err = userService.Create(domain.UserCreate{
		Name:     "test2",
		Email:    "email2@com.com",
		Password: password,
	})
	if err != nil {
		t.Errorf("Create() = %s", err.Error())
	}

	// Attempt to update second user to use an existing email
	_, err = userService.Update(data.ID.Hex(), domain.UserUpdate{
		Email: "email@com.com",
	})
	if err != util.ErrorDuplicateKey {
		t.Errorf("Update() = %s; want = %s", err.Error(), util.ErrorDuplicateKey.Error())
	}
}

// TestDelete validates user deletion and handling of non-existent user deletion.
func TestDelete(t *testing.T) {
	userRepo := mock.NewMockUserRepository()
	userService := CreateUserService(userRepo)
	password := "test1234"

	// Ensure initial user count is 0
	_, err := userService.GetCount()
	if err != nil {
		t.Errorf("GetCount() = %s; want 0", err.Error())
	}

	// Create a user
	data, err := userService.Create(domain.UserCreate{
		Name:     "test",
		Email:    "email@com.com",
		Password: password,
	})
	if err != nil {
		t.Errorf("Create() = %s", err.Error())
	}

	id := data.ID.Hex()

	// Delete the user
	err = userService.Delete(id)
	if err != nil {
		t.Errorf("Delete(%s) = %s;", id, err.Error())
	}

	// Attempt to delete the same user again
	err = userService.Delete(id)
	if err != mongo.ErrNoDocuments {
		t.Errorf("Delete(%s) = %s; want = %s", id, err.Error(), mongo.ErrNoDocuments.Error())
	}
}
