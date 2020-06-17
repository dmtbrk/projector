package persistence

import (
	"context"
	"database/sql"
	"log"

	"github.com/ortymid/projector/projector"
)

type SQLiteUserRepo struct {
	db *sql.DB
}

func NewSQLiteUserRepo(db *sql.DB) SQLiteUserRepo {
	return SQLiteUserRepo{db: db}
}

func (repo SQLiteUserRepo) All(ctx context.Context) ([]projector.User, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT name FROM users")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	users := []projector.User{}

	for rows.Next() {
		var name string

		err := rows.Scan(&name)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		// uid, err := uuid.Parse(rawUUID)
		// if err != nil {
		// 	log.Println(err)
		// 	return nil, err
		// }

		user, err := projector.NewUser(name)
		// err = u.Validate()
		if err != nil {
			log.Println(err)
			return nil, err
		}

		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return users, nil
}

func (repo SQLiteUserRepo) Create(ctx context.Context, user projector.User) (err error) {
	// if user.ID != 0 {
	// 	_, err = repo.db.ExecContext(ctx, "INSERT INTO users (id, name) VALUES (?, ?)", user.ID, user.Name)
	// } else {
	// 	_, err = repo.db.ExecContext(ctx, "INSERT INTO users (name) VALUES (?)", user.Name)
	// }
	_, err = repo.db.ExecContext(ctx, "INSERT INTO users (name) VALUES (?)", user.Name)
	if err != nil {
		log.Println("SQLiteUserRepo.Create:", err)
	}
	return err
}
