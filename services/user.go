package services

import (
	"context"
	"errors"

	"github.com/ortymid/projector/models"
	"github.com/ortymid/projector/persistence/repos"
)

var ErrUserExist = errors.New("projector: user already exists")

type UserService struct {
	ur repos.UserRepo
}

func NewUserService(ur repos.UserRepo) UserService {
	return UserService{ur: ur}
}

func (srv UserService) CreateUser(name string) (user models.User, err error) {
	user.Name = name
	if err = user.Validate(); err != nil {
		return user, err
	}
	user, err = srv.ur.Create(context.TODO(), user)
	return user, err
}
