package services

import (
	"context"

	"github.com/ortymid/projector/models"
)

// const defaultColumnName = "New Column"

type ColumnRepo interface {
	AllByBoard(context.Context, models.Board) ([]models.Column, error)
	Create(context.Context, models.Board, models.Column) (models.Column, error)
}

type ColumnService struct {
	cr ColumnRepo
}

func NewColumnService(cr ColumnRepo) ColumnService {
	return ColumnService{cr: cr}
}

func (srv ColumnService) CreateColumn(b models.Board, name string) (c models.Column, err error) {
	c.Name = name
	if err = c.Validate(); err != nil {
		return c, err
	}
	return srv.cr.Create(context.TODO(), b, c)
}
