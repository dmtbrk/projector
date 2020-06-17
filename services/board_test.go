package services_test

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ortymid/projector/mock_services"
	"github.com/ortymid/projector/models"
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

func TestCreateBoardDefault(t *testing.T) {
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
		brMock brMock
		crMock crMock
		u      models.User
		name   string
		desc   string
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
				brMock: brMock{
					expectCtx:   gomock.Any(),
					expectUser:  gomock.Eq(models.User{ID: 1, Name: "user1"}),
					expectBoard: gomock.Eq(models.Board{Name: "board1", Description: "desc1"}),
					returnBoard: models.Board{ID: 1, Name: "board1", Description: "desc1", UserID: 1},
					returnErr:   nil,
				},
				crMock: crMock{
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

			br := mock_services.NewMockBoardRepo(ctrl)
			br.EXPECT().Create(tt.args.brMock.expectCtx, tt.args.brMock.expectUser, tt.args.brMock.expectBoard).Return(tt.args.brMock.returnBoard, tt.args.brMock.returnErr)

			cr := mock_services.NewMockColumnRepo(ctrl)
			cr.EXPECT().Create(tt.args.crMock.expectCtx, tt.args.crMock.expectBoard, tt.args.crMock.expectColumn).Return(tt.args.crMock.returnColumn, tt.args.crMock.returnErr)

			srv := services.NewBoardService(br, cr)
			gotC, err := srv.CreateBoardDefault(tt.args.u, tt.args.name, tt.args.desc)
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
