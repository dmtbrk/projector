package repos

import (
	"context"

	"github.com/ortymid/projector/models"
)

type TaskRepo interface {
	WithTx(context.Context, Tx, func(TaskRepo) error) (Tx, error)
	AllByColumn(context.Context, models.Column) ([]models.Task, error)
	AllByBoard(context.Context, models.Board) ([]models.Task, error)
	Create(context.Context, models.Column, models.Task) (models.Task, error)
}
