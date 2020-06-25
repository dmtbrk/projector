package repos

import (
	"context"

	"github.com/ortymid/projector/models"
)

type CommentRepo interface {
	WithTx(context.Context, Tx, func(CommentRepo) error) (Tx, error)
	AllByTask(context.Context, models.Task) ([]models.Comment, error)
	Create(context.Context, models.Task, models.Comment) (models.Comment, error)
}
