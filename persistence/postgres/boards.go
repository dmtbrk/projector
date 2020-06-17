package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/ortymid/projector/models"
)

const defaultBoardsTable = "boards"

type BoardRepo struct {
	db    *sql.DB
	table string
}

func NewBoardRepo(db *sql.DB, table string) BoardRepo {
	if table == "" {
		table = defaultBoardsTable
	}
	return BoardRepo{db: db, table: table}
}

func (repo BoardRepo) AllByUser(ctx context.Context, user models.User) (boards []models.Board, err error) {
	if _, ok := user.ID.(int); !ok {
		return boards, errors.New("postgres: BoardRepo.AllByUser: wrong user.ID type, expecting integer")
	}
	query := fmt.Sprintf(`
		SELECT b.id, b.name, b.description, b.user_id
		FROM %s b
		WHERE b.user_id = %d
	`, repo.table, user.ID)
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		err = fmt.Errorf("postgres: BoardRepo.AllByUser: %w", err)
		log.Println(err)
		return boards, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, userID int
		var name, desc string

		err := rows.Scan(&id, &name, &desc, &userID)
		if err != nil {
			err = fmt.Errorf("postgres: BoardRepo.AllByUser: %w", err)
			log.Println(err)
			return boards, err
		}

		board := models.Board{
			ID:          id,
			Name:        name,
			Description: desc,
			UserID:      userID,
		}
		err = board.Validate()
		if err != nil {
			err = fmt.Errorf("postgres: BoardRepo.AllByUser: %w", err)
			log.Println(err)
			return boards, err
		}

		boards = append(boards, board)
	}
	if err := rows.Err(); err != nil {
		err = fmt.Errorf("postgres: BoardRepo.AllByUser: %w", err)
		log.Println(err)
		return boards, err
	}

	return boards, err
}

func (repo BoardRepo) Create(ctx context.Context, u models.User, b models.Board) (models.Board, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, description, user_id) VALUES ($1, $2, $3) RETURNING id", repo.table)
	err := repo.db.QueryRowContext(ctx, query, b.Name, b.Description, u.ID).Scan(&id)
	if err != nil {
		err = fmt.Errorf("postgres: BoardRepo.Create: %w", err)
		log.Println(err)
		return b, err
	}
	b.ID = id
	b.UserID = u.ID
	return b, nil
}
