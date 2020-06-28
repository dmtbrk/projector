package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/ortymid/projector/models"
	"github.com/ortymid/projector/persistence/repos"
)

type ColumnRepo struct {
	db    *sql.DB
	tx    *sql.Tx
	table string
}

func NewColumnRepo(db *sql.DB, table string) ColumnRepo {
	return ColumnRepo{db: db, table: table}
}

func (repo ColumnRepo) WithTx(ctx context.Context, tx repos.Tx, f func(repos.ColumnRepo) error) (repos.Tx, error) {
	var sqlTx *sql.Tx
	var cr repos.ColumnRepo
	if tx == nil {
		sqlTx, err := repo.db.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
		}
		tx = sqlTx
	}
	sqlTx, ok := tx.(*sql.Tx)
	if !ok {
		return nil, errors.New("WithTx: wrong concrete tx type, expecting *sql.Tx")
	}
	cr = ColumnRepo{tx: sqlTx, table: repo.table}
	err := f(cr)
	return tx, err
}

func (repo ColumnRepo) AllByBoard(ctx context.Context, board models.Board) (cols []models.Column, err error) {
	if _, ok := board.ID.(int); !ok {
		return cols, errors.New("postgres: BoardRepo.AllByBoard: wrong board.ID type, expecting integer")
	}
	query := fmt.Sprintf(`
		SELECT c.id, c.board_id, c.name, c.order_index
		FROM %s c
		WHERE c.board_id = %d
		ORDER BY c.order_index ASC
	`, repo.table, board.ID)
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		err = fmt.Errorf("postgres: BoardRepo.AllByboard: %w", err)
		log.Println(err)
		return cols, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, boardID, orderIndex int
		var name string

		err := rows.Scan(&id, &boardID, &name, &orderIndex)
		if err != nil {
			err = fmt.Errorf("postgres: BoardRepo.AllByboard: %w", err)
			log.Println(err)
			return cols, err
		}

		col := models.Column{
			ID:         id,
			BoardID:    boardID,
			Name:       name,
			OrderIndex: orderIndex,
		}
		err = col.Validate()
		if err != nil {
			err = fmt.Errorf("postgres: BoardRepo.AllByboard: %w", err)
			log.Println(err)
			return cols, err
		}

		cols = append(cols, col)
	}
	if err := rows.Err(); err != nil {
		err = fmt.Errorf("postgres: BoardRepo.AllByboard: %w", err)
		log.Println(err)
		return cols, err
	}

	return cols, err
}

func (repo ColumnRepo) Create(ctx context.Context, b models.Board, c models.Column) (models.Column, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (board_id, name, order_index) VALUES ($1, $2, $3) RETURNING id", repo.table)
	var row *sql.Row
	if repo.tx != nil {
		row = repo.tx.QueryRowContext(ctx, query, b.ID, c.Name, c.OrderIndex)
	} else {
		row = repo.db.QueryRowContext(ctx, query, b.ID, c.Name, c.OrderIndex)
	}
	err := row.Scan(&id)
	if err != nil {
		err = fmt.Errorf("postgres.ColumnRepo.Create: %w", err)
		log.Println(err)
		return c, err
	}
	c.ID = id
	c.BoardID = b.ID
	return c, nil
}
