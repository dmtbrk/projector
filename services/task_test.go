package services_test

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ortymid/projector/mock_services"
	"github.com/ortymid/projector/models"
	"github.com/ortymid/projector/services"
)

func TestTaskService_CreateTask(t *testing.T) {
	lowerBoundName := makeTestString(1)
	upperBoundName := makeTestString(500)
	lowerBoundDesc := makeTestString(0)
	upperBoundDesc := makeTestString(5000)

	type trMock struct {
		expectCtx    interface{}
		expectColumn interface{}
		expectTask   interface{}
		returnTask   interface{}
		returnErr    interface{}
		notExpect    bool
	}
	type fields struct {
		tr trMock
	}
	type args struct {
		c    models.Column
		name string
		desc string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantT   models.Task
		wantErr bool
	}{
		{
			name: "Should create Task with Name and Description",
			fields: fields{
				tr: trMock{
					expectCtx:    gomock.Any(),
					expectColumn: gomock.Eq(models.Column{ID: 1, Name: "col1", BoardID: 1}),
					expectTask:   gomock.Eq(models.Task{Name: "task1", Description: "desc"}),
					returnTask:   models.Task{ID: 1, Name: "task1", Description: "desc", ColumnID: 1},
					returnErr:    nil,
				},
			},
			args: args{
				c:    models.Column{ID: 1, Name: "col1", BoardID: 1},
				name: "task1",
				desc: "desc",
			},
			wantT:   models.Task{ID: 1, Name: "task1", Description: "desc", ColumnID: 1},
			wantErr: false,
		},
		{
			name: "Should create Task with lower bound Name",
			fields: fields{
				tr: trMock{
					expectCtx:    gomock.Any(),
					expectColumn: gomock.Eq(models.Column{ID: 1, Name: "col1", BoardID: 1}),
					expectTask:   gomock.Eq(models.Task{Name: lowerBoundName, Description: "desc"}),
					returnTask:   models.Task{ID: 1, Name: lowerBoundName, Description: "desc", ColumnID: 1},
					returnErr:    nil,
				},
			},
			args: args{
				c:    models.Column{ID: 1, Name: "col1", BoardID: 1},
				name: lowerBoundName,
				desc: "desc",
			},
			wantT:   models.Task{ID: 1, Name: lowerBoundName, Description: "desc", ColumnID: 1},
			wantErr: false,
		},
		{
			name: "Should create Task with upper bound Name",
			fields: fields{
				tr: trMock{
					expectCtx:    gomock.Any(),
					expectColumn: gomock.Eq(models.Column{ID: 1, Name: "col1", BoardID: 1}),
					expectTask:   gomock.Eq(models.Task{Name: upperBoundName, Description: "desc"}),
					returnTask:   models.Task{ID: 1, Name: upperBoundName, Description: "desc", ColumnID: 1},
					returnErr:    nil,
				},
			},
			args: args{
				c:    models.Column{ID: 1, Name: "col1", BoardID: 1},
				name: upperBoundName,
				desc: "desc",
			},
			wantT:   models.Task{ID: 1, Name: upperBoundName, Description: "desc", ColumnID: 1},
			wantErr: false,
		},
		{
			name: "Should create Task with lower bound Description",
			fields: fields{
				tr: trMock{
					expectCtx:    gomock.Any(),
					expectColumn: gomock.Eq(models.Column{ID: 1, Name: "col1", BoardID: 1}),
					expectTask:   gomock.Eq(models.Task{Name: "task1", Description: lowerBoundDesc}),
					returnTask:   models.Task{ID: 1, Name: "task1", Description: lowerBoundDesc, ColumnID: 1},
					returnErr:    nil,
				},
			},
			args: args{
				c:    models.Column{ID: 1, Name: "col1", BoardID: 1},
				name: "task1",
				desc: lowerBoundDesc,
			},
			wantT:   models.Task{ID: 1, Name: "task1", Description: lowerBoundDesc, ColumnID: 1},
			wantErr: false,
		},
		{
			name: "Should create Task with upper bound Description",
			fields: fields{
				tr: trMock{
					expectCtx:    gomock.Any(),
					expectColumn: gomock.Eq(models.Column{ID: 1, Name: "col1", BoardID: 1}),
					expectTask:   gomock.Eq(models.Task{Name: "task1", Description: upperBoundDesc}),
					returnTask:   models.Task{ID: 1, Name: "task1", Description: upperBoundDesc, ColumnID: 1},
					returnErr:    nil,
				},
			},
			args: args{
				c:    models.Column{ID: 1, Name: "col1", BoardID: 1},
				name: "task1",
				desc: upperBoundDesc,
			},
			wantT:   models.Task{ID: 1, Name: "task1", Description: upperBoundDesc, ColumnID: 1},
			wantErr: false,
		},
		{
			name: "Should error for Task with empty name",
			fields: fields{
				tr: trMock{notExpect: true},
			},
			args: args{
				c:    models.Column{ID: 1, Name: "col1", BoardID: 1},
				name: "",
				desc: "desc",
			},
			wantT:   models.Task{Name: "", Description: "desc"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			tr := mock_services.NewMockTaskRepo(ctrl)
			if !tt.fields.tr.notExpect {
				tr.EXPECT().Create(tt.fields.tr.expectCtx, tt.fields.tr.expectColumn, tt.fields.tr.expectTask).Return(tt.fields.tr.returnTask, tt.fields.tr.returnErr)
			} else {
				tr.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).MaxTimes(0)
			}
			srv := services.NewTaskService(tr)
			gotT, err := srv.CreateTask(tt.args.c, tt.args.name, tt.args.desc)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotT, tt.wantT) {
				t.Errorf("CreateTask() = %v, want %v", gotT, tt.wantT)
			}
		})
	}
}
