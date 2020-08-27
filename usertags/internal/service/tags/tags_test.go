package tags

import (
	"articles/usertags/internal/dto"
	"articles/usertags/internal/model"
	"errors"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
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
	u := &model.User{}
	if w {
		u.Tags = m.tags
	}
	return u, nil
}

func (m mockStore) AddTagsToUser(*model.User) error {
	if m.fail {
		return errors.New("")
	}
	return nil
}

func TestTagsService_Add(t *testing.T) {
	type fields struct {
		store store
	}
	type args struct {
		email string
		t     dto.AddTagsPayload
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
				email: "",
				t: dto.AddTagsPayload{
					Tags: nil,
				},
			},
			wantErr: false,
		}, {
			name: "fail store",
			fields: fields{
				store: mockStore{fail: true},
			},
			args: args{
				email: "",
				t: dto.AddTagsPayload{
					Tags: nil,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Tags{
				store: tt.fields.store,
			}

			err := l.Add(tt.args.email, tt.args.t)
			if !tt.wantErr {
				require.Nil(t, err, "failed to execute Add()")
			}
		})
	}
}

func TestTags_Get(t *testing.T) {
	type fields struct {
		store store
	}
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    dto.GetTagsPayload
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				store: mockStore{tags: []model.Tag{{Keyword: "tag1"}, {Keyword: "tag2"}}},
			},
			args: args{
				email: "user@test.com",
			},
			want: dto.GetTagsPayload{
				Tags: []string{"tag1", "tag2"},
			},
			wantErr: false,
		}, {
			name: "fail store",
			fields: fields{
				store: mockStore{fail: true},
			},
			args: args{
				email: "user@test.com",
			},
			want:    dto.GetTagsPayload{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Tags{
				store: tt.fields.store,
			}
			got, err := l.Get(tt.args.email)
			if !tt.wantErr {
				require.Nil(t, err, "failed to execute Get()")
			}

			assert.Equal(t, got, tt.want, "response did not match expected output")
		})
	}
}
