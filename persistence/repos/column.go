package repos

import (
	"context"

	"github.com/ortymid/projector/models"
)

type ColumnRepo interface {
	WithTx(context.Context, Tx, func(ColumnRepo) error) (Tx, error)
	AllByBoard(context.Context, models.Board) ([]models.Column, error)
	Create(context.Context, models.Board, models.Column) (models.Column, error)
}
