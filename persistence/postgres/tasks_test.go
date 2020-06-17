package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/ortymid/projector/models"
)

var testTableTasks = "tasks"

func createTableTasks(t *testing.T, db *sql.DB) {
	_, err := db.Exec(fmt.Sprintf(`
	CREATE TABLE %s (
		id			SERIAL PRIMARY KEY,
		name		VARCHAR(255) NOT NULL UNIQUE CHECK (char_length(name) > 0),
		description VARCHAR(5000),
		column_id 	INTEGER NOT NULL REFERENCES %s(id) ON DELETE CASCADE
	)`, testTableTasks, testTableColumns))
	if err != nil {
		t.Fatal(err)
	}
}

func dropTableTasks(t *testing.T, db *sql.DB) {
	_, err := db.Exec(fmt.Sprintf(`DROP TABLE %s`, testTableTasks))
	if err != nil {
		t.Fatal(err)
	}
}

func insertIntoTasks(t *testing.T, db *sql.DB, tasks []models.Task) {
	query := fmt.Sprintf("INSERT INTO %s (id, name, desctiption, column_id) VALUES ($1, $2, $3, $4)", testTableTasks)
	for _, task := range tasks {
		_, err := db.Exec(query, task.ID, task.Name, task.Description, task.ColumnID)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestTaskRepo_Create(t *testing.T) {
	db := setupDB(t)
	defer teardownDB(t, db)

	type data struct {
		users   []models.User
		boards  []models.Board
		columns []models.Column
		tasks   []models.Task
	}
	type fields struct {
		db    *sql.DB
		table string
	}
	type args struct {
		ctx context.Context
		c   models.Column
		t   models.Task
	}
	tests := []struct {
		name    string
		data    data
		fields  fields
		args    args
		want    models.Task
		wantErr bool
	}{
		{
			name: "Should create Task",
			data: data{
				users: []models.User{
					{ID: 1, Name: "user1"},
				},
				boards: []models.Board{
					{ID: 1, Name: "board1", Description: "desc1", UserID: 1},
				},
				columns: []models.Column{
					{ID: 1, Name: "col1", BoardID: 1},
				},
				tasks: []models.Task{},
			},
			fields: fields{
				db:    db,
				table: testTableTasks,
			},
			args: args{
				ctx: context.TODO(),
				c:   models.Column{ID: 1, Name: "col1", BoardID: 1},
				t:   models.Task{Name: "task1", Description: "desc"},
			},
			want:    models.Task{ID: 1, Name: "task1", Description: "desc", ColumnID: 1},
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

			createTableTasks(t, db)
			insertIntoTasks(t, db, tt.data.tasks)
			defer dropTableTasks(t, db)

			repo := TaskRepo{
				db:    tt.fields.db,
				table: tt.fields.table,
			}
			got, err := repo.Create(tt.args.ctx, tt.args.c, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TaskRepo.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
