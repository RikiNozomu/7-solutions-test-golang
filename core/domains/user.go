package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `json:"-"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type UserCreate struct {
	Name     string `json:"name" binding:"required,min=8,max=24"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=24,nospace,nospecial,havelower,haveupper,havenumber"`
}

type UserUpdate struct {
	Name  string `json:"name" binding:"omitempty,min=8,max=24"`
	Email string `json:"email" binding:"omitempty,email"`
}

type UserRepository interface {
	Create(user User) error
	GetOne(key string, value any) (*User, error)
	GetAll() ([]User, error)
	GetCount() (int64, error)
	Update(user User) error
	Delete(id string) error
}
