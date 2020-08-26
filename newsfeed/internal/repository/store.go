package repository

import (
	"articles/newsfeed/internal/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	client *mongo.Client
}

func NewStore(cln *mongo.Client) Store {
	return Store{
		client: cln,
	}
}

func (u Store) GetCards(tags []string, matchAll bool) ([]model.Card, error) {
	col := u.client.Database("news").Collection("cards")

	op := "$in"
	if matchAll{
		op = "$all"
	}

	tagsIf := make([]interface{}, len(tags))
	for i := range tags {
		tagsIf[i] = tags[i]
	}

	filter := bson.D{{
		"tags",
		bson.D{{
			op,
			bson.A(tagsIf),
		}},
	}}

	findOptions := options.Find()

	var cards []model.Card

	cur, err := col.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var elem model.Card
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		cards = append(cards, elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}

func (u Store) AddCard(card *model.Card) error {
	col := u.client.Database("news").Collection("cards")
	_, err := col.InsertOne(context.TODO(), card)
	return err
}
