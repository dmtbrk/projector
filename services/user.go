package services

import (
	"context"
	"errors"

	"github.com/ortymid/projector/models"
)

var ErrUserExist = errors.New("projector: user already exists")

type UserRepo interface {
	All(context.Context) ([]models.User, error)
	Create(context.Context, models.User) (models.User, error)
}

type UserService struct {
	ur UserRepo
}

func NewUserService(repo UserRepo) UserService {
	return UserService{ur: repo}
}

func (srv UserService) CreateUser(name string) (user models.User, err error) {
	user.Name = name
	if err = user.Validate(); err != nil {
		return user, err
	}
	user, err = srv.ur.Create(context.TODO(), user)
	return user, err
}
