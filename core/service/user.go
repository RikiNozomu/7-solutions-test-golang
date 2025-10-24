package service

import (
	domain "7-solutions-test-golang/core/domains"
	util "7-solutions-test-golang/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	userRepo domain.UserRepository
}

func CreateUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Create(user domain.UserCreate) (*domain.User, error) {
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	createData := domain.User{
		ID:        primitive.NewObjectID(),
		Name:      user.Name,
		Email:     user.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}
	err = s.userRepo.Create(createData)
	if err != nil {
		return nil, err
	}
	return &createData, nil
}

func (s *UserService) Get(id string) (*domain.User, error) {
	return s.userRepo.GetOne("_id", id)
}

func (s *UserService) GetByEmail(email string) (*domain.User, error) {
	return s.userRepo.GetOne("email", email)
}

func (s *UserService) GetAll() ([]domain.User, error) {
	return s.userRepo.GetAll()
}

func (s *UserService) GetCount() (int64, error) {
	return s.userRepo.GetCount()
}

func (s *UserService) Update(id string, user domain.UserUpdate) (*domain.User, error) {
	oldData, err := s.userRepo.GetOne("_id", id)
	if err != nil {
		return nil, err
	}

	if len(user.Email) > 0 {
		oldData.Email = user.Email
	}

	if len(user.Name) > 0 {
		oldData.Name = user.Name
	}
	return oldData, s.userRepo.Update(*oldData)
}

func (s *UserService) Delete(id string) error {
	return s.userRepo.Delete(id)
}
