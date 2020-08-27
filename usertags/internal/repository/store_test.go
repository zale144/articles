// +build integration

package repository

import (
	"articles/usertags/internal/config"
	"articles/usertags/internal/model"
	"articles/usertags/internal/pkg/db"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestMain(m *testing.M) {

	var err error
	if err := config.Configure(); err != nil {
		log.Fatal(err)
	}

	dbc, err = db.Setup()
	if err != nil {
		log.Fatal(err)
	}

	m.Run()

	dbc.DropTableIfExists(&model.UserTags{}, &model.User{}, &model.Tag{})
}

var dbc *gorm.DB

func TestStore_AddTagsToUser(t *testing.T) {

	user := &model.User{
		Name:     "user1",
		Email:    "user1@test.com",
		Password: "asdf",
	}

	dbc.Create(user)

	user.Tags = []model.Tag{{Keyword: "tag1"}, {Keyword: "tag2"}}

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		user model.User
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
				db: dbc,
			},
			args: args{
				user: *user,
			},
			wantErr: false,
		}, {
			name: "add dupes",
			fields: fields{
				db: dbc,
			},
			args: args{
				user: func() model.User {
					u := *user
					u.Tags = append(u.Tags, model.Tag{
						Keyword: "tag3",
					})
					return u
				}(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Store{
				db: tt.fields.db,
			}

			err := u.AddTagsToUser(&tt.args.user)
			if !tt.wantErr {
				require.Nil(t, err, "error executing AddTagsToUser()")
			}

			exists := &model.User{}

			err = dbc.Model(user).Preload("Tags").Where(tt.args.user).First(exists).Error
			if !tt.wantErr {
				require.Nil(t, err, "error retrieving user")
				require.NotNil(t, exists)

				assert.Equal(t, exists.Name, tt.args.user.Name, "the saved user's Name does not match the retrieved one")
				assert.Equal(t, exists.Password, tt.args.user.Password, "the saved user's Password does not match the retrieved one")
				assert.Equal(t, exists.Email, tt.args.user.Email, "the saved user's Email does not match the retrieved one")
				assert.Equal(t, len(exists.Tags), len(tt.args.user.Tags), "the saved user's Tags do not match the retrieved one")
			}
		})
	}
}

func TestStore_CreateUser(t *testing.T) {

	user := model.User{
		Name:     "user2",
		Email:    "user2@test.com",
		Password: "asdf",
	}

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		user model.User
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
				db: dbc,
			},
			args: args{
				user: user,
			},
			wantErr: false,
		}, {
			name: "duplicate",
			fields: fields{
				db: dbc,
			},
			args: args{
				user: user,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Store{
				db: tt.fields.db,
			}

			err := u.CreateUser(&tt.args.user)
			if !tt.wantErr {
				require.Nil(t, err, "error executing CreateUser()")
			}

			exists := &model.User{}

			err = dbc.Model(user).Where(tt.args.user).First(&exists).Error
			if !tt.wantErr {
				require.Nil(t, err, "error retrieving user")
				require.NotNil(t, exists)

				assert.Equal(t, exists.Name, tt.args.user.Name, "the saved user's Name does not match the retrieved one")
				assert.Equal(t, exists.Password, tt.args.user.Password, "the saved user's Password does not match the retrieved one")
				assert.Equal(t, exists.Email, tt.args.user.Email, "the saved user's Email does not match the retrieved one")
				assert.Equal(t, len(exists.Tags), len(tt.args.user.Tags), "the saved user's Tags do not match the retrieved one")
			}
		})
	}
}

func TestStore_GetUser(t *testing.T) {

	user := &model.User{
		Name:     "user3",
		Email:    "user3@test.com",
		Password: "asdf",
	}

	dbc.Create(user)

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		email    string
		withTags bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "success: no tags",
			fields: fields{
				db: dbc,
			},
			args: args{
				email:    user.Email,
				withTags: false,
			},
			want:    user,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Store{
				db: tt.fields.db,
			}

			got, err := u.GetUser(tt.args.email, tt.args.withTags)
			if !tt.wantErr {
				require.Nil(t, err, "error executing GetUser()")
			}
			require.NotNil(t, got, "user must not be nil")

			assert.Equal(t, got.Name, user.Name, "the saved user's Name does not match the retrieved one")
			assert.Equal(t, got.Password, user.Password, "the saved user's Password does not match the retrieved one")
			assert.Equal(t, got.Email, user.Email, "the saved user's Email does not match the retrieved one")
			assert.Equal(t, got.Tags, user.Tags, "the saved user's Tags do not match the retrieved one")
		})
	}
}
