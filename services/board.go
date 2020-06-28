package services

import (
	"context"
	"fmt"

	"github.com/ortymid/projector/models"
	"github.com/ortymid/projector/persistence/repos"
)

type BoardService struct {
	br repos.BoardRepo
	cr repos.ColumnRepo
}

func NewBoardService(br repos.BoardRepo, cr repos.ColumnRepo) BoardService {
	return BoardService{br: br, cr: cr}
}

// func (srv BoardService) CreateBoard(ctx context.Context, u models.User, name, desc string) (b models.Board, err error) {
// 	tx, err := srv.br.BeginTx()
// 	if err != nil {
// 		err = fmt.Errorf("CreateBoard: %w", err)
// 		return b, err
// 	}

// 	b, err = models.NewBoard(name, desc)
// 	if err != nil {
// 		err = fmt.Errorf("CreateBoard: %w", err)
// 		return b, err
// 	}

// 	b, err = srv.br.Create(ctx, u, b)
// 	if err != nil {
// 		err = fmt.Errorf("CreateBoard: %w", err)
// 		if rErr := tx.Rollback(); rErr != nil {
// 			err = fmt.Errorf("CreateBoard: unable to rollback: %w", rErr)
// 		}
// 		return b, err
// 	}

// 	tx, err = srv.cr.BeginTx(tx)
// 	if err != nil {
// 		err = fmt.Errorf("CreateBoard: unable to begin transaction: %w", err)
// 		return b, err
// 	}

// 	c, err := models.NewColumn("New Column", 0)
// 	_, err = srv.cr.Create(ctx, b, c)
// 	if err != nil {
// 		if rErr := tx.Rollback(); rErr != nil {
// 			err = fmt.Errorf("CreateBoard: unable to rollback: %w", rErr)
// 		} else {
// 			err = fmt.Errorf("CreateBoard: %w", err)
// 		}
// 		return b, err
// 	}

// 	if err = tx.Commit(); err != nil {
// 		err = fmt.Errorf("CreateBoard: unable to commit: %w", err)
// 	}
// 	return b, err
// }

func (srv BoardService) CreateBoard(ctx context.Context, u models.User, name, desc string) (b models.Board, err error) {
	tx, err := srv.br.WithTx(ctx, nil, func(br repos.BoardRepo) error {
		b, err = models.NewBoard(name, desc)
		if err != nil {
			return err
		}
		b, err = br.Create(ctx, u, b)
		return err
	})
	if err != nil {
		err = fmt.Errorf("CreateBoard: %w", err)
		if rErr := tx.Rollback(); rErr != nil {
			err = fmt.Errorf("unable to rollback: %w", rErr)
		}
		return b, err
	}

	tx, err = srv.cr.WithTx(ctx, tx, func(cr repos.ColumnRepo) error {
		c, err := models.NewColumn("New Column", 0)
		if err != nil {
			return err
		}
		_, err = cr.Create(ctx, b, c)
		return nil
	})
	if err != nil {
		err = fmt.Errorf("CreateBoard: %w", err)
		if rErr := tx.Rollback(); rErr != nil {
			err = fmt.Errorf("unable to rollback: %w", rErr)
		}
		return b, err
	}

	if err = tx.Commit(); err != nil {
		err = fmt.Errorf("unable to commit: %w", err)
	}

	return b, err
}

// func (srv BoardService) CreateBoard(u models.User, name, desc string) (b models.Board, err error) {
// 	b.Name = name
// 	b.Description = desc
// 	if err = b.Validate(); err != nil {
// 		return b, err
// 	}
// 	b, err = srv.br.Create(context.TODO(), u, b)
// 	return b, err
// }

// func (srv BoardService) CreateBoard(u models.User, name, desc string) (b models.Board, err error) {
// 	t := Tx{}

// 	b, err = srv.CreateBoard(u, name, desc)
// 	if err != nil {
// 		if err := t.Rollback(); err != nil {
// 			return b, err
// 		}
// 		return b, err
// 	}
// 	colSrv := NewColumnService(srv.cr)
// 	_, err = colSrv.CreateColumn(b, "New Column")

// 	if err != nil {
// 		if err := t.Rollback(); err != nil {
// 			return b, err
// 		}
// 	} else {
// 		if err := t.Commmit(); err != nil {
// 			return b, err
// 		}
// 	}
// 	return b, err
// }
