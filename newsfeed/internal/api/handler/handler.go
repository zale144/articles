package handler

import (
	"articles/newsfeed/internal/api/handler/cards"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	GetCards,
	GetCardsByTags,
	AddCard echo.HandlerFunc
}

func NewHandler(cSvc cards.CardService) Handler {
	return Handler{
		GetCards:       cards.Get(cSvc),
		GetCardsByTags: cards.GetByTags(cSvc),
		AddCard:        cards.Add(cSvc),
	}
}
