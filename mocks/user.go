package mock

import (
	domain "7-solutions-test-golang/core/domains"
	util "7-solutions-test-golang/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MockUserRepository struct {
	users []domain.User
}

func NewMockUserRepository() domain.UserRepository {
	return &MockUserRepository{users: []domain.User{}}
}

func findUserByEmail(email string, users []domain.User) (*domain.User, int) {
	for index, user := range users {
		if user.Email == email {
			return &user, index
		}
	}
	return nil, -1
}

func findUserByID(id string, users []domain.User) (*domain.User, int) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, -1
	}

	for index, user := range users {
		if user.ID == objectID {
			return &user, index
		}
	}
	return nil, -1
}

// Create implements domain.UserRepository.
func (m *MockUserRepository) Create(user domain.User) error {
	if _, index := findUserByEmail(user.Email, m.users); index > -1 {
		return util.ErrorDuplicateKey
	}
	m.users = append(m.users, user)
	return nil
}

// GetByEmail implements domain.UserRepository.
func (m *MockUserRepository) GetOne(key string, value any) (*domain.User, error) {
	str, ok := value.(string)
	if !ok {
		return nil, mongo.ErrNoDocuments
	}
	switch key {
	case "_id":
		user, index := findUserByID(str, m.users)
		if index == -1 {
			return nil, mongo.ErrNoDocuments
		}
		return user, nil
	case "email":
		user, index := findUserByEmail(str, m.users)
		if index == -1 {
			return nil, mongo.ErrNoDocuments
		}
		return user, nil
	default:
		return nil, mongo.ErrNoDocuments
	}
}

// GetAll implements domain.UserRepository.
func (m *MockUserRepository) GetAll() ([]domain.User, error) {
	if len(m.users) <= 0 {
		return []domain.User{}, nil
	}
	return m.users, nil
}

// GetCount implements domain.UserRepository.
func (m *MockUserRepository) GetCount() (int64, error) {
	return int64(len(m.users)), nil
}

// Update implements domain.UserRepository.
func (m *MockUserRepository) Update(user domain.User) error {
	data, index := findUserByID(user.ID.Hex(), m.users)
	if index == -1 {
		return mongo.ErrNoDocuments
	}

	if _, index := findUserByEmail(user.Email, append(m.users[:index], m.users[index+1:]...)); index > -1 {
		return util.ErrorDuplicateKey
	}

	data.Email = user.Email
	data.Name = user.Name
	m.users[index] = *data

	return nil
}

// Delete implements domain.UserRepository.
func (m *MockUserRepository) Delete(id string) error {
	_, index := findUserByID(id, m.users)
	if index == -1 {
		return mongo.ErrNoDocuments
	}
	m.users = append(m.users[:index], m.users[index+1:]...)
	return nil
}
