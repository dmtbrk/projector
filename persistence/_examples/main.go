package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ortymid/projector/persistence/postgres"
	"github.com/ortymid/projector/projector"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	type repoClose struct {
		name  string
		repo  projector.UserRepo
		close func() error
	}

	sqliteRepo, sqliteClose := createSQLiteRepo()
	mongoRepo, mongoClose := createMongoRepo()
	postgresRepo, postgresClose := createPostgresRepo()
	streamRepo, streamClose := createStreamRepo()

	repos := []repoClose{
		{"sqlite", sqliteRepo, sqliteClose},
		{"postgres", postgresRepo, postgresClose},
		{"mongo", mongoRepo, mongoClose},
		{"stream", streamRepo, streamClose},
	}

	for _, rc := range repos {

		ctx := context.Background()

		users, err := rc.repo.All(ctx)
		if err != nil {
			panic(err)
		}

		fmt.Println(rc.name, "before:", users)

		for i := 0; i < 10; i++ {
			user := projector.User{Name: "user" + strconv.Itoa(i)}
			if err != nil {
				panic(err)
			}
			// fmt.Println(user)

			err = rc.repo.Create(ctx, user)
			if err != nil {
				panic(err)
			}
		}

		users, err = rc.repo.All(ctx)
		if err != nil {
			panic(err)
		}

		fmt.Println(rc.name, "after:", users)

		rc.close()
	}
}

func createSQLiteRepo() (repo projector.UserRepo, close func() error) {
	os.Remove("test.sqlite")

	db, err := sql.Open("sqlite3", "test.sqlite")
	if err != nil {
		panic(err)
	}

	MustExec(db, `
		CREATE TABLE users (
			id INTEGER NOT NULL PRIMARY KEY,
			name TEXT NOT NULL
		)`)

	if err := db.Ping(); err != nil {
		panic(err)
	}

	repo = persistence.NewSQLiteUserRepo(db)

	return repo, db.Close
}

func createPostgresRepo() (repo projector.UserRepo, close func() error) {
	db, err := sql.Open("postgres", "sslmode=disable user=postgres")
	if err != nil {
		panic(err)
	}

	MustExec(db, `DROP TABLE IF EXISTS users`)
	MustExec(db, `
		CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			name VARCHAR NOT NULL
		)`)

	if err := db.Ping(); err != nil {
		panic(err)
	}

	repo = postgres.NewUserRepo(db)

	return repo, db.Close
}

func createMongoRepo() (repo projector.UserRepo, close func() error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	db := client.Database("test")
	err = db.Drop(ctx)
	if err != nil {
		panic(err)
	}
	cancel()

	repo = persistence.NewMongoUserRepo(db)
	close = func() error {
		return client.Disconnect(context.Background())
	}

	return repo, close
}

func createStreamRepo() (repo projector.UserRepo, close func() error) {
	file, err := os.Create("test.json")
	if err != nil {
		panic(err)
	}

	repo = persistence.NewStreamUserRepo(file)

	return repo, file.Close
}

func MustExec(db *sql.DB, query string, args ...interface{}) sql.Result {
	result, err := db.Exec(query, args...)
	if err != nil {
		panic(err)
	}
	return result
}
