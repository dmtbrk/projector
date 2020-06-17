package services_test

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ortymid/projector/mock_services"
	"github.com/ortymid/projector/models"
	"github.com/ortymid/projector/services"
)

func TestCommentService_CreateComment(t *testing.T) {
	lowerBoundText := makeTestString(1)
	upperBoundText := makeTestString(500)

	type crMock struct {
		expectCtx     interface{}
		expectTask    interface{}
		expectComment interface{}
		returnComment interface{}
		returnErr     interface{}
		notExpect     bool
	}
	type fields struct {
		cr crMock
	}
	type args struct {
		t    models.Task
		text string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantC   models.Comment
		wantErr bool
	}{
		{
			name: "Should create Comment with lower bound Text",
			fields: fields{
				cr: crMock{
					expectCtx:     gomock.Any(),
					expectTask:    gomock.Eq(models.Task{ID: 1, Name: "task1", Description: "desc", ColumnID: 1}),
					expectComment: gomock.Eq(models.Comment{Text: lowerBoundText}),
					returnComment: models.Comment{ID: 1, Text: lowerBoundText, TaskID: 1},
					returnErr:     nil,
				},
			},
			args: args{
				t:    models.Task{ID: 1, Name: "task1", Description: "desc", ColumnID: 1},
				text: lowerBoundText,
			},
			wantC:   models.Comment{ID: 1, Text: lowerBoundText, TaskID: 1},
			wantErr: false,
		},
		{
			name: "Should create Comment with upper bound Text",
			fields: fields{
				cr: crMock{
					expectCtx:     gomock.Any(),
					expectTask:    gomock.Eq(models.Task{ID: 1, Name: "task1", Description: "desc", ColumnID: 1}),
					expectComment: gomock.Eq(models.Comment{Text: upperBoundText}),
					returnComment: models.Comment{ID: 1, Text: upperBoundText, TaskID: 1},
					returnErr:     nil,
				},
			},
			args: args{
				t:    models.Task{ID: 1, Name: "task1", Description: "desc", ColumnID: 1},
				text: upperBoundText,
			},
			wantC:   models.Comment{ID: 1, Text: upperBoundText, TaskID: 1},
			wantErr: false,
		},
		{
			name: "Should error for Comment with empty Text",
			fields: fields{
				cr: crMock{notExpect: true},
			},
			args: args{
				t:    models.Task{ID: 1, Name: "task1", Description: "desc", ColumnID: 1},
				text: "",
			},
			wantC:   models.Comment{Text: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			cr := mock_services.NewMockCommentRepo(ctrl)
			if !tt.fields.cr.notExpect {
				cr.EXPECT().Create(tt.fields.cr.expectCtx, tt.fields.cr.expectTask, tt.fields.cr.expectComment).Return(tt.fields.cr.returnComment, tt.fields.cr.returnErr)
			} else {
				cr.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).MaxTimes(0)
			}

			srv := services.NewCommentService(cr)
			gotC, err := srv.CreateComment(tt.args.t, tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("CreateComment() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}
