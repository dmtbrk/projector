package services

import (
	"context"

	"github.com/ortymid/projector/models"
	"github.com/ortymid/projector/persistence/repos"
)

type CommentService struct {
	cr repos.CommentRepo
}

func NewCommentService(cr repos.CommentRepo) CommentService {
	return CommentService{cr: cr}
}

func (srv CommentService) CreateComment(t models.Task, text string) (c models.Comment, err error) {
	c.Text = text
	if err = c.Validate(); err != nil {
		return c, err
	}
	return srv.cr.Create(context.TODO(), t, c)
}
