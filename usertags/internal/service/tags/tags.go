package tags

import (
	"articles/usertags/internal/dto"
	"articles/usertags/internal/model"
)

type Tags struct {
	store store
}

func NewTagService(store store) Tags {
	return Tags{
		store: store,
	}
}

type store interface {
	GetUser(string, bool) (*model.User, error)
	AddTagsToUser(*model.User) error
}

func (l Tags) Add(email string, t dto.AddTagsPayload) error {

	var tags []model.Tag

	for _, t := range t.Tags {
		tags = append(tags, model.Tag{
			Keyword: t,
		})
	}

	user, err := l.store.GetUser(email, false)
	if err != nil {
		return err
	}

	user.Tags = append(user.Tags, tags...)

	return l.store.AddTagsToUser(user)
}

func (l Tags) Get(email string) (dto.GetTagsPayload, error) {

	rsp := dto.GetTagsPayload{}

	user, err := l.store.GetUser(email, true)
	if err != nil {
		return rsp, err
	}

	rsp.Tags = make([]string, len(user.Tags))

	for i := range user.Tags {
		rsp.Tags[i] = user.Tags[i].Keyword
	}

	return rsp, nil
}
