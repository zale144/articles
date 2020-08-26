package cards

import (
	"articles/newsfeed/internal/dto"
	"articles/newsfeed/internal/model"
	"articles/pb"
	"context"
)

type CardsService struct {
	store  store
	client pb.TagsServiceClient
}

func NewCardsService(str store, cln pb.TagsServiceClient) CardsService {
	return CardsService{
		store: str,
		client: cln,
	}
}

type store interface {
	GetCards(tags []string, matchAll bool) ([]model.Card, error)
	AddCard(card *model.Card) error
}

func (c CardsService) GetByUser(email string) (dto.GetCardsPayload, error) {
	rsp := dto.GetCardsPayload{}

	tagRsp, err := c.client.GetUserTags(context.TODO(), &pb.UserTagsReq{
		Email: email,
	})
	if err != nil {
		return rsp, err
	}

	cards, err := c.store.GetCards(tagRsp.Tags, false)
	if err != nil {
		return rsp, err
	}

	rsp.Cards = make([]dto.Card, len(cards))

	for i := range cards {
		rsp.Cards[i] = dto.Card(cards[i])
	}

	return rsp, nil
}

func (c CardsService) GetByTags(tags []string) (dto.GetCardsPayload, error) {
	rsp := dto.GetCardsPayload{}

	cards, err := c.store.GetCards(tags, true)
	if err != nil {
		return rsp, err
	}

	rsp.Cards = make([]dto.Card, len(cards))

	for i := range cards {
		rsp.Cards[i] = dto.Card(cards[i])
	}

	return rsp, nil
}

func (c CardsService) Add(crd dto.Card) error {
	card := model.Card(crd)
	return c.store.AddCard(&card)
}
