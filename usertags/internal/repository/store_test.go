// +build integration

package repository

import (
	"articles/usertags/internal/config"
	"articles/usertags/internal/model"
	"articles/usertags/internal/pkg/db"
	"github.com/jinzhu/gorm"
	"log"
	"reflect"
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

			if err := u.AddTagsToUser(&tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("AddTagsToUser() error = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr {
				var exists []model.User

				dbc.Model(user).Preload("Tags").Where(tt.args.user).Find(&exists)

				if l := len(exists); l == 0 {
					t.Errorf("Find() len() = %v, wantLen %v", l, 1)
				} else {
					if reflect.DeepEqual(exists[0].Tags, tt.args.user.Tags) {
						t.Errorf("Find() got = %v, want %v", exists[0].Tags, tt.args.user.Tags)
					}
				}
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
			if err := u.CreateUser(&tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr {
				var exists []model.User

				dbc.Model(user).Where(tt.args.user).Find(&exists)

				if l := len(exists); l == 0 {
					t.Errorf("Find() len() = %v, wantLen %v", l, 1)
				} else if !tt.wantErr {
					if exists[0].Email != tt.args.user.Email {
						t.Errorf("Find() got = %v, want %v", exists[0].Email, tt.args.user.Email)
					}
				}
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
			want: user,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Store{
				db: tt.fields.db,
			}

			got, err := u.GetUser(tt.args.email, tt.args.withTags)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if !tt.wantErr && got != nil {

				if got.Email != tt.want.Email {
					t.Errorf("GetUser().Email got = %v, want %v", got.Email, tt.want.Email)
				}

				if !reflect.DeepEqual(got.Email, tt.want.Email) {
					t.Errorf("GetUser().Tags got = %v, want %v", got.Tags, tt.want.Tags)
				}
			}
		})
	}
}
