package tags

import (
	"articles/usertags/internal/dto"
	"context"
	"errors"
	"github.com/zale144/articles/pb"
)

type Tags struct {
	tagSrvc tagService
}

func NewTags(tagSrvc tagService) Tags {
	return Tags{
		tagSrvc: tagSrvc,
	}
}

type tagService interface {
	Get(string) (dto.GetTagsPayload, error)
}

func (u Tags) GetUserTags(ctx context.Context, in *pb.UserTagsReq) (rsp *pb.UserTagsRsp, err error) {

	rsp = &pb.UserTagsRsp{}

	if len(in.Email) == 0 {
		err = errors.New("no email provided")
		return
	}

	tagsRsp, err := u.tagSrvc.Get(in.Email)
	if err != nil {
		err = errors.New("failed to find user tags")
		return
	}

	var tags = make([]string, len(tagsRsp.Tags))

	for i := range tagsRsp.Tags {
		tags[i] = tagsRsp.Tags[i]
	}

	rsp.Tags = tags
	return
}
