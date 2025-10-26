package mock

import (
	domain "7-solutions-test-golang/core/domains"
	util "7-solutions-test-golang/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// MockUserRepository is an in-memory implementation of domain.UserRepository.
type MockUserRepository struct {
	users []domain.User
}

// NewMockUserRepository initializes a new empty mock repository.
func NewMockUserRepository() domain.UserRepository {
	return &MockUserRepository{users: []domain.User{}}
}

// findUserByEmail searches for a user by email in the mock slice.
func findUserByEmail(email string, users []domain.User) (*domain.User, int) {
	for index, user := range users {
		if user.Email == email {
			return &user, index
		}
	}
	return nil, -1
}

// findUserByID searches for a user by id in the mock slice.
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

// Create adds a new user to the mock repository.
func (m *MockUserRepository) Create(user domain.User) error {
	if _, index := findUserByEmail(user.Email, m.users); index > -1 {
		return util.ErrorDuplicateKey
	}
	m.users = append(m.users, user)
	return nil
}

// GetOne retrieves a user by key ("_id" or "email").
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

// GetAll returns all users in the mock repository.
func (m *MockUserRepository) GetAll() ([]domain.User, error) {
	if len(m.users) <= 0 {
		return []domain.User{}, nil
	}
	return m.users, nil
}

// GetCount returns the number of users in the mock repository.
func (m *MockUserRepository) GetCount() (int64, error) {
	return int64(len(m.users)), nil
}

// Update modifies an existing user by ID.
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

// Delete removes a user by ID from the mock repository.
func (m *MockUserRepository) Delete(id string) error {
	_, index := findUserByID(id, m.users)
	if index == -1 {
		return mongo.ErrNoDocuments
	}
	m.users = append(m.users[:index], m.users[index+1:]...)
	return nil
}
