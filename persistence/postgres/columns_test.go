package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/ortymid/projector/models"
)

var testTableColumns = "columns"

func createTableColumns(t *testing.T, db *sql.DB) {
	_, err := db.Exec(fmt.Sprintf(`
	CREATE TABLE %s (
		id			SERIAL PRIMARY KEY,
		name		VARCHAR(255) NOT NULL UNIQUE CHECK (char_length(name) > 0),
		board_id 	INTEGER NOT NULL REFERENCES %s(id) ON DELETE CASCADE,
		order_index INTEGER NOT NULL UNIQUE
	)`, testTableColumns, testTableBoards))
	if err != nil {
		t.Fatal(err)
	}
}

func dropTableColumns(t *testing.T, db *sql.DB) {
	_, err := db.Exec(fmt.Sprintf(`DROP TABLE %s`, testTableColumns))
	if err != nil {
		t.Fatal(err)
	}
}

func insertIntoColumns(t *testing.T, db *sql.DB, cols []models.Column) {
	query := fmt.Sprintf("INSERT INTO %s (id, board_id, name, order_index) VALUES ($1, $2, $3, $4)", testTableColumns)
	for _, col := range cols {
		_, err := db.Exec(query, col.ID, col.BoardID, col.Name, col.OrderIndex)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestColumnRepo_AllByBoard(t *testing.T) {
	db := setupDB(t)
	defer teardownDB(t, db)

	type data struct {
		users   []models.User
		boards  []models.Board
		columns []models.Column
	}
	type fields struct {
		db    *sql.DB
		table string
	}
	type args struct {
		ctx   context.Context
		board models.Board
	}
	tests := []struct {
		name     string
		data     data
		fields   fields
		args     args
		wantCols []models.Column
		wantErr  bool
	}{
		{
			name: "Should return already sorted Columns in order",
			data: data{
				users: []models.User{
					{ID: 1, Name: "user1"},
				},
				boards: []models.Board{
					{ID: 1, Name: "board1", Description: "desc", UserID: 1},
				},
				columns: []models.Column{
					{ID: 1, BoardID: 1, Name: "col1", OrderIndex: 1},
					{ID: 2, BoardID: 1, Name: "col2", OrderIndex: 2},
					{ID: 3, BoardID: 1, Name: "col3", OrderIndex: 3},
				},
			},
			fields: fields{db: db, table: testTableColumns},
			args: args{
				ctx:   context.TODO(),
				board: models.Board{ID: 1, Name: "board1", Description: "desc", UserID: 1},
			},
			wantCols: []models.Column{
				{ID: 1, BoardID: 1, Name: "col1", OrderIndex: 1},
				{ID: 2, BoardID: 1, Name: "col2", OrderIndex: 2},
				{ID: 3, BoardID: 1, Name: "col3", OrderIndex: 3},
			},
			wantErr: false,
		},
		{
			name: "Should return not sorted Columns in order",
			data: data{
				users: []models.User{
					{ID: 1, Name: "user1"},
				},
				boards: []models.Board{
					{ID: 1, Name: "board1", Description: "desc", UserID: 1},
				},
				columns: []models.Column{
					{ID: 1, BoardID: 1, Name: "col3", OrderIndex: 3},
					{ID: 2, BoardID: 1, Name: "col1", OrderIndex: 1},
					{ID: 3, BoardID: 1, Name: "col2", OrderIndex: 2},
				},
			},
			fields: fields{db: db, table: testTableColumns},
			args: args{
				ctx:   context.TODO(),
				board: models.Board{ID: 1, Name: "board1", Description: "desc", UserID: 1},
			},
			wantCols: []models.Column{
				{ID: 2, BoardID: 1, Name: "col1", OrderIndex: 1},
				{ID: 3, BoardID: 1, Name: "col2", OrderIndex: 2},
				{ID: 1, BoardID: 1, Name: "col3", OrderIndex: 3},
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

			createTableColumns(t, db)
			insertIntoColumns(t, db, tt.data.columns)
			defer dropTableColumns(t, db)

			repo := ColumnRepo{
				db:    tt.fields.db,
				table: tt.fields.table,
			}
			gotCols, err := repo.AllByBoard(tt.args.ctx, tt.args.board)
			if (err != nil) != tt.wantErr {
				t.Errorf("ColumnRepo.AllByBoard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCols, tt.wantCols) {
				t.Errorf("ColumnRepo.AllByBoard() = %v, want %v", gotCols, tt.wantCols)
			}
		})
	}
}

func TestColumnRepo_Create(t *testing.T) {
	db := setupDB(t)
	defer teardownDB(t, db)

	type data struct {
		users   []models.User
		boards  []models.Board
		columns []models.Column
	}
	type fields struct {
		db    *sql.DB
		table string
	}
	type args struct {
		ctx context.Context
		b   models.Board
		c   models.Column
	}
	tests := []struct {
		name    string
		data    data
		fields  fields
		args    args
		want    models.Column
		wantErr bool
	}{
		{
			name: "Should create column",
			data: data{
				users: []models.User{
					{ID: 1, Name: "user1"},
				},
				boards: []models.Board{
					{ID: 1, Name: "board1", Description: "desc1", UserID: 1},
				},
				columns: []models.Column{},
			},
			fields: fields{
				db:    db,
				table: testTableColumns,
			},
			args: args{
				ctx: context.TODO(),
				b:   models.Board{ID: 1, Name: "board1", Description: "desc1", UserID: 1},
				c:   models.Column{Name: "col1"},
			},
			want:    models.Column{ID: 1, Name: "col1", BoardID: 1},
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

			createTableColumns(t, db)
			insertIntoColumns(t, db, tt.data.columns)
			defer dropTableColumns(t, db)

			repo := ColumnRepo{
				db:    tt.fields.db,
				table: tt.fields.table,
			}
			got, err := repo.Create(tt.args.ctx, tt.args.b, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("ColumnRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ColumnRepo.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
