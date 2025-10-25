package service

import (
	domain "7-solutions-test-golang/core/domains"
	mock "7-solutions-test-golang/mocks"
	util "7-solutions-test-golang/utils"
	"testing"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestCreate(t *testing.T) {
	userRepo := mock.NewMockUserRepository()
	userService := CreateUserService(userRepo)
	password := "test1234"

	_, err := userService.GetCount()
	if err != nil {
		t.Errorf("GetAll() = %s; want 0", err.Error())
	}

	_, err = userService.Create(domain.UserCreate{
		Name:     "test",
		Email:    "email@com.com",
		Password: password,
	})
	if err != nil {
		t.Errorf("Create() = %s", err.Error())
	}

	count, err := userService.GetCount()
	if err != nil {
		t.Errorf("GetAll() = %s; want 1", err.Error())
	}

	_, err = userService.Create(domain.UserCreate{
		Name:     "test",
		Email:    "email@com.com",
		Password: password,
	})
	if err != util.ErrorDuplicateKey {
		t.Errorf("Create() = %s; want = %s", err.Error(), util.ErrorDuplicateKey.Error())
	}

	if count != 1 {
		t.Errorf("GetAll() = %d; want 1", count)
	}
}
func TestUpdate(t *testing.T) {
	userRepo := mock.NewMockUserRepository()
	userService := CreateUserService(userRepo)
	password := "test1234"

	_, err := userService.GetCount()
	if err != nil {
		t.Errorf("GetCount() = %s; want 0", err.Error())
	}

	data, err := userService.Create(domain.UserCreate{
		Name:     "test",
		Email:    "email@com.com",
		Password: password,
	})
	if err != nil {
		t.Errorf("Create() = %s", err.Error())
	}

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

	data, err = userService.Create(domain.UserCreate{
		Name:     "test2",
		Email:    "email2@com.com",
		Password: password,
	})
	if err != nil {
		t.Errorf("Create() = %s", err.Error())
	}

	_, err = userService.Update(data.ID.Hex(), domain.UserUpdate{
		Email: "email@com.com",
	})

	if err != util.ErrorDuplicateKey {
		t.Errorf("Update() = %s; want = %s", err.Error(), util.ErrorDuplicateKey.Error())
	}
}
func TestDelete(t *testing.T) {
	userRepo := mock.NewMockUserRepository()
	userService := CreateUserService(userRepo)

	password := "test1234"

	_, err := userService.GetCount()
	if err != nil {
		t.Errorf("GetCount() = %s; want 0", err.Error())
	}

	data, err := userService.Create(domain.UserCreate{
		Name:     "test",
		Email:    "email@com.com",
		Password: password,
	})
	if err != nil {
		t.Errorf("Create() = %s", err.Error())
	}

	id := data.ID.Hex()

	err = userService.Delete(id)

	if err != nil {
		t.Errorf("Delete(%s) = %s;", id, err.Error())
	}

	err = userService.Delete(id)

	if err != mongo.ErrNoDocuments {
		t.Errorf("Delete(%s) = %s; want = %s", id, err.Error(), mongo.ErrNoDocuments.Error())
	}
}
