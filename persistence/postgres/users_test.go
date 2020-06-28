package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/ortymid/projector/models"
)

var testTableUsers = "users"

func createTableUsers(t *testing.T, db *sql.DB) {
	_, err := db.Exec(fmt.Sprintf(`
		CREATE TABLE %s (
			id SERIAL PRIMARY KEY,
			name VARCHAR NOT NULL UNIQUE CHECK (char_length(name) > 0)
		)
	`, testTableUsers))
	if err != nil {
		t.Fatal(err)
	}
}

func dropTableUsers(t *testing.T, db *sql.DB) {
	_, err := db.Exec(fmt.Sprintf("DROP TABLE %s", testTableUsers))
	if err != nil {
		t.Fatal(err)
	}
}

func insertIntoUsers(t *testing.T, db *sql.DB, users []models.User) {
	query := fmt.Sprintf("INSERT INTO %s (id, name) VALUES ($1, $2)", testTableUsers)
	for _, user := range users {
		_, err := db.Exec(query, user.ID, user.Name)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestUserRepo_All(t *testing.T) {
	db := setupDB(t)
	defer teardownDB(t, db)

	var nilUsers []models.User

	type fields struct {
		db    *sql.DB
		table string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		data    []models.User
		want    []models.User
		wantErr bool
	}{
		{
			name:    "No users",
			fields:  fields{db: db, table: testTableUsers},
			args:    args{ctx: context.TODO()},
			data:    []models.User{},
			want:    nilUsers,
			wantErr: false,
		},
		{
			name:    "One user",
			fields:  fields{db: db, table: testTableUsers},
			args:    args{ctx: context.TODO()},
			data:    []models.User{{ID: 1, Name: "user1"}},
			want:    []models.User{{ID: 1, Name: "user1"}},
			wantErr: false,
		},
		{
			name:    "Two users",
			fields:  fields{db: db, table: testTableUsers},
			args:    args{ctx: context.TODO()},
			data:    []models.User{{ID: 1, Name: "user1"}, {ID: 2, Name: "user2"}},
			want:    []models.User{{ID: 1, Name: "user1"}, {ID: 2, Name: "user2"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createTableUsers(t, db)
			insertIntoUsers(t, db, tt.data)
			defer dropTableUsers(t, db)

			repo := UserRepo{
				db:    tt.fields.db,
				table: tt.fields.table,
			}
			got, err := repo.All(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepo.All() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepo.All() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepo_Create(t *testing.T) {
	db := setupDB(t)
	defer teardownDB(t, db)

	type fields struct {
		db    *sql.DB
		table string
	}
	type args struct {
		ctx context.Context
		u   models.User
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		data     []models.User
		wantUser models.User
		wantErr  bool
	}{
		{
			name:     "Success with not existing user",
			fields:   fields{db, testTableUsers},
			args:     args{ctx: context.TODO(), u: models.User{Name: "user1"}},
			data:     []models.User{},
			wantUser: models.User{ID: 1, Name: "user1"},
			wantErr:  false,
		},
		{
			name:   "Error with existing user",
			fields: fields{db, testTableUsers},
			args:   args{ctx: context.TODO(), u: models.User{Name: "user1"}},
			data: []models.User{
				{ID: 1, Name: "user1"},
			},
			wantUser: models.User{Name: "user1"},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createTableUsers(t, db)
			insertIntoUsers(t, db, tt.data)
			defer dropTableUsers(t, db)

			repo := UserRepo{
				db:    tt.fields.db,
				table: tt.fields.table,
			}
			gotUser, err := repo.Create(tt.args.ctx, tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("UserRepo.Create() = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}

// func TestUserRepo_UserByBoard(t *testing.T) {
// 	db := setupDB(t)
// 	defer teardownDB(t, db)

// 	type data struct {
// 		users  []models.User
// 		boards []models.Board
// 	}

// 	type fields struct {
// 		db         *sql.DB
// 		table      string
// 		usersTable string
// 	}
// 	type args struct {
// 		ctx   context.Context
// 		board models.Board
// 	}
// 	tests := []struct {
// 		name     string
// 		fields   fields
// 		args     args
// 		data     data
// 		wantUser models.User
// 		wantErr  bool
// 	}{
// 		{
// 			name:   "Success with one user and one board",
// 			fields: fields{db: db, table: testTableBoards, usersTable: testTableUsers},
// 			args:   args{ctx: context.TODO(), board: models.Board{ID: 1, Name: "board1", UserID: 1}},
// 			data: data{
// 				users: []models.User{
// 					{ID: 1, Name: "user1"},
// 				},
// 				boards: []models.Board{
// 					{ID: 1, Name: "board1", Description: "desc", UserID: 1},
// 				},
// 			},
// 			wantUser: models.User{ID: 1, Name: "user1"},
// 			wantErr:  false,
// 		},
// 		{
// 			name:   "Success with two users and one board",
// 			fields: fields{db: db, table: testTableBoards, usersTable: testTableUsers},
// 			args:   args{ctx: context.TODO(), board: models.Board{ID: 1, Name: "board1", UserID: 1}},
// 			data: data{
// 				users: []models.User{
// 					{ID: 1, Name: "user1"},
// 					{ID: 2, Name: "user2"},
// 				},
// 				boards: []models.Board{
// 					{ID: 1, Name: "board1", Description: "desc", UserID: 2},
// 				},
// 			},
// 			wantUser: models.User{ID: 2, Name: "user2"},
// 			wantErr:  false,
// 		},
// 		{
// 			name:   "Success with one user and two boards",
// 			fields: fields{db: db, table: testTableBoards, usersTable: testTableUsers},
// 			args:   args{ctx: context.TODO(), board: models.Board{ID: 2, Name: "board2", UserID: 1}},
// 			data: data{
// 				users: []models.User{
// 					{ID: 1, Name: "user1"},
// 				},
// 				boards: []models.Board{
// 					{ID: 1, Name: "board1", Description: "desc", UserID: 1},
// 					{ID: 2, Name: "board2", Description: "desc", UserID: 1},
// 				},
// 			},
// 			wantUser: models.User{ID: 1, Name: "user1"},
// 			wantErr:  false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			createTableUsers(t, db)
// 			insertIntoUsers(t, db, tt.data.users)
// 			defer dropTableUsers(t, db)
// 			createTableBoards(t, db)
// 			insertIntoBoards(t, db, tt.data.boards)
// 			defer dropTableBoards(t, db)

// 			repo := BoardRepo{
// 				db:         tt.fields.db,
// 				table:      tt.fields.table,
// 				usersTable: tt.fields.usersTable,
// 			}
// 			gotUser, err := repo.UserByBoard(tt.args.ctx, tt.args.board)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("BoardRepo.UserByBoard() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(gotUser, tt.wantUser) {
// 				t.Errorf("BoardRepo.UserByBoard() = %v, want %v", gotUser, tt.wantUser)
// 			}
// 		})
// 	}
// }
