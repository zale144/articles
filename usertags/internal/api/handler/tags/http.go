package tags

import (
	"articles/usertags/internal/dto"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
)

type TagService interface {
	Add(email string, tags dto.AddTagsPayload) error
	Get(email string) (dto.GetTagsPayload, error)
}

func Add(svc TagService) echo.HandlerFunc {
	return func(c echo.Context) error {

		i := new(dto.AddTagsPayload)

		if err := c.Bind(i); err != nil {
			return err
		}

		if err := c.Validate(i); err != nil {
			return c.JSON(http.StatusBadRequest, dto.ResponseMessage{
				Message: err.(validator.ValidationErrors).Error(),
			})
		}

		email := c.Get("email").(string)

		if err := svc.Add(email, *i); err != nil {
			return c.JSON(http.StatusBadRequest, dto.ResponseMessage{
				Message: errors.Wrap(err, "could not add tag").Error(),
			})
		}

		return c.JSON(http.StatusOK, dto.ResponseMessage{
			Message: "tags added successfully",
		})
	}
}

func Get(svc TagService) echo.HandlerFunc {
	return func(c echo.Context) error {

		email := c.Get("email").(string)

		tags, err := svc.Get(email)
		if err != nil {
			return c.JSON(http.StatusBadRequest, dto.ResponseMessage{
				Message: errors.Wrap(err, "could not get tags").Error(),
			})
		}

		return c.JSON(http.StatusOK, tags)
	}
}
