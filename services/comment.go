package services

import (
	"context"

	"github.com/ortymid/projector/models"
)

type CommentRepo interface {
	AllByTask(context.Context, models.Task) ([]models.Comment, error)
	Create(context.Context, models.Task, models.Comment) (models.Comment, error)
}

type CommentService struct {
	cr CommentRepo
}

func NewCommentService(cr CommentRepo) CommentService {
	return CommentService{cr: cr}
}

func (srv CommentService) CreateComment(t models.Task, text string) (c models.Comment, err error) {
	c.Text = text
	if err = c.Validate(); err != nil {
		return c, err
	}
	return srv.cr.Create(context.TODO(), t, c)
}
