package services_test

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ortymid/projector/mock_services"
	"github.com/ortymid/projector/models"
	"github.com/ortymid/projector/services"
)

func TestUserService_CreateUser(t *testing.T) {
	lowerBoundName := makeTestString(1)
	upperBoundName := makeTestString(500)

	type urMock struct {
		expectCtx  interface{}
		expectUser interface{}
		returnUser interface{}
		returnErr  interface{}
		notExpect  bool
	}
	type fields struct {
		ur urMock
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantU   models.User
		wantErr bool
	}{
		{
			name: "Should create User with lower bound Name",
			fields: fields{
				ur: urMock{
					expectCtx:  gomock.Any(),
					expectUser: gomock.Eq(models.User{Name: lowerBoundName}),
					returnUser: models.User{ID: 1, Name: lowerBoundName},
					returnErr:  nil,
				},
			},
			args: args{
				name: lowerBoundName,
			},
			wantU:   models.User{ID: 1, Name: lowerBoundName},
			wantErr: false,
		},
		{
			name: "Should create User with upper bound Name",
			fields: fields{
				ur: urMock{
					expectCtx:  gomock.Any(),
					expectUser: gomock.Eq(models.User{Name: upperBoundName}),
					returnUser: models.User{ID: 1, Name: upperBoundName},
					returnErr:  nil,
				},
			},
			args: args{
				name: upperBoundName,
			},
			wantU:   models.User{ID: 1, Name: upperBoundName},
			wantErr: false,
		},
		{
			name: "Should return error for User with empty Name",
			fields: fields{
				ur: urMock{notExpect: true},
			},
			args: args{
				name: "",
			},
			wantU:   models.User{Name: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			ur := mock_services.NewMockUserRepo(ctrl)
			if !tt.fields.ur.notExpect {
				ur.EXPECT().Create(tt.fields.ur.expectCtx, tt.fields.ur.expectUser).Return(tt.fields.ur.returnUser, tt.fields.ur.returnErr)
			} else {
				ur.EXPECT().Create(gomock.Any(), gomock.Any()).MaxTimes(0)
			}
			srv := services.NewUserService(ur)
			gotU, err := srv.CreateUser(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotU, tt.wantU) {
				t.Errorf("CreateUser() = %v, want %v", gotU, tt.wantU)
			}
		})
	}
}
