package users

import (
	"articles/usertags/internal/dto"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
)

type UserService interface {
	Register(user dto.RegisterPayload) error
	Login(user dto.LoginPayload) (string, error)
}

func Register(svc UserService) echo.HandlerFunc {
	return func(c echo.Context) error {

		i := new(dto.RegisterPayload)

		if err := c.Bind(i); err != nil {
			return err
		}

		if err := c.Validate(i); err != nil {
			return c.JSON(http.StatusBadRequest, dto.ResponseMessage{
				Message: err.(validator.ValidationErrors).Error(),
			})
		}

		if err := svc.Register(*i); err != nil {
			return c.JSON(http.StatusBadRequest, dto.ResponseMessage{
				Message: errors.Wrap(err, "could not register user").Error(),
			})
		}

		return c.JSON(http.StatusOK, dto.ResponseMessage{
			Message: "registration successful",
		})
	}
}

func Login(svc UserService) echo.HandlerFunc {
	return func(c echo.Context) error {

		i := new(dto.LoginPayload)

		if err := c.Bind(i); err != nil {
			return err
		}

		if err := c.Validate(i); err != nil {
			return c.JSON(http.StatusBadRequest, dto.ResponseMessage{
				Message: err.(validator.ValidationErrors).Error(),
			})
		}

		token, err := svc.Login(*i)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, dto.ResponseMessage{
				Message: "could not login user",
			})
		}

		return c.JSON(http.StatusOK, dto.LoginResponse{
			Token: token,
		})
	}
}
