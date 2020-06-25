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

const defaultUsersTable = "users"

type UserRepo struct {
	db    *sql.DB
	tx    *sql.Tx
	table string
}

func NewUserRepo(db *sql.DB, table string) UserRepo {
	if table == "" {
		table = defaultUsersTable
	}
	return UserRepo{db: db, table: table}
}

func (repo UserRepo) WithTx(ctx context.Context, tx repos.Tx, f func(repos.UserRepo) error) (repos.Tx, error) {
	var sqlTx *sql.Tx
	var ur repos.UserRepo
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
	ur = UserRepo{tx: sqlTx, table: repo.table}
	err := f(ur)
	return tx, err
}

func (repo UserRepo) All(ctx context.Context) (users []models.User, err error) {
	query := fmt.Sprintf("SELECT id, name FROM %s", repo.table)
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		err = fmt.Errorf("postgres: UserRepo.All: %w", err)
		log.Println(err)
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string

		err = rows.Scan(&id, &name)
		if err != nil {
			err = fmt.Errorf("postgres: UserRepo.All: %w", err)
			log.Println(err)
			return users, err
		}

		user := models.User{
			ID:   id,
			Name: name,
		}
		err = user.Validate()
		if err != nil {
			err = fmt.Errorf("postgres: UserRepo.All: %w", err)
			log.Println(err)
			return users, err
		}
		user.ID = id

		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		err = fmt.Errorf("postgres: UserRepo.All: %w", err)
		log.Println(err)
		return users, err
	}

	return users, err
}

func (repo UserRepo) Create(ctx context.Context, user models.User) (models.User, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id", repo.table)
	var row *sql.Row
	if repo.tx != nil {
		row = repo.tx.QueryRowContext(ctx, query, user.Name)
	} else {
		row = repo.db.QueryRowContext(ctx, query, user.Name)
	}
	err := row.Scan(&id)
	if err != nil {
		err = fmt.Errorf("postgres.UserRepo.Create: %w", err)
		log.Println(err)
		return user, err
	}
	user.ID = id
	return user, nil
}
