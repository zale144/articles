package service

import (
	"articles/usertags/internal/dto"
	"articles/usertags/internal/model"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

type mockStore struct {
	fail bool
	tags []model.Tag
}

func (m mockStore) GetUser(s string, w bool) (*model.User, error) {
	if m.fail {
		return nil, errors.New("")
	}
	u := &model.User{
		Password: encryptPassword("secret"),
	}
	if w {
		u.Tags = m.tags
	}
	return u, nil
}

func encryptPassword(pass string) string {
	password, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	return string(password)
}

func (m mockStore) CreateUser(user *model.User) error {
	if m.fail {
		return errors.New("")
	}
	return nil
}

func TestUser_Login(t *testing.T) {
	type fields struct {
		store store
	}
	type args struct {
		userDto dto.LoginPayload
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				store: mockStore{},
			},
			args: args{
				userDto: dto.LoginPayload{
					Email:    "user@test.com",
					Password: "secret",
				},
			},
			wantErr: false,
		}, {
			name: "fail store",
			fields: fields{
				store: mockStore{fail: true},
			},
			args:    args{},
			wantErr: true,
		}, {
			name: "fail password",
			fields: fields{
				store: mockStore{},
			},
			args: args{
				userDto: dto.LoginPayload{
					Email:    "user@test.com",
					Password: "wrongpass",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := User{
				store: tt.fields.store,
			}
			got, err := r.Login(tt.args.userDto)
			if !tt.wantErr {
				require.Nil(t, err, "error while logging in", err)
				return
			}

			assert.Empty(t, got, "Token is empty")
		})
	}
}

func TestUser_Register(t *testing.T) {
	type fields struct {
		store store
	}
	type args struct {
		userDto dto.RegisterPayload
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				store: mockStore{},
			},
			args: args{
				userDto: dto.RegisterPayload{
					Name: "user1",
					LoginPayload: dto.LoginPayload{
						Email:    "user@test.com",
						Password: "secret",
					},
				},
			},
			wantErr: false,
		}, {
			name: "fail store",
			fields: fields{
				store: mockStore{fail: true},
			},
			args: args{
				userDto: dto.RegisterPayload{
					Name: "user1",
					LoginPayload: dto.LoginPayload{
						Email:    "user@test.com",
						Password: "secret",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := User{
				store: tt.fields.store,
			}
			err := r.Register(tt.args.userDto)
			if !tt.wantErr {
				require.Nil(t, err, "error while registering", err)
			}
		})
	}
}
