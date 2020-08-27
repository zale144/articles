package cards

import (
	"articles/newsfeed/internal/dto"
	"articles/newsfeed/internal/model"
	"context"
	"errors"
	"github.com/zale144/articles/pb"
	"google.golang.org/grpc"
	"reflect"
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
				store:  mockStore{
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
			if err := c.Add(tt.args.crd); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
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
				store:  mockStore{},
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
				store:  mockStore{
					fail: true,
				},
			},
			want: dto.GetCardsPayload{},
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
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByTags() got = %v, want %v", got, tt.want)
			}
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
		Tags:                 nil,
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
				store:  mockStore{
					fail: true,
				},
				client: mockClient{},
			},
			args: args{
				email: "user@test.com",
			},
			want: dto.GetCardsPayload{},
			wantErr: true,
		}, {
			name: "fail client",
			fields: fields{
				store:  mockStore{},
				client: mockClient{
					fail: true,
				},
			},
			want: dto.GetCardsPayload{},
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
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByUser() got = %v, want %v", got, tt.want)
			}
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
