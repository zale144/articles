package cards

import (
	"articles/newsfeed/internal/dto"
	"articles/newsfeed/internal/model"
	"context"
	"errors"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"github.com/zale144/articles/pb"
	"google.golang.org/grpc"
	"testing"
	"time"
)

type mockStore struct {
	fail bool
}

func (m mockStore) GetCards(tags []string, matchAll bool) ([]model.Card, error) {
	if m.fail {
		return nil, errors.New("")
	}
	return fakeCards, nil
}

func (m mockStore) AddCard(card *model.Card) error {
	if m.fail {
		return errors.New("")
	}
	return nil
}

func TestCardsService_Add(t *testing.T) {
	type fields struct {
		store  store
		client pb.TagsServiceClient
	}
	type args struct {
		crd dto.Card
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
				store:  mockStore{},
				client: nil,
			},
			args: args{
				crd: dto.Card{
					Title:     "card1",
					Timestamp: time.Now(),
					Tags:      []string{"tag1", "tag2"},
				},
			},
			wantErr: false,
		}, {
			name: "fail store",
			fields: fields{
				store: mockStore{
					fail: true,
				},
				client: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CardsService{
				store:  tt.fields.store,
				client: tt.fields.client,
			}

			err := c.Add(tt.args.crd)
			if !tt.wantErr {
				require.Nil(t, err, "failed to add card", err)
			}
		})
	}
}

func TestCardsService_GetByTags(t *testing.T) {
	type fields struct {
		store  store
		client pb.TagsServiceClient
	}
	type args struct {
		tags []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    dto.GetCardsPayload
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				store: mockStore{},
			},
			args: args{
				tags: []string{"tag1", "tag2"},
			},
			want: dto.GetCardsPayload{
				Cards: cardsDto(),
			},
			wantErr: false,
		}, {
			name: "fail store",
			fields: fields{
				store: mockStore{
					fail: true,
				},
			},
			want:    dto.GetCardsPayload{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CardsService{
				store:  tt.fields.store,
				client: tt.fields.client,
			}

			got, err := c.GetByTags(tt.args.tags)
			if !tt.wantErr {
				require.Nil(t, err, "failed to execute GetByTags()")
			}

			assert.Equal(t, got, tt.want, "response does not match expected output")
		})
	}
}

type mockClient struct {
	fail bool
}

func (c mockClient) GetUserTags(ctx context.Context, in *pb.UserTagsReq, opts ...grpc.CallOption) (*pb.UserTagsRsp, error) {
	if c.fail {
		return nil, errors.New("")
	}
	return &pb.UserTagsRsp{
		Tags: nil,
	}, nil
}

func TestCardsService_GetByUser(t *testing.T) {
	type fields struct {
		store  store
		client pb.TagsServiceClient
	}
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    dto.GetCardsPayload
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				store:  mockStore{},
				client: mockClient{},
			},
			args: args{
				email: "user@test.com",
			},
			want: dto.GetCardsPayload{
				Cards: cardsDto(),
			},
			wantErr: false,
		}, {
			name: "fail store",
			fields: fields{
				store: mockStore{
					fail: true,
				},
				client: mockClient{},
			},
			args: args{
				email: "user@test.com",
			},
			want:    dto.GetCardsPayload{},
			wantErr: true,
		}, {
			name: "fail client",
			fields: fields{
				store: mockStore{},
				client: mockClient{
					fail: true,
				},
			},
			want:    dto.GetCardsPayload{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CardsService{
				store:  tt.fields.store,
				client: tt.fields.client,
			}

			got, err := c.GetByUser(tt.args.email)
			if !tt.wantErr {
				require.Nil(t, err, "failed to execute GetByUser()")
			}

			assert.Equal(t, got, tt.want, "response does not match expected output")
		})
	}
}

func cardsDto() (c []dto.Card) {
	for _, f := range fakeCards {
		c = append(c, dto.Card(f))
	}
	return
}

var fakeCards = []model.Card{
	{
		Title:     "card1",
		Timestamp: time.Now().Local(),
		Tags:      []string{"tag1", "tag2"},
	}, {
		Title:     "card2",
		Timestamp: time.Now().Local(),
		Tags:      []string{"tag3", "tag3"},
	},
}
