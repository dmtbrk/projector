package repos

import (
	"context"

	"github.com/ortymid/projector/models"
)

type UserRepo interface {
	WithTx(context.Context, Tx, func(UserRepo) error) (Tx, error)
	All(context.Context) ([]models.User, error)
	Create(context.Context, models.User) (models.User, error)
}
