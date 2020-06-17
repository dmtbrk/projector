package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/ortymid/projector/models"
)

const defaultUsersTable = "users"

type UserRepo struct {
	db    *sql.DB
	table string
}

func NewUserRepo(db *sql.DB, table string) UserRepo {
	if table == "" {
		table = defaultUsersTable
	}
	return UserRepo{db: db, table: table}
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
	err := repo.db.QueryRowContext(ctx, query, user.Name).Scan(&id)
	if err != nil {
		err = fmt.Errorf("postgres.UserRepo.Create: %w", err)
		log.Println(err)
		return user, err
	}
	user.ID = id
	return user, nil
}
