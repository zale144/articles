package handler

import (
	"articles/usertags/internal/api/handler/tags"
	"articles/usertags/internal/api/handler/users"
	"github.com/labstack/echo/v4"
)

func NewHandler(uSvc users.UserService, tSvc tags.TagService) Handler {
	return Handler{
		Register: users.Register(uSvc),
		Login: users.Login(uSvc),
		GetTag: tags.Get(tSvc),
		AddTag: tags.Add(tSvc),
	}
}

type Handler struct {
	Login,
	Register,
	GetTag,
	AddTag echo.HandlerFunc
}
