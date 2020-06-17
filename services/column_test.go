package services_test

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ortymid/projector/mock_services"
	"github.com/ortymid/projector/models"
	"github.com/ortymid/projector/services"
)

func TestColumnService_CreateColumn(t *testing.T) {
	type crMock struct {
		expectCtx    interface{}
		expectBoard  interface{}
		expectColumn interface{}
		returnColumn interface{}
		returnErr    interface{}
	}
	type fields struct {
		cr crMock
	}
	type args struct {
		b    models.Board
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantC   models.Column
		wantErr bool
	}{
		{
			name: "Should create Column",
			fields: fields{
				cr: crMock{
					expectCtx:    gomock.Any(),
					expectBoard:  gomock.Eq(models.Board{ID: 1, Name: "board1", Description: "desc", UserID: 1}),
					expectColumn: gomock.Eq(models.Column{Name: "col1"}),
					returnColumn: models.Column{ID: 1, Name: "col1", BoardID: 1},
					returnErr:    nil,
				},
			},
			args: args{
				b:    models.Board{ID: 1, Name: "board1", Description: "desc", UserID: 1},
				name: "col1",
			},
			wantC:   models.Column{ID: 1, Name: "col1", BoardID: 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			cr := mock_services.NewMockColumnRepo(ctrl)
			cr.EXPECT().Create(tt.fields.cr.expectCtx, tt.fields.cr.expectBoard, tt.fields.cr.expectColumn).Return(tt.fields.cr.returnColumn, tt.fields.cr.returnErr)

			srv := services.NewColumnService(cr)
			gotC, err := srv.CreateColumn(tt.args.b, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateColumn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("CreateColumn() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

// var nilColumn = Column{}

// type mockColumnRepo struct {
// 	allByBoard    []Column
// 	allByBoardErr error
// 	create        Column
// 	createErr     error
// }

// func (repo mockColumnRepo) AllByBoard(ctx context.Context, b Board) ([]Column, error) {
// 	return repo.allByBoard, repo.allByBoardErr
// }

// func (repo mockColumnRepo) Create(ctx context.Context, b Board, c Column) (Column, error) {
// 	if repo.create == nilColumn {
// 		c.ID = 1
// 		c.BoardID = b.ID
// 		return c, repo.createErr
// 	}
// 	return repo.create, repo.createErr
// }

// func TestColumn_MoveTask(t *testing.T) {
// 	type fields struct {
// 		Name  string
// 		Tasks []*Task
// 	}
// 	type args struct {
// 		from int
// 		to   int
// 	}
// 	tests := []struct {
// 		name       string
// 		fields     fields
// 		args       args
// 		wantFields fields
// 		wantErr    bool
// 	}{
// 		{
// 			"One position down: from 0 to 1, len 2",
// 			fields{Name: "test", Tasks: []*Task{
// 				{Name: "test1", Description: ""},
// 				{Name: "test2", Description: ""},
// 			}},
// 			args{0, 1},
// 			fields{Name: "test", Tasks: []*Task{
// 				{Name: "test2", Description: ""},
// 				{Name: "test1", Description: ""},
// 			}},
// 			false,
// 		},
// 		{
// 			"One position up: from 1 to 0, len 2",
// 			fields{Name: "test", Tasks: []*Task{
// 				{Name: "test1", Description: ""},
// 				{Name: "test2", Description: ""},
// 			}},
// 			args{0, 1},
// 			fields{Name: "test", Tasks: []*Task{
// 				{Name: "test2", Description: ""},
// 				{Name: "test1", Description: ""},
// 			}},
// 			false,
// 		},
// 		{
// 			"Two positions down: from 0 to 2, len 3",
// 			fields{Name: "test", Tasks: []*Task{
// 				{Name: "test1", Description: ""},
// 				{Name: "test2", Description: ""},
// 				{Name: "test3", Description: ""},
// 			}},
// 			args{0, 2},
// 			fields{Name: "test", Tasks: []*Task{
// 				{Name: "test2", Description: ""},
// 				{Name: "test3", Description: ""},
// 				{Name: "test1", Description: ""},
// 			}},
// 			false,
// 		},
// 		{
// 			"Two positions up: from 2 to 0, len 3",
// 			fields{Name: "test", Tasks: []*Task{
// 				{Name: "test1", Description: ""},
// 				{Name: "test2", Description: ""},
// 				{Name: "test3", Description: ""},
// 			}},
// 			args{2, 0},
// 			fields{Name: "test", Tasks: []*Task{
// 				{Name: "test3", Description: ""},
// 				{Name: "test1", Description: ""},
// 				{Name: "test2", Description: ""},
// 			}},
// 			false,
// 		},
// 		{
// 			"To out of range",
// 			fields{Name: "test", Tasks: []*Task{
// 				{Name: "test1", Description: ""},
// 				{Name: "test2", Description: ""},
// 			}},
// 			args{0, 3},
// 			fields{Name: "test", Tasks: []*Task{
// 				{Name: "test1", Description: ""},
// 				{Name: "test2", Description: ""},
// 			}},
// 			true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &Column{
// 				Name:  tt.fields.Name,
// 				Tasks: tt.fields.Tasks,
// 			}
// 			err := c.MoveTask(tt.args.from, tt.args.to)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Column.MoveTask() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			want := &Column{
// 				Name:  tt.wantFields.Name,
// 				Tasks: tt.wantFields.Tasks,
// 			}
// 			if !reflect.DeepEqual(c, want) {
// 				t.Errorf("Column.MoveTask() = %v, want %v", c, want)
// 			}
// 		})
// 	}
// }
