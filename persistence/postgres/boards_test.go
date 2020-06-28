package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/ortymid/projector/models"
)

var testTableBoards = "boards"

func createTableBoards(t *testing.T, db *sql.DB) {
	_, err := db.Exec(fmt.Sprintf(`
	CREATE TABLE %s (
		id			SERIAL PRIMARY KEY,
		name		VARCHAR(500) NOT NULL CHECK (char_length(name) > 0),
		description	VARCHAR(1000),
		user_id 	INTEGER NOT NULL REFERENCES %s(id) ON DELETE CASCADE
	)`, testTableBoards, testTableUsers))
	if err != nil {
		t.Fatal(err)
	}
}

func dropTableBoards(t *testing.T, db *sql.DB) {
	_, err := db.Exec(fmt.Sprintf(`DROP TABLE %s`, testTableBoards))
	if err != nil {
		t.Fatal(err)
	}
}

func insertIntoBoards(t *testing.T, db *sql.DB, boards []models.Board) {
	query := fmt.Sprintf("INSERT INTO %s (id, name, description, user_id) VALUES ($1, $2, $3, $4)", testTableBoards)
	for _, board := range boards {
		_, err := db.Exec(query, board.ID, board.Name, board.Description, board.UserID)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestBoardRepo_AllByUser(t *testing.T) {
	db := setupDB(t)
	defer teardownDB(t, db)

	type data struct {
		users  []models.User
		boards []models.Board
	}

	type fields struct {
		db    *sql.DB
		table string
	}
	type args struct {
		ctx  context.Context
		user models.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		data    data
		want    []models.Board
		wantErr bool
	}{
		{
			name:    "No boards",
			fields:  fields{db: db, table: testTableBoards},
			args:    args{ctx: context.TODO(), user: models.User{ID: 1, Name: "user1"}},
			data:    data{users: nil, boards: nil},
			want:    nil,
			wantErr: false,
		},
		{
			name:   "One board",
			fields: fields{db: db, table: testTableBoards},
			args:   args{ctx: context.TODO(), user: models.User{ID: 1, Name: "user1"}},
			data: data{users: []models.User{
				{ID: 1, Name: "user1"},
			}, boards: []models.Board{
				{ID: 1, Name: "board1", Description: "desc", UserID: 1},
			}},
			want: []models.Board{
				{ID: 1, Name: "board1", Description: "desc", UserID: 1},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createTableUsers(t, db)
			insertIntoUsers(t, db, tt.data.users)
			defer dropTableUsers(t, db)

			createTableBoards(t, db)
			insertIntoBoards(t, db, tt.data.boards)
			defer dropTableBoards(t, db)

			repo := BoardRepo{
				db:    tt.fields.db,
				table: tt.fields.table,
			}
			got, err := repo.AllByUser(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("BoardRepo.AllByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoardRepo.AllByUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoardRepo_Create(t *testing.T) {
	db := setupDB(t)
	defer teardownDB(t, db)

	type data struct {
		users  []models.User
		boards []models.Board
	}
	type fields struct {
		db    *sql.DB
		table string
	}
	type args struct {
		ctx context.Context
		u   models.User
		b   models.Board
	}
	tests := []struct {
		name    string
		data    data
		fields  fields
		args    args
		want    models.Board
		wantErr bool
	}{
		{
			name: "Should create board",
			data: data{
				users: []models.User{
					{ID: 1, Name: "user1"},
				},
				boards: []models.Board{},
			},
			fields: fields{db: db, table: testTableBoards},
			args: args{
				ctx: context.TODO(),
				u:   models.User{ID: 1, Name: "user1"},
				b:   models.Board{Name: "board1", Description: "desc1"},
			},
			want:    models.Board{ID: 1, Name: "board1", Description: "desc1", UserID: 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createTableUsers(t, db)
			insertIntoUsers(t, db, tt.data.users)
			defer dropTableUsers(t, db)
			createTableBoards(t, db)
			insertIntoBoards(t, db, tt.data.boards)
			defer dropTableBoards(t, db)

			repo := BoardRepo{
				db:    tt.fields.db,
				table: tt.fields.table,
			}
			got, err := repo.Create(tt.args.ctx, tt.args.u, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("BoardRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoardRepo.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
