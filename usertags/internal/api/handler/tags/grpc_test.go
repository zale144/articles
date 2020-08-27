package tags

import (
	"articles/usertags/internal/dto"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zale144/articles/pb"
	"testing"
)

type mockUserTagsService struct {
	fail bool
	tags dto.GetTagsPayload
}

func (m mockUserTagsService) Get(string) (dto.GetTagsPayload, error) {
	if m.fail {
		return dto.GetTagsPayload{}, errors.New("")
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
					tags: dto.GetTagsPayload{
						Tags: []string{"tag1", "tag2"},
					},
				},
			},
			args: args{
				ctx: nil,
				in: &pb.UserTagsReq{
					Email: "user@test.com",
				},
			},
			wantRsp: &pb.UserTagsRsp{
				Tags: []string{"tag1", "tag2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Tags{
				tagSrvc: tt.fields.tagSrvc,
			}
			gotRsp, err := u.GetUserTags(tt.args.ctx, tt.args.in)
			if !tt.wantErr {
				require.Nil(t, err, "failed to execute GetUserTags()")
			}

			assert.Equal(t, gotRsp, tt.wantRsp, "response does not match expected output")
		})
	}
}
