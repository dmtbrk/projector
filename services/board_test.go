package services_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ortymid/projector/mock_repos"
	"github.com/ortymid/projector/models"
	"github.com/ortymid/projector/persistence/repos"
	"github.com/ortymid/projector/services"
)

// func TestBoardService_CreateBoard(t *testing.T) {
// 	type brMock struct {
// 		expectCtx   interface{}
// 		expectUser  interface{}
// 		expectBoard interface{}
// 		returnBoard interface{}
// 		returnErr   interface{}
// 	}
// 	type fields struct {
// 		brMock brMock
// 	}
// 	type args struct {
// 		u    models.User
// 		name string
// 		desc string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantB   models.Board
// 		wantErr bool
// 	}{
// 		{
// 			name: "Should create board",
// 			fields: fields{
// 				brMock: brMock{
// 					expectCtx:   gomock.Any(),
// 					expectUser:  gomock.Eq(models.User{ID: 1, Name: "user1"}),
// 					expectBoard: gomock.Eq(models.Board{Name: "board1", Description: "desc1"}),
// 					returnBoard: models.Board{ID: 1, Name: "board1", Description: "desc1", UserID: 1},
// 					returnErr:   nil,
// 				},
// 			},
// 			args: args{
// 				u:    models.User{ID: 1, Name: "user1"},
// 				name: "board1",
// 				desc: "desc1",
// 			},
// 			wantB:   models.Board{ID: 1, Name: "board1", Description: "desc1", UserID: 1},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)

// 			br := mock_services.NewMockBoardRepo(ctrl)
// 			br.EXPECT().Create(tt.args.brMock.expectCtx, tt.args.brMock.expectUser, tt.args.brMock.expectBoard).Return(tt.args.brMock.returnBoard, tt.args.brMock.returnErr)

// 			srv := services.NewBoardService(br)
// 			gotC, err := srv.CreateBoard(tt.args.u, tt.args.name, tt.args.desc)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("CreateColumn() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(gotC, tt.wantB) {
// 				t.Errorf("CreateColumn() = %v, want %v", gotC, tt.wantB)
// 			}
// 		})
// 	}
// }

func TestBoardService_CreateBoard(t *testing.T) {
	type brMock struct {
		expectCtx   interface{}
		expectUser  interface{}
		expectBoard interface{}
		returnBoard interface{}
		returnErr   interface{}
	}
	type crMock struct {
		expectCtx    interface{}
		expectBoard  interface{}
		expectColumn interface{}
		returnColumn interface{}
		returnErr    interface{}
	}
	type args struct {
		br   brMock
		cr   crMock
		ctx  context.Context
		u    models.User
		name string
		desc string
	}
	tests := []struct {
		name    string
		args    args
		wantB   models.Board
		wantErr bool
	}{
		{
			name: "Should create Board with default Column",
			args: args{
				br: brMock{
					expectCtx:   gomock.Any(),
					expectUser:  gomock.Eq(models.User{ID: 1, Name: "user1"}),
					expectBoard: gomock.Eq(models.Board{Name: "board1", Description: "desc1"}),
					returnBoard: models.Board{ID: 1, Name: "board1", Description: "desc1", UserID: 1},
					returnErr:   nil,
				},
				cr: crMock{
					expectCtx:    gomock.Any(),
					expectBoard:  gomock.Eq(models.Board{ID: 1, Name: "board1", Description: "desc1", UserID: 1}),
					expectColumn: gomock.Eq(models.Column{Name: "New Column"}),
					returnColumn: models.Column{ID: 1, Name: "New Column", BoardID: 1},
					returnErr:    nil,
				},
				u:    models.User{ID: 1, Name: "user1"},
				name: "board1",
				desc: "desc1",
			},
			wantB:   models.Board{ID: 1, Name: "board1", Description: "desc1", UserID: 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			tx := mock_repos.NewMockTx(ctrl)
			tx.EXPECT().Commit().Return(nil)

			var err error
			br := mock_repos.NewMockBoardRepo(ctrl)
			br.EXPECT().WithTx(gomock.Any(), gomock.Eq(nil), gomock.Any()).Do(func(ctx context.Context, tx repos.Tx, f func(repos.BoardRepo) error) {
				err = f(br)
			}).Return(tx, err)
			br.EXPECT().Create(tt.args.br.expectCtx, tt.args.br.expectUser, tt.args.br.expectBoard).Return(tt.args.br.returnBoard, tt.args.br.returnErr)

			cr := mock_repos.NewMockColumnRepo(ctrl)
			cr.EXPECT().WithTx(gomock.Any(), gomock.Eq(tx), gomock.Any()).Do(func(ctx context.Context, tx repos.Tx, f func(repos.ColumnRepo) error) {
				err = f(cr)
			}).Return(tx, err)
			cr.EXPECT().Create(tt.args.cr.expectCtx, tt.args.cr.expectBoard, tt.args.cr.expectColumn).Return(tt.args.cr.returnColumn, tt.args.cr.returnErr)

			srv := services.NewBoardService(br, cr)
			gotC, err := srv.CreateBoard(tt.args.ctx, tt.args.u, tt.args.name, tt.args.desc)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateColumn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantB) {
				t.Errorf("CreateColumn() = %v, want %v", gotC, tt.wantB)
			}
		})
	}
}
