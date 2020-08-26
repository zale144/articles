package tags

import (
	"articles/pb"
	"context"
	"errors"
	"reflect"
	"testing"
)

type mockUserTagsService struct {
	fail    bool
	tags []string
}

func (m mockUserTagsService) GetUserTags(string) ([]string, error) {
	if m.fail {
		return nil, errors.New("")
	}
	return m.tags, nil
}

func TestUser_GetUserTags(t *testing.T) {
	type fields struct {
		tagSrvc tagService
	}
	type args struct {
		ctx context.Context
		in  *pb.UserTagsReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRsp *pb.UserTagsRsp
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				tagSrvc: mockUserTagsService{
					tags: []string{"tag1", "tag2"},
				},
			},
			args: args{
				ctx: nil,
				in: &pb.UserTagsReq{
					Email: "user@test.com",
				},
			},
			wantRsp: &pb.UserTagsRsp{
				Tags:  []string{"tag1", "tag2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Tags{
				tagSrvc: tt.fields.tagSrvc,
			}
			gotRsp, err := u.GetUserTags(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRsp, tt.wantRsp) {
				t.Errorf("GetUserTags() gotRsp = %v, want %v", gotRsp, tt.wantRsp)
			}
		})
	}
}
