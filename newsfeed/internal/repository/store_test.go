// +build integration

package repository

import (
	"articles/newsfeed/internal/config"
	"articles/newsfeed/internal/model"
	"articles/newsfeed/internal/pkg/db"
	"context"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
	"time"
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

	dbc.Database(viper.GetString(config.DBName)).Collection("cards").DeleteMany(context.TODO(), options.Delete())

	m.Run()
}

var dbc *mongo.Client

func TestStore_AddCard(t *testing.T) {

	type fields struct {
		client *mongo.Client
	}
	type args struct {
		card model.Card
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
				client: dbc,
			},
			args: args{
				card: model.Card{
					Title:     "card1",
					Timestamp: time.Time{},
					Tags:      []string{"tag1", "tag2"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Store{
				client: tt.fields.client,
			}
			err := u.AddCard(&tt.args.card)
			if !tt.wantErr {
				require.Nil(t, err, "failed to execute AddCard()")
			}

			res := u.client.Database(viper.GetString(config.DBName)).Collection("cards").FindOne(context.TODO(), bson.D{}, options.FindOne())
			require.NotNil(t, res, "card not found")

			card := model.Card{}
			err = res.Decode(&card)
			require.Nil(t, err, "error decoding result from database")
			assert.Equal(t, tt.args.card, card, "the saved card does not match the retrieved one")
		})
	}
}

func TestStore_GetCards(t *testing.T) {

	card := &model.Card{
		Title:     "card3",
		Timestamp: time.Time{},
		Tags:      []string{"tag5", "tag6"},
	}

	if _, err := dbc.Database(viper.GetString(config.DBName)).Collection("cards").
		InsertOne(context.TODO(), card, options.InsertOne()); err != nil {
		t.Fatal(err)
	}

	type fields struct {
		client *mongo.Client
	}
	type args struct {
		tags     []string
		matchAll bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []model.Card
	}{
		{
			name: "success: match any",
			fields: fields{
				client: dbc,
			},
			args: args{
				tags:     []string{"tag5", "tag8"},
				matchAll: false,
			},
			want: []model.Card{
				{
					Title:     "card3",
					Timestamp: time.Time{},
					Tags:      []string{"tag5", "tag6"},
				},
			},
		}, {
			name: "fail: match any",
			fields: fields{
				client: dbc,
			},
			args: args{
				tags:     []string{"tag8", "tag7"},
				matchAll: false,
			},
			want: nil,
		}, {
			name: "success: match all",
			fields: fields{
				client: dbc,
			},
			args: args{
				tags:     []string{"tag5", "tag6"},
				matchAll: true,
			},
			want: []model.Card{
				{
					Title:     "card3",
					Timestamp: time.Time{},
					Tags:      []string{"tag5", "tag6"},
				},
			},
		}, {
			name: "fail: match all",
			fields: fields{
				client: dbc,
			},
			args: args{
				tags:     []string{"tag5", "tag7"},
				matchAll: true,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Store{
				client: tt.fields.client,
			}

			got, err := u.GetCards(tt.args.tags, tt.args.matchAll)
			require.Nil(t, err, "failed to execute GetCards()")

			assert.Equal(t, tt.want, got, "the retrieved cards do not match the expected ones")
		})
	}
}
