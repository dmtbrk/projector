package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/ortymid/projector/models"
)

type TaskRepo struct {
	db    *sql.DB
	table string
}

func NewTaskRepo(db *sql.DB, table string) TaskRepo {
	return TaskRepo{db: db, table: table}
}

func (repo TaskRepo) AllByColumn(ctx context.Context, col models.Column) (tasks []models.Task, err error) {
	if _, ok := col.ID.(int); !ok {
		return tasks, errors.New("postgres.TaskRepo.AllByBoard: wrong board.ID type, expecting integer")
	}
	query := fmt.Sprintf(`
		SELECT t.id, t.name, t.description, t.column_id
		FROM %s t
		WHERE t.column_id = %d
	`, repo.table, col.ID)
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		err = fmt.Errorf("postgres.TaskRepo.AllByboard: %w", err)
		log.Println(err)
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, columnID int
		var name, desc string

		err := rows.Scan(&id, &name, &desc, &columnID)
		if err != nil {
			err = fmt.Errorf("postgres.TaskRepo.AllByboard: %w", err)
			log.Println(err)
			return tasks, err
		}

		task := models.Task{
			ID:          id,
			Name:        name,
			Description: desc,
			ColumnID:    columnID,
		}
		err = task.Validate()
		if err != nil {
			err = fmt.Errorf("postgres.TaskRepo.AllByboard: %w", err)
			log.Println(err)
			return tasks, err
		}

		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		err = fmt.Errorf("postgres.TaskRepo.AllByboard: %w", err)
		log.Println(err)
		return tasks, err
	}

	return tasks, err
}

func (repo TaskRepo) AllByBoard(ctx context.Context, b models.Board) ([]models.Task, error) {
	return nil, nil
}

func (repo TaskRepo) Create(ctx context.Context, c models.Column, t models.Task) (models.Task, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, description, column_id) VALUES ($1, $2, $3) RETURNING id", repo.table)
	err := repo.db.QueryRowContext(ctx, query, t.Name, t.Description, c.ID).Scan(&id)
	if err != nil {
		err = fmt.Errorf("postgres.TaskRepo.Create: %w", err)
		log.Println(err)
		return t, err
	}
	t.ID = id
	t.ColumnID = c.ID
	return t, nil
}
