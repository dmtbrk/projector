package services

import (
	"context"

	"github.com/ortymid/projector/models"
)

type BoardRepo interface {
	AllByUser(context.Context, models.User) ([]models.Board, error)
	Create(context.Context, models.User, models.Board) (models.Board, error)
}

type BoardService struct {
	br BoardRepo
	cr ColumnRepo
}

func NewBoardService(br BoardRepo, cr ColumnRepo) BoardService {
	return BoardService{br: br, cr: cr}
}

func (srv BoardService) CreateBoard(u models.User, name, desc string) (b models.Board, err error) {
	b.Name = name
	b.Description = desc
	if err = b.Validate(); err != nil {
		return b, err
	}
	b, err = srv.br.Create(context.TODO(), u, b)
	return b, err
}

func (srv BoardService) CreateBoardDefault(u models.User, name, desc string) (b models.Board, err error) {
	if b, err = srv.CreateBoard(u, name, desc); err != nil {
		return b, err
	}
	colSrv := NewColumnService(srv.cr)
	_, err = colSrv.CreateColumn(b, "New Column")
	return b, err
}
