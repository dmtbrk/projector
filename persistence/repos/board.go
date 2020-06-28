package repos

import (
	"context"

	"github.com/ortymid/projector/models"
)

type BoardRepo interface {
	WithTx(context.Context, Tx, func(BoardRepo) error) (Tx, error)
	AllByUser(context.Context, models.User) ([]models.Board, error)
	Create(context.Context, models.User, models.Board) (models.Board, error)
}
