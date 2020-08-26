package cards

import (
	"articles/newsfeed/internal/dto"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
)

type CardService interface {
	GetByUser(email string) (dto.GetCardsPayload, error)
	GetByTags(tags []string) (dto.GetCardsPayload, error)
	Add(card dto.Card) error
}

func Get(svc CardService) echo.HandlerFunc {
	return func(c echo.Context) error {

		email := c.Get("email").(string)

		cards, err := svc.GetByUser(email)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ResponseMessage{
				Message: errors.Wrap(err, "could not get cards").Error(),
			})
		}

		return c.JSON(http.StatusOK, cards)
	}
}

func GetByTags(svc CardService) echo.HandlerFunc {
	return func(c echo.Context) error {

		i := new(dto.GetCardsRequest)

		if err := c.Bind(i); err != nil {
			return err
		}

		if err := c.Validate(i); err != nil {
			return c.JSON(http.StatusBadRequest, dto.ResponseMessage{
				Message: err.(validator.ValidationErrors).Error(),
			})
		}

		cards, err := svc.GetByTags(i.Tags)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ResponseMessage{
				Message: errors.Wrap(err, "could not get cards").Error(),
			})
		}

		return c.JSON(http.StatusOK, cards)
	}
}

func Add(svc CardService) echo.HandlerFunc {
	return func(c echo.Context) error {

		i := new(dto.Card)

		if err := c.Bind(i); err != nil {
			return err
		}

		if err := c.Validate(i); err != nil {
			return c.JSON(http.StatusBadRequest, dto.ResponseMessage{
				Message: err.(validator.ValidationErrors).Error(),
			})
		}

		if err := svc.Add(*i); err != nil {
			return c.JSON(http.StatusBadRequest, dto.ResponseMessage{
				Message: errors.Wrap(err, "could not add card").Error(),
			})
		}

		return c.JSON(http.StatusOK, dto.ResponseMessage{
			Message: "card added successfully",
		})
	}
}
