package repository

import (
	"context"
	"github.com/ashkan-maleki/kiwi-eats/internal/user/repository/entity"
)

type User struct{}

func NewUser() *User {
	return &User{}
}

func (u *User) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	// TODO: Implement login logic
	return nil, nil
}

func (u *User) CreateUser(ctx context.Context, user *entity.User) (string, error) {
	// TODO: Implement login logic
	return "", nil
}
