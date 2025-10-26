package service

import (
	domain "7-solutions-test-golang/core/domains"
	util "7-solutions-test-golang/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	userRepo domain.UserRepository // Interface to interact with user data storage
}

// CreateUserService initializes a new UserService with the given repository.
func CreateUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// Create hashes the password and creates a new user record.
func (s *UserService) Create(user domain.UserCreate) (*domain.User, error) {
	// Securely hash the user's password
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	// Prepare user data for insertion
	createData := domain.User{
		ID:        primitive.NewObjectID(),
		Name:      user.Name,
		Email:     user.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}

	// Save to repository
	err = s.userRepo.Create(createData)
	if err != nil {
		return nil, err
	}

	return &createData, nil
}

// Get retrieves a user by their ID.
func (s *UserService) Get(id string) (*domain.User, error) {
	return s.userRepo.GetOne("_id", id)
}

// GetByEmail retrieves a user by their email address.
func (s *UserService) GetByEmail(email string) (*domain.User, error) {
	return s.userRepo.GetOne("email", email)
}

// GetAll returns all users in the system.
func (s *UserService) GetAll() ([]domain.User, error) {
	return s.userRepo.GetAll()
}

// GetCount returns the total number of users.
func (s *UserService) GetCount() (int64, error) {
	return s.userRepo.GetCount()
}

// Update modifies an existing user's name and/or email.
func (s *UserService) Update(id string, user domain.UserUpdate) (*domain.User, error) {
	// Fetch existing user data
	oldData, err := s.userRepo.GetOne("_id", id)
	if err != nil {
		return nil, err
	}

	// Apply updates if provided
	if len(user.Email) > 0 {
		oldData.Email = user.Email
	}
	if len(user.Name) > 0 {
		oldData.Name = user.Name
	}

	// Save updated user
	return oldData, s.userRepo.Update(*oldData)
}

// Delete removes a user by ID.
func (s *UserService) Delete(id string) error {
	return s.userRepo.Delete(id)
}
