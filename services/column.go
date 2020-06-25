package services

import (
	"context"

	"github.com/ortymid/projector/models"
	"github.com/ortymid/projector/persistence/repos"
)

// const defaultColumnName = "New Column"

type ColumnService struct {
	cr repos.ColumnRepo
}

func NewColumnService(cr repos.ColumnRepo) ColumnService {
	return ColumnService{cr: cr}
}

func (srv ColumnService) CreateColumn(b models.Board, name string) (c models.Column, err error) {
	c.Name = name
	if err = c.Validate(); err != nil {
		return c, err
	}
	return srv.cr.Create(context.TODO(), b, c)
}
