package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ortymid/projector/persistence/postgres"
	"github.com/ortymid/projector/services"
)

func main() {
	db := setupDB()
	defer teardownDB(db)

	ur := postgres.NewUserRepo(db, "users")
	br := postgres.NewBoardRepo(db, "boards")
	cr := postgres.NewColumnRepo(db, "columns")

	us := services.NewUserService(ur)
	bs := services.NewBoardService(br, cr)
	// cs := services.NewColumnService(cr)

	u, err := us.CreateUser("user1")
	if err != nil {
		panic(err)
	}
	b, err := bs.CreateBoard(context.TODO(), u, "b1", "desc")
	if err != nil {
		panic(err)
	}
	// c, _ := cs.CreateColumn(b, "New Column")
	cols, err := cr.AllByBoard(context.TODO(), b)
	if err != nil {
		panic(err)
	}

	fmt.Println(u, b, cols)
}

func setupDB() *sql.DB {
	db, err := sql.Open("postgres", "sslmode=disable user=postgres")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
	CREATE TABLE users (
		id		SERIAL PRIMARY KEY,
		name 	VARCHAR NOT NULL UNIQUE CHECK (char_length(name) > 0)
	)`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`
	CREATE TABLE boards (
		id			SERIAL PRIMARY KEY,
		name		VARCHAR(500) NOT NULL CHECK (char_length(name) > 0),
		description	VARCHAR(1000),
		user_id 	INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE
	)`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`
	CREATE TABLE columns (
		id			SERIAL PRIMARY KEY,
		board_id 	INTEGER NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
		name		VARCHAR(255) NOT NULL UNIQUE CHECK (char_length(name) > 0),
		order_index INTEGER NOT NULL
	)`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`
	CREATE TABLE tasks (
		id			SERIAL PRIMARY KEY,
		name		VARCHAR(255) NOT NULL UNIQUE CHECK (char_length(name) > 0),
		description VARCHAR(5000),
		column_id 	INTEGER NOT NULL REFERENCES columns(id) ON DELETE CASCADE
	)`)
	if err != nil {
		panic(err)
	}

	return db
}

func teardownDB(db *sql.DB) {
	_, err := db.Exec("DROP TABLE tasks")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("DROP TABLE columns")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("DROP TABLE boards")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("DROP TABLE users")
	if err != nil {
		panic(err)
	}

	err = db.Close()
	if err != nil {
		panic(err)
	}
}
